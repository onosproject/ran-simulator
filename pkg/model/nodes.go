// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package model

import (
	"github.com/onosproject/ran-simulator/api/types"
	"sync"
)

// SimNode represents a simulated RAN E2 node
type SimNode struct {
	ECGI          types.ECGI
	ServiceModels []int32
	// Cell types.Cell // for now just encapsulate gRPC API cell; to be decoupled
}

// SimNodes represents a collection of simulated RAN E2 nodes
type SimNodes struct {
	nodes map[types.ECGI]*SimNode
	lock  *sync.RWMutex
}

// NewSimNodes creates a new and empty collection of simulated RAN E2 nodes
func NewSimNodes() *SimNodes {
	return &SimNodes{
		nodes: make(map[types.ECGI]*SimNode),
		lock:  &sync.RWMutex{},
	}
}

// Add a new simulated E2 node to the collection
func (n *SimNodes) Add(node *SimNode) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.nodes[node.ECGI] = node
}

// Get the simulated E2 node with the specified EGGI; returns nil of no such node
func (n *SimNodes) Get(ecgi types.ECGI) *SimNode {
	n.lock.RLock()
	defer n.lock.Unlock()
	return n.nodes[ecgi]
}

// GetAll simulated E2 nodes in the collection
func (n *SimNodes) GetAll() []*SimNode {
	n.lock.RLock()
	defer n.lock.Unlock()
	all := make([]*SimNode, len(n.nodes))
	for _, v := range n.nodes {
		all = append(all, v)
	}
	return all
}

// Remove the simulated E2 node with the specified EGGI; returns nil of no such node
func (n *SimNodes) Remove(ecgi types.ECGI) *SimNode {
	n.lock.Lock()
	defer n.lock.Unlock()
	node := n.nodes[ecgi]
	delete(n.nodes, ecgi)
	return node
}
