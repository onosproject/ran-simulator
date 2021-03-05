// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package control

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
)

// Control defines required fields for creating control acknowledge and failure responses
type Control struct {
	reqID         int32
	ricInstanceID int32
	ranFuncID     int32
	ricCallPrID   types.RicCallProcessID
	ricCtrlStatus e2apies.RiccontrolStatus
	ricCtrlOut    []byte
	cause         e2apies.Cause
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
func WithCause(cause e2apies.Cause) func(control *Control) {
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

	ricRequestID := e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes29{
		Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: control.reqID,
			RicInstanceId:  control.ricInstanceID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ranFunctionID := e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes5{
		Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RanfunctionId{
			Value: control.ranFuncID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ricCallProcessID := e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes20{
		Id:          int32(v1beta2.ProtocolIeIDRiccallProcessID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2ap_commondatatypes.RiccallProcessId{
			Value: []byte(control.ricCallPrID),
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
	}

	ricControlStatus := e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes24{
		Id:          int32(v1beta2.ProtocolIeIDRiccontrolStatus),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value:       control.ricCtrlStatus,
		Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ricControlOutcome := e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes32{
		Id:          int32(v1beta2.ProtocolIeIDRiccontrolOutcome),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2ap_commondatatypes.RiccontrolOutcome{
			Value: control.ricCtrlOut,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
	}

	response = &e2appducontents.RiccontrolAcknowledge{
		ProtocolIes: &e2appducontents.RiccontrolAcknowledgeIes{
			E2ApProtocolIes29: &ricRequestID,      // RIC Requestor & RIC Instance ID
			E2ApProtocolIes5:  &ranFunctionID,     // RAN function ID
			E2ApProtocolIes20: &ricCallProcessID,  // RIC Call Process ID
			E2ApProtocolIes24: &ricControlStatus,  // RIC Control Status
			E2ApProtocolIes32: &ricControlOutcome, // RIC Control Outcome
		},
	}

	return response, nil

}

// BuildControlFailure builds e2ap control failure message
func (control *Control) BuildControlFailure() (response *e2appducontents.RiccontrolFailure, err error) {
	ricRequestID := e2appducontents.RiccontrolFailureIes_RiccontrolFailureIes29{
		Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: control.reqID,
			RicInstanceId:  control.ricInstanceID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ranFunctionID := e2appducontents.RiccontrolFailureIes_RiccontrolFailureIes5{
		Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RanfunctionId{
			Value: control.ranFuncID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ricCallProcessID := e2appducontents.RiccontrolFailureIes_RiccontrolFailureIes20{
		Id:          int32(v1beta2.ProtocolIeIDRiccallProcessID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2ap_commondatatypes.RiccallProcessId{
			Value: []byte(control.ricCallPrID),
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
	}

	ricCause := e2appducontents.RiccontrolFailureIes_RiccontrolFailureIes1{
		Id:          int32(v1beta2.ProtocolIeIDCause),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
		Value:       &control.cause,
		Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ricControlOutcome := e2appducontents.RiccontrolFailureIes_RiccontrolFailureIes32{
		Id:          int32(v1beta2.ProtocolIeIDRiccontrolOutcome),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2ap_commondatatypes.RiccontrolOutcome{
			Value: control.ricCtrlOut,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
	}

	response = &e2appducontents.RiccontrolFailure{
		ProtocolIes: &e2appducontents.RiccontrolFailureIes{
			E2ApProtocolIes29: &ricRequestID,      // RIC Requestor & RIC Instance ID
			E2ApProtocolIes5:  &ranFunctionID,     // RAN function ID
			E2ApProtocolIes20: &ricCallProcessID,  // RIC Call Process ID
			E2ApProtocolIes1:  &ricCause,          // Cause
			E2ApProtocolIes32: &ricControlOutcome, // RIC Control Outcome
		},
	}

	return response, nil

}
