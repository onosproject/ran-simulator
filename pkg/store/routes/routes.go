// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package routes

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/onosproject/ran-simulator/pkg/store/watcher"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
)

var log = liblog.GetLogger()

// Store tracks a collection of routes used to simulate UE mobility
type Store interface {
	// Len returns the number of active routes
	Len(ctx context.Context) int

	// Add adds a new route
	Add(ctx context.Context, route *model.Route) error

	// Get retrieves the route for the specified IMSI
	Get(ctx context.Context, imsi types.IMSI) (*model.Route, error)

	// Start sets the route to its start condition
	Start(ctx context.Context, imsi types.IMSI, speedAvg uint32, speedStdDev uint32) error

	// Advance advances to the next waypoint in the specified direction, reversing at route end-points
	Advance(ctx context.Context, imsi types.IMSI) error

	// Delete destroy the specified UE route
	Delete(ctx context.Context, imsi types.IMSI) (*model.Route, error)

	// List returns an array of all routes
	List(ctx context.Context) []*model.Route

	// Watch watches the route events using the supplied channel
	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error

	// Clear removes all routes; no events will be generated
	Clear(ctx context.Context)
}

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

type store struct {
	mu       sync.RWMutex
	routes   map[types.IMSI]*model.Route
	watchers *watcher.Watchers
}

// NewRouteRegistry creates a new route registry
func NewRouteRegistry() Store {
	log.Infof("Creating route registry")
	watchers := watcher.NewWatchers()
	store := &store{
		mu:       sync.RWMutex{},
		routes:   make(map[types.IMSI]*model.Route),
		watchers: watchers,
	}
	log.Infof("Created route registry")
	return store
}

// Clear removes all routes; no events will be generated
func (s *store) Clear(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id := range s.routes {
		delete(s.routes, id)
	}
}

func (s *store) Len(ctx context.Context) int {
	return len(s.routes)
}

func (s *store) Add(ctx context.Context, route *model.Route) error {
	if len(route.Points) < 2 {
		return errors.New(errors.NotFound, "route must have at least two points")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.routes[route.IMSI]; ok {
		return errors.New(errors.NotFound, "route for IMSI already exists")
	}

	s.routes[route.IMSI] = route
	cellEvent := event.Event{
		Key:   route.IMSI,
		Value: route,
		Type:  Created,
	}
	s.watchers.Send(cellEvent)
	return nil
}

// Get gets a UE based on a given imsi
func (s *store) Get(ctx context.Context, imsi types.IMSI) (*model.Route, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if route, ok := s.routes[imsi]; ok {
		return route, nil
	}

	return nil, errors.New(errors.NotFound, "route not found")
}

func (s *store) Start(ctx context.Context, imsi types.IMSI, speedAvg uint32, speedStdDev uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if route, ok := s.routes[imsi]; ok {
		route.NextPoint = 1
		route.Reverse = false
		route.SpeedAvg = speedAvg
		route.SpeedStdDev = speedStdDev
		updateEvent := event.Event{
			Key:   imsi,
			Value: route,
			Type:  Updated,
		}
		s.watchers.Send(updateEvent)
		return nil
	}
	return errors.New(errors.NotFound, "route not found")
}

func (s *store) Advance(ctx context.Context, imsi types.IMSI) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if route, ok := s.routes[imsi]; ok {
		if !route.Reverse && route.NextPoint+1 >= uint32(len(route.Points)) {
			route.Reverse = true
		} else if route.Reverse && route.NextPoint == 0 {
			route.Reverse = false
		}

		if !route.Reverse {
			route.NextPoint = route.NextPoint + 1
		} else {
			route.NextPoint = route.NextPoint - 1
		}
		updateEvent := event.Event{
			Key:   imsi,
			Value: route,
			Type:  Updated,
		}
		s.watchers.Send(updateEvent)
		return nil
	}
	return errors.New(errors.NotFound, "route not found")
}

// Delete deletes a UE based on a given imsi
func (s *store) Delete(ctx context.Context, imsi types.IMSI) (*model.Route, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if route, ok := s.routes[imsi]; ok {
		delete(s.routes, imsi)
		deleteEvent := event.Event{
			Key:   imsi,
			Value: route,
			Type:  Deleted,
		}
		s.watchers.Send(deleteEvent)
		return route, nil
	}
	return nil, errors.New(errors.NotFound, "route not found")
}

func (s *store) List(ctx context.Context) []*model.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Route, 0, len(s.routes))
	for _, route := range s.routes {
		list = append(list, route)
	}
	return list
}

func (s *store) Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error {
	log.Debug("Watching route changes")
	replay := len(options) > 0 && options[0].Replay

	id := uuid.New()
	err := s.watchers.AddWatcher(id, ch)
	if err != nil {
		log.Error(err)
		close(ch)
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

	if replay {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, route := range s.routes {
				ch <- event.Event{
					Key:   route.IMSI,
					Value: route,
					Type:  None,
				}
			}
		}()
	}

	return nil
}
