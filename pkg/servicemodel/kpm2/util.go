// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm2

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/labelinfo"
	"google.golang.org/protobuf/proto"
)

func (sm *Client) getActionDefinition(actionList []*e2appducontents.RicactionToBeSetupItemIes, ricActionsAccepted []*e2aptypes.RicActionID) ([]*e2smkpmv2.E2SmKpmActionDefinition, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	var actionDefinitions []*e2smkpmv2.E2SmKpmActionDefinition
	for _, action := range actionList {
		for _, acceptedActionID := range ricActionsAccepted {
			if action.Value.RicActionId.Value == int32(*acceptedActionID) {
				actionDefinitionBytes := action.Value.RicActionDefinition.Value
				actionDefinitionProtoBytes, err := modelPlugin.ActionDefinitionASN1toProto(actionDefinitionBytes)
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

func (sm *Client) createInfoLabelList() (*e2smkpmv2.LabelInfoList, error) {
	plmnID := ransimtypes.NewUint24(uint32(sm.ServiceModel.Model.PlmnID))
	var fiveQI int32 = 10
	var qfi int32 = 62
	var qci int32 = 15
	var qciMin int32 = 1
	var qciMax int32 = 15
	var arpMax int32 = 15
	var arpMin int32 = 10
	var bitrateRange int32 = 251
	var layerMuMimo int32 = 5
	var distX int32 = 123
	var distY int32 = 456
	var distZ int32 = 789
	startEndIndication := e2smkpmv2.StartEndInd_START_END_IND_START
	sst := []byte{0x01}
	sd := []byte{0x01, 0x02, 0x03}

	// Creates label information
	labelInfo, err := labelinfo.NewLabelInfo(labelinfo.WithFiveQI(fiveQI),
		labelinfo.WithPlmnID(plmnID.Value()),
		labelinfo.WithArpMin(arpMin),
		labelinfo.WithArpMax(arpMax),
		labelinfo.WithBitRateRange(bitrateRange),
		labelinfo.WithDistX(distX),
		labelinfo.WithDistY(distY),
		labelinfo.WithDistZ(distZ),
		labelinfo.WithLayerMuMimo(layerMuMimo),
		labelinfo.WithQCI(qci),
		labelinfo.WithQCIMin(qciMin),
		labelinfo.WithQCIMax(qciMax),
		labelinfo.WithQFI(qfi),
		labelinfo.WithStartEndIndication(startEndIndication),
		labelinfo.WithSD(sd),
		labelinfo.WithSST(sst))

	if err != nil {
		log.Warn(err)
		return nil, err
	}
	labelInfoItem, err := labelInfo.Build()
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	labelInfoList := e2smkpmv2.LabelInfoList{
		Value: make([]*e2smkpmv2.LabelInfoItem, 0),
	}
	labelInfoList.Value = append(labelInfoList.Value, labelInfoItem)

	return &labelInfoList, nil
}
