// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package model

import (
	"sync"
)

// ECGI is a type alias for ECGI ID
type ECGI string

// SMID is a type alias for service model ID
type SMID int32

// SimNode represents a simulated RAN E2 node
type SimNode struct {
	ECGI          ECGI
	ServiceModels []SMID
	// TODO: add other simulation attributes, i.e. locations, signal strength, neighbours, etc.
}

// SimNodes represents a collection of simulated RAN E2 nodes
type SimNodes struct {
	nodes map[ECGI]*SimNode
	lock  *sync.RWMutex
}

// NewSimNodes creates a new and empty collection of simulated RAN E2 nodes
func NewSimNodes() *SimNodes {
	return &SimNodes{
		nodes: make(map[ECGI]*SimNode),
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
func (n *SimNodes) Get(ecgi ECGI) *SimNode {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.nodes[ecgi]
}

// GetAll simulated E2 nodes in the collection
func (n *SimNodes) GetAll() []*SimNode {
	n.lock.RLock()
	defer n.lock.RUnlock()
	all := make([]*SimNode, 0)
	for _, v := range n.nodes {
		all = append(all, v)
	}
	return all
}

// Remove the simulated E2 node with the specified EGGI; returns nil of no such node
func (n *SimNodes) Remove(ecgi ECGI) *SimNode {
	n.lock.Lock()
	defer n.lock.Unlock()
	node := n.nodes[ecgi]
	delete(n.nodes, ecgi)
	return node
}
