// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mobility

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/v1/e2sm-mho"
	"github.com/onosproject/ran-simulator/pkg/model"
	"math/rand"
)

// RrcStateChangeProbability determines the rate of change of RRC states in ransim
var RrcStateChangeProbability float64 = 0.2

// RrcStateChangeVariance provides non-determinism in enforcing the UeCountPerCell
var RrcStateChangeVariance float64 = 0.9

// UeCountPerCellDefault is the default number of RRC Connected UEs per cell
var UeCountPerCellDefault uint = 15

// RrcCtrl is the RRC controller
type RrcCtrl struct {
	RrcUpdateChan  chan model.UE
	ueCountPerCell uint
}

// NewRrcCtrl returns a new RRC Controller
func NewRrcCtrl(ueCountPerCell uint) RrcCtrl {
	if ueCountPerCell == 0 {
		ueCountPerCell = UeCountPerCellDefault
	}
	return RrcCtrl{
		RrcUpdateChan:  make(chan model.UE),
		ueCountPerCell: ueCountPerCell,
	}
}

func (d *driver) totalUeCount(ctx context.Context, ncgi types.NCGI) uint {
	cell, err := d.cellStore.Get(ctx, ncgi)
	if err != nil {
		log.Error(err)
		return 0
	}
	return uint(cell.RrcConnectedCount + cell.RrcIdleCount)
}

//func (d *driver) idleUeCount(ctx context.Context, ncgi types.NCGI) uint {
//	cell, err := d.cellStore.Get(ctx, ncgi)
//	if err != nil {
//		log.Error(err)
//		return 0
//	}
//	return uint(cell.RrcIdleCount)
//}

func (d *driver) connectedUeCount(ctx context.Context, ncgi types.NCGI) uint {
	cell, err := d.cellStore.Get(ctx, ncgi)
	if err != nil {
		log.Error(err)
		return 0
	}
	return uint(cell.RrcConnectedCount)
}

func (d *driver) updateRrc(ctx context.Context, imsi types.IMSI) {
	var rrcStateChanged bool

	if rand.Float64() < RrcStateChangeProbability {
		ue, err := d.ueStore.Get(ctx, imsi)
		if err != nil {
			log.Error(err)
			return
		}

		if ue.RrcState == e2sm_mho.Rrcstatus_RRCSTATUS_IDLE {
			rrcStateChanged, err = d.rrcConnected(ctx, imsi, RrcStateChangeVariance)
		} else if ue.RrcState == e2sm_mho.Rrcstatus_RRCSTATUS_CONNECTED {
			rrcStateChanged, err = d.rrcIdle(ctx, imsi, RrcStateChangeVariance)
		} else { // Ignore e2sm_mho.Rrcstatus_RRCSTATUS_INACTIVE
			return
		}

		if err == nil && d.hoLogic != "local" && rrcStateChanged {
			// TODO - check subscription for RRC state changes
			d.rrcCtrl.RrcUpdateChan <- *ue
		}

	}

}

func (d *driver) rrcIdle(ctx context.Context, imsi types.IMSI, p float64) (bool, error) {
	var rrcStateChanged bool = false

	ue, err := d.ueStore.Get(ctx, imsi)
	if err != nil {
		return false, err
	}

	if d.totalUeCount(ctx, ue.Cell.NCGI) > d.rrcCtrl.ueCountPerCell {
		r := rand.Float64()
		if d.connectedUeCount(ctx, ue.Cell.NCGI) > d.rrcCtrl.ueCountPerCell {
			if r < p {
				rrcStateChanged = true
			}
		} else {
			if r < 1-p {
				rrcStateChanged = true
			}
		}
	} else {
		rrcStateChanged = true
	}

	if rrcStateChanged {
		log.Debugf("RRC state change imsi:%d from CONNECTED to IDLE", imsi)
		ue.RrcState = e2sm_mho.Rrcstatus_RRCSTATUS_IDLE
		d.cellStore.IncrementRrcIdleCount(ctx, ue.Cell.NCGI)
		d.cellStore.DecrementRrcConnectedCount(ctx, ue.Cell.NCGI)
	}

	return rrcStateChanged, err

}

func (d *driver) rrcConnected(ctx context.Context, imsi types.IMSI, p float64) (bool, error) {
	var rrcStateChanged bool = false

	ue, err := d.ueStore.Get(ctx, imsi)
	if err != nil {
		return false, err
	}

	if d.totalUeCount(ctx, ue.Cell.NCGI) > d.rrcCtrl.ueCountPerCell {
		r := rand.Float64()
		if d.connectedUeCount(ctx, ue.Cell.NCGI) > d.rrcCtrl.ueCountPerCell {
			if r < 1-p {
				rrcStateChanged = true
			}
		} else {
			if r < p {
				rrcStateChanged = true
			}
		}
	} else {
		rrcStateChanged = true
	}

	if rrcStateChanged {
		log.Debugf("RRC state change imsi:%d from IDLE to CONNECTED", imsi)
		ue.RrcState = e2sm_mho.Rrcstatus_RRCSTATUS_CONNECTED
		d.cellStore.IncrementRrcConnectedCount(ctx, ue.Cell.NCGI)
		d.cellStore.DecrementRrcIdleCount(ctx, ue.Cell.NCGI)
	}

	return rrcStateChanged, err

}
