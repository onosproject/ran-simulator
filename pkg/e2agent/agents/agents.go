// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package agents

import (
	"context"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/e2agent"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/store/agents"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
)

var log = logging.GetLogger("e2agent", "agents")

// E2Agents represents a collection of E2 agents to allow centralized management
type E2Agents struct {
	agentStore          agents.Store
	modelPluginRegistry *modelplugins.ModelPluginRegistry
	nodeStore           nodes.Store
	ueStore             ues.Store
	model               *model.Model
}

// Agents agents interface
type Agents interface {
	Start() error

	Stop() error
}

func (agents *E2Agents) processNodeEvents() {
	ch := make(chan event.Event)
	err := agents.nodeStore.Watch(context.Background(), ch)
	if err != nil {
		log.Error(err)
	}
	for nodeEvent := range ch {
		switch nodeEvent.Type {
		case nodes.Created:
			// TODO start e2 agent
			log.Debugf("New node is added: %v", nodeEvent.Key)
		case nodes.Deleted:
			// TODO stop e2agent
			log.Debugf("Node %v is deleted: %v", nodeEvent.Key)
		case nodes.Updated:
			log.Debugf("Node %v is updated: %v", nodeEvent.Key)

		}
	}
}

// NewE2Agents creates a new collection of E2 agents from the specified list of nodes
func NewE2Agents(m *model.Model, modelPluginRegistry *modelplugins.ModelPluginRegistry,
	nodeStore nodes.Store, ueStore ues.Store) (*E2Agents, error) {
	agentStore := agents.NewStore()
	e2agents := &E2Agents{
		agentStore:          agentStore,
		nodeStore:           nodeStore,
		modelPluginRegistry: modelPluginRegistry,
		model:               m,
		ueStore:             ueStore,
	}

	for _, node := range m.Nodes {
		e2Node, err := e2agent.NewE2Agent(node, m, modelPluginRegistry, nodeStore, ueStore)
		if err != nil {
			return nil, err
		}
		err = agentStore.Add(node.EnbID, e2Node)
		if err != nil {
			log.Error(err)
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
