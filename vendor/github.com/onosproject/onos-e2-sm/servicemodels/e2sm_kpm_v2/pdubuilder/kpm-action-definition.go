// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_v2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
)

func CreateE2SmKpmActionDefinitionFormat1(ricStyleType int32,
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

func CreateE2SmKpmActionDefinitionFormat2(ricStyleType int32,
	actionDefinitionFormat2 *e2sm_kpm_v2.E2SmKpmActionDefinitionFormat2) (*e2sm_kpm_v2.E2SmKpmActionDefinition, error) {

	e2SmKpmPdu := e2sm_kpm_v2.E2SmKpmActionDefinition{
		RicStyleType: &e2sm_kpm_v2.RicStyleType{
			Value: ricStyleType,
		},
		E2SmKpmActionDefinition: &e2sm_kpm_v2.E2SmKpmActionDefinition_ActionDefinitionFormat2{
			ActionDefinitionFormat2: actionDefinitionFormat2,
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmActionDefinition %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}

func CreateE2SmKpmActionDefinitionFormat3(ricStyleType int32,
	actionDefinitionFormat3 *e2sm_kpm_v2.E2SmKpmActionDefinitionFormat3) (*e2sm_kpm_v2.E2SmKpmActionDefinition, error) {

	e2SmKpmPdu := e2sm_kpm_v2.E2SmKpmActionDefinition{
		RicStyleType: &e2sm_kpm_v2.RicStyleType{
			Value: ricStyleType,
		},
		E2SmKpmActionDefinition: &e2sm_kpm_v2.E2SmKpmActionDefinition_ActionDefinitionFormat3{
			ActionDefinitionFormat3: actionDefinitionFormat3,
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

func CreateActionDefinitionFormat2(ueID string, actionDefinitionFormat1 *e2sm_kpm_v2.E2SmKpmActionDefinitionFormat1) (*e2sm_kpm_v2.E2SmKpmActionDefinitionFormat2, error) {

	actionDefinitionFormat2 := e2sm_kpm_v2.E2SmKpmActionDefinitionFormat2{
		UeId: &e2sm_kpm_v2.UeIdentity{
			Value: ueID,
		},
		SubscriptInfo: actionDefinitionFormat1,
	}

	if err := actionDefinitionFormat2.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmActionDefinitionFormat2 %s", err.Error())
	}

	return &actionDefinitionFormat2, nil
}

func CreateActionDefinitionFormat3(cellObjID string, measCondList *e2sm_kpm_v2.MeasurementCondList,
	granularity int32, subID int64) (*e2sm_kpm_v2.E2SmKpmActionDefinitionFormat3, error) {

	actionDefinitionFormat3 := e2sm_kpm_v2.E2SmKpmActionDefinitionFormat3{
		CellObjId: &e2sm_kpm_v2.CellObjectId{
			Value: cellObjID,
		},
		MeasCondList: measCondList,
		GranulPeriod: &e2sm_kpm_v2.GranularityPeriod{
			Value: granularity,
		},
		SubscriptId: &e2sm_kpm_v2.SubscriptionId{
			Value: subID,
		},
	}

	if err := actionDefinitionFormat3.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmActionDefinitionFormat3 %s", err.Error())
	}

	return &actionDefinitionFormat3, nil
}

func CreateMeasurementInfoItem(measType *e2sm_kpm_v2.MeasurementType, labelInfoList *e2sm_kpm_v2.LabelInfoList) (*e2sm_kpm_v2.MeasurementInfoItem, error) {

	item := e2sm_kpm_v2.MeasurementInfoItem{
		MeasType: measType,
	}

	// optional instance
	if labelInfoList != nil {
		item.LabelInfoList = labelInfoList
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

func CreateLabelInfoItem(plmnID []byte, sst []byte, sd []byte) (*e2sm_kpm_v2.LabelInfoItem, error) {

	labelInfoItem := e2sm_kpm_v2.LabelInfoItem{
		MeasLabel: &e2sm_kpm_v2.MeasurementLabel{},
	}

	if plmnID != nil {
		if len(plmnID) != 3 {
			return nil, fmt.Errorf("error: Plmn ID should be 3 chars")
		}
		labelInfoItem.MeasLabel.PlmnId = &e2sm_kpm_v2.PlmnIdentity{
			Value: plmnID,
		}
	}
	if sst != nil {
		if len(sst) != 1 {
			return nil, fmt.Errorf("error: SST should be 1 chars")
		}
		labelInfoItem.MeasLabel.SliceId = &e2sm_kpm_v2.Snssai{
			SSt: sst,
		}
		if sd != nil {
			if len(sd) != 3 {
				return nil, fmt.Errorf("error: SD should be 3 chars")
			}
			labelInfoItem.MeasLabel.SliceId.SD = sd
		}
	}

	labelInfoItem.MeasLabel.FiveQi = &e2sm_kpm_v2.FiveQi{
		Value: -1, // Not valid value, indicates this item not present in message - handled later in CGo encoding
	}
	labelInfoItem.MeasLabel.QCi = &e2sm_kpm_v2.Qci{
		Value: -1, // Not valid value, indicates this item not present in message - handled later in CGo encoding
	}
	labelInfoItem.MeasLabel.QCimax = &e2sm_kpm_v2.Qci{
		Value: -1, // Not valid value, indicates this item not present in message - handled later in CGo encoding
	}
	labelInfoItem.MeasLabel.QCimin = &e2sm_kpm_v2.Qci{
		Value: -1, // Not valid value, indicates this item not present in message - handled later in CGo encoding
	}
	labelInfoItem.MeasLabel.ARpmin = &e2sm_kpm_v2.Arp{
		Value: -1, // Not valid value, indicates this item not present in message - handled later in CGo encoding
	}
	labelInfoItem.MeasLabel.ARpmax = &e2sm_kpm_v2.Arp{
		Value: -1, // Not valid value, indicates this item not present in message - handled later in CGo encoding
	}
	labelInfoItem.MeasLabel.BitrateRange = -1     // Not valid value, indicates this item not present in message - handled later in CGo encoding
	labelInfoItem.MeasLabel.LayerMuMimo = -1      // Not valid value, indicates this item not present in message - handled later in CGo encoding
	labelInfoItem.MeasLabel.SUm = -1              // Not valid value, indicates this item not present in message - handled later in CGo encoding
	labelInfoItem.MeasLabel.DistBinX = -1         // Not valid value, indicates this item not present in message - handled later in CGo encoding
	labelInfoItem.MeasLabel.DistBinY = -1         // Not valid value, indicates this item not present in message - handled later in CGo encoding
	labelInfoItem.MeasLabel.DistBinZ = -1         // Not valid value, indicates this item not present in message - handled later in CGo encoding
	labelInfoItem.MeasLabel.PreLabelOverride = -1 // Not valid value, indicates this item not present in message - handled later in CGo encoding
	labelInfoItem.MeasLabel.StartEndInd = -1      // Not valid value, indicates this item not present in message - handled later in CGo encoding

	//if err := labelInfoItem.Validate(); err != nil {
	//	return nil, fmt.Errorf("error validating LabelInfoItem %s", err.Error())
	//}

	return &labelInfoItem, nil
}
func CreateMeasurementCondItem(measType *e2sm_kpm_v2.MeasurementType, measCondList *e2sm_kpm_v2.MatchingCondList) (*e2sm_kpm_v2.MeasurementCondItem, error) {

	measCondItem := e2sm_kpm_v2.MeasurementCondItem{
		MeasType:     measType,
		MatchingCond: measCondList,
	}

	if err := measCondItem.Validate(); err != nil {
		return nil, fmt.Errorf("error validating MeasurementCondItem %s", err.Error())
	}
	return &measCondItem, nil
}

func CreateMatchingCondItemMeasLabel(measLabel *e2sm_kpm_v2.MeasurementLabel) (*e2sm_kpm_v2.MatchingCondItem, error) {

	res := e2sm_kpm_v2.MatchingCondItem{
		MatchingCondItem: &e2sm_kpm_v2.MatchingCondItem_MeasLabel{
			MeasLabel: measLabel,
		},
	}

	if err := res.Validate(); err != nil {
		return nil, fmt.Errorf("error validating MatchingCondItem (MeasLabel) %s", err.Error())
	}
	return &res, nil
}

func CreateMatchingCondItemTestCondInfo(testCondInfo *e2sm_kpm_v2.TestCondInfo) (*e2sm_kpm_v2.MatchingCondItem, error) {

	res := e2sm_kpm_v2.MatchingCondItem{
		MatchingCondItem: &e2sm_kpm_v2.MatchingCondItem_TestCondInfo{
			TestCondInfo: testCondInfo,
		},
	}

	if err := res.Validate(); err != nil {
		return nil, fmt.Errorf("error validating MatchingCondItem (TestCondInfo) %s", err.Error())
	}
	return &res, nil
}

func CreateTestCondInfo(tct *e2sm_kpm_v2.TestCondType, tce e2sm_kpm_v2.TestCondExpression, tcv *e2sm_kpm_v2.TestCondValue) (*e2sm_kpm_v2.TestCondInfo, error) {

	tci := e2sm_kpm_v2.TestCondInfo{
		TestValue: tcv,
		TestExpr:  tce,
		TestType:  tct,
	}

	if err := tci.Validate(); err != nil {
		return nil, fmt.Errorf("error validating TestCondInfo (TestCondInfo) %s", err.Error())
	}
	return &tci, nil
}

func CreateTestCondTypeGBR() *e2sm_kpm_v2.TestCondType {

	return &e2sm_kpm_v2.TestCondType{
		TestCondType: &e2sm_kpm_v2.TestCondType_GBr{
			GBr: e2sm_kpm_v2.GBR_GBR_TRUE,
		},
	}
}

func CreateTestCondTypeAMBR() *e2sm_kpm_v2.TestCondType {

	return &e2sm_kpm_v2.TestCondType{
		TestCondType: &e2sm_kpm_v2.TestCondType_AMbr{
			AMbr: e2sm_kpm_v2.AMBR_AMBR_TRUE,
		},
	}
}

func CreateTestCondTypeIsStat() *e2sm_kpm_v2.TestCondType {

	return &e2sm_kpm_v2.TestCondType{
		TestCondType: &e2sm_kpm_v2.TestCondType_IsStat{
			IsStat: e2sm_kpm_v2.ISSTAT_ISSTAT_TRUE,
		},
	}
}

func CreateTestCondTypeIsCatM() *e2sm_kpm_v2.TestCondType {

	return &e2sm_kpm_v2.TestCondType{
		TestCondType: &e2sm_kpm_v2.TestCondType_IsCatM{
			IsCatM: e2sm_kpm_v2.ISCATM_ISCATM_TRUE,
		},
	}
}

func CreateTestCondTypeRSRP() *e2sm_kpm_v2.TestCondType {

	return &e2sm_kpm_v2.TestCondType{
		TestCondType: &e2sm_kpm_v2.TestCondType_RSrp{
			RSrp: e2sm_kpm_v2.RSRP_RSRP_TRUE,
		},
	}
}

func CreateTestCondTypeRSRQ() *e2sm_kpm_v2.TestCondType {

	return &e2sm_kpm_v2.TestCondType{
		TestCondType: &e2sm_kpm_v2.TestCondType_RSrq{
			RSrq: e2sm_kpm_v2.RSRQ_RSRQ_TRUE,
		},
	}
}

func CreateTestCondValueInt(val int64) *e2sm_kpm_v2.TestCondValue {

	return &e2sm_kpm_v2.TestCondValue{
		TestCondValue: &e2sm_kpm_v2.TestCondValue_ValueInt{
			ValueInt: val,
		},
	}
}

func CreateTestCondValueEnum(val int64) *e2sm_kpm_v2.TestCondValue {

	return &e2sm_kpm_v2.TestCondValue{
		TestCondValue: &e2sm_kpm_v2.TestCondValue_ValueEnum{
			ValueEnum: val,
		},
	}
}

func CreateTestCondValueBool(val bool) *e2sm_kpm_v2.TestCondValue {

	return &e2sm_kpm_v2.TestCondValue{
		TestCondValue: &e2sm_kpm_v2.TestCondValue_ValueBool{
			ValueBool: val,
		},
	}
}

func CreateTestCondValueBitS(val *e2sm_kpm_v2.BitString) *e2sm_kpm_v2.TestCondValue {

	return &e2sm_kpm_v2.TestCondValue{
		TestCondValue: &e2sm_kpm_v2.TestCondValue_ValueBitS{
			ValueBitS: val,
		},
	}
}

func CreateTestCondValueOctS(val string) *e2sm_kpm_v2.TestCondValue {

	return &e2sm_kpm_v2.TestCondValue{
		TestCondValue: &e2sm_kpm_v2.TestCondValue_ValueOctS{
			ValueOctS: val,
		},
	}
}

func CreateTestCondValuePrtS(val string) *e2sm_kpm_v2.TestCondValue {

	return &e2sm_kpm_v2.TestCondValue{
		TestCondValue: &e2sm_kpm_v2.TestCondValue_ValuePrtS{
			ValuePrtS: val,
		},
	}
}
