// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cellglobalid

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
)

// GlobalNRCGIID cell global NRCGI ID
type GlobalNRCGIID struct {
	plmnID   *ransimtypes.Uint24
	nrCellID *e2smkpmv2.BitString
}

// NewGlobalNRCGIID creates new global NRCGI ID
func NewGlobalNRCGIID(options ...func(*GlobalNRCGIID)) *GlobalNRCGIID {
	nrcgiid := &GlobalNRCGIID{}
	for _, option := range options {
		option(nrcgiid)
	}

	return nrcgiid
}

// WithPlmnID sets plmn ID
func WithPlmnID(plmnID *ransimtypes.Uint24) func(nrcgiid *GlobalNRCGIID) {
	return func(nrcgid *GlobalNRCGIID) {
		nrcgid.plmnID = plmnID

	}
}

// WithNRCellID sets NRCellID
func WithNRCellID(nrCellID *e2smkpmv2.BitString) func(nrcgiid *GlobalNRCGIID) {
	return func(nrcgid *GlobalNRCGIID) {
		nrcgid.nrCellID = nrCellID
	}
}

// Build builds a global NRCGI ID
func (gNRCGIID *GlobalNRCGIID) Build() (*e2smkpmv2.CellGlobalId, error) {
	return &e2smkpmv2.CellGlobalId{
		CellGlobalId: &e2smkpmv2.CellGlobalId_NrCgi{
			NrCgi: &e2smkpmv2.Nrcgi{
				PLmnIdentity: &e2smkpmv2.PlmnIdentity{
					Value: gNRCGIID.plmnID.ToBytes(),
				},
				NRcellIdentity: &e2smkpmv2.NrcellIdentity{
					Value: gNRCGIID.nrCellID,
				},
			},
		},
	}, nil
}
