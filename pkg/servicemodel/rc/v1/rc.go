// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"context"
	"github.com/gogo/protobuf/proto"
	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/pdubuilder"
	e2smrc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/servicemodel"
	e2smcommonies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-common-ies"
	e2smrcies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc/v1/e2sm-rc-ies"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/mobility"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	controlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"
	"github.com/onosproject/rrm-son-lib/pkg/handover"
	"github.com/onosproject/rrm-son-lib/pkg/model/id"
	meastype "github.com/onosproject/rrm-son-lib/pkg/model/measurement/type"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger()

// Client rc service model client
type Client struct {
	ServiceModel   *registry.ServiceModel
	mobilityDriver mobility.Driver
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store,
	ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store,
	a3Chan chan handover.A3HandoverDecision, mobilityDriver mobility.Driver) (registry.ServiceModel, error) {
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
		A3Chan:        a3Chan,
	}

	rcClient := &Client{
		ServiceModel:   &rcSm,
		mobilityDriver: mobilityDriver,
	}

	rcSm.Client = rcClient

	rcRANFuncDescPDU, err := pdubuilder.CreateE2SmRcRanfunctionDefinition(modelFullName, modelOID, "RAN Control")
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// Event Trigger style 1: Message Event
	ricEventTriggerStyle1, err := pdubuilder.CreateRanfunctionDefinitionEventTriggerStyleItem(1, eventTriggerStyle1, 1)
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
	ricEventTriggerStyleList = append(ricEventTriggerStyleList, ricEventTriggerStyle1)
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
	ranFunctionDefinitionInsertItem, err := pdubuilder.CreateRanfunctionDefinitionInsertIndicationItem(ricInsertIndicationIDForMHO, "Handover Control Request")
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
	insertStyleItem3, err := pdubuilder.CreateRanfunctionDefinitionInsertItem(3, "Connected Mode Mobility Control Request", 1, 3, 2, 5, 1)
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

	// Create Policy List
	ranFunctionDefinitionPolicyList := make([]*e2smrcies.RanfunctionDefinitionPolicyItem, 0)

	// policy item
	ranFunctionDefinitionPolicyItem, err := pdubuilder.CreateRanfunctionDefinitionPolicyItem(ricPolicyStyleType3, ricPolicyStyleName, 1)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	// policy actions
	ranFunctionDefinitionPolicyActionItem, err := pdubuilder.CreateRanfunctionDefinitionPolicyActionItem(ricPolicyActionIDForMLB, ricPolicyActionNameForMLB, ricActionDefinitionFormatTypeForMLB)
	if err != nil {
		return registry.ServiceModel{}, err
	}

	//associated ran parameters for policy action
	policyActionRANParameterItem1 := &e2smrcies.PolicyActionRanparameterItem{
		RanParameterId: &e2smrcies.RanparameterId{
			Value: TargetPrimaryCellIDRANParameterID,
		},
		RanParameterName: &e2smrcies.RanparameterName{
			Value: TargetPrimaryCellIDRANParameterName,
		},
	}
	policyActionRANParameterItem2 := &e2smrcies.PolicyActionRanparameterItem{
		RanParameterId: &e2smrcies.RanparameterId{
			Value: CellSpecificOffsetRANParameterID,
		},
		RanParameterName: &e2smrcies.RanparameterName{
			Value: CellSpecificOffsetRANParameterName,
		},
	}
	policyActionRANParameterList := make([]*e2smrcies.PolicyActionRanparameterItem, 0)
	policyActionRANParameterList = append(policyActionRANParameterList, policyActionRANParameterItem1)
	policyActionRANParameterList = append(policyActionRANParameterList, policyActionRANParameterItem2)

	//associated ran parameters for policy condition
	policyConditionRANParameterItem1 := &e2smrcies.PolicyConditionRanparameterItem{
		RanParameterId: &e2smrcies.RanparameterId{
			Value: TargetPrimaryCellIDRANParameterID,
		},
		RanParameterName: &e2smrcies.RanparameterName{
			Value: TargetPrimaryCellIDRANParameterName,
		},
	}
	policyConditionRANParameterItem2 := &e2smrcies.PolicyConditionRanparameterItem{
		RanParameterId: &e2smrcies.RanparameterId{
			Value: CellSpecificOffsetRANParameterID,
		},
		RanParameterName: &e2smrcies.RanparameterName{
			Value: CellSpecificOffsetRANParameterName,
		},
	}
	policyConditionRANParameterList := make([]*e2smrcies.PolicyConditionRanparameterItem, 0)
	policyConditionRANParameterList = append(policyConditionRANParameterList, policyConditionRANParameterItem1)
	policyConditionRANParameterList = append(policyConditionRANParameterList, policyConditionRANParameterItem2)

	ranFunctionDefinitionPolicyActionItem.SetRanPolicyActionParametersList(policyActionRANParameterList)
	ranFunctionDefinitionPolicyActionItem.SetRanPolicyConditionParametersList(policyConditionRANParameterList)

	ranFunctionDefinitionPolicyActionList := make([]*e2smrcies.RanfunctionDefinitionPolicyActionItem, 0)
	ranFunctionDefinitionPolicyActionList = append(ranFunctionDefinitionPolicyActionList, ranFunctionDefinitionPolicyActionItem)

	ranFunctionDefinitionPolicyItem.SetRicPolicyActionList(ranFunctionDefinitionPolicyActionList)

	ranFunctionDefinitionPolicyList = append(ranFunctionDefinitionPolicyList, ranFunctionDefinitionPolicyItem)
	ranFunctionDefinitionPolicy, err := pdubuilder.CreateRanfunctionDefinitionPolicy(ranFunctionDefinitionPolicyList)
	if err != nil {
		return registry.ServiceModel{}, err
	}
	// Creates RAN function definition control list
	controlItemList := make([]*e2smrcies.RanfunctionDefinitionControlItem, 0)
	// for PCI
	// Creates control action list
	controlActionList1 := make([]*e2smrcies.RanfunctionDefinitionControlActionItem, 0)
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
	controlActionRANParameterItem2 := &e2smrcies.ControlActionRanparameterItem{
		RanParameterId: &e2smrcies.RanparameterId{
			Value: 2,
		},
		RanParameterName: &e2smrcies.RanparameterName{
			Value: "Serving Cell CGI",
		},
	}

	ranControlActionParametersList1 = append(ranControlActionParametersList1, controlActionRANParameterItem1)
	ranControlActionParametersList1 = append(ranControlActionParametersList1, controlActionRANParameterItem2)

	controlActionItem1.SetRanControlActionParametersList(ranControlActionParametersList1)
	controlActionList1 = append(controlActionList1, controlActionItem1)

	controlItem1, err := pdubuilder.CreateRanfunctionDefinitionControlItem(controlStyleType200, "PCI Control", 1, 1, 1)
	if err != nil {
		return registry.ServiceModel{}, err
	}
	controlItem1.SetRicControlActionList(controlActionList1)

	// For MHO
	// Creates control action list
	controlActionList2 := make([]*e2smrcies.RanfunctionDefinitionControlActionItem, 0)

	controlActionItem2, err := pdubuilder.CreateRanfunctionDefinitionControlActionItem(1, "Handover Control")
	if err != nil {
		return registry.ServiceModel{}, err
	}

	ranControlActionParametersList2 := make([]*e2smrcies.ControlActionRanparameterItem, 0)
	controlActionRANParameterItem3 := &e2smrcies.ControlActionRanparameterItem{
		RanParameterId: &e2smrcies.RanparameterId{
			Value: 1,
		},
		RanParameterName: &e2smrcies.RanparameterName{
			Value: "Target Primary Cell ID",
		},
	}
	ranControlActionParametersList2 = append(ranControlActionParametersList2, controlActionRANParameterItem3)

	controlActionItem2.SetRanControlActionParametersList(ranControlActionParametersList2)
	controlActionList2 = append(controlActionList2, controlActionItem2)

	controlItem2, err := pdubuilder.CreateRanfunctionDefinitionControlItem(controlStyleType3, "Connected Mode Mobility", 1, 1, 1)
	if err != nil {
		return registry.ServiceModel{}, err
	}
	controlItem2.SetRicControlActionList(controlActionList2)

	controlItemList = append(controlItemList, controlItem1)
	controlItemList = append(controlItemList, controlItem2)

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
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_CONTROL_MESSAGE_INVALID,
			},
		}
		failure, err = controlutils.NewControl(
			controlutils.WithRanFuncID(*ranFuncID),
			controlutils.WithRequestID(*reqID),
			controlutils.WithRicInstanceID(*ricInstanceID),
			controlutils.WithCause(cause),
			controlutils.WithRicControlOutcome(nil)).BuildControlFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	controlHeader, err := getControlHeader(request)
	if err != nil {
		log.Error(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		failure, err = controlutils.NewControl(
			controlutils.WithRanFuncID(*ranFuncID),
			controlutils.WithRequestID(*reqID),
			controlutils.WithRicInstanceID(*ricInstanceID),
			controlutils.WithCause(cause),
			controlutils.WithRicControlOutcome(nil)).BuildControlFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	log.Debugf("RC control header: %v", controlHeader)
	log.Debugf("RC control message: %v", controlMessage)

	// Check if the control request is for changing the PCI value to change it PCI
	err = c.handleControlMessage(ctx, controlHeader, controlMessage)
	if err != nil {
		log.Error(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_CONTROL_MESSAGE_INVALID,
			},
		}
		failure, err = controlutils.NewControl(
			controlutils.WithRanFuncID(*ranFuncID),
			controlutils.WithRequestID(*reqID),
			controlutils.WithRicInstanceID(*ricInstanceID),
			controlutils.WithCause(cause),
			controlutils.WithRicControlOutcome(nil)).BuildControlFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

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
	log.Debugf("Received subscription request message %+v", request.ProtocolIes)
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
	log.Infof("Action Definitions map: %+v", actionDefinitionsMaps)

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

	for _, action := range actionList {
		if action.GetValue().GetRicactionToBeSetupItem().GetRicActionType() == e2apies.RicactionType_RICACTION_TYPE_REPORT {
			log.Debugf("Processing Report Action for e2 Node %v", c.ServiceModel.Node.GnbID)
			err := c.processReportAction(ctx, subscription, eventTriggers)
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
		} else if action.GetValue().GetRicactionToBeSetupItem().GetRicActionType() == e2apies.RicactionType_RICACTION_TYPE_INSERT {
			log.Debugf("Processing Insert Action for e2 Node %v", c.ServiceModel.Node.GnbID)
			err := c.processInsertAction(ctx, subscription, eventTriggers)
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
		} else if action.GetValue().GetRicactionToBeSetupItem().GetRicActionType() == e2apies.RicactionType_RICACTION_TYPE_POLICY {
			log.Debugf("Processing Policy Action for e2 Node %v", c.ServiceModel.Node.GnbID)
			err := c.processPolicyAction(ctx, subscription, eventTriggers, actionDefinitionsMaps)
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
		}
	}

	return response, nil, nil

}

// RICSubscriptionDelete implements subscription delete handler for RC service model
func (c *Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Debugf("Received delete subscription request message %+v", request.ProtocolIes)
	log.Infof("RIC subscription delete request is received for e2 node %d and service model %s", c.ServiceModel.Node.GnbID, c.ServiceModel.ModelName)
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
	sub, err := c.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
		subdeleteutils.WithRequestID(*reqID),
		subdeleteutils.WithRanFuncID(*ranFuncID),
		subdeleteutils.WithRicInstanceID(*ricInstanceID))
	subDeleteResponse, err := subscriptionDelete.BuildSubscriptionDeleteResponse()
	if err != nil {
		return nil, nil, err
	}
	// Stops the goroutine sending the indication messages
	if sub.Ticker != nil {
		sub.Ticker.Stop()
	}
	return subDeleteResponse, nil, nil
}

