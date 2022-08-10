// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package format2

import (
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/pdubuilder"
	e2smrcsm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/servicemodel"
	e2smcommonies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-common-ies"
	e2smrcies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-rc-ies"
	"google.golang.org/protobuf/proto"
)

// Header indication header format 1
type Header struct {
	ueID               *e2smcommonies.Ueid
	ricInsertStyleType int32
	insertIndicationID int32
}

// NewIndicationHeader creates a new indication header
func NewIndicationHeader(options ...func(header *Header)) *Header {
	header := &Header{}
	for _, option := range options {
		option(header)
	}

	return header
}

// WithUEID sets UE ID
func WithUEID(ueID *e2smcommonies.Ueid) func(header *Header) {
	return func(header *Header) {
		header.ueID = ueID
	}
}

// WithRICInsertStyleType sets RIC Insert Style Type
func WithRICInsertStyleType(style int32) func(header *Header) {
	return func(header *Header) {
		header.ricInsertStyleType = style
	}
}

// WithInsertIndicationID sets Insert Indication ID
func WithInsertIndicationID(id int32) func(header *Header) {
	return func(header *Header) {
		header.insertIndicationID = id
	}
}

// Build builds indication header format 2
func (h *Header) Build() (*e2smrcies.E2SmRcIndicationHeader, error) {
	header, err := pdubuilder.CreateE2SmRcIndicationHeaderFormat2(h.ueID, h.ricInsertStyleType, h.insertIndicationID)
	if err != nil {
		return nil, err
	}
	return header, nil
}

// ToAsn1Bytes converts Header to ASN.1 bytes
func (h *Header) ToAsn1Bytes() ([]byte, error) {
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
