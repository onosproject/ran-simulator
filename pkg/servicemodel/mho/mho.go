// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mho

import (
	"context"
	"encoding/binary"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
	"github.com/onosproject/rrm-son-lib/pkg/model/id"
	"strconv"
	"time"

	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/pdubuilder"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/v1/e2sm-mho"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/mobility"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	e2apIndicationUtils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/indication"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"
	indHdr "github.com/onosproject/ran-simulator/pkg/utils/e2sm/mho/indication/header"
	indMsgFmt1 "github.com/onosproject/ran-simulator/pkg/utils/e2sm/mho/indication/message_format1"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/mho/ranfundesc"
	"google.golang.org/protobuf/proto"
)

var log = logging.GetLogger("sm", "mho")

const (
	fileFormatVersion1 string = "version1"
	//senderName         string = "RAN Simulator"
	//senderType         string = ""
	//vendorName         string = "ONF"
)

// Client mho service model client
type Client struct {
	ServiceModel *registry.ServiceModel
	subscription *subutils.Subscription
	context      context.Context
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model,
	modelPluginRegistry modelplugins.ModelRegistry,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store,
	ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store,
	measChan chan device.UE) (registry.ServiceModel, error) {
	modelName := e2smtypes.ShortName(modelFullName)
	mhoSm := registry.ServiceModel{
		RanFunctionID:       registry.Mho,
		ModelName:           modelName,
		Revision:            1,
		OID:                 modelOID,
		Version:             version,
		ModelPluginRegistry: modelPluginRegistry,
		Node:                node,
		Model:               model,
		Subscriptions:       subStore,
		Nodes:               nodeStore,
		UEs:                 ueStore,
		CellStore:           cellStore,
		MetricStore:         metricStore,
		MeasChan:            measChan,
	}

	mhoClient := &Client{
		ServiceModel: &mhoSm,
	}

	mhoSm.Client = mhoClient

	var ranFunctionShortName = modelFullName
	var ranFunctionE2SmOid = modelOID
	var ranFunctionDescription = "MHO"
	var ranFunctionInstance int32 = 3
	var ricEventStyleType int32 = 1
	var ricEventStyleName = "Periodic and On Change Report"
	var ricEventFormatType int32 = 1
	var ricReportStyleType int32 = 1
	var ricReportStyleName = "PCI and NRT update for eNB"
	var ricIndicationHeaderFormatType int32 = 1
	var ricIndicationMessageFormatType int32 = 1

	ricEventTriggerStyleList := make([]*e2sm_mho.RicEventTriggerStyleList, 0)
	ricEventTriggerItem1 := pdubuilder.CreateRicEventTriggerStyleItem(ricEventStyleType, ricEventStyleName, ricEventFormatType)
	ricEventTriggerStyleList = append(ricEventTriggerStyleList, ricEventTriggerItem1)

	ricReportStyleList := make([]*e2sm_mho.RicReportStyleList, 0)
	ricReportStyleItem1 := pdubuilder.CreateRicReportStyleItem(ricReportStyleType, ricReportStyleName, ricIndicationHeaderFormatType,
		ricIndicationMessageFormatType)
	ricReportStyleList = append(ricReportStyleList, ricReportStyleItem1)

	ranFuncDescPdu, err := ranfundesc.NewRANFunctionDescription(
		ranfundesc.WithRANFunctionDescription(ranFunctionDescription),
		ranfundesc.WithRANFunctionInstance(ranFunctionInstance),
		ranfundesc.WithRANFunctionShortName(ranFunctionShortName),
		ranfundesc.WithRANFunctionE2SmOID(ranFunctionE2SmOid),
		ranfundesc.WithRICEventTriggerStyleList(ricEventTriggerStyleList),
		ranfundesc.WithRICReportStyleList(ricReportStyleList)).
		Build()

	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	protoBytes, err := proto.Marshal(ranFuncDescPdu)
	log.Debug("Proto bytes of MHO service model Ran Function Description:", protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}
	mhoModelPlugin, err := modelPluginRegistry.GetPlugin(modelOID)
	if mhoModelPlugin == nil {
		log.Debug("model plugin names:", modelPluginRegistry.GetPlugins())
		return registry.ServiceModel{}, errors.New(errors.Invalid, "model plugin is nil: %v", err)
	}
	ranFuncDescBytes, err := mhoModelPlugin.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	mhoSm.Description = ranFuncDescBytes
	return mhoSm, nil
}

