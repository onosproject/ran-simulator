// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package indication

import (
	e2smcccsm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/servicemodel"
	e2smccc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/v1/e2sm-ccc-ies"
	"google.golang.org/protobuf/proto"
)

// Message indication message format 1 fields for ccc service model
type Message struct {
	listOfConfigurationsReported *e2smccc.ListOfConfigurationsReported
}

// NewIndicationMessage creates a new indication message
func NewIndicationMessage(options ...func(message *Message)) *Message {
	msg := &Message{}
	for _, option := range options {
		option(msg)
	}

	return msg
}

// WithConfigurationsReported sets measurement info list
func WithConfigurationsReported(listOfConfigurationsReported *e2smccc.ListOfConfigurationsReported) func(msg *Message) {
	return func(msg *Message) {
		msg.listOfConfigurationsReported = listOfConfigurationsReported
	}
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

	var cccServiceModel e2smcccsm.CCCServiceModel

	indicationMessageAsn1Bytes, err := cccServiceModel.IndicationMessageProtoToASN1(indicationMessageProtoBytes)
	if err != nil {
		return nil, err
	}

	return indicationMessageAsn1Bytes, nil
}

// Build builds indication message format 1 for ccc service model
func (message *Message) Build() (*e2smccc.E2SmCCcRIcIndicationMessage, error) {
	e2SmCccPdu := e2smccc.E2SmCCcRIcIndicationMessage{
		IndicationMessageFormat: &e2smccc.IndicationMessageFormat{
			IndicationMessageFormat: &e2smccc.IndicationMessageFormat_E2SmCccIndicationMessageFormat1{
				E2SmCccIndicationMessageFormat1: &e2smccc.E2SmCCcIndicationMessageFormat1{
					ListOfConfigurationStructuresReported: message.listOfConfigurationsReported,
				},
			},
		},
	}

	if err := e2SmCccPdu.Validate(); err != nil {
		return nil, err
	}

	return &e2SmCccPdu, nil
}
