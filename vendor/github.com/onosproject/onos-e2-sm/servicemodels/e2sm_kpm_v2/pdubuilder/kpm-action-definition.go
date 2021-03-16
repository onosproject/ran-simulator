// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_v2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-ies"
)

func CreateE2SmKpmActionDefinition(ricStyleType int32,
	actionDefinition *e2sm_kpm_v2.E2SmKpmActionDefinitionFormat1) (*e2sm_kpm_v2.E2SmKpmActionDefinition, error) {

	e2SmKpmPdu := e2sm_kpm_v2.E2SmKpmActionDefinition{
		RicStyleType: &e2sm_kpm_v2.RicStyleType{
			Value: ricStyleType,
		},
		E2SmKpmActionDefinition: &e2sm_kpm_v2.E2SmKpmActionDefinition_ActionDefinitionFormat1{
			ActionDefinitionFormat1: actionDefinition,
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmActionDefinition %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}

func CreateActionDefinitionFormat1(cellObjID string, measInfoList *e2sm_kpm_v2.MeasurementInfoList,
	granularity int32, subID int64) (*e2sm_kpm_v2.E2SmKpmActionDefinitionFormat1, error) {

	actionDefinitionFormat1 := e2sm_kpm_v2.E2SmKpmActionDefinitionFormat1{
		CellObjId: &e2sm_kpm_v2.CellObjectId{
			Value: cellObjID,
		},
		MeasInfoList: measInfoList,
		GranulPeriod: &e2sm_kpm_v2.GranularityPeriod{
			Value: granularity,
		},
		SubscriptId: &e2sm_kpm_v2.SubscriptionId{
			Value: subID,
		},
	}

	if err := actionDefinitionFormat1.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmActionDefinitionFormat1 %s", err.Error())
	}

	return &actionDefinitionFormat1, nil
}

func CreateMeasurementInfoItem(measType *e2sm_kpm_v2.MeasurementType, labelInfoList *e2sm_kpm_v2.LabelInfoList) (*e2sm_kpm_v2.MeasurementInfoItem, error) {

	item := e2sm_kpm_v2.MeasurementInfoItem{
		MeasType:      measType,
		LabelInfoList: labelInfoList,
	}

	if err := item.Validate(); err != nil {
		return nil, fmt.Errorf("error validating MeasurementInfoItem %s", err.Error())
	}

	return &item, nil
}

func CreateMeasurementTypeMeasID(measTypeID int32) (*e2sm_kpm_v2.MeasurementType, error) {
	measType := e2sm_kpm_v2.MeasurementType{
		MeasurementType: &e2sm_kpm_v2.MeasurementType_MeasId{
			MeasId: &e2sm_kpm_v2.MeasurementTypeId{
				Value: measTypeID,
			},
		},
	}

	if err := measType.Validate(); err != nil {
		return nil, fmt.Errorf("error validating MeasurementType %s", err.Error())
	}

	return &measType, nil
}

func CreateMeasurementTypeMeasName(measName string) (*e2sm_kpm_v2.MeasurementType, error) {
	measType := e2sm_kpm_v2.MeasurementType{
		MeasurementType: &e2sm_kpm_v2.MeasurementType_MeasName{
			MeasName: &e2sm_kpm_v2.MeasurementTypeName{
				Value: measName,
			},
		},
	}

	if err := measType.Validate(); err != nil {
		return nil, fmt.Errorf("error validating MeasurementType %s", err.Error())
	}

	return &measType, nil
}

func CreateLabelInfoItem(plmnID []byte, sst []byte, sd []byte, fiveQI int32, qci int32, qciMax int32, qciMin int32,
	arpMax int32, arpMin int32, bitrateRange int32, layerMuMimo int32, distX int32, distY int32, distZ int32,
	startEndIndication e2sm_kpm_v2.StartEndInd) (*e2sm_kpm_v2.LabelInfoItem, error) {

	if len(sst) != 1 {
		return nil, fmt.Errorf("error: SST should be 1 chars")
	}
	if len(sd) != 3 {
		return nil, fmt.Errorf("error: SD should be 3 chars")
	}
	if len(plmnID) != 3 {
		return nil, fmt.Errorf("error: Plmn ID should be 3 chars")
	}
	if arpMax < 1 && arpMax > 15 {
		return nil, fmt.Errorf("error: ARP values must be in rang [1, 15]")
	}
	if arpMin < 1 && arpMin > 15 {
		return nil, fmt.Errorf("error: ARP values must be in rang [1, 15]")
	}

	labelInfoItem := e2sm_kpm_v2.LabelInfoItem{
		MeasLabel: &e2sm_kpm_v2.MeasurementLabel{
			PlmnId: &e2sm_kpm_v2.PlmnIdentity{
				Value: plmnID,
			},
			SliceId: &e2sm_kpm_v2.Snssai{
				SD:  sd,
				SSt: sst,
			},
			FiveQi: &e2sm_kpm_v2.FiveQi{
				Value: fiveQI,
			},
			QCi: &e2sm_kpm_v2.Qci{
				Value: qci,
			},
			QCimax: &e2sm_kpm_v2.Qci{
				Value: qciMax,
			},
			QCimin: &e2sm_kpm_v2.Qci{
				Value: qciMin,
			},
			ARpmax: &e2sm_kpm_v2.Arp{
				Value: arpMax,
			},
			ARpmin: &e2sm_kpm_v2.Arp{
				Value: arpMin,
			},
			BitrateRange:     bitrateRange,
			LayerMuMimo:      layerMuMimo,
			SUm:              e2sm_kpm_v2.SUM_SUM_TRUE,
			DistBinX:         distX,
			DistBinY:         distY,
			DistBinZ:         distZ,
			PreLabelOverride: e2sm_kpm_v2.PreLabelOverride_PRE_LABEL_OVERRIDE_TRUE,
			StartEndInd:      startEndIndication,
		},
	}

	if err := labelInfoItem.Validate(); err != nil {
		return nil, fmt.Errorf("error validating LabelInfoItem %s", err.Error())
	}

	return &labelInfoItem, nil
}
