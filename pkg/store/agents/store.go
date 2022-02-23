// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package agents

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"

	"sync"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/e2agent"
)

var log = logging.GetLogger()

// E2Agents e2 agents
type E2Agents struct {
	agents map[types.GnbID]e2agent.E2Agent
	mu     sync.RWMutex
}

// NewStore creates a new e2 agents store
func NewStore() *E2Agents {
	return &E2Agents{
		agents: make(map[types.GnbID]e2agent.E2Agent),
		mu:     sync.RWMutex{},
	}
}

// Get gets an e2 agent
func (e *E2Agents) Get(id types.GnbID) (e2agent.E2Agent, error) {
	log.Debug("Getting e2 agent with ID:", id)
	e.mu.RLock()
	defer e.mu.RUnlock()
	if val, ok := e.agents[id]; ok {
		return val, nil
	}
	return nil, errors.New(errors.NotFound, "e2 agent has not been found")
}

// Add adds an e2 agent
func (e *E2Agents) Add(id types.GnbID, agent e2agent.E2Agent) error {
	log.Debug("Adding e2 agent with ID:", id)
	e.mu.Lock()
	defer e.mu.Unlock()
	if id == 0 {
		return errors.New(errors.Invalid, "E2 node ID cannot be empty or zero")
	}
	e.agents[id] = agent
	return nil

}

// Remove removes an e2 agent
func (e *E2Agents) Remove(id types.GnbID) error {
	log.Debug("Removing e2 agent with ID:", id)
	e.mu.Lock()
	defer e.mu.Unlock()
	if id == 0 {
		return errors.New(errors.Invalid, "E2 node ID cannot be empty or zero")
	}
	delete(e.agents, id)
	return nil
}

// List list e2 agents
func (e *E2Agents) List() (map[types.GnbID]e2agent.E2Agent, error) {
	return e.agents, nil
}

// Store e2 agents store interface
type Store interface {
	// Add an e2 agent
	Add(types.GnbID, e2agent.E2Agent) error

	// Remove an e2 agent
	Remove(types.GnbID) error

	// List list all of the e2 agents
	List() (map[types.GnbID]e2agent.E2Agent, error)

	// Get gets an e2 agent
	Get(types.GnbID) (e2agent.E2Agent, error)
}

var _ Store = &E2Agents{}
