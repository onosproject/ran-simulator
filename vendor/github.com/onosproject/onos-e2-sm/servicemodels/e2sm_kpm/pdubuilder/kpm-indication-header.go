// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
)

func CreateE2SmKpmIndicationHeader(plmnID string, gNbCuUpID int64, gNbDuID int64, plmnIDnrcgi string,
	sst string, sd string, fiveQi int32, qCi int32) (*e2sm_kpm_ies.E2SmKpmIndicationHeader, error) {
	if len(plmnID) != 3 {
		return nil, fmt.Errorf("error: Plmn ID should be 3 chars")
	}
	if len(plmnIDnrcgi) != 3 {
		return nil, fmt.Errorf("error: Plmn ID NRcgi should be 3 chars")
	}
	if len(sst) != 1 {
		return nil, fmt.Errorf("error: SSt should be 1 char")
	}
	if len(sd) != 3 {
		return nil, fmt.Errorf("error: SD should be 3 chars")
	}

	e2SmKpmPdu := e2sm_kpm_ies.E2SmKpmIndicationHeader{
		E2SmKpmIndicationHeader: &e2sm_kpm_ies.E2SmKpmIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2sm_kpm_ies.E2SmKpmIndicationHeaderFormat1{
				IdGlobalKpmnodeId: &e2sm_kpm_ies.GlobalKpmnodeId{
					GlobalKpmnodeId: &e2sm_kpm_ies.GlobalKpmnodeId_GNb{
						GNb: &e2sm_kpm_ies.GlobalKpmnodeGnbId{
							GlobalGNbId: &e2sm_kpm_ies.GlobalgNbId{
								PlmnId: &e2sm_kpm_ies.PlmnIdentity{
									Value: []byte(plmnID),
								},
								GnbId: &e2sm_kpm_ies.GnbIdChoice{
									GnbIdChoice: &e2sm_kpm_ies.GnbIdChoice_GnbId{
										GnbId: &e2sm_kpm_ies.BitString{
											Value: 0x9bcd4, //uint64
											Len:   22,      //uint32
										},
									},
								},
							},
							GNbCuUpId: &e2sm_kpm_ies.GnbCuUpId{
								Value: gNbCuUpID, //int64
							},
							GNbDuId: &e2sm_kpm_ies.GnbDuId{
								Value: gNbDuID, //int64
							},
						},
					},
				},
				NRcgi: &e2sm_kpm_ies.Nrcgi{
					PLmnIdentity: &e2sm_kpm_ies.PlmnIdentity{
						Value: []byte(plmnIDnrcgi),
					},
					NRcellIdentity: &e2sm_kpm_ies.NrcellIdentity{
						Value: &e2sm_kpm_ies.BitString{
							Value: 0x9bcd4abef, //uint64
							Len:   36,          //uint32
						},
					},
				},
				PLmnIdentity: &e2sm_kpm_ies.PlmnIdentity{
					Value: []byte(plmnID),
				},
				SliceId: &e2sm_kpm_ies.Snssai{
					SSt: []byte(sst),
					SD:  []byte(sd),
				},
				FiveQi: fiveQi, //int32
				Qci:    qCi,    //int32
			},
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmPDU %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}
