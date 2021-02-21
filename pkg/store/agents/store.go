// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package agents

import (
	"sync"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/e2agent"
)

// E2Agents e2 agents
type E2Agents struct {
	agents map[types.EnbID]e2agent.E2Agent
	mu     sync.RWMutex
}

// NewStore creates a new e2 agents store
func NewStore() *E2Agents {
	return &E2Agents{
		agents: make(map[types.EnbID]e2agent.E2Agent),
		mu:     sync.RWMutex{},
	}
}

// Get gets an e2 agent
func (e *E2Agents) Get(id types.EnbID) (e2agent.E2Agent, error) {
	panic("implement me")
}

// Add adds an e2 agent
func (e *E2Agents) Add(id types.EnbID, agent e2agent.E2Agent) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if id == 0 {
		return errors.New(errors.Invalid, "E2 node ID cannot be empty or zero")
	}
	e.agents[id] = agent
	return nil

}

// Remove removes an e2 agent
func (e *E2Agents) Remove(id types.EnbID) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if id == 0 {
		return errors.New(errors.Invalid, "E2 node ID cannot be empty or zero")
	}
	delete(e.agents, id)
	return nil
}

// List list e2 agents
func (e *E2Agents) List() ([]e2agent.E2Agent, error) {
	return nil, nil
}

type Store interface {
	// Add an e2 agent
	Add(types.EnbID, e2agent.E2Agent) error

	// Remove an e2 agent
	Remove(types.EnbID) error

	// List list all of the e2 agents
	List() ([]e2agent.E2Agent, error)

	// Get gets an e2 agent
	Get(types.EnbID) (e2agent.E2Agent, error)
}

var _ Store = &E2Agents{}
