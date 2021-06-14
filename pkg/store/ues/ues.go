// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package ues

import (
	"context"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/onosproject/ran-simulator/pkg/store/watcher"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
)

const (
	minIMSI = 1000000
	maxIMSI = 9999999
)

var log = liblog.GetLogger("store", "ues")

// Store tracks inventory of user-equipment for the simulation
type Store interface {
	// SetUECount updates the UE count and creates or deletes new UEs as needed
	SetUECount(ctx context.Context, count uint)

	// Len returns the number of active UEs
	Len(ctx context.Context) int

	// LenPerCell returns the number of active UEs per cell
	LenPerCell(ctx context.Context, cellECGI uint64) int

	// MaxUEsPerCell returns the maximum number of active UEs per cell
	MaxUEsPerCell(ctx context.Context, cellECGI uint64) int

	// SetMaxUEsPerCell sets the maximum number of active UEs per cell
	SetMaxUEsPerCell(ctx context.Context, cellECGI uint64, maxNumUEs int)

	// UpdateMaxUEsPerCell updates the maximum number of active UEs for all cells
	UpdateMaxUEsPerCell(ctx context.Context)

	// CreateUEs creates the specified number of UEs
	CreateUEs(ctx context.Context, count uint)

	// Get retrieves the UE with the specified IMSI
	Get(ctx context.Context, imsi types.IMSI) (*model.UE, error)

	// Delete destroy the specified UE
	Delete(ctx context.Context, imsi types.IMSI) (*model.UE, error)

	// MoveToCell update the cell affiliation of the specified UE
	MoveToCell(ctx context.Context, imsi types.IMSI, ecgi types.ECGI, strength float64) error

	// MoveToCoordinate updates the UEs geo location and compass heading
	MoveToCoordinate(ctx context.Context, imsi types.IMSI, location model.Coordinate, heading uint32) error

	// UpdateCells updates the visible cells and their signal strength
	UpdateCells(ctx context.Context, imsi types.IMSI, cells []*model.UECell) error

	// UpdateCell updates the serving cell
	UpdateCell(ctx context.Context, imsi types.IMSI, cell *model.UECell) error

	// ListAllUEs returns an array of all UEs
	ListAllUEs(ctx context.Context) []*model.UE

	// ListUEs returns an array of all UEs associated with the specified cell
	ListUEs(ctx context.Context, ecgi types.ECGI) []*model.UE

	// Watch watches the UE inventory events using the supplied channel
	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error
}

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

type store struct {
	mu        sync.RWMutex
	ues       map[types.IMSI]*model.UE
	maxUEs    map[uint64]int
	cellStore cells.Store
	watchers  *watcher.Watchers
}

// NewUERegistry creates a new user-equipment registry primed with the specified number of UEs to start.
// UEs will be semi-randomly distributed between the specified cells
func NewUERegistry(count uint, cellStore cells.Store) Store {
	log.Infof("Creating registry from model with %d UEs", count)
	watchers := watcher.NewWatchers()
	store := &store{
		mu:        sync.RWMutex{},
		ues:       make(map[types.IMSI]*model.UE),
		maxUEs:    make(map[uint64]int),
		cellStore: cellStore,
		watchers:  watchers,
	}
	ctx := context.Background()
	store.CreateUEs(ctx, count)
	log.Infof("Created registry primed with %d UEs", len(store.ues))
	return store
}

func (s *store) SetUECount(ctx context.Context, count uint) {
	delta := len(s.ues) - int(count)
	if delta < 0 {
		s.CreateUEs(ctx, uint(-delta))
	} else if delta > 0 {
		s.removeSomeUEs(ctx, delta)
	}
	s.UpdateMaxUEsPerCell(ctx)
}

func (s *store) Len(ctx context.Context) int {
	return len(s.ues)
}

func (s *store) LenPerCell(ctx context.Context, cellECGI uint64) int {
	result := 0
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ue := range s.ues {
		if uint64(ue.Cell.ECGI) == cellECGI {
			result++
		}
	}
	return result
}

func (s *store) MaxUEsPerCell(ctx context.Context, cellECGI uint64) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.maxUEs[cellECGI]
	if !ok {
		return 0
	}
	return result
}

func (s *store) SetMaxUEsPerCell(ctx context.Context, cellECGI uint64, maxNumUEs int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxUEs[cellECGI] = maxNumUEs
}

func (s *store) UpdateMaxUEsPerCell(ctx context.Context) {
	cNumUEsMap := make(map[uint64]int)
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ue := range s.ues {
		if _, ok := s.maxUEs[uint64(ue.Cell.ECGI)]; !ok {
			cNumUEsMap[uint64(ue.Cell.ECGI)] = 1
			continue
		}
		cNumUEsMap[uint64(ue.Cell.ECGI)]++
	}

	log.Debugf("[before] cNumUEsMap: %v", cNumUEsMap)
	log.Debugf("[before] maxUEs: %v", s.maxUEs)

	// compare
	for k, v := range cNumUEsMap {
		oNumUEs, ok := s.maxUEs[k]
		if !ok || v > oNumUEs {
			s.maxUEs[k] = v
			continue
		}
	}

	log.Debugf("[after] maxUEs: %v", s.maxUEs)
}

