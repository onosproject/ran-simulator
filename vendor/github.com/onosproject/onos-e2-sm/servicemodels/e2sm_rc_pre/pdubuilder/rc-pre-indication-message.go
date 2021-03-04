// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

func CreateE2SmRcPreIndicationMsg(plmnIDBytes []byte, cellID *e2sm_rc_pre_ies.BitString, dlEarfcn int32,
	cellSize e2sm_rc_pre_ies.CellSize, pci int32, pciPool *e2sm_rc_pre_ies.PciRange, neighbors *e2sm_rc_pre_ies.Nrt) (*e2sm_rc_pre_ies.E2SmRcPreIndicationMessage, error) {
	if len(plmnIDBytes) != 3 {
		return nil, fmt.Errorf("error: Plmn ID should be 3 bytes, actual length is %d", len(plmnIDBytes))
	}

	e2SmIindicationMsg := e2sm_rc_pre_ies.E2SmRcPreIndicationMessage_IndicationMessageFormat1{
		IndicationMessageFormat1: &e2sm_rc_pre_ies.E2SmRcPreIndicationMessageFormat1{
			Neighbors: make([]*e2sm_rc_pre_ies.Nrt, 0),
			PciPool:   make([]*e2sm_rc_pre_ies.PciRange, 0),
			DlEarfcn: &e2sm_rc_pre_ies.Earfcn{
				Value: dlEarfcn,
			},
			CellSize: cellSize,
			Pci: &e2sm_rc_pre_ies.Pci{
				Value: pci,
			},
		},
	}

	e2SmIindicationMsg.IndicationMessageFormat1.Cgi = &e2sm_rc_pre_ies.CellGlobalId{
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
	}

	e2SmIindicationMsg.IndicationMessageFormat1.PciPool = append(e2SmIindicationMsg.IndicationMessageFormat1.PciPool, pciPool)
	e2SmIindicationMsg.IndicationMessageFormat1.Neighbors = append(e2SmIindicationMsg.IndicationMessageFormat1.Neighbors, neighbors)

	E2SmRcPrePdu := e2sm_rc_pre_ies.E2SmRcPreIndicationMessage{
		E2SmRcPreIndicationMessage: &e2SmIindicationMsg,
	}

	if err := E2SmRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &E2SmRcPrePdu, nil
}
