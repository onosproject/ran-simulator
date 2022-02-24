// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measurement

import (
	"context"
	"math"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	utils "github.com/onosproject/ran-simulator/pkg/utils/measurement"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
	"github.com/onosproject/rrm-son-lib/pkg/model/id"
	"github.com/onosproject/rrm-son-lib/pkg/model/measurement"
	meastype "github.com/onosproject/rrm-son-lib/pkg/model/measurement/type"
)

var logConverter = logging.GetLogger()

var qoffsetRanges = utils.QOffsetRanges{
	{Min: math.MinInt32, Max: -24, Value: meastype.QOffsetMinus24dB},
	{Min: -24, Max: -22, Value: meastype.QOffsetMinus22dB},
	{Min: -22, Max: -22, Value: meastype.QOffsetMinus20dB},
	{Min: -20, Max: -20, Value: meastype.QOffsetMinus18dB},
	{Min: -18, Max: -18, Value: meastype.QOffsetMinus16dB},
	{Min: -16, Max: -16, Value: meastype.QOffsetMinus14dB},
	{Min: -14, Max: -14, Value: meastype.QOffsetMinus12dB},
	{Min: -12, Max: -12, Value: meastype.QOffsetMinus10dB},
	{Min: -10, Max: -10, Value: meastype.QOffsetMinus8dB},
	{Min: -8, Max: -6, Value: meastype.QOffsetMinus6dB},
	{Min: -6, Max: -5, Value: meastype.QOffsetMinus5dB},
	{Min: -5, Max: -4, Value: meastype.QOffsetMinus4dB},
	{Min: -4, Max: -3, Value: meastype.QOffsetMinus3dB},
	{Min: -3, Max: -2, Value: meastype.QOffsetMinus2dB},
	{Min: -2, Max: -1, Value: meastype.QOffsetMinus1dB},
	{Min: -1, Max: 0, Value: meastype.QOffset0dB},
	{Min: 0, Max: 1, Value: meastype.QOffset1dB},
	{Min: 1, Max: 2, Value: meastype.QOffset2dB},
	{Min: 2, Max: 3, Value: meastype.QOffset3dB},
	{Min: 3, Max: 4, Value: meastype.QOffset4dB},
	{Min: 4, Max: 5, Value: meastype.QOffset5dB},
	{Min: 5, Max: 6, Value: meastype.QOffset6dB},
	{Min: 6, Max: 8, Value: meastype.QOffset8dB},
	{Min: 8, Max: 10, Value: meastype.QOffset10dB},
	{Min: 10, Max: 12, Value: meastype.QOffset12dB},
	{Min: 12, Max: 14, Value: meastype.QOffset14dB},
	{Min: 14, Max: 16, Value: meastype.QOffset16dB},
	{Min: 16, Max: 18, Value: meastype.QOffset18dB},
	{Min: 18, Max: 20, Value: meastype.QOffset20dB},
	{Min: 20, Max: 22, Value: meastype.QOffset22dB},
	{Min: 22, Max: math.MaxInt32, Value: meastype.QOffset24dB},
}

var tttRanges = utils.TimeToTriggerRanges{
	{Min: math.MinInt32, Max: 0, Value: meastype.TTT0ms},
	{Min: 40, Max: 64, Value: meastype.TTT40ms},
	{Min: 64, Max: 80, Value: meastype.TTT64ms},
	{Min: 80, Max: 100, Value: meastype.TTT80ms},
	{Min: 100, Max: 128, Value: meastype.TTT100ms},
	{Min: 128, Max: 160, Value: meastype.TTT128ms},
	{Min: 160, Max: 256, Value: meastype.TTT160ms},
	{Min: 256, Max: 320, Value: meastype.TTT256ms},
	{Min: 320, Max: 480, Value: meastype.TTT320ms},
	{Min: 480, Max: 512, Value: meastype.TTT480ms},
	{Min: 512, Max: 640, Value: meastype.TTT512ms},
	{Min: 640, Max: 1024, Value: meastype.TTT640ms},
	{Min: 1024, Max: 1280, Value: meastype.TTT1024ms},
	{Min: 1280, Max: 2560, Value: meastype.TTT1280ms},
	{Min: 2560, Max: 5120, Value: meastype.TTT2560ms},
	{Min: 5120, Max: math.MaxInt32, Value: meastype.TTT5120ms},
}

