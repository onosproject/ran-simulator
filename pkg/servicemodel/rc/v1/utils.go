// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/pdubuilder"
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

func createRANParametersInsertStyle3List() ([]*e2smrcies.InsertIndicationRanparameterItem, error) {
	// RAN Parameters for Insert Style 3
	insertRANParametersStyle3List := make([]*e2smrcies.InsertIndicationRanparameterItem, 0)
	ranParameter1, err := pdubuilder.CreateInsertIndicationRanparameterItem(1, "Target Primary Cell ID")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter1)

	ranParameter2, err := pdubuilder.CreateInsertIndicationRanparameterItem(2, "Target Cell")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter2)

	ranParameter3, err := pdubuilder.CreateInsertIndicationRanparameterItem(3, "NR Cell")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter3)

	ranParameter4, err := pdubuilder.CreateInsertIndicationRanparameterItem(4, "NR CGI")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter4)

	ranParameter5, err := pdubuilder.CreateInsertIndicationRanparameterItem(7, "List of PDU sessions for handover")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter5)

	ranParameter6, err := pdubuilder.CreateInsertIndicationRanparameterItem(13, "List of DRBs for handover")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter6)
	return insertRANParametersStyle3List, nil

}

func createRANParametersReportStyle3List() ([]*e2smrcies.ReportRanparameterItem, error) {
	// RAN Parameters for Report Style 3
	reportParametersStyle3List := make([]*e2smrcies.ReportRanparameterItem, 0)

	return reportParametersStyle3List, nil
}

func createRANParametersReportStyle2List() ([]*e2smrcies.ReportRanparameterItem, error) {
	// RAN Parameters for Report Style 2
	reportParametersStyle2List := make([]*e2smrcies.ReportRanparameterItem, 0)
	ranParameter1, err := pdubuilder.CreateReportRanparameterItem(1, "Current UE ID")
	if err != nil {
		return nil, err
	}

	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter1)

	ranParameter2, err := pdubuilder.CreateReportRanparameterItem(21001, "S-NSSAI")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter2)

	ranParameter3, err := pdubuilder.CreateReportRanparameterItem(21002, "SST")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter3)

	ranParameter4, err := pdubuilder.CreateReportRanparameterItem(21003, "SD")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter4)

	ranParameter5, err := pdubuilder.CreateReportRanparameterItem(27108, "Best Neighboring Cell")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter5)

	ranParameter6, err := pdubuilder.CreateReportRanparameterItem(21528, "List of Neighbor cells")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter6)

	ranParameter7, err := pdubuilder.CreateReportRanparameterItem(10102, "Cell Results")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter7)

	ranParameter8, err := pdubuilder.CreateReportRanparameterItem(10103, "SSB Results")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter8)

	ranParameter9, err := pdubuilder.CreateReportRanparameterItem(12501, "RSRP")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter9)

	ranParameter10, err := pdubuilder.CreateReportRanparameterItem(12502, "RSRQ")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter10)

	ranParameter11, err := pdubuilder.CreateReportRanparameterItem(12503, "SINR")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter11)
	return reportParametersStyle2List, nil
}
