// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package messageformat2

import (
	"fmt"

	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

// Message indication message format 2 fields for kpm v2 service model
type Message struct {
	subscriptionID int64
	cellObjID      string
	granularity    int32
	measCondUEList *e2smkpmv2.MeasurementCondUeidList
	measData       *e2smkpmv2.MeasurementData
}

// NewIndicationMessage creates a new indication message
func NewIndicationMessage(options ...func(message *Message)) *Message {
	msg := &Message{}
	for _, option := range options {
		option(msg)
	}

	return msg
}

// WithSubscriptionID sets subscription id
func WithSubscriptionID(subscriptionID int64) func(msg *Message) {
	return func(msg *Message) {
		msg.subscriptionID = subscriptionID
	}
}

// WithCellObjID sets cell object ID
func WithCellObjID(cellObjID string) func(msg *Message) {
	return func(msg *Message) {
		msg.cellObjID = cellObjID
	}
}

// WithGranularity sets granularity
func WithGranularity(granularity int32) func(msg *Message) {
	return func(msg *Message) {
		msg.granularity = granularity
	}
}

// WithMeasCondUEList sets measurement ue list
func WithMeasCondUEList(measCondUEList *e2smkpmv2.MeasurementCondUeidList) func(msg *Message) {
	return func(msg *Message) {
		msg.measCondUEList = measCondUEList
	}
}

// WithMeasData sets measurement data
func WithMeasData(measData *e2smkpmv2.MeasurementData) func(msg *Message) {
	return func(msg *Message) {
		msg.measData = measData
	}
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

// Build builds indication message format 2 for kpm v2 service model
func (message *Message) Build() (*e2smkpmv2.E2SmKpmIndicationMessage, error) {
	e2SmKpmPdu := e2smkpmv2.E2SmKpmIndicationMessage{
		E2SmKpmIndicationMessage: &e2smkpmv2.E2SmKpmIndicationMessage_IndicationMessageFormat2{
			IndicationMessageFormat2: &e2smkpmv2.E2SmKpmIndicationMessageFormat2{
				SubscriptId: &e2smkpmv2.SubscriptionId{
					Value: message.subscriptionID,
				},
				CellObjId: &e2smkpmv2.CellObjectId{
					Value: message.cellObjID,
				},
				GranulPeriod: &e2smkpmv2.GranularityPeriod{
					Value: message.granularity,
				},
				MeasCondUeidList: message.measCondUEList,
				MeasData:         message.measData,
			},
		},
	}

	if err := e2SmKpmPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmKpmPDU %s", err.Error())
	}
	return &e2SmKpmPdu, nil
}
