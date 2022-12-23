// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/pdubuilder"
	e2smrc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/servicemodel"
	e2smcommonies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-common-ies"
	e2smrcies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-rc-ies"
	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/utils"
	indicationutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/indication"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/v1/indication/headers/format1"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/v1/indication/headers/format2"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/v1/indication/messages/format3"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/v1/indication/messages/format5"
	"google.golang.org/protobuf/proto"
	"math"
	"strconv"
)

func getActionDefinitionMap(actionList []*e2appducontents.RicactionToBeSetupItemIes, ricActionsAccepted []*e2aptypes.RicActionID) (map[*e2aptypes.RicActionID]*e2smrcies.E2SmRcActionDefinition, error) {
	actionDefinitionsMap := make(map[*e2aptypes.RicActionID]*e2smrcies.E2SmRcActionDefinition)
	for _, action := range actionList {
		for _, actionID := range ricActionsAccepted {
			if action.GetValue().GetRicactionToBeSetupItem().GetRicActionId().GetValue() == int32(*actionID) {
				actionDefinitionBytes := action.GetValue().GetRicactionToBeSetupItem().GetRicActionDefinition().GetValue()
				var rcServiceModel e2smrc.RCServiceModel

				actionDefinitionProtoBytes, err := rcServiceModel.ActionDefinitionASN1toProto(actionDefinitionBytes)
				if err != nil {
					return nil, err
				}

				actionDefinition := &e2smrcies.E2SmRcActionDefinition{}
				err = proto.Unmarshal(actionDefinitionProtoBytes, actionDefinition)
				if err != nil {
					return nil, err
				}
				actionDefinitionsMap[actionID] = actionDefinition
			}
		}
	}
	return actionDefinitionsMap, nil
}

func getEventTrigger(request *e2appducontents.RicsubscriptionRequest) (*e2smrcies.E2SmRcEventTrigger, error) {
	var eventTriggerAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			eventTriggerAsnBytes = v.GetValue().GetRicsubscriptionDetails().GetRicEventTriggerDefinition().GetValue()
			break
		}
	}

	var rcServiceModel e2smrc.RCServiceModel
	eventTriggerProtoBytes, err := rcServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return nil, err
	}
	eventTriggerDefinition := &e2smrcies.E2SmRcEventTrigger{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return nil, err
	}

	return eventTriggerDefinition, nil
}

func getControlMessage(request *e2appducontents.RiccontrolRequest) (*e2smrcies.E2SmRcControlMessage, error) {
	var rcServiceModel e2smrc.RCServiceModel
	var controlMessageAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolMessage) {
			controlMessageAsnBytes = v.GetValue().GetRiccontrolMessage().GetValue()
			break
		}
	}
	controlMessageProtoBytes, err := rcServiceModel.ControlMessageASN1toProto(controlMessageAsnBytes)
	if err != nil {
		return nil, err
	}
	controlMessage := &e2smrcies.E2SmRcControlMessage{}
	err = proto.Unmarshal(controlMessageProtoBytes, controlMessage)

	if err != nil {
		return nil, err
	}
	return controlMessage, nil
}

func getControlHeader(request *e2appducontents.RiccontrolRequest) (*e2smrcies.E2SmRcControlHeader, error) {
	var rcServiceModel e2smrc.RCServiceModel
	var controlHeaderAsnBytes []byte
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRiccontrolHeader) {
			controlHeaderAsnBytes = v.GetValue().GetRiccontrolHeader().GetValue()
			break
		}
	}
	controlHeaderProtoBytes, err := rcServiceModel.ControlHeaderASN1toProto(controlHeaderAsnBytes)
	if err != nil {
		return nil, err
	}
	controlHeader := &e2smrcies.E2SmRcControlHeader{}
	err = proto.Unmarshal(controlHeaderProtoBytes, controlHeader)
	if err != nil {
		return nil, err
	}

	return controlHeader, nil
}

