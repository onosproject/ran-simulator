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

func (d *driver) updateRrc(ctx context.Context, imsi types.IMSI) error {
	var err error
	var rrcStateChanged bool

	if rand.Float64() < RrcStateChangeProbability {
		ue, err := d.ueStore.Get(ctx, imsi)
		if err != nil {
			return err
		}

		if ue.RrcState == e2sm_mho.Rrcstatus_RRCSTATUS_IDLE {
			rrcStateChanged, err = d.rrcConnected(ctx, imsi, RrcStateChangeVariance)
		} else if ue.RrcState == e2sm_mho.Rrcstatus_RRCSTATUS_CONNECTED {
			rrcStateChanged, err = d.rrcIdle(ctx, imsi, RrcStateChangeVariance)
		} else { // Ignore e2sm_mho.Rrcstatus_RRCSTATUS_INACTIVE
			return nil
		}

		if err == nil && d.hoLogic != "local" && rrcStateChanged {
			select {
			case d.rrcCtrl.RrcUpdateChan <- *ue:
				// TODO - increment counters
			default:
				log.Debugf("RRC state changed but not reported imsi:%v state:%v", imsi, ue.RrcState)
			}
		}

	}

	return err

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
		if err = ueDetach(ctx, ue, d.ueStore, d.cellStore); err == nil {
			log.Debugf("RRC state change imsi:%d from CONNECTED to IDLE", imsi)
			ue.RrcState = e2sm_mho.Rrcstatus_RRCSTATUS_IDLE
			d.cellStore.IncrementRrcIdleCount(ctx, ue.Cell.NCGI)
			d.cellStore.DecrementRrcConnectedCount(ctx, ue.Cell.NCGI)
		}
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
		if err = ueAttach(ctx, ue, d.ueStore, d.cellStore); err == nil {
			ue.RrcState = e2sm_mho.Rrcstatus_RRCSTATUS_CONNECTED
			d.cellStore.IncrementRrcConnectedCount(ctx, ue.Cell.NCGI)
			d.cellStore.DecrementRrcIdleCount(ctx, ue.Cell.NCGI)
		}
	}

	return rrcStateChanged, err

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

	return ueStore.UpdateCell(ctx, ue.IMSI, newUECell)

}

func ueDetach(ctx context.Context, ue *model.UE, ueStore ues.Store, cellStore cells.Store) error {
	// TODO
	return nil
}
