// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_v2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
)

func CreateE2SmKpmRanfunctionDescription(rfSn string, rfE2SMoid string, rfd string, rfi int32, rknl []*e2sm_kpm_v2.RicKpmnodeItem,
	retsl []*e2sm_kpm_v2.RicEventTriggerStyleItem, rrsl []*e2sm_kpm_v2.RicReportStyleItem) (*e2sm_kpm_v2.E2SmKpmRanfunctionDescription, error) {

	e2SmKpmPdu := e2sm_kpm_v2.E2SmKpmRanfunctionDescription{
		RanFunctionName: &e2sm_kpm_v2.RanfunctionName{
			RanFunctionShortName:   rfSn,
			RanFunctionE2SmOid:     rfE2SMoid,
			RanFunctionDescription: rfd,
			RanFunctionInstance:    rfi,
		},
		RicKpmNodeList:           rknl,
		RicEventTriggerStyleList: retsl,
		RicReportStyleList:       rrsl,
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmRanfunctionDescription %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}

func CreateRicKpmnodeItem(globalKpmnodeID *e2sm_kpm_v2.GlobalKpmnodeId, cmol []*e2sm_kpm_v2.CellMeasurementObjectItem) *e2sm_kpm_v2.RicKpmnodeItem {

	return &e2sm_kpm_v2.RicKpmnodeItem{
		RicKpmnodeType:            globalKpmnodeID,
		CellMeasurementObjectList: cmol,
	}
}

func CreateCellMeasurementObjectItem(cellObjID string, cellGlobalID *e2sm_kpm_v2.CellGlobalId) *e2sm_kpm_v2.CellMeasurementObjectItem {

	return &e2sm_kpm_v2.CellMeasurementObjectItem{
		CellObjectId: &e2sm_kpm_v2.CellObjectId{
			Value: cellObjID,
		},
		CellGlobalId: cellGlobalID,
	}
}

func CreateCellGlobalIDNRCGI(plmnID []byte, cellIDBits36 uint64) (*e2sm_kpm_v2.CellGlobalId, error) {

	if len(plmnID) != 3 {
		return nil, fmt.Errorf("PlmnID should be 3 chars")
	}

	if cellIDBits36&0x000000000fffffff > 0 {
		return nil, fmt.Errorf("bits should be at the left - not expecting anything in last 28 bits")
	}
	bs := e2sm_kpm_v2.BitString{
		Value: cellIDBits36,
		Len:   36,
	}

	return &e2sm_kpm_v2.CellGlobalId{
		CellGlobalId: &e2sm_kpm_v2.CellGlobalId_NrCgi{
			NrCgi: &e2sm_kpm_v2.Nrcgi{
				PLmnIdentity: &e2sm_kpm_v2.PlmnIdentity{
					Value: plmnID,
				},
				NRcellIdentity: &e2sm_kpm_v2.NrcellIdentity{
					Value: &bs,
				},
			},
		},
	}, nil
}

func CreateCellGlobalIDEUTRACGI(plmnID []byte, bs *e2sm_kpm_v2.BitString) (*e2sm_kpm_v2.CellGlobalId, error) {

	if len(plmnID) != 3 {
		return nil, fmt.Errorf("PlmnID should be 3 chars")
	}

	return &e2sm_kpm_v2.CellGlobalId{
		CellGlobalId: &e2sm_kpm_v2.CellGlobalId_EUtraCgi{
			EUtraCgi: &e2sm_kpm_v2.Eutracgi{
				PLmnIdentity: &e2sm_kpm_v2.PlmnIdentity{
					Value: plmnID,
				},
				EUtracellIdentity: &e2sm_kpm_v2.EutracellIdentity{
					Value: bs,
				},
			},
		},
	}, nil
}

func CreateRicEventTriggerStyleItem(ricStyleType int32, ricStyleName string, ricFormatType int32) *e2sm_kpm_v2.RicEventTriggerStyleItem {

	return &e2sm_kpm_v2.RicEventTriggerStyleItem{
		RicEventTriggerStyleType: &e2sm_kpm_v2.RicStyleType{
			Value: ricStyleType,
		},
		RicEventTriggerStyleName: &e2sm_kpm_v2.RicStyleName{
			Value: ricStyleName,
		},
		RicEventTriggerFormatType: &e2sm_kpm_v2.RicFormatType{
			Value: ricFormatType,
		},
	}
}

func CreateRicReportStyleItem(ricStyleType int32, ricStyleName string, ricFormatType int32,
	measInfoActionList *e2sm_kpm_v2.MeasurementInfoActionList, indHdrFormatType int32,
	indMsgFormatType int32) *e2sm_kpm_v2.RicReportStyleItem {

	return &e2sm_kpm_v2.RicReportStyleItem{
		RicReportStyleType: &e2sm_kpm_v2.RicStyleType{
			Value: ricStyleType,
		},
		RicReportStyleName: &e2sm_kpm_v2.RicStyleName{
			Value: ricStyleName,
		},
		RicActionFormatType: &e2sm_kpm_v2.RicFormatType{
			Value: ricFormatType,
		},
		MeasInfoActionList: measInfoActionList,
		RicIndicationHeaderFormatType: &e2sm_kpm_v2.RicFormatType{
			Value: indHdrFormatType,
		},
		RicIndicationMessageFormatType: &e2sm_kpm_v2.RicFormatType{
			Value: indMsgFormatType,
		},
	}
}

func CreateMeasurementInfoActionItem(measTypeName string, measTypeID int32) *e2sm_kpm_v2.MeasurementInfoActionItem {

	return &e2sm_kpm_v2.MeasurementInfoActionItem{
		MeasName: &e2sm_kpm_v2.MeasurementTypeName{
			Value: measTypeName,
		},
		MeasId: &e2sm_kpm_v2.MeasurementTypeId{
			Value: measTypeID,
		},
	}
}
