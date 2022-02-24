// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package agents

import (
	"context"

	"github.com/onosproject/rrm-son-lib/pkg/handover"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/mobility"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"

	"github.com/onosproject/ran-simulator/pkg/store/cells"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/e2agent"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/agents"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
)

var log = logging.GetLogger()

// E2Agents represents a collection of E2 agents to allow centralized management
type E2Agents struct {
	agentStore     agents.Store
	nodeStore      nodes.Store
	ueStore        ues.Store
	cellStore      cells.Store
	metricStore    metrics.Store
	model          *model.Model
	a3Chan         chan handover.A3HandoverDecision
	mobilityDriver mobility.Driver
}

// Agents agents interface
type Agents interface {
	Start() error

	Stop() error
}

func (agents *E2Agents) processNodeEvents() {
	ch := make(chan event.Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := agents.nodeStore.Watch(ctx, ch)
	if err != nil {
		log.Error(err)
	}
	for nodeEvent := range ch {
		log.Debug("Received Node event:", nodeEvent)
		switch nodeEvent.Type {
		case nodes.Created:
			node := nodeEvent.Value.(*model.Node)
			log.Debugf("Starting e2 agent %d", nodeEvent.Key.(types.GnbID))
			e2Node, err := e2agent.NewE2Agent(*node, agents.model, agents.nodeStore, agents.ueStore,
				agents.cellStore, agents.metricStore, agents.a3Chan, agents.mobilityDriver)
			if err != nil {
				log.Error(err)
				continue
			}
			err = agents.agentStore.Add(node.GnbID, e2Node)
			if err != nil {
				log.Error(err)
			}

			err = e2Node.Start()
			if err != nil {
				log.Error(err)
				err = agents.agentStore.Remove(node.GnbID)
				if err != nil {
					log.Error(err)
				}
			}
			err = agents.nodeStore.SetStatus(context.Background(), node.GnbID, "Running")
			if err != nil {
				log.Error(err)
			}

		case nodes.Deleted:
			log.Debugf("Stopping e2 agent %d", nodeEvent.Key.(types.GnbID))
			node := nodeEvent.Value.(*model.Node)
			e2Node, err := agents.agentStore.Get(node.GnbID)
			if err != nil {
				log.Error(err)
				continue
			}
			err = e2Node.Stop()
			if err != nil {
				log.Error(err)
				continue
			}

			err = agents.nodeStore.SetStatus(context.Background(), node.GnbID, "Stopped")
			if err != nil {
				log.Error(err)
			}

			err = agents.agentStore.Remove(node.GnbID)
			if err != nil {
				log.Error(err)
			}

		}
	}
}

// NewE2Agents creates a new collection of E2 agents from the specified list of nodes
func NewE2Agents(m *model.Model,
	nodeStore nodes.Store, ueStore ues.Store, cellStore cells.Store, metricStore metrics.Store,
	a3Chan chan handover.A3HandoverDecision, mobilityDriver mobility.Driver) (*E2Agents, error) {
	agentStore := agents.NewStore()
	e2agents := &E2Agents{
		agentStore:     agentStore,
		nodeStore:      nodeStore,
		model:          m,
		ueStore:        ueStore,
		cellStore:      cellStore,
		metricStore:    metricStore,
		a3Chan:         a3Chan,
		mobilityDriver: mobilityDriver,
	}

	for _, node := range m.Nodes {
		e2Node, err := e2agent.NewE2Agent(node, m, nodeStore, ueStore, cellStore, metricStore, a3Chan, mobilityDriver)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		err = agentStore.Add(node.GnbID, e2Node)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = nodeStore.SetStatus(context.Background(), node.GnbID, "Running")
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	go e2agents.processNodeEvents()
	return e2agents, nil
}

// Start all simulated node agents
func (agents *E2Agents) Start() error {
	log.Info("Starting E2 Agents")
	agentList, err := agents.agentStore.List()
	if err != nil {
		log.Error(err)
		return err
	}
	for id, agent := range agentList {
		log.Debug("Starting agent with e2 node ID:", id)
		err := agent.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop all simulated node agents
func (agents *E2Agents) Stop() error {
	log.Info("Stopping E2 Agents")
	agentList, err := agents.agentStore.List()
	if err != nil {
		log.Error(err)
		return err
	}
	for id, agent := range agentList {
		log.Debug("Stopping agent with e2 node ID:", id)
		err := agent.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

var _ Agents = &E2Agents{}
