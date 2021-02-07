// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

func CreateE2SmRcPreActionDefinition(actionDef int32) (*e2sm_rc_pre_ies.E2SmRcPreActionDefinition, error) {

	e2SmRcPrePdu := e2sm_rc_pre_ies.E2SmRcPreActionDefinition{
		RicStyleType: &e2sm_rc_pre_ies.RicStyleType{
			Value: actionDef,
		},
	}

	if err := e2SmRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmRcPrePDU %s", err.Error())
	}
	return &e2SmRcPrePdu, nil
}
