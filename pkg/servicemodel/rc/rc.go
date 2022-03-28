// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package rc

import (
	"context"
	"time"

	"github.com/onosproject/ran-simulator/pkg/utils"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/ranfundesc"

	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/controloutcome"

	"github.com/onosproject/ran-simulator/pkg/store/metrics"

	"github.com/onosproject/ran-simulator/pkg/store/cells"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	e2smrcpresm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/servicemodel"
	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/v2/e2sm-rc-pre-v2-go"

	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"

	controlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"

	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"

	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"

	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/pdubuilder"
	"github.com/onosproject/ran-simulator/pkg/model"
	"google.golang.org/protobuf/proto"

	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger()

// Client rc service model client
type Client struct {
	ServiceModel *registry.ServiceModel
}

func (sm *Client) reportPeriodicIndication(ctx context.Context, interval uint32, subscription *subutils.Subscription) error {
	log.Debugf("Starting periodic report with interval %d ms", interval)
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	intervalDuration := time.Duration(interval)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}
	sub.Ticker = time.NewTicker(intervalDuration * time.Millisecond)
	for {
		select {
		case <-sub.Ticker.C:
			log.Debug("Sending periodic indication report for subscription:", sub.ID)
			err = sm.sendRicIndication(ctx, subscription)
			if err != nil {
				log.Error("creating indication message is failed", err)
				return err
			}

		case <-sub.E2Channel.Context().Done():
			sub.Ticker.Stop()
			return nil
		}
	}
}

