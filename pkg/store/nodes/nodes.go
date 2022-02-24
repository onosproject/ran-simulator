// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package nodes

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	"github.com/onosproject/ran-simulator/pkg/store/watcher"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
)

var log = liblog.GetLogger()

// Store tracks inventory of simulated E2 nodes.
type Store interface {
	// Add adds the specified node to the registry
	Add(ctx context.Context, node *model.Node) error

	// Get retrieves the node with the specified GnbID
	Get(ctx context.Context, gnbID types.GnbID) (*model.Node, error)

	// Update updates the node
	Update(ctx context.Context, node *model.Node) error

	// Delete deletes the node with the specified GnbID
	Delete(ctx context.Context, gnbID types.GnbID) (*model.Node, error)

	// Watch watches the node inventory events using the supplied channel
	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error

	// List lists the nodes
	List(ctx context.Context) ([]*model.Node, error)

	// Len returns number of nodes
	Len(ctx context.Context) (int, error)

	// SetsStatus changes the E2 node agent status value
	SetStatus(ctx context.Context, gnbID types.GnbID, status string) error

	// PruneCell  the node that has the specified cell
	PruneCell(ctx context.Context, ncgi types.NCGI) error

	// Load add all nodes from the specified node map; no events will be generated
	Load(ctx context.Context, nodes map[string]model.Node)

	// Clear removes all nodes; no events will be generated
	Clear(ctx context.Context)
}

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

type store struct {
	mu       sync.RWMutex
	nodes    map[types.GnbID]*model.Node
	watchers *watcher.Watchers
}

// NewNodeRegistry creates a new store abstraction from the specified fixed node map.
func NewNodeRegistry(nodes map[string]model.Node) Store {
	log.Infof("Creating registry from model with %d nodes", len(nodes))
	watchers := watcher.NewWatchers()
	reg := &store{
		mu:       sync.RWMutex{},
		nodes:    make(map[types.GnbID]*model.Node),
		watchers: watchers,
	}

	reg.Load(context.Background(), nodes)

	log.Infof("Created registry primed with %d nodes", len(reg.nodes))
	return reg
}

// Load add all nodes from the specified node map; no events will be generated
func (s *store) Load(ctx context.Context, nodes map[string]model.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Copy the nodes into our own map
	for _, n := range nodes {
		node := n // avoids scopelint issue
		s.nodes[node.GnbID] = &node
	}
}

// Clear removes all nodes; no events will be generated
func (s *store) Clear(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id := range s.nodes {
		delete(s.nodes, id)
	}
}

// Add adds a new node
func (s *store) Add(ctx context.Context, node *model.Node) error {
	log.Debugf("Adding node with ID: %d", node.GnbID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodes[node.GnbID]; ok {
		return errors.New(errors.NotFound, "node with GnbID already exists")
	}

	s.nodes[node.GnbID] = node
	addEvent := event.Event{
		Key:   node.GnbID,
		Value: node,
		Type:  Created,
	}
	s.watchers.Send(addEvent)
	return nil

}

// Get gets a node based on a given ID
func (s *store) Get(ctx context.Context, gnbID types.GnbID) (*model.Node, error) {
	log.Debugf("Getting node with ID: %d", gnbID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if node, ok := s.nodes[gnbID]; ok {
		return node, nil
	}

	return nil, errors.New(errors.NotFound, "node not found")
}

// Update updates a node
func (s *store) Update(ctx context.Context, node *model.Node) error {
	log.Debugf("Updating node with ID:%d", node.GnbID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodes[node.GnbID]; ok {
		s.nodes[node.GnbID] = node
		updateEvent := event.Event{
			Key:   node.GnbID,
			Value: node,
			Type:  Updated,
		}

		s.watchers.Send(updateEvent)
		return nil
	}

	return errors.New(errors.NotFound, "node not found")
}

// PruneCell prunes a cell
func (s *store) PruneCell(ctx context.Context, ncgi types.NCGI) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, node := range s.nodes {
		for i, e := range node.Cells {
			if e == ncgi {
				node.Cells = removeNCGI(node.Cells, i)
				updateEvent := event.Event{
					Key:   node.GnbID,
					Value: node,
					Type:  Updated,
				}
				s.watchers.Send(updateEvent)
				return nil
			}
		}
	}
	return nil
}

func (s *store) SetStatus(ctx context.Context, gnbID types.GnbID, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if node, ok := s.nodes[gnbID]; ok {
		node.Status = status
		return nil
	}
	return errors.New(errors.NotFound, "node not found")
}

// Delete deletes a node
func (s *store) Delete(ctx context.Context, gnbID types.GnbID) (*model.Node, error) {
	log.Debugf("Deleting node %d:", gnbID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if node, ok := s.nodes[gnbID]; ok {
		delete(s.nodes, gnbID)
		deleteEvent := event.Event{
			Key:   node.GnbID,
			Value: node,
			Type:  Deleted,
		}
		s.watchers.Send(deleteEvent)
		return node, nil
	}
	return nil, errors.New(errors.NotFound, "node not found")
}

// Watch
func (s *store) Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error {
	log.Debug("Watching node changes")
	replay := len(options) > 0 && options[0].Replay
	id := uuid.New()
	err := s.watchers.AddWatcher(id, ch)
	if err != nil {
		log.Error(err)
		close(ch)
		return err
	}
	go func() {
		<-ctx.Done()
		err = s.watchers.RemoveWatcher(id)
		if err != nil {
			log.Error(err)
		}
		close(ch)
	}()

	if replay {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, node := range s.nodes {
				ch <- event.Event{
					Key:   node.GnbID,
					Value: node,
					Type:  None,
				}
			}
		}()
	}
	return nil
}

// List list of nodes
func (s *store) List(ctx context.Context) ([]*model.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Node, 0, len(s.nodes))
	for _, node := range s.nodes {
		list = append(list, node)
	}
	return list, nil
}

// Len number of nodes
func (s *store) Len(ctx context.Context) (int, error) {
	return len(s.nodes), nil
}

func removeNCGI(ecgis []types.NCGI, i int) []types.NCGI {
	ecgis[len(ecgis)-1], ecgis[i] = ecgis[i], ecgis[len(ecgis)-1]
	return ecgis[:len(ecgis)-1]
}
