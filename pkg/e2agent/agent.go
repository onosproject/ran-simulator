// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"context"
	"fmt"
	"github.com/onosproject/ran-simulator/pkg/mobility"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/mho"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
	"time"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm2"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/onosproject/ran-simulator/pkg/store/metrics"

	"github.com/onosproject/ran-simulator/pkg/store/cells"

	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"

	controlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/rc"

	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"

	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"

	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"

	"github.com/cenkalti/backoff"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm"
	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/setup"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"

	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap101"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

var log = logging.GetLogger("e2agent")

// E2Agent is an E2 agent
type E2Agent interface {
	// Start starts the agent
	Start() error

	// Stop stops the agent
	Stop() error
}

// e2Agent is an E2 agent
type e2Agent struct {
	node      model.Node
	model     *model.Model
	channel   e2.ClientChannel
	registry  *registry.ServiceModelRegistry
	subStore  *subscriptions.Subscriptions
	nodeStore nodes.Store
	ueStore   ues.Store
	cellStore cells.Store
}

// NewE2Agent creates a new E2 agent
func NewE2Agent(node model.Node, model *model.Model, modelPluginRegistry modelplugins.ModelRegistry,
	nodeStore nodes.Store, ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store,
	measChan chan device.UE, mobilityDriver mobility.Driver) (E2Agent, error) {
	log.Info("Creating New E2 Agent for node with eNbID:", node.GnbID)
	reg := registry.NewServiceModelRegistry()

	// Each new e2 agent has its own subscription store
	subStore := subscriptions.NewStore()
	sms := node.ServiceModels
	for _, smID := range sms {
		serviceModel, err := model.GetServiceModel(smID)
		if err != nil {
			return nil, err
		}
		switch registry.RanFunctionID(serviceModel.ID) {
		case registry.Kpm:
			kpmSm, err := kpm.NewServiceModel(node, model, modelPluginRegistry,
				subStore, nodeStore, ueStore)
			if err != nil {
				return nil, err
			}
			err = reg.RegisterServiceModel(kpmSm)
			if err != nil {
				log.Error(err)
				return nil, err
			}
		case registry.Rc:
			rcSm, err := rc.NewServiceModel(node, model, modelPluginRegistry,
				subStore, nodeStore, ueStore, cellStore, metricStore)
			if err != nil {
				return nil, err
			}
			err = reg.RegisterServiceModel(rcSm)
			if err != nil {
				log.Error(err)
				return nil, err
			}
		case registry.Kpm2:
			log.Info("KPM2 service model for node with eNbID:", node.GnbID)
			kpm2Sm, err := kpm2.NewServiceModel(node, model, modelPluginRegistry,
				subStore, nodeStore, ueStore)
			if err != nil {
				log.Info("Failure creating KPM2 service model for eNbID:", node.GnbID)
				return nil, err
			}
			err = reg.RegisterServiceModel(kpm2Sm)
			if err != nil {
				log.Info("Failure registering KPM2 service model for eNbID:", node.GnbID)
				log.Error(err)
				return nil, err
			}
		case registry.Mho:
			log.Info("MHO service model for node with eNbID:", node.GnbID)
			mhoSm, err := mho.NewServiceModel(node, model, modelPluginRegistry, subStore, nodeStore, ueStore, cellStore,
				metricStore, measChan, mobilityDriver)
			if err != nil {
				log.Info("Failure creating MHO service model for eNbID:", node.GnbID)
				return nil, err
			}
			err = reg.RegisterServiceModel(mhoSm)
			if err != nil {
				log.Info("Failure registering MHO service model for eNbID:", node.GnbID)
				log.Error(err)
				return nil, err
			}
		}
	}
	return &e2Agent{
		node:      node,
		registry:  reg,
		model:     model,
		subStore:  subStore,
		nodeStore: nodeStore,
		ueStore:   ueStore,
		cellStore: cellStore,
	}, nil
}