func (c *Client) processMLBLogic(ctx context.Context, adFormat2List []*e2smrcies.E2SmRcActionDefinitionFormat2Item) error {
	for _, pc := range adFormat2List {
		ncgi, err := c.extractNCGIFromPrintableNCGI(pc.GetRicPolicyAction())
		if err != nil {
			return err
		}
		ocnInt, err := c.extractOcn(pc.GetRicPolicyAction())
		if err != nil {
			log.Error(err)
			return err
		}

		ocn := meastype.QOffsetRange(ocnInt)

		for _, id := range c.ServiceModel.Node.Cells {
			// id: serving cell ID
			log.Debugf("MLB: sCell NCGI: %v / NCGI: %v / Ocn: %d", id, ncgi, ocnInt)

			sCell, err := c.ServiceModel.CellStore.Get(ctx, id)
			if err != nil {
				log.Errorf("NCGI (%v) is not in cell store", id)
				continue
			}
			if _, ok := sCell.MeasurementParams.NCellIndividualOffsets[ncgi]; !ok {
				log.Errorf("the cell NCGI (%v) is not a neighbor of the cell NCGI (%v)", ncgi, id)
				continue
			}
			log.Infof("Cell (%v) Ocn in the cell (%v) is set from %v to %v", ncgi, id, sCell.MeasurementParams.NCellIndividualOffsets[ncgi], ocn.GetValue().(int))
			sCell.MeasurementParams.NCellIndividualOffsets[ncgi] = int32(ocn.GetValue().(int))
		}
	}
	return nil
}

