// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/store/agents"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
)

// E2Agents represents a collection of E2 agents to allow centralized management
type E2Agents struct {
	agentStore agents.Store
}

// NewE2Agents creates a new collection of E2 agents from the specified list of nodes
func NewE2Agents(m *model.Model, modelPluginRegistry *modelplugins.ModelPluginRegistry,
	nodeStore nodes.NodeRegistry, ueStore ues.UERegistry) (*E2Agents, error) {
	agentStore := agents.NewStore()
	e2agents := &E2Agents{
		agentStore: agentStore,
	}

	for _, node := range m.Nodes {
		e2Node, err := NewE2Agent(node, m, modelPluginRegistry, nodeStore, ueStore)
		if err != nil {
			return nil, err
		}
		err = agentStore.Add(node.EnbID, e2Node)
		if err != nil {
			log.Error(err)
		}
	}
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
