// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/onosproject/ran-simulator/api/types"
	"math/rand"
	"sync"
)

const (
	minIMSI = 10000
	maxIMSI = 99999
)

// UEType represents type of user-equipment
type UEType string

// UECell represents UE-cell relationship
type UECell struct {
	ID       types.GEnbID
	Ecgi     types.ECGI // Auxiliary form of association
	Strength float64
}

// UE represents user-equipment, i.e. phone, IoT device, etc.
type UE struct {
	IMSI     types.IMSI
	Type     UEType
	Location Coordinate
	Rotation uint32

	Cell  *UECell
	CRNTI types.CRNTI
	Cells []*UECell

	IsAdmitted bool
	// Metrics
}

// UERegistry tracks inventory of user-equipment for the simulation
type UERegistry interface {
	// SetUECount updates the UE count and creates or deletes new UEs as needed
	SetUECount(count uint)

	// CreateUEs creates the specified number of UEs
	CreateUEs(count uint)

	// DestroyUE destroy the specified UE
	DestroyUE(IMSI types.IMSI)

	// ListAllUEs returns an array of all UEs
	ListAllUEs() []*UE

	// MoveUE update the cell affiliation of the specified UE
	MoveUE(IMSI types.IMSI, genbID types.GEnbID, strength float64)

	// ListUEs returns an array of all UEs associated with the specified cell
	ListUEs(genbID types.GEnbID) []*UE

	// GetNumUes returns number of active UEs
	GetNumUes() int
}

type ueRegistry struct {
	lock sync.RWMutex
	ues  map[types.IMSI]*UE
}

// NewUERegistry creates a new user-equipment registry primed with the specified number of UEs to start
func NewUERegistry(count uint) UERegistry {
	reg := &ueRegistry{
		lock: sync.RWMutex{},
		ues:  make(map[types.IMSI]*UE),
	}
	reg.CreateUEs(count)
	return reg
}

func (r *ueRegistry) SetUECount(count uint) {
	delta := len(r.ues) - int(count)
	if delta < 0 {
		r.CreateUEs(uint(-delta))
	} else if delta > 0 {
		r.removeSomeUEs(delta)
	}
}

func (r *ueRegistry) removeSomeUEs(count int) {
	c := count
	for IMSI := range r.ues {
		if c == 0 {
			break
		}
		r.DestroyUE(IMSI)
		c = c - 1
	}
}

func (r *ueRegistry) CreateUEs(count uint) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i := uint(0); i < count; i++ {
		// FIXME: fill in with more sensible values
		ue := &UE{
			IMSI:       types.IMSI(rand.Int63n(maxIMSI-minIMSI) + minIMSI),
			Type:       "phone",
			Location:   Coordinate{0, 0},
			Rotation:   0,
			Cell:       &UECell{},
			CRNTI:      types.CRNTI(90125 + i),
			Cells:      nil,
			IsAdmitted: false,
		}
		r.ues[ue.IMSI] = ue
	}
}

func (r *ueRegistry) DestroyUE(IMSI types.IMSI) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.ues, IMSI)
}

func (r *ueRegistry) ListAllUEs() []*UE {
	r.lock.RLock()
	defer r.lock.RUnlock()
	list := make([]*UE, 0, len(r.ues))
	for _, ue := range r.ues {
		list = append(list, ue)
	}
	return list
}

func (r *ueRegistry) MoveUE(IMSI types.IMSI, genbID types.GEnbID, strength float64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	ue := r.ues[IMSI]
	if ue != nil {
		ue.Cell.ID = genbID
		ue.Cell.Strength = strength
	}
}

func (r *ueRegistry) ListUEs(genbID types.GEnbID) []*UE {
	r.lock.RLock()
	defer r.lock.RUnlock()
	list := make([]*UE, 0, len(r.ues))
	for _, ue := range r.ues {
		if ue.Cell.ID == genbID {
			list = append(list, ue)
		}
	}
	return list
}

func (r *ueRegistry) GetNumUes() int {
	return len(r.ues)
}