func (c *Client) processPolicyAction(ctx context.Context, subscription *subutils.Subscription, eventTriggers *e2smrcies.E2SmRcEventTrigger, actionDefinitionsMaps map[*e2aptypes.RicActionID]*e2smrcies.E2SmRcActionDefinition) error {
	eventTriggerFormats := eventTriggers.GetRicEventTriggerFormats()
	switch eventTrigger := eventTriggerFormats.RicEventTriggerFormats.(type) {
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat1:
		// Process RIC Event trigger definition IE style 1: Message Event
		messageList := eventTrigger.EventTriggerFormat1.GetMessageList()
		for _, m := range messageList {
			ueEventList := m.GetAssociatedUeevent().GetUeEventList()
			for _, e := range ueEventList {
				if e.GetUeEventId().Value == A3MeasurementReportUEEventID {
					log.Debugf("Processing event trigger format 1: Message Event - A3 measurement report received (UE Event ID: %d)", A3MeasurementReportUEEventID)
					log.Debugf("Action definition maps: %+v", actionDefinitionsMaps)
					for _, ad := range actionDefinitionsMaps {
						if adFormat2 := ad.GetRicActionDefinitionFormats().GetActionDefinitionFormat2(); adFormat2 != nil {
							err := c.processMLBLogic(ctx, adFormat2.GetRicPolicyConditionsList())
							if err != nil {
								log.Warn(err)
							}
						}
					}
				}
			}
		}
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat2:
		// TODO Process RIC Event trigger definition IE style 2: Call Process Breakpoint
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat3:
		// TODO Process RIC Event trigger definition IE style 3
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat4:
		// TODO Process RIC Event trigger definition IE style 4
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat5:
		// TODO Process RIC Event trigger definition IE style 5

	}
	return nil
}

