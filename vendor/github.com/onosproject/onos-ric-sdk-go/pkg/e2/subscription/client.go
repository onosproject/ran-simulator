// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"io"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"google.golang.org/grpc/status"

	subapi "github.com/onosproject/onos-api/go/onos/e2sub/subscription"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("e2", "subscription", "client")

// Context is a subscription context
type Context interface {
	io.Closer

	// ID returns the subscription identifier
	ID() subapi.ID

	// Err returns the subscription error channel
	Err() <-chan error
}

// Client is an E2 subscription service client interface
type Client interface {
	io.Closer

	// Add adds a subscription
	Add(ctx context.Context, subscription *subapi.Subscription) error

	// Remove removes a subscription
	Remove(ctx context.Context, subscription *subapi.Subscription) error

	// Get returns a subscription based on a given subscription ID
	Get(ctx context.Context, id subapi.ID) (*subapi.Subscription, error)

	// List returns the list of existing subscriptions
	List(ctx context.Context) ([]subapi.Subscription, error)

	// Watch watches the subscription changes
	Watch(ctx context.Context, ch chan<- subapi.Event) error
}

// NewClient creates a new subscribe service client
func NewClient(conn *grpc.ClientConn) Client {
	cl := subapi.NewE2SubscriptionServiceClient(conn)
	return &subscriptionClient{
		client: cl,
	}
}

// subscriptionClient subscription client
type subscriptionClient struct {
	client subapi.E2SubscriptionServiceClient
}

// Add adds a subscription
func (c *subscriptionClient) Add(ctx context.Context, subscription *subapi.Subscription) error {
	req := &subapi.AddSubscriptionRequest{
		Subscription: subscription,
	}

	_, err := c.client.AddSubscription(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return errors.FromStatus(stat)
		}
		return err
	}

	return nil

}

// Remove removes a subscription
func (c *subscriptionClient) Remove(ctx context.Context, subscription *subapi.Subscription) error {
	req := &subapi.RemoveSubscriptionRequest{
		ID: subscription.ID,
	}

	_, err := c.client.RemoveSubscription(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return errors.FromStatus(stat)
		}
		return err
	}

	return nil
}

// Get returns information about a subscription
func (c *subscriptionClient) Get(ctx context.Context, id subapi.ID) (*subapi.Subscription, error) {
	req := &subapi.GetSubscriptionRequest{
		ID: id,
	}

	resp, err := c.client.GetSubscription(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return nil, errors.FromStatus(stat)
		}
		return nil, err
	}

	return resp.Subscription, nil
}

// List returns the list of all subscriptions
func (c *subscriptionClient) List(ctx context.Context) ([]subapi.Subscription, error) {
	req := &subapi.ListSubscriptionsRequest{}

	resp, err := c.client.ListSubscriptions(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return nil, errors.FromStatus(stat)
		}
		return nil, err
	}

	return resp.Subscriptions, nil
}

// Watch watches for changes in the set of subscriptions
func (c *subscriptionClient) Watch(ctx context.Context, ch chan<- subapi.Event) error {
	req := subapi.WatchSubscriptionsRequest{}
	stream, err := c.client.WatchSubscriptions(ctx, &req)
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
				log.Error("An error occurred in receiving Subscription changes", err)
			} else {
				ch <- resp.Event
			}
		}
	}()
	return nil
}

// Close closes the client connection
func (c *subscriptionClient) Close() error {
	return nil
}

var _ Client = &subscriptionClient{}