func (a *e2Agent) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	ranFuncID := registry.RanFunctionID(controlutils.GetRanFunctionID(request))
	log.Debugf("Received Control Request %+v for ran function %d", request, ranFuncID)
	sm, err := a.registry.GetServiceModel(ranFuncID)
	if err != nil {
		log.Warn(err)
		// TODO If the target E2 Node receives a RIC CONTROL REQUEST message
		//  which contains a RAN Function ID IE that was not previously announced as a s
		//  supported RAN function in the E2 Setup procedure or the RIC Service Update procedure,
		//  or the E2 Node does not support the specific RIC Control procedure action, then
		//  the target E2 Node shall ignore message and send an ERROR INDICATION message to the Near-RT RIC.

		return nil, nil, err
	}
	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		response, failure, err = client.RICControl(ctx, request)
	case registry.Rc:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICControl(ctx, request)
	case registry.Mho:
		client := sm.Client.(*mho.Mho)
		response, failure, err = client.RICControl(ctx, request)
	}
	if err != nil {
		return nil, nil, err
	}

	return response, failure, err
}

func (a *e2Agent) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	ranFuncID := registry.RanFunctionID(subutils.GetRanFunctionID(request))
	log.Debugf("Received Subscription Request %v for ran function %d", request, ranFuncID)
	sm, err := a.registry.GetServiceModel(ranFuncID)
	id := subscriptions.NewID(subutils.GetRicInstanceID(request),
		subutils.GetRequesterID(request),
		subutils.GetRanFunctionID(request))

	if err != nil {
		log.Warn(err)
		// If the target E2 Node receives a RIC SUBSCRIPTION REQUEST
		//  message which contains a RAN Function ID IE that was not previously
		//  announced as a supported RAN function in the E2 Setup procedure or
		//  the RIC Service Update procedure, the target E2 Node shall send the RIC SUBSCRIPTION FAILURE message
		//  to the Near-RT RIC with an appropriate cause value.
		var ricActionsAccepted []*e2aptypes.RicActionID
		ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
		actionList := subutils.GetRicActionToBeSetupList(request)
		reqID := subutils.GetRequesterID(request)
		ranFuncID := subutils.GetRanFunctionID(request)
		ricInstanceID := subutils.GetRicInstanceID(request)

		for _, action := range actionList {
			actionID := e2aptypes.RicActionID(action.Value.RicActionId.Value)
			cause := &e2apies.Cause{
				Cause: &e2apies.Cause_RicRequest{
					RicRequest: e2apies.CauseRic_CAUSE_RIC_RAN_FUNCTION_ID_INVALID,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(reqID),
			subutils.WithRanFuncID(ranFuncID),
			subutils.WithRicInstanceID(ricInstanceID),
			subutils.WithActionsAccepted(ricActionsAccepted),
			subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))
		failure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}
	subscription, err := subscriptions.NewSubscription(id, request, a.channel)
	if err != nil {
		return response, failure, err
	}
	err = a.subStore.Add(subscription)
	if err != nil {
		return response, failure, err
	}

	// TODO - Assumes ono-to-one mapping between ran function and server model
	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Rc:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Kpm2:
		client := sm.Client.(*kpm2.Client)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Mho:
		client := sm.Client.(*mho.Mho)
		response, failure, err = client.RICSubscription(ctx, request)

	}
	// Ric subscription is failed
	if err != nil {
		return response, failure, err
	}

	return response, failure, err
}

