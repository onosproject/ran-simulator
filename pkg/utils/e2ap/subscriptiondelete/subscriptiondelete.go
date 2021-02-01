// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscriptiondelete

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2ap-commondatatypes"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
)

// SubscriptionDelete required fields for creating subscription delete response and failure
type SubscriptionDelete struct {
	reqID         int32
	ricInstanceID int32
	ranFuncID     int32
	cause         *e2apies.Cause
	// TODO add more fields including cause of failure
}

// NewSubscriptionDelete creates a new instance of subscription delete
func NewSubscriptionDelete(options ...func(subscriptionDelete *SubscriptionDelete)) (*SubscriptionDelete, error) {
	subscriptionDelete := &SubscriptionDelete{}

	for _, option := range options {
		option(subscriptionDelete)
	}
	return subscriptionDelete, nil
}

// GetRanFuncID returns subscription ran function ID
func (s *SubscriptionDelete) GetRanFuncID() int32 {
	return s.ranFuncID
}

// GetRicInstanceID returns subscription RicInstance ID
func (s *SubscriptionDelete) GetRicInstanceID() int32 {
	return s.ricInstanceID
}

// GetReqID returns subscription request ID
func (s *SubscriptionDelete) GetReqID() int32 {
	return s.reqID
}

// WithRequestID sets request ID
func WithRequestID(reqID int32) func(subscriptionDelete *SubscriptionDelete) {
	return func(subscriptionDelete *SubscriptionDelete) {
		subscriptionDelete.reqID = reqID
	}
}

// WithRanFuncID sets ran function ID
func WithRanFuncID(ranFuncID int32) func(subscriptionDelete *SubscriptionDelete) {
	return func(subscriptionDelete *SubscriptionDelete) {
		subscriptionDelete.ranFuncID = ranFuncID
	}
}

// WithRicInstanceID sets ric instance ID
func WithRicInstanceID(ricInstanceID int32) func(subscriptionDelete *SubscriptionDelete) {
	return func(subscriptionDelete *SubscriptionDelete) {
		subscriptionDelete.ricInstanceID = ricInstanceID
	}
}

// WithCause sets cause of subscription delete failure
func WithCause(cause *e2apies.Cause) func(subscriptionDelete *SubscriptionDelete) {
	return func(subscriptionDelete *SubscriptionDelete) {
		subscriptionDelete.cause = cause
	}
}

// CreateSubscriptionDeleteFailure creates e2 subscription delete failure
func CreateSubscriptionDeleteFailure(subscriptionDelete *SubscriptionDelete) (response *e2appducontents.RicsubscriptionDeleteFailure) {
	ricRequestID := e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes29{
		Id:          int32(v1beta1.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: subscriptionDelete.reqID,         // sequence from e2ap-v01.00.asn1:1126
			RicInstanceId:  subscriptionDelete.ricInstanceID, // sequence from e2ap-v01.00.asn1:1127
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ranFunctionID := e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes5{
		Id:          int32(v1beta1.ProtocolIeIDRanfunctionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RanfunctionId{
			Value: subscriptionDelete.ranFuncID, // range of Integer from e2ap-v01.00.asn1:1050, value from line 1277
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}
	causeOfFailure := e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes1{
		Id:          int32(v1beta1.ProtocolIeIDCause),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
		Value:       subscriptionDelete.cause,
		Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	resp := &e2appducontents.RicsubscriptionDeleteFailure{
		ProtocolIes: &e2appducontents.RicsubscriptionDeleteFailureIes{
			E2ApProtocolIes29: &ricRequestID,  //RIC request ID
			E2ApProtocolIes5:  &ranFunctionID, //RAN function ID
			E2ApProtocolIes1:  &causeOfFailure,
		},
	}

	return resp
}

// CreateSubscriptionDeleteResponse creates e2 subscription delete response
func CreateSubscriptionDeleteResponse(subscriptionDelete *SubscriptionDelete) (response *e2appducontents.RicsubscriptionDeleteResponse) {
	ricRequestID := e2appducontents.RicsubscriptionDeleteResponseIes_RicsubscriptionDeleteResponseIes29{
		Id:          int32(v1beta1.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: subscriptionDelete.reqID,         // sequence from e2ap-v01.00.asn1:1126
			RicInstanceId:  subscriptionDelete.ricInstanceID, // sequence from e2ap-v01.00.asn1:1127
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ranFunctionID := e2appducontents.RicsubscriptionDeleteResponseIes_RicsubscriptionDeleteResponseIes5{
		Id:          int32(v1beta1.ProtocolIeIDRanfunctionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RanfunctionId{
			Value: subscriptionDelete.ranFuncID, // range of Integer from e2ap-v01.00.asn1:1050, value from line 1277
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	resp := &e2appducontents.RicsubscriptionDeleteResponse{
		ProtocolIes: &e2appducontents.RicsubscriptionDeleteResponseIes{
			E2ApProtocolIes29: &ricRequestID,  //RIC request ID
			E2ApProtocolIes5:  &ranFunctionID, //RAN function ID
		},
	}

	return resp

}