// MeasReportConverter is an abstraction of measurement report converter
type MeasReportConverter interface {
	// Convert forms measurement report defined in rrm-son-lib from model.ue defined in ransim
	Convert(ctx context.Context, ue *model.UE) device.UE
}

type measReportConverter struct {
	cellStore cells.Store
	ueStore   ues.Store
}

// NewMeasReportConverter returns the measurement report converter object
func NewMeasReportConverter(cellStore cells.Store, ueStore ues.Store) MeasReportConverter {
	return &measReportConverter{
		cellStore: cellStore,
		ueStore:   ueStore,
	}
}

func (c *measReportConverter) Convert(ctx context.Context, ue *model.UE) device.UE {
	ueid := id.NewUEID(uint64(ue.IMSI), uint32(ue.CRNTI), uint64(ue.Cell.NCGI))
	sCellInStore, err := c.cellStore.Get(ctx, ue.Cell.NCGI)
	if err != nil {
		logConverter.Errorf("Can't get serving cell from cell store: %v", err)
	}
	sCell := device.NewCell(id.NewECGI(uint64(sCellInStore.NCGI)),
		c.convertA3Offset(sCellInStore.MeasurementParams.EventA3Params.A3Offset),
		c.convertHysteresis(sCellInStore.MeasurementParams.Hysteresis),
		c.convertQOffset(sCellInStore.MeasurementParams.PCellIndividualOffset),
		c.convertQOffset(sCellInStore.MeasurementParams.FrequencyOffset),
		c.convertTimeToTrigger(sCellInStore.MeasurementParams.TimeToTrigger))

	var csCells []device.Cell
	measurements := make(map[string]measurement.Measurement)
	sCellMeas := measurement.NewMeasEventA3(id.NewECGI(uint64(sCellInStore.NCGI)), measurement.RSRP(ue.Cell.Strength))
	measurements[sCellMeas.GetCellID().String()] = sCellMeas

	for _, ueCell := range ue.Cells {
		tmpCellInStore, _ := c.cellStore.Get(ctx, ueCell.NCGI)
		if err != nil {
			logConverter.Errorf("Can't get candidate serving cell from cell storeL: %v", err)
		}

		var csCellIndividualOffset int32
		if _, ok := sCellInStore.MeasurementParams.NCellIndividualOffsets[ueCell.NCGI]; !ok {
			csCellIndividualOffset = 0
		} else {
			csCellIndividualOffset = sCellInStore.MeasurementParams.NCellIndividualOffsets[ueCell.NCGI]
		}

		csCells = append(csCells, device.NewCell(id.NewECGI(uint64(tmpCellInStore.NCGI)),
			c.convertA3Offset(tmpCellInStore.MeasurementParams.EventA3Params.A3Offset),
			c.convertHysteresis(tmpCellInStore.MeasurementParams.Hysteresis),
			c.convertQOffset(csCellIndividualOffset),
			c.convertQOffset(tmpCellInStore.MeasurementParams.FrequencyOffset),
			c.convertTimeToTrigger(tmpCellInStore.MeasurementParams.TimeToTrigger)))

		tmpCsCell := measurement.NewMeasEventA3(id.NewECGI(uint64(tmpCellInStore.NCGI)), measurement.RSRP(ueCell.Strength))
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
	return qoffsetRanges.Search(qoffset)
}

func (c *measReportConverter) convertTimeToTrigger(ttt int32) meastype.TimeToTriggerRange {
	return tttRanges.Search(ttt)
}