func createRANParametersInsertStyle3List() ([]*e2smrcies.InsertIndicationRanparameterItem, error) {
	// RAN Parameters for Insert Style 3
	insertRANParametersStyle3List := make([]*e2smrcies.InsertIndicationRanparameterItem, 0)
	ranParameter1, err := pdubuilder.CreateInsertIndicationRanparameterItem(1, "Target Primary Cell ID")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter1)

	ranParameter2, err := pdubuilder.CreateInsertIndicationRanparameterItem(2, "Target Cell")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter2)

	ranParameter3, err := pdubuilder.CreateInsertIndicationRanparameterItem(3, "NR Cell")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter3)

	ranParameter4, err := pdubuilder.CreateInsertIndicationRanparameterItem(4, "NR CGI")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter4)

	ranParameter5, err := pdubuilder.CreateInsertIndicationRanparameterItem(7, "List of PDU sessions for handover")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter5)

	ranParameter6, err := pdubuilder.CreateInsertIndicationRanparameterItem(13, "List of DRBs for handover")
	if err != nil {
		return nil, err
	}
	insertRANParametersStyle3List = append(insertRANParametersStyle3List, ranParameter6)
	return insertRANParametersStyle3List, nil

}

func createRANParametersReportStyle3List() ([]*e2smrcies.ReportRanparameterItem, error) {
	// RAN Parameters for Report Style 3
	reportParametersStyle3List := make([]*e2smrcies.ReportRanparameterItem, 0)

	return reportParametersStyle3List, nil
}

func createRANParametersReportStyle2List() ([]*e2smrcies.ReportRanparameterItem, error) {
	// RAN Parameters for Report Style 2
	reportParametersStyle2List := make([]*e2smrcies.ReportRanparameterItem, 0)
	ranParameter1, err := pdubuilder.CreateReportRanparameterItem(1, "Current UE ID")
	if err != nil {
		return nil, err
	}

	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter1)

	ranParameter2, err := pdubuilder.CreateReportRanparameterItem(21001, "S-NSSAI")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter2)

	ranParameter3, err := pdubuilder.CreateReportRanparameterItem(21002, "SST")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter3)

	ranParameter4, err := pdubuilder.CreateReportRanparameterItem(21003, "SD")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter4)

	ranParameter5, err := pdubuilder.CreateReportRanparameterItem(27108, "Best Neighboring Cell")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter5)

	ranParameter6, err := pdubuilder.CreateReportRanparameterItem(21528, "List of Neighbor cells")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter6)

	ranParameter7, err := pdubuilder.CreateReportRanparameterItem(10102, "Cell Results")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter7)

	ranParameter8, err := pdubuilder.CreateReportRanparameterItem(10103, "SSB Results")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter8)

	ranParameter9, err := pdubuilder.CreateReportRanparameterItem(12501, "RSRP")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter9)

	ranParameter10, err := pdubuilder.CreateReportRanparameterItem(12502, "RSRQ")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter10)

	ranParameter11, err := pdubuilder.CreateReportRanparameterItem(12503, "SINR")
	if err != nil {
		return nil, err
	}
	reportParametersStyle2List = append(reportParametersStyle2List, ranParameter11)
	return reportParametersStyle2List, nil
}

func (c *Client) getCellPCI(ctx context.Context, ncgi ransimtypes.NCGI) (int32, error) {
	cell, err := c.ServiceModel.CellStore.Get(ctx, ncgi)
	if err != nil {
		return 0, err
	}

	return int32(cell.PCI), nil
}

func (c *Client) getARFCN(ctx context.Context, ncgi ransimtypes.NCGI) (int32, error) {
	cell, err := c.ServiceModel.CellStore.Get(ctx, ncgi)
	if err != nil {
		return 0, err
	}

	return int32(cell.Earfcn), nil
}

