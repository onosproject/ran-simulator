// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package kpm2

import (
	e2smkpmv2sm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/servicemodel"
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"google.golang.org/protobuf/proto"
)

func (sm *Client) getActionDefinition(actionList []*e2appducontents.RicactionToBeSetupItemIes, ricActionsAccepted []*e2aptypes.RicActionID) ([]*e2smkpmv2.E2SmKpmActionDefinition, error) {
	var actionDefinitions []*e2smkpmv2.E2SmKpmActionDefinition
	for _, action := range actionList {
		for _, acceptedActionID := range ricActionsAccepted {
			if action.GetValue().GetRicactionToBeSetupItem().GetRicActionId().GetValue() == int32(*acceptedActionID) {
				actionDefinitionBytes := action.GetValue().GetRicactionToBeSetupItem().GetRicActionDefinition().GetValue()
				var kpm2ServiceModel e2smkpmv2sm.Kpm2ServiceModel

				actionDefinitionProtoBytes, err := kpm2ServiceModel.ActionDefinitionASN1toProto(actionDefinitionBytes)
				if err != nil {
					log.Warn(err)
					return nil, err
				}

				actionDefinition := &e2smkpmv2.E2SmKpmActionDefinition{}
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

// getReportPeriod extracts report period
func (sm *Client) getReportPeriod(request *e2appducontents.RicsubscriptionRequest) (int64, error) {
	var eventTriggerAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			eventTriggerAsnBytes = v.GetValue().GetRicsubscriptionDetails().GetRicEventTriggerDefinition().GetValue()
			break
		}
	}
	var kpm2ServiceModel e2smkpmv2sm.Kpm2ServiceModel
	eventTriggerProtoBytes, err := kpm2ServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return 0, err
	}
	eventTriggerDefinition := &e2smkpmv2.E2SmKpmEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return 0, err
	}
	reportPeriod := eventTriggerDefinition.GetEventDefinitionFormats().GetEventDefinitionFormat1().GetReportingPeriod()
	return reportPeriod, nil
}
