// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

func CreateE2SmRcPreControlMessage(RANparameterID int32, RANparameterName string, RANparameterValue int32) (*e2sm_rc_pre_ies.E2SmRcPreControlMessage, error) {

	e2smRcPreMsgFormat1 := e2sm_rc_pre_ies.E2SmRcPreControlMessageFormat1{
		ParameterType: &e2sm_rc_pre_ies.RanparameterDefItem{
			RanParameterId: &e2sm_rc_pre_ies.RanparameterId{
				Value: RANparameterID,
			},
			RanParameterName: &e2sm_rc_pre_ies.RanparameterName{
				Value: RANparameterName,
			},
			RanParameterType: e2sm_rc_pre_ies.RanparameterType_RANPARAMETER_TYPE_INTEGER,
		},
		ParameterVal: &e2sm_rc_pre_ies.RanparameterValue{
			RanparameterValue: &e2sm_rc_pre_ies.RanparameterValue_ValueInt{
				ValueInt: RANparameterValue,
			},
		},
	}
	e2smRcPrePdu := e2sm_rc_pre_ies.E2SmRcPreControlMessage{
		E2SmRcPreControlMessage: &e2sm_rc_pre_ies.E2SmRcPreControlMessage_ControlMessage{
			ControlMessage: &e2smRcPreMsgFormat1,
		},
	}

	if err := e2smRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &e2smRcPrePdu, nil
}