func (c *Client) createRICIndicationFormat3(ctx context.Context, cells []ransimtypes.NCGI, subscription *subutils.Subscription, e2NodeInfoChangeID int32) (*e2appducontents.Ricindication, error) {
	headerFormat1 := format1.NewIndicationHeader(format1.WithEventConditionID(e2NodeInfoChangeID))
	indicationHeaderAsn1Bytes, err := headerFormat1.ToAsn1Bytes()
	if err != nil {
		return nil, err
	}
	plmnID := c.getPlmnID()

	cellInfoList := make([]*e2smrcies.E2SmRcIndicationMessageFormat3Item, 0)
	for _, ncgi := range cells {
		cell, err := c.ServiceModel.CellStore.Get(ctx, ncgi)
		if err != nil {
			return nil, err
		}

		nci := ransimtypes.GetNCI(ncgi)
		nrCGI, err := pdubuilder.CreateNrCgi(plmnID.Value().ToBytes(), &asn1.BitString{
			Value: utils.Uint64ToBitString(uint64(nci), 36),
			Len:   36,
		})
		if err != nil {
			return nil, err
		}

		cgi, err := pdubuilder.CreateCgiNRCgi(nrCGI)
		if err != nil {
			return nil, err
		}

		pci, err := c.getCellPCI(ctx, ncgi)
		if err != nil {
			return nil, err
		}
		cellPCI, err := pdubuilder.CreateServingCellPciNR(pci)
		if err != nil {
			return nil, err
		}

		earfcn, err := c.getARFCN(ctx, ncgi)
		if err != nil {
			return nil, err
		}
		cellArfcn, err := pdubuilder.CreateServingCellArfcnNR(earfcn)
		if err != nil {
			return nil, err
		}

		neighborCellList := make([]*e2smrcies.NeighborCellItem, 0)
		for _, neighborNCGI := range cell.Neighbors {
			neighborNci := ransimtypes.GetNCI(neighborNCGI)
			neighborNrCGI, err := pdubuilder.CreateNrCgi(plmnID.Value().ToBytes(), &asn1.BitString{
				Value: utils.Uint64ToBitString(uint64(neighborNci), 36),
				Len:   36,
			})
			if err != nil {
				return nil, err
			}

			neighborPci, err := c.getCellPCI(ctx, neighborNCGI)
			if err != nil {
				return nil, err
			}

			neighborEarfcn, err := c.getARFCN(ctx, neighborNCGI)
			if err != nil {
				return nil, err
			}

			neighborNrArfcn, err := pdubuilder.CreateNrArfcn(neighborEarfcn)
			if err != nil {
				return nil, err
			}

			nrFrequencyBandList := &e2smcommonies.NrfrequencyBandList{
				Value: make([]*e2smcommonies.NrfrequencyBandItem, 0),
			}

			supportedSulbandList := make([]*e2smcommonies.SupportedSulfreqBandItem, 0)
			frequencyBandItem, err := pdubuilder.CreateNrfrequencyBandItem(1, &e2smcommonies.SupportedSulbandList{
				Value: supportedSulbandList,
			})
			if err != nil {
				return nil, err
			}

			nrFrequencyBandList.Value = append(nrFrequencyBandList.Value, frequencyBandItem)

			nrFrequencyInfo, err := pdubuilder.CreateNrfrequencyInfo(neighborNrArfcn, nrFrequencyBandList)
			if err != nil {
				return nil, err
			}
			nrFrequencyInfo.SetFrequencyShift7P5Khz(pdubuilder.CreateNrfrequencyShift7P5KhzTrue())

			neighborCellItem, err := pdubuilder.CreateNeighborCellItemRanTypeChoiceNr(neighborNrCGI, neighborPci, []byte{0xFF, 0xFF, 0xFF}, pdubuilder.CreateNRModeInfoFDD(),
				nrFrequencyInfo, pdubuilder.CreateX2XNEstablishedTrue(), pdubuilder.CreateHOValidatedTrue(), 1)
			if err != nil {
				return nil, err
			}
			neighborCellList = append(neighborCellList, neighborCellItem)

		}

		neighborRelationTable, err := pdubuilder.CreateNeighborRelationInfo(cellPCI, cellArfcn, &e2smrcies.NeighborCellList{
			Value: neighborCellList,
		})
		if err != nil {
			return nil, err
		}

		item, err := pdubuilder.CreateE2SmRcIndicationMessageFormat3Item(cgi)
		if err != nil {
			return nil, err
		}
		item.SetNeighborRelationTable(neighborRelationTable).SetCellDeleted(false)
		cellInfoList = append(cellInfoList, item)

	}

	messageFormat3 := format3.NewIndicationMessage(format3.WithMessageItems(cellInfoList))
	indicationMessageAsn1Bytes, err := messageFormat3.ToAsn1Bytes()
	if err != nil {
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
		return nil, err
	}
	return ricIndication, nil
}

func (c *Client) getPlmnID() ransimtypes.Uint24 {
	plmnIDUint24 := ransimtypes.Uint24{}
	plmnIDUint24.Set(uint32(c.ServiceModel.Model.PlmnID))
	return plmnIDUint24
}

