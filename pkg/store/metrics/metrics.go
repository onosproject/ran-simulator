// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"context"
	"sync"

	"github.com/google/uuid"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/onosproject/ran-simulator/pkg/store/watcher"
)

var log = liblog.GetLogger()

// Store tracks arbitrary (named) scalar metrics on per entity (cell, node) basis.
type Store interface {
	// ListEntities retrieves all entities that presently have metrics associated with them
	ListEntities(ctx context.Context) ([]uint64, error)

	// Set applies the specified metric value on the given entity
	Set(ctx context.Context, entityID uint64, name string, value interface{}) error

	// Get retrieves the specified metric value on the given entity
	Get(ctx context.Context, entityID uint64, name string) (interface{}, bool)

	// Delete removes the specified metric
	Delete(ctx context.Context, entityID uint64, name string) error

	// DeleteAll removes all metrics for the specified entity
	DeleteAll(ctx context.Context, entityID uint64) error

	// Get retrieves all metrics of the specified entity as a map
	List(ctx context.Context, entityID uint64) (map[string]interface{}, error)

	// WatchMetrics monitors changes to the metrics
	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error

	// Clear clears all metrics; no events will be generated
	Clear(ctx context.Context)
}

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
}

type store struct {
	mu       sync.RWMutex
	metrics  map[Key]interface{}
	watchers *watcher.Watchers
}

// NewMetricsStore returns a newly created metric store
func NewMetricsStore() Store {
	log.Infof("Creating metrics store")
	watchers := watcher.NewWatchers()
	return &store{
		mu:       sync.RWMutex{},
		metrics:  make(map[Key]interface{}),
		watchers: watchers,
	}
}

// Clear clears all metrics; no events will be generated
func (s *store) Clear(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key := range s.metrics {
		delete(s.metrics, key)
	}
}

// ListEntities retrieves all entities that presently have metrics associated with them
func (s *store) ListEntities(ctx context.Context) ([]uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idmap := make(map[uint64]uint64)
	for k := range s.metrics {
		id := k.EntityID
		idmap[id] = id
	}

	entities := make([]uint64, 0, len(idmap))
	for k := range idmap {
		entities = append(entities, k)
	}

	return entities, nil
}

// Generate composite key from entity ID and metric name
func key(entityID uint64, name string) Key {
	return Key{
		EntityID: entityID,
		Name:     name,
	}
}

func metricEvent(key Key, value interface{}, eventType interface{}) event.Event {
	return event.Event{
		Key:   key,
		Value: value,
		Type:  eventType,
	}
}

// Set applies the specified metric value on the given entity
func (s *store) Set(ctx context.Context, entityID uint64, name string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(entityID, name)
	s.metrics[k] = value
	s.watchers.Send(metricEvent(k, value, Updated))
	return nil
}

// Get retrieves the specified metric value on the given entity
func (s *store) Get(ctx context.Context, entityID uint64, name string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.metrics[key(entityID, name)]; ok {
		return v, ok
	}
	return nil, false
}

// Delete removes the specified metric
func (s *store) Delete(ctx context.Context, entityID uint64, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(entityID, name)
	delete(s.metrics, k)
	s.watchers.Send(metricEvent(k, nil, Deleted))
	return nil
}

// DeleteAll removes all metrics for the specified entity
func (s *store) DeleteAll(ctx context.Context, entityID uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.metrics {
		if k.EntityID == entityID {
			delete(s.metrics, k)
			s.watchers.Send(metricEvent(k, v, Deleted))
		}
	}
	return nil
}

// Get retrieves all metrics of the specified entity as a map
func (s *store) List(ctx context.Context, entityID uint64) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	metrics := make(map[string]interface{})
	for k, v := range s.metrics {
		if k.EntityID == entityID {
			metrics[k.Name] = v
		}
	}
	return metrics, nil
}

// WatchMetrics monitors changes to the metrics
func (s *store) Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error {
	log.Debug("Watching metric changes")
	id := uuid.New()
	err := s.watchers.AddWatcher(id, ch)
	if err != nil {
		log.Error(err)
		return err
	}
	go func() {
		<-ctx.Done()
		err = s.watchers.RemoveWatcher(id)
		if err != nil {
			log.Error(err)
		}
		close(ch)

	}()
	return nil
}
