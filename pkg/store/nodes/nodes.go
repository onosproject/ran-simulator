// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package nodes

import (
	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"sync"
)

// NodeRegistry tracks inventory of simulated E2 nodes.
type NodeRegistry interface {
	// AddNode adds the specified node to the registry
	AddNode(node *model.Node) error

	// GetNode retrieves the node with the specified EnbID
	GetNode(enbID types.EnbID) (*model.Node, error)

	// UpdateNode updates the node
	UpdateNode(node *model.Node) error

	// DeleteNode deletes the node with the specified EnbID
	DeleteNode(enbID types.EnbID) (*model.Node, error)

	// WatchNodes watches the node inventory events using the supplied channel
	WatchNodes(ch chan<- NodeEvent, options ...WatchOptions)
}

// NodeEvent represents a change in the node inventory
type NodeEvent struct {
	Node *model.Node
	Type uint8
}

type nodeWatcher struct {
	id uuid.UUID
	ch chan<- NodeEvent
}

func (nr *nodeRegistry) notify(node *model.Node, eventType uint8) {
	event := NodeEvent{
		Node: node,
		Type: eventType,
	}
	for _, watcher := range nr.watchers {
		watcher.ch <- event
	}
}

type nodeRegistry struct {
	lock     sync.RWMutex
	nodes    map[types.EnbID]*model.Node
	watchers []nodeWatcher
}

// NewNodeRegistry creates a new store abstraction from the specified fixed node map.
func NewNodeRegistry(nodes map[string]model.Node) NodeRegistry {
	reg := &nodeRegistry{
		lock:     sync.RWMutex{},
		nodes:    make(map[types.EnbID]*model.Node),
		watchers: make([]nodeWatcher, 0, 8),
	}

	// Copy the nodes into our own map
	for _, n := range nodes {
		node := n // avoids scopelint issue
		reg.nodes[node.EnbID] = &node
	}

	return reg
}

const (
	// NONE indicates no change event
	NONE uint8 = 0

	// ADDED indicates new node was added
	ADDED uint8 = 1

	// UPDATED indicates an existing node was updated
	UPDATED uint8 = 2

	// DELETED indicates a node was deleted
	DELETED uint8 = 3
)

func (nr *nodeRegistry) AddNode(node *model.Node) error {
	nr.lock.Lock()
	defer nr.lock.Unlock()
	if _, ok := nr.nodes[node.EnbID]; ok {
		return errors.New(errors.NotFound, "node with EnbID already exists")
	}

	nr.nodes[node.EnbID] = node
	go nr.notify(node, ADDED)
	return nil

}

func (nr *nodeRegistry) GetNode(enbID types.EnbID) (*model.Node, error) {
	nr.lock.RLock()
	defer nr.lock.RUnlock()
	if node, ok := nr.nodes[enbID]; ok {
		return node, nil
	}

	return nil, errors.New(errors.NotFound, "node not found")
}

func (nr *nodeRegistry) UpdateNode(node *model.Node) error {
	nr.lock.Lock()
	defer nr.lock.Unlock()
	if _, ok := nr.nodes[node.EnbID]; ok {
		nr.nodes[node.EnbID] = node
		nr.notify(node, UPDATED)
		return nil
	}

	return errors.New(errors.NotFound, "node not found")
}

func (nr *nodeRegistry) DeleteNode(enbID types.EnbID) (*model.Node, error) {
	nr.lock.Lock()
	defer nr.lock.Unlock()
	if node, ok := nr.nodes[enbID]; ok {
		delete(nr.nodes, enbID)
		nr.notify(node, DELETED)
		return node, nil
	}
	return nil, errors.New(errors.NotFound, "node not found")
}

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

func (nr *nodeRegistry) WatchNodes(ch chan<- NodeEvent, options ...WatchOptions) {
	monitor := len(options) == 0 || options[0].Monitor
	replay := len(options) > 0 && options[0].Replay
	go func() {
		watcher := nodeWatcher{
			id: uuid.UUID{},
			ch: ch,
		}
		if monitor {
			nr.lock.RLock()
			nr.watchers = append(nr.watchers, watcher)
			nr.lock.RUnlock()
		}

		if replay {
			nr.lock.RLock()
			defer nr.lock.RUnlock()
			for _, node := range nr.nodes {
				ch <- NodeEvent{Node: node, Type: NONE}
			}
			if !monitor {
				close(ch)
			}
		}
	}()
}
