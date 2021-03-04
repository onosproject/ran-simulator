// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package rc

import (
	"context"
	"strconv"

	indicationutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/indication"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	rcindicationhdr "github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/indication/header"
	rcindicationmsg "github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/indication/message"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/nrt"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/pcirange"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/rc/pciload"

	"github.com/onosproject/onos-api/go/onos/ransim/types"

	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	uint24 "github.com/onosproject/ran-simulator/pkg/types"
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

func (sm *Client) getPlmnID() uint24.Uint24 {
	plmnIDUint24 := uint24.Uint24{}
	plmnIDUint24.Set(uint32(sm.ServiceModel.Model.PlmnID))
	return plmnIDUint24
}

func (sm *Client) toCellSizeEnum(cellSize string) e2sm_rc_pre_ies.CellSize {
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

func (sm *Client) getCellPci(ctx context.Context, ecgi types.ECGI) (int32, error) {
	cellPci, found := sm.ServiceModel.MetricStore.Get(ctx, uint64(ecgi), "pci")
	if !found {
		return 0, errors.New(errors.NotFound, "pci value is not found for cell:", ecgi)
	}
	return int32(cellPci.(uint32)), nil
}

func (sm *Client) getEarfcn(ctx context.Context, ecgi types.ECGI) (int32, error) {
	earfcn, found := sm.ServiceModel.MetricStore.Get(ctx, uint64(ecgi), "earfcn")
	if !found {
		return 0, errors.New(errors.NotFound, "earfc value is not found for cell:", ecgi)
	}

	return int32(earfcn.(uint32)), nil
}

func (sm *Client) getCellSize(ctx context.Context, ecgi types.ECGI) (string, error) {
	cellSize, found := sm.ServiceModel.MetricStore.Get(ctx, uint64(ecgi), "cellSize")
	if !found {
		return "", errors.New(errors.NotFound, "cell size value is not found for  neighbour  cell:", ecgi)
	}
	return cellSize.(string), nil
}

func (sm *Client) getPciPool(ctx context.Context, ecgi types.ECGI) ([]pciload.PciRange, error) {
	pciPool, found := sm.ServiceModel.MetricStore.Get(ctx, uint64(ecgi), "pcipool")
	if !found {
		return nil, errors.New(errors.NotFound, "cell size value is not found for  neighbour  cell:", ecgi)
	}
	return pciPool.([]pciload.PciRange), nil
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
	eventTriggerDefinition := &e2sm_rc_pre_ies.E2SmRcPreEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return 0, err
	}
	reportPeriod := eventTriggerDefinition.GetEventDefinitionFormat1().ReportingPeriodMs
	return reportPeriod, nil
}

// createRicIndication creates ric indication  for each cell in the node
func (sm *Client) createRicIndication(ctx context.Context, ecgi types.ECGI, subscription *subutils.Subscription) (*e2appducontents.Ricindication, error) {
	plmnID := sm.getPlmnID()
	var neighbourList []*e2sm_rc_pre_ies.Nrt
	neighbourList = make([]*e2sm_rc_pre_ies.Nrt, 0)
	cell, err := sm.ServiceModel.CellStore.Get(ctx, ecgi)
	if err != nil {
		return nil, err
	}
	cellPci, err := sm.getCellPci(ctx, ecgi)
	if err != nil {
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
	for index, neighbourEcgi := range cell.Neighbors {
		neighbourCellPci, err := sm.getCellPci(ctx, neighbourEcgi)
		if err != nil {
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
		neighbour, err := nrt.NewNeighbour(
			nrt.WithNrIndex(int32(index)),
			nrt.WithPci(neighbourCellPci),
			nrt.WithEutraCellIdentity(uint64(neighbourEcgi)),
			nrt.WithEarfcn(neighbourEarfcn),
			nrt.WithCellSize(sm.toCellSizeEnum(neighbourCellSize)),
			nrt.WithPlmnID(plmnID.Value())).Build()
		if err == nil {
			neighbourList = append(neighbourList, neighbour)
		}
	}

	pciRanges, err := sm.getPciPool(ctx, ecgi)
	if err != nil {
		return nil, err
	}

	var pciPool []*e2smrcpreies.PciRange
	for _, pciRangeValue := range pciRanges {
		pciRange, err := pcirange.NewPciRange(pcirange.WithLowerPci(int32(pciRangeValue.Min)),
			pcirange.WithUpperPci(int32(pciRangeValue.Max))).Build()
		if err != nil {
			return nil, err
		}
		pciPool = append(pciPool, pciRange)
	}

	// Creates RC indication header
	header := rcindicationhdr.NewIndicationHeader(
		rcindicationhdr.WithPlmnID(plmnID.Value()),
		rcindicationhdr.WithEutracellIdentity(uint64(cell.ECGI)))

	// Creates RC indication message
	message := rcindicationmsg.NewIndicationMessage(rcindicationmsg.WithPlmnID(plmnID.Value()),
		rcindicationmsg.WithCellSize(sm.toCellSizeEnum(cellSize)),
		rcindicationmsg.WithEarfcn(earfcn),
		rcindicationmsg.WithEutraCellIdentity(uint64(cell.ECGI)),
		rcindicationmsg.WithPci(cellPci),
		rcindicationmsg.WithNeighbours(neighbourList),
		rcindicationmsg.WithPciPool(pciPool))

	rcModelPlugin := sm.ServiceModel.ModelPluginRegistry.ModelPlugins[sm.ServiceModel.ModelFullName]
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

func toEcgi(plmnIDByteArray []byte, eci uint64) string {
	plmnID := uint24.Uint24ToUint32(plmnIDByteArray)
	plmnIDString := strconv.FormatUint(uint64(plmnID), 10)
	eciString := strconv.FormatUint(eci, 10)
	ecgi := plmnIDString + eciString
	return ecgi
}