func (s *store) removeSomeUEs(ctx context.Context, count int) {
	c := count
	for imsi := range s.ues {
		if c == 0 {
			break
		}
		_, _ = s.Delete(ctx, imsi)
		c = c - 1
	}
}

func (s *store) CreateUEs(ctx context.Context, count uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := uint(0); i < count; i++ {
		imsi := types.IMSI(rand.Int63n(maxIMSI-minIMSI) + minIMSI)
		if _, ok := s.ues[imsi]; ok {
			// FIXME: more robust check for duplicates
			imsi = types.IMSI(rand.Int63n(maxIMSI-minIMSI) + minIMSI)
		}

		randomCell, err := s.cellStore.GetRandomCell()
		if err != nil {
			log.Error(err)
		}
		ecgi := randomCell.ECGI
		ue := &model.UE{
			IMSI:     imsi,
			Type:     "phone",
			Location: model.Coordinate{Lat: 0, Lng: 0},
			Heading:  0,
			Cell: &model.UECell{
				ID:       types.GnbID(ecgi), // placeholder
				ECGI:     ecgi,
				Strength: rand.Float64() * 100,
			},
			CRNTI:      types.CRNTI(90125 + i),
			Cells:      nil,
			IsAdmitted: false,
		}
		s.ues[ue.IMSI] = ue
	}
}

// Get gets a UE based on a given imsi
func (s *store) Get(ctx context.Context, imsi types.IMSI) (*model.UE, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if node, ok := s.ues[imsi]; ok {
		return node, nil
	}

	return nil, errors.New(errors.NotFound, "UE not found")
}

// Delete deletes a UE based on a given imsi
func (s *store) Delete(ctx context.Context, imsi types.IMSI) (*model.UE, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ue, ok := s.ues[imsi]; ok {
		delete(s.ues, imsi)
		deleteEvent := event.Event{
			Key:   imsi,
			Value: ue,
			Type:  Deleted,
		}
		s.watchers.Send(deleteEvent)
		return ue, nil
	}
	return nil, errors.New(errors.NotFound, "UE not found")
}

func (s *store) ListAllUEs(ctx context.Context) []*model.UE {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.UE, 0, len(s.ues))
	for _, ue := range s.ues {
		list = append(list, ue)
	}
	return list
}

func (s *store) MoveToCell(ctx context.Context, imsi types.IMSI, ecgi types.ECGI, strength float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ue, ok := s.ues[imsi]; ok {
		ue.Cell.ECGI = ecgi
		ue.Cell.Strength = strength
		updateEvent := event.Event{
			Key:   ue.IMSI,
			Value: ue,
			Type:  Updated,
		}
		s.watchers.Send(updateEvent)
		return nil
	}
	return errors.New(errors.NotFound, "UE not found")
}

func (s *store) MoveToCoordinate(ctx context.Context, imsi types.IMSI, location model.Coordinate, heading uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ue, ok := s.ues[imsi]; ok {
		ue.Location = location
		ue.Heading = heading
		updateEvent := event.Event{
			Key:   ue.IMSI,
			Value: ue,
			Type:  Updated,
		}
		s.watchers.Send(updateEvent)
		return nil
	}
	return errors.New(errors.NotFound, "UE not found")
}

func (s *store) UpdateCells(ctx context.Context, imsi types.IMSI, cells []*model.UECell) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ue, ok := s.ues[imsi]; ok {
		ue.Cells = cells
		updateEvent := event.Event{
			Key:   ue.IMSI,
			Value: ue,
			Type:  Updated,
		}
		s.watchers.Send(updateEvent)
		return nil
	}
	return errors.New(errors.NotFound, "UE not found")
}

func (s *store) UpdateCell(ctx context.Context, imsi types.IMSI, cell *model.UECell) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ue, ok := s.ues[imsi]; ok {
		ue.Cell = cell
		updateEvent := event.Event{
			Key:   ue.IMSI,
			Value: ue,
			Type:  Updated,
		}
		s.watchers.Send(updateEvent)
		return nil
	}

	return errors.New(errors.NotFound, "UE not found")
}

func (s *store) ListUEs(ctx context.Context, ecgi types.ECGI) []*model.UE {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.UE, 0, len(s.ues))
	for _, ue := range s.ues {
		if ue.Cell.ECGI == ecgi {
			list = append(list, ue)
		}
	}
	return list
}

func (s *store) Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error {
	log.Debug("Watching ue changes")
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
			for _, ue := range s.ues {
				ch <- event.Event{
					Key:   ue.IMSI,
					Value: ue,
					Type:  None,
				}
			}
		}()
	}

	return nil
}
