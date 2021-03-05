// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

func CreateE2SmRcPreIndicationHeader(plmnIDBytes []byte, cellID *e2sm_rc_pre_ies.BitString) (*e2sm_rc_pre_ies.E2SmRcPreIndicationHeader, error) {
	if len(plmnIDBytes) != 3 {
		return nil, fmt.Errorf("error: Plmn ID should be 3 chars")
	}

	E2SmRcPrePdu := e2sm_rc_pre_ies.E2SmRcPreIndicationHeader{
		E2SmRcPreIndicationHeader: &e2sm_rc_pre_ies.E2SmRcPreIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2sm_rc_pre_ies.E2SmRcPreIndicationHeaderFormat1{
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
			},
		},
	}

	if err := E2SmRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmRcPrePDU %s", err.Error())
	}
	return &E2SmRcPrePdu, nil
}
