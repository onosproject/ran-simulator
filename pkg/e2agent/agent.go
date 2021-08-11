// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"net"

	"github.com/onosproject/ran-simulator/pkg/mobility"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/mho"
	"github.com/onosproject/ran-simulator/pkg/store/channels"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm2"

	"github.com/onosproject/ran-simulator/pkg/store/metrics"

	"github.com/onosproject/ran-simulator/pkg/store/cells"

	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/rc"

	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm"
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
	node         model.Node
	model        *model.Model
	registry     *registry.ServiceModelRegistry
	subStore     *subscriptions.Subscriptions
	nodeStore    nodes.Store
	ueStore      ues.Store
	cellStore    cells.Store
	channelStore channels.Store
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
		case registry.Rcpre2:
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

func (a *e2Agent) Start() error {
	if len(a.node.Controllers) == 0 {
		return errors.NewInvalid("no controller is associated with this node")
	}
	controller, err := a.model.GetController(a.node.Controllers[0])
	if err != nil {
		return err
	}

	controllerAddresses, err := net.LookupHost(controller.Address)
	if err != nil {
		return err
	}
	ricAddress := RICAddress{
		ipAddress: net.ParseIP(controllerAddresses[0]),
		port:      uint64(controller.Port),
	}
	channelStore := channels.NewStore()
	a.channelStore = channelStore

	e2Instance := NewE2Instance(WithNode(a.node),
		WithModel(a.model),
		WithSMRegistry(a.registry),
		WithSubStore(a.subStore),
		WithRICAddress(ricAddress),
		WithChannelStore(channelStore))

	err = e2Instance.Start()
	if err != nil {
		return err
	}
	return nil
}

func (a *e2Agent) Stop() error {
	log.Debugf("Stopping e2 agent with ID %d:", a.node.GnbID)
	return nil
}

var _ E2Agent = &e2Agent{}
