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
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"sync"
)

// Type is the set type
const Type primitive.Type = "Set"

// Client provides an API for creating Sets
type Client interface {
	// GetSet gets the Set instance of the given name
	GetSet(ctx context.Context, name string) (Set, error)
}

// Set provides a distributed set data structure
// The set values are defines as strings. To store more complex types in the set, encode values to strings e.g.
// using base 64 encoding.
type Set interface {
	primitive.Primitive

	// Add adds a value to the set
	Add(ctx context.Context, value string) (bool, error)

	// Remove removes a value from the set
	// A bool indicating whether the set contained the given value will be returned
	Remove(ctx context.Context, value string) (bool, error)

	// Contains returns a bool indicating whether the set contains the given value
	Contains(ctx context.Context, value string) (bool, error)

	// Len gets the set size in number of elements
	Len(ctx context.Context) (int, error)

	// Clear removes all values from the set
	Clear(ctx context.Context) error

	// Elements lists the elements in the set
	Elements(ctx context.Context, ch chan<- string) error

	// Watch watches the set for changes
	// This is a non-blocking method. If the method returns without error, set events will be pushed onto
	// the given channel.
	Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error
}

// EventType is the type of a set event
type EventType string

const (
	// EventNone indicates that the event is not in reaction to a state change
	EventNone EventType = ""

	// EventAdded indicates a value was added to the set
	EventAdded EventType = "added"

	// EventRemoved indicates a value was removed from the set
	EventRemoved EventType = "removed"
)

// Event is a set change event
type Event struct {
	// Type is the change event type
	Type EventType

	// Value is the value that changed
	Value string
}

// New creates a new partitioned set primitive
func New(ctx context.Context, name primitive.Name, partitions []*primitive.Session) (Set, error) {
	results, err := util.ExecuteOrderedAsync(len(partitions), func(i int) (interface{}, error) {
		return newPartition(ctx, name, partitions[i])
	})
	if err != nil {
		return nil, err
	}

	sets := make([]Set, len(results))
	for i, result := range results {
		sets[i] = result.(Set)
	}

	return &set{
		name:       name,
		partitions: sets,
	}, nil
}

// set is the partitioned implementation of Set
type set struct {
	name       primitive.Name
	partitions []Set
}

func (s *set) Name() primitive.Name {
	return s.name
}

func (s *set) getPartition(key string) (Set, error) {
	i, err := util.GetPartitionIndex(key, len(s.partitions))
	if err != nil {
		return nil, err
	}
	return s.partitions[i], nil
}

func (s *set) Add(ctx context.Context, value string) (bool, error) {
	partition, err := s.getPartition(value)
	if err != nil {
		return false, err
	}
	return partition.Add(ctx, value)
}

func (s *set) Remove(ctx context.Context, value string) (bool, error) {
	partition, err := s.getPartition(value)
	if err != nil {
		return false, err
	}
	return partition.Remove(ctx, value)
}

func (s *set) Contains(ctx context.Context, value string) (bool, error) {
	partition, err := s.getPartition(value)
	if err != nil {
		return false, err
	}
	return partition.Contains(ctx, value)
}

func (s *set) Len(ctx context.Context) (int, error) {
	results, err := util.ExecuteAsync(len(s.partitions), func(i int) (interface{}, error) {
		return s.partitions[i].Len(ctx)
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

func (s *set) Elements(ctx context.Context, ch chan<- string) error {
	n := len(s.partitions)
	wg := sync.WaitGroup{}
	wg.Add(n)

	go func() {
		wg.Wait()
		close(ch)
	}()

	return util.IterAsync(n, func(i int) error {
		partitionCh := make(chan string)
		go func() {
			for kv := range partitionCh {
				ch <- kv
			}
			wg.Done()
		}()
		return s.partitions[i].Elements(ctx, partitionCh)
	})
}

func (s *set) Clear(ctx context.Context) error {
	return util.IterAsync(len(s.partitions), func(i int) error {
		return s.partitions[i].Clear(ctx)
	})
}

func (s *set) Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error {
	n := len(s.partitions)
	wg := sync.WaitGroup{}
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
		return s.partitions[i].Watch(ctx, partitionCh, opts...)
	})
}

func (s *set) Close(ctx context.Context) error {
	return util.IterAsync(len(s.partitions), func(i int) error {
		return s.partitions[i].Close(ctx)
	})
}

func (s *set) Delete(ctx context.Context) error {
	return util.IterAsync(len(s.partitions), func(i int) error {
		return s.partitions[i].Delete(ctx)
	})
}
