// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
)

// Subscription defines required fields for creating subscription response and failures
type Subscription struct {
	reqID                 int32
	ricInstanceID         int32
	ranFuncID             int32
	ricActionsAccepted    []*types.RicActionID
	ricActionsNotAdmitted map[types.RicActionID]*e2apies.Cause
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

// BuildSubscriptionFailure builds e2ap subscription failure
func (subscription *Subscription) BuildSubscriptionFailure() (response *e2appducontents.RicsubscriptionFailure, err error) {
	ricRequestID := e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes29{
		Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: subscription.reqID,
			RicInstanceId:  subscription.ricInstanceID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ranFunctionID := e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes5{
		Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RanfunctionId{
			Value: subscription.ranFuncID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ricActionNotAdmittedList := e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes18{
		Id:          int32(v1beta2.ProtocolIeIDRicactionsNotAdmitted),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2appducontents.RicactionNotAdmittedList{
			Value: make([]*e2appducontents.RicactionNotAdmittedItemIes, 0),
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	for ricActionID, cause := range subscription.ricActionsNotAdmitted {
		ranaItemIe := e2appducontents.RicactionNotAdmittedItemIes{
			Id:          int32(v1beta2.ProtocolIeIDRicactionNotAdmittedItem),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Value: &e2appducontents.RicactionNotAdmittedItem{
				RicActionId: &e2apies.RicactionId{
					Value: int32(ricActionID),
				},
				Cause: cause,
			},
			Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
		ricActionNotAdmittedList.GetValue().Value = append(ricActionNotAdmittedList.GetValue().Value, &ranaItemIe)
	}

	resp := &e2appducontents.RicsubscriptionFailure{
		ProtocolIes: &e2appducontents.RicsubscriptionFailureIes{
			E2ApProtocolIes5:  &ranFunctionID,
			E2ApProtocolIes18: &ricActionNotAdmittedList,
			E2ApProtocolIes29: &ricRequestID,
		},
	}

	return resp, nil
}

// BuildSubscriptionResponse builds e2ap subscription response
func (subscription *Subscription) BuildSubscriptionResponse() (response *e2appducontents.RicsubscriptionResponse, err error) {
	ricRequestID := e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes29{
		Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RicrequestId{
			RicRequestorId: subscription.reqID,
			RicInstanceId:  subscription.ricInstanceID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ranFunctionID := e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes5{
		Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.RanfunctionId{
			Value: subscription.ranFuncID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	ricActionAdmit := e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes17{
		Id:          int32(v1beta2.ProtocolIeIDRicactionsAdmitted),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2appducontents.RicactionAdmittedList{
			Value: make([]*e2appducontents.RicactionAdmittedItemIes, 0),
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	for _, raa := range subscription.ricActionsAccepted {
		raaIe := &e2appducontents.RicactionAdmittedItemIes{
			Id:          int32(v1beta2.ProtocolIeIDRicactionAdmittedItem),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Value: &e2appducontents.RicactionAdmittedItem{
				RicActionId: &e2apies.RicactionId{
					Value: int32(*raa),
				},
			},
			Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
		ricActionAdmit.GetValue().Value = append(ricActionAdmit.GetValue().Value, raaIe)
	}

	resp := &e2appducontents.RicsubscriptionResponse{
		ProtocolIes: &e2appducontents.RicsubscriptionResponseIes{
			E2ApProtocolIes29: &ricRequestID,  //RIC request ID
			E2ApProtocolIes5:  &ranFunctionID, //RAN function ID
			E2ApProtocolIes17: &ricActionAdmit,
		},
	}

	return resp, nil
}
