// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscriptiontask

import (
	"io"

	epapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"
	subapi "github.com/onosproject/onos-api/go/onos/e2sub/subscription"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"google.golang.org/grpc/status"

	subtaskapi "github.com/onosproject/onos-api/go/onos/e2sub/task"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("e2", "subscription", "client")

// ListOption is an option for filtering List calls
type ListOption interface {
	applyList(*listOptions)
}

type listOptions struct {
	subscriptionID subapi.ID
	endpointID     epapi.ID
}

// WatchOption is an option for filtering Watch calls
type WatchOption interface {
	applyWatch(*watchOptions)
}

type watchOptions struct {
	subscriptionID subapi.ID
	endpointID     epapi.ID
}

// FilterOption is an option for filtering List/Watch calls
type FilterOption interface {
	ListOption
	WatchOption
}

// WithSubscriptionID creates an option for filtering by subscription ID
func WithSubscriptionID(id subapi.ID) FilterOption {
	return &filterSubscriptionOption{
		subID: id,
	}
}

type filterSubscriptionOption struct {
	subID subapi.ID
}

func (o *filterSubscriptionOption) applyList(options *listOptions) {
	options.subscriptionID = o.subID
}

func (o *filterSubscriptionOption) applyWatch(options *watchOptions) {
	options.subscriptionID = o.subID
}

// WithEndpointID creates an option for filtering by endpoint ID
func WithEndpointID(id epapi.ID) FilterOption {
	return &filterEndpointOption{
		epID: id,
	}
}

type filterEndpointOption struct {
	epID epapi.ID
}

func (o *filterEndpointOption) applyList(options *listOptions) {
	options.endpointID = o.epID
}

func (o *filterEndpointOption) applyWatch(options *watchOptions) {
	options.endpointID = o.epID
}

// Client is an E2 subscription service client interface
type Client interface {
	io.Closer

	// Get returns a subscription based on a given subscription ID
	Get(ctx context.Context, id subtaskapi.ID) (*subtaskapi.SubscriptionTask, error)

	// List returns the list of existing subscriptions
	List(ctx context.Context, opts ...ListOption) ([]subtaskapi.SubscriptionTask, error)

	// Watch watches the subscription changes
	Watch(ctx context.Context, ch chan<- subtaskapi.Event, opts ...WatchOption) error
}

// NewClient creates a new subscribe task service client
func NewClient(conn *grpc.ClientConn) Client {
	cl := subtaskapi.NewE2SubscriptionTaskServiceClient(conn)
	return &subscriptionTaskClient{
		client: cl,
	}
}

// subscriptionTaskClient subscription client
type subscriptionTaskClient struct {
	client subtaskapi.E2SubscriptionTaskServiceClient
}

// Get returns information about a subscription
func (c *subscriptionTaskClient) Get(ctx context.Context, id subtaskapi.ID) (*subtaskapi.SubscriptionTask, error) {
	req := &subtaskapi.GetSubscriptionTaskRequest{
		ID: id,
	}

	resp, err := c.client.GetSubscriptionTask(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return nil, errors.FromStatus(stat)
		}
		return nil, err
	}

	return resp.Task, nil
}

// List returns the list of all subscriptions
func (c *subscriptionTaskClient) List(ctx context.Context, opts ...ListOption) ([]subtaskapi.SubscriptionTask, error) {
	options := &listOptions{}
	for _, opt := range opts {
		opt.applyList(options)
	}

	req := &subtaskapi.ListSubscriptionTasksRequest{
		SubscriptionID: options.subscriptionID,
		EndpointID:     options.endpointID,
	}

	resp, err := c.client.ListSubscriptionTasks(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return nil, errors.FromStatus(stat)
		}
		return nil, err
	}
	return resp.Tasks, nil
}

// Watch watches for changes in the set of subscriptions
func (c *subscriptionTaskClient) Watch(ctx context.Context, ch chan<- subtaskapi.Event, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.applyWatch(options)
	}

	req := subtaskapi.WatchSubscriptionTasksRequest{
		SubscriptionID: options.subscriptionID,
		EndpointID:     options.endpointID,
	}

	stream, err := c.client.WatchSubscriptionTasks(ctx, &req)
	if err != nil {
		defer close(ch)
		stat, ok := status.FromError(err)
		if ok {
			return errors.FromStatus(stat)
		}
		return err
	}

	go func() {
		defer close(ch)
		for {
			resp, err := stream.Recv()
			if err == io.EOF || err == context.Canceled {
				break
			}

			if err != nil {
				stat, ok := status.FromError(err)
				if ok {
					err = errors.FromStatus(stat)
					if errors.IsCanceled(err) || errors.IsTimeout(err) {
						break
					}
				}
				log.Error("An error occurred in receiving SubscriptionTask changes", err)
			} else {
				ch <- resp.Event
			}
		}
	}()
	return nil
}

// Close closes the client connection
func (c *subscriptionTaskClient) Close() error {
	return nil
}

var _ Client = &subscriptionTaskClient{}
