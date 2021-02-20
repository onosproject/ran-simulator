// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cells

import (
	"github.com/onosproject/onos-lib-go/pkg/errors"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"math/rand"
	"reflect"
	"sync"
)

var log = liblog.GetLogger("store", "cells")

// CellRegistry tracks inventory of simulated cells.
type CellRegistry interface {
	// AddCell adds the specified cell to the registry
	AddCell(cell *model.Cell) error

	// GetCell retrieves the cell with the specified ECGI
	GetCell(ecgi types.ECGI) (*model.Cell, error)

	// UpdateCell updates the cell
	UpdateCell(Cell *model.Cell) error

	// DeleteCell deletes the cell with the specified ECGI
	DeleteCell(ecgi types.ECGI) (*model.Cell, error)

	// WatchCells watches the cell inventory events using the supplied channel
	WatchCells(ch chan<- CellEvent, options ...WatchOptions)

	// GetRandomCell retrieves a random cell from the registry
	GetRandomCell() *model.Cell
}

// CellEvent represents a change in the cell inventory
type CellEvent struct {
	Cell *model.Cell
	Type uint8
}

// WatchOptions allows tailoring the WatchCells behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

type cellWatcher struct {
	ch chan<- CellEvent
}

func (r *cellRegistry) notify(cell *model.Cell, eventType uint8) {
	event := CellEvent{
		Cell: cell,
		Type: eventType,
	}
	for _, watcher := range r.watchers {
		watcher.ch <- event
	}
}

type cellRegistry struct {
	lock     sync.RWMutex
	cells    map[types.ECGI]*model.Cell
	watchers []cellWatcher
}

// NewCellRegistry creates a new store abstraction from the specified fixed cell map.
func NewCellRegistry(cells map[string]model.Cell) CellRegistry {
	log.Infof("Creating registry from model with %d cells", len(cells))
	reg := &cellRegistry{
		lock:     sync.RWMutex{},
		cells:    make(map[types.ECGI]*model.Cell),
		watchers: make([]cellWatcher, 0, 8),
	}

	// Copy the Cells into our own map
	for _, c := range cells {
		cell := c // avoids scopelint issue
		reg.cells[cell.ECGI] = &cell
	}

	log.Infof("Created registry primed with %d cells", len(reg.cells))
	return reg
}

func (r *cellRegistry) AddCell(cell *model.Cell) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.cells[cell.ECGI]; ok {
		return errors.New(errors.NotFound, "cell with EnbID already exists")
	}

	r.cells[cell.ECGI] = cell
	r.notify(cell, ADDED)
	return nil

}

func (r *cellRegistry) GetCell(ecgi types.ECGI) (*model.Cell, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if Cell, ok := r.cells[ecgi]; ok {
		return Cell, nil
	}

	return nil, errors.New(errors.NotFound, "cell not found")
}

func (r *cellRegistry) UpdateCell(cell *model.Cell) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.cells[cell.ECGI]; ok {
		r.cells[cell.ECGI] = cell
		r.notify(cell, UPDATED)
		return nil
	}

	return errors.New(errors.NotFound, "cell not found")
}

func (r *cellRegistry) DeleteCell(ecgi types.ECGI) (*model.Cell, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if cell, ok := r.cells[ecgi]; ok {
		delete(r.cells, ecgi)
		r.notify(cell, DELETED)
		return cell, nil
	}
	return nil, errors.New(errors.NotFound, "cell not found")
}

const (
	// NONE indicates no change event
	NONE uint8 = 0

	// ADDED indicates new cell was added
	ADDED uint8 = 1

	// UPDATED indicates an existing cell was updated
	UPDATED uint8 = 2

	// DELETED indicates a cell was deleted
	DELETED uint8 = 3
)

func (r *cellRegistry) WatchCells(ch chan<- CellEvent, options ...WatchOptions) {
	log.Infof("WatchCells: %v (#%d)\n", options, len(r.cells))
	monitor := len(options) == 0 || options[0].Monitor
	replay := len(options) > 0 && options[0].Replay
	go func() {
		watcher := cellWatcher{ch: ch}
		if monitor {
			r.lock.RLock()
			r.watchers = append(r.watchers, watcher)
			r.lock.RUnlock()
		}

		if replay {
			r.lock.RLock()
			defer r.lock.RUnlock()
			for _, cell := range r.cells {
				ch <- CellEvent{Cell: cell, Type: NONE}
			}
			if !monitor {
				close(ch)
			}
		}
	}()
}

func (r *cellRegistry) GetRandomCell() *model.Cell {
	keys := reflect.ValueOf(r.cells).MapKeys()
	ecgi := types.ECGI(keys[rand.Intn(len(keys))].Uint())
	return r.cells[ecgi]
}
