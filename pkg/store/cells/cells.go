// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cells

import (
	"context"
	"math/rand"
	"reflect"
	"sync"

	"github.com/google/uuid"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	"github.com/onosproject/ran-simulator/pkg/store/watcher"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
)

var log = liblog.GetLogger("store", "cells")

// Store tracks inventory of simulated cells.
type Store interface {
	// AddCell adds the specified cell to the registry
	Add(ctx context.Context, cell *model.Cell) error

	// GetCell retrieves the cell with the specified ECGI
	Get(ctx context.Context, ecgi types.ECGI) (*model.Cell, error)

	// UpdateCell updates the cell
	Update(ctx context.Context, Cell *model.Cell) error

	// DeleteCell deletes the cell with the specified ECGI
	Delete(ctx context.Context, ecgi types.ECGI) (*model.Cell, error)

	// WatchCells watches the cell inventory events using the supplied channel
	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error

	// GetRandomCell retrieves a random cell from the registry
	GetRandomCell() (*model.Cell, error)
}

// WatchOptions allows tailoring the WatchCells behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

type store struct {
	lock      sync.RWMutex
	cells     map[types.ECGI]*model.Cell
	nodeStore nodes.Store
	watchers  *watcher.Watchers
}

// NewCellRegistry creates a new store abstraction from the specified fixed cell map.
func NewCellRegistry(cells map[string]model.Cell, nodeStore nodes.Store) Store {
	log.Infof("Creating registry from model with %d cells", len(cells))
	watchers := watcher.NewWatchers()
	reg := &store{
		lock:      sync.RWMutex{},
		cells:     make(map[types.ECGI]*model.Cell),
		nodeStore: nodeStore,
		watchers:  watchers,
	}

	// Copy the Cells into our own map
	for _, c := range cells {
		cell := c // avoids scopelint issue
		reg.cells[cell.ECGI] = &cell
	}

	log.Infof("Created registry primed with %d cells", len(reg.cells))
	return reg
}

// Add adds a cell
func (s *store) Add(ctx context.Context, cell *model.Cell) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.cells[cell.ECGI]; ok {
		return errors.New(errors.NotFound, "cell with EnbID already exists")
	}

	s.cells[cell.ECGI] = cell
	cellEvent := event.Event{
		Key:   cell.ECGI,
		Value: cell,
		Type:  Created,
	}
	s.watchers.Send(cellEvent)
	return nil

}

// Get gets a cell
func (s *store) Get(ctx context.Context, ecgi types.ECGI) (*model.Cell, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if cell, ok := s.cells[ecgi]; ok {
		return cell, nil
	}

	return nil, errors.New(errors.NotFound, "cell not found")
}

// Update updates a cell
func (s *store) Update(ctx context.Context, cell *model.Cell) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.cells[cell.ECGI]; ok {
		s.cells[cell.ECGI] = cell
		cellEvent := event.Event{
			Key:   cell.ECGI,
			Value: cell,
			Type:  Updated,
		}
		s.watchers.Send(cellEvent)
		return nil
	}

	return errors.New(errors.NotFound, "cell not found")
}

// Delete deletes a cell
func (s *store) Delete(ctx context.Context, ecgi types.ECGI) (*model.Cell, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if cell, ok := s.cells[ecgi]; ok {
		delete(s.cells, ecgi)
		deleteEvent := event.Event{
			Key:   cell.ECGI,
			Value: cell,
			Type:  Deleted,
		}
		s.watchers.Send(deleteEvent)
		err := s.nodeStore.PruneCell(ctx, ecgi)
		if err != nil {
			return nil, err
		}
		return cell, nil
	}
	return nil, errors.New(errors.NotFound, "cell not found")
}

// Watch watch cell events
func (s *store) Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error {
	log.Infof("WatchCells: %v (#%d)\n", options, len(s.cells))
	replay := len(options) > 0 && options[0].Replay
	id := uuid.New()
	err := s.watchers.AddWatcher(id, ch)
	if err != nil {
		log.Error(err)
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
		go func() {
			for _, cell := range s.cells {
				ch <- event.Event{
					Key:   cell.ECGI,
					Value: cell,
					Type:  None,
				}
			}
		}()

	}
	return nil
}

func (s *store) GetRandomCell() (*model.Cell, error) {
	keys := reflect.ValueOf(s.cells).MapKeys()
	ecgi := types.ECGI(keys[rand.Intn(len(keys))].Uint())
	return s.cells[ecgi], nil
}
