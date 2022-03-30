// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

import (
	"context"

	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"
	e2smmhosm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/onosproject/rrm-son-lib/pkg/handover"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/pdubuilder"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/mobility"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	controlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/mho/ranfundesc"
	"google.golang.org/protobuf/proto"
)

var log = logging.GetLogger("sm", "mho")

// Mho represents the MHO service model
type Mho struct {
	ServiceModel   *registry.ServiceModel
	rrcUpdateChan  chan model.UE
	mobilityDriver mobility.Driver
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store,
	ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store,
	a3Chan chan handover.A3HandoverDecision, mobilityDriver mobility.Driver) (registry.ServiceModel, error) {
	modelName := e2smtypes.ShortName(modelFullName)
	mhoSm := registry.ServiceModel{
		RanFunctionID: registry.Mho,
		ModelName:     modelName,
		Revision:      1,
		OID:           modelOID,
		Version:       version,
		Node:          node,
		Model:         model,
		Subscriptions: subStore,
		Nodes:         nodeStore,
		UEs:           ueStore,
		CellStore:     cellStore,
		MetricStore:   metricStore,
		A3Chan:        a3Chan,
	}

	mho := &Mho{
		ServiceModel: &mhoSm,
	}

	mhoSm.Client = mho

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
	ricEventTriggerItem1, err := pdubuilder.CreateRicEventTriggerStyleItem(ricEventStyleType, ricEventStyleName, ricEventFormatType)
	if err != nil {
		return registry.ServiceModel{}, err
	}
	ricEventTriggerStyleList = append(ricEventTriggerStyleList, ricEventTriggerItem1)

	ricReportStyleList := make([]*e2sm_mho.RicReportStyleList, 0)
	ricReportStyleItem1, err := pdubuilder.CreateRicReportStyleItem(ricReportStyleType, ricReportStyleName, ricIndicationHeaderFormatType,
		ricIndicationMessageFormatType)
	if err != nil {
		return registry.ServiceModel{}, err
	}
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

	var mhosm e2smmhosm.MhoServiceModel
	ranFuncDescBytes, err := mhosm.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	mhoSm.Description = ranFuncDescBytes
	return mhoSm, nil
}

// E2ConnectionUpdate implements connection update handler
func (m *Mho) E2ConnectionUpdate(ctx context.Context, request *e2appducontents.E2ConnectionUpdate) (response *e2appducontents.E2ConnectionUpdateAcknowledge, failure *e2appducontents.E2ConnectionUpdateFailure, err error) {
	return nil, nil, errors.NewNotSupported("E2 connection update is not supported")
}

