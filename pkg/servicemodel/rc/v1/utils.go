// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	e2smrc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/servicemodel"
	e2smrcies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-rc-ies"
	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"google.golang.org/protobuf/proto"
)

func getActionDefinitionMap(actionList []*e2appducontents.RicactionToBeSetupItemIes, ricActionsAccepted []*e2aptypes.RicActionID) (map[*e2aptypes.RicActionID]*e2smrcies.E2SmRcActionDefinition, error) {
	actionDefinitionsMap := make(map[*e2aptypes.RicActionID]*e2smrcies.E2SmRcActionDefinition)
	for _, action := range actionList {
		for _, actionID := range ricActionsAccepted {
			if action.GetValue().GetRicactionToBeSetupItem().GetRicActionId().GetValue() == int32(*actionID) {
				actionDefinitionBytes := action.GetValue().GetRicactionToBeSetupItem().GetRicActionDefinition().GetValue()
				var rcServiceModel e2smrc.RCServiceModel

				actionDefinitionProtoBytes, err := rcServiceModel.ActionDefinitionASN1toProto(actionDefinitionBytes)
				if err != nil {
					return nil, err
				}

				actionDefinition := &e2smrcies.E2SmRcActionDefinition{}
				err = proto.Unmarshal(actionDefinitionProtoBytes, actionDefinition)
				if err != nil {
					return nil, err
				}
				actionDefinitionsMap[actionID] = actionDefinition
			}
		}
	}
	return actionDefinitionsMap, nil
}

func getEventTrigger(request *e2appducontents.RicsubscriptionRequest) (*e2smrcies.E2SmRcEventTrigger, error) {
	var eventTriggerAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			eventTriggerAsnBytes = v.GetValue().GetRicsubscriptionDetails().GetRicEventTriggerDefinition().GetValue()
			break
		}
	}

	var rcServiceModel e2smrc.RCServiceModel
	eventTriggerProtoBytes, err := rcServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return nil, err
	}
	eventTriggerDefinition := &e2smrcies.E2SmRcEventTrigger{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return nil, err
	}

	return eventTriggerDefinition, nil
}

func getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2smrcies.E2SmRcControlMessage, error) {
	var rcServiceModel e2smrc.RCServiceModel
	var controlMessageAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolMessage) {
			controlMessageAsnBytes = v.GetValue().GetRiccontrolMessage().GetValue()
			break
		}
	}
	controlMessageProtoBytes, err := rcServiceModel.ControlMessageASN1toProto(controlMessageAsnBytes)
	if err != nil {
		return nil, err
	}
	controlMessage := &e2smrcies.E2SmRcControlMessage{}
	err = proto.Unmarshal(controlMessageProtoBytes, controlMessage)

	if err != nil {
		return nil, err
	}
	return controlMessage, nil
}

func getControlHeader(request *e2appducontents.RiccontrolRequest) (*e2smrcies.E2SmRcControlHeader, error) {
	var rcServiceModel e2smrc.RCServiceModel
	var controlHeaderAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolHeader) {
			controlHeaderAsnBytes = v.GetValue().GetRiccontrolHeader().GetValue()
			break
		}
	}
	controlHeaderProtoBytes, err := rcServiceModel.ControlHeaderASN1toProto(controlHeaderAsnBytes)
	if err != nil {
		return nil, err
	}
	controlHeader := &e2smrcies.E2SmRcControlHeader{}
	err = proto.Unmarshal(controlHeaderProtoBytes, controlHeader)
	if err != nil {
		return nil, err
	}

	return controlHeader, nil
}
