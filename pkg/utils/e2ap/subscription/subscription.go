// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
)

// Subscription defines required fields for creating subscription response and failures
type Subscription struct {
	reqID                 int32
	ricInstanceID         int32
	ranFuncID             int32
	ricActionsAccepted    []*types.RicActionID
	ricActionsNotAdmitted map[types.RicActionID]*e2apies.Cause
	cause                 *e2apies.Cause
}

// NewSubscription creates a new instance of subscription
func NewSubscription(options ...func(*Subscription)) *Subscription {
	subscription := &Subscription{}

	for _, option := range options {
		option(subscription)
	}

	return subscription

}

// GetRanFuncID returns subscription ran function ID
func (subscription *Subscription) GetRanFuncID() int32 {
	return subscription.ranFuncID
}

// GetRicInstanceID returns subscription RicInstance ID
func (subscription *Subscription) GetRicInstanceID() int32 {
	return subscription.ricInstanceID
}

// GetReqID returns subscription request ID
func (subscription *Subscription) GetReqID() int32 {
	return subscription.reqID
}

// WithRequestID sets request ID
func WithRequestID(reqID int32) func(*Subscription) {
	return func(subscription *Subscription) {
		subscription.reqID = reqID
	}
}

// WithRanFuncID sets ran function ID
func WithRanFuncID(ranFuncID int32) func(*Subscription) {
	return func(subscription *Subscription) {
		subscription.ranFuncID = ranFuncID
	}
}

// WithRicInstanceID sets ric instance ID
func WithRicInstanceID(ricInstanceID int32) func(*Subscription) {
	return func(subscription *Subscription) {
		subscription.ricInstanceID = ricInstanceID
	}
}

// WithActionsAccepted sets accepted actions
func WithActionsAccepted(ricActionsAccepted []*types.RicActionID) func(*Subscription) {
	return func(subscription *Subscription) {
		subscription.ricActionsAccepted = ricActionsAccepted
	}
}

// WithActionsNotAdmitted sets not admitted actions
func WithActionsNotAdmitted(ricActionsNotAdmitted map[types.RicActionID]*e2apies.Cause) func(*Subscription) {
	return func(subscription *Subscription) {
		subscription.ricActionsNotAdmitted = ricActionsNotAdmitted
	}
}

// WithCause sets subscription failure cause
func WithCause(cause *e2apies.Cause) func(subscription *Subscription) {
	return func(subscription *Subscription) {
		subscription.cause = cause
	}
}

// BuildSubscriptionFailure builds e2ap subscription failure
func (subscription *Subscription) BuildSubscriptionFailure() (response *e2appducontents.RicsubscriptionFailure, err error) {

	rfID := types.RanFunctionID(subscription.ranFuncID)
	rrID := types.RicRequest{
		RequestorID: types.RicRequestorID(subscription.reqID),
		InstanceID:  types.RicInstanceID(subscription.ricInstanceID),
	}

	resp := &e2appducontents.RicsubscriptionFailure{
		ProtocolIes: make([]*e2appducontents.RicsubscriptionFailureIes, 0),
	}
	resp.SetRanFunctionID(&rfID).SetRicRequestID(&rrID).SetCause(subscription.cause)

	return resp, nil
}

// BuildSubscriptionResponse builds e2ap subscription response
func (subscription *Subscription) BuildSubscriptionResponse() (response *e2appducontents.RicsubscriptionResponse, err error) {

	rfID := types.RanFunctionID(subscription.ranFuncID)
	rrID := types.RicRequest{
		RequestorID: types.RicRequestorID(subscription.reqID),
		InstanceID:  types.RicInstanceID(subscription.ricInstanceID),
	}

	resp := &e2appducontents.RicsubscriptionResponse{
		ProtocolIes: make([]*e2appducontents.RicsubscriptionResponseIes, 0),
	}
	resp.SetRicRequestID(&rrID).SetRanFunctionID(&rfID).SetRicActionAdmitted(subscription.ricActionsAccepted).
		SetRicActionNotAdmitted(subscription.ricActionsNotAdmitted)

	return resp, nil
}
