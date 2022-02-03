// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package header

import (
	e2smrcpresm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/servicemodel"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"

	"github.com/onosproject/ran-simulator/pkg/utils"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"google.golang.org/protobuf/proto"

	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/v2/e2sm-rc-pre-v2-go"
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

// WithNRcellIdentity sets NRcellIdentity
func WithNRcellIdentity(nRcellIdentity uint64) func(header *Header) {
	return func(header *Header) {
		header.eutraCellIdentity = nRcellIdentity
	}
}

// Build builds indication header for rc service model
func (header *Header) Build() (*e2smrcpreies.E2SmRcPreIndicationHeader, error) {
	E2SmRcPrePdu := e2smrcpreies.E2SmRcPreIndicationHeader{
		E2SmRcPreIndicationHeader: &e2smrcpreies.E2SmRcPreIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2smrcpreies.E2SmRcPreIndicationHeaderFormat1{
				Cgi: &e2smrcpreies.CellGlobalId{
					CellGlobalId: &e2smrcpreies.CellGlobalId_NrCgi{
						NrCgi: &e2smrcpreies.Nrcgi{
							PLmnIdentity: &e2smrcpreies.PlmnIdentity{
								Value: header.plmnID.ToBytes(),
							},
							NRcellIdentity: &e2smrcpreies.NrcellIdentity{
								Value: &asn1.BitString{
									Value: utils.Uint64ToBitString(header.eutraCellIdentity, 36),
									Len:   36,
								},
							},
						},
					},
				},
			},
		},
	}

	//ToDo - return it back once the Validation is functional again
	//if err := E2SmRcPrePdu.Validate(); err != nil {
	//	return nil, fmt.Errorf("error validating E2SmRcPrePDU %s", err.Error())
	//}
	return &E2SmRcPrePdu, nil

}

// ToAsn1Bytes converts header to asn1 bytes
func (header *Header) ToAsn1Bytes() ([]byte, error) {
	// Creating an indication header
	indicationHeader, err := header.Build()
	if err != nil {
		return nil, err
	}

	indicationHeaderProtoBytes, err := proto.Marshal(indicationHeader)
	if err != nil {
		return nil, err
	}

	var rcPreServiceModel e2smrcpresm.RcPreServiceModel
	indicationHeaderAsn1Bytes, err := rcPreServiceModel.IndicationHeaderProtoToASN1(indicationHeaderProtoBytes)

	if err != nil {
		return nil, err
	}
	return indicationHeaderAsn1Bytes, nil
}
