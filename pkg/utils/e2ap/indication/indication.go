// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package indication

import (
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
)

// Indication indication data struct
type Indication struct {
	reqID             int32
	ricInstanceID     int32
	ranFuncID         int32
	indicationHeader  []byte
	indicationMessage []byte
	ricCallProcessID  []byte
	// TODO add ric action ID and ric indication sn
}

// NewIndication creates a new indication
func NewIndication(options ...func(*Indication)) *Indication {
	indication := &Indication{}

	for _, option := range options {
		option(indication)
	}

	return indication

}

// WithRequestID sets request ID
func WithRequestID(reqID int32) func(*Indication) {
	return func(indication *Indication) {
		indication.reqID = reqID
	}
}

// WithRanFuncID sets ran function ID
func WithRanFuncID(ranFuncID int32) func(*Indication) {
	return func(indication *Indication) {
		indication.ranFuncID = ranFuncID
	}
}

// WithRicInstanceID sets ric instance ID
func WithRicInstanceID(ricInstanceID int32) func(*Indication) {
	return func(indication *Indication) {
		indication.ricInstanceID = ricInstanceID
	}
}

// WithIndicationHeader sets indication header
func WithIndicationHeader(indicationHeader []byte) func(*Indication) {
	return func(indication *Indication) {
		indication.indicationHeader = indicationHeader
	}
}

// WithIndicationMessage sets indication message
func WithIndicationMessage(indicationMessage []byte) func(*Indication) {
	return func(indication *Indication) {
		indication.indicationMessage = indicationMessage
	}
}

// WithRicCallProcessID sets RIC call process ID
func WithRicCallProcessID(ricCallProcessID []byte) func(*Indication) {
	return func(indication *Indication) {
		indication.ricCallProcessID = ricCallProcessID
	}
}

// Build builds e2ap indication message
func (indication *Indication) Build() (e2Indication *e2appducontents.Ricindication, err error) {
	rrID := types.RicRequest{
		RequestorID: types.RicRequestorID(indication.reqID),
		InstanceID:  types.RicInstanceID(indication.ricInstanceID),
	}
	ricIndication := &e2appducontents.Ricindication{
		ProtocolIes: make([]*e2appducontents.RicindicationIes, 0),
	}
	ricIndication.SetRicRequestID(rrID).SetRanFunctionID(types.RanFunctionID(indication.ranFuncID)).
		SetRicActionID(2).
		SetRicIndicationSN(3).SetRicIndicationType(e2apies.RicindicationType_RICINDICATION_TYPE_REPORT).
		SetRicIndicationHeader(indication.indicationHeader).SetRicIndicationMessage(indication.indicationMessage).
		SetRicCallProcessID(indication.ricCallProcessID)

	return ricIndication, nil
}