func float_decoder(data int32) float32 {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)
	bits := binary.LittleEndian.Uint32(buf.Bytes())
	res := math.Float32frombits(bits)
	log.Debugf("data : %v", data)
	log.Debugf("res : %v", res)
	return res
}

// List of RAN Parameters
// > RAN Parameter ID: 1
// > Target Primary Cell ID structure
// >> RAN Parameter ID: 2
// >> Choice Target Cell structure
//
//	>>> [Choice 1] RAN Parameter ID: 3
//	>>> [Choice 1] NR Cell structure
//	 >>>> RAN Parameter ID: 4
//	 >>>> NR CGI Element with key flag false
//	>>> [Choice 2] RAN Parameter ID: 5
//	>>> [Choice 2] E-UTRA Cell structure
//	 >>>> RAN Parameter ID: 6
//	 >>>> E-UTRA CGI Element with key flag false
func (c *Client) createRICIndicationFormat5Header2(ctx context.Context, subscription *subutils.Subscription, ueID *e2smcommonies.Ueid, ricInsertStyleType int32, insertIndicationID int32, targetNCGI ransimtypes.NCGI) (*e2appducontents.Ricindication, error) {
	headerFormat2 := format2.NewIndicationHeader(format2.WithUEID(ueID),
		format2.WithRICInsertStyleType(ricInsertStyleType),
		format2.WithInsertIndicationID(insertIndicationID))
	indicationHeaderAsn1Byte, err := headerFormat2.ToAsn1Bytes()
	if err != nil {
		return nil, err
	}

	nrcgiRanParamValuePrint, err := pdubuilder.CreateRanparameterValuePrintableString(fmt.Sprintf("%x", targetNCGI))
	if err != nil {
		return nil, err
	}
	nrCgiRanParamValue, err := pdubuilder.CreateRanparameterValueTypeChoiceElementFalse(nrcgiRanParamValuePrint)
	if err != nil {
		return nil, err
	}
	nrCgiRanParamValueItem, err := pdubuilder.CreateRanparameterStructureItem(NRCGIRANParameterID, nrCgiRanParamValue)
	if err != nil {
		return nil, err
	}
	nrCellRanParamValue, err := pdubuilder.CreateRanParameterStructure([]*e2smrcies.RanparameterStructureItem{nrCgiRanParamValueItem})
	if err != nil {
		return nil, err
	}
	nrCellRanParamValueType, err := pdubuilder.CreateRanparameterValueTypeChoiceStructure(nrCellRanParamValue)
	if err != nil {
		return nil, err
	}
	nrCellRanParamValueItem, err := pdubuilder.CreateRanparameterStructureItem(NRCellRANParameterID, nrCellRanParamValueType)
	if err != nil {
		return nil, err
	}
	targetCellRanParamValue, err := pdubuilder.CreateRanParameterStructure([]*e2smrcies.RanparameterStructureItem{nrCellRanParamValueItem})
	if err != nil {
		return nil, err
	}
	targetCellRanParamValueType, err := pdubuilder.CreateRanparameterValueTypeChoiceStructure(targetCellRanParamValue)
	if err != nil {
		return nil, err
	}
	targetCellRanParamValueItem, err := pdubuilder.CreateRanparameterStructureItem(TargetCellRANParameterID, targetCellRanParamValueType)
	if err != nil {
		return nil, err
	}
	targetPrimaryCellIDRanParamValue, err := pdubuilder.CreateRanParameterStructure([]*e2smrcies.RanparameterStructureItem{targetCellRanParamValueItem})
	if err != nil {
		return nil, err
	}
	targetPrimaryCellIDRanParamValueType, err := pdubuilder.CreateRanparameterValueTypeChoiceStructure(targetPrimaryCellIDRanParamValue)
	if err != nil {
		return nil, err
	}
	targetPrimaryCellIDRanParamValueItem, err := pdubuilder.CreateE2SmRcIndicationMessageFormat5Item(TargetPrimaryCellIDRANParameterID, targetPrimaryCellIDRanParamValueType)
	if err != nil {
		return nil, err
	}

	rpl := []*e2smrcies.E2SmRcIndicationMessageFormat5Item{targetPrimaryCellIDRanParamValueItem}
	messageFormat5 := format5.NewIndicationMessage(format5.WithMessageItems(rpl))

	indicationMessageAsn1Byte, err := messageFormat5.ToAsn1Bytes()
	if err != nil {
		return nil, err
	}

	indication := indicationutils.NewIndication(
		indicationutils.WithRicInstanceID(subscription.GetRicInstanceID()),
		indicationutils.WithRanFuncID(subscription.GetRanFuncID()),
		indicationutils.WithRequestID(subscription.GetReqID()),
		indicationutils.WithIndicationHeader(indicationHeaderAsn1Byte),
		indicationutils.WithIndicationMessage(indicationMessageAsn1Byte))
	// TODO add indicationutils.WithRicCallProcessID([]byte(0)))

	ricIndication, err := indication.Build()
	if err != nil {
		return nil, err
	}

	return ricIndication, nil
}

