// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package rc

import (
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/types"
	"google.golang.org/protobuf/proto"
)

func (sm *Client) getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2sm_rc_pre_ies.E2SmRcPreControlMessage, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		return nil, err
	}
	controlMessageProtoBytes, err := modelPlugin.ControlMessageASN1toProto(request.ProtocolIes.E2ApProtocolIes23.Value.Value)
	if err != nil {
		return nil, err
	}
	controlMessage := &e2sm_rc_pre_ies.E2SmRcPreControlMessage{}
	err = proto.Unmarshal(controlMessageProtoBytes, controlMessage)

	if err != nil {
		return nil, err
	}
	return controlMessage, nil
}

func (sm *Client) getControlHeader(request *e2appducontents.RiccontrolRequest) (*e2sm_rc_pre_ies.E2SmRcPreControlHeader, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		return nil, err
	}
	controlHeaderProtoBytes, err := modelPlugin.ControlHeaderASN1toProto(request.ProtocolIes.E2ApProtocolIes22.Value.Value)
	if err != nil {
		return nil, err
	}
	controlHeader := &e2sm_rc_pre_ies.E2SmRcPreControlHeader{}
	err = proto.Unmarshal(controlHeaderProtoBytes, controlHeader)
	if err != nil {
		return nil, err
	}

	return controlHeader, nil
}

// getEventTriggerType extracts event trigger type
func (sm *Client) getEventTriggerType(request *e2appducontents.RicsubscriptionRequest) (e2sm_rc_pre_ies.RcPreTriggerType, error) {
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
	eventTriggerDefinition := &e2sm_rc_pre_ies.E2SmRcPreEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return -1, err
	}
	eventTriggerType := eventTriggerDefinition.GetEventDefinitionFormat1().TriggerType
	return eventTriggerType, nil
}

func (sm *Client) getModelPlugin() (modelplugins.ModelPlugin, error) {
	if modelPlugin, ok := sm.ServiceModel.ModelPluginRegistry.ModelPlugins[modelFullName]; ok {
		return modelPlugin, nil
	}
	return nil, errors.New(errors.NotFound, "model plugin for model %s not found", modelFullName)
}

func (sm *Client) getPlmnID() types.Uint24 {
	plmnIDUint24 := types.Uint24{}
	plmnIDUint24.Set(uint32(sm.ServiceModel.Model.PlmnID))
	return plmnIDUint24
}

func (sm *Client) getCellSize(cellSize string) e2sm_rc_pre_ies.CellSize {
	switch cellSize {
	case "ENTERPRISE":
		return e2sm_rc_pre_ies.CellSize_CELL_SIZE_ENTERPRISE
	case "FEMTO":
		return e2sm_rc_pre_ies.CellSize_CELL_SIZE_FEMTO
	case "MACRO":
		return e2sm_rc_pre_ies.CellSize_CELL_SIZE_MACRO
	case "OUTDOOR_SMALL":
		return e2sm_rc_pre_ies.CellSize_CELL_SIZE_OUTDOOR_SMALL
	default:
		return e2sm_rc_pre_ies.CellSize_CELL_SIZE_ENTERPRISE
	}
}
