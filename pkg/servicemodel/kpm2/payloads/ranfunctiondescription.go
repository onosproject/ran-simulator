// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package payloads

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/pdubuilder"
	e2sm_kpm_v2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

var log = logging.GetLogger("sm", "kpm2")

const (
	ricStyleType           = 1
	ricStyleName           = "Periodic Report"
	ricFormatType          = 5
	ricIndMsgFormat        = 1
	ricIndHdrFormat        = 1
	ranFunctionDescription = "KPM 2.0 Monitor"
	ranFunctionShortName   = "ORAN-E2SM-KPM"
	ranFunctionE2SmOid     = "1.3.6.1.4.1.53148.1.2.2.2"
	ranFunctionInstance    = 1
)

// RanFunctionDescriptionBytes - breaking out the construction of the RFD
func RanFunctionDescriptionBytes(modelPlmnID ransimtypes.PlmnID, modelPluginRegistry modelplugins.ModelRegistry) ([]byte, error) {
	plmnID := ransimtypes.NewUint24(uint32(modelPlmnID)).ToBytes()
	bs := e2sm_kpm_v2.BitString{
		Value: 0x9bcd4,
		Len:   22,
	}
	// TODO - Fix hardcoded cellID
	cellGlobalID, err := pdubuilder.CreateCellGlobalIDNRCGI(plmnID, 0xabcdef012<<28) // 36 bit
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// TODO - Fix hardcoded cellObjID
	var cellObjID = "ONF"
	cellMeasObjItem := pdubuilder.CreateCellMeasurementObjectItem(cellObjID, cellGlobalID)

	// TODO - Fix hardcoded IDs
	var gnbCuUpID int64 = 12345
	var gnbDuID int64 = 6789
	globalKpmnodeID, err := pdubuilder.CreateGlobalKpmnodeIDgNBID(&bs, plmnID, gnbCuUpID, gnbDuID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cmol := make([]*e2sm_kpm_v2.CellMeasurementObjectItem, 0)
	cmol = append(cmol, cellMeasObjItem)

	kpmNodeItem := pdubuilder.CreateRicKpmnodeItem(globalKpmnodeID, cmol)

	rknl := make([]*e2sm_kpm_v2.RicKpmnodeItem, 0)
	rknl = append(rknl, kpmNodeItem)

	retsi := pdubuilder.CreateRicEventTriggerStyleItem(ricStyleType, ricStyleName, ricFormatType)

	retsl := make([]*e2sm_kpm_v2.RicEventTriggerStyleItem, 0)
	retsl = append(retsl, retsi)

	measInfoActionList := e2sm_kpm_v2.MeasurementInfoActionList{
		Value: make([]*e2sm_kpm_v2.MeasurementInfoActionItem, 0),
	}

	var measTypeName1 = "RRC.ConnEstabAtt.Tot"
	var measTypeID1 int32 = 1
	measInfoActionItem1 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName1, measTypeID1)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem1)

	var measTypeName2 = "RRC.ConnEstabSucc.Tot"
	var measTypeID2 int32 = 2
	measInfoActionItem2 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName2, measTypeID2)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem2)

	var measTypeName3 = "RRC.ConnReEstabAtt.Tot"
	var measTypeID3 int32 = 3
	measInfoActionItem3 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName3, measTypeID3)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem3)

	var measTypeName4 = "RRC.ConnReEstabAtt.reconfigFail"
	var measTypeID4 int32 = 4
	measInfoActionItem4 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName4, measTypeID4)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem4)

	var measTypeName5 = "RRC.ConnReEstabAtt.HOFail"
	var measTypeID5 int32 = 5
	measInfoActionItem5 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName5, measTypeID5)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem5)

	var measTypeName6 = "RRC.ConnReEstabAtt.Other"
	var measTypeID6 int32 = 6
	measInfoActionItem6 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName6, measTypeID6)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem6)

	var measTypeName7 = "RRC.Conn.Avg"
	var measTypeID7 int32 = 7
	measInfoActionItem7 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName7, measTypeID7)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem7)

	var measTypeName8 = "RRC.Conn.Max"
	var measTypeID8 int32 = 8
	measInfoActionItem8 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName8, measTypeID8)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem8)

	rrsi := pdubuilder.CreateRicReportStyleItem(ricStyleType, ricStyleName, ricFormatType, &measInfoActionList, ricIndHdrFormat, ricIndMsgFormat)

	rrsl := make([]*e2sm_kpm_v2.RicReportStyleItem, 0)
	rrsl = append(rrsl, rrsi)

	ranFuncDescPdu, err := pdubuilder.CreateE2SmKpmRanfunctionDescription(
		ranFunctionShortName,
		ranFunctionE2SmOid,
		ranFunctionDescription,
		ranFunctionInstance,
		rknl,
		retsl,
		rrsl)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	protoBytes, err := proto.Marshal(ranFuncDescPdu)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	kpmModelPlugin, _ := modelPluginRegistry.GetPlugin(ranFunctionE2SmOid)
	if kpmModelPlugin == nil {
		return nil, errors.New(errors.Invalid, "model plugin is nil")
	}
	ranFuncDescBytes, err := kpmModelPlugin.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return ranFuncDescBytes, nil
}