func (c *Client) handleControlMessage(ctx context.Context, controlHeader *e2smrcies.E2SmRcControlHeader, controlMessage *e2smrcies.E2SmRcControlMessage) error {
	headerFormat1 := controlHeader.GetRicControlHeaderFormats().GetControlHeaderFormat1()
	headerFormat2 := controlHeader.GetRicControlHeaderFormats().GetControlHeaderFormat2()
	messageFormat1 := controlMessage.GetRicControlMessageFormats().GetControlMessageFormat1()

	if headerFormat1 != nil {
		// handler for header format 1
		// handler for message format 1
		if messageFormat1 != nil {
			if headerFormat1.GetRicStyleType().Value == controlStyleType200 && headerFormat1.GetRicControlActionId().Value == controlActionID1 {
				// for PCI change
				err := c.checkAndSetPCI(ctx, messageFormat1)
				if err != nil {
					return err
				}
			} else if headerFormat1.GetRicStyleType().Value == controlStyleType3 && headerFormat1.GetRicControlActionId().Value == controlActionID1 {
				// for MHO
				err := c.runHandover(ctx, headerFormat1, messageFormat1)
				if err != nil {
					return err
				}
			}
		}
	} else if headerFormat2 != nil {
		// TODO write handler for header format2
		log.Error("header format 2 handler is not implemented yet")
	}

	return nil
}

func (c *Client) runHandover(ctx context.Context, controlHeader *e2smrcies.E2SmRcControlHeaderFormat1, controlMessage *e2smrcies.E2SmRcControlMessageFormat1) error {
	ueID := controlHeader.GetUeId()
	ue, err := c.ServiceModel.UEs.GetWithGNbUeID(ctx, ueID.GetGNbUeid())
	if err != nil {
		log.Error(err)
		return err
	}

	for _, ranParameter := range controlMessage.GetRanPList() {
		ranParameterID := ranParameter.GetRanParameterId().Value
		if ranParameterID == TargetPrimaryCellIDRANParameterID {
			targetPrimaryCellIDString := ranParameter.GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()[0].
				GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()[0].
				GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()[0].
				GetRanParameterValueType().GetRanPChoiceElementFalse().GetRanParameterValue().GetValuePrintableString()
			targetPrimaryCellID, err := strconv.ParseUint(targetPrimaryCellIDString, 16, 64)
			if err != nil {
				log.Error(err)
				return err
			}
			ncgi := ransimtypes.NCGI(targetPrimaryCellID)
			tCell := &model.UECell{
				ID:   ransimtypes.GnbID(ncgi),
				NCGI: ncgi,
			}
			c.mobilityDriver.Handover(ctx, ue.IMSI, tCell)
		}
	}
	return nil
}

func (c *Client) extractNCGIFromPrintableNCGI(pa *e2smrcies.RicPolicyAction) (ransimtypes.NCGI, error) {
	for _, rp := range pa.GetRanParametersList() {
		if rp.GetRanParameterId().Value == 1 {
			targetPrimaryCellIDString := rp.GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()[0].
				GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()[0].
				GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()[0].
				GetRanParameterValueType().GetRanPChoiceElementFalse().GetRanParameterValue().GetValuePrintableString()
			log.Debugf("targetPrimaryCellIDString %+v", targetPrimaryCellIDString)
			log.Debugf("RAN Parameter %+v", rp)
			targetPrimaryCellID, err := strconv.ParseUint(targetPrimaryCellIDString, 16, 64)
			if err != nil {
				return 0, err
			}

			return ransimtypes.NCGI(targetPrimaryCellID), nil
		}
	}
	return 0, errors.NewNotFound("RanParameter 1 for target primary cell ID not found")
}

