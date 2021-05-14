// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mho

import (
	"context"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"

	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/v1/e2sm-mho"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

//func (sm *Client) getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2sm_mho.E2SmMhoControlMessage, error) {
//	modelPlugin, err := sm.getModelPlugin()
//	if err != nil {
//		return nil, err
//	}
//	controlMessageProtoBytes, err := modelPlugin.ControlMessageASN1toProto(request.ProtocolIes.E2ApProtocolIes23.Value.Value)
//	if err != nil {
//		return nil, err
//	}
//	controlMessage := &e2sm_mho.E2SmMhoControlMessage{}
//	err = proto.Unmarshal(controlMessageProtoBytes, controlMessage)
//
//	if err != nil {
//		return nil, err
//	}
//	return controlMessage, nil
//}

//func (sm *Client) getControlHeader(request *e2appducontents.RiccontrolRequest) (*e2sm_mho.E2SmMhoControlHeader, error) {
//	modelPlugin, err := sm.getModelPlugin()
//	if err != nil {
//		return nil, err
//	}
//	controlHeaderProtoBytes, err := modelPlugin.ControlHeaderASN1toProto(request.ProtocolIes.E2ApProtocolIes22.Value.Value)
//	if err != nil {
//		return nil, err
//	}
//	controlHeader := &e2sm_mho.E2SmMhoControlHeader{}
//	err = proto.Unmarshal(controlHeaderProtoBytes, controlHeader)
//	if err != nil {
//		return nil, err
//	}
//
//	return controlHeader, nil
//}

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

//func (sm *Client) getPlmnID() ransimtypes.Uint24 {
//	plmnIDUint24 := ransimtypes.Uint24{}
//	plmnIDUint24.Set(uint32(sm.ServiceModel.Model.PlmnID))
//	return plmnIDUint24
//}

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

// createRicIndication creates ric indication  for each cell in the node
func (sm *Client) createRicIndication(ctx context.Context, ecgi ransimtypes.ECGI, subscription *subutils.Subscription) (*e2appducontents.Ricindication, error) {
	//plmnID := sm.getPlmnID()
	//var neighbourList []*e2sm_mho.Nrt
	//neighbourList = make([]*e2sm_mho.Nrt, 0)
	//cell, err := sm.ServiceModel.CellStore.Get(ctx, ecgi)
	//if err != nil {
	//	return nil, err
	//}
	//cellPci, err := sm.getCellPci(ctx, ecgi)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//earfcn, err := sm.getEarfcn(ctx, ecgi)
	//if err != nil {
	//	return nil, err
	//}

	//cellSize, err := sm.getCellSize(ctx, ecgi)
	//if err != nil {
	//	return nil, err
	//}
	//for _, neighbourEcgi := range cell.Neighbors {
	//	neighbourCellPci, err := sm.getCellPci(ctx, neighbourEcgi)
	//	if err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}
	//	neighbourEarfcn, err := sm.getEarfcn(ctx, neighbourEcgi)
	//	if err != nil {
	//		return nil, err
	//	}
	//	neighbourCellSize, err := sm.getCellSize(ctx, neighbourEcgi)
	//	if err != nil {
	//		return nil, err
	//	}
	//	neighbourEci := ransimtypes.GetECI(uint64(neighbourEcgi))
	//	neighbour, err := nrt.NewNeighbour(
	//		nrt.WithPci(neighbourCellPci),
	//		nrt.WithEutraCellIdentity(uint64(neighbourEci)),
	//		nrt.WithEarfcn(neighbourEarfcn),
	//		nrt.WithCellSize(sm.toCellSizeEnum(neighbourCellSize)),
	//		nrt.WithPlmnID(plmnID.Value())).Build()
	//	if err == nil {
	//		neighbourList = append(neighbourList, neighbour)
	//	}
	//}

	//cellEci := ransimtypes.GetECI(uint64(cell.ECGI))
	//// Creates MHO indication header
	//header := mhoindicationhdr.NewIndicationHeader(
	//	mhoindicationhdr.WithPlmnID(plmnID.Value()),
	//	mhoindicationhdr.WithEutracellIdentity(uint64(cellEci)))

	//// Creates MHO indication message

	//message := mhoindicationmsg.NewIndicationMessage(mhoindicationmsg.WithPlmnID(plmnID.Value()),
	//	mhoindicationmsg.WithCellSize(sm.toCellSizeEnum(cellSize)),
	//	mhoindicationmsg.WithEarfcn(earfcn),
	//	mhoindicationmsg.WithEutraCellIdentity(uint64(cellEci)),
	//	mhoindicationmsg.WithPci(cellPci),
	//	mhoindicationmsg.WithNeighbours(neighbourList))

	//mhoModelPlugin, _ := sm.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(sm.ServiceModel.OID))
	//indicationHeaderAsn1Bytes, err := header.ToAsn1Bytes(mhoModelPlugin)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}

	//indicationMessageAsn1Bytes, err := message.ToAsn1Bytes(mhoModelPlugin)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}

	//// Creates e2 indication
	//indication := indicationutils.NewIndication(
	//	indicationutils.WithRicInstanceID(subscription.GetRicInstanceID()),
	//	indicationutils.WithRanFuncID(subscription.GetRanFuncID()),
	//	indicationutils.WithRequestID(subscription.GetReqID()),
	//	indicationutils.WithIndicationHeader(indicationHeaderAsn1Bytes),
	//	indicationutils.WithIndicationMessage(indicationMessageAsn1Bytes))

	//ricIndication, err := indication.Build()
	//if err != nil {
	//	log.Error("creating indication message is failed", err)
	//	return nil, err
	//}
	//return ricIndication, nil

	return nil, nil
}
