// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package counter

import (
	"context"
	api "github.com/atomix/api/proto/atomix/counter"
	"github.com/atomix/api/proto/atomix/headers"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"google.golang.org/grpc"
)

// Type is the counter type
const Type primitive.Type = "Counter"

// Client provides an API for creating Counters
type Client interface {
	// GetCounter gets the Counter instance of the given name
	GetCounter(ctx context.Context, name string) (Counter, error)
}

// Counter provides a distributed atomic counter
type Counter interface {
	primitive.Primitive

	// Get gets the current value of the counter
	Get(ctx context.Context) (int64, error)

	// Set sets the value of the counter
	Set(ctx context.Context, value int64) error

	// Increment increments the counter by the given delta
	Increment(ctx context.Context, delta int64) (int64, error)

	// Decrement decrements the counter by the given delta
	Decrement(ctx context.Context, delta int64) (int64, error)
}

// New creates a new counter for the given partitions
func New(ctx context.Context, name primitive.Name, partitions []*primitive.Session) (Counter, error) {
	i, err := util.GetPartitionIndex(name.Name, len(partitions))
	if err != nil {
		return nil, err
	}

	instance, err := primitive.NewInstance(ctx, name, partitions[i], &primitiveHandler{})
	if err != nil {
		return nil, err
	}

	return &counter{
		name:     name,
		instance: instance,
	}, nil
}

// counter is the single partition implementation of Counter
type counter struct {
	name     primitive.Name
	instance *primitive.Instance
}

func (c *counter) Name() primitive.Name {
	return c.name
}

func (c *counter) Get(ctx context.Context) (int64, error) {
	response, err := c.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewCounterServiceClient(conn)
		request := &api.GetRequest{
			Header: header,
		}
		response, err := client.Get(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}
	return response.(*api.GetResponse).Value, nil
}

func (c *counter) Set(ctx context.Context, value int64) error {
	_, err := c.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewCounterServiceClient(conn)
		request := &api.SetRequest{
			Header: header,
			Value:  value,
		}
		response, err := client.Set(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	return err
}

func (c *counter) Increment(ctx context.Context, delta int64) (int64, error) {
	response, err := c.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewCounterServiceClient(conn)
		request := &api.IncrementRequest{
			Header: header,
			Delta:  delta,
		}
		response, err := client.Increment(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}
	return response.(*api.IncrementResponse).NextValue, nil
}

func (c *counter) Decrement(ctx context.Context, delta int64) (int64, error) {
	response, err := c.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewCounterServiceClient(conn)
		request := &api.DecrementRequest{
			Header: header,
			Delta:  delta,
		}
		response, err := client.Decrement(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}
	return response.(*api.DecrementResponse).NextValue, nil
}

func (c *counter) Close(ctx context.Context) error {
	return c.instance.Close(ctx)
}

func (c *counter) Delete(ctx context.Context) error {
	return c.instance.Delete(ctx)
}