// RICSubscription implements subscription handler for MHO service model
func (sm *Client) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Infof("Ric Subscription Request is received for service model %v and e2 node with ID:%d", sm.ServiceModel.ModelName, sm.ServiceModel.Node.GnbID)
	log.Debugf("MHO subscription, request: %v", request)
	var ricActionsAccepted []*e2aptypes.RicActionID
	ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
	actionList := subutils.GetRicActionToBeSetupList(request)
	reqID := subutils.GetRequesterID(request)
	ranFuncID := subutils.GetRanFunctionID(request)
	ricInstanceID := subutils.GetRicInstanceID(request)

	log.Debugf("MHO subscription, action list: %v", actionList)
	log.Debugf("MHO subscription, requester id: %v", reqID)
	log.Debugf("MHO subscription, ran func id: %v", ranFuncID)
	log.Debugf("MHO subscription, ric instance id: %v", ricInstanceID)

	for _, action := range actionList {
		log.Debugf("MHO subscription action: %v", action)
		actionID := e2aptypes.RicActionID(action.Value.RicActionId.Value)
		actionType := action.Value.RicActionType
		// mho service model supports report and insert action and should be added to the
		// list of accepted actions
		if actionType == e2apies.RicactionType_RICACTION_TYPE_REPORT ||
			actionType == e2apies.RicactionType_RICACTION_TYPE_INSERT {
			ricActionsAccepted = append(ricActionsAccepted, &actionID)
		}
		// mho service model does not support POLICY actions and
		// should be added into the list of not admitted actions
		if actionType == e2apies.RicactionType_RICACTION_TYPE_POLICY {
			cause := &e2apies.Cause{
				Cause: &e2apies.Cause_RicRequest{
					RicRequest: e2apies.CauseRic_CAUSE_RIC_ACTION_NOT_SUPPORTED,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
	}
	sm.subscription = subutils.NewSubscription(
		subutils.WithRequestID(reqID),
		subutils.WithRanFuncID(ranFuncID),
		subutils.WithRicInstanceID(ricInstanceID),
		subutils.WithActionsAccepted(ricActionsAccepted),
		subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		log.Warn("MHO subscription failed: no actions are accepted")
		subscriptionFailure, err := sm.subscription.BuildSubscriptionFailure()
		if err != nil {
			log.Error(err)
			return nil, nil, err
		}
		log.Warnf("MHO subscription failed, no actions accepted: %v", actionList)
		return nil, subscriptionFailure, nil
	}

	response, err = sm.subscription.BuildSubscriptionResponse()
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	eventTriggerType, err := sm.getEventTriggerType(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	sm.context = ctx

	log.Debugf("MHO subscription event trigger type: %v", eventTriggerType)
	switch eventTriggerType {
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_PERIODIC:
		log.Debug("Received periodic report subscription request")
		interval, err := sm.getReportPeriod(request)
		if err != nil {
			log.Error(err)
			return nil, nil, err
		}
		go sm.reportPeriodicIndication(interval)
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_UPON_RCV_MEAS_REPORT:
		log.Debug("Received MHO_TRIGGER_TYPE_UPON_RCV_MEAS_REPORT subscription request")
		go sm.processEventA3MeasReport()
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_UPON_CHANGE_RRC_STATUS:
		log.Debug("Received MHO_TRIGGER_TYPE_UPON_CHANGE_RRC_STATUS subscription request")
	default:
		log.Errorf("MHO subscription failed, invalid event trigger type: %v", eventTriggerType)
	}

	log.Debug("MHO subscription response: %v", response)
	return response, nil, nil
}

// RICSubscriptionDelete implements subscription delete handler for MHO service model
func (sm *Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Infof("Ric subscription delete request is received for service model %v and e2 node with ID: %d", sm.ServiceModel.ModelName, sm.ServiceModel.Node.GnbID)
	reqID := subdeleteutils.GetRequesterID(request)
	ranFuncID := subdeleteutils.GetRanFunctionID(request)
	ricInstanceID := subdeleteutils.GetRicInstanceID(request)
	subID := subscriptions.NewID(ricInstanceID, reqID, ranFuncID)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerAsnBytes := sub.Details.RicEventTriggerDefinition.Value
	mhoModelPlugin, _ := sm.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(sm.ServiceModel.OID))
	eventTriggerProtoBytes, err := mhoModelPlugin.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerDefinition := &e2sm_mho.E2SmMhoEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerType := eventTriggerDefinition.GetEventDefinitionFormat1().TriggerType
	subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
		subdeleteutils.WithRequestID(reqID),
		subdeleteutils.WithRanFuncID(ranFuncID),
		subdeleteutils.WithRicInstanceID(ricInstanceID))
	response, err = subscriptionDelete.BuildSubscriptionDeleteResponse()
	if err != nil {
		return nil, nil, err
	}

	switch eventTriggerType {
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_PERIODIC:
		log.Debug("Stopping the periodic report subscription")
		sub.Ticker.Stop()
	}

	return response, nil, nil
}

// RICControl implements control handler for MHO service model
func (sm *Client) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	log.Infof("Control Request is received for service model %v and e2 node ID: %d", sm.ServiceModel.ModelName, sm.ServiceModel.Node.GnbID)

	controlHeader, err := sm.getControlHeader(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	// TODO - check MHO command
	log.Debugf("MHO control header: %v", controlHeader)

	controlMessage, err := sm.getControlMessage(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	log.Debugf("MHO control message: %v", controlMessage)

	imsi, err := strconv.Atoi(controlMessage.GetControlMessageFormat1().GetUedId().GetValue())
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	ue, err := sm.ServiceModel.UEs.Get(ctx, types.IMSI(imsi))
	if err != nil {
		log.Errorf("UE: %v not found err: %v", ue, err)
		return nil, nil, err
	}
	log.Debugf("imsi: %v", imsi)

	plmnIDBytes := controlMessage.GetControlMessageFormat1().GetTargetCgi().GetNrCgi().GetPLmnIdentity().GetValue()
	plmnID := ransimtypes.Uint24ToUint32(plmnIDBytes)
	nci := controlMessage.GetControlMessageFormat1().GetTargetCgi().GetNrCgi().GetNRcellIdentity().GetValue().GetValue()
	log.Debugf("ECI is %d and PLMN ID is %d", nci, plmnID)
	tCellNcgi := ransimtypes.ToNCGI(ransimtypes.PlmnID(plmnID), ransimtypes.NCI(nci))

	tCell := &model.UECell{
		ID:   types.GnbID(tCellNcgi),
		NCGI: tCellNcgi,
	}

	mobility.DoHandover(ctx, types.IMSI(imsi), tCell, sm.ServiceModel.UEs, sm.ServiceModel.CellStore)

	return nil, nil, nil
}

func (sm *Client) reportPeriodicIndication(interval int32) {
	log.Debugf("Starting periodic report with interval %d ms", interval)
	subID := subscriptions.NewID(sm.subscription.GetRicInstanceID(), sm.subscription.GetReqID(), sm.subscription.GetRanFuncID())
	intervalDuration := time.Duration(interval)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return
	}
	sub.Ticker = time.NewTicker(intervalDuration * time.Millisecond)
	for {
		select {
		case <-sub.Ticker.C:
			log.Debug("Sending periodic indication report for subscription:", sub.ID)
			err = sm.sendRicIndication()
			if err != nil {
				log.Error("Failure sending indication message: ", err)
			}

		case <-sub.E2Channel.Context().Done():
			sub.Ticker.Stop()
			return
		}
	}
}

func (sm *Client) sendRicIndication() error {
	node := sm.ServiceModel.Node
	// Creates and sends an indication message for each cell in the node
	for _, ncgi := range node.Cells {
		log.Debugf("Send MHO indications for cell ncgi:%d", ncgi)
		for _, ue := range sm.ServiceModel.UEs.ListUEs(sm.context, ncgi) {
			log.Debugf("Send MHO indications for cell ncgi:%d, IMSI:%d", ncgi, ue.IMSI)
			err := sm.sendRicIndicationFormat1(ncgi, ue)
			if err != nil {
				log.Warn(err)
				continue
			}
		}
	}
	return nil
}

func (sm *Client) sendRicIndicationFormat1(ncgi ransimtypes.NCGI, ue *model.UE) error {
	subID := subscriptions.NewID(sm.subscription.GetRicInstanceID(), sm.subscription.GetReqID(), sm.subscription.GetRanFuncID())
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	indicationHeaderBytes, err := sm.createIndicationHeaderBytes(ncgi, fileFormatVersion1)
	if err != nil {
		return err
	}

	indicationMessageBytes, err := sm.createIndicationMsgFormat1(ue)
	if err != nil {
		return err
	}
	if indicationMessageBytes == nil {
		return nil
	}

	indication := e2apIndicationUtils.NewIndication(
		e2apIndicationUtils.WithRicInstanceID(sm.subscription.GetRicInstanceID()),
		e2apIndicationUtils.WithRanFuncID(sm.subscription.GetRanFuncID()),
		e2apIndicationUtils.WithRequestID(sm.subscription.GetReqID()),
		e2apIndicationUtils.WithIndicationHeader(indicationHeaderBytes),
		e2apIndicationUtils.WithIndicationMessage(indicationMessageBytes))

	ricIndication, err := indication.Build()
	if err != nil {
		return err
	}

	err = sub.E2Channel.RICIndication(sm.context, ricIndication)
	if err != nil {
		return err
	}

	return nil
}

func (sm *Client) createIndicationHeaderBytes(ncgi ransimtypes.NCGI, fileFormatVersion string) ([]byte, error) {

	cell, _ := sm.ServiceModel.CellStore.Get(sm.context, ncgi)
	plmnID := ransimtypes.NewUint24(uint32(sm.ServiceModel.Model.PlmnID))
	timestamp := make([]byte, 4)
	binary.BigEndian.PutUint32(timestamp, uint32(time.Now().Unix()))
	header := indHdr.NewIndicationHeader(
		indHdr.WithPlmnID(*plmnID),
		indHdr.WithNrcellIdentity(uint64(ransimtypes.GetECI(uint64(cell.NCGI)))))

	mhoModelPlugin, err := sm.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(sm.ServiceModel.OID))
	if err != nil {
		return nil, err
	}

	indicationHeaderAsn1Bytes, err := header.MhoToAsn1Bytes(mhoModelPlugin)
	if err != nil {
		return nil, err
	}

	return indicationHeaderAsn1Bytes, nil
}

