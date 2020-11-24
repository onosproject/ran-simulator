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

package _map //nolint:golint

import (
	"context"
	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"google.golang.org/grpc"
)

func newPartition(ctx context.Context, name primitive.Name, session *primitive.Session, opts ...Option) (Map, error) {
	options := &options{}
	for _, opt := range opts {
		opt.apply(options)
	}

	instance, err := primitive.NewInstance(ctx, name, session, &primitiveHandler{})
	if err != nil {
		return nil, err
	}
	var partition Map = &mapPartition{
		name:     name,
		instance: instance,
	}
	if options.cached {
		cached, err := newCachingMap(partition, options.cacheSize)
		if err != nil {
			return nil, err
		}
		partition = cached
	}
	return partition, nil
}

type mapPartition struct {
	name     primitive.Name
	instance *primitive.Instance
}

func (m *mapPartition) Name() primitive.Name {
	return m.name
}

func (m *mapPartition) Put(ctx context.Context, key string, value []byte, opts ...PutOption) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewMapServiceClient(conn)
		request := &api.PutRequest{
			Header: header,
			Key:    key,
			Value:  value,
		}
		for i := range opts {
			opts[i].beforePut(request)
		}
		response, err := client.Put(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterPut(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.PutResponse)
	return &Entry{
		Key:     key,
		Value:   value,
		Version: Version(response.Header.Index),
		Created: response.Created,
		Updated: response.Updated,
	}, nil
}

func (m *mapPartition) Get(ctx context.Context, key string, opts ...GetOption) (*Entry, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewMapServiceClient(conn)
		request := &api.GetRequest{
			Header: header,
			Key:    key,
		}
		for i := range opts {
			opts[i].beforeGet(request)
		}
		response, err := client.Get(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterGet(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.GetResponse)
	return &Entry{
		Key:     key,
		Value:   response.Value,
		Version: Version(response.Version),
		Created: response.Created,
		Updated: response.Updated,
	}, nil
}

func (m *mapPartition) Remove(ctx context.Context, key string, opts ...RemoveOption) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewMapServiceClient(conn)
		request := &api.RemoveRequest{
			Header: header,
			Key:    key,
		}
		for i := range opts {
			opts[i].beforeRemove(request)
		}
		response, err := client.Remove(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterRemove(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.RemoveResponse)
	return &Entry{
		Key:     key,
		Value:   response.PreviousValue,
		Version: Version(response.PreviousVersion),
	}, nil
}

func (m *mapPartition) Len(ctx context.Context) (int, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewMapServiceClient(conn)
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
	return int(r.(*api.SizeResponse).Size_), nil
}

func (m *mapPartition) Clear(ctx context.Context) error {
	_, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewMapServiceClient(conn)
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

func (m *mapPartition) Entries(ctx context.Context, ch chan<- *Entry) error {
	stream, err := m.instance.DoQueryStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewMapServiceClient(conn)
		request := &api.EntriesRequest{
			Header: header,
		}
		return client.Entries(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.MapService_EntriesClient).Recv()
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
			response := event.(*api.EntriesResponse)
			ch <- &Entry{
				Key:     response.Key,
				Value:   response.Value,
				Version: Version(response.Version),
				Created: response.Created,
				Updated: response.Updated,
			}
		}
	}()
	return nil
}

func (m *mapPartition) Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error {
	stream, err := m.instance.DoCommandStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewMapServiceClient(conn)
		request := &api.EventRequest{
			Header: header,
		}
		for _, opt := range opts {
			opt.beforeWatch(request)
		}
		return client.Events(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.MapService_EventsClient).Recv()
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
			var version Version
			var t EventType
			switch response.Type {
			case api.EventResponse_NONE:
				t = EventNone
				version = Version(response.Version)
			case api.EventResponse_INSERTED:
				t = EventInserted
				version = Version(response.Version)
			case api.EventResponse_UPDATED:
				t = EventUpdated
				version = Version(response.Version)
			case api.EventResponse_REMOVED:
				t = EventRemoved
				version = Version(response.Header.Index)
			}
			ch <- &Event{
				Type: t,
				Entry: &Entry{
					Key:     response.Key,
					Value:   response.Value,
					Version: version,
					Created: response.Created,
					Updated: response.Updated,
				},
			}
		}
	}()
	return nil
}

func (m *mapPartition) Close(ctx context.Context) error {
	return m.instance.Close(ctx)
}

func (m *mapPartition) Delete(ctx context.Context) error {
	return m.instance.Delete(ctx)
}