func (c *Client) processInsertAction(ctx context.Context, subscription *subutils.Subscription, eventTriggers *e2smrcies.E2SmRcEventTrigger) error {
	eventTriggerFormats := eventTriggers.GetRicEventTriggerFormats()
	switch eventTrigger := eventTriggerFormats.RicEventTriggerFormats.(type) {
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat1:
		// TODO: Process RIC Event trigger definition IE style 1: Message Event
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat2:
		// Process RIC Event trigger definition IE style 2: Call Process Breakpoint
		callProcessTypeID := eventTrigger.EventTriggerFormat2.GetRicCallProcessTypeId().Value
		callBreakPointID := eventTrigger.EventTriggerFormat2.GetRicCallProcessBreakpointId().Value

		// handover
		if callProcessTypeID == CallProcessTypeIDMobilityManagement && callBreakPointID == CallBreakpointIDHandoverPreparation {
			log.Debug("Processing event trigger format 2: Call Process Breakpoint - Mobility Management / Handover Preparation")
			go func() {
				err := c.insertOnA3MeasurementReceived(ctx, subscription)
				if err != nil {
					log.Warn(err)
					// TODO we should propagate this error back
					return
				}
			}()
		}

	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat3:
		// TODO Process RIC Event trigger definition IE style 3
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat4:
		// TODO Process RIC Event trigger definition IE style 4
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat5:
		// TODO Process RIC Event trigger definition IE style 5
	}

	return nil
}

