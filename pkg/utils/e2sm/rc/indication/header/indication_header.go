// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package header

import (
	"fmt"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"

	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

// Header indication header for rc service model
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

// Build builds indication header for rc service model
func (header *Header) Build() (*e2smrcpreies.E2SmRcPreIndicationHeader, error) {
	E2SmRcPrePdu := e2smrcpreies.E2SmRcPreIndicationHeader{
		E2SmRcPreIndicationHeader: &e2smrcpreies.E2SmRcPreIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2smrcpreies.E2SmRcPreIndicationHeaderFormat1{
				Cgi: &e2smrcpreies.CellGlobalId{
					CellGlobalId: &e2smrcpreies.CellGlobalId_EUtraCgi{
						EUtraCgi: &e2smrcpreies.Eutracgi{
							PLmnIdentity: &e2smrcpreies.PlmnIdentity{
								Value: header.plmnID.ToBytes(),
							},
							EUtracellIdentity: &e2smrcpreies.EutracellIdentity{
								Value: &e2smrcpreies.BitString{
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

	if err := E2SmRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmRcPrePDU %s", err.Error())
	}
	return &E2SmRcPrePdu, nil

}

// ToAsn1Bytes converts header to asn1 bytes
func (header *Header) ToAsn1Bytes(modelPlugin modelplugins.ModelPlugin) ([]byte, error) {
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
