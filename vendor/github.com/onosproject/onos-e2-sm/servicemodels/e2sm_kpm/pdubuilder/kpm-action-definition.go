// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
)

func CreateE2SmKpmActionDefinition(actionDef int32) (*e2sm_kpm_ies.E2SmKpmActionDefinition, error) {

	e2SmKpmPdu := e2sm_kpm_ies.E2SmKpmActionDefinition{
		RicStyleType: &e2sm_kpm_ies.RicStyleType{
			Value: actionDef,
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmPDU %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}
