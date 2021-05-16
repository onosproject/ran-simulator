// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package rc

import (
	"context"
	"strconv"

	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"

	indicationutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/indication"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	rcindicationhdr "github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/indication/header"
	rcindicationmsg "github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/indication/message"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/nrt"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v2/e2sm-rc-pre-v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

func (sm *Client) getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2smrcpreies.E2SmRcPreControlMessage, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		return nil, err
	}
	controlMessageProtoBytes, err := modelPlugin.ControlMessageASN1toProto(request.ProtocolIes.E2ApProtocolIes23.Value.Value)
	if err != nil {
		return nil, err
	}
	controlMessage := &e2smrcpreies.E2SmRcPreControlMessage{}
	err = proto.Unmarshal(controlMessageProtoBytes, controlMessage)

	if err != nil {
		return nil, err
	}
	return controlMessage, nil
}

func (sm *Client) getControlHeader(request *e2appducontents.RiccontrolRequest) (*e2smrcpreies.E2SmRcPreControlHeader, error) {
	modelPlugin, err := sm.getModelPlugin()
	if err != nil {
		return nil, err
	}
	controlHeaderProtoBytes, err := modelPlugin.ControlHeaderASN1toProto(request.ProtocolIes.E2ApProtocolIes22.Value.Value)
	if err != nil {
		return nil, err
	}
	controlHeader := &e2smrcpreies.E2SmRcPreControlHeader{}
	err = proto.Unmarshal(controlHeaderProtoBytes, controlHeader)
	if err != nil {
		return nil, err
	}

	return controlHeader, nil
}

// getEventTriggerType extracts event trigger type
func (sm *Client) getEventTriggerType(request *e2appducontents.RicsubscriptionRequest) (e2smrcpreies.RcPreTriggerType, error) {
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
	eventTriggerDefinition := &e2smrcpreies.E2SmRcPreEventTriggerDefinition{}
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

func (sm *Client) getPlmnID() ransimtypes.Uint24 {
	plmnIDUint24 := ransimtypes.Uint24{}
	plmnIDUint24.Set(uint32(sm.ServiceModel.Model.PlmnID))
	return plmnIDUint24
}

func (sm *Client) toCellSizeEnum(cellSize string) e2smrcpreies.CellSize {
	switch cellSize {
	case "ENTERPRISE":
		return e2smrcpreies.CellSize_CELL_SIZE_ENTERPRISE
	case "FEMTO":
		return e2smrcpreies.CellSize_CELL_SIZE_FEMTO
	case "MACRO":
		return e2smrcpreies.CellSize_CELL_SIZE_MACRO
	case "OUTDOOR_SMALL":
		return e2smrcpreies.CellSize_CELL_SIZE_OUTDOOR_SMALL
	default:
		return e2smrcpreies.CellSize_CELL_SIZE_ENTERPRISE
	}
}

func (sm *Client) getCellPci(ctx context.Context, ecgi ransimtypes.ECGI) (int32, error) {
	cellPci, found := sm.ServiceModel.MetricStore.Get(ctx, uint64(ecgi), "pci")
	if !found {
		return 0, errors.New(errors.NotFound, "pci value is not found for cell:", ecgi)
	}
	// TODO we should handle this properly in metric store
	switch cellPci := cellPci.(type) {
	case uint32:
		return int32(cellPci), nil
	case int32:
		return cellPci, nil
	case int64:
		return int32(cellPci), nil
	case uint64:
		return int32(cellPci), nil
	case uint8:
		return int32(cellPci), nil
	case int8:
		return int32(cellPci), nil
	case int16:
		return int32(cellPci), nil
	case uint16:
		return int32(cellPci), nil
	case string:
		val, err := strconv.Atoi(cellPci)
		if err != nil {
			return 0, err
		}
		return int32(val), nil
	default:
		return 0, nil
	}

}

func (sm *Client) getEarfcn(ctx context.Context, ecgi ransimtypes.ECGI) (int32, error) {
	earfcn, found := sm.ServiceModel.MetricStore.Get(ctx, uint64(ecgi), "earfcn")
	if !found {
		return 0, errors.New(errors.NotFound, "earfc value is not found for cell:", ecgi)
	}

	return int32(earfcn.(uint32)), nil
}

func (sm *Client) getCellSize(ctx context.Context, ecgi ransimtypes.ECGI) (string, error) {
	cellSize, found := sm.ServiceModel.MetricStore.Get(ctx, uint64(ecgi), "cellSize")
	if !found {
		return "", errors.New(errors.NotFound, "cell size value is not found for  neighbour  cell:", ecgi)
	}
	return cellSize.(string), nil
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
	eventTriggerDefinition := &e2smrcpreies.E2SmRcPreEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return 0, err
	}
	reportPeriod := eventTriggerDefinition.GetEventDefinitionFormat1().ReportingPeriodMs
	return reportPeriod, nil
}

