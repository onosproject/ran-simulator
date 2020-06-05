// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2

import (
	"context"
	"fmt"

	"github.com/onosproject/onos-ric/api/sb/e2ap/elements"
	"github.com/onosproject/onos-ric/api/sb/e2ap/subscription"
)

// RicSubscribeDelete implements E2 subscription delete procedure
func (s *Server) RicSubscribeDelete(ctx context.Context, req *subscription.RicSubscriptionDeleteRequest) (*subscription.RicSubscriptionDeleteResponse, error) {

	return nil, fmt.Errorf("not yet implemented")
}

// RicSubscribe - implements E2 subscription procedure
func (s *Server) RicSubscribe(ctx context.Context, req *subscription.RicSubscriptionRequest) (*subscription.RicSubscriptionResponse, error) {

	// Determine the target function using the
	// information in the RAN Function ID IE
	// and configure the requested event trigger
	// using information in the RIC Subscription Details IE.

	// If one or more Report, Insert and/or Policy RIC service actions are included in the
	// RIC Subscription Details IE then the target function shall validate
	// the event trigger and requested action sequence and,
	// if accepted, store the required RIC Request ID,
	// RIC Event Trigger Definition IE and sequence of RIC Action ID IE,
	// RIC Action Type IE, RIC Action Definition, RIC Subsequent Action IE.

	actionAdmittedList := &elements.ActionAdmittedList{}
	actionNotAdmittedList := &elements.ActionNotAdmittedList{}
	actions := req.Details.Actions

	eventTrigger := req.Details.EventTriggerDefinition

	if eventTrigger == nil {
		return &subscription.RicSubscriptionResponse{
			Response: &subscription.RicSubscriptionResponse_SubscriptionFailure{
				SubscriptionFailure: &subscription.SubscriptionFailure{
					MessageType:           elements.MessageType_UNSUCCESSFUL_OUTCOME,
					RanFunctionId:         req.RanFunctionId,
					RicRequestId:          req.RicRequestId,
					ActionNotAdmittedList: actionNotAdmittedList,
				},
			},
		}, fmt.Errorf("event trigger definition cannot be null")
	}

	// Validate number of requested actions
	if len(actions) > MaxNumActions {
		return &subscription.RicSubscriptionResponse{
			Response: &subscription.RicSubscriptionResponse_SubscriptionFailure{
				SubscriptionFailure: &subscription.SubscriptionFailure{
					MessageType:           elements.MessageType_UNSUCCESSFUL_OUTCOME,
					RanFunctionId:         req.RanFunctionId,
					RicRequestId:          req.RicRequestId,
					ActionNotAdmittedList: actionNotAdmittedList,
				},
			},
		}, fmt.Errorf("number of actions must be less than %d", MaxNumActions)
	}

	// List of admitted actions
	var actionAdmittedIDs []*elements.RicActionID
	// List of not admitted actions
	var actionNotAdmittedIDs []*elements.RicActionID

	// For now we include REPORT action in the list of admitted actions
	// and add others in list of not admitted actions
	// TODO support INSERT and POLICY in the list of admitted actions
	for _, action := range actions {
		switch action.RicActionType {
		case elements.RicActionType_REPORT:
			actionAdmittedIDs = append(actionAdmittedIDs, action.RicActionId)
			// TODO send indication report using indication channel

		case elements.RicActionType_INSERT:
			actionNotAdmittedIDs = append(actionNotAdmittedIDs, action.RicActionId)
		case elements.RicActionType_POLICY:
			actionNotAdmittedIDs = append(actionNotAdmittedIDs, action.RicActionId)

		}
	}

	actionAdmittedList = &elements.ActionAdmittedList{
		ActionAdmittedIds: actionAdmittedIDs,
	}

	response := &subscription.RicSubscriptionResponse{
		Response: &subscription.RicSubscriptionResponse_SubscriptionResponse{
			SubscriptionResponse: &subscription.SubscriptionResponse{
				MessageType:           elements.MessageType_SUCCESSFUL_OUTCOME,
				RanFunctionId:         req.RanFunctionId,
				RicRequestId:          req.RicRequestId,
				ActionAdmittedList:    actionAdmittedList,
				ActionNotAdmittedList: actionNotAdmittedList,
			},
		},
	}

	return response, nil
}
