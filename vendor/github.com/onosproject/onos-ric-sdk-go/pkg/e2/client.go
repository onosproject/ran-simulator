// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2

import (
	"context"
	"fmt"

	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/encoding"

	"github.com/google/uuid"
	epapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"
	subapi "github.com/onosproject/onos-api/go/onos/e2sub/subscription"
	subtaskapi "github.com/onosproject/onos-api/go/onos/e2sub/task"
	e2tapi "github.com/onosproject/onos-api/go/onos/e2t/e2"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric-sdk-go/pkg/app"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/connection"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/endpoint"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/indication"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/subscription"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/subscriptiontask"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/termination"
)

var log = logging.GetLogger("e2")

const defaultServicePort = 5150

// Config is an E2 client configuration
type Config struct {
	// AppID is the application identifier
	AppID app.ID
	// InstanceID is the application instance identifier
	InstanceID app.InstanceID
	// SubscriptionService is the subscription service configuration
	SubscriptionService ServiceConfig
}

// ServiceConfig is an E2 service configuration
type ServiceConfig struct {
	// Host is the service host
	Host string
	// Port is the service port
	Port int
}

// GetHost gets the service host
func (c ServiceConfig) GetHost() string {
	return c.Host
}

// GetPort gets the service port
func (c ServiceConfig) GetPort() int {
	if c.Port == 0 {
		return defaultServicePort
	}
	return c.Port
}

// Client is an E2 client
type Client interface {
	// Subscribe creates a subscription from the given SubscriptionDetails
	// The Subscribe method will block until the subscription is successfully registered.
	// The context.Context represents the lifecycle of this initial subscription process.
	// Once the subscription has been created and the method returns, indications will be written
	// to the given channel.
	// If the subscription is successful, a subscription.Context will be returned. The subscription
	// context can be used to cancel the subscription by calling Close() on the subscription.Context.
	Subscribe(ctx context.Context, details subapi.SubscriptionDetails, ch chan<- indication.Indication) (subscription.Context, error)
}

// NewClient creates a new E2 client
func NewClient(config Config) (Client, error) {
	uuid.SetNodeID([]byte(fmt.Sprintf("%s:%s", config.AppID, config.InstanceID)))
	conns := connection.NewManager()
	subConn, err := conns.Connect(fmt.Sprintf("%s:%d", config.SubscriptionService.GetHost(), config.SubscriptionService.GetPort()))
	if err != nil {
		return nil, err
	}
	return &e2Client{
		config:     config,
		epClient:   endpoint.NewClient(subConn),
		subClient:  subscription.NewClient(subConn),
		taskClient: subscriptiontask.NewClient(subConn),
		conns:      conns,
	}, nil
}

// e2Client is the default E2 client implementation
type e2Client struct {
	config     Config
	epClient   endpoint.Client
	subClient  subscription.Client
	taskClient subscriptiontask.Client
	conns      *connection.Manager
}

