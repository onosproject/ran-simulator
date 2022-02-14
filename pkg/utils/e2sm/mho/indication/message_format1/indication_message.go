// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package messageformat1

import (
	"fmt"
	e2smmhosm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/servicemodel"
	mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"
	e2smv2ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-v2-ies"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"google.golang.org/protobuf/proto"
)

// Message indication message fields for MHO service model
type Message struct {
	ueID       int64
	MeasReport []*mho.E2SmMhoMeasurementReportItem
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
func WithUeID(ueID int64) func(message *Message) {
	return func(message *Message) {
		message.ueID = ueID
	}
}

// WithMeasReport sets measReport
func WithMeasReport(measReport []*mho.E2SmMhoMeasurementReportItem) func(message *Message) {
	return func(message *Message) {
		message.MeasReport = measReport
	}
}

// Build builds indication message for MHO service model
func (message *Message) Build() (*mho.E2SmMhoIndicationMessage, error) {
	e2SmIndicationMsg := mho.E2SmMhoIndicationMessage_IndicationMessageFormat1{
		IndicationMessageFormat1: &mho.E2SmMhoIndicationMessageFormat1{
			UeId: &e2smv2ies.Ueid{
				Ueid: &e2smv2ies.Ueid_GNbUeid{
					GNbUeid: &e2smv2ies.UeidGnb{
						AmfUeNgapId: &e2smv2ies.AmfUeNgapId{
							Value: message.ueID,
						},
						// ToDo - move out GUAMI hardcoding
						Guami: &e2smv2ies.Guami{
							PLmnidentity: &e2smv2ies.PlmnIdentity{
								Value: []byte{0xAA, 0xBB, 0xCC},
							},
							AMfregionId: &e2smv2ies.AmfregionId{
								Value: &asn1.BitString{
									Value: []byte{0xDD},
									Len:   8,
								},
							},
							AMfsetId: &e2smv2ies.AmfsetId{
								Value: &asn1.BitString{
									Value: []byte{0xCC, 0xC0},
									Len:   10,
								},
							},
							AMfpointer: &e2smv2ies.Amfpointer{
								Value: &asn1.BitString{
									Value: []byte{0xFC},
									Len:   6,
								},
							},
						},
					},
				},
			},
			MeasReport: message.MeasReport,
		},
	}

	E2SmMhoPdu := mho.E2SmMhoIndicationMessage{
		E2SmMhoIndicationMessage: &e2SmIndicationMsg,
	}

	if err := E2SmMhoPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
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
