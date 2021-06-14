// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package messageformat2

import (
	"fmt"

	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/v1/e2sm-mho"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"

	"google.golang.org/protobuf/proto"
)

// Message indication message fields for MHO service model
type Message struct {
	ueID      string
	RrcStatus e2sm_mho.Rrcstatus
}

// NewIndicationMessage creates a new indication message
func NewIndicationMessage(options ...func(msg *Message)) *Message {
	msg := &Message{}
	for _, option := range options {
		option(msg)
	}

	return msg
}

// WithUeID sets ueID
func WithUeID(ueID string) func(message *Message) {
	return func(message *Message) {
		message.ueID = ueID
	}
}

// WithRrcStatus sets RrcStatus
func WithRrcStatus(rrcStatus e2sm_mho.Rrcstatus) func(message *Message) {
	return func(message *Message) {
		message.RrcStatus = rrcStatus
	}
}

// Build builds indication message for MHO service model
func (message *Message) Build() (*e2sm_mho.E2SmMhoIndicationMessage, error) {
	e2SmIndicationMsg := e2sm_mho.E2SmMhoIndicationMessage_IndicationMessageFormat2{
		IndicationMessageFormat2: &e2sm_mho.E2SmMhoIndicationMessageFormat2{
			UeId: &e2sm_mho.UeIdentity{
				Value: message.ueID,
			},
			RrcStatus: message.RrcStatus,
		},
	}

	E2SmMhoPdu := e2sm_mho.E2SmMhoIndicationMessage{
		E2SmMhoIndicationMessage: &e2SmIndicationMsg,
	}

	if err := E2SmMhoPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &E2SmMhoPdu, nil

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