func (a *e2Agent) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	ranFuncID := registry.RanFunctionID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	log.Debugf("Received Subscription Delete Request %v for ran function ID %d", request, ranFuncID)
	subID := subscriptions.NewID(subdeleteutils.GetRicInstanceID(request),
		subdeleteutils.GetRequesterID(request),
		subdeleteutils.GetRanFunctionID(request))
	_, err = a.subStore.Get(subID)
	if err != nil {
		log.Warn(err)
		//  If the target E2 Node receives a RIC SUBSCRIPTION DELETE REQUEST
		//  message containing RIC Request ID IE that is not known, the target
		//  E2 Node shall send the RIC SUBSCRIPTION DELETE FAILURE message
		//  to the Near-RT RIC. The message shall contain the Cause IE with an appropriate value.
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRic_CAUSE_RIC_REQUEST_ID_UNKNOWN,
			},
		}
		subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
			subdeleteutils.WithRanFuncID(subdeleteutils.GetRanFunctionID(request)),
			subdeleteutils.WithRequestID(subdeleteutils.GetRequesterID(request)),
			subdeleteutils.WithRicInstanceID(subdeleteutils.GetRicInstanceID(request)),
			subdeleteutils.WithCause(cause))
		failure, err := subscriptionDelete.BuildSubscriptionDeleteFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil

	}

	sm, err := a.registry.GetServiceModel(ranFuncID)
	if err != nil {
		log.Warn(err)
		//  If the target E2 Node receives a RIC SUBSCRIPTION DELETE REQUEST message contains a
		//  RAN Function ID IE that was not previously announced as a supported RAN function
		//  in the E2 Setup procedure or the RIC Service Update procedure, the target E2 Node
		//  shall send the RIC SUBSCRIPTION DELETE FAILURE message to the Near-RT RIC.
		//  The message shall contain with an appropriate cause value.
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRic_CAUSE_RIC_RAN_FUNCTION_ID_INVALID,
			},
		}
		subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
			subdeleteutils.WithRanFuncID(subdeleteutils.GetRanFunctionID(request)),
			subdeleteutils.WithRequestID(subdeleteutils.GetRequesterID(request)),
			subdeleteutils.WithRicInstanceID(subdeleteutils.GetRicInstanceID(request)),
			subdeleteutils.WithCause(cause))
		failure, err := subscriptionDelete.BuildSubscriptionDeleteFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	case registry.Rc:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	case registry.Kpm2:
		client := sm.Client.(*kpm2.Client)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)

	}
	// Ric subscription delete procedure is failed so we are not going to update subscriptions store
	if err != nil {
		log.Warn(err)
		return response, failure, err
	}

	err = a.subStore.Remove(subID)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	return response, failure, err
}

func (a *e2Agent) Start() error {
	if len(a.node.Controllers) == 0 {
		return errors.New(errors.Invalid, "no controller is associated with this node")
	}

	log.Infof("E2 node %d is starting; attempting to connect", a.node.GnbID)
	b := newExpBackoff()

	// Attempt to connect to the E2T controller; use exponential back-off retry
	count := 0
	connectNotify := func(err error, t time.Duration) {
		count++
		log.Infof("E2 node %d failed to connect; retry after %v; attempt %d", a.node.GnbID, b.GetElapsedTime(), count)
	}

	err := backoff.RetryNotify(a.connect, b, connectNotify)
	if err != nil {
		return err
	}
	log.Infof("E2 node %d connected; attempting setup", a.node.GnbID)

	// Attempt to negotiate E2 setup procedure; use exponential back-off retry
	count = 0
	setupNotify := func(err error, t time.Duration) {
		count++
		log.Infof("E2 node %d failed setup procedure; retry after %v; attempt %d", a.node.GnbID, b.GetElapsedTime(), count)
	}

	err = backoff.RetryNotify(a.setup, b, setupNotify)
	log.Infof("E2 node %d completed connection setup", a.node.GnbID)
	return err
}

func (a *e2Agent) connect() error {
	controller, err := a.model.GetController(a.node.Controllers[0])
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", controller.Address, controller.Port)
	channel, err := e2.Connect(context.TODO(), addr,
		func(channel e2.ClientChannel) e2.ClientInterface {
			return a
		},
	)

	if err != nil {
		return err
	}
	a.channel = channel
	return nil
}

func (a *e2Agent) setup() error {
	plmnID := ransimtypes.NewUint24(uint32(a.model.PlmnID))
	setupRequest := setup.NewSetupRequest(
		setup.WithRanFunctions(a.registry.GetRanFunctions()),
		setup.WithPlmnID(plmnID.Value()),
		setup.WithE2NodeID(uint64(a.node.GnbID)))

	e2SetupRequest, err := setupRequest.Build()

	if err != nil {
		log.Error(err)
		return err
	}
	_, e2SetupFailure, err := a.channel.E2Setup(context.Background(), e2SetupRequest)
	if err != nil {
		log.Error(err)
		return errors.NewUnknown("E2 setup failed: %v", err)
	} else if e2SetupFailure != nil {
		err := errors.NewInvalid("E2 setup failed")
		log.Error(err)
		return err
	}
	return nil
}

func (a *e2Agent) Stop() error {
	log.Debugf("Stopping e2 agent with ID %d:", a.node.GnbID)

	if a.channel != nil {
		return a.channel.Close()
	}
	return nil
}

var _ E2Agent = &e2Agent{}

var _ e2.ClientInterface = &e2Agent{}
