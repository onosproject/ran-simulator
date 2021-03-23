// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm

import (
	e2sm_kpm_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

func getReportPeriods() map[string]int32 {
	return map[string]int32{
		"RT_PERIOD_IE_MS10":    10,
		"RT_PERIOD_IE_MS20":    20,
		"RT_PERIOD_IE_MS32":    32,
		"RT_PERIOD_IE_MS40":    40,
		"RT_PERIOD_IE_MS60":    60,
		"RT_PERIOD_IE_MS64":    64,
		"RT_PERIOD_IE_MS70":    70,
		"RT_PERIOD_IE_MS80":    80,
		"RT_PERIOD_IE_MS128":   128,
		"RT_PERIOD_IE_MS160":   160,
		"RT_PERIOD_IE_MS256":   256,
		"RT_PERIOD_IE_MS320":   320,
		"RT_PERIOD_IE_MS512":   512,
		"RT_PERIOD_IE_MS640":   640,
		"RT_PERIOD_IE_MS1024":  1024,
		"RT_PERIOD_IE_MS1280":  1280,
		"RT_PERIOD_IE_MS2048":  2048,
		"RT_PERIOD_IE_MS2560":  2560,
		"RT_PERIOD_IE_MS5120":  5120,
		"RT_PERIOD_IE_MS10240": 10240,
	}
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
	eventTriggerDefinition := &e2sm_kpm_ies.E2SmKpmEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return 0, err
	}
	reportPeriod := eventTriggerDefinition.GetEventDefinitionFormat1().PolicyTestList[0].ReportPeriodIe.Enum().String()
	interval := getReportPeriods()[reportPeriod]
	return interval, nil
}

func (sm *Client) getModelPlugin() (modelplugins.ServiceModel, error) {
	modelPlugin, err := sm.ServiceModel.ModelPluginRegistry.GetPlugin(modelOID)
	if err != nil {
		return nil, errors.New(errors.NotFound, "model plugin for model %s not found", modelName)
	}

	return modelPlugin, nil
}
