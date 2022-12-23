// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package e2agent

import (
	"context"
	"net"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm2"

	"github.com/onosproject/ran-simulator/pkg/e2agent/addressing"

	"github.com/onosproject/ran-simulator/pkg/e2agent/connection"

	"github.com/onosproject/ran-simulator/pkg/mobility"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/mho"
	"github.com/onosproject/ran-simulator/pkg/store/connections"
	"github.com/onosproject/rrm-son-lib/pkg/handover"

	"github.com/onosproject/ran-simulator/pkg/store/metrics"

	"github.com/onosproject/ran-simulator/pkg/store/cells"

	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/rc"
	rcv1 "github.com/onosproject/ran-simulator/pkg/servicemodel/rc/v1"

	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	connectionController "github.com/onosproject/ran-simulator/pkg/controller/connection"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

var log = logging.GetLogger()

// E2Agent is an E2 agent
type E2Agent interface {
	// Start starts the agent
	Start() error

	// Stop stops the agent
	Stop() error
}

// e2Agent is an E2 agent
type e2Agent struct {
	node            model.Node
	model           *model.Model
	registry        *registry.ServiceModelRegistry
	subStore        *subscriptions.Subscriptions
	nodeStore       nodes.Store
	ueStore         ues.Store
	cellStore       cells.Store
	connectionStore connections.Store
}

// NewE2Agent creates a new E2 agent
func NewE2Agent(node model.Node, model *model.Model,
	nodeStore nodes.Store, ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store,
	a3Chan chan handover.A3HandoverDecision, mobilityDriver mobility.Driver) (E2Agent, error) {
	log.Info("Creating New E2 Agent for node with e2 Node ID:", node.GnbID)
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
		case registry.Rcpre2:
			log.Infof("Registering RC PRE service model for node with e2 Node ID: %v", node.GnbID)
			rcSm, err := rc.NewServiceModel(node, model,
				subStore, nodeStore, ueStore, cellStore, metricStore)
			if err != nil {
				log.Errorf("Failure creating RC PRE service model for e2 node ID: %v, %s", node.GnbID, err.Error())
				return nil, err
			}
			err = reg.RegisterServiceModel(rcSm)
			if err != nil {
				log.Errorf("Failure registering RC PRE service model for e2 Node ID: %v, %s", node.GnbID, err.Error())
				return nil, err
			}
		case registry.Kpm2:
			log.Infof("Registering KPM2 service model for node with e2 Node ID: %v", node.GnbID)
			kpm2Sm, err := kpm2.NewServiceModel(node, model,
				subStore, nodeStore, ueStore)
			if err != nil {
				log.Errorf("Failure creating KPM2 service model for e2 node ID: %v, %s", node.GnbID, err.Error())
				return nil, err
			}
			err = reg.RegisterServiceModel(kpm2Sm)
			if err != nil {
				log.Errorf("Failure registering KPM2 service model for e2 Node ID: %v, %s", node.GnbID, err.Error())
				return nil, err
			}
		case registry.Mho:
			log.Infof("Registering MHO service model for node with e2 Node ID: %v", node.GnbID)
			mhoSm, err := mho.NewServiceModel(node, model, subStore, nodeStore, ueStore, cellStore,
				metricStore, a3Chan, mobilityDriver)
			if err != nil {
				log.Errorf("Failure creating MHO service model for e2 Node ID: %v, %s", node.GnbID, err.Error())
				return nil, err
			}
			err = reg.RegisterServiceModel(mhoSm)
			if err != nil {
				log.Errorf("Failure registering MHO service model for e2 Node ID: %s, %s", node.GnbID, err.Error())
				return nil, err
			}
		case registry.Rc:
			log.Infof("Registering RC service model for e2 node ID:%v", node.GnbID)
			rcv1Sm, err := rcv1.NewServiceModel(node, model, subStore, nodeStore, ueStore, cellStore, metricStore,
				a3Chan, mobilityDriver)
			if err != nil {
				log.Errorf("Failure creating RC service model for e2 Node ID: %v, %s", node.GnbID, err.Error())
				return nil, err
			}
			err = reg.RegisterServiceModel(rcv1Sm)
			if err != nil {
				log.Errorf("Failure registering RC service model for e2 Node ID: %v, %s", node.GnbID, err.Error())
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
	ricAddress := addressing.RICAddress{
		IPAddress: net.ParseIP(controllerAddresses[0]),
		Port:      uint64(controller.Port),
	}
	connectionStore := connections.NewStore()
	a.connectionStore = connectionStore

	c := connectionController.NewController(connectionStore, a.node, a.model, a.registry, a.subStore, a.cellStore)
	err = c.Start()
	if err != nil {
		return err
	}

	e2Connection := connection.NewE2Connection(connection.WithNode(a.node),
		connection.WithModel(a.model),
		connection.WithSMRegistry(a.registry),
		connection.WithSubStore(a.subStore),
		connection.WithRICAddress(ricAddress),
		connection.WithConnectionStore(connectionStore),
		connection.WithCellStore(a.cellStore))

	err = e2Connection.Setup()
	if err != nil {
		return err
	}
	return nil
}

func (a *e2Agent) Stop() error {
	log.Debugf("Stopping e2 agent with ID %d:", a.node.GnbID)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conns := a.connectionStore.List(context.Background())
	log.Debugf("List of Connections: %+v", conns)
	for _, conn := range conns {
		if conn.Client != nil {
			log.Debugf("Closing connection: %+v", conn.ID)
			err := conn.Client.Close()
			if err != nil {
				return err
			}
			err = a.connectionStore.Remove(ctx, conn.ID)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

var _ E2Agent = &e2Agent{}
