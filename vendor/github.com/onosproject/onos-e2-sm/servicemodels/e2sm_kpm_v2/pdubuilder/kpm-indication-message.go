// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_v2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-ies"
)

func CreateE2SmKpmIndicationMessage(subscriptionID int64, cellObjID string, granularity int32,
	measInfoList *e2sm_kpm_v2.MeasurementInfoList, measData *e2sm_kpm_v2.MeasurementData) (*e2sm_kpm_v2.E2SmKpmIndicationMessage, error) {

	e2SmKpmPdu := e2sm_kpm_v2.E2SmKpmIndicationMessage{
		E2SmKpmIndicationMessage: &e2sm_kpm_v2.E2SmKpmIndicationMessage_IndicationMessageFormat1{
			IndicationMessageFormat1: &e2sm_kpm_v2.E2SmKpmIndicationMessageFormat1{
				SubscriptId: &e2sm_kpm_v2.SubscriptionId{
					Value: subscriptionID,
				},
				CellObjId: &e2sm_kpm_v2.CellObjectId{
					Value: cellObjID,
				},
				GranulPeriod: &e2sm_kpm_v2.GranularityPeriod{
					Value: granularity,
				},
				MeasInfoList: measInfoList,
				MeasData:     measData,
			},
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmPDU %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}

func CreateMeasurementRecordItemInteger(integer int64) *e2sm_kpm_v2.MeasurementRecordItem {

	return &e2sm_kpm_v2.MeasurementRecordItem{
		MeasurementRecordItem: &e2sm_kpm_v2.MeasurementRecordItem_Integer{
			Integer: integer,
		},
	}
}

func CreateMeasurementRecordItemReal(real float64) *e2sm_kpm_v2.MeasurementRecordItem {

	return &e2sm_kpm_v2.MeasurementRecordItem{
		MeasurementRecordItem: &e2sm_kpm_v2.MeasurementRecordItem_Real{
			Real: real,
		},
	}
}

func CreateMeasurementRecordItemNoValue() *e2sm_kpm_v2.MeasurementRecordItem {

	return &e2sm_kpm_v2.MeasurementRecordItem{
		MeasurementRecordItem: &e2sm_kpm_v2.MeasurementRecordItem_NoValue{
			NoValue: 0,
		},
	}
}
