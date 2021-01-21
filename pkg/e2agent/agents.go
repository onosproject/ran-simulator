// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"github.com/onosproject/ran-simulator/pkg/model"
)

// E2Agents represents a collection of E2 agents to allow centralized management
type E2Agents struct {
	Agents map[model.Ecgi]E2Agent
}

// NewE2Agents creates a new collection of E2 agents from the specified list of nodes
func NewE2Agents(m *model.Model) *E2Agents {
	agents := &E2Agents{
		Agents: make(map[model.Ecgi]E2Agent),
	}

	for _, node := range m.Nodes {
		agents.Agents[node.Ecgi] = NewE2Agent(node, m)
	}
	return agents
}

// Start all simulated node agents
func (agents *E2Agents) Start() error {
	for _, a := range agents.Agents {
		log.Info("start agent:")
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
