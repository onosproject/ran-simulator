// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package message

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	e2smrcpresm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/servicemodel"
	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/v2/e2sm-rc-pre-v2-go"
	"google.golang.org/protobuf/proto"
)

// Message indication message fields for rc service model
type Message struct {
	plmnID     ransimtypes.Uint24
	earfcn     int32
	cellSize   e2smrcpreies.CellSize
	pci        int32
	neighbours []*e2smrcpreies.Nrt
}

// NewIndicationMessage creates a new indication message
func NewIndicationMessage(options ...func(msg *Message)) *Message {
	msg := &Message{}
	for _, option := range options {
		option(msg)
	}

	return msg
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID ransimtypes.Uint24) func(message *Message) {
	return func(message *Message) {
		message.plmnID = plmnID

	}
}

// WithEarfcn sets 	earfcn
func WithEarfcn(earfcn int32) func(message *Message) {
	return func(message *Message) {
		message.earfcn = earfcn
	}
}

// WithCellSize sets cell size
func WithCellSize(cellSize e2smrcpreies.CellSize) func(message *Message) {
	return func(message *Message) {
		message.cellSize = cellSize
	}
}

// WithPci sets pci
func WithPci(pci int32) func(message *Message) {
	return func(message *Message) {
		message.pci = pci
	}
}

// WithNeighbours sets neighbours
func WithNeighbours(neighbours []*e2smrcpreies.Nrt) func(message *Message) {
	return func(message *Message) {
		message.neighbours = neighbours
	}
}

// Build builds indication message for RC service model
func (message *Message) Build() (*e2smrcpreies.E2SmRcPreIndicationMessage, error) {
	e2SmIindicationMsg := e2smrcpreies.E2SmRcPreIndicationMessage_IndicationMessageFormat1{
		IndicationMessageFormat1: &e2smrcpreies.E2SmRcPreIndicationMessageFormat1{
			Neighbors: make([]*e2smrcpreies.Nrt, 0),
		},
	}

	e2SmIindicationMsg.IndicationMessageFormat1.DlArfcn = &e2smrcpreies.Arfcn{
		Arfcn: &e2smrcpreies.Arfcn_EArfcn{
			EArfcn: &e2smrcpreies.Earfcn{
				Value: message.earfcn,
			},
		},
	}

	e2SmIindicationMsg.IndicationMessageFormat1.CellSize = message.cellSize
	e2SmIindicationMsg.IndicationMessageFormat1.Pci = &e2smrcpreies.Pci{
		Value: message.pci,
	}

	e2SmIindicationMsg.IndicationMessageFormat1.Neighbors = message.neighbours

	E2SmRcPrePdu := e2smrcpreies.E2SmRcPreIndicationMessage{
		E2SmRcPreIndicationMessage: &e2SmIindicationMsg,
	}

	//ToDo - return it back once the Validation is functional again
	//if err := E2SmRcPrePdu.Validate(); err != nil {
	//	return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	//}
	return &E2SmRcPrePdu, nil

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

	var rcPreServiceModel e2smrcpresm.RcPreServiceModel
	indicationMessageAsn1Bytes, err := rcPreServiceModel.IndicationMessageProtoToASN1(indicationMessageProtoBytes)
	if err != nil {
		return nil, err
	}

	return indicationMessageAsn1Bytes, nil
}
