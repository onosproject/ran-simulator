// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package messageformat1

import (
	e2smmhosm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/servicemodel"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v1/e2sm-mho-go"
	"google.golang.org/protobuf/proto"
)

// Message indication message fields for MHO service model
type Message struct {
	ueID       string
	MeasReport []*e2sm_mho.E2SmMhoMeasurementReportItem
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

// WithMeasReport sets measReport
func WithMeasReport(measReport []*e2sm_mho.E2SmMhoMeasurementReportItem) func(message *Message) {
	return func(message *Message) {
		message.MeasReport = measReport
	}
}

// Build builds indication message for MHO service model
func (message *Message) Build() (*e2sm_mho.E2SmMhoIndicationMessage, error) {
	e2SmIndicationMsg := e2sm_mho.E2SmMhoIndicationMessage_IndicationMessageFormat1{
		IndicationMessageFormat1: &e2sm_mho.E2SmMhoIndicationMessageFormat1{
			UeId: &e2sm_mho.UeIdentity{
				Value: []byte(message.ueID),
			},
			MeasReport: message.MeasReport,
		},
	}

	E2SmMhoPdu := e2sm_mho.E2SmMhoIndicationMessage{
		E2SmMhoIndicationMessage: &e2SmIndicationMsg,
	}

	//ToDo - return it back once the Validation is functional again
	//if err := E2SmMhoPdu.Validate(); err != nil {
	//	return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	//}
	return &E2SmMhoPdu, nil

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

	var mhoServiceModel e2smmhosm.MhoServiceModel
	indicationMessageAsn1Bytes, err := mhoServiceModel.IndicationMessageProtoToASN1(indicationMessageProtoBytes)
	if err != nil {
		return nil, err
	}

	return indicationMessageAsn1Bytes, nil
}
