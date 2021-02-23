// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package metrics

import (
	"context"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/onosproject/ran-simulator/pkg/store/watcher"
	"sync"
)

// MetricsStore tracks arbitrary (named) scalar metrics (as float64 values) on per entity (cell, node) basis.
type MetricsStore interface {
	// ListEntities retrieves all entities that presently have metrics associated with them
	ListEntities(ctx context.Context) []uint64

	// Set applies the specified metric value on the given entity
	Set(ctx context.Context, entityID uint64, key string, value float64)

	// Get retrieves the specified metric value on the given entity
	Get(ctx context.Context, entityID uint64, key string) (float64, error)

	// Delete removes the specified metric
	Delete(ctx context.Context, entityID uint64, key string) error

	// DeleteAll removes all metrics for the specified entity
	DeleteAll(ctx context.Context, entityID uint64, key string) error

	// Get retrieves all metrics of the specified entity as a map
	List(ctx context.Context, entityID uint64, key string, value float64) (map[string]float64, error)

	// WatchMetrics monitors changes to the metrics
	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions)
}

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay   bool
	Monitor  bool
	Entities []uint64
}

// Structure to track named metrics as floats for a single entity
type metrics struct {
	mu      sync.RWMutex
	metrics map[string]float64
}

type store struct {
	mu       sync.RWMutex
	metrics  map[uint64]metrics
	watchers *watcher.Watchers
}
