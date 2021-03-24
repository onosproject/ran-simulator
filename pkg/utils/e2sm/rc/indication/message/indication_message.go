// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package message

import (
	"fmt"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"

	"google.golang.org/protobuf/proto"
)

// Message indication message fields for rc service model
type Message struct {
	plmnID            ransimtypes.Uint24
	eutraCellIdentity uint64
	earfcn            int32
	cellSize          e2smrcpreies.CellSize
	pci               int32
	neighbours        []*e2smrcpreies.Nrt
	pciPool           []*e2smrcpreies.PciRange
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

// WithEutraCellIdentity sets eutraCellIdentity
func WithEutraCellIdentity(eutraCellIdentity uint64) func(message *Message) {
	return func(message *Message) {
		message.eutraCellIdentity = eutraCellIdentity
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

// WithPciPool sets pciPool
func WithPciPool(pciPool []*e2smrcpreies.PciRange) func(message *Message) {
	return func(message *Message) {
		message.pciPool = pciPool
	}
}

// Build builds indication message for RC service model
func (message *Message) Build() (*e2smrcpreies.E2SmRcPreIndicationMessage, error) {
	e2SmIindicationMsg := e2smrcpreies.E2SmRcPreIndicationMessage_IndicationMessageFormat1{
		IndicationMessageFormat1: &e2smrcpreies.E2SmRcPreIndicationMessageFormat1{
			Neighbors: make([]*e2smrcpreies.Nrt, 0),
			PciPool:   make([]*e2smrcpreies.PciRange, 0),
		},
	}

	e2SmIindicationMsg.IndicationMessageFormat1.Cgi = &e2smrcpreies.CellGlobalId{
		CellGlobalId: &e2smrcpreies.CellGlobalId_EUtraCgi{
			EUtraCgi: &e2smrcpreies.Eutracgi{
				PLmnIdentity: &e2smrcpreies.PlmnIdentity{
					Value: message.plmnID.ToBytes(),
				},
				EUtracellIdentity: &e2smrcpreies.EutracellIdentity{
					Value: &e2smrcpreies.BitString{
						Value: message.eutraCellIdentity, //uint64
						Len:   28,                        //uint32
					},
				},
			},
		},
	}
	e2SmIindicationMsg.IndicationMessageFormat1.DlEarfcn = &e2smrcpreies.Earfcn{
		Value: message.earfcn,
	}

	e2SmIindicationMsg.IndicationMessageFormat1.CellSize = message.cellSize
	e2SmIindicationMsg.IndicationMessageFormat1.Pci = &e2smrcpreies.Pci{
		Value: message.pci,
	}

	e2SmIindicationMsg.IndicationMessageFormat1.PciPool = message.pciPool
	e2SmIindicationMsg.IndicationMessageFormat1.Neighbors = message.neighbours

	E2SmRcPrePdu := e2smrcpreies.E2SmRcPreIndicationMessage{
		E2SmRcPreIndicationMessage: &e2SmIindicationMsg,
	}

	if err := E2SmRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &E2SmRcPrePdu, nil

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
