// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package subscriptiondelete

import (
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
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

	rrID := types.RicRequest{
		RequestorID: types.RicRequestorID(subscriptionDelete.reqID),
		InstanceID:  types.RicInstanceID(subscriptionDelete.ricInstanceID),
	}

	resp := &e2appducontents.RicsubscriptionDeleteFailure{
		ProtocolIes: make([]*e2appducontents.RicsubscriptionDeleteFailureIes, 0),
	}
	resp.SetRicRequestID(&rrID).SetRanFunctionID(types.RanFunctionID(subscriptionDelete.ranFuncID)).SetCause(subscriptionDelete.cause)

	return resp, nil
}

// BuildSubscriptionDeleteResponse builds subscription delete response
func (subscriptionDelete *SubscriptionDelete) BuildSubscriptionDeleteResponse() (response *e2appducontents.RicsubscriptionDeleteResponse, err error) {

	rrID := types.RicRequest{
		RequestorID: types.RicRequestorID(subscriptionDelete.reqID),
		InstanceID:  types.RicInstanceID(subscriptionDelete.ricInstanceID),
	}

	resp := &e2appducontents.RicsubscriptionDeleteResponse{
		ProtocolIes: make([]*e2appducontents.RicsubscriptionDeleteResponseIes, 0),
	}
	resp.SetRicRequestID(&rrID).SetRanFunctionID(types.RanFunctionID(subscriptionDelete.ranFuncID))

	return resp, nil

}
