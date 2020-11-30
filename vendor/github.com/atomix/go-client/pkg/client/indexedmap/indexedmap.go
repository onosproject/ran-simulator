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

package indexedmap

import (
	"context"
	"errors"
	"fmt"
	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/indexedmap"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"google.golang.org/grpc"
	"time"
)

// Type is the indexedmap type
const Type primitive.Type = "IndexedMap"

// Index is the index of an entry
type Index uint64

// Version is the version of an entry
type Version uint64

// Client provides an API for creating IndexedMaps
type Client interface {
	// GetIndexedMap gets the IndexedMap instance of the given name
	GetIndexedMap(ctx context.Context, name string) (IndexedMap, error)
}

// IndexedMap is a distributed linked map
type IndexedMap interface {
	primitive.Primitive

	// Append appends the given key/value to the map
	Append(ctx context.Context, key string, value []byte) (*Entry, error)

	// Put appends the given key/value to the map
	Put(ctx context.Context, key string, value []byte) (*Entry, error)

	// Set sets the given index in the map
	Set(ctx context.Context, index Index, key string, value []byte, opts ...SetOption) (*Entry, error)

	// Get gets the value of the given key
	Get(ctx context.Context, key string, opts ...GetOption) (*Entry, error)

	// GetIndex gets the entry at the given index
	GetIndex(ctx context.Context, index Index, opts ...GetOption) (*Entry, error)

	// FirstIndex gets the first index in the map
	FirstIndex(ctx context.Context) (Index, error)

	// LastIndex gets the last index in the map
	LastIndex(ctx context.Context) (Index, error)

	// PrevIndex gets the index before the given index
	PrevIndex(ctx context.Context, index Index) (Index, error)

	// NextIndex gets the index after the given index
	NextIndex(ctx context.Context, index Index) (Index, error)

	// FirstEntry gets the first entry in the map
	FirstEntry(ctx context.Context) (*Entry, error)

	// LastEntry gets the last entry in the map
	LastEntry(ctx context.Context) (*Entry, error)

	// PrevEntry gets the entry before the given index
	PrevEntry(ctx context.Context, index Index) (*Entry, error)

	// NextEntry gets the entry after the given index
	NextEntry(ctx context.Context, index Index) (*Entry, error)

	// Replace replaces the given key with the given value
	Replace(ctx context.Context, key string, value []byte, opts ...ReplaceOption) (*Entry, error)

	// ReplaceIndex replaces the given index with the given value
	ReplaceIndex(ctx context.Context, index Index, value []byte, opts ...ReplaceOption) (*Entry, error)

	// Remove removes a key from the map
	Remove(ctx context.Context, key string, opts ...RemoveOption) (*Entry, error)

	// RemoveIndex removes an index from the map
	RemoveIndex(ctx context.Context, index Index, opts ...RemoveOption) (*Entry, error)

	// Len returns the number of entries in the map
	Len(ctx context.Context) (int, error)

	// Clear removes all entries from the map
	Clear(ctx context.Context) error

	// Entries lists the entries in the map
	// This is a non-blocking method. If the method returns without error, key/value paids will be pushed on to the
	// given channel and the channel will be closed once all entries have been read from the map.
	Entries(ctx context.Context, ch chan<- *Entry) error

	// Watch watches the map for changes
	// This is a non-blocking method. If the method returns without error, map events will be pushed onto
	// the given channel in the order in which they occur.
	Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error
}

// Entry is an indexed key/value pair
type Entry struct {
	// Index is the unique, monotonically increasing, globally unique index of the entry. The index is static
	// for the lifetime of a key.
	Index Index

	// Version is the unique, monotonically increasing version number for the key/value pair. The version is
	// suitable for use in optimistic locking.
	Version Version

	// Key is the key of the pair
	Key string

	// Value is the value of the pair
	Value []byte

	// Created is the time at which the key was created
	Created time.Time

	// Updated is the time at which the key was last updated
	Updated time.Time
}

func (kv Entry) String() string {
	return fmt.Sprintf("key: %s\nvalue: %s\nversion: %d", kv.Key, string(kv.Value), kv.Version)
}

// EventType is the type of a map event
type EventType string

const (
	// EventNone indicates the event is not a change event
	EventNone EventType = ""

	// EventInserted indicates a key was newly created in the map
	EventInserted EventType = "inserted"

	// EventUpdated indicates the value of an existing key was changed
	EventUpdated EventType = "updated"

	// EventRemoved indicates a key was removed from the map
	EventRemoved EventType = "removed"
)

// Event is a map change event
type Event struct {
	// Type indicates the change event type
	Type EventType

	// Entry is the event entry
	Entry *Entry
}

// New creates a new IndexedMap primitive
func New(ctx context.Context, name primitive.Name, partitions []*primitive.Session) (IndexedMap, error) {
	i, err := util.GetPartitionIndex(name.Name, len(partitions))
	if err != nil {
		return nil, err
	}
	return newIndexedMap(ctx, name, partitions[i])
}

