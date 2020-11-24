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

package log

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/log"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"google.golang.org/grpc"
)

// Type is the log type
const Type primitive.Type = "Log"

// Index is the index of an entry
type Index uint64

// Client provides an API for creating IndexedMaps
type Client interface {
	// GetLog gets the log instance of the given name
	GetLog(ctx context.Context, name string) (Log, error)
}

// Log is a distributed log
type Log interface {
	primitive.Primitive

	// Appends appends the given value to the end of the log
	Append(ctx context.Context, value []byte) (*Entry, error)

	// Get gets the value of the given index
	Get(ctx context.Context, index Index, opts ...GetOption) (*Entry, error)

	// FirstIndex gets the first index in the log
	FirstIndex(ctx context.Context) (Index, error)

	// LastIndex gets the last index in the log
	LastIndex(ctx context.Context) (Index, error)

	// PrevIndex gets the index before the given index
	PrevIndex(ctx context.Context, index Index) (Index, error)

	// NextIndex gets the index after the given index
	NextIndex(ctx context.Context, index Index) (Index, error)

	// FirstEntry gets the first entry in the log
	FirstEntry(ctx context.Context) (*Entry, error)

	// LastEntry gets the last entry in the log
	LastEntry(ctx context.Context) (*Entry, error)

	// PrevEntry gets the entry before the given index
	PrevEntry(ctx context.Context, index Index) (*Entry, error)

	// NextEntry gets the entry after the given index
	NextEntry(ctx context.Context, index Index) (*Entry, error)

	// Remove removes an entry from the log
	Remove(ctx context.Context, index Index, opts ...RemoveOption) (*Entry, error)

	// Size returns the number of entries in the log
	Size(ctx context.Context) (int, error)

	// Clear removes all entries from the log
	Clear(ctx context.Context) error

	// Entries lists the entries in the log
	// This is a non-blocking method. If the method returns without error, index/value paids will be pushed on to the
	// given channel and the channel will be closed once all entries have been read from the map.
	Entries(ctx context.Context, ch chan<- *Entry) error

	// Watch watches the log for changes
	// This is a non-blocking method. If the method returns without error, log events will be pushed onto
	// the given channel in the order in which they occur.
	Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error
}

// Entry is an indexed key/value pair
type Entry struct {
	// Index is the unique, monotonically increasing, globally unique index of the entry. The index is static
	// for the lifetime of a key.
	Index Index

	// Value is the value of the pair
	Value []byte

	// Timestamp
	Timestamp time.Time
}

func (kv Entry) String() string {
	return fmt.Sprintf("index: %d\nvalue: %s\n", kv.Index, string(kv.Value))
}

// EventType is the type of a log event
type EventType string

const (
	// EventNone indicates the event is not a change event
	EventNone EventType = ""

	// EventAppended indicates an entry was appended to the log
	EventAppended EventType = "appended"

	// EventRemoved indicates an entry was removed from the log
	EventRemoved EventType = "removed"
)

// Event is a log change event
type Event struct {
	// Type indicates the change event type
	Type EventType

	// Entry is the event entry
	Entry *Entry
}

// New creates a new log primitive
func New(ctx context.Context, name primitive.Name, partitions []*primitive.Session) (Log, error) {
	i, err := util.GetPartitionIndex(name.Name, len(partitions))
	if err != nil {
		return nil, err
	}
	return newLog(ctx, name, partitions[i])
}

// newLog creates a new Log for the given partition
func newLog(ctx context.Context, name primitive.Name, partition *primitive.Session) (*log, error) {
	instance, err := primitive.NewInstance(ctx, name, partition, &primitiveHandler{})
	if err != nil {
		return nil, err
	}
	return &log{
		name:     name,
		instance: instance,
	}, nil
}

// log is the default single-partition implementation of Log
type log struct {
	name     primitive.Name
	instance *primitive.Instance
}

func (l *log) Name() primitive.Name {
	return l.name
}

