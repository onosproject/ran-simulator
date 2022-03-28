// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package rc

import (
	"context"
	"fmt"

	e2smrcpresm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/servicemodel"
	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"

	meastype "github.com/onosproject/rrm-son-lib/pkg/model/measurement/type"

	"github.com/onosproject/ran-simulator/pkg/model"

	indicationutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/indication"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	rcindicationhdr "github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/indication/header"
	rcindicationmsg "github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/indication/message"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/nrt"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/v2/e2sm-rc-pre-v2-go"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"google.golang.org/protobuf/proto"
)

func (sm *Client) getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2smrcpreies.E2SmRcPreControlMessage, error) {
	var rcPreServiceModel e2smrcpresm.RcPreServiceModel
	var controlMessageAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolMessage) {
			controlMessageAsnBytes = v.GetValue().GetRiccontrolMessage().GetValue()
			break
		}
	}
	controlMessageProtoBytes, err := rcPreServiceModel.ControlMessageASN1toProto(controlMessageAsnBytes)
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
	var rcPreServiceModel e2smrcpresm.RcPreServiceModel
	var controlHeaderAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolHeader) {
			controlHeaderAsnBytes = v.GetValue().GetRiccontrolHeader().GetValue()
			break
		}
	}
	controlHeaderProtoBytes, err := rcPreServiceModel.ControlHeaderASN1toProto(controlHeaderAsnBytes)
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
	var eventTriggerAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			eventTriggerAsnBytes = v.GetValue().GetRicsubscriptionDetails().GetRicEventTriggerDefinition().GetValue()
			break
		}
	}

	var rcPreServiceModel e2smrcpresm.RcPreServiceModel
	eventTriggerProtoBytes, err := rcPreServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return -1, err
	}
	eventTriggerDefinition := &e2smrcpreies.E2SmRcPreEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return -1, err
	}
	eventTriggerType := eventTriggerDefinition.GetEventDefinitionFormats().GetEventDefinitionFormat1().TriggerType
	return eventTriggerType, nil
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

func (sm *Client) getCellPCI(ctx context.Context, ncgi ransimtypes.NCGI) (int32, error) {
	cell, err := sm.ServiceModel.CellStore.Get(ctx, ncgi)
	if err != nil {
		return 0, err
	}

	return int32(cell.PCI), nil
}

func (sm *Client) getEARFCN(ctx context.Context, ncgi ransimtypes.NCGI) (int32, error) {
	cell, err := sm.ServiceModel.CellStore.Get(ctx, ncgi)
	if err != nil {
		return 0, err
	}

	return int32(cell.Earfcn), nil
}

func (sm *Client) getCellSize(ctx context.Context, ncgi ransimtypes.NCGI) (string, error) {
	cell, err := sm.ServiceModel.CellStore.Get(ctx, ncgi)
	if err != nil {
		return "", err
	}
	return cell.CellType.String(), nil
}

// getReportPeriod extracts report period
func (sm *Client) getReportPeriod(request *e2appducontents.RicsubscriptionRequest) (uint32, error) {
	var eventTriggerAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			eventTriggerAsnBytes = v.GetValue().GetRicsubscriptionDetails().GetRicEventTriggerDefinition().GetValue()
			break
		}
	}

	var rcPreServiceModel e2smrcpresm.RcPreServiceModel
	eventTriggerProtoBytes, err := rcPreServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return 0, err
	}
	eventTriggerDefinition := &e2smrcpreies.E2SmRcPreEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return 0, err
	}
	reportPeriod := eventTriggerDefinition.GetEventDefinitionFormats().GetEventDefinitionFormat1().ReportingPeriodMs

	if reportPeriod == nil {
		return 0, fmt.Errorf("no reporting period was set, obtained %v", reportPeriod)
	}
	rp := uint32(*reportPeriod)

	return rp, nil
}

// createRicIndication creates ric indication  for each cell in the node
func (sm *Client) createRicIndication(ctx context.Context, ncgi ransimtypes.NCGI, subscription *subutils.Subscription) (*e2appducontents.Ricindication, error) {
	plmnID := sm.getPlmnID()
	var neighbourList []*e2smrcpreies.Nrt
	neighbourList = make([]*e2smrcpreies.Nrt, 0)
	cell, err := sm.ServiceModel.CellStore.Get(ctx, ncgi)
	if err != nil {
		return nil, err
	}
	cellPci, err := sm.getCellPCI(ctx, ncgi)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	earfcn, err := sm.getEARFCN(ctx, ncgi)
	if err != nil {
		return nil, err
	}

	cellSize, err := sm.getCellSize(ctx, ncgi)
	if err != nil {
		return nil, err
	}
	for _, neighbourNcgi := range cell.Neighbors {
		neighbourCellPci, err := sm.getCellPCI(ctx, neighbourNcgi)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		neighbourEarfcn, err := sm.getEARFCN(ctx, neighbourNcgi)
		if err != nil {
			return nil, err
		}
		neighbourCellSize, err := sm.getCellSize(ctx, neighbourNcgi)
		if err != nil {
			return nil, err
		}
		neighbourEci := ransimtypes.GetNCI(neighbourNcgi)
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

	cellEci := ransimtypes.GetNCI(cell.NCGI)
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

	indicationHeaderAsn1Bytes, err := header.ToAsn1Bytes()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	indicationMessageAsn1Bytes, err := message.ToAsn1Bytes()
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

func setPCI(parameterName string, parameterValue interface{}, cell *model.Cell) {
	if parameterName == "pci" {
		switch parameterValue := parameterValue.(type) {
		case int32:
			cell.PCI = uint32(parameterValue)
		case uint32:
			cell.PCI = parameterValue
		case int64:
			cell.PCI = uint32(parameterValue)
		case uint64:
			cell.PCI = uint32(parameterValue)
		}
	}
}

func (sm *Client) setHandoverOcn(ctx context.Context, parameterName string, parameterValue interface{}, cell *model.Cell) {
	var ocnRc meastype.QOffsetRange
	nCellNCGI := cell.NCGI

	if parameterName == "ocn_rc" {
		switch parameterValue := parameterValue.(type) {
		case int32:
			ocnRc = meastype.QOffsetRange(parameterValue)
		case uint32:
			ocnRc = meastype.QOffsetRange(parameterValue)
		case int64:
			ocnRc = meastype.QOffsetRange(parameterValue)
		case uint64:
			ocnRc = meastype.QOffsetRange(parameterValue)
		}

		for _, ncgi := range sm.ServiceModel.Node.Cells {
			if ncgi == nCellNCGI {
				continue
			}
			sCell, err := sm.ServiceModel.CellStore.Get(ctx, ncgi)
			if err != nil {
				log.Errorf("NCGI (%v) is not in cell store")
			}
			if _, ok := sCell.MeasurementParams.NCellIndividualOffsets[nCellNCGI]; !ok {
				log.Errorf("the cell NCGI (%v) is not a neighbor of the cell NCGI (%v)", nCellNCGI, ncgi)
				continue
			}
			log.Debugf("Cell (%v) Ocn in the cell (%v) is set from %v to %v", cell.NCGI, ncgi, sCell.MeasurementParams.NCellIndividualOffsets[nCellNCGI], ocnRc.GetValue().(int))
			sCell.MeasurementParams.NCellIndividualOffsets[nCellNCGI] = int32(ocnRc.GetValue().(int))
		}
	}

}
