// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm

import (
	"fmt"

	e2sm_kpm_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
)

type IndicationHeader struct {
	plmnID      string
	gNbCuUpID   int64
	gNbDuID     int64
	plmnIDnrcgi string
	sst         string
	sd          string
	fiveQi      int32
	qCi         int32
	gnbId       uint64
}

// NewIndicationHeader creates a new indication header
func NewIndicationHeader(options ...func(header *IndicationHeader)) (*IndicationHeader, error) {
	header := &IndicationHeader{}
	for _, option := range options {
		option(header)
	}

	return header, nil
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID string) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.plmnID = plmnID

	}
}

// WithGnbCuUpID sets gNBCuCpID
func WithGNbCuUpID(gNbCuUpID int64) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.gNbCuUpID = gNbCuUpID

	}
}

// WithGNbDuID sets gNbDuID
func WithGNbDuID(gNbDuID int64) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.gNbDuID = gNbDuID
	}
}

// WithPlmnIDnrcgi sets plmnIDnrcgi
func WithPlmnIDnrcgi(plmnIDnrcgi string) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.plmnIDnrcgi = plmnIDnrcgi
	}
}

// WithSst sets sst
func WithSst(sst string) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.sst = sst
	}
}

// WithSd sets sd
func WithSd(sd string) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.sd = sd
	}
}

// WithFiveQi sets fiveQi
func WithFiveQi(fiveQi int32) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.fiveQi = fiveQi
	}
}

// WithQci sets Qci
func WithQci(qCi int32) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.qCi = qCi
	}
}

// WithGnbId sets E2 global node ID
func WithGnbId(gnbId uint64) func(header *IndicationHeader) {
	return func(header *IndicationHeader) {
		header.gnbId = gnbId
	}
}

// CreateIndicationHeader creates indication header for kpm service model
func CreateIndicationHeader(header *IndicationHeader) (*e2sm_kpm_ies.E2SmKpmIndicationHeader, error) {
	e2SmKpmPdu := &e2sm_kpm_ies.E2SmKpmIndicationHeader{
		E2SmKpmIndicationHeader: &e2sm_kpm_ies.E2SmKpmIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2sm_kpm_ies.E2SmKpmIndicationHeaderFormat1{
				IdGlobalKpmnodeId: &e2sm_kpm_ies.GlobalKpmnodeId{
					GlobalKpmnodeId: &e2sm_kpm_ies.GlobalKpmnodeId_GNb{
						GNb: &e2sm_kpm_ies.GlobalKpmnodeGnbId{
							GlobalGNbId: &e2sm_kpm_ies.GlobalgNbId{
								PlmnId: &e2sm_kpm_ies.PlmnIdentity{
									Value: []byte(header.plmnID),
								},
								GnbId: &e2sm_kpm_ies.GnbIdChoice{
									GnbIdChoice: &e2sm_kpm_ies.GnbIdChoice_GnbId{
										GnbId: &e2sm_kpm_ies.BitString{
											Value: header.gnbId, //uint64
											Len:   22,           //uint32
										},
									},
								},
							},
							GNbCuUpId: &e2sm_kpm_ies.GnbCuUpId{
								Value: header.gNbCuUpID, //int64
							},
							GNbDuId: &e2sm_kpm_ies.GnbDuId{
								Value: header.gNbDuID, //int64
							},
						},
					},
				},
				NRcgi: &e2sm_kpm_ies.Nrcgi{
					PLmnIdentity: &e2sm_kpm_ies.PlmnIdentity{
						Value: []byte(header.plmnIDnrcgi),
					},
					NRcellIdentity: &e2sm_kpm_ies.NrcellIdentity{
						Value: &e2sm_kpm_ies.BitString{
							Value: header.gnbId, //uint64
							Len:   36,           //uint32
						},
					},
				},
				PLmnIdentity: &e2sm_kpm_ies.PlmnIdentity{
					Value: []byte(header.plmnID),
				},
				SliceId: &e2sm_kpm_ies.Snssai{
					SSt: []byte(header.sst),
					SD:  []byte(header.sd),
				},
				FiveQi: header.fiveQi, //int32
				Qci:    header.qCi,    //int32
			},
		},
	}
	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmPDU %s", err.Error())
	}

	return e2SmKpmPdu, nil
}