func (c *Client) extractOcn(pa *e2smrcies.RicPolicyAction) (int, error) {
	for _, rp := range pa.GetRanParametersList() {
		if rp.GetRanParameterId().Value == 10201 {
			log.Debugf("extracted Ocn %+v", rp.GetRanParameterValueType().GetRanPChoiceElementFalse().GetRanParameterValue().GetValueInt())
			return int(rp.GetRanParameterValueType().GetRanPChoiceElementFalse().GetRanParameterValue().GetValueInt()), nil
		}
	}
	return 0, errors.NewNotFound("RanParameter 10201 for Ocn not found")
}

// checkAndSetPCI check if the control header and message including the required info for changing the PCI value for a specific cell
func (c *Client) checkAndSetPCI(ctx context.Context, controlMessage *e2smrcies.E2SmRcControlMessageFormat1) error {
	var pciValue int64
	for _, ranParameter := range controlMessage.GetRanPList() {
		var ncgi ransimtypes.NCGI
		// Extracts NR PCI ran parameter
		ranParameterID := ranParameter.GetRanParameterId().Value
		if ranParameterID == PCIRANParameterID {
			ranParameterValue := ranParameter.GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()[0].GetRanParameterValueType().GetRanPChoiceElementFalse()
			if ranParameterValue != nil {
				pciValue = ranParameterValue.GetRanParameterValue().GetValueInt()
			} else {
				return errors.NewInvalid("PCI ran parameter is not set")
			}
		}
		// Extracts NCGI ran parameter
		if ranParameterID == NCGIRANParameterID {
			ncgiStruct := ranParameter.GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()[0].GetRanParameterValueType().GetRanPChoiceStructure()
			if ncgiStruct != nil {
				ncgiFields := ncgiStruct.GetRanParameterStructure().GetSequenceOfRanParameters()
				if len(ncgiFields) == 2 {
					plmnIDField := ncgiFields[0]
					var plmnID ransimtypes.PlmnID
					var nci ransimtypes.NCI
					if plmnIDField != nil {
						plmnIDBitString := plmnIDField.GetRanParameterValueType().GetRanPChoiceElementFalse().GetRanParameterValue().GetValueOctS()
						plmnID = ransimtypes.PlmnID(ransimtypes.Uint24ToUint32(plmnIDBitString))

					} else {
						return errors.NewInvalid("plmn ID ran parameter is not set")
					}
					nciField := ncgiFields[1]
					if nciField != nil {
						nciBitString := nciField.GetRanParameterValueType().GetRanPChoiceElementFalse().GetRanParameterValue().GetValueBitS()
						nci = ransimtypes.NCI(utils.BitStringToUint64(nciBitString.GetValue(), int(nciBitString.GetLen())))
					} else {
						return errors.NewInvalid("NCI ran parameter is not set")
					}
					ncgi = ransimtypes.ToNCGI(plmnID, nci)
					cell, err := c.ServiceModel.CellStore.Get(ctx, ncgi)
					if err != nil {
						return err
					}
					cell.PCI = uint32(pciValue)
					err = c.ServiceModel.CellStore.Update(ctx, cell)
					if err != nil {
						return err
					}
				}
			} else {
				return errors.NewInvalid("NCGI ran parameter is not set")
			}
		}
		if ranParameterID == NSRANParameterID {
			var control_values []float32
			ranParameter := ranParameter.GetRanParameterValueType().GetRanPChoiceStructure().GetRanParameterStructure().GetSequenceOfRanParameters()
			if ranParameter != nil {
				for index := 0; index < len(ranParameter); index++ {
					control_value := int32(ranParameter[index].GetRanParameterValueType().GetRanPChoiceElementFalse().GetRanParameterValue().GetValueInt())
					convert_control_value := float_decoder(control_value)
					control_values = append(control_values, convert_control_value)
				}
			} else {
				return errors.NewInvalid("Can not get control values")
			}
			log.Infof("control values : %v", control_values)
		}
	}
	return nil
}
