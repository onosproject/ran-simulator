// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"

	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"

	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"

	"github.com/cenkalti/backoff"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm"
	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/setup"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

var log = logging.GetLogger("agent")

const (
	backoffInterval = 10 * time.Millisecond
	maxBackoffTime  = 5 * time.Second
)

// E2Agent is an E2 agent
type E2Agent interface {
	// Start starts the agent
	Start() error

	// Stop stops the agent
	Stop() error
}

// NewE2Agent creates a new E2 agent
func NewE2Agent(node model.Node, model *model.Model, modelPluginRegistry *modelplugins.ModelPluginRegistry) (E2Agent, error) {
	log.Info("Creating New E2 Agent for node with eNbID:", node.EnbID)
	reg := registry.NewServiceModelRegistry()
	sms := node.ServiceModels
	for _, smID := range sms {
		serviceModel, err := model.GetServiceModel(smID)
		if err != nil {
			return nil, err
		}
		switch registry.RanFunctionID(serviceModel.ID) {
		case registry.Kpm:
			sm, err := kpm.NewServiceModel(node, model, modelPluginRegistry)
			if err != nil {
				return nil, err
			}
			err = reg.RegisterServiceModel(sm)
			if err != nil {
				log.Error(err)
				return nil, err
			}
		}
	}

	return &e2Agent{
		node:     node,
		registry: reg,
		model:    model,
		subStore: subscriptions.NewStore(),
	}, nil
}

// e2Agent is an E2 agent
type e2Agent struct {
	node     model.Node
	model    *model.Model
	channel  e2.ClientChannel
	registry *registry.ServiceModelRegistry
	subStore *subscriptions.Subscriptions
}

func (a *e2Agent) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	ranFuncID := registry.RanFunctionID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	sm, err := a.registry.GetServiceModel(ranFuncID)
	if err != nil {
		return nil, nil, err
	}
	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		client.ServiceModel = &sm
		response, failure, err = client.RICControl(ctx, request)
	default:
		return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

	}
	return response, failure, err
}

func (a *e2Agent) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Debugf("Received Subscription Request %v", request)
	ranFuncID := registry.RanFunctionID(subutils.GetRanFunctionID(request))
	sm, err := a.registry.GetServiceModel(ranFuncID)
	id := subscriptions.NewID(request.ProtocolIes.E2ApProtocolIes29.Value.RicInstanceId,
		request.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId,
		request.ProtocolIes.E2ApProtocolIes5.Value.Value)

	if err != nil {
		// If the target E2 Node receives a RIC SUBSCRIPTION REQUEST
		//  message which contains a RAN Function ID IE that was not previously
		//  announced as a supported RAN function in the E2 Setup procedure or
		//  the RIC Service Update procedure, the target E2 Node shall send the RIC SUBSCRIPTION FAILURE message
		//  to the Near-RT RIC with an appropriate cause value.
		var ricActionsAccepted []*types.RicActionID
		ricActionsNotAdmitted := make(map[types.RicActionID]*e2apies.Cause)
		actionList := subutils.GetRicActionToBeSetupList(request)
		reqID := subutils.GetRequesterID(request)
		ranFuncID := subutils.GetRanFunctionID(request)
		ricInstanceID := subutils.GetRicInstanceID(request)

		for _, action := range actionList {
			actionID := types.RicActionID(action.Value.RicActionId.Value)
			cause := &e2apies.Cause{
				Cause: &e2apies.Cause_RicRequest{
					RicRequest: e2apies.CauseRic_CAUSE_RIC_RAN_FUNCTION_ID_INVALID,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
		subscription, _ := subutils.NewSubscription(
			subutils.WithRequestID(reqID),
			subutils.WithRanFuncID(ranFuncID),
			subutils.WithRicInstanceID(ricInstanceID),
			subutils.WithActionsAccepted(ricActionsAccepted),
			subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))
		failure := subutils.CreateSubscriptionFailure(subscription)
		return nil, failure, err
	}
	subscription, err := subscriptions.NewSubscription(id, request, a.channel)
	if err != nil {
		return response, failure, err
	}
	err = a.subStore.Add(subscription)
	if err != nil {
		return response, failure, err
	}

	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		client.Subscriptions = a.subStore
		client.ServiceModel = &sm
		response, failure, err = client.RICSubscription(ctx, request)
	}
	// Ric subscription is failed so we are not going to update the store
	if err != nil {
		return response, failure, err
	}

	return response, failure, err
}

