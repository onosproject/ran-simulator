// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"context"
	"fmt"

	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/utils/setup"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

var log = logging.GetLogger("agent")

// E2Agent is an E2 agent
type E2Agent interface {
	// Start starts the agent
	Start() error

	// Stop stops the agent
	Stop() error
}

// NewE2Agent creates a new E2 agent
func NewE2Agent(node model.Node, model *model.Model) (E2Agent, error) {
	reg := registry.NewServiceModelRegistry()
	sms := node.ServiceModels
	for _, smID := range sms {
		serviceModel, err := model.GetServiceModel(smID)
		if err != nil {
			return nil, err
		}
		switch registry.ID(serviceModel.ID) {
		case registry.Kpm:
			err := reg.RegisterServiceModel(kpm.GetConfig())
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
	}, nil
}

// e2Agent is an E2 agent
type e2Agent struct {
	node     model.Node
	model    *model.Model
	channel  e2.ClientChannel
	registry *registry.ServiceModelRegistry
}

func (a *e2Agent) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	ranFuncID := registry.ID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	switch ranFuncID {
	case registry.Kpm:
		var kpmService kpm.ServiceModel
		err = a.registry.GetServiceModel(ranFuncID, &kpmService)
		if err != nil {
			return nil, nil, err
		}
		return kpmService.RICControl(ctx, request)

	}
	return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

}

func (a *e2Agent) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	ranFuncID := registry.ID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	switch ranFuncID {
	case registry.Kpm:
		var kpmService kpm.ServiceModel
		kpmService.Channel = a.channel
		err = a.registry.GetServiceModel(ranFuncID, &kpmService)
		if err != nil {
			return nil, nil, err
		}
		return kpmService.RICSubscription(ctx, request)

	}
	return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

}

func (a *e2Agent) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	ranFuncID := registry.ID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	switch ranFuncID {
	case registry.Kpm:
		var kpmService kpm.ServiceModel
		err = a.registry.GetServiceModel(ranFuncID, &kpmService)
		if err != nil {
			return nil, nil, err
		}
		return kpmService.RICSubscriptionDelete(ctx, request)

	}
	return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

}

func (a *e2Agent) Start() error {
	if len(a.node.Controllers) == 0 {
		return errors.New(errors.Invalid, "no controller is associated with this node")
	}

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

	setupRequest, err := setup.NewSetupRequest(
		setup.WithRanFunctions(a.registry.GetRanFunctions()),
		setup.WithPlmnID("onf"))

	if err != nil {
		return err
	}

	e2SetupRequest := setup.CreateSetupRequest(setupRequest)
	_, e2SetupFailure, err := channel.E2Setup(context.Background(), e2SetupRequest)
	if err != nil {
		log.Error(err)
		return errors.NewUnknown("E2 setup failed: %v", err)
	} else if e2SetupFailure != nil {
		err := errors.NewInvalid("E2 setup failed")
		log.Error(err)
		return err
	}

	a.channel = channel
	return nil
}

func (a *e2Agent) Stop() error {
	if a.channel != nil {
		return a.channel.Close()
	}
	return nil
}

var _ E2Agent = &e2Agent{}

var _ e2.ClientInterface = &e2Agent{}
