// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package indicationerror

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
)

// ErrorIndication required fields for creating error indication error message
type ErrorIndication struct {
	reqID           int32
	ricInstanceID   int32
	ranFuncID       int32
	cause           e2apies.Cause
	failureProcCode int32
	failureTrigMsg  e2ap_commondatatypes.TriggeringMessage
	critDiags       []*types.CritDiag
	failureCrit     e2ap_commondatatypes.Criticality
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
func WithCause(cause e2apies.Cause) func(*ErrorIndication) {
	return func(errorIndication *ErrorIndication) {
		errorIndication.cause = cause
	}
}

// Build builds an error indication message
func (e *ErrorIndication) Build() (*e2appducontents.ErrorIndication, error) {
	ricRequestID := e2appducontents.ErrorIndicationIes_ErrorIndicationIes29{
		Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: e.reqID,
			RicInstanceId:  e.ricInstanceID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
	}

	ranFunctionID := e2appducontents.ErrorIndicationIes_ErrorIndicationIes5{
		Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RanfunctionId{
			Value: e.ranFuncID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
	}

	errorCause := e2appducontents.ErrorIndicationIes_ErrorIndicationIes1{
		Id:          int32(v1beta2.ProtocolIeIDCause),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
		Value:       &e.cause,
		Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
	}

	criticalityDiagnostics := e2appducontents.ErrorIndicationIes_ErrorIndicationIes2{
		Id:          int32(v1beta2.ProtocolIeIDCriticalityDiagnostics),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
		Value: &e2apies.CriticalityDiagnostics{
			ProcedureCode: &e2ap_commondatatypes.ProcedureCode{
				Value: e.failureProcCode,
			},
			TriggeringMessage:    e.failureTrigMsg,
			ProcedureCriticality: e.failureCrit,
			RicRequestorId: &e2apies.RicrequestId{
				RicRequestorId: e.reqID,
				RicInstanceId:  e.ricInstanceID,
			},
			IEsCriticalityDiagnostics: &e2apies.CriticalityDiagnosticsIeList{
				Value: make([]*e2apies.CriticalityDiagnosticsIeItem, 0),
			},
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
	}

	for _, critDiag := range e.critDiags {
		criticDiagnostics := e2apies.CriticalityDiagnosticsIeItem{
			IEcriticality: critDiag.IECriticality,
			IEId: &e2ap_commondatatypes.ProtocolIeId{
				Value: int32(critDiag.IEId), // value were taken from e2ap-v01.00.asn1:1278
			},
			TypeOfError: critDiag.TypeOfError,
		}
		criticalityDiagnostics.Value.IEsCriticalityDiagnostics.Value = append(criticalityDiagnostics.Value.IEsCriticalityDiagnostics.Value, &criticDiagnostics)
	}

	errorIndication := &e2appducontents.ErrorIndication{
		ProtocolIes: &e2appducontents.ErrorIndicationIes{
			E2ApProtocolIes29: &ricRequestID,           // RIC Requestor & RIC Instance ID
			E2ApProtocolIes5:  &ranFunctionID,          // RAN function ID
			E2ApProtocolIes1:  &errorCause,             // Cause
			E2ApProtocolIes2:  &criticalityDiagnostics, // CriticalityDiagnostics
		},
	}

	return errorIndication, nil
}
