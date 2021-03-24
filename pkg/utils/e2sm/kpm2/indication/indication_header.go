// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package indication

import (
	"github.com/onosproject/onos-lib-go/pkg/errors"

	"google.golang.org/protobuf/proto"

	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	// "google.golang.org/protobuf/proto"

	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
)

// Header indication header for kpm service model
type Header struct {
	timeStamp         []byte
	fileFormatVersion string
	senderName        string
	senderType        string
	vendorName        string
	globalKpmNodeID   *e2smkpmv2.GlobalKpmnodeId
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

// WithFileFormatVersion sets file format version
func WithFileFormatVersion(fileFormatVersion string) func(header *Header) {
	return func(header *Header) {
		header.fileFormatVersion = fileFormatVersion
	}
}

// WithSenderName sets the sender name
func WithSenderName(senderName string) func(header *Header) {
	return func(header *Header) {
		header.senderName = senderName
	}
}

// WithSenderType sets the sender type
func WithSenderType(senderType string) func(header *Header) {
	return func(header *Header) {
		header.senderType = senderType
	}
}

// WithVendorName sets the vendor name
func WithVendorName(vendorName string) func(header *Header) {
	return func(header *Header) {
		header.vendorName = vendorName
	}
}

// WithGlobalKpmNodeID sets the global kpm node ID
func WithGlobalKpmNodeID(globalKpmNodeID *e2smkpmv2.GlobalKpmnodeId) func(header *Header) {
	return func(header *Header) {
		header.globalKpmNodeID = globalKpmNodeID
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

// Build builds kpm v2 indication header message
func (header *Header) Build() (*e2smkpmv2.E2SmKpmIndicationHeader, error) {
	e2SmKpmPdu := e2smkpmv2.E2SmKpmIndicationHeader{
		E2SmKpmIndicationHeader: &e2smkpmv2.E2SmKpmIndicationHeader_IndicationHeaderFormat1{
			IndicationHeaderFormat1: &e2smkpmv2.E2SmKpmIndicationHeaderFormat1{
				ColletStartTime: &e2smkpmv2.TimeStamp{
					Value: header.timeStamp,
				},
				FileFormatversion: header.fileFormatVersion,
				SenderName:        header.senderName,
				SenderType:        header.senderType,
				VendorName:        header.vendorName,
				KpmNodeId:         header.globalKpmNodeID,
			},
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, errors.New(errors.Invalid, err.Error())
	}
	return &e2SmKpmPdu, nil
}
