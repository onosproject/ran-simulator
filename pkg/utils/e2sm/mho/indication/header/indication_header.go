// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package header

import (
	"fmt"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"

	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/v1/e2sm-mho"
)

// Header indication header for mho service model
type Header struct {
	plmnID            ransimtypes.Uint24
	eutraCellIdentity uint64
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

// WithEutracellIdentity sets eutraCellIdentity
func WithEutracellIdentity(eutraCellIdentity uint64) func(header *Header) {
	return func(header *Header) {
		header.eutraCellIdentity = eutraCellIdentity
	}
}

// Build builds indication header for mho service model
func (header *Header) Build() (*e2sm_mho.E2SmMhoIndicationHeader, error) {
	E2SmMhoPdu := e2sm_mho.E2SmMhoIndicationHeader{
		E2SmMhoIndicationHeader: &e2sm_mho.E2SmMhoIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2sm_mho.E2SmMhoIndicationHeaderFormat1{
				Cgi: &e2sm_mho.CellGlobalId{
					CellGlobalId: &e2sm_mho.CellGlobalId_EUtraCgi{
						EUtraCgi: &e2sm_mho.Eutracgi{
							PLmnIdentity: &e2sm_mho.PlmnIdentity{
								Value: header.plmnID.ToBytes(),
							},
							EUtracellIdentity: &e2sm_mho.EutracellIdentity{
								Value: &e2sm_mho.BitString{
									Value: header.eutraCellIdentity, //uint64
									Len:   28,                       //uint32
								},
							},
						},
					},
				},
			},
		},
	}

	if err := E2SmMhoPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmMhoPDU %s", err.Error())
	}
	return &E2SmMhoPdu, nil

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
