// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package nrt

import (
	"fmt"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v2/e2sm-rc-pre-v2"
)

// Neighbour neighbour fields for nrt message
type Neighbour struct {
	plmnID            ransimtypes.Uint24
	eutraCellIdentity uint64
	earfcn            int32
	pci               int32
	cellSize          e2smrcpreies.CellSize
}

// NewNeighbour creates a new neighbour message
func NewNeighbour(options ...func(nrb *Neighbour)) *Neighbour {
	nrb := &Neighbour{}
	for _, option := range options {
		option(nrb)
	}

	return nrb
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID ransimtypes.Uint24) func(neighbour *Neighbour) {
	return func(neighbour *Neighbour) {
		neighbour.plmnID = plmnID
	}
}

// WithEutraCellIdentity sets eutraCellIdentity
func WithEutraCellIdentity(eutraCellIdentity uint64) func(neighbour *Neighbour) {
	return func(neighbour *Neighbour) {
		neighbour.eutraCellIdentity = eutraCellIdentity
	}
}

// WithEarfcn sets earfcn
func WithEarfcn(earfcn int32) func(neighbour *Neighbour) {
	return func(neighbour *Neighbour) {
		neighbour.earfcn = earfcn
	}
}

// WithPci sets pci
func WithPci(pci int32) func(neighbour *Neighbour) {
	return func(neighbour *Neighbour) {
		neighbour.pci = pci
	}
}

// WithCellSize sets cell size
func WithCellSize(cellSize e2smrcpreies.CellSize) func(neighbour *Neighbour) {
	return func(neighbour *Neighbour) {
		neighbour.cellSize = cellSize
	}
}

// Build builds Nrt message for RC service model
func (neighbour *Neighbour) Build() (*e2smrcpreies.Nrt, error) {
	nrtMsg := &e2smrcpreies.Nrt{
		Cgi: &e2smrcpreies.CellGlobalId{
			CellGlobalId: &e2smrcpreies.CellGlobalId_EUtraCgi{
				EUtraCgi: &e2smrcpreies.Eutracgi{
					PLmnIdentity: &e2smrcpreies.PlmnIdentity{
						Value: neighbour.plmnID.ToBytes(),
					},
					EUtracellIdentity: &e2smrcpreies.EutracellIdentity{
						Value: &e2smrcpreies.BitString{
							Value: neighbour.eutraCellIdentity, //uint64
							Len:   28,                          //uint32
						},
					},
				},
			},
		},
		Pci: &e2smrcpreies.Pci{
			Value: neighbour.pci,
		},
		CellSize: neighbour.cellSize,
		DlArfcn: &e2smrcpreies.Arfcn{
			Arfcn: &e2smrcpreies.Arfcn_EArfcn{
				EArfcn: &e2smrcpreies.Earfcn{
					Value: neighbour.earfcn,
				},
			},
		},
	}
	if err := nrtMsg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return nrtMsg, nil
}
