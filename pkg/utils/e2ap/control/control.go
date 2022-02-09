// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package control

import (
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
)

// Control defines required fields for creating control acknowledge and failure responses
type Control struct {
	reqID         int32
	ricInstanceID int32
	ranFuncID     int32
	ricCallPrID   types.RicCallProcessID
	ricCtrlStatus e2apies.RiccontrolStatus
	ricCtrlOut    []byte
	cause         *e2apies.Cause
}

// NewControl creates a new instance of control
func NewControl(options ...func(control *Control)) *Control {
	control := &Control{}

	for _, option := range options {
		option(control)
	}

	return control
}

// GetRanFuncID returns control ran function ID
func (control *Control) GetRanFuncID() int32 {
	return control.ranFuncID
}

// GetRicInstanceID returns control RicInstance ID
func (control *Control) GetRicInstanceID() int32 {
	return control.ricInstanceID
}

// GetReqID returns control request ID
func (control *Control) GetReqID() int32 {
	return control.reqID
}

// WithRequestID sets request ID
func WithRequestID(reqID int32) func(control *Control) {
	return func(control *Control) {
		control.reqID = reqID
	}
}

// WithRanFuncID sets ran function ID
func WithRanFuncID(ranFuncID int32) func(control *Control) {
	return func(control *Control) {
		control.ranFuncID = ranFuncID
	}
}

// WithRicCallProcessID sets ric call process ID
func WithRicCallProcessID(ricCallPrID types.RicCallProcessID) func(control *Control) {
	return func(control *Control) {
		control.ricCallPrID = ricCallPrID
	}
}

// WithRicInstanceID sets ric instance ID
func WithRicInstanceID(ricInstanceID int32) func(control *Control) {
	return func(control *Control) {
		control.ricInstanceID = ricInstanceID
	}
}

// WithCause sets failure cause
func WithCause(cause *e2apies.Cause) func(control *Control) {
	return func(control *Control) {
		control.cause = cause
	}
}

// WithRicControlStatus sets ric control status
func WithRicControlStatus(ricCtrlStatus e2apies.RiccontrolStatus) func(control *Control) {
	return func(control *Control) {
		control.ricCtrlStatus = ricCtrlStatus
	}
}

// WithRicControlOutcome sets ric control outcome
func WithRicControlOutcome(ricCtrlOut []byte) func(control *Control) {
	return func(control *Control) {
		control.ricCtrlOut = ricCtrlOut
	}
}

// BuildControlAcknowledge builds e2ap control acknowledge message
func (control *Control) BuildControlAcknowledge() (response *e2appducontents.RiccontrolAcknowledge, err error) {

	response = &e2appducontents.RiccontrolAcknowledge{
		ProtocolIes: make([]*e2appducontents.RiccontrolAcknowledgeIes, 0),
	}

	ricRequestID := types.RicRequest{
		RequestorID: types.RicRequestorID(control.reqID),
		InstanceID:  types.RicInstanceID(control.ricInstanceID),
	}
	ranFunctionID := types.RanFunctionID(control.ranFuncID)

	response.SetRicRequestID(ricRequestID).SetRanFunctionID(&ranFunctionID).SetRicCallProcessID(control.ricCallPrID)

	if len(control.ricCtrlOut) != 0 {
		response.SetRicControlOutcome(control.ricCtrlOut)
	}

	return response, nil

}

// BuildControlFailure builds e2ap control failure message
func (control *Control) BuildControlFailure() (response *e2appducontents.RiccontrolFailure, err error) {
	response = &e2appducontents.RiccontrolFailure{
		ProtocolIes: make([]*e2appducontents.RiccontrolFailureIes, 0),
	}
	ricRequestID := types.RicRequest{
		RequestorID: types.RicRequestorID(control.reqID),
		InstanceID:  types.RicInstanceID(control.ricInstanceID),
	}
	ranFuncID := types.RanFunctionID(control.ranFuncID)
	response.SetRicRequestID(ricRequestID).SetRanFunctionID(&ranFuncID).SetRicCallProcessID(control.ricCallPrID).SetCause(control.cause)

	if len(control.ricCtrlOut) != 0 {
		response.SetRicControlOutcome(control.ricCtrlOut)
	}

	return response, nil

}
