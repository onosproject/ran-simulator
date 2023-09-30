// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package indication

import (
	e2smcccsm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/servicemodel"
	"google.golang.org/protobuf/proto"

	e2smccc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/v1/e2sm-ccc-ies"
)

// Header indication header for kpm service model
type Header struct {
	timeStamp        []byte
	indicationReason *e2smccc.IndicationReason
}

// NewIndicationHeader creates a new indication header
func NewIndicationHeader(options ...func(header *Header)) *Header {
	header := &Header{}
	for _, option := range options {
		option(header)
	}

	return header
}

// WithTimeStamp sets timestamp
func WithTimeStamp(timeStamp []byte) func(header *Header) {
	return func(header *Header) {
		header.timeStamp = timeStamp
	}
}

// WithGlobalKpmNodeID sets the global kpm node ID
func WithIndicationReason(indicationReason e2smccc.IndicationReason) func(header *Header) {
	return func(header *Header) {
		header.indicationReason = &indicationReason
	}
}

// Build builds ccc indication header message
func (header *Header) Build() (*e2smccc.E2SmCCcRIcIndicationHeader, error) {
	e2SmCccPdu := e2smccc.E2SmCCcRIcIndicationHeader{
		IndicationHeaderFormat: &e2smccc.IndicationHeaderFormat{
			IndicationHeaderFormat: &e2smccc.IndicationHeaderFormat_E2SmCccIndicationHeaderFormat1{
				E2SmCccIndicationHeaderFormat1: &e2smccc.E2SmCCcIndicationHeaderFormat1{
					IndicationReason: *header.indicationReason,
					EventTime:        header.timeStamp,
				},
			},
		},
	}

	if err := e2SmCccPdu.Validate(); err != nil {
		return nil, err
	}
	return &e2SmCccPdu, nil
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
	var cccServiceModel e2smcccsm.CCCServiceModel

	indicationHeaderAsn1Bytes, err := cccServiceModel.IndicationHeaderProtoToASN1(indicationHeaderProtoBytes)

	if err != nil {
		return nil, err
	}
	return indicationHeaderAsn1Bytes, nil
}
