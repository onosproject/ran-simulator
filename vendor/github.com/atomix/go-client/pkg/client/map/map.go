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
	"fmt"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"math"
	"sync"
	"time"
)

// Type is the map type
const Type primitive.Type = "Map"

// Client provides an API for creating Maps
type Client interface {
	// GetMap gets the Map instance of the given name
	GetMap(ctx context.Context, name string, opts ...Option) (Map, error)
}

// Map is a distributed set of keys and values
type Map interface {
	primitive.Primitive

	// Put sets a key/value pair in the map
	Put(ctx context.Context, key string, value []byte, opts ...PutOption) (*Entry, error)

	// Get gets the value of the given key
	Get(ctx context.Context, key string, opts ...GetOption) (*Entry, error)

	// Remove removes a key from the map
	Remove(ctx context.Context, key string, opts ...RemoveOption) (*Entry, error)

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

// Version is an entry version
type Version uint64

// Entry is a versioned key/value pair
type Entry struct {
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

// New creates a new partitioned Map
func New(ctx context.Context, name primitive.Name, sessions []*primitive.Session, opts ...Option) (Map, error) {
	options := &options{}
	for _, opt := range opts {
		opt.apply(options)
	}

	results, err := util.ExecuteOrderedAsync(len(sessions), func(i int) (interface{}, error) {
		if options.cached {
			return newPartition(ctx, name, sessions[i], WithCache(int(math.Max(float64(options.cacheSize/len(sessions)), 1))))
		}
		return newPartition(ctx, name, sessions[i])
	})
	if err != nil {
		return nil, err
	}

	maps := make([]Map, len(results))
	for i, result := range results {
		maps[i] = result.(Map)
	}

	return &_map{
		name:       name,
		partitions: maps,
	}, nil
}

// _map is the default single-partition implementation of Map
type _map struct {
	name       primitive.Name
	partitions []Map
}

func (m *_map) Name() primitive.Name {
	return m.name
}

func (m *_map) getPartition(key string) (Map, error) {
	i, err := util.GetPartitionIndex(key, len(m.partitions))
	if err != nil {
		return nil, err
	}
	return m.partitions[i], nil
}

func (m *_map) Put(ctx context.Context, key string, value []byte, opts ...PutOption) (*Entry, error) {
	session, err := m.getPartition(key)
	if err != nil {
		return nil, err
	}
	return session.Put(ctx, key, value, opts...)
}

func (m *_map) Get(ctx context.Context, key string, opts ...GetOption) (*Entry, error) {
	session, err := m.getPartition(key)
	if err != nil {
		return nil, err
	}
	entry, err := session.Get(ctx, key, opts...)
	if err != nil {
		return nil, err
	} else if entry.Value == nil {
		return nil, nil
	}
	return entry, nil
}

func (m *_map) Remove(ctx context.Context, key string, opts ...RemoveOption) (*Entry, error) {
	session, err := m.getPartition(key)
	if err != nil {
		return nil, err
	}
	return session.Remove(ctx, key, opts...)
}

func (m *_map) Len(ctx context.Context) (int, error) {
	results, err := util.ExecuteAsync(len(m.partitions), func(i int) (interface{}, error) {
		return m.partitions[i].Len(ctx)
	})
	if err != nil {
		return 0, err
	}

	total := 0
	for _, result := range results {
		total += result.(int)
	}
	return total, nil
}

func (m *_map) Entries(ctx context.Context, ch chan<- *Entry) error {
	n := len(m.partitions)
	wg := sync.WaitGroup{}
	wg.Add(n)

	go func() {
		wg.Wait()
		close(ch)
	}()

	return util.IterAsync(n, func(i int) error {
		partitionCh := make(chan *Entry)
		go func() {
			for kv := range partitionCh {
				ch <- kv
			}
			wg.Done()
		}()
		return m.partitions[i].Entries(ctx, partitionCh)
	})
}

func (m *_map) Clear(ctx context.Context) error {
	return util.IterAsync(len(m.partitions), func(i int) error {
		return m.partitions[i].Clear(ctx)
	})
}

func (m *_map) Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error {
	n := len(m.partitions)
	wg := &sync.WaitGroup{}
	wg.Add(n)

	go func() {
		wg.Wait()
		close(ch)
	}()

	return util.IterAsync(n, func(i int) error {
		partitionCh := make(chan *Event)
		go func() {
			for event := range partitionCh {
				ch <- event
			}
			wg.Done()
		}()
		return m.partitions[i].Watch(ctx, partitionCh, opts...)
	})
}

func (m *_map) Close(ctx context.Context) error {
	return util.IterAsync(len(m.partitions), func(i int) error {
		return m.partitions[i].Close(ctx)
	})
}

func (m *_map) Delete(ctx context.Context) error {
	return util.IterAsync(len(m.partitions), func(i int) error {
		return m.partitions[i].Delete(ctx)
	})
}
