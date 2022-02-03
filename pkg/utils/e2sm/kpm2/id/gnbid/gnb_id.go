// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package gnbid

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
)

// GlobalGNBID global gNB ID
type GlobalGNBID struct {
	plmnID      ransimtypes.Uint24
	gNBIDChoice *asn1.BitString
	gNBCuUpID   int64
	gNBDuID     int64
}

// NewGlobalGNBID creates new global gnb ID
func NewGlobalGNBID(options ...func(*GlobalGNBID)) *GlobalGNBID {
	gNBID := &GlobalGNBID{}
	for _, option := range options {
		option(gNBID)
	}

	return gNBID
}

// WithPlmnID sets plmn ID
func WithPlmnID(plmnID ransimtypes.Uint24) func(gNBID *GlobalGNBID) {
	return func(gNBID *GlobalGNBID) {
		gNBID.plmnID = plmnID

	}
}

// WithGNBIDChoice sets gNBID choice
func WithGNBIDChoice(gnbIDChoice *asn1.BitString) func(gNBID *GlobalGNBID) {
	return func(gNBID *GlobalGNBID) {
		gNBID.gNBIDChoice = gnbIDChoice
	}
}

// WithGNBCuUpID sets gNB CuUp ID
func WithGNBCuUpID(gNBCuUpID int64) func(gNBID *GlobalGNBID) {
	return func(gNBID *GlobalGNBID) {
		gNBID.gNBCuUpID = gNBCuUpID
	}
}

// WithGNBDuID sets gNB DuID
func WithGNBDuID(gNBDuID int64) func(gNBID *GlobalGNBID) {
	return func(gNBID *GlobalGNBID) {
		gNBID.gNBDuID = gNBDuID
	}
}

// Build builds a global gNB ID
func (gNBID *GlobalGNBID) Build() (*e2smkpmv2.GlobalKpmnodeId, error) {
	return &e2smkpmv2.GlobalKpmnodeId{
		GlobalKpmnodeId: &e2smkpmv2.GlobalKpmnodeId_GNb{
			GNb: &e2smkpmv2.GlobalKpmnodeGnbId{
				GlobalGNbId: &e2smkpmv2.GlobalgNbId{
					GnbId: &e2smkpmv2.GnbIdChoice{
						GnbIdChoice: &e2smkpmv2.GnbIdChoice_GnbId{
							GnbId: gNBID.gNBIDChoice,
						},
					},
					PlmnId: &e2smkpmv2.PlmnIdentity{
						Value: gNBID.plmnID.ToBytes(),
					},
				},
				GNbCuUpId: &e2smkpmv2.GnbCuUpId{
					Value: gNBID.gNBCuUpID,
				},
				GNbDuId: &e2smkpmv2.GnbDuId{
					Value: gNBID.gNBDuID,
				},
			},
		},
	}, nil
}