func (c *Client) processReportAction(ctx context.Context, subscription *subutils.Subscription, eventTriggers *e2smrcies.E2SmRcEventTrigger) error {
	eventTriggerFormats := eventTriggers.GetRicEventTriggerFormats()
	switch eventTrigger := eventTriggerFormats.RicEventTriggerFormats.(type) {
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat1:
		// TODO Process RIC Event trigger definition IE style 1: Message Event
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat2:
		// TODO Process RIC Event trigger definition IE style 2: Call Process Breakpoint
	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat3:
		// Process RIC Event trigger definition IE style 3: E2 Node Information Change
		e2NodeInfoChangeList := eventTrigger.EventTriggerFormat3.GetE2NodeInfoChangeList()
		for _, e2NodeChange := range e2NodeInfoChangeList {
			e2NodeInfoChangeID := e2NodeChange.E2NodeInfoChangeId
			if e2NodeInfoChangeID == 1 {
				log.Debugf("Processing event trigger format 3: cell configuration change for e2 Node %v", c.ServiceModel.Node.GnbID)
				go func(e *e2smrcies.E2SmRcEventTriggerFormat3Item) {
					err := c.reportOnCellConfigurationChange(ctx, subscription, e)
					if err != nil {
						log.Warn(err)
						// TODO we should propagate this error back
						return
					}
				}(e2NodeChange)

			} else if e2NodeInfoChangeID == 2 {
				log.Debug("Processing event trigger format 3: cell neighbor relation change for e2 node %v", c.ServiceModel.Node.GnbID)
				go func(e *e2smrcies.E2SmRcEventTriggerFormat3Item) {
					err := c.reportOnCellNeighborRelationChange(ctx, subscription, e)
					if err != nil {
						log.Warn(err)
						// TODO we should propagate this error back
						return
					}
				}(e2NodeChange)

			} else {
				return errors.NewNotSupported("E2 node information change ID %d is not supported", e2NodeInfoChangeID)
			}
		}

	case *e2smrcies.RicEventTriggerFormats_EventTriggerFormat4:
		// TODO Process RIC Event trigger definition IE style 4: UE Information Change

	}

	return nil

}

func (c *Client) reportOnCellConfigurationChange(ctx context.Context, subscription *subutils.Subscription, e2NodeChange *e2smrcies.E2SmRcEventTriggerFormat3Item) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := c.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	node := c.ServiceModel.Node
	cellInfoList := e2NodeChange.AssociatedCellInfo.GetCellInfoList()
	cellList := make([]ransimtypes.NCGI, 0)
	if len(cellInfoList) == 0 {
		cellList = node.Cells
	} // TODO else create a list of cells based on cell info list to report cell changes just for those requested cells
	cellEventCh := make(chan event.Event)
	err = c.ServiceModel.CellStore.Watch(context.Background(), cellEventCh)
	if err != nil {
		return err
	}

	// Sends an initial indication message
	err = c.sendRICIndicationFormat3(ctx, cellList, subscription, e2NodeChange.GetE2NodeInfoChangeId())
	if err != nil {
		return err
	}

	for {
		select {
		case cellEvent := <-cellEventCh:
			log.Debugf("A Cell change event is occurred %v", cellEvent)
			cellEventType := cellEvent.Type.(cells.CellEvent)
			if cellEventType == cells.Updated {
				err = c.sendRICIndicationFormat3(ctx, cellList, subscription, e2NodeChange.GetE2NodeInfoChangeId())
				if err != nil {
					log.Error(err)
					continue
				}

			}
		case <-sub.E2Channel.Context().Done():
			log.Debugf("E2 channel is closed for subscription: %v", subID)
			return nil

		}
	}
}

