// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
)

func CreateE2SmKpmEventTriggerDefinition(rtPeriod int32) (*e2sm_kpm_ies.E2SmKpmEventTriggerDefinition, error) {
	if rtPeriod < 0 && rtPeriod > 19 {
		return nil, fmt.Errorf("reportPeriodIe is out of range. Should be from 0 to 19")
	}

	policyTestItem := &e2sm_kpm_ies.TriggerConditionIeItem{
		ReportPeriodIe: e2sm_kpm_ies.RtPeriodIe(rtPeriod),
	}

	eventDefinitionFormat1 := &e2sm_kpm_ies.E2SmKpmEventTriggerDefinitionFormat1{
		PolicyTestList: make([]*e2sm_kpm_ies.TriggerConditionIeItem, 0),
	}
	eventDefinitionFormat1.PolicyTestList = append(eventDefinitionFormat1.PolicyTestList, policyTestItem)

	e2SmKpmPdu := e2sm_kpm_ies.E2SmKpmEventTriggerDefinition{
		E2SmKpmEventTriggerDefinition: &e2sm_kpm_ies.E2SmKpmEventTriggerDefinition_EventDefinitionFormat1{
			EventDefinitionFormat1: eventDefinitionFormat1,
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmPDU %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}