// RICSubscription implements subscription handler for MHO service model
func (m *Mho) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Infof("Ric Subscription Request is received for service model %v and e2 node with ID:%d", m.ServiceModel.ModelName, m.ServiceModel.Node.GnbID)
	log.Debugf("MHO subscription, request: %v", request)
	var ricActionsAccepted []*e2aptypes.RicActionID
	ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
	actionList := subutils.GetRicActionToBeSetupList(request)
	reqID, err := subutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := subutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := subutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("MHO subscription, action list: %v", actionList)
	log.Debugf("MHO subscription, requester id: %v", reqID)
	log.Debugf("MHO subscription, ran func id: %v", ranFuncID)
	log.Debugf("MHO subscription, ric instance id: %v", ricInstanceID)

	for _, action := range actionList {
		log.Debugf("MHO subscription action: %v", action)
		actionID := e2aptypes.RicActionID(action.GetValue().GetRicactionToBeSetupItem().GetRicActionId().GetValue())
		actionType := action.GetValue().GetRicactionToBeSetupItem().GetRicActionType()
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
					RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_ACTION_NOT_SUPPORTED,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
	}

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		log.Warn("no action is accepted")
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_ACTION_NOT_SUPPORTED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, subscriptionFailure, nil
	}

	eventTriggerType, err := m.getEventTriggerType(request)
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, subscriptionFailure, err
		}
		return nil, subscriptionFailure, nil
	}

	subscription := subutils.NewSubscription(
		subutils.WithRequestID(*reqID),
		subutils.WithRanFuncID(*ranFuncID),
		subutils.WithRicInstanceID(*ricInstanceID),
		subutils.WithActionsAccepted(ricActionsAccepted),
		subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))
	response, err = subscription.BuildSubscriptionResponse()
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, subscriptionFailure, err
		}
		return nil, subscriptionFailure, nil
	}

	log.Debugf("MHO subscription event trigger type: %v", eventTriggerType)
	switch eventTriggerType {
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_PERIODIC:
		log.Infof("Received periodic report subscription request")
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			interval, err := m.getReportPeriod(request)
			if err != nil {
				log.Error(err)
				return
			}
			m.reportPeriodicIndication(ctx, interval, subscription)
		}()
	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_UPON_RCV_MEAS_REPORT:
		log.Infof("Received MHO_TRIGGER_TYPE_UPON_RCV_MEAS_REPORT subscription request")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if m.mobilityDriver.GetHoLogic() == "local" {
			m.mobilityDriver.SetHoLogic("mho")
		}

		go m.processEventA3MeasReport(ctx, subscription)

	case e2sm_mho.MhoTriggerType_MHO_TRIGGER_TYPE_UPON_CHANGE_RRC_STATUS:
		log.Infof("Received MHO_TRIGGER_TYPE_UPON_CHANGE_RRC_STATUS subscription request")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		m.rrcUpdateChan = make(chan model.UE)
		go m.processRrcUpdate(ctx, subscription)
		m.mobilityDriver.AddRrcChan(m.rrcUpdateChan)

	default:
		log.Errorf("MHO subscription failed, invalid event trigger type: %v", eventTriggerType)
	}

	log.Debug("MHO subscription response: %v", response)
	return response, nil, nil
}

// RICSubscriptionDelete implements subscription delete handler for MHO service model
func (m *Mho) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Infof("Ric subscription delete request is received for service model %v and e2 node with ID: %d", m.ServiceModel.ModelName, m.ServiceModel.Node.GnbID)
	reqID, err := subdeleteutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := subdeleteutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := subdeleteutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}
	subID := subscriptions.NewID(*ricInstanceID, *reqID, *ranFuncID)
	sub, err := m.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerAsnBytes := sub.Details.RicEventTriggerDefinition.Value

	var mhoServiceModel e2smmhosm.MhoServiceModel
	eventTriggerProtoBytes, err := mhoServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerDefinition := &e2sm_mho.E2SmMhoEventTriggerDefinition{}
	err = proto.Unmarshal(eventTriggerProtoBytes, eventTriggerDefinition)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerType := eventTriggerDefinition.GetEventDefinitionFormats().GetEventDefinitionFormat1().TriggerType
	subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
		subdeleteutils.WithRequestID(*reqID),
		subdeleteutils.WithRanFuncID(*ranFuncID),
		subdeleteutils.WithRicInstanceID(*ricInstanceID))
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

	go func() {

		// ToDo - should be reconsidered (not locked on GNb and AmfNGap)
		imsi := controlMessage.GetControlMessageFormat1().GetUedId().GetGNbUeid().GetAmfUeNgapId().GetValue()

		plmnIDBytes := controlMessage.GetControlMessageFormat1().GetTargetCgi().GetNRCgi().GetPLmnidentity().GetValue()
		plmnID := ransimtypes.Uint24ToUint32(plmnIDBytes)
		nci := utils.NewNCellIDWithBytes(controlMessage.GetControlMessageFormat1().GetTargetCgi().GetNRCgi().GetNRcellIdentity().GetValue().GetValue())
		tCellNcgi := ransimtypes.ToNCGI(ransimtypes.PlmnID(plmnID), ransimtypes.NCI(nci.Uint64()))
		tCell := &model.UECell{
			ID:   types.GnbID(tCellNcgi),
			NCGI: tCellNcgi,
		}
		m.mobilityDriver.Handover(ctx, types.IMSI(imsi), tCell)
	}()

	reqID, err := controlutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := controlutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := controlutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}
	response, err = controlutils.NewControl(
		controlutils.WithRanFuncID(*ranFuncID),
		controlutils.WithRequestID(*reqID),
		controlutils.WithRicInstanceID(*ricInstanceID)).BuildControlAcknowledge()
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	return response, nil, nil
}
