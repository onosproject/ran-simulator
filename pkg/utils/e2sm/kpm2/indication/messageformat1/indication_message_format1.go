// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package messageformat1

import (
	e2smkpmv2sm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/servicemodel"
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
	"google.golang.org/protobuf/proto"
)

// Message indication message format 1 fields for kpm v2 service model
type Message struct {
	subscriptionID int64
	cellObjID      string
	granularity    uint32
	measInfoList   *e2smkpmv2.MeasurementInfoList
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
func WithGranularity(granularity uint32) func(msg *Message) {
	return func(msg *Message) {
		msg.granularity = granularity
	}
}

// WithMeasData sets measurements data
func WithMeasData(measData *e2smkpmv2.MeasurementData) func(msg *Message) {
	return func(msg *Message) {
		msg.measData = measData
	}
}

// WithMeasInfoList sets measurement info list
func WithMeasInfoList(measInfoList *e2smkpmv2.MeasurementInfoList) func(msg *Message) {
	return func(msg *Message) {
		msg.measInfoList = measInfoList
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

	var kpm2ServiceModel e2smkpmv2sm.Kpm2ServiceModel

	indicationMessageAsn1Bytes, err := kpm2ServiceModel.IndicationMessageProtoToASN1(indicationMessageProtoBytes)
	if err != nil {
		return nil, err
	}

	return indicationMessageAsn1Bytes, nil

}

// Build builds indication message format 1 for kpm service model
func (message *Message) Build() (*e2smkpmv2.E2SmKpmIndicationMessage, error) {
	e2SmKpmPdu := e2smkpmv2.E2SmKpmIndicationMessage{
		IndicationMessageFormats: &e2smkpmv2.IndicationMessageFormats{
			E2SmKpmIndicationMessage: &e2smkpmv2.IndicationMessageFormats_IndicationMessageFormat1{
				IndicationMessageFormat1: &e2smkpmv2.E2SmKpmIndicationMessageFormat1{
					SubscriptId: &e2smkpmv2.SubscriptionId{
						Value: message.subscriptionID,
					},
					CellObjId: &e2smkpmv2.CellObjectId{
						Value: message.cellObjID,
					},
					GranulPeriod: &e2smkpmv2.GranularityPeriod{
						Value: int64(message.granularity),
					},
					MeasInfoList: message.measInfoList,
					MeasData:     message.measData,
				},
			},
		},
	}

	// FIXME: Add back when ready
	//if err := e2SmKpmPdu.Validate(); err != nil {
	//	return nil, errors.New(errors.Invalid, err.Error())
	//}

	return &e2SmKpmPdu, nil
}
