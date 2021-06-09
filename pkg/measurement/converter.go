// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package measurement

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
	"github.com/onosproject/rrm-son-lib/pkg/model/id"
	"github.com/onosproject/rrm-son-lib/pkg/model/measurement"
	meastype "github.com/onosproject/rrm-son-lib/pkg/model/measurement/type"
)

var logConverter = logging.GetLogger("measurement", "converter")

type MeasReportConverter interface {
	// Convert forms measurement report defined in rrm-son-lib from model.ue defined in ransim
	Convert(ctx context.Context, ue *model.UE) device.UE
}

type measReportConverter struct {
	cellStore cells.Store
	ueStore   ues.Store
}

func NewMeasReportConverter(cellStore cells.Store, ueStore ues.Store) MeasReportConverter {
	return &measReportConverter{
		cellStore: cellStore,
		ueStore:   ueStore,
	}
}

func (c *measReportConverter) Convert(ctx context.Context, ue *model.UE) device.UE {
	ueid := id.NewUEID(uint64(ue.IMSI), uint32(ue.CRNTI), uint64(ue.Cell.ECGI))
	sCellInStore, err := c.cellStore.Get(ctx, ue.Cell.ECGI)
	if err != nil {
		logConverter.Errorf("Can't get serving cell from cell store: %v", err)
	}
	sCell := device.NewCell(id.NewECGI(uint64(sCellInStore.ECGI)),
		c.convertA3Offset(sCellInStore.EventA3Params.A3CellOffset),
		c.convertHysteresis(sCellInStore.EventA3Params.A3Hysteresis),
		c.convertQOffset(sCellInStore.EventA3Params.A3CellOffset),
		c.convertQOffset(sCellInStore.EventA3Params.A3FrequencyOffset),
		c.convertTimeToTrigger(sCellInStore.EventA3Params.A3TimeToTrigger))

	var csCells []device.Cell
	measurements := make(map[string]measurement.Measurement)
	sCellMeas := measurement.NewMeasEventA3(id.NewECGI(uint64(sCellInStore.ECGI)), measurement.RSRP(ue.Cell.Strength))
	measurements[sCellMeas.GetCellID().String()] = sCellMeas

	for _, ueCell := range ue.Cells {
		tmpCellInStore, _ := c.cellStore.Get(ctx, ueCell.ECGI)
		if err != nil {
			logConverter.Errorf("Can't get candidate serving cell from cell storeL: %v", err)
		}

		csCells = append(csCells, device.NewCell(id.NewECGI(uint64(tmpCellInStore.ECGI)),
			c.convertA3Offset(tmpCellInStore.EventA3Params.A3CellOffset),
			c.convertHysteresis(tmpCellInStore.EventA3Params.A3Hysteresis),
			c.convertQOffset(tmpCellInStore.EventA3Params.A3CellOffset),
			c.convertQOffset(tmpCellInStore.EventA3Params.A3FrequencyOffset),
			c.convertTimeToTrigger(tmpCellInStore.EventA3Params.A3TimeToTrigger)))

		tmpCsCell := measurement.NewMeasEventA3(id.NewECGI(uint64(tmpCellInStore.ECGI)), measurement.RSRP(ueCell.Strength))
		measurements[tmpCsCell.GetCellID().String()] = tmpCsCell
	}

	report := device.NewUE(ueid, sCell, csCells)
	report.SetMeasurements(measurements)
	return report
}

func (c *measReportConverter) convertA3Offset(a3Offset int32) meastype.A3OffsetRange {
	return meastype.A3OffsetRange(a3Offset)
}

func (c *measReportConverter) convertHysteresis(hyst int32) meastype.HysteresisRange {
	return meastype.HysteresisRange(hyst)
}

func (c *measReportConverter) convertQOffset(qoffset int32) meastype.QOffsetRange {
	if qoffset <= -24 {
		return meastype.QOffsetMinus24dB
	} else if qoffset <= -22 {
		return meastype.QOffsetMinus22dB
	} else if qoffset <= -20 {
		return meastype.QOffsetMinus20dB
	} else if qoffset <= -18 {
		return meastype.QOffsetMinus18dB
	} else if qoffset <= -16 {
		return meastype.QOffsetMinus16dB
	} else if qoffset <= -14 {
		return meastype.QOffsetMinus14dB
	} else if qoffset <= -12 {
		return meastype.QOffsetMinus12dB
	} else if qoffset <= -10 {
		return meastype.QOffsetMinus10dB
	} else if qoffset <= -8 {
		return meastype.QOffsetMinus8dB
	} else if qoffset <= -6 {
		return meastype.QOffsetMinus6dB
	} else if qoffset <= -5 {
		return meastype.QOffsetMinus5dB
	} else if qoffset <= -4 {
		return meastype.QOffsetMinus4dB
	} else if qoffset <= -3 {
		return meastype.QOffsetMinus3dB
	} else if qoffset <= -2 {
		return meastype.QOffsetMinus2dB
	} else if qoffset <= -1 {
		return meastype.QOffsetMinus1dB
	} else if qoffset <= 0 {
		return meastype.QOffset0dB
	} else if qoffset <= 1 {
		return meastype.QOffset1dB
	} else if qoffset <= 2 {
		return meastype.QOffset2dB
	} else if qoffset <= 3 {
		return meastype.QOffset3dB
	} else if qoffset <= 4 {
		return meastype.QOffset4dB
	} else if qoffset <= 5 {
		return meastype.QOffset5dB
	} else if qoffset <= 6 {
		return meastype.QOffset6dB
	} else if qoffset <= 8 {
		return meastype.QOffset8dB
	} else if qoffset <= 10 {
		return meastype.QOffset10dB
	} else if qoffset <= 12 {
		return meastype.QOffset12dB
	} else if qoffset <= 14 {
		return meastype.QOffset14dB
	} else if qoffset <= 16 {
		return meastype.QOffset16dB
	} else if qoffset <= 18 {
		return meastype.QOffset18dB
	} else if qoffset <= 20 {
		return meastype.QOffset20dB
	} else if qoffset <= 22 {
		return meastype.QOffset22dB
	} else {
		return meastype.QOffset24dB
	}
}

func (c *measReportConverter) convertTimeToTrigger(ttt int32) meastype.TimeToTriggerRange {
	if ttt <= 0 {
		return meastype.TTT0ms
	} else if ttt <= 40 {
		return meastype.TTT40ms
	} else if ttt <= 64 {
		return meastype.TTT64ms
	} else if ttt <= 80 {
		return meastype.TTT80ms
	} else if ttt <= 100 {
		return meastype.TTT100ms
	} else if ttt <= 128 {
		return meastype.TTT128ms
	} else if ttt <= 160 {
		return meastype.TTT160ms
	} else if ttt <= 256 {
		return meastype.TTT256ms
	} else if ttt <= 320 {
		return meastype.TTT320ms
	} else if ttt <= 480 {
		return meastype.TTT480ms
	} else if ttt <= 512 {
		return meastype.TTT512ms
	} else if ttt <= 640 {
		return meastype.TTT640ms
	} else if ttt <= 1024 {
		return meastype.TTT1024ms
	} else if ttt <= 1280 {
		return meastype.TTT1280ms
	} else if ttt <= 2560 {
		return meastype.TTT2560ms
	} else {
		return meastype.TTT5120ms
	}
}