func (c *Client) reportOnCellNeighborRelationChange(ctx context.Context, subscription *subutils.Subscription, e2NodeChange *e2smrcies.E2SmRcEventTriggerFormat3Item) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := c.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	node := c.ServiceModel.Node
	cellInfoList := e2NodeChange.AssociatedCellInfo.GetCellInfoList()
	cellList := make([]ransimtypes.NCGI, 0)
	if len(cellInfoList) == 0 {
		cellList = node.Cells
	} // TODO else create a list of cells based on cell info list to report cell changes just for those requested cells
	cellEventCh := make(chan event.Event)
	err = c.ServiceModel.CellStore.Watch(context.Background(), cellEventCh)
	if err != nil {
		return err
	}

	// Sends an initial indication message
	err = c.sendRICIndicationFormat3(ctx, cellList, subscription, e2NodeChange.GetE2NodeInfoChangeId())
	if err != nil {
		return err
	}

	for {
		select {
		case cellEvent := <-cellEventCh:
			log.Debugf("A Cell change event is occurred %v", cellEvent)
			cellEventType := cellEvent.Type.(cells.CellEvent)
			if cellEventType == cells.UpdatedNeighbors {
				err = c.sendRICIndicationFormat3(ctx, cellList, subscription, e2NodeChange.GetE2NodeInfoChangeId())
				if err != nil {
					log.Error(err)
					continue
				}

			}
		case <-sub.E2Channel.Context().Done():
			log.Debugf("E2 channel is closed for subscription: %v", subID)
			return nil

		}
	}
}

