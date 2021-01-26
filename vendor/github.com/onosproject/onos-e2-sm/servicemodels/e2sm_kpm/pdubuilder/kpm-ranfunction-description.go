// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
)

func CreateE2SmKpmRanfunctionDescriptionMsg(ranFunctionShortName string, ranFunctionE2SmOid string, ranFunctionDescription string,
	ranFunctionInstance int32, ricEventStyleType int32, ricEventStyleName string, ricEventFormatType int32,
	ricReportStyleType int32, ricReportStyleName string, ricIndicationHeaderFormatType int32,
	ricIndicationMessageFormatType int32) (*e2sm_kpm_ies.E2SmKpmRanfunctionDescription, error) {

	ranfunctionItem := e2sm_kpm_ies.E2SmKpmRanfunctionDescription_E2SmKpmRanfunctionItem001{
		RicEventTriggerStyleList: make([]*e2sm_kpm_ies.RicEventTriggerStyleList, 0),
		RicReportStyleList:       make([]*e2sm_kpm_ies.RicReportStyleList, 0),
	}

	ricEventTriggerStyleList := e2sm_kpm_ies.RicEventTriggerStyleList{
		RicEventTriggerStyleType: &e2sm_kpm_ies.RicStyleType{
			Value: ricEventStyleType, //int32
		},
		RicEventTriggerStyleName: &e2sm_kpm_ies.RicStyleName{
			Value: ricEventStyleName, //string
		},
		RicEventTriggerFormatType: &e2sm_kpm_ies.RicFormatType{
			Value: ricEventFormatType, //int32
		},
	}
	ranfunctionItem.RicEventTriggerStyleList = append(ranfunctionItem.RicEventTriggerStyleList, &ricEventTriggerStyleList)

	ricReportStyleList := e2sm_kpm_ies.RicReportStyleList{
		RicReportStyleType: &e2sm_kpm_ies.RicStyleType{
			Value: ricReportStyleType, //int32
		},
		RicReportStyleName: &e2sm_kpm_ies.RicStyleName{
			Value: ricReportStyleName, //string
		},
		RicIndicationHeaderFormatType: &e2sm_kpm_ies.RicFormatType{
			Value: ricIndicationHeaderFormatType, //int32
		},
		RicIndicationMessageFormatType: &e2sm_kpm_ies.RicFormatType{
			Value: ricIndicationMessageFormatType, //int32
		},
	}
	ranfunctionItem.RicReportStyleList = append(ranfunctionItem.RicReportStyleList, &ricReportStyleList)

	e2smKpmPdu := e2sm_kpm_ies.E2SmKpmRanfunctionDescription{
		RanFunctionName: &e2sm_kpm_ies.RanfunctionName{
			RanFunctionShortName:   ranFunctionShortName,   //string
			RanFunctionE2SmOid:     ranFunctionE2SmOid,     //sting
			RanFunctionDescription: ranFunctionDescription, //string
			RanFunctionInstance:    ranFunctionInstance,    //int32
		},
		E2SmKpmRanfunctionItem: &ranfunctionItem,
	}

	if err := e2smKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &e2smKpmPdu, nil
}
