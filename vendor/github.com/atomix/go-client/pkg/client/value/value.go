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

package value

import (
	"context"
	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/value"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"google.golang.org/grpc"
)

// Type is the value type
const Type primitive.Type = "Value"

// Client provides an API for creating Values
type Client interface {
	// GetValue gets the Value instance of the given name
	GetValue(ctx context.Context, name string) (Value, error)
}

// Value provides a simple atomic value
type Value interface {
	primitive.Primitive

	// Set sets the current value and returns the version
	Set(ctx context.Context, value []byte, opts ...SetOption) (uint64, error)

	// Get gets the current value and version
	Get(ctx context.Context) ([]byte, uint64, error)

	// Watch watches the value for changes
	Watch(ctx context.Context, ch chan<- *Event) error
}

// EventType is the type of a set event
type EventType string

const (
	// EventUpdated indicates the value was updated
	EventUpdated EventType = "updated"
)

// Event is a value change event
type Event struct {
	// Type is the change event type
	Type EventType

	// Value is the updated value
	Value []byte

	// Version is the updated version
	Version uint64
}

// New creates a new Lock primitive for the given partitions
// The value will be created in one of the given partitions.
func New(ctx context.Context, name primitive.Name, partitions []*primitive.Session) (Value, error) {
	i, err := util.GetPartitionIndex(name.Name, len(partitions))
	if err != nil {
		return nil, err
	}
	return newValue(ctx, name, partitions[i])
}

// newValue creates a new Value primitive for the given partition
func newValue(ctx context.Context, name primitive.Name, session *primitive.Session) (*value, error) {
	instance, err := primitive.NewInstance(ctx, name, session, &primitiveHandler{})
	if err != nil {
		return nil, err
	}
	return &value{
		name:     name,
		instance: instance,
	}, nil
}

// value is the single partition implementation of Lock
type value struct {
	name     primitive.Name
	instance *primitive.Instance
}

func (v *value) Name() primitive.Name {
	return v.name
}

func (v *value) Set(ctx context.Context, value []byte, opts ...SetOption) (uint64, error) {
	request := &api.SetRequest{}
	for i := range opts {
		opts[i].beforeSet(request)
	}

	r, err := v.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewValueServiceClient(conn)
		request := &api.SetRequest{
			Header: header,
			Value:  value,
		}
		for i := range opts {
			opts[i].beforeSet(request)
		}
		response, err := client.Set(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterSet(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}

	response := r.(*api.SetResponse)
	return response.Version, nil
}

func (v *value) Get(ctx context.Context) ([]byte, uint64, error) {
	r, err := v.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewValueServiceClient(conn)
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
		return nil, 0, err
	}

	response := r.(*api.GetResponse)
	return response.Value, response.Version, nil
}

func (v *value) Watch(ctx context.Context, ch chan<- *Event) error {
	stream, err := v.instance.DoCommandStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewValueServiceClient(conn)
		request := &api.EventRequest{
			Header: header,
		}
		return client.Events(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.ValueService_EventsClient).Recv()
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range stream {
			response := event.(*api.EventResponse)
			ch <- &Event{
				Type:    EventUpdated,
				Value:   response.NewValue,
				Version: response.NewVersion,
			}
		}
	}()
	return nil
}

func (v *value) Close(ctx context.Context) error {
	return v.instance.Close(ctx)
}

func (v *value) Delete(ctx context.Context) error {
	return v.instance.Delete(ctx)
}
