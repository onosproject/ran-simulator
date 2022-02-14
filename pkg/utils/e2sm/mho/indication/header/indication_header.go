// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package header

import (
	"fmt"
	e2smmhosm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/servicemodel"
	e2smv2ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-v2-ies"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"google.golang.org/protobuf/proto"

	mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("sm", "mho")

// Header indication header for mho service model
type Header struct {
	plmnID         ransimtypes.Uint24
	nrCellIdentity []byte
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

// WithNrcellIdentity sets nrCellIdentity
func WithNrcellIdentity(nrCellIdentity []byte) func(header *Header) {
	return func(header *Header) {
		header.nrCellIdentity = nrCellIdentity
	}
}

// Build builds indication header for mho service model
func (header *Header) Build() (*mho.E2SmMhoIndicationHeader, error) {
	E2SmMhoPdu := mho.E2SmMhoIndicationHeader{
		E2SmMhoIndicationHeader: &mho.E2SmMhoIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &mho.E2SmMhoIndicationHeaderFormat1{
				Cgi: &e2smv2ies.Cgi{
					Cgi: &e2smv2ies.Cgi_NRCgi{
						NRCgi: &e2smv2ies.NrCgi{
							PLmnidentity: &e2smv2ies.PlmnIdentity{
								Value: header.plmnID.ToBytes(),
							},
							NRcellIdentity: &e2smv2ies.NrcellIdentity{
								Value: &asn1.BitString{
									//ToDo - should be of type []byte
									Value: header.nrCellIdentity, //uint64
									Len:   36,                    //uint32
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

// MhoToAsn1Bytes converts header to asn1 bytes
func (header *Header) MhoToAsn1Bytes() ([]byte, error) {
	log.Debug("MhoToAsn1Bytes")
	// Creating an indication header
	indicationHeader, err := header.Build()
	if err != nil {
		return nil, err
	}

	indicationHeaderProtoBytes, err := proto.Marshal(indicationHeader)
	if err != nil {
		return nil, err
	}

	var mhoServiceModel e2smmhosm.MhoServiceModel
	indicationHeaderAsn1Bytes, err := mhoServiceModel.IndicationHeaderProtoToASN1(indicationHeaderProtoBytes)

	if err != nil {
		return nil, err
	}
	return indicationHeaderAsn1Bytes, nil
}
