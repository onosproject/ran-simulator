// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mobility

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"math"
)

// DoHandover handovers ue to target cell
func DoHandover(ctx context.Context, imsi types.IMSI, tCell *model.UECell, ueStore ues.Store, cellStore cells.Store) {
	err := ueStore.UpdateCell(ctx, imsi, tCell)
	if err != nil {
		log.Warn("Unable to update UE %d cell info", imsi)
	}

	// after changing serving cell, calculate channel quality/signal strength again
	UpdateUESignalStrength(ctx, imsi, ueStore, cellStore)

	log.Debugf("HO is done successfully: %v to %v", imsi, tCell)
}

// UpdateUESignalStrength updates UE signal strength
func UpdateUESignalStrength(ctx context.Context, imsi types.IMSI, ueStore ues.Store, cellStore cells.Store) {
	ue, err := ueStore.Get(ctx, imsi)
	if err != nil {
		log.Warn("Unable to find UE %d", imsi)
		return
	}

	// update RSRP from serving cell
	err = UpdateUESignalStrengthServCell(ctx, ue, ueStore, cellStore)
	if err != nil {
		log.Warnf("For UE %v: %v", *ue, err)
		return
	}

	// update RSRP from candidate serving cells
	err = UpdateUESignalStrengthCandServCells(ctx, ue, ueStore, cellStore)
	if err != nil {
		log.Warnf("For UE %v: %v", *ue, err)
		return
	}

	log.Debugf("for UE [%v]: sCell strength - %v, "+
		"csCell1 strength - %v "+
		"csCell2 strength - %v "+
		"csCell3 strength - %v", ue.IMSI, ue.Cell.Strength, ue.Cells[0].Strength,
		ue.Cells[1].Strength, ue.Cells[2].Strength)
}

// UpdateUESignalStrengthCandServCells updates UE signal strength for serving and candidate cells
func UpdateUESignalStrengthCandServCells(ctx context.Context, ue *model.UE, ueStore ues.Store, cellStore cells.Store) error {
	cellList, err := cellStore.List(ctx)
	if err != nil {
		return fmt.Errorf("Unable to get all cells")
	}
	var csCellList []*model.UECell
	for _, cell := range cellList {
		rsrp := StrengthAtLocation(ue.Location, *cell)
		if math.IsNaN(rsrp) {
			continue
		}
		if ue.Cell.ECGI == cell.ECGI {
			continue
		}
		ueCell := &model.UECell{
			ID:       types.GnbID(cell.ECGI),
			ECGI:     cell.ECGI,
			Strength: rsrp,
		}
		csCellList = SortUECells(append(csCellList, ueCell), 3) // hardcoded: to be parameterized for the future
	}
	err = ueStore.UpdateCells(ctx, ue.IMSI, csCellList)
	if err != nil {
		log.Warn("Unable to update UE %d cells info", ue.IMSI)
	}

	return nil
}

// UpdateUESignalStrengthServCell  updates UE signal strength for serving cell
func UpdateUESignalStrengthServCell(ctx context.Context, ue *model.UE, ueStore ues.Store, cellStore cells.Store) error {
	sCell, err := cellStore.Get(ctx, ue.Cell.ECGI)
	if err != nil {
		return fmt.Errorf("Unable to find serving cell %d", ue.Cell.ECGI)
	}

	strength := StrengthAtLocation(ue.Location, *sCell)

	newUECell := &model.UECell{
		ID:       ue.Cell.ID,
		ECGI:     ue.Cell.ECGI,
		Strength: strength,
	}

	err = ueStore.UpdateCell(ctx, ue.IMSI, newUECell)
	if err != nil {
		log.Warn("Unable to update UE %d cell info", ue.IMSI)
	}
	return nil
}

// SortUECells sorts ue cells
func SortUECells(ueCells []*model.UECell, numAdjCells int) []*model.UECell {
	// bubble sort
	for i := 0; i < len(ueCells)-1; i++ {
		for j := 0; j < len(ueCells)-i-1; j++ {
			if ueCells[j].Strength < ueCells[j+1].Strength {
				ueCells[j], ueCells[j+1] = ueCells[j+1], ueCells[j]
			}
		}
	}
	if len(ueCells) >= numAdjCells {
		return ueCells[0:numAdjCells]
	}
	return ueCells
}
