// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package pdubuilder

import (
	"fmt"
	e2sm_kpm_v2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-ies"
)

func CreateE2SmKpmIndicationHeader(timeStamp []byte, fileFormatVersion string, senderName string, senderType string,
	vendorName string, globalKpmNodeID *e2sm_kpm_v2.GlobalKpmnodeId) (*e2sm_kpm_v2.E2SmKpmIndicationHeader, error) {

	if len(timeStamp) != 4 {
		return nil, fmt.Errorf("TimeStamp should be 4 chars")
	}

	e2SmKpmPdu := e2sm_kpm_v2.E2SmKpmIndicationHeader{
		E2SmKpmIndicationHeader: &e2sm_kpm_v2.E2SmKpmIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2sm_kpm_v2.E2SmKpmIndicationHeaderFormat1{
				ColletStartTime: &e2sm_kpm_v2.TimeStamp{
					Value: timeStamp,
				},
				FileFormatversion: fileFormatVersion,
				SenderName:        senderName,
				SenderType:        senderType,
				VendorName:        vendorName,
				KpmNodeId:         globalKpmNodeID,
			},
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmPDU %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}

func CreateGlobalKpmnodeIDgNBID(bs *e2sm_kpm_v2.BitString, plmnID []byte, gnbCuUpID int64,
	gnbDuID int64) (*e2sm_kpm_v2.GlobalKpmnodeId, error) {

	if len(plmnID) != 3 {
		return nil, fmt.Errorf("PlmnID should be 3 chars")
	}

	return &e2sm_kpm_v2.GlobalKpmnodeId{
		GlobalKpmnodeId: &e2sm_kpm_v2.GlobalKpmnodeId_GNb{
			GNb: &e2sm_kpm_v2.GlobalKpmnodeGnbId{
				GlobalGNbId: &e2sm_kpm_v2.GlobalgNbId{
					GnbId: &e2sm_kpm_v2.GnbIdChoice{
						GnbIdChoice: &e2sm_kpm_v2.GnbIdChoice_GnbId{
							GnbId: bs,
						},
					},
					PlmnId: &e2sm_kpm_v2.PlmnIdentity{
						Value: plmnID,
					},
				},
				GNbCuUpId: &e2sm_kpm_v2.GnbCuUpId{
					Value: gnbCuUpID,
				},
				GNbDuId: &e2sm_kpm_v2.GnbDuId{
					Value: gnbDuID,
				},
			},
		},
	}, nil
}

func CreateGlobalKpmnodeIDenGNbID(bsValue uint64, bsLen uint32, plmnID []byte, gnbCuUpID int64,
	gnbDuID int64) (*e2sm_kpm_v2.GlobalKpmnodeId, error) {

	if len(plmnID) != 3 {
		return nil, fmt.Errorf("PlmnID should be 3 chars")
	}

	return &e2sm_kpm_v2.GlobalKpmnodeId{
		GlobalKpmnodeId: &e2sm_kpm_v2.GlobalKpmnodeId_EnGNb{
			EnGNb: &e2sm_kpm_v2.GlobalKpmnodeEnGnbId{
				GlobalGNbId: &e2sm_kpm_v2.GlobalenGnbId{
					GNbId: &e2sm_kpm_v2.EngnbId{
						EngnbId: &e2sm_kpm_v2.EngnbId_GNbId{
							GNbId: &e2sm_kpm_v2.BitString{
								Value: bsValue,
								Len:   bsLen, // should be 22 to 32
							},
						},
					},
					PLmnIdentity: &e2sm_kpm_v2.PlmnIdentity{
						Value: plmnID,
					},
				},
				GNbCuUpId: &e2sm_kpm_v2.GnbCuUpId{
					Value: gnbCuUpID,
				},
				GNbDuId: &e2sm_kpm_v2.GnbDuId{
					Value: gnbDuID,
				},
			},
		},
	}, nil
}

func CreateGlobalKpmnodeIDngENbID(bs *e2sm_kpm_v2.BitString, plmnID []byte, shortMacroEnbID *e2sm_kpm_v2.BitString,
	longMacroEnbID *e2sm_kpm_v2.BitString, gnbDuID int64) (*e2sm_kpm_v2.GlobalKpmnodeId, error) {

	if len(plmnID) != 3 {
		return nil, fmt.Errorf("PlmnID should be 3 chars")
	}

	return &e2sm_kpm_v2.GlobalKpmnodeId{
		GlobalKpmnodeId: &e2sm_kpm_v2.GlobalKpmnodeId_NgENb{
			NgENb: &e2sm_kpm_v2.GlobalKpmnodeNgEnbId{
				GlobalNgENbId: &e2sm_kpm_v2.GlobalngeNbId{
					EnbId: &e2sm_kpm_v2.EnbIdChoice{
						EnbIdChoice: &e2sm_kpm_v2.EnbIdChoice_EnbIdMacro{
							EnbIdMacro: bs,
						},
					},
					PlmnId: &e2sm_kpm_v2.PlmnIdentity{
						Value: plmnID,
					},
					ShortMacroENbId: shortMacroEnbID,
					LongMacroENbId:  longMacroEnbID,
				},
				GNbDuId: &e2sm_kpm_v2.GnbDuId{
					Value: gnbDuID,
				},
			},
		},
	}, nil
}

func CreateGlobalKpmnodeIDeNBID(bs *e2sm_kpm_v2.BitString, plmnID []byte) (*e2sm_kpm_v2.GlobalKpmnodeId, error) {

	if len(plmnID) != 3 {
		return nil, fmt.Errorf("PlmnID should be 3 chars")
	}

	return &e2sm_kpm_v2.GlobalKpmnodeId{
		GlobalKpmnodeId: &e2sm_kpm_v2.GlobalKpmnodeId_ENb{
			ENb: &e2sm_kpm_v2.GlobalKpmnodeEnbId{
				GlobalENbId: &e2sm_kpm_v2.GlobalEnbId{
					ENbId: &e2sm_kpm_v2.EnbId{
						EnbId: &e2sm_kpm_v2.EnbId_HomeENbId{
							HomeENbId: bs,
						},
					},
					PLmnIdentity: &e2sm_kpm_v2.PlmnIdentity{
						Value: plmnID,
					},
				},
			},
		},
	}, nil
}
