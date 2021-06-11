// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mho

import (
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/v1/e2sm-mho"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

func (sm *Client) getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2sm_mho.E2SmMhoControlMessage, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		return nil, err
	}
	controlMessageProtoBytes, err := modelPlugin.ControlMessageASN1toProto(request.ProtocolIes.E2ApProtocolIes23.Value.Value)
	if err != nil {
		return nil, err
	}
	controlMessage := &e2sm_mho.E2SmMhoControlMessage{}
	err = proto.Unmarshal(controlMessageProtoBytes, controlMessage)

	if err != nil {
		return nil, err
	}
	return controlMessage, nil
}

func (sm *Client) getControlHeader(request *e2appducontents.RiccontrolRequest) (*e2sm_mho.E2SmMhoControlHeader, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		return nil, err
	}
	controlHeaderProtoBytes, err := modelPlugin.ControlHeaderASN1toProto(request.ProtocolIes.E2ApProtocolIes22.Value.Value)
	if err != nil {
		return nil, err
	}
	controlHeader := &e2sm_mho.E2SmMhoControlHeader{}
	err = proto.Unmarshal(controlHeaderProtoBytes, controlHeader)
	if err != nil {
		return nil, err
	}

	return controlHeader, nil
}

// getEventTriggerType extracts event trigger type
func (sm *Client) getEventTriggerType(request *e2appducontents.RicsubscriptionRequest) (e2sm_mho.MhoTriggerType, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	eventTriggerAsnBytes := request.ProtocolIes.E2ApProtocolIes30.Value.RicEventTriggerDefinition.Value
	eventTriggerProtoBytes, err := modelPlugin.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return -1, err
	}
	eventTriggerDefinition := &e2sm_mho.E2SmMhoEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return -1, err
	}
	eventTriggerType := eventTriggerDefinition.GetEventDefinitionFormat1().TriggerType
	return eventTriggerType, nil
}

func (sm *Client) getModelPlugin() (modelplugins.ServiceModel, error) {
	modelPlugin, err := sm.ServiceModel.ModelPluginRegistry.GetPlugin(modelOID)
	if err != nil {
		return nil, errors.New(errors.NotFound, "model plugin for model %s not found", modelFullName)
	}

	return modelPlugin, nil
}

// getReportPeriod extracts report period
func (sm *Client) getReportPeriod(request *e2appducontents.RicsubscriptionRequest) (int32, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	eventTriggerAsnBytes := request.ProtocolIes.E2ApProtocolIes30.Value.RicEventTriggerDefinition.Value
	eventTriggerProtoBytes, err := modelPlugin.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return 0, err
	}
	eventTriggerDefinition := &e2sm_mho.E2SmMhoEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return 0, err
	}
	reportPeriod := eventTriggerDefinition.GetEventDefinitionFormat1().ReportingPeriodMs
	return reportPeriod, nil
}