// newIndexedMap creates a new IndexedMap for the given partition
func newIndexedMap(ctx context.Context, name primitive.Name, partition *primitive.Session) (*indexedMap, error) {
	instance, err := primitive.NewInstance(ctx, name, partition, &primitiveHandler{})
	if err != nil {
		return nil, err
	}
	return &indexedMap{
		name:     name,
		instance: instance,
	}, nil
}

// indexedMap is the default single-partition implementation of Map
type indexedMap struct {
	name     primitive.Name
	instance *primitive.Instance
}

func (m *indexedMap) Name() primitive.Name {
	return m.name
}

func (m *indexedMap) Append(ctx context.Context, key string, value []byte) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.PutRequest{
			Header:  header,
			Key:     key,
			Value:   value,
			IfEmpty: true,
		}
		response, err := client.Put(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.PutResponse)
	if response.Status == api.ResponseStatus_OK {
		return &Entry{
			Index:   Index(response.Index),
			Key:     key,
			Value:   value,
			Version: Version(response.Header.Index),
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return &Entry{
			Index:   Index(response.Index),
			Key:     key,
			Value:   value,
			Version: Version(response.PreviousVersion),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
}

func (m *indexedMap) Put(ctx context.Context, key string, value []byte) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.PutRequest{
			Header: header,
			Key:    key,
			Value:  value,
		}
		response, err := client.Put(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.PutResponse)
	if response.Status == api.ResponseStatus_OK {
		return &Entry{
			Index:   Index(response.Index),
			Key:     key,
			Value:   value,
			Version: Version(response.Header.Index),
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return &Entry{
			Index:   Index(response.Index),
			Key:     key,
			Value:   value,
			Version: Version(response.PreviousVersion),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
}

func (m *indexedMap) Set(ctx context.Context, index Index, key string, value []byte, opts ...SetOption) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.PutRequest{
			Header: header,
			Index:  uint64(index),
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
	if response.Status == api.ResponseStatus_OK {
		return &Entry{
			Index:   Index(response.Index),
			Key:     key,
			Value:   value,
			Version: Version(response.Header.Index),
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return &Entry{
			Index:   Index(response.Index),
			Key:     key,
			Value:   value,
			Version: Version(response.PreviousVersion),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
}

func (m *indexedMap) Get(ctx context.Context, key string, opts ...GetOption) (*Entry, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
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
	if response.Version != 0 {
		return &Entry{
			Index:   Index(response.Index),
			Key:     response.Key,
			Value:   response.Value,
			Version: Version(response.Version),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
	return nil, nil
}

func (m *indexedMap) GetIndex(ctx context.Context, index Index, opts ...GetOption) (*Entry, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.GetRequest{
			Header: header,
			Index:  uint64(index),
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
	if response.Version != 0 {
		return &Entry{
			Index:   Index(response.Index),
			Key:     response.Key,
			Value:   response.Value,
			Version: Version(response.Version),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
	return nil, nil
}

func (m *indexedMap) FirstIndex(ctx context.Context) (Index, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.FirstEntryRequest{
			Header: header,
		}
		response, err := client.FirstEntry(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}

	response := r.(*api.FirstEntryResponse)
	if response.Version != 0 {
		return Index(response.Index), nil
	}
	return 0, nil
}

func (m *indexedMap) LastIndex(ctx context.Context) (Index, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.LastEntryRequest{
			Header: header,
		}
		response, err := client.LastEntry(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}

	response := r.(*api.LastEntryResponse)
	if response.Version != 0 {
		return Index(response.Index), nil
	}
	return 0, nil
}

func (m *indexedMap) PrevIndex(ctx context.Context, index Index) (Index, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.PrevEntryRequest{
			Header: header,
			Index:  uint64(index),
		}
		response, err := client.PrevEntry(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}

	response := r.(*api.PrevEntryResponse)
	if response.Version != 0 {
		return Index(response.Index), nil
	}
	return 0, nil
}

func (m *indexedMap) NextIndex(ctx context.Context, index Index) (Index, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.NextEntryRequest{
			Header: header,
			Index:  uint64(index),
		}
		response, err := client.NextEntry(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}

	response := r.(*api.NextEntryResponse)
	if response.Version != 0 {
		return Index(response.Index), nil
	}
	return 0, nil
}

func (m *indexedMap) FirstEntry(ctx context.Context) (*Entry, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.FirstEntryRequest{
			Header: header,
		}
		response, err := client.FirstEntry(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.FirstEntryResponse)
	if response.Version != 0 {
		return &Entry{
			Index:   Index(response.Index),
			Key:     response.Key,
			Value:   response.Value,
			Version: Version(response.Version),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
	return nil, err
}

func (m *indexedMap) LastEntry(ctx context.Context) (*Entry, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.LastEntryRequest{
			Header: header,
		}
		response, err := client.LastEntry(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.LastEntryResponse)
	if response.Version != 0 {
		return &Entry{
			Index:   Index(response.Index),
			Key:     response.Key,
			Value:   response.Value,
			Version: Version(response.Version),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
	return nil, err
}

func (m *indexedMap) PrevEntry(ctx context.Context, index Index) (*Entry, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.PrevEntryRequest{
			Header: header,
			Index:  uint64(index),
		}
		response, err := client.PrevEntry(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.PrevEntryResponse)
	if response.Version != 0 {
		return &Entry{
			Index:   Index(response.Index),
			Key:     response.Key,
			Value:   response.Value,
			Version: Version(response.Version),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
	return nil, err
}

func (m *indexedMap) NextEntry(ctx context.Context, index Index) (*Entry, error) {
	r, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.NextEntryRequest{
			Header: header,
			Index:  uint64(index),
		}
		response, err := client.NextEntry(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.NextEntryResponse)
	if response.Version != 0 {
		return &Entry{
			Index:   Index(response.Index),
			Key:     response.Key,
			Value:   response.Value,
			Version: Version(response.Version),
			Created: response.Created,
			Updated: response.Updated,
		}, nil
	}
	return nil, err
}

func (m *indexedMap) Replace(ctx context.Context, key string, value []byte, opts ...ReplaceOption) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.ReplaceRequest{
			Header:   header,
			Key:      key,
			NewValue: value,
		}
		for i := range opts {
			opts[i].beforeReplace(request)
		}
		response, err := client.Replace(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterReplace(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.ReplaceResponse)
	if response.Status == api.ResponseStatus_OK {
		return &Entry{
			Index:   Index(response.Index),
			Key:     key,
			Value:   value,
			Version: Version(response.Header.Index),
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return nil, nil
	}
}

func (m *indexedMap) ReplaceIndex(ctx context.Context, index Index, value []byte, opts ...ReplaceOption) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.ReplaceRequest{
			Header:   header,
			Index:    uint64(index),
			NewValue: value,
		}
		for i := range opts {
			opts[i].beforeReplace(request)
		}
		response, err := client.Replace(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterReplace(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.ReplaceResponse)
	if response.Status == api.ResponseStatus_OK {
		return &Entry{
			Index:   Index(response.Index),
			Key:     response.Key,
			Value:   response.PreviousValue,
			Version: Version(response.PreviousVersion),
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return nil, nil
	}
}

func (m *indexedMap) Remove(ctx context.Context, key string, opts ...RemoveOption) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
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
	if response.Status == api.ResponseStatus_OK {
		return &Entry{
			Index:   Index(response.Index),
			Key:     key,
			Value:   response.PreviousValue,
			Version: Version(response.PreviousVersion),
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return nil, nil
	}
}

func (m *indexedMap) RemoveIndex(ctx context.Context, index Index, opts ...RemoveOption) (*Entry, error) {
	r, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.RemoveRequest{
			Header: header,
			Index:  uint64(index),
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
	if response.Status == api.ResponseStatus_OK {
		return &Entry{
			Index:   Index(response.Index),
			Key:     response.Key,
			Value:   response.PreviousValue,
			Version: Version(response.PreviousVersion),
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return nil, nil
	}
}

func (m *indexedMap) Len(ctx context.Context) (int, error) {
	response, err := m.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
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

func (m *indexedMap) Clear(ctx context.Context) error {
	_, err := m.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
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

func (m *indexedMap) Entries(ctx context.Context, ch chan<- *Entry) error {
	stream, err := m.instance.DoQueryStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.EntriesRequest{
			Header: header,
		}
		return client.Entries(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.IndexedMapService_EntriesClient).Recv()
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
				Index:   Index(response.Index),
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

func (m *indexedMap) Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error {
	stream, err := m.instance.DoCommandStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewIndexedMapServiceClient(conn)
		request := &api.EventRequest{
			Header: header,
		}
		for _, opt := range opts {
			opt.beforeWatch(request)
		}
		return client.Events(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.IndexedMapService_EventsClient).Recv()
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

			// If this is a normal event (not a handshake response), write the event to the watch channel
			var t EventType
			switch response.Type {
			case api.EventResponse_NONE:
				t = EventNone
			case api.EventResponse_INSERTED:
				t = EventInserted
			case api.EventResponse_UPDATED:
				t = EventUpdated
			case api.EventResponse_REMOVED:
				t = EventRemoved
			}
			ch <- &Event{
				Type: t,
				Entry: &Entry{
					Index:   Index(response.Index),
					Key:     response.Key,
					Value:   response.Value,
					Version: Version(response.Version),
					Created: response.Created,
					Updated: response.Updated,
				},
			}
		}
	}()
	return nil
}

func (m *indexedMap) Close(ctx context.Context) error {
	return m.instance.Close(ctx)
}

func (m *indexedMap) Delete(ctx context.Context) error {
	return m.instance.Delete(ctx)
}
