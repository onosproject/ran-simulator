// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mobility

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/v1/e2sm-mho"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"math"
	"math/rand"
)

// RrcStateChangeProbability determines the rate of change of RRC states in ransim
var RrcStateChangeProbability float64 = 0.02

// RrcCtrl is the RRC controller
type RrcCtrl struct {
	RrcUpdateChan chan model.UE
}

// NewRrcCtrl returns a new RRC Controller
func (d *driver) NewRrcCtrl() *RrcCtrl {
	return &RrcCtrl{}
}

func (d *driver) updateRrc(ctx context.Context, imsi types.IMSI) error {
	var randomBoolean = rand.Float64() < RrcStateChangeProbability

	if randomBoolean {
		ue, err := d.ueStore.Get(ctx, imsi)
		if err != nil {
			return err
		}

		if ue.RrcState == e2sm_mho.Rrcstatus_RRCSTATUS_IDLE {
			err = d.rrcConnected(ctx, imsi)
		} else if ue.RrcState == e2sm_mho.Rrcstatus_RRCSTATUS_CONNECTED {
			err = d.rrcIdle(ctx, imsi)
		} else { // Ignore e2sm_mho.Rrcstatus_RRCSTATUS_INACTIVE
			return nil
		}
		if err != nil {
			return err
		}

		if d.hoLogic != "local" {
			d.rrcCtrl.RrcUpdateChan <- *ue
		}

	}

	return nil

}

func (d *driver) rrcIdle(ctx context.Context, imsi types.IMSI) error {
	ue, err := d.ueStore.Get(ctx, imsi)
	if err != nil {
		return err
	}
	log.Debugf("RRC state change imsi:%d from CONNECTED to IDLE", imsi)
	ue.RrcState = e2sm_mho.Rrcstatus_RRCSTATUS_IDLE

	//Detach UE
	return ueDetach(ctx, ue, d.ueStore, d.cellStore)
}

func (d *driver) rrcConnected(ctx context.Context, imsi types.IMSI) error {
	ue, err := d.ueStore.Get(ctx, imsi)
	if err != nil {
		return err
	}
	log.Debugf("RRC state change imsi:%d from IDLE to CONNECTED", imsi)
	ue.RrcState = e2sm_mho.Rrcstatus_RRCSTATUS_CONNECTED

	// Attach UE to nearest cell
	return ueAttach(ctx, ue, d.ueStore, d.cellStore)
}

func ueAttach(ctx context.Context, ue *model.UE, ueStore ues.Store, cellStore cells.Store) error {
	cellList, err := cellStore.List(ctx)
	if err != nil {
		return err
	}
	var servCell *model.Cell
	maxRsrp := -math.MaxFloat64
	for _, cell := range cellList {
		rsrp := StrengthAtLocation(ue.Location, *cell)
		if math.IsNaN(rsrp) {
			continue
		}
		if rsrp > maxRsrp {
			servCell = cell
			maxRsrp = rsrp
		}
	}
	newUECell := &model.UECell{
		ID:       types.GnbID(servCell.NCGI),
		NCGI:     servCell.NCGI,
		Strength: maxRsrp,
	}

	err = ueStore.UpdateCell(ctx, ue.IMSI, newUECell)
	if err != nil {
		return err
	}

	return nil
}

func ueDetach(ctx context.Context, ue *model.UE, ueStore ues.Store, cellStore cells.Store) error {
	// TODO
	return nil
}
