// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/subscriptiondelete"

	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"

	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/utils/setup"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/subscription"

	"github.com/cenkalti/backoff"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm"

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
	log.Info("Creating New E2 Agent for node")
	reg := registry.NewServiceModelRegistry()
	sms := node.ServiceModels
	for _, smID := range sms {
		serviceModel, err := model.GetServiceModel(smID)
		if err != nil {
			return nil, err
		}
		switch registry.RanFunctionID(serviceModel.ID) {
		case registry.Kpm:
			sm, err := kpm.NewServiceModel(modelPluginRegistry)
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
		client.Channel = a.channel
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
	if err != nil {
		return nil, nil, err
	}

	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		client.Channel = a.channel
		client.ServiceModel = &sm
		response, failure, err = client.RICSubscription(ctx, request)
	default:
		return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

	}
	// Ric subscription is failed so we are not going to update the store
	if err != nil {
		return response, failure, err
	}
	err = a.subStore.Add(subscriptions.NewSubscription(request))
	if err != nil {
		return response, failure, err
	}
	return response, failure, err
}

func (a *e2Agent) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Debugf("Received Subscription Delete Request %v", request)
	ranFuncID := registry.RanFunctionID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)

	sm, err := a.registry.GetServiceModel(ranFuncID)
	if err != nil {
		return nil, nil, err
	}

	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		client.Channel = a.channel
		client.ServiceModel = &sm
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	default:
		return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

	}
	// Ric subscription delete procedure is failed so we are not going to update subscriptions store
	if err != nil {
		return response, failure, err
	}

	id := subscriptions.NewID(subdeleteutils.GetRicInstanceID(request),
		subdeleteutils.GetRequesterID(request),
		subdeleteutils.GetRanFunctionID(request))
	err = a.subStore.Remove(id)
	if err != nil {
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

	log.Infof("%s is starting; attempting to connect", a.node.Ecgi)
	b := newExpBackoff()

	// Attempt to connect to the E2T controller; use exponential back-off retry
	count := 0
	connectNotify := func(err error, t time.Duration) {
		count++
		log.Infof("%s failed to connect; retry after %v; attempt %d", a.node.Ecgi, b.GetElapsedTime(), count)
	}

	log.Infof("%s connected; attempting setup", a.node.Ecgi)

	err := backoff.RetryNotify(a.connect, b, connectNotify)
	if err != nil {
		return err
	}

	// Attempt to negotiate E2 setup procedure; use exponential back-off retry
	count = 0
	setupNotify := func(err error, t time.Duration) {
		count++
		log.Infof("%s failed setup procedure; retry after %v; attempt %d", a.node.Ecgi, b.GetElapsedTime(), count)
	}

	err = backoff.RetryNotify(a.setup, b, setupNotify)

	log.Infof("%s completed connection setup", a.node.Ecgi)
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
	setupRequest, err := setup.NewSetupRequest(
		setup.WithRanFunctions(a.registry.GetRanFunctions()),
		setup.WithPlmnID("onf"),
		setup.WithE2NodeID(nodeID(a.node.Ecgi)))

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

func nodeID(ecgi model.Ecgi) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(ecgi))
	return h.Sum64()
}

func (a *e2Agent) Stop() error {
	if a.channel != nil {
		return a.channel.Close()
	}
	return nil
}

var _ E2Agent = &e2Agent{}

var _ e2.ClientInterface = &e2Agent{}
