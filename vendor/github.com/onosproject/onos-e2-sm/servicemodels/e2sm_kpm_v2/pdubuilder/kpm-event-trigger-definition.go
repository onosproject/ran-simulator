// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_v2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
)

func CreateE2SmKpmEventTriggerDefinition(rtPeriod int32) (*e2sm_kpm_v2.E2SmKpmEventTriggerDefinition, error) {
	if rtPeriod < 0 && rtPeriod > 19 {
		return nil, fmt.Errorf("reportPeriodIe is out of range. Should be from 0 to 19")
	}

	e2SmKpmPdu := e2sm_kpm_v2.E2SmKpmEventTriggerDefinition{
		E2SmKpmEventTriggerDefinition: &e2sm_kpm_v2.E2SmKpmEventTriggerDefinition_EventDefinitionFormat1{
			EventDefinitionFormat1: &e2sm_kpm_v2.E2SmKpmEventTriggerDefinitionFormat1{
				ReportingPeriod: rtPeriod,
			},
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmPDU %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}