func (c *e2Client) Subscribe(ctx context.Context, details subapi.SubscriptionDetails, ch chan<- indication.Indication) (subscription.Context, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	sub := &subapi.Subscription{
		ID:      subapi.ID(id.String()),
		AppID:   subapi.AppID(c.config.AppID),
		Details: &details,
	}

	client := &subContext{
		e2Client: c,
		sub:      sub,
	}
	err = client.subscribe(ctx, ch)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// subContext is an implementation of subscription.Context for managing a subscription
// The subscription context is responsible for creating the subscription, monitoring the task
// service for changes, and initializing the subscription stream with the appropriate E2 termination
// when the subscription is assigned by the subscription service.
type subContext struct {
	*e2Client
	sub    *subapi.Subscription
	cancel context.CancelFunc
}

// subscribe activates the subscription context
// The given context.Context is the context within which the subscription must be created,
// not the lifetime of the subscription. If the subscription cannot be created within the
// given context.Context, it will be deleted.
// Once the subscription has been created, the client tracks assignment of the subscription
// to E2 terminations by watching the subscription task service.
func (c *subContext) subscribe(ctx context.Context, indCh chan<- indication.Indication) error {
	// Add the subscription to the subscription service
	err := c.subClient.Add(ctx, c.sub)
	if err != nil {
		return err
	}

	// Watch the subscription task service to determine assignment of the subscription to E2 terminations
	watchCh := make(chan subtaskapi.Event)
	watchCtx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	err = c.taskClient.Watch(watchCtx, watchCh, subscriptiontask.WithSubscriptionID(c.sub.ID))
	if err != nil {
		return err
	}

	// The subscription is considered activated, and task events are processed in a separate goroutine.
	go c.processTaskEvents(watchCh, indCh)
	return nil
}

// processTaskEvents processes changes to subscription tasks related to the subscription
// When a task associated with this subscription is created, connect to the associated E2 termination
// and open a stream for indications. When the task is reassigned to a new termination point, clean
// up the prior stream and open a new stream to the new E2 termination point.
func (c *subContext) processTaskEvents(eventCh <-chan subtaskapi.Event, indCh chan<- indication.Indication) {
	// After the context is closed and the associated Watch call is canceled, the eventCh will be closed.
	// The indications channel is closed to indicate the subscription has been cleaned up.
	defer close(indCh)

	var prevCancel context.CancelFunc
	var prevEndpoint epapi.ID
	for event := range eventCh {
		// Only interested in tasks related to this subscription
		if event.Task.SubscriptionID != c.sub.ID {
			continue
		}

		// If the stream is already open for the associated E2 endpoint, skip the event
		if event.Task.EndpointID == prevEndpoint {
			continue
		}

		// If the task was assigned to a new endpoint, close the prior stream and open a new one.
		// If the task was unassigned, close the prior stream and wait for a new event.
		if event.Type == subtaskapi.EventType_NONE || event.Type == subtaskapi.EventType_CREATED {
			if prevCancel != nil {
				prevCancel()
			}
			ctx, cancel := context.WithCancel(context.Background())
			go func(epID epapi.ID) {
				defer cancel()
				err := c.openStream(ctx, epID, indCh)
				if err != nil {
					log.Error(err)
				}
			}(event.Task.EndpointID)
			prevEndpoint = event.Task.EndpointID
			prevCancel = cancel
		} else if event.Type == subtaskapi.EventType_REMOVED {
			prevEndpoint = ""
			if prevCancel != nil {
				prevCancel()
				prevCancel = nil
			}
		}
	}
}

// openStream opens a new stream to the given endpoint
// The client will lookup the endpoint address via the E2T endpoint service. If a valid endpoint is found,
// the client will connect to the E2 termination point and initialize the subscription stream with a StreamRequest.
// The stream lifetime is controlled by the given context.Context. When the context is closed, the stream and any
// associated state will be closed and cleaned up.
func (c *subContext) openStream(ctx context.Context, epID epapi.ID, indCh chan<- indication.Indication) error {
	response, err := c.epClient.Get(ctx, epID)
	if err != nil {
		return err
	}

	conn, err := c.conns.Connect(fmt.Sprintf("%s:%d", response.IP, response.Port))
	if err != nil {
		return err
	}

	client := termination.NewClient(conn)
	responseCh := make(chan e2tapi.StreamResponse)
	requestCh, err := client.Stream(ctx, responseCh)
	if err != nil {
		return err
	}

	requestCh <- e2tapi.StreamRequest{
		AppID:          e2tapi.AppID(c.config.AppID),
		InstanceID:     e2tapi.InstanceID(c.config.InstanceID),
		SubscriptionID: e2tapi.SubscriptionID(c.sub.ID),
	}

	for response := range responseCh {
		indCh <- indication.Indication{
			EncodingType: encoding.Type(response.Header.EncodingType),
			Payload: indication.Payload{
				Header:  response.Header.IndicationHeader,
				Message: response.IndicationMessage,
			},
		}
	}
	return nil
}

// Close closes the subscription context
// When the subscription context is closed, any existing streams will be closed and the subscription will be
// removed from the subscription service. Propagation of the subscription delete is asynchronous.
func (c *subContext) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	return c.subClient.Remove(context.Background(), c.sub)
}

var _ subscription.Context = &subContext{}