func (c *Client) sendRICIndicationFormat3(ctx context.Context, cells []ransimtypes.NCGI, subscription *subutils.Subscription, e2NodeInfoChangeID int32) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := c.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}
	// Report all Cell changes using Indication message format 3
	// Creates and sends an indication message for each cell in the node
	ricIndication, err := c.createRICIndicationFormat3(ctx, cells, subscription, e2NodeInfoChangeID)
	if err != nil {
		log.Error(err)
		return err
	}
	err = sub.E2Channel.RICIndication(ctx, ricIndication)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (c *Client) insertOnA3MeasurementReceived(ctx context.Context, subscription *subutils.Subscription) error {
	if c.mobilityDriver.GetHoLogic() == "local" {
		c.mobilityDriver.SetHoLogic("mho")
	}

	log.Info("Start RC insert service for A3 measurement report received")
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := c.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		log.Error(err)
		return err
	}

	for {
		select {
		case <-sub.E2Channel.Context().Done():
			log.Debugf("E2 channel is closed for subscription: %v", subID)
			return nil
		case report := <-c.ServiceModel.A3Chan:
			log.Debugf("received event a3 measurement report: %v", report)
			log.Debugf("Send upon-rcv-meas-report indication for cell ecgi:%d, IMSI:%s",
				report.UE.GetSCell().GetID().GetID().(id.ECGI), report.UE.GetID().String())
			imsi := report.UE.GetID().GetID().(id.UEID).IMSI
			ue, err := c.ServiceModel.UEs.Get(ctx, ransimtypes.IMSI(imsi))
			//gNbID := utils.NewGNbID(uint64(ransimtypes.GetGnbID(uint64(ue.Cell.NCGI))), 22)

			if err != nil {
				log.Warn(err)
				continue
			}

			ueID := &e2smcommonies.Ueid{
				Ueid: &e2smcommonies.Ueid_GNbUeid{
					GNbUeid: &e2smcommonies.UeidGnb{
						AmfUeNgapId: &e2smcommonies.AmfUeNgapId{
							Value: int64(ue.AmfUeNgapID),
						},
						// ToDo - move out GUAMI hardcoding
						Guami: &e2smcommonies.Guami{
							PLmnidentity: &e2smcommonies.Plmnidentity{
								Value: c.getPlmnID().ToBytes(),
							},
							AMfregionId: &e2smcommonies.AmfregionId{
								Value: &asn1.BitString{
									Value: []byte{0xDD},
									Len:   8,
								},
							},
							AMfsetId: &e2smcommonies.AmfsetId{
								Value: &asn1.BitString{
									Value: []byte{0xCC, 0xC0},
									Len:   10,
								},
							},
							AMfpointer: &e2smcommonies.Amfpointer{
								Value: &asn1.BitString{
									Value: []byte{0xFC},
									Len:   6,
								},
							},
						},
						// remain below as comments. all those are optional fields but raising errors - todo: need to check later
						//GNbCuUeF1ApIdList:   &e2smcommonies.UeidGnbCuF1ApIdList{Value: []*e2smcommonies.UeidGnbCuCpF1ApIdItem{{GNbCuUeF1ApId: &e2smcommonies.GnbCuUeF1ApId{Value: 0}}}},
						//GNbCuCpUeE1ApIdList: &e2smcommonies.UeidGnbCuCpE1ApIdList{Value: []*e2smcommonies.UeidGnbCuCpE1ApIdItem{{GNbCuCpUeE1ApId: &e2smcommonies.GnbCuCpUeE1ApId{Value: 0}}}},
						//RanUeid:             &e2smcommonies.Ranueid{Value: []byte{0, 0, 0, 0, 0, 0, 0, 0}}, // TODO update it with C-RNTI
						//MNgRanUeXnApId:      &e2smcommonies.NgRannodeUexnApid{Value: 0},
						//GlobalGnbId: &e2smcommonies.GlobalGnbId{
						//	PLmnidentity: &e2smcommonies.Plmnidentity{Value: c.getPlmnID().ToBytes()},
						//	GNbId:        &e2smcommonies.GnbId{GnbId: &e2smcommonies.GnbId_GNbId{GNbId: &asn1.BitString{Value: gNbID.IDByte.Bytes(gNbID.Length), Len: uint32(gNbID.Length)}}},
						//},
						//GlobalNgRannodeId: &e2smcommonies.GlobalNgrannodeId{
						//	GlobalNgrannodeId: &e2smcommonies.GlobalNgrannodeId_GNb{
						//		GNb: &e2smcommonies.GlobalGnbId{
						//			PLmnidentity: &e2smcommonies.Plmnidentity{Value: c.getPlmnID().ToBytes()},
						//			GNbId:        &e2smcommonies.GnbId{GnbId: &e2smcommonies.GnbId_GNbId{GNbId: &asn1.BitString{Value: gNbID.IDByte.Bytes(gNbID.Length), Len: uint32(gNbID.Length)}}},
						//		},
						//	},
						//},
					},
				},
			}

			err = c.sendRICIndicationFormat5Header2(ctx, subscription, ueID, ricInsertStyleType3, ricInsertIndicationIDForMHO, ransimtypes.NCGI(report.TargetCell.GetID().GetID().(id.ECGI)))
			if err != nil {
				log.Warn(err)
				continue
			}
		}
	}
}

func (c *Client) sendRICIndicationFormat5Header2(ctx context.Context, subscription *subutils.Subscription, ueID *e2smcommonies.Ueid, ricInsertStyleType int32, insertIndicationID int32, targetNCGI ransimtypes.NCGI) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := c.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}
	ricIndication, err := c.createRICIndicationFormat5Header2(ctx, subscription, ueID, ricInsertStyleType, insertIndicationID, targetNCGI)
	if err != nil {
		return err
	}
	err = sub.E2Channel.RICIndication(ctx, ricIndication)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
