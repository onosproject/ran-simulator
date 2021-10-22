// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package mho

import (
	"fmt"
	e2smmhosm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/servicemodel"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v1/e2sm-mho-go"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"google.golang.org/protobuf/proto"
)

func (m *Mho) getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2sm_mho.E2SmMhoControlMessage, error) {
	var mhoServiceModel e2smmhosm.MhoServiceModel
	controlMessageProtoBytes, err := mhoServiceModel.ControlMessageASN1toProto(request.ProtocolIes.E2ApProtocolIes23.Value.Value)
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

func (m *Mho) getControlHeader(request *e2appducontents.RiccontrolRequest) (*e2sm_mho.E2SmMhoControlHeader, error) {
	var mhoServiceModel e2smmhosm.MhoServiceModel
	controlHeaderProtoBytes, err := mhoServiceModel.ControlHeaderASN1toProto(request.ProtocolIes.E2ApProtocolIes22.Value.Value)
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
func (m *Mho) getEventTriggerType(request *e2appducontents.RicsubscriptionRequest) (e2sm_mho.MhoTriggerType, error) {
	eventTriggerAsnBytes := request.ProtocolIes.E2ApProtocolIes30.Value.RicEventTriggerDefinition.Value

	var mhoServiceModel e2smmhosm.MhoServiceModel
	eventTriggerProtoBytes, err := mhoServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return -1, err
	}
	eventTriggerDefinition := &e2sm_mho.E2SmMhoEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return -1, err
	}
	eventTriggerType := eventTriggerDefinition.GetEventDefinitionFormats().GetEventDefinitionFormat1().TriggerType
	return eventTriggerType, nil
}

// getReportPeriod extracts report period
func (m *Mho) getReportPeriod(request *e2appducontents.RicsubscriptionRequest) (int32, error) {
	eventTriggerAsnBytes := request.ProtocolIes.E2ApProtocolIes30.Value.RicEventTriggerDefinition.Value

	var mhoServiceModel e2smmhosm.MhoServiceModel
	eventTriggerProtoBytes, err := mhoServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return 0, err
	}
	eventTriggerDefinition := &e2sm_mho.E2SmMhoEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return 0, err
	}
	reportPeriod := eventTriggerDefinition.GetEventDefinitionFormats().GetEventDefinitionFormat1().ReportingPeriodMs
	if reportPeriod == nil {
		return 0, fmt.Errorf("no reporting period was set, obtained %v", reportPeriod)
	}
	rp := *reportPeriod

	return rp, nil
}
