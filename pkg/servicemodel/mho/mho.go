// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mho

import (
	"context"
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
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/mho/ranfundesc"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
	"google.golang.org/protobuf/proto"
	"strconv"
)

var log = logging.GetLogger("sm", "mho")

// Mho represents the MHO service model
type Mho struct {
	ServiceModel   *registry.ServiceModel
	subscription   *subutils.Subscription
	context        context.Context
	rrcUpdateChan  chan model.UE
	mobilityDriver mobility.Driver
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model,
	modelPluginRegistry modelplugins.ModelRegistry,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store,
	ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store,
	measChan chan device.UE, mobilityDriver mobility.Driver) (registry.ServiceModel, error) {
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

	mho := &Mho{
		ServiceModel: &mhoSm,
	}

	mhoSm.Client = mho

	mho.rrcUpdateChan = mobilityDriver.GetRrcCtrl().RrcUpdateChan
	mho.mobilityDriver = mobilityDriver

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
func (m *Mho) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Infof("Ric Subscription Request is received for service model %v and e2 node with ID:%d", m.ServiceModel.ModelName, m.ServiceModel.Node.GnbID)
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
	m.subscription = subutils.NewSubscription(
		subutils.WithRequestID(reqID),
		subutils.WithRanFuncID(ranFuncID),
		subutils.WithRicInstanceID(ricInstanceID),
		subutils.WithActionsAccepted(ricActionsAccepted),
		subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		log.Warn("MHO subscription failed: no actions are accepted")
		subscriptionFailure, err := m.subscription.BuildSubscriptionFailure()
		if err != nil {
			log.Error(err)
			return nil, nil, err
		}
		log.Warnf("MHO subscription failed, no actions accepted: %v", actionList)
		return nil, subscriptionFailure, nil
	}

	response, err = m.subscription.BuildSubscriptionResponse()
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	eventTriggerType, err := m.getEventTriggerType(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	m.context = ctx

	log.Debugf("MHO subscription event trigger type: %v", eventTriggerType)
	switch eventTriggerType {
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_PERIODIC:
		log.Infof("Received periodic report subscription request")
		interval, err := m.getReportPeriod(request)
		if err != nil {
			log.Error(err)
			return nil, nil, err
		}
		go m.reportPeriodicIndication(interval)
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_UPON_RCV_MEAS_REPORT:
		log.Infof("Received MHO_TRIGGER_TYPE_UPON_RCV_MEAS_REPORT subscription request")
		go m.processEventA3MeasReport()
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_UPON_CHANGE_RRC_STATUS:
		log.Infof("Received MHO_TRIGGER_TYPE_UPON_CHANGE_RRC_STATUS subscription request")
		go m.processRrcUpdate()
	default:
		log.Errorf("MHO subscription failed, invalid event trigger type: %v", eventTriggerType)
	}

	log.Debug("MHO subscription response: %v", response)
	return response, nil, nil
}

// RICSubscriptionDelete implements subscription delete handler for MHO service model
func (m *Mho) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Infof("Ric subscription delete request is received for service model %v and e2 node with ID: %d", m.ServiceModel.ModelName, m.ServiceModel.Node.GnbID)
	reqID := subdeleteutils.GetRequesterID(request)
	ranFuncID := subdeleteutils.GetRanFunctionID(request)
	ricInstanceID := subdeleteutils.GetRicInstanceID(request)
	subID := subscriptions.NewID(ricInstanceID, reqID, ranFuncID)
	sub, err := m.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerAsnBytes := sub.Details.RicEventTriggerDefinition.Value
	mhoModelPlugin, _ := m.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(m.ServiceModel.OID))
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
func (m *Mho) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	log.Infof("Control Request is received for service model %v and e2 node ID: %d", m.ServiceModel.ModelName, m.ServiceModel.Node.GnbID)

	controlHeader, err := m.getControlHeader(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	// TODO - check MHO command
	log.Debugf("MHO control header: %v", controlHeader)

	controlMessage, err := m.getControlMessage(request)
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
	ue, err := m.ServiceModel.UEs.Get(ctx, types.IMSI(imsi))
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

	m.mobilityDriver.Handover(ctx, types.IMSI(imsi), tCell)

	return nil, nil, nil
}
