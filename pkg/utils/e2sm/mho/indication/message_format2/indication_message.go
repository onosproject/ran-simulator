// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package messageformat2

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
	ueID        int64
	RrcStatus   mho.Rrcstatus
	plmnID      *e2smv2ies.PlmnIdentity
	amfRegionID *e2smv2ies.AmfregionId
	amfSetID    *e2smv2ies.AmfsetId
	amfPointer  *e2smv2ies.Amfpointer
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

// WithRrcStatus sets RrcStatus
func WithRrcStatus(rrcStatus mho.Rrcstatus) func(message *Message) {
	return func(message *Message) {
		message.RrcStatus = rrcStatus
	}
}

// WithGuami sets GUAMI
func WithGuami(plmnid uint64, amfRegionID uint32, amfSetID uint32, amfPointer uint32) func(message *Message) {
	return func(message *Message) {
		message.plmnID = &e2smv2ies.PlmnIdentity{
			Value: []byte{byte(plmnid & 0xFF0000 >> 16), byte(plmnid & 0xFF00 >> 8), byte(plmnid & 0xFF)},
		}

		message.amfRegionID = &e2smv2ies.AmfregionId{
			Value: &asn1.BitString{
				Len:   8,
				Value: []byte{byte(amfRegionID & 0xFF)},
			},
		}

		message.amfSetID = &e2smv2ies.AmfsetId{
			Value: &asn1.BitString{
				Len:   10,
				Value: []byte{byte((amfSetID << 6) & 0xFF00 >> 8), byte((amfSetID << 6) & 0xFF)},
			},
		}

		message.amfPointer = &e2smv2ies.Amfpointer{
			Value: &asn1.BitString{
				Len:   6,
				Value: []byte{byte((amfPointer << 2) & 0xFF)},
			},
		}
	}
}

// Build builds indication message for MHO service model
func (message *Message) Build() (*mho.E2SmMhoIndicationMessage, error) {
	e2SmIndicationMsg := mho.E2SmMhoIndicationMessage_IndicationMessageFormat2{
		IndicationMessageFormat2: &mho.E2SmMhoIndicationMessageFormat2{
			UeId: &e2smv2ies.Ueid{
				Ueid: &e2smv2ies.Ueid_GNbUeid{
					GNbUeid: &e2smv2ies.UeidGnb{
						AmfUeNgapId: &e2smv2ies.AmfUeNgapId{
							Value: message.ueID,
						},
						// ToDo - move out GUAMI hardcoding
						Guami: &e2smv2ies.Guami{
							PLmnidentity: message.plmnID,
							AMfregionId:  message.amfRegionID,
							AMfsetId:     message.amfSetID,
							AMfpointer:   message.amfPointer,
						},
					},
				},
			},
			RrcStatus: message.RrcStatus,
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
