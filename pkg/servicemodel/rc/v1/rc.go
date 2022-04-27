// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"context"
	"github.com/gogo/protobuf/proto"
	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/pdubuilder"
	e2smrc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/servicemodel"
	e2smrcies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-rc-ies"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	controlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger()

// Client rc service model client
type Client struct {
	ServiceModel *registry.ServiceModel
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store,
	ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store) (registry.ServiceModel, error) {
	var rcsm e2smrc.RCServiceModel
	modelName := e2smtypes.ShortName(modelFullName)
	rcSm := registry.ServiceModel{
		RanFunctionID: registry.Rc,
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

	rcRANFuncDescPDU, err := pdubuilder.CreateE2SmRcRanfunctionDefinition(modelFullName, modelOID, "RAN Control")
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Event Trigger style 2: Call Process Breakpoint
	ricEventTriggerStyle2, err := pdubuilder.CreateRanfunctionDefinitionEventTriggerStyleItem(2, eventTriggerStyle2, 2)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Event trigger style 3: E2 Node Information Change
	ricEventTriggerStyle3, err := pdubuilder.CreateRanfunctionDefinitionEventTriggerStyleItem(3, eventTriggerStyle3, 3)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Create event trigger style list
	ricEventTriggerStyleList := make([]*e2smrcies.RanfunctionDefinitionEventTriggerStyleItem, 0)
	ricEventTriggerStyleList = append(ricEventTriggerStyleList, ricEventTriggerStyle2)
	ricEventTriggerStyleList = append(ricEventTriggerStyleList, ricEventTriggerStyle3)
	ranFunctionDefinitionEventTrigger, err := pdubuilder.CreateRanfunctionDefinitionEventTrigger(ricEventTriggerStyleList)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Create Report Style 2: Call Process Outcome, this style is used to report the outcome of an ongoing call process.
	reportStyleItem2, err := pdubuilder.CreateRanfunctionDefinitionReportItem(2, "Call Process Outcome", 2, 1, 1, 2)
	if err != nil {
		return registry.ServiceModel{}, err
	}
	reportParametersReportStyle2List, err := createRANParametersReportStyle2List()
	if err != nil {
		return registry.ServiceModel{}, err
	}
	reportStyleItem2.SetRanReportParametersList(reportParametersReportStyle2List)

	// Create Report Style 3:  E2 Node Information. This style is used to report E2 Node information, Serving Cell Configuration and Neighbour Relation related information.
	reportStyleItem3, err := pdubuilder.CreateRanfunctionDefinitionReportItem(3, "E2 Node information", 3, 1, 1, 3)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	reportParametersReportStyle3List, err := createRANParametersReportStyle3List()
	if err != nil {
		return registry.ServiceModel{}, err
	}
	reportStyleItem3.SetRanReportParametersList(reportParametersReportStyle3List)

	// Add report styles to report style list
	reportStyleList := make([]*e2smrcies.RanfunctionDefinitionReportItem, 0)
	reportStyleList = append(reportStyleList, reportStyleItem2)
	reportStyleList = append(reportStyleList, reportStyleItem3)
	ranFunctionDefinitionReport, err := pdubuilder.CreateRanfunctionDefinitionReport(reportStyleList)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Create RAN Function Definition Insert Indication item (RIC Indication 1 for Handover Control Request)
	ranFunctionDefinitionInsertItem, err := pdubuilder.CreateRanfunctionDefinitionInsertIndicationItem(1, "Handover Control Request")
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Create List of RAN parameters for RIC Indication 1
	insertParametersInsertStyle3List, err := createRANParametersInsertStyle3List()
	if err != nil {
		return registry.ServiceModel{}, err
	}
	ranFunctionDefinitionInsertItem.SetRanInsertIndicationParametersList(insertParametersInsertStyle3List)

	// Create Insert Style 3: Connected Mode Mobility Control Request
	insertStyleItem3, err := pdubuilder.CreateRanfunctionDefinitionInsertItem(3, "Connected Mode Mobility Control Request", 2, 3, 2, 5, 1)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	insertIndicationList := make([]*e2smrcies.RanfunctionDefinitionInsertIndicationItem, 0)
	insertIndicationList = append(insertIndicationList, ranFunctionDefinitionInsertItem)
	insertStyleItem3.SetRicInsertIndicationList(insertIndicationList)

	//  Add insert styles to insert style list
	insertStyleList := make([]*e2smrcies.RanfunctionDefinitionInsertItem, 0)
	insertStyleList = append(insertStyleList, insertStyleItem3)
	ranFunctionDefinitionInsert, err := pdubuilder.CreateRanfunctionDefinitionInsert(insertStyleList)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Creates RAN function definition policy list
	ranFunctionDefinitionPolicyList := make([]*e2smrcies.RanfunctionDefinitionPolicyItem, 0)
	// Creates policy items
	ranFunctionPolicyItem1, err := pdubuilder.CreateRanfunctionDefinitionPolicyItem(3, "Connected Mobile Mobility", 2)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	ranFunctionDefinitionPolicyList = append(ranFunctionDefinitionPolicyList, ranFunctionPolicyItem1)
	ranFunctionDefinitionPolicy, err := pdubuilder.CreateRanfunctionDefinitionPolicy(ranFunctionDefinitionPolicyList)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Creates RAN function definition control list
	controlItemList := make([]*e2smrcies.RanfunctionDefinitionControlItem, 0)
	// Creates control action list
	controlActionList := make([]*e2smrcies.RanfunctionDefinitionControlActionItem, 0)
	controlActionItem1, err := pdubuilder.CreateRanfunctionDefinitionControlActionItem(1, "PCI Control")
	if err != nil {
		return registry.ServiceModel{}, err
	}

	ranControlActionParametersList1 := make([]*e2smrcies.ControlActionRanparameterItem, 0)
	controlActionRANParameterItem1 := &e2smrcies.ControlActionRanparameterItem{
		RanParameterId: &e2smrcies.RanparameterId{
			Value: 1,
		},
		RanParameterName: &e2smrcies.RanparameterName{
			Value: "Serving Cell NR PCI",
		},
	}
	ranControlActionParametersList1 = append(ranControlActionParametersList1, controlActionRANParameterItem1)

	controlActionItem1.SetRanControlActionParametersList(ranControlActionParametersList1)
	controlActionList = append(controlActionList, controlActionItem1)
	controlItem1, err := pdubuilder.CreateRanfunctionDefinitionControlItem(9, "PCI Control", 1, 1, 1)
	if err != nil {
		return registry.ServiceModel{}, err
	}
	controlItem1.SetRicControlActionList(controlActionList)
	controlItemList = append(controlItemList, controlItem1)

	ranFunctionDefinitionControl, err := pdubuilder.CreateRanfunctionDefinitionControl(controlItemList)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Sets RAN function report definition
	rcRANFuncDescPDU.SetRanFunctionDefinitionReport(ranFunctionDefinitionReport)
	// Sets RAN function event trigger definition
	rcRANFuncDescPDU.SetRanFunctionDefinitionEventTrigger(ranFunctionDefinitionEventTrigger)
	// Sets RAN function insert definition
	rcRANFuncDescPDU.SetRanFunctionDefinitionInsert(ranFunctionDefinitionInsert)
	// Sets RAN function policy definition
	rcRANFuncDescPDU.SetRanFunctionDefinitionPolicy(ranFunctionDefinitionPolicy)

	// Sets RAN function control definition
	rcRANFuncDescPDU.SetRanFunctionDefinitionControl(ranFunctionDefinitionControl)

	protoBytes, err := proto.Marshal(rcRANFuncDescPDU)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	ranFuncDescBytes, err := rcsm.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	rcSm.Description = ranFuncDescBytes
	return rcSm, nil

}

// E2ConnectionUpdate implements connection update handler
func (c *Client) E2ConnectionUpdate(ctx context.Context, request *e2appducontents.E2ConnectionUpdate) (response *e2appducontents.E2ConnectionUpdateAcknowledge, failure *e2appducontents.E2ConnectionUpdateFailure, err error) {
	return nil, nil, errors.NewNotSupported("E2 connection update is not supported")

}

// RICControl implements control handler for RC service model
func (c *Client) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	log.Infof("Control Request is received for service model %v and e2 node ID: %d", c.ServiceModel.ModelName, c.ServiceModel.Node.GnbID)
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

	controlMessage, err := getControlMessage(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	controlHeader, err := getControlHeader(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	log.Debugf("RC control header: %v", controlHeader)
	log.Debugf("RC control message: %v", controlMessage)

	// TODO Add control outcome if needed
	response, err = controlutils.NewControl(
		controlutils.WithRanFuncID(*ranFuncID),
		controlutils.WithRequestID(*reqID),
		controlutils.WithRicInstanceID(*ricInstanceID),
		controlutils.WithRicControlOutcome(nil)).BuildControlAcknowledge()
	if err != nil {
		return nil, nil, err
	}
	return response, nil, nil
}

// RICSubscription implements subscription handler for RC service model
func (c *Client) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Infof("Ric Subscription Request is received for service model %v and e2 node with ID: %v", c.ServiceModel.ModelName, c.ServiceModel.Node.GnbID)
	var ricActionsAccepted []*e2aptypes.RicActionID
	actionList := subutils.GetRicActionToBeSetupList(request)
	ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
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
		// rc service model supports report, policy, and inserts action and should be added to the
		// list of accepted actions
		if actionType == e2apies.RicactionType_RICACTION_TYPE_REPORT ||
			actionType == e2apies.RicactionType_RICACTION_TYPE_INSERT ||
			actionType == e2apies.RicactionType_RICACTION_TYPE_POLICY {
			ricActionsAccepted = append(ricActionsAccepted, &actionID)
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

	// Process RC event triggers
	eventTriggers, err := getEventTrigger(request)
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

	// Process RC action Definitions to create a map of action ID and action definition
	actionDefinitionsMaps, err := getActionDefinitionMap(actionList, ricActionsAccepted)
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_INCONSISTENT_ACTION_SUBSEQUENT_ACTION_SEQUENCE,
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
	log.Debugf("Action Definitions map: %+v", actionDefinitionsMaps)

	eventTriggerFormats := eventTriggers.GetRicEventTriggerFormats()
	switch eventTriggerFormats.RicEventTriggerFormats.(type) {
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat1:
		// TODO Process RIC Event trigger definition IE style 1: Message Event
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat2:
		// TODO Process RIC Event trigger definition IE style 2: Call Process Breakpoint
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat3:
		// TODO Process RIC Event trigger definition IE style 3: E2 Node Information Change
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat4:
		// TODO Process RIC Event trigger definition IE style 4: UE Information Change

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

	return response, nil, nil

}

// RICSubscriptionDelete implements subscription delete handler for RC service model
func (c *Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	//TODO implement me
	log.Info("implement me")
	return nil, nil, nil
}
