// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

import (
	"encoding/hex"

	e2smmhosm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/servicemodel"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"
	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func (m *Mho) getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2sm_mho.E2SmMhoControlMessage, error) {
	var mhoServiceModel e2smmhosm.MhoServiceModel
	var controlMessageAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolMessage) {
			controlMessageAsnBytes = v.GetValue().GetRiccontrolMessage().GetValue()
			break
		}
	}
	log.Debugf("ControlMessage is\n%v", hex.Dump(controlMessageAsnBytes))

	controlMessageProtoBytes, err := mhoServiceModel.ControlMessageASN1toProto(controlMessageAsnBytes)
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
	var controlHeaderAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolHeader) {
			controlHeaderAsnBytes = v.GetValue().GetRiccontrolHeader().GetValue()
			break
		}
	}
	log.Debugf("ControlHeader is\n%v", hex.Dump(controlHeaderAsnBytes))

	controlHeaderProtoBytes, err := mhoServiceModel.ControlHeaderASN1toProto(controlHeaderAsnBytes)
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
	var eventTriggerAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			eventTriggerAsnBytes = v.GetValue().GetRicsubscriptionDetails().GetRicEventTriggerDefinition().GetValue()
			break
		}
	}

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
	var eventTriggerAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			eventTriggerAsnBytes = v.GetValue().GetRicsubscriptionDetails().GetRicEventTriggerDefinition().GetValue()
			break
		}
	}

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
		return 0, errors.NewNotFound("no reporting period was set, obtained %v", reportPeriod)
	}
	rp := *reportPeriod

	return rp, nil
}