func (sm *Client) sendRicIndication(ctx context.Context, subscription *subutils.Subscription) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	node := sm.ServiceModel.Node
	// Creates and sends an indication message for each cell in the node
	for _, ncgi := range node.Cells {
		ricIndication, err := sm.createRicIndication(ctx, ncgi, subscription)
		if err != nil {
			log.Error(err)
			return err
		}
		err = sub.E2Channel.RICIndication(ctx, ricIndication)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (sm *Client) reportIndicationOnChange(ctx context.Context, subscription *subutils.Subscription) error {
	log.Debugf("Sending report indication on change from node: %d", sm.ServiceModel.Node.GnbID)
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}
	cellEventCh := make(chan event.Event)
	nodeCells := sm.ServiceModel.Node.Cells
	err = sm.ServiceModel.CellStore.Watch(context.Background(), cellEventCh)
	if err != nil {
		return err
	}

	// Sends the first indication message
	err = sm.sendRicIndication(ctx, subscription)
	if err != nil {
		return err
	}

	for {
		select {
		case cellEvent := <-cellEventCh:
			log.Debug("Received cell event:", cellEvent)
			cellEventType := cellEvent.Type.(cells.CellEvent)
			if cellEventType == cells.UpdatedNeighbors || cellEventType == cells.Updated {
				cell := cellEvent.Value.(*model.Cell)
				for _, nodeCell := range nodeCells {
					if nodeCell == cell.NCGI {
						err = sm.sendRicIndication(ctx, subscription)
						if err != nil {
							log.Error(err)
						}
					}
				}
			}

		case <-sub.E2Channel.Context().Done():
			log.Debug("E2 channel context is done")
			return nil
		}
	}
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store,
	ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store) (registry.ServiceModel, error) {
	modelName := e2smtypes.ShortName(modelFullName)
	rcSm := registry.ServiceModel{
		RanFunctionID: registry.Rcpre2,
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
	}

	rcClient := &Client{
		ServiceModel: &rcSm,
	}

	rcSm.Client = rcClient

	var ranFunctionShortName = modelFullName
	var ranFunctionE2SmOid = modelOID
	var ranFunctionDescription = "RC PRE"
	var ranFunctionInstance int32 = 3
	var ricEventStyleType int32 = 1
	var ricEventStyleName = "Periodic and On Change Report"
	var ricEventFormatType int32 = 1
	var ricReportStyleType int32 = 1
	var ricReportStyleName = "PCI and NRT update for eNB"
	var ricIndicationHeaderFormatType int32 = 1
	var ricIndicationMessageFormatType int32 = 1

	ricEventTriggerStyleList := make([]*e2smrcpreies.RicEventTriggerStyleList, 0)
	ricEventTriggerItem1 := pdubuilder.CreateRicEventTriggerStyleItem(ricEventStyleType, ricEventStyleName, ricEventFormatType)
	ricEventTriggerStyleList = append(ricEventTriggerStyleList, ricEventTriggerItem1)

	ricReportStyleList := make([]*e2smrcpreies.RicReportStyleList, 0)
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
	log.Debug("Proto bytes of RC service model Ran Function Description:", protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	var rcsm e2smrcpresm.RcPreServiceModel
	ranFuncDescBytes, err := rcsm.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	rcSm.Description = ranFuncDescBytes
	return rcSm, nil
}

// E2ConnectionUpdate implements connection update handler
func (sm *Client) E2ConnectionUpdate(ctx context.Context, request *e2appducontents.E2ConnectionUpdate) (response *e2appducontents.E2ConnectionUpdateAcknowledge, failure *e2appducontents.E2ConnectionUpdateFailure, err error) {
	return nil, nil, errors.NewNotSupported("E2 connection update is not supported")
}

// RICControl implements control handler for RC service model
func (sm *Client) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	log.Infof("Control Request is received for service model %v and e2 node ID: %d", sm.ServiceModel.ModelName, sm.ServiceModel.Node.GnbID)
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

	controlMessage, err := sm.getControlMessage(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	controlHeader, err := sm.getControlHeader(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	log.Debugf("RC control header: %v", controlHeader)
	log.Debugf("RC control message: %v", controlMessage)

	plmnIDBytes := controlHeader.GetControlHeaderFormat1().Cgi.GetNrCgi().PLmnIdentity.Value
	nci := utils.BitStringToUint64(controlHeader.GetControlHeaderFormat1().GetCgi().GetNrCgi().NRcellIdentity.Value.GetValue(), 36)
	plmnID := ransimtypes.Uint24ToUint32(plmnIDBytes)
	log.Debugf("NCI is %d and PLMN ID is %d", nci, plmnID)

	ncgi := ransimtypes.ToNCGI(ransimtypes.PlmnID(plmnID), ransimtypes.NCI(nci))
	parameterName := controlMessage.GetControlMessage().ParameterType.RanParameterName.Value
	parameterID := controlMessage.GetControlMessage().ParameterType.RanParameterId.Value
	cell, err := sm.ServiceModel.CellStore.Get(ctx, ncgi)
	if err != nil {
		log.Debugf("Ran parameter for entity %d not found", ncgi)
		outcomeAsn1Bytes, err := controloutcome.NewControlOutcome(
			controloutcome.WithRanParameterID(parameterID)).
			ToAsn1Bytes()
		if err != nil {
			return nil, nil, err
		}
		failure, err = controlutils.NewControl(
			controlutils.WithRanFuncID(*ranFuncID),
			controlutils.WithRequestID(*reqID),
			controlutils.WithRicInstanceID(*ricInstanceID),
			controlutils.WithRicControlOutcome(outcomeAsn1Bytes)).BuildControlFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	var parameterValue interface{}
	switch controlMessage.GetControlMessage().ParameterType.RanParameterType {
	case e2smrcpreies.RanparameterType_RANPARAMETER_TYPE_INTEGER:
		parameterValue = controlMessage.GetControlMessage().GetParameterVal().GetValueInt()
	case e2smrcpreies.RanparameterType_RANPARAMETER_TYPE_ENUMERATED:
		parameterValue = controlMessage.GetControlMessage().GetParameterVal().GetValueEnum()
	case e2smrcpreies.RanparameterType_RANPARAMETER_TYPE_PRINTABLE_STRING:
		parameterValue = controlMessage.GetControlMessage().GetParameterVal().GetValuePrtS()
	}
	setPCI(parameterName, parameterValue, cell)
	sm.setHandoverOcn(ctx, parameterName, parameterValue, cell)

	err = sm.ServiceModel.CellStore.Update(ctx, cell)
	if err != nil {
		outcomeAsn1Bytes, err := controloutcome.NewControlOutcome(
			controloutcome.WithRanParameterID(parameterID)).
			ToAsn1Bytes()
		if err != nil {
			return nil, nil, err
		}
		failure, err = controlutils.NewControl(
			controlutils.WithRanFuncID(*ranFuncID),
			controlutils.WithRequestID(*reqID),
			controlutils.WithRicInstanceID(*ricInstanceID),
			controlutils.WithRicControlOutcome(outcomeAsn1Bytes)).BuildControlFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	outcomeAsn1Bytes, err := controloutcome.NewControlOutcome(
		controloutcome.WithRanParameterID(parameterID)).
		ToAsn1Bytes()
	if err != nil {
		return nil, nil, err
	}

	response, err = controlutils.NewControl(
		controlutils.WithRanFuncID(*ranFuncID),
		controlutils.WithRequestID(*reqID),
		controlutils.WithRicInstanceID(*ricInstanceID),
		controlutils.WithRicControlOutcome(outcomeAsn1Bytes)).BuildControlAcknowledge()
	if err != nil {
		return nil, nil, err
	}
	return response, nil, nil
}

// RICSubscription implements subscription handler for RC service model
func (sm *Client) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Infof("Ric Subscription Request is received for service model %v and e2 node with ID:%d", sm.ServiceModel.ModelName, sm.ServiceModel.Node.GnbID)
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

	for _, action := range actionList {
		actionID := e2aptypes.RicActionID(action.GetValue().GetRicactionToBeSetupItem().GetRicActionId().GetValue())
		actionType := action.GetValue().GetRicactionToBeSetupItem().GetRicActionType()
		// rc service model supports report and insert action and should be added to the
		// list of accepted actions
		if actionType == e2apies.RicactionType_RICACTION_TYPE_REPORT ||
			actionType == e2apies.RicactionType_RICACTION_TYPE_INSERT {
			ricActionsAccepted = append(ricActionsAccepted, &actionID)
		}
		// rc service model does not support POLICY actions and
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

	eventTriggerType, err := sm.getEventTriggerType(request)
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
			return nil, subscriptionFailure, nil
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
			return nil, subscriptionFailure, nil
		}
		return nil, subscriptionFailure, nil
	}

	switch eventTriggerType {
	case e2smrcpreies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_UPON_CHANGE:
		log.Debug("Received on change report subscription request")
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err = sm.reportIndicationOnChange(ctx, subscription)
			if err != nil {
				return
			}
		}()
	case e2smrcpreies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_PERIODIC:
		log.Debug("Received periodic report subscription request")
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			interval, err := sm.getReportPeriod(request)
			if err != nil {
				log.Error(err)
				return
			}
			err = sm.reportPeriodicIndication(ctx, interval, subscription)
			if err != nil {
				return
			}
		}()

	}

	return response, nil, nil
}

// RICSubscriptionDelete implements subscription delete handler for RC service model
func (sm *Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Infof("Ric subscription delete request is received for service model %v and e2 node with ID: %d", sm.ServiceModel.ModelName, sm.ServiceModel.Node.GnbID)
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
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerAsnBytes := sub.Details.RicEventTriggerDefinition.Value

	var rcPreServiceModel e2smrcpresm.RcPreServiceModel
	eventTriggerProtoBytes, err := rcPreServiceModel.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerDefinition := &e2smrcpreies.E2SmRcPreEventTriggerDefinition{}
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
	case e2smrcpreies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_PERIODIC:
		log.Debug("Stopping the periodic report subscription")
		sub.Ticker.Stop()
	case e2smrcpreies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_UPON_CHANGE:
		// TODO stop on change event trigger
	}

	return response, nil, nil
}
