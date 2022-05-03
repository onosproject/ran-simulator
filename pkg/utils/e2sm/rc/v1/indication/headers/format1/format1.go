// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package format1

import (
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/pdubuilder"
	e2smrcsm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/servicemodel"
	e2smrcies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-rc-ies"
	"google.golang.org/protobuf/proto"
)

// HeaderFormat1 indication header format 1
type HeaderFormat1 struct {
}

// NewIndicationHeader creates a new indication header
func NewIndicationHeader(options ...func(header *HeaderFormat1)) *HeaderFormat1 {
	header := &HeaderFormat1{}
	for _, option := range options {
		option(header)
	}

	return header
}

func (h *HeaderFormat1) Build() (*e2smrcies.E2SmRcIndicationHeader, error) {
	header, err := pdubuilder.CreateE2SmRcIndicationHeaderFormat1()
	if err != nil {
		return nil, err
	}
	return header, nil

}

// ToAsn1Bytes converts header to asn1 bytes
func (h *HeaderFormat1) ToAsn1Bytes() ([]byte, error) {
	// Creating an indication header
	indicationHeader, err := h.Build()
	if err != nil {
		return nil, err
	}

	indicationHeaderProtoBytes, err := proto.Marshal(indicationHeader)
	if err != nil {
		return nil, err
	}

	var rcServiceModel e2smrcsm.RCServiceModel
	indicationHeaderAsn1Bytes, err := rcServiceModel.IndicationHeaderProtoToASN1(indicationHeaderProtoBytes)

	if err != nil {
		return nil, err
	}
	return indicationHeaderAsn1Bytes, nil
}
