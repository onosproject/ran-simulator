// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm

import (
	"context"

	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
)

var _ servicemodel.ServiceModel = &ServiceModel{}

// ServiceModel kpm service model struct
type ServiceModel struct {
}

// RICControl ....
func (sm *ServiceModel) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	panic("implement me")

}

// RICSubscription ...
func (sm *ServiceModel) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {

	var ricActionsAccepted []*types.RicActionID
	var ricActionsNotAdmitted map[types.RicActionID]*e2apies.Cause
	actionList := request.ProtocolIes.E2ApProtocolIes30.Value.RicActionToBeSetupList.Value

	reqID := request.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId
	ranFuncID := request.ProtocolIes.E2ApProtocolIes5.Value.Value
	ricInstanceID := request.ProtocolIes.E2ApProtocolIes29.Value.RicInstanceId

	for _, action := range actionList {
		actionID := types.RicActionID(action.Value.RicActionId.Value)
		actionType := action.Value.RicActionType
		if actionType == e2apies.RicactionType_RICACTION_TYPE_REPORT {
			ricActionsAccepted = append(ricActionsAccepted, &actionID)
		}
		// TODO handle not admitted actions
	}
	subscription, _ := NewSubscription(
		WithRequestID(reqID),
		WithRanFuncID(ranFuncID),
		WithRicInstanceID(ricInstanceID),
		WithActionsAccepted(ricActionsAccepted),
		WithActionsNotAdmitted(ricActionsNotAdmitted))

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		subscriptionFailure := createSubscriptionFailure(subscription)
		return nil, subscriptionFailure, nil
	}

	subscriptionResponse := createSubscriptionResponse(subscription)
	return subscriptionResponse, nil, nil

}

// RICSubscriptionDelete ...
func (sm *ServiceModel) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	panic("implement me")
}
