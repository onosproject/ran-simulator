// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscriptiondelete

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
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
func NewSubscriptionDelete(options ...func(subscriptionDelete *SubscriptionDelete)) *SubscriptionDelete {
	subscriptionDelete := &SubscriptionDelete{}

	for _, option := range options {
		option(subscriptionDelete)
	}
	return subscriptionDelete
}

// GetRanFuncID returns subscription ran function ID
func (subscriptionDelete *SubscriptionDelete) GetRanFuncID() int32 {
	return subscriptionDelete.ranFuncID
}

// GetRicInstanceID returns subscription RicInstance ID
func (subscriptionDelete *SubscriptionDelete) GetRicInstanceID() int32 {
	return subscriptionDelete.ricInstanceID
}

// GetReqID returns subscription request ID
func (subscriptionDelete *SubscriptionDelete) GetReqID() int32 {
	return subscriptionDelete.reqID
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

// BuildSubscriptionDeleteFailure builds subscription delete failure
func (subscriptionDelete *SubscriptionDelete) BuildSubscriptionDeleteFailure() (response *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	ricRequestID := e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes29{
		Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: subscriptionDelete.reqID,         // sequence from e2ap-v01.00.asn1:1126
			RicInstanceId:  subscriptionDelete.ricInstanceID, // sequence from e2ap-v01.00.asn1:1127
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ranFunctionID := e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes5{
		Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RanfunctionId{
			Value: subscriptionDelete.ranFuncID, // range of Integer from e2ap-v01.00.asn1:1050, value from line 1277
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}
	causeOfFailure := e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes1{
		Id:          int32(v1beta2.ProtocolIeIDCause),
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

	return resp, nil
}

// BuildSubscriptionDeleteResponse builds subscription delete response
func (subscriptionDelete *SubscriptionDelete) BuildSubscriptionDeleteResponse() (response *e2appducontents.RicsubscriptionDeleteResponse, err error) {
	ricRequestID := e2appducontents.RicsubscriptionDeleteResponseIes_RicsubscriptionDeleteResponseIes29{
		Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: subscriptionDelete.reqID,         // sequence from e2ap-v01.00.asn1:1126
			RicInstanceId:  subscriptionDelete.ricInstanceID, // sequence from e2ap-v01.00.asn1:1127
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ranFunctionID := e2appducontents.RicsubscriptionDeleteResponseIes_RicsubscriptionDeleteResponseIes5{
		Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
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

	return resp, nil

}
