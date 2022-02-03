// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package indication

import (
	"fmt"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"

	e2smkpmies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
)

// Header indication header for kpm service model
type Header struct {
	plmnID      ransimtypes.Uint24
	gNbCuUpID   int64
	gNbDuID     int64
	plmnIDnrcgi ransimtypes.Uint24
	sst         string
	sd          string
	fiveQi      int32
	qCi         int32
	gnbID       uint64
}

// NewIndicationHeader creates a new indication header
func NewIndicationHeader(options ...func(header *Header)) *Header {
	header := &Header{}
	for _, option := range options {
		option(header)
	}

	return header
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID ransimtypes.Uint24) func(header *Header) {
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
func WithPlmnIDnrcgi(plmnIDnrcgi ransimtypes.Uint24) func(header *Header) {
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

// ToAsn1Bytes converts header to asn1 bytes
func (header *Header) ToAsn1Bytes(modelPlugin modelplugins.ServiceModel) ([]byte, error) {
	// Creating an indication header
	indicationHeader, err := header.Build()
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

// Build builds kpm indication header message
func (header *Header) Build() (*e2smkpmies.E2SmKpmIndicationHeader, error) {
	e2SmKpmPdu := &e2smkpmies.E2SmKpmIndicationHeader{
		E2SmKpmIndicationHeader: &e2smkpmies.E2SmKpmIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2smkpmies.E2SmKpmIndicationHeaderFormat1{
				IdGlobalKpmnodeId: &e2smkpmies.GlobalKpmnodeId{
					GlobalKpmnodeId: &e2smkpmies.GlobalKpmnodeId_GNb{
						GNb: &e2smkpmies.GlobalKpmnodeGnbId{
							GlobalGNbId: &e2smkpmies.GlobalgNbId{
								PlmnId: &e2smkpmies.PlmnIdentity{
									Value: header.plmnID.ToBytes(),
								},
								GnbId: &e2smkpmies.GnbIdChoice{
									GnbIdChoice: &e2smkpmies.GnbIdChoice_GnbId{
										GnbId: &e2smkpmies.BitString{
											Value: header.gnbID, //uint64
											Len:   22,           //uint32
										},
									},
								},
							},
							GNbCuUpId: &e2smkpmies.GnbCuUpId{
								Value: header.gNbCuUpID, //int64
							},
							GNbDuId: &e2smkpmies.GnbDuId{
								Value: header.gNbDuID, //int64
							},
						},
					},
				},
				NRcgi: &e2smkpmies.Nrcgi{
					PLmnIdentity: &e2smkpmies.PlmnIdentity{
						Value: header.plmnIDnrcgi.ToBytes(),
					},
					NRcellIdentity: &e2smkpmies.NrcellIdentity{
						Value: &e2smkpmies.BitString{
							Value: header.gnbID, //uint64
							Len:   36,           //uint32
						},
					},
				},
				PLmnIdentity: &e2smkpmies.PlmnIdentity{
					Value: header.plmnID.ToBytes(),
				},
				SliceId: &e2smkpmies.Snssai{
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
