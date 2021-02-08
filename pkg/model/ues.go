// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"math/rand"
	"sync"
)

const (
	minImsi = 10000
	maxImsi = 99999
)

// UEType represents type of user-equipment
type UEType string

// UECell represents UE-cell relationship
type UECell struct {
	ID       GEnbID
	Ecgi     Ecgi // Auxiliary form of association
	Strength float64
}

// UE represents user-equipment, i.e. phone, IoT device, etc.
type UE struct {
	Imsi     Imsi
	Type     UEType
	Location Coordinate
	Rotation uint32

	Cell  *UECell
	Crnti Crnti
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
	DestroyUE(imsi Imsi)

	// ListAllUEs returns an array of all UEs
	ListAllUEs() []*UE

	// MoveUE update the cell affiliation of the specified UE
	MoveUE(imsi Imsi, genbID GEnbID, strength float64)

	// ListUEs returns an array of all UEs associated with the specified cell
	ListUEs(genbID GEnbID) []*UE

	// GetNumUes returns number of active UEs
	GetNumUes() int
}

type registry struct {
	lock sync.RWMutex
	ues  map[Imsi]*UE
}

// NewUERegistry creates a new user-equipment registry primed with the specified number of UEs to start
func NewUERegistry(count uint) UERegistry {
	reg := &registry{
		lock: sync.RWMutex{},
		ues:  make(map[Imsi]*UE),
	}
	reg.CreateUEs(count)
	return reg
}

func (r *registry) SetUECount(count uint) {
	delta := len(r.ues) - int(count)
	if delta < 0 {
		r.CreateUEs(uint(-delta))
	} else if delta > 0 {
		r.removeSomeUEs(delta)
	}
}

func (r *registry) removeSomeUEs(count int) {
	c := count
	for imsi := range r.ues {
		if c == 0 {
			break
		}
		r.DestroyUE(imsi)
		c = c - 1
	}
}

func (r *registry) CreateUEs(count uint) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i := uint(0); i < count; i++ {
		// FIXME: fill in with more sensible values
		ue := &UE{
			Imsi:       Imsi(rand.Int63n(maxImsi-minImsi) + minImsi),
			Type:       "phone",
			Location:   Coordinate{0, 0},
			Rotation:   0,
			Cell:       &UECell{},
			Crnti:      "90125",
			Cells:      nil,
			IsAdmitted: false,
		}
		r.ues[ue.Imsi] = ue
	}
}

func (r *registry) DestroyUE(imsi Imsi) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.ues, imsi)
}

func (r *registry) ListAllUEs() []*UE {
	r.lock.RLock()
	defer r.lock.RUnlock()
	list := make([]*UE, 0, len(r.ues))
	for _, ue := range r.ues {
		list = append(list, ue)
	}
	return list
}

func (r *registry) MoveUE(imsi Imsi, genbID GEnbID, strength float64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	ue := r.ues[imsi]
	if ue != nil {
		ue.Cell.ID = genbID
		ue.Cell.Strength = strength
	}
}

func (r *registry) ListUEs(genbID GEnbID) []*UE {
	r.lock.RLock()
	defer r.lock.RUnlock()
	list := make([]*UE, 0, len(r.ues))
	for _, ue := range r.ues {
		if ue.Cell.ID.EnbID == genbID.EnbID && ue.Cell.ID.PlmnID == genbID.PlmnID {
			list = append(list, ue)
		}
	}
	return list
}

func (r *registry) GetNumUes() int {
	return len(r.ues)
}
