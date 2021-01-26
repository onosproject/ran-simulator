// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package endpoint

import (
	"io"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"google.golang.org/grpc/status"

	regapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("e2", "endpoint", "client")

// Client provides an E2 end-point client interface
type Client interface {
	io.Closer

	// Add adds a TerminationEndpoint
	Add(ctx context.Context, endPoint *regapi.TerminationEndpoint) error

	// Remove removes a TerminationEndpoint
	Remove(ctx context.Context, endPoint *regapi.TerminationEndpoint) error

	// Get returns a TerminationEndpoint based on a given TerminationEndpoint ID
	Get(ctx context.Context, id regapi.ID) (*regapi.TerminationEndpoint, error)

	// List returns the list of existing TerminationEndpoints
	List(ctx context.Context) ([]regapi.TerminationEndpoint, error)

	// Watch watches the TerminationEndpoint changes
	Watch(ctx context.Context, ch chan<- regapi.Event) error
}

// NewClient creates a new termination endpoint service client
func NewClient(conn *grpc.ClientConn) Client {
	cl := regapi.NewE2RegistryServiceClient(conn)
	return &endpointClient{
		client: cl,
	}
}

// endpointClient TerminationEndpoint client
type endpointClient struct {
	client regapi.E2RegistryServiceClient
}

// Add adds a new E2 termination end-point
func (c *endpointClient) Add(ctx context.Context, endPoint *regapi.TerminationEndpoint) error {
	req := &regapi.AddTerminationRequest{
		Endpoint: endPoint,
	}

	_, err := c.client.AddTermination(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return errors.FromStatus(stat)
		}
		return err
	}

	return nil

}

// Remove removes an E2 termination end-point
func (c *endpointClient) Remove(ctx context.Context, endPoint *regapi.TerminationEndpoint) error {
	req := &regapi.RemoveTerminationRequest{
		ID: endPoint.ID,
	}

	_, err := c.client.RemoveTermination(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return errors.FromStatus(stat)
		}
		return err
	}

	return nil
}

// Get returns information about an E2 termination end-point
func (c *endpointClient) Get(ctx context.Context, id regapi.ID) (*regapi.TerminationEndpoint, error) {
	req := &regapi.GetTerminationRequest{
		ID: id,
	}

	resp, err := c.client.GetTermination(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return nil, errors.FromStatus(stat)
		}
		return nil, err
	}

	return resp.Endpoint, nil
}

// List returns the list of currently registered E2 termination end-points
func (c *endpointClient) List(ctx context.Context) ([]regapi.TerminationEndpoint, error) {
	req := &regapi.ListTerminationsRequest{}

	resp, err := c.client.ListTerminations(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			return nil, errors.FromStatus(stat)
		}
		return nil, err
	}

	return resp.Endpoints, nil
}

// Watch watches for changes in the inventory of available E2T termination end-points
func (c *endpointClient) Watch(ctx context.Context, ch chan<- regapi.Event) error {
	req := regapi.WatchTerminationsRequest{}
	stream, err := c.client.WatchTerminations(ctx, &req)
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
				log.Error("An error occurred in receiving Endpoint changes", err)
			} else {
				ch <- resp.Event
			}
		}
	}()
	return nil
}

// Close closes the client connection
func (c *endpointClient) Close() error {
	return nil
}

var _ Client = &endpointClient{}
