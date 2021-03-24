// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package indication

import (
	e2smkpmies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

// Message indication message fields for kpm v2 service model
type Message struct {
}

// NewIndicationMessage creates a new indication message
func NewIndicationMessage(options ...func(message *Message)) *Message {
	msg := &Message{}
	for _, option := range options {
		option(msg)
	}

	return msg
}

// ToAsn1Bytes converts to Asn1 bytes
func (message *Message) ToAsn1Bytes(modelPlugin modelplugins.ServiceModel) ([]byte, error) {
	indicationMessage, err := message.Build()
	if err != nil {
		return nil, err
	}
	indicationMessageProtoBytes, err := proto.Marshal(indicationMessage)
	if err != nil {
		return nil, err
	}

	indicationMessageAsn1Bytes, err := modelPlugin.IndicationMessageProtoToASN1(indicationMessageProtoBytes)
	if err != nil {
		return nil, err
	}

	return indicationMessageAsn1Bytes, nil

}

// Build builds indication message for kpm service model
func (message *Message) Build() (*e2smkpmies.E2SmKpmIndicationMessage, error) {

	return nil, nil
}
