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

package set

import (
	"context"
	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/set"
	"github.com/atomix/go-client/pkg/client/primitive"
	"google.golang.org/grpc"
)

func newPartition(ctx context.Context, name primitive.Name, session *primitive.Session) (Set, error) {
	sess, err := primitive.NewInstance(ctx, name, session, &primitiveHandler{})
	if err != nil {
		return nil, err
	}
	return &setPartition{
		name:     name,
		instance: sess,
	}, nil
}

type setPartition struct {
	name     primitive.Name
	instance *primitive.Instance
}

func (s *setPartition) Name() primitive.Name {
	return s.name
}

func (s *setPartition) Add(ctx context.Context, value string) (bool, error) {
	r, err := s.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewSetServiceClient(conn)
		request := &api.AddRequest{
			Header: header,
			Value:  value,
		}
		response, err := client.Add(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return false, err
	}

	response := r.(*api.AddResponse)
	return response.Added, nil
}

func (s *setPartition) Remove(ctx context.Context, value string) (bool, error) {
	r, err := s.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewSetServiceClient(conn)
		request := &api.RemoveRequest{
			Header: header,
			Value:  value,
		}
		response, err := client.Remove(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return false, err
	}

	response := r.(*api.RemoveResponse)
	return response.Removed, nil
}

func (s *setPartition) Contains(ctx context.Context, value string) (bool, error) {
	response, err := s.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewSetServiceClient(conn)
		request := &api.ContainsRequest{
			Header: header,
			Value:  value,
		}
		response, err := client.Contains(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return false, err
	}
	return response.(*api.ContainsResponse).Contains, nil
}

func (s *setPartition) Len(ctx context.Context) (int, error) {
	response, err := s.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewSetServiceClient(conn)
		request := &api.SizeRequest{
			Header: header,
		}
		response, err := client.Size(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}
	return int(response.(*api.SizeResponse).Size_), nil
}

func (s *setPartition) Clear(ctx context.Context) error {
	_, err := s.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewSetServiceClient(conn)
		request := &api.ClearRequest{
			Header: header,
		}
		response, err := client.Clear(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	return err
}

func (s *setPartition) Elements(ctx context.Context, ch chan<- string) error {
	stream, err := s.instance.DoQueryStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewSetServiceClient(conn)
		request := &api.IterateRequest{
			Header: header,
		}
		return client.Iterate(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.SetService_IterateClient).Recv()
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
			ch <- event.(*api.IterateResponse).Value
		}
	}()
	return nil
}

func (s *setPartition) Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error {
	stream, err := s.instance.DoCommandStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewSetServiceClient(conn)
		request := &api.EventRequest{
			Header: header,
		}
		for _, opt := range opts {
			opt.beforeWatch(request)
		}
		return client.Events(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.SetService_EventsClient).Recv()
		if err != nil {
			return nil, nil, err
		}
		for _, opt := range opts {
			opt.afterWatch(response)
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
			var t EventType
			switch response.Type {
			case api.EventResponse_NONE:
				t = EventNone
			case api.EventResponse_ADDED:
				t = EventAdded
			case api.EventResponse_REMOVED:
				t = EventRemoved
			}

			ch <- &Event{
				Type:  t,
				Value: response.Value,
			}
		}
	}()
	return nil
}

func (s *setPartition) Close(ctx context.Context) error {
	return s.instance.Close(ctx)
}

func (s *setPartition) Delete(ctx context.Context) error {
	return s.instance.Delete(ctx)
}