// createRicIndication creates ric indication  for each cell in the node
func (sm *Client) createRicIndication(ctx context.Context, ecgi ransimtypes.ECGI, subscription *subutils.Subscription) (*e2appducontents.Ricindication, error) {
	plmnID := sm.getPlmnID()
	var neighbourList []*e2smrcpreies.Nrt
	neighbourList = make([]*e2smrcpreies.Nrt, 0)
	cell, err := sm.ServiceModel.CellStore.Get(ctx, ecgi)
	if err != nil {
		return nil, err
	}
	cellPci, err := sm.getCellPci(ctx, ecgi)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	earfcn, err := sm.getEarfcn(ctx, ecgi)
	if err != nil {
		return nil, err
	}

	cellSize, err := sm.getCellSize(ctx, ecgi)
	if err != nil {
		return nil, err
	}
	for _, neighbourEcgi := range cell.Neighbors {
		neighbourCellPci, err := sm.getCellPci(ctx, neighbourEcgi)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		neighbourEarfcn, err := sm.getEarfcn(ctx, neighbourEcgi)
		if err != nil {
			return nil, err
		}
		neighbourCellSize, err := sm.getCellSize(ctx, neighbourEcgi)
		if err != nil {
			return nil, err
		}
		neighbourEci := ransimtypes.GetECI(uint64(neighbourEcgi))
		neighbour, err := nrt.NewNeighbour(
			nrt.WithPci(neighbourCellPci),
			nrt.WithNrcellIdentity(uint64(neighbourEci)),
			nrt.WithEarfcn(neighbourEarfcn),
			nrt.WithCellSize(sm.toCellSizeEnum(neighbourCellSize)),
			nrt.WithPlmnID(plmnID.Value())).Build()
		if err == nil {
			neighbourList = append(neighbourList, neighbour)
		}
	}

	cellEci := ransimtypes.GetECI(uint64(cell.ECGI))
	// Creates RC indication header
	header := rcindicationhdr.NewIndicationHeader(
		rcindicationhdr.WithPlmnID(plmnID.Value()),
		rcindicationhdr.WithNRcellIdentity(uint64(cellEci)))
	// Creates RC indication message

	message := rcindicationmsg.NewIndicationMessage(rcindicationmsg.WithPlmnID(plmnID.Value()),
		rcindicationmsg.WithCellSize(sm.toCellSizeEnum(cellSize)),
		rcindicationmsg.WithEarfcn(earfcn),
		rcindicationmsg.WithPci(cellPci),
		rcindicationmsg.WithNeighbours(neighbourList))

	rcModelPlugin, _ := sm.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(sm.ServiceModel.OID))
	indicationHeaderAsn1Bytes, err := header.ToAsn1Bytes(rcModelPlugin)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	indicationMessageAsn1Bytes, err := message.ToAsn1Bytes(rcModelPlugin)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Creates e2 indication
	indication := indicationutils.NewIndication(
		indicationutils.WithRicInstanceID(subscription.GetRicInstanceID()),
		indicationutils.WithRanFuncID(subscription.GetRanFuncID()),
		indicationutils.WithRequestID(subscription.GetReqID()),
		indicationutils.WithIndicationHeader(indicationHeaderAsn1Bytes),
		indicationutils.WithIndicationMessage(indicationMessageAsn1Bytes))

	ricIndication, err := indication.Build()
	if err != nil {
		log.Error("creating indication message is failed", err)
		return nil, err
	}
	return ricIndication, nil
}
