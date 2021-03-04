// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

func CreateE2SmRcPreControlHeader(controlMessagePriority int32, plmnIDBytes []byte, cellID *e2sm_rc_pre_ies.BitString) (*e2sm_rc_pre_ies.E2SmRcPreControlHeader, error) {

	e2smRcPreFormat1 := e2sm_rc_pre_ies.E2SmRcPreControlHeaderFormat1{
		Cgi: &e2sm_rc_pre_ies.CellGlobalId{
			CellGlobalId: &e2sm_rc_pre_ies.CellGlobalId_EUtraCgi{
				EUtraCgi: &e2sm_rc_pre_ies.Eutracgi{
					PLmnIdentity: &e2sm_rc_pre_ies.PlmnIdentity{
						Value: plmnIDBytes,
					},
					EUtracellIdentity: &e2sm_rc_pre_ies.EutracellIdentity{
						Value: cellID,
					},
				},
			},
		},
		RcCommand: e2sm_rc_pre_ies.RcPreCommand_RC_PRE_COMMAND_SET_PARAMETERS,
		RicControlMessagePriority: &e2sm_rc_pre_ies.RicControlMessagePriority{
			Value: controlMessagePriority,
		},
	}
	e2smRcPrePdu := e2sm_rc_pre_ies.E2SmRcPreControlHeader{
		E2SmRcPreControlHeader: &e2sm_rc_pre_ies.E2SmRcPreControlHeader_ControlHeaderFormat1{
			ControlHeaderFormat1: &e2smRcPreFormat1,
		},
	}

	if err := e2smRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &e2smRcPrePdu, nil
}
