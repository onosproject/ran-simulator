// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm2

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

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
	eventTriggerDefinition := &e2smkpmv2.E2SmKpmEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return 0, err
	}
	reportPeriod := eventTriggerDefinition.GetEventDefinitionFormat1().GetReportingPeriod()
	return reportPeriod, nil
}

func (sm *Client) getModelPlugin() (modelplugins.ServiceModel, error) {
	modelPlugin, err := sm.ServiceModel.ModelPluginRegistry.GetPlugin(ranFunctionE2SmOid)
	if err != nil {
		return nil, errors.New(errors.NotFound, "model plugin for model %s not found", ranFunctionShortName)
	}

	return modelPlugin, nil
}