func (l *log) Append(ctx context.Context, value []byte) (*Entry, error) {
	r, err := l.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
		request := &api.AppendRequest{
			Header: header,
			Value:  value,
		}
		response, err := client.Append(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}

	response := r.(*api.AppendResponse)
	if response.Status == api.ResponseStatus_OK {
		return &Entry{
			Index: Index(response.Index),
			Value: value,
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return &Entry{
			Index:     Index(response.Index),
			Value:     value,
			Timestamp: response.Timestamp,
		}, nil
	}
}

func (l *log) Get(ctx context.Context, index Index, opts ...GetOption) (*Entry, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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

	return &Entry{
		Index:     Index(response.Index),
		Value:     response.Value,
		Timestamp: response.Timestamp,
	}, nil

}

func (l *log) GetIndex(ctx context.Context, index Index, opts ...GetOption) (*Entry, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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
	return &Entry{
		Index:     Index(response.Index),
		Value:     response.Value,
		Timestamp: response.Timestamp,
	}, nil

}

func (l *log) FirstIndex(ctx context.Context) (Index, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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
	return Index(response.Index), nil
}

func (l *log) LastIndex(ctx context.Context) (Index, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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
	return Index(response.Index), nil

}

func (l *log) PrevIndex(ctx context.Context, index Index) (Index, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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
	return Index(response.Index), nil

}

func (l *log) NextIndex(ctx context.Context, index Index) (Index, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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
	return Index(response.Index), nil
}

func (l *log) FirstEntry(ctx context.Context) (*Entry, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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

	return &Entry{
		Index:     Index(response.Index),
		Value:     response.Value,
		Timestamp: response.Timestamp,
	}, nil

}

func (l *log) LastEntry(ctx context.Context) (*Entry, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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

	return &Entry{
		Index:     Index(response.Index),
		Value:     response.Value,
		Timestamp: response.Timestamp,
	}, nil

}

func (l *log) PrevEntry(ctx context.Context, index Index) (*Entry, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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
	return &Entry{
		Index: Index(response.Index),
		Value: response.Value,
	}, nil

}

func (l *log) NextEntry(ctx context.Context, index Index) (*Entry, error) {
	r, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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
	return &Entry{
		Index:     Index(response.Index),
		Value:     response.Value,
		Timestamp: response.Timestamp,
	}, nil
}

func (l *log) Remove(ctx context.Context, index Index, opts ...RemoveOption) (*Entry, error) {
	r, err := l.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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
			Index: Index(response.Index),
			Value: response.PreviousValue,
		}, nil
	} else if response.Status == api.ResponseStatus_PRECONDITION_FAILED {
		return nil, errors.New("write condition failed")
	} else if response.Status == api.ResponseStatus_WRITE_LOCK {
		return nil, errors.New("write lock failed")
	} else {
		return nil, nil
	}
}

func (l *log) Size(ctx context.Context) (int, error) {
	response, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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

func (l *log) Clear(ctx context.Context) error {
	_, err := l.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLogServiceClient(conn)
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

func (l *log) Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error {
	stream, err := l.instance.DoCommandStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewLogServiceClient(conn)
		request := &api.EventRequest{
			Header: header,
		}
		for _, opt := range opts {
			opt.beforeWatch(request)
		}
		return client.Events(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.LogService_EventsClient).Recv()
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
			case api.EventResponse_APPENDED:
				t = EventAppended
			case api.EventResponse_REMOVED:
				t = EventRemoved
			}
			ch <- &Event{
				Type: t,
				Entry: &Entry{
					Index:     Index(response.Index),
					Value:     response.Value,
					Timestamp: response.Timestamp,
				},
			}
		}
	}()
	return nil
}

func (l *log) Entries(ctx context.Context, ch chan<- *Entry) error {
	stream, err := l.instance.DoQueryStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewLogServiceClient(conn)
		request := &api.EntriesRequest{
			Header: header,
		}
		return client.Entries(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.LogService_EntriesClient).Recv()
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
				Index:     Index(response.Index),
				Value:     response.Value,
				Timestamp: response.Timestamp,
			}
		}
	}()
	return nil
}

func (l *log) Close(ctx context.Context) error {
	return l.instance.Close(ctx)
}

func (l *log) Delete(ctx context.Context) error {
	return l.instance.Delete(ctx)
}
