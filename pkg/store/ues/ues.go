// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ues

import (
	"context"
	"fmt"
	e2smcommonies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-common-ies"
	"math/rand"
	"sync"

	mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"

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

var log = liblog.GetLogger()

// Store tracks inventory of user-equipment for the simulation
type Store interface {
	// SetUECount updates the UE count and creates or deletes new UEs as needed
	SetUECount(ctx context.Context, count uint)

	// Len returns the number of active UEs
	Len(ctx context.Context) int

	// LenPerCell returns the number of active UEs per cell
	LenPerCell(ctx context.Context, cellNCGI uint64) int

	// MaxUEsPerCell returns the maximum number of active UEs per cell
	MaxUEsPerCell(ctx context.Context, cellNCGI uint64) int

	// SetMaxUEsPerCell sets the maximum number of active UEs per cell
	SetMaxUEsPerCell(ctx context.Context, cellNCGI uint64, maxNumUEs int)

	// UpdateMaxUEsPerCell updates the maximum number of active UEs for all cells
	UpdateMaxUEsPerCell(ctx context.Context)

	// CreateUEs creates the specified number of UEs
	CreateUEs(ctx context.Context, count uint)

	// Get retrieves the UE with the specified IMSI
	Get(ctx context.Context, imsi types.IMSI) (*model.UE, error)

	// GetWithGNbUeID retrieves the UE with the gNB UE ID
	GetWithGNbUeID(ctx context.Context, gNBUeID *e2smcommonies.UeidGnb) (*model.UE, error)

	// Delete destroy the specified UE
	Delete(ctx context.Context, imsi types.IMSI) (*model.UE, error)

	// MoveToCell update the cell affiliation of the specified UE
	MoveToCell(ctx context.Context, imsi types.IMSI, ncgi types.NCGI, strength float64) error

	// MoveToCoordinate updates the UEs geo location and compass heading
	MoveToCoordinate(ctx context.Context, imsi types.IMSI, location model.Coordinate, heading uint32) error

	// UpdateUE updates the 5QI value of the UE
	UpdateUE(ctx context.Context, imsi types.IMSI, fiveQi int, isChanged bool) error

	// UpdateCells updates the visible cells and their signal strength
	UpdateCells(ctx context.Context, imsi types.IMSI, cells []*model.UECell) error

	// UpdateCell updates the serving cell
	UpdateCell(ctx context.Context, imsi types.IMSI, cell *model.UECell) error

	// ListAllUEs returns an array of all UEs
	ListAllUEs(ctx context.Context) []*model.UE

	// ListUEs returns an array of all UEs associated with the specified cell
	ListUEs(ctx context.Context, ncgi types.NCGI) []*model.UE

	// Watch watches the UE inventory events using the supplied channel
	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error
}

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

type store struct {
	mu              sync.RWMutex
	ues             map[types.IMSI]*model.UE
	maxUEs          map[uint64]int
	cellStore       cells.Store
	watchers        *watcher.Watchers
	initialRrcState string
}

// NewUERegistry creates a new user-equipment registry primed with the specified number of UEs to start.
// UEs will be semi-randomly distributed between the specified cells
func NewUERegistry(count uint, cellStore cells.Store, initialRrcState string) Store {
	log.Infof("Creating registry from model with %d UEs", count)
	watchers := watcher.NewWatchers()
	store := &store{
		mu:              sync.RWMutex{},
		ues:             make(map[types.IMSI]*model.UE),
		maxUEs:          make(map[uint64]int),
		cellStore:       cellStore,
		watchers:        watchers,
		initialRrcState: initialRrcState,
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
}

func (s *store) Len(ctx context.Context) int {
	return len(s.ues)
}

func (s *store) LenPerCell(ctx context.Context, cellNCGI uint64) int {
	result := 0
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ue := range s.ues {
		if uint64(ue.Cell.NCGI) == cellNCGI {
			result++
		}
	}
	return result
}

func (s *store) MaxUEsPerCell(ctx context.Context, cellNCGI uint64) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.maxUEs[cellNCGI]
	if !ok {
		return 0
	}
	return result
}

func (s *store) SetMaxUEsPerCell(ctx context.Context, cellNCGI uint64, maxNumUEs int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxUEs[cellNCGI] = maxNumUEs
}

func (s *store) UpdateMaxUEsPerCell(ctx context.Context) {
	cNumUEsMap := make(map[uint64]int)
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ue := range s.ues {
		if _, ok := s.maxUEs[uint64(ue.Cell.NCGI)]; !ok {
			cNumUEsMap[uint64(ue.Cell.NCGI)] = 1
			continue
		}
		cNumUEsMap[uint64(ue.Cell.NCGI)]++
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

func randomBoolean() bool {
	return rand.Float32() < 0.5
}

func (s *store) CreateUEs(ctx context.Context, count uint) {
	s.mu.Lock()
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
		ncgi := randomCell.NCGI
		var rrcState mho.Rrcstatus
		if s.initialRrcState == "connected" || s.initialRrcState == "idle" {
			if s.initialRrcState == "idle" {
				rrcState = mho.Rrcstatus_RRCSTATUS_IDLE
				s.cellStore.IncrementRrcIdleCount(ctx, ncgi)
			} else {
				rrcState = mho.Rrcstatus_RRCSTATUS_CONNECTED
				s.cellStore.IncrementRrcConnectedCount(ctx, ncgi)
			}
		} else {
			if randomBoolean() {
				rrcState = mho.Rrcstatus_RRCSTATUS_IDLE
				s.cellStore.IncrementRrcIdleCount(ctx, ncgi)
			} else {
				rrcState = mho.Rrcstatus_RRCSTATUS_CONNECTED
				s.cellStore.IncrementRrcConnectedCount(ctx, ncgi)
			}
		}
		ue := &model.UE{
			IMSI:        imsi,
			AmfUeNgapID: types.AmfUENgapID(i + 1000),
			Type:        "phone",
			Location:    model.Coordinate{Lat: 0, Lng: 0},
			Heading:     0,
			Cell: &model.UECell{
				ID:       types.GnbID(ncgi), // placeholder
				NCGI:     ncgi,
				Strength: rand.Float64() * 100,
			},
			CRNTI:      types.CRNTI(90125 + i),
			Cells:      nil,
			IsAdmitted: false,
			RrcState:   rrcState,
		}
		s.ues[ue.IMSI] = ue
	}
	s.mu.Unlock()
	s.UpdateMaxUEsPerCell(ctx)
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

func (s *store) GetWithGNbUeID(ctx context.Context, gNBUeID *e2smcommonies.UeidGnb) (*model.UE, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	amfUeNgapID := gNBUeID.AmfUeNgapId.GetValue()
	for _, ue := range s.ues {
		// TODO add GUAMI - currently RAN simulator only supports single AMF, it should be fine
		// TODO for the future, GUAMI should be considered here
		if int64(ue.AmfUeNgapID) == amfUeNgapID {
			return ue, nil
		}
	}
	return nil, errors.NewNotFound(fmt.Sprintf("the UE having gNB UE ID %v Not found", gNBUeID))
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

func (s *store) MoveToCell(ctx context.Context, imsi types.IMSI, ncgi types.NCGI, strength float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ue, ok := s.ues[imsi]; ok {
		ue.Cell.NCGI = ncgi
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

func (s *store) UpdateUE(ctx context.Context, imsi types.IMSI, fiveQi int, isChanged bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ue, ok := s.ues[imsi]; ok {
		ue.FiveQi = fiveQi
		updateEvent := event.Event{
			Key:   ue.IMSI,
			Value: ue,
			Type:  Updated,
		}
		log.Debugf("Updating UE %v with 5QI value %v", ue.IMSI, ue.FiveQi)
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

func (s *store) ListUEs(ctx context.Context, ncgi types.NCGI) []*model.UE {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.UE, 0, len(s.ues))
	for _, ue := range s.ues {
		if ue.Cell.NCGI == ncgi {
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
