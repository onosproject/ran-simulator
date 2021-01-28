// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import "sync"

// Imsi represents a unique UE identifier
type Imsi uint64

// UEType represents type of user-equipment
type UEType string

// Crnti is a tower-specific UE identifier
type Crnti string

// Coordinate represents a geographical location
type Coordinate struct {
	Lat float64
	Lng float64
}

// UETower represents UE-tower relationship
type UETower struct {
	Ecgi     Ecgi
	Strength float64
}

// UE represents user-equipment, i.e. phone, IoT device, etc.
type UE struct {
	Imsi     Imsi
	Type     UEType
	Location Coordinate
	Rotation uint32

	Tower  *UETower
	Crnti  Crnti
	Towers []*UETower

	IsAdmitted bool
	// Metrics
}

// UERegistry tracks inventory of user-equipment for the simulation
type UERegistry interface {
	// CreateUEs creates the specified number of UEs
	CreateUEs(count uint)

	// DestroyUE destroy the specified UE
	DestroyUE(imsi Imsi)

	// ListAllUEs returns an array of all UEs
	ListAllUEs() []*UE

	// MoveUE update the tower affiliation of the specified UE
	MoveUE(imsi Imsi, ecgi Ecgi, strength float64)

	// ListUEs returns an array of all UEs associated with the specified tower
	ListUEs(ecgi Ecgi) []*UE
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

func (r *registry) CreateUEs(count uint) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i := uint(0); i < count; i++ {

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

func (r *registry) MoveUE(imsi Imsi, ecgi Ecgi, strength float64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	ue := r.ues[imsi]
	if ue != nil {
		ue.Tower.Ecgi = ecgi
		ue.Tower.Strength = strength
	}
}

func (r *registry) ListUEs(ecgi Ecgi) []*UE {
	r.lock.RLock()
	defer r.lock.RUnlock()
	list := make([]*UE, 0, len(r.ues))
	for _, ue := range r.ues {
		if ue.Tower.Ecgi == ecgi {
			list = append(list, ue)
		}
	}
	return list
}
