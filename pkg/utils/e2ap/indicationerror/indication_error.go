// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package indicationerror

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
)

// ErrorIndication required fields for creating error indication error message
type ErrorIndication struct {
	reqID           int32
	ricInstanceID   int32
	ranFuncID       int32
	cause           *e2apies.Cause
	failureProcCode int32
	failureTrigMsg  *e2ap_commondatatypes.TriggeringMessage
	critDiags       []*types.CritDiag
	failureCrit     *e2ap_commondatatypes.Criticality
}

// NewErrorIndication creates a new error indication
func NewErrorIndication(options ...func(errorIndication *ErrorIndication)) *ErrorIndication {
	errorIndication := &ErrorIndication{}

	for _, option := range options {
		option(errorIndication)
	}

	return errorIndication
}

// WithRequestID sets request ID
func WithRequestID(reqID int32) func(*ErrorIndication) {
	return func(errorIndication *ErrorIndication) {
		errorIndication.reqID = reqID
	}
}

// WithRanFuncID sets ran function ID
func WithRanFuncID(ranFuncID int32) func(*ErrorIndication) {
	return func(errorIndication *ErrorIndication) {
		errorIndication.ranFuncID = ranFuncID
	}
}

// WithRicInstanceID sets ric instance ID
func WithRicInstanceID(ricInstanceID int32) func(*ErrorIndication) {
	return func(errorIndication *ErrorIndication) {
		errorIndication.ricInstanceID = ricInstanceID
	}
}

// WithFailureProcCode sets failure proc code
func WithFailureProcCode(failureProcCode int32) func(*ErrorIndication) {
	return func(errorIndication *ErrorIndication) {
		errorIndication.failureProcCode = failureProcCode
	}
}

// WithCause sets cause of error
func WithCause(cause *e2apies.Cause) func(*ErrorIndication) {
	return func(errorIndication *ErrorIndication) {
		errorIndication.cause = cause
	}
}

// Build builds an error indication message
func (e *ErrorIndication) Build() (*e2appducontents.ErrorIndication, error) {
	rrID := types.RicRequest{
		RequestorID: types.RicRequestorID(e.reqID),
		InstanceID:  types.RicInstanceID(e.ricInstanceID),
	}
	rfID := types.RanFunctionID(e.ranFuncID)
	errorIndication := &e2appducontents.ErrorIndication{
		ProtocolIes: make([]*e2appducontents.ErrorIndicationIes, 0),
	}
	failureProcCode := v2.ProcedureCodeT(e.failureProcCode)
	errorIndication.SetRicRequestID(&rrID).SetRanFunctionID(&rfID).SetCause(e.cause).
		SetCriticalityDiagnostics(&failureProcCode, e.failureCrit, e.failureTrigMsg, &rrID, e.critDiags)

	return errorIndication, nil
}
