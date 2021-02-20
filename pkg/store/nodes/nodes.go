// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package nodes

import (
	"github.com/onosproject/onos-lib-go/pkg/errors"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"sync"
)

var log = liblog.GetLogger("store", "nodes")

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

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

type nodeWatcher struct {
	ch chan<- NodeEvent
}

func (r *nodeRegistry) notify(node *model.Node, eventType uint8) {
	event := NodeEvent{
		Node: node,
		Type: eventType,
	}
	for _, watcher := range r.watchers {
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
	log.Infof("Creating registry from model with %d nodes", len(nodes))
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

	log.Infof("Created registry primed with %d nodes", len(reg.nodes))
	return reg
}

func (r *nodeRegistry) AddNode(node *model.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.nodes[node.EnbID]; ok {
		return errors.New(errors.NotFound, "node with EnbID already exists")
	}

	r.nodes[node.EnbID] = node
	r.notify(node, ADDED)
	return nil

}

func (r *nodeRegistry) GetNode(enbID types.EnbID) (*model.Node, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if node, ok := r.nodes[enbID]; ok {
		return node, nil
	}

	return nil, errors.New(errors.NotFound, "node not found")
}

func (r *nodeRegistry) UpdateNode(node *model.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.nodes[node.EnbID]; ok {
		r.nodes[node.EnbID] = node
		r.notify(node, UPDATED)
		return nil
	}

	return errors.New(errors.NotFound, "node not found")
}

func (r *nodeRegistry) DeleteNode(enbID types.EnbID) (*model.Node, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if node, ok := r.nodes[enbID]; ok {
		delete(r.nodes, enbID)
		r.notify(node, DELETED)
		return node, nil
	}
	return nil, errors.New(errors.NotFound, "node not found")
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

func (r *nodeRegistry) WatchNodes(ch chan<- NodeEvent, options ...WatchOptions) {
	log.Infof("WatchNodes: %v (#%d)\n", options, len(r.nodes))
	monitor := len(options) == 0 || options[0].Monitor
	replay := len(options) > 0 && options[0].Replay
	go func() {
		watcher := nodeWatcher{ch: ch}
		if monitor {
			r.lock.RLock()
			r.watchers = append(r.watchers, watcher)
			r.lock.RUnlock()
		}

		if replay {
			r.lock.RLock()
			defer r.lock.RUnlock()
			for _, node := range r.nodes {
				ch <- NodeEvent{Node: node, Type: NONE}
			}
			if !monitor {
				close(ch)
			}
		}
	}()
}
