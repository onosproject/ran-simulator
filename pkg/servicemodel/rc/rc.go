// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package rc

import (
	"context"
	"time"

	"github.com/onosproject/ran-simulator/pkg/store/metrics"

	"github.com/onosproject/ran-simulator/pkg/store/cells"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"

	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"

	controlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"

	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"

	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"

	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/pdubuilder"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"

	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger("sm", "rc")

const (
	modelFullName = "e2sm_rc_pre-v1"
	version       = "v1"
	// TODO - Replace with OID for rc-pre service model
	modelOID      = "1.3.6.1.4.1.1.1.2.99"
)

// Client rc service model client
type Client struct {
	ServiceModel *registry.ServiceModel
}

func (sm *Client) reportPeriodicIndication(ctx context.Context, interval int32, subscription *subutils.Subscription) error {
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
			_ = sm.sendRicIndication(ctx, subscription)

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
	for _, ecgi := range node.Cells {
		ricIndication, err := sm.createRicIndication(ctx, ecgi, subscription)
		if err != nil {
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
	log.Debugf("Sending report indication on change from node: %d", sm.ServiceModel.Node.EnbID)
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}
	ch := make(chan event.Event)
	nodeCells := sm.ServiceModel.Node.Cells
	err = sm.ServiceModel.CellStore.Watch(context.Background(), ch)
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
		case cellEvent := <-ch:
			cellEventType := cellEvent.Type.(cells.CellEvent)
			if cellEventType == cells.UpdatedNeighbors {
				cell := cellEvent.Value.(*model.Cell)
				for _, nodeCell := range nodeCells {
					if nodeCell == cell.ECGI {
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
	modelPluginRegistry *modelplugins.ModelPluginRegistry,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store,
	ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store) (registry.ServiceModel, error) {
	modelFullName := modelplugins.ModelFullName(modelFullName)
	rcSm := registry.ServiceModel{
		RanFunctionID:       registry.Rc,
		ModelFullName:       modelFullName,
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
	}

	rcClient := &Client{
		ServiceModel: &rcSm,
	}

	rcSm.Client = rcClient

	var ranFunctionShortName = string(modelFullName)
	var ranFunctionE2SmOid = "OID124"
	var ranFunctionDescription = "RC PRE"
	var ranFunctionInstance int32 = 3
	var ricEventStyleType int32 = 1
	var ricEventStyleName = "Periodic and On Change Report"
	var ricEventFormatType int32 = 1
	var ricReportStyleType int32 = 1
	var ricReportStyleName = "PCI and NRT update for eNB"
	var ricIndicationHeaderFormatType int32 = 1
	var ricIndicationMessageFormatType int32 = 1
	ranFuncDescPdu, err := pdubuilder.CreateE2SmRcPreRanfunctionDescriptionMsg(ranFunctionShortName, ranFunctionE2SmOid, ranFunctionDescription,
		ranFunctionInstance, ricEventStyleType, ricEventStyleName, ricEventFormatType, ricReportStyleType, ricReportStyleName,
		ricIndicationHeaderFormatType, ricIndicationMessageFormatType)

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
	rcModelPlugin := modelPluginRegistry.ModelPlugins[modelFullName]
	if rcModelPlugin == nil {
		log.Debug("model plugin names:", modelPluginRegistry.ModelPlugins)
		return registry.ServiceModel{}, errors.New(errors.Invalid, "model plugin is nil")
	}
	ranFuncDescBytes, err := rcModelPlugin.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	rcSm.Description = ranFuncDescBytes
	return rcSm, nil
}

// RICControl implements control handler for RC service model
func (sm *Client) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	log.Infof("Control Request is received for service model %v and e2 node ID: %d", sm.ServiceModel.ModelFullName, sm.ServiceModel.Node.EnbID)
	reqID := controlutils.GetRequesterID(request)
	ranFuncID := controlutils.GetRanFunctionID(request)
	ricInstanceID := controlutils.GetRicInstanceID(request)

	controlMessage, err := sm.getControlMessage(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	log.Debugf("Control Message Proto: %+v", controlMessage)

	controlHeader, err := sm.getControlHeader(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	log.Debugf("Control Header Proto: %+v", controlHeader)
	// TODO implement RC control logic

	response, _ = controlutils.NewControl(
		controlutils.WithRanFuncID(ranFuncID),
		controlutils.WithRequestID(reqID),
		controlutils.WithRicInstanceID(ricInstanceID),
		controlutils.WithRicControlOutcome(e2aptypes.RicControlOutcome("OK"))).BuildControlAcknowledge()
	return response, nil, err
}

// RICSubscription implements subscription handler for RC service model
func (sm *Client) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Infof("Ric Subscription Request is received for service model %v and e2 node with ID:%d", sm.ServiceModel.ModelFullName, sm.ServiceModel.Node.EnbID)
	var ricActionsAccepted []*e2aptypes.RicActionID
	ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
	actionList := subutils.GetRicActionToBeSetupList(request)
	reqID := subutils.GetRequesterID(request)
	ranFuncID := subutils.GetRanFunctionID(request)
	ricInstanceID := subutils.GetRicInstanceID(request)

	for _, action := range actionList {
		actionID := e2aptypes.RicActionID(action.Value.RicActionId.Value)
		actionType := action.Value.RicActionType
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
					RicRequest: e2apies.CauseRic_CAUSE_RIC_ACTION_NOT_SUPPORTED,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
	}
	subscription := subutils.NewSubscription(
		subutils.WithRequestID(reqID),
		subutils.WithRanFuncID(ranFuncID),
		subutils.WithRicInstanceID(ricInstanceID),
		subutils.WithActionsAccepted(ricActionsAccepted),
		subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, subscriptionFailure, nil
	}

	response, err = subscription.BuildSubscriptionResponse()
	if err != nil {
		return nil, nil, err
	}

	eventTriggerType, err := sm.getEventTriggerType(request)
	if err != nil {
		return nil, nil, err
	}

	switch eventTriggerType {
	case e2sm_rc_pre_ies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_UPON_CHANGE:
		log.Debugf("Event trigger is on change")
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err = sm.reportIndicationOnChange(ctx, subscription)
			if err != nil {
				return
			}
		}()
	case e2sm_rc_pre_ies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_PERIODIC:
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
	log.Infof("Ric subscription delete request is received for service model %v and e2 node with ID: %d", sm.ServiceModel.ModelFullName, sm.ServiceModel.Node.EnbID)
	reqID := subdeleteutils.GetRequesterID(request)
	ranFuncID := subdeleteutils.GetRanFunctionID(request)
	ricInstanceID := subdeleteutils.GetRicInstanceID(request)
	subID := subscriptions.NewID(ricInstanceID, reqID, ranFuncID)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerAsnBytes := sub.Details.RicEventTriggerDefinition.Value
	rcModelPlugin := sm.ServiceModel.ModelPluginRegistry.ModelPlugins[sm.ServiceModel.ModelFullName]
	eventTriggerProtoBytes, err := rcModelPlugin.EventTriggerDefinitionASN1toProto(eventTriggerAsnBytes)
	if err != nil {
		return nil, nil, err
	}
	eventTriggerDefinition := &e2sm_rc_pre_ies.E2SmRcPreEventTriggerDefinition{}
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
	case e2sm_rc_pre_ies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_PERIODIC:
		sub.Ticker.Stop()
	case e2sm_rc_pre_ies.RcPreTriggerType_RC_PRE_TRIGGER_TYPE_UPON_CHANGE:
		// TODO stop on change event trigger
	}

	return response, nil, nil
}