func (sm *Client) createIndicationMsgFormat1(ue *model.UE) ([]byte, error) {
	log.Debugf("Create MHO Indication message ueID: %d", ue.IMSI)

	plmnID := ransimtypes.NewUint24(uint32(sm.ServiceModel.Model.PlmnID))
	measReport := make([]*e2sm_mho.E2SmMhoMeasurementReportItem, 0)

	if len(ue.Cells) == 0 {
		log.Infof("no neighbor cells found for ueID:%d", ue.IMSI)
		return nil, nil
	}

	// add serving cell to measReport
	measReport = append(measReport, &e2sm_mho.E2SmMhoMeasurementReportItem{
		Cgi: &e2sm_mho.CellGlobalId{
			CellGlobalId: &e2sm_mho.CellGlobalId_NrCgi{
				NrCgi: &e2sm_mho.Nrcgi{
					PLmnIdentity: &e2sm_mho.PlmnIdentity{
						Value: plmnID.ToBytes(),
					},
					NRcellIdentity: &e2sm_mho.NrcellIdentity{
						Value: &e2sm_mho.BitString{
							Value: uint64(ransimtypes.GetECI(uint64(ue.Cell.NCGI))),
							Len:   36,
						},
					},
				},
			},
		},
		Rsrp: &e2sm_mho.Rsrp{
			Value: int32(ue.Cell.Strength),
		},
	})

	for _, cell := range ue.Cells {
		measReport = append(measReport, &e2sm_mho.E2SmMhoMeasurementReportItem{
			Cgi: &e2sm_mho.CellGlobalId{
				CellGlobalId: &e2sm_mho.CellGlobalId_NrCgi{
					NrCgi: &e2sm_mho.Nrcgi{
						PLmnIdentity: &e2sm_mho.PlmnIdentity{
							Value: plmnID.ToBytes(),
						},
						NRcellIdentity: &e2sm_mho.NrcellIdentity{
							Value: &e2sm_mho.BitString{
								Value: uint64(ransimtypes.GetECI(uint64(cell.NCGI))),
								Len:   36,
							},
						},
					},
				},
			},
			Rsrp: &e2sm_mho.Rsrp{
				Value: int32(cell.Strength),
			},
		})
	}

	ueID := strconv.Itoa(int(ue.IMSI))

	log.Debugf("MHO measurement report for ueID %s: %v", ueID, measReport)

	indicationMessage := indMsgFmt1.NewIndicationMessage(
		indMsgFmt1.WithUeID(ueID),
		indMsgFmt1.WithMeasReport(measReport))

	log.Debugf("MHO indication message for ueID %s: %v", ueID, indicationMessage)

	mhoModelPlugin, err := sm.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(sm.ServiceModel.OID))
	if err != nil {
		return nil, err
	}
	indicationMessageBytes, err := indicationMessage.ToAsn1Bytes(mhoModelPlugin)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationMessageBytes, nil
}

func (sm *Client) processEventA3MeasReport() {
	log.Info("Start processing event a3 measurement report")
	for report := range sm.ServiceModel.MeasChan {
		log.Debugf("received event a3 measurement report: %v", report)
		log.Debugf("Send upon-rcv-meas-report indication for cell ecgi:%d, IMSI:%s",
			report.GetSCell().GetID().GetID().(id.ECGI), report.GetID().String())
		ecgi := report.GetSCell().GetID().GetID().(id.ECGI)
		imsi := report.GetID().GetID().(id.UEID).IMSI
		ue, err := sm.ServiceModel.UEs.Get(sm.context, types.IMSI(imsi))
		if err != nil {
			log.Warn(err)
			continue
		}
		err = sm.sendRicIndicationFormat1(ransimtypes.NCGI(ecgi), ue)
		if err != nil {
			log.Warn(err)
			continue
		}
	}
}
