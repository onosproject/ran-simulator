// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

func CreateE2SmRcPreEventTriggerDefinition(rtPeriod int32) (*e2sm_rc_pre_ies.E2SmRcPreEventTriggerDefinition, error) {

	eventDefinitionFormat1 := &e2sm_rc_pre_ies.E2SmRcPreEventTriggerDefinitionFormat1{
		TriggerType:       e2sm_rc_pre_ies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_UPON_CHANGE,
		ReportingPeriodMs: rtPeriod,
	}

	E2SmRcPrePdu := e2sm_rc_pre_ies.E2SmRcPreEventTriggerDefinition{
		E2SmRcPreEventTriggerDefinitionEventDefinitionFormats: &e2sm_rc_pre_ies.E2SmRcPreEventTriggerDefinition_EventDefinitionFormat1{
			EventDefinitionFormat1: eventDefinitionFormat1,
		},
	}

	if err := E2SmRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmRcPrePDU %s", err.Error())
	}
	return &E2SmRcPrePdu, nil
}
