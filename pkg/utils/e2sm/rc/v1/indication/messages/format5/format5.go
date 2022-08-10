// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package format5

import (
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/pdubuilder"
	e2smrcsm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/servicemodel"
	e2smrcies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-rc-ies"
	"google.golang.org/protobuf/proto"
)

// Message indication message fields for rc service model
type Message struct {
	indicationMessageItems []*e2smrcies.E2SmRcIndicationMessageFormat5Item
}

// NewIndicationMessage creates a new indication message
func NewIndicationMessage(options ...func(msg *Message)) *Message {
	msg := &Message{}
	for _, option := range options {
		option(msg)
	}

	return msg
}

// WithMessageItems sets indication message items
func WithMessageItems(indicationMessageItems []*e2smrcies.E2SmRcIndicationMessageFormat5Item) func(message *Message) {
	return func(message *Message) {
		message.indicationMessageItems = indicationMessageItems
	}
}

// Build builds indication message for RC service model
func (message *Message) Build() (*e2smrcies.E2SmRcIndicationMessage, error) {
	indicationMessage, err := pdubuilder.CreateE2SmRcIndicationMessageFormat5(message.indicationMessageItems)
	if err != nil {
		return nil, err
	}
	return indicationMessage, nil

}

// ToAsn1Bytes converts to Asn1 bytes
func (message *Message) ToAsn1Bytes() ([]byte, error) {
	indicationMessage, err := message.Build()
	if err != nil {
		return nil, err
	}
	indicationMessageProtoBytes, err := proto.Marshal(indicationMessage)
	if err != nil {
		return nil, err
	}

	var rcServiceModel e2smrcsm.RCServiceModel
	indicationMessageAsn1Bytes, err := rcServiceModel.IndicationMessageProtoToASN1(indicationMessageProtoBytes)
	if err != nil {
		return nil, err
	}

	return indicationMessageAsn1Bytes, nil
}
