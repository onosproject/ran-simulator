// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

func CreateE2SmRcPreControlOutcome(RANparameterID int32, RANparameterValue int32) (*e2sm_rc_pre_ies.E2SmRcPreControlOutcome, error) {

	e2smRcPreOutcomeFormat1 := e2sm_rc_pre_ies.E2SmRcPreControlOutcomeFormat1{
		OutcomeElementList: make([]*e2sm_rc_pre_ies.RanparameterItem, 0),
	}
	outcomeElementList := &e2sm_rc_pre_ies.RanparameterItem{
		RanParameterId: &e2sm_rc_pre_ies.RanparameterId{
			Value: RANparameterID,
		},
		RanParameterValue: &e2sm_rc_pre_ies.RanparameterValue{
			RanparameterValue: &e2sm_rc_pre_ies.RanparameterValue_ValueInt{
				ValueInt: RANparameterValue,
			},
		},
	}
	e2smRcPreOutcomeFormat1.OutcomeElementList = append(e2smRcPreOutcomeFormat1.OutcomeElementList, outcomeElementList)
	e2smRcPrePdu := e2sm_rc_pre_ies.E2SmRcPreControlOutcome{
		E2SmRcPreControlOutcome: &e2sm_rc_pre_ies.E2SmRcPreControlOutcome_ControlOutcomeFormat1{
			ControlOutcomeFormat1: &e2smRcPreOutcomeFormat1,
		},
	}

	if err := e2smRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &e2smRcPrePdu, nil
}
