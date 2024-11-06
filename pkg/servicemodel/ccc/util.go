// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package ccc

import (
	e2smcccsm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/servicemodel"
	e2smccc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/v1/e2sm-ccc-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"google.golang.org/protobuf/proto"
)

// getReportInterval extracts event trigger type
func (sm *Client) getReportInterval(request *e2appducontents.RicsubscriptionRequest) (int32, error) {
	var eventTriggerAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			eventTriggerAsnBytes = v.GetValue().GetRicsubscriptionDetails().GetRicEventTriggerDefinition().GetValue()
			break
		}
	}

	var cccServiceModel e2smcccsm.CCCServiceModel
	eventTriggerProtoBytes, err := cccServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return 0, err
	}
	eventTriggerDefinition := &e2smccc.E2SmCCcRIceventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return -1, err
	}
	eventTriggerType := eventTriggerDefinition.GetEventTriggerDefinitionFormat().GetE2SmCccEventTriggerDefinitionFormat3().GetPeriod()
	return eventTriggerType, nil
}

func (sm *Client) getActionDefinition(actionList []*e2appducontents.RicactionToBeSetupItemIes, ricActionsAccepted []*e2aptypes.RicActionID) ([]*e2smccc.E2SmCCcRIcactionDefinition, error) {
	var actionDefinitions []*e2smccc.E2SmCCcRIcactionDefinition
	for _, action := range actionList {
		for _, acceptedActionID := range ricActionsAccepted {
			if action.GetValue().GetRicactionToBeSetupItem().GetRicActionId().GetValue() == int32(*acceptedActionID) {
				actionDefinitionBytes := action.GetValue().GetRicactionToBeSetupItem().GetRicActionDefinition().GetValue()
				var cccServiceModel e2smcccsm.CCCServiceModel

				actionDefinitionProtoBytes, err := cccServiceModel.ActionDefinitionASN1toProto(actionDefinitionBytes)
				if err != nil {
					log.Warn(err)
					return nil, err
				}

				actionDefinition := &e2smccc.E2SmCCcRIcactionDefinition{}
				err = proto.Unmarshal(actionDefinitionProtoBytes, actionDefinition)
				if err != nil {
					log.Warn(err)
					return nil, err
				}

				actionDefinitions = append(actionDefinitions, actionDefinition)

			}
		}
	}
	return actionDefinitions, nil
}

func (sm *Client) getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2smccc.E2SmCCcRIcControlMessage, error) {
	var cccServiceModel e2smcccsm.CCCServiceModel
	var controlMessageAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolMessage) {
			controlMessageAsnBytes = v.GetValue().GetRiccontrolMessage().GetValue()
			break
		}
	}
	controlMessageProtoBytes, err := cccServiceModel.ControlMessageASN1toProto(controlMessageAsnBytes)
	if err != nil {
		return nil, err
	}
	controlMessage := &e2smccc.E2SmCCcRIcControlMessage{}
	err = proto.Unmarshal(controlMessageProtoBytes, controlMessage)

	if err != nil {
		return nil, err
	}
	return controlMessage, nil
}

func (sm *Client) getControlHeader(request *e2appducontents.RiccontrolRequest) (*e2smccc.E2SmCCcRIcControlHeader, error) {
	var cccServiceModel e2smcccsm.CCCServiceModel
	var controlHeaderAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolHeader) {
			controlHeaderAsnBytes = v.GetValue().GetRiccontrolHeader().GetValue()
			break
		}
	}
	controlHeaderProtoBytes, err := cccServiceModel.ControlHeaderASN1toProto(controlHeaderAsnBytes)
	if err != nil {
		return nil, err
	}
	controlHeader := &e2smccc.E2SmCCcRIcControlHeader{}
	err = proto.Unmarshal(controlHeaderProtoBytes, controlHeader)
	if err != nil {
		return nil, err
	}

	return controlHeader, nil
}
