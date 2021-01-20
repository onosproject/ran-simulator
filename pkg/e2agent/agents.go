// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

// E2Agents represents a collection of E2 agents to allow centralized management
type E2Agents struct {
	Agents map[model.ECGI]E2Agent
}

// NewE2Agents creates a new collection of E2 agents from the specified list of nodes
func NewE2Agents(nodes []*model.SimNode, reg *registry.ServiceModelRegistry, controllers []*model.Controller) *E2Agents {
	agents := &E2Agents{
		Agents: make(map[model.ECGI]E2Agent),
	}

	for _, node := range nodes {
		agents.Agents[node.ECGI] = NewE2Agent(node, reg, controllers)
	}
	return agents
}

// Start all simulated node agents
func (agents *E2Agents) Start() error {
	for _, a := range agents.Agents {
		err := a.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop all simulated node agents
func (agents *E2Agents) Stop() error {
	for _, a := range agents.Agents {
		err := a.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
