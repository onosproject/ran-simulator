// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package indication

import (
	"fmt"

	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"

	e2sm_kpm_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
)

// Header indication header for kpm service model
type Header struct {
	plmnID      string
	gNbCuUpID   int64
	gNbDuID     int64
	plmnIDnrcgi string
	sst         string
	sd          string
	fiveQi      int32
	qCi         int32
	gnbID       uint64
}

// NewIndicationHeader creates a new indication header
func NewIndicationHeader(options ...func(header *Header)) (*Header, error) {
	header := &Header{}
	for _, option := range options {
		option(header)
	}

	return header, nil
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID string) func(header *Header) {
	return func(header *Header) {
		header.plmnID = plmnID

	}
}

// WithGnbCuUpID sets gNBCuCpID
func WithGnbCuUpID(gNbCuUpID int64) func(header *Header) {
	return func(header *Header) {
		header.gNbCuUpID = gNbCuUpID

	}
}

// WithGnbDuID sets gNbDuID
func WithGnbDuID(gNbDuID int64) func(header *Header) {
	return func(header *Header) {
		header.gNbDuID = gNbDuID
	}
}

// WithPlmnIDnrcgi sets plmnIDnrcgi
func WithPlmnIDnrcgi(plmnIDnrcgi string) func(header *Header) {
	return func(header *Header) {
		header.plmnIDnrcgi = plmnIDnrcgi
	}
}

// WithSst sets sst
func WithSst(sst string) func(header *Header) {
	return func(header *Header) {
		header.sst = sst
	}
}

// WithSd sets sd
func WithSd(sd string) func(header *Header) {
	return func(header *Header) {
		header.sd = sd
	}
}

// WithFiveQi sets fiveQi
func WithFiveQi(fiveQi int32) func(header *Header) {
	return func(header *Header) {
		header.fiveQi = fiveQi
	}
}

// WithQci sets Qci
func WithQci(qCi int32) func(header *Header) {
	return func(header *Header) {
		header.qCi = qCi
	}
}

// WithGnbID sets E2 global node ID
func WithGnbID(gnbID uint64) func(header *Header) {
	return func(header *Header) {
		header.gnbID = gnbID
	}
}

// CreateIndicationHeaderAsn1Bytes creates ASN.1 bytes from a protobuf encoded indication header
func CreateIndicationHeaderAsn1Bytes(modelPlugin modelplugins.ModelPlugin, header *Header) ([]byte, error) {
	// Creating an indication header
	indicationHeader, err := CreateIndicationHeader(header)
	if err != nil {
		return nil, err
	}

	indicationHeaderProtoBytes, err := proto.Marshal(indicationHeader)
	if err != nil {
		return nil, err
	}

	indicationHeaderAsn1Bytes, err := modelPlugin.IndicationHeaderProtoToASN1(indicationHeaderProtoBytes)

	if err != nil {
		return nil, err
	}
	return indicationHeaderAsn1Bytes, nil
}

// CreateIndicationHeader creates indication header for kpm service model
func CreateIndicationHeader(header *Header) (*e2sm_kpm_ies.E2SmKpmIndicationHeader, error) {
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
											Value: header.gnbID, //uint64
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
							Value: header.gnbID, //uint64
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