func (a *e2Agent) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Debugf("Received Subscription Delete Request %v", request)
	ranFuncID := registry.RanFunctionID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	subID := subscriptions.NewID(subdeleteutils.GetRicInstanceID(request),
		subdeleteutils.GetRequesterID(request),
		subdeleteutils.GetRanFunctionID(request))

	sm, err := a.registry.GetServiceModel(ranFuncID)
	if err != nil {
		log.Error(err)
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
		subscriptionDelete, _ := subdeleteutils.NewSubscriptionDelete(
			subdeleteutils.WithRanFuncID(subdeleteutils.GetRanFunctionID(request)),
			subdeleteutils.WithRequestID(subdeleteutils.GetRequesterID(request)),
			subdeleteutils.WithRicInstanceID(subdeleteutils.GetRicInstanceID(request)),
			subdeleteutils.WithCause(cause))
		failure := subdeleteutils.CreateSubscriptionDeleteFailure(subscriptionDelete)
		return nil, failure, err
	}

	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		client.Subscriptions = a.subStore
		client.ServiceModel = &sm
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	}
	// Ric subscription delete procedure is failed so we are not going to update subscriptions store
	if err != nil {
		return response, failure, err
	}

	err = a.subStore.Remove(subID)
	log.Info("Remove from store:", subID, a.subStore)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	return response, failure, err
}

func newExpBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = backoffInterval
	// MaxInterval caps the RetryInterval
	b.MaxInterval = maxBackoffTime
	// Never stops retrying
	b.MaxElapsedTime = 0
	return b
}

func (a *e2Agent) Start() error {
	if len(a.node.Controllers) == 0 {
		return errors.New(errors.Invalid, "no controller is associated with this node")
	}

	log.Infof("%s is starting; attempting to connect", a.node.EnbID)
	b := newExpBackoff()

	// Attempt to connect to the E2T controller; use exponential back-off retry
	count := 0
	connectNotify := func(err error, t time.Duration) {
		count++
		log.Infof("%s failed to connect; retry after %v; attempt %d", a.node.EnbID, b.GetElapsedTime(), count)
	}

	err := backoff.RetryNotify(a.connect, b, connectNotify)
	if err != nil {
		return err
	}
	log.Infof("%s connected; attempting setup", a.node.EnbID)

	// Attempt to negotiate E2 setup procedure; use exponential back-off retry
	count = 0
	setupNotify := func(err error, t time.Duration) {
		count++
		log.Infof("%s failed setup procedure; retry after %v; attempt %d", a.node.EnbID, b.GetElapsedTime(), count)
	}

	err = backoff.RetryNotify(a.setup, b, setupNotify)

	log.Infof("%s completed connection setup", a.node.EnbID)
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
	e2GlobalID, err := nodeID(a.model.PlmnID, a.node.EnbID)
	if err != nil {
		return err
	}
	setupRequest, err := setup.NewSetupRequest(
		setup.WithRanFunctions(a.registry.GetRanFunctions()),
		setup.WithPlmnID(string(a.model.PlmnID)),
		setup.WithE2NodeID(e2GlobalID))

	if err != nil {
		return err
	}

	e2SetupRequest := setup.CreateSetupRequest(setupRequest)
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

func nodeID(plmndID model.PlmnID, enbID model.EnbID) (uint64, error) {
	gEnbID := model.GEnbID{
		PlmnID: plmndID,
		EnbID:  enbID,
	}
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(gEnbID)
	if err != nil {
		return 0, err
	}

	h := fnv.New64a()
	_, _ = h.Write(buf.Bytes())
	return h.Sum64(), nil
}

func (a *e2Agent) Stop() error {
	if a.channel != nil {
		return a.channel.Close()
	}
	return nil
}

var _ E2Agent = &e2Agent{}

var _ e2.ClientInterface = &e2Agent{}
