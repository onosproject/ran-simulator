// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connections

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/onosproject/ran-simulator/pkg/store/watcher"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	"github.com/onosproject/ran-simulator/pkg/store/event"
)

var log = logging.GetLogger()

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

// NewConnectionID creates a new connection ID
func NewConnectionID(ricAddress string, ricPort uint64) ConnectionID {
	return ConnectionID{
		ricIPAddress: ricAddress,
		ricPort:      ricPort,
	}
}

// GetRICIPAddress gets RIC IP address
func (ch ConnectionID) GetRICIPAddress() string {
	return ch.ricIPAddress
}

// GetRICPort gets RIC port
func (ch ConnectionID) GetRICPort() uint64 {
	return ch.ricPort
}

// Connections data structure for storing connections
type Connections struct {
	connections map[ConnectionID]*Connection
	mu          sync.RWMutex
	watchers    *watcher.Watchers
}

// NewStore creates a new e2 agents store
func NewStore() *Connections {
	watchers := watcher.NewWatchers()
	return &Connections{
		connections: make(map[ConnectionID]*Connection),
		mu:          sync.RWMutex{},
		watchers:    watchers,
	}
}

// Add adds a connection to connection store
func (c *Connections) Add(ctx context.Context, id ConnectionID, connection *Connection) error {
	log.Infof("Adding a connection with connection ID: %v", id)
	c.mu.Lock()
	defer c.mu.Unlock()
	if id.ricIPAddress == "" || id.ricPort == 0 {
		return errors.NewInvalid("ric address or port number is invalid")
	}
	c.connections[id] = connection
	addChannelEvent := event.Event{
		Key:   id,
		Value: connection,
		Type:  Created,
	}

	c.watchers.Send(addChannelEvent)
	return nil

}

// Remove removes a connection from connection store
func (c *Connections) Remove(ctx context.Context, id ConnectionID) error {
	log.Infof("Removing a connection with connection ID: %v", id)
	c.mu.Lock()
	defer c.mu.Unlock()
	removeChannelEvent := event.Event{
		Key:   id,
		Value: c.connections[id],
		Type:  Deleted,
	}
	delete(c.connections, id)
	c.watchers.Send(removeChannelEvent)

	return nil
}

// Get gets connection based on a given connection ID
func (c *Connections) Get(ctx context.Context, id ConnectionID) (*Connection, error) {
	log.Debugf("Getting a connection with connection ID: %v", id)
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.connections[id]; ok {
		return val, nil
	}
	return nil, errors.NewNotFound("connection with ID %v not found", id)
}

// List list all of the available connections
func (c *Connections) List(ctx context.Context) []*Connection {
	c.mu.RLock()
	defer c.mu.RUnlock()
	connections := make([]*Connection, 0)
	for _, conn := range c.connections {
		connections = append(connections, conn)
	}

	return connections
}

// Update update a connection
func (c *Connections) Update(ctx context.Context, connection *Connection) error {
	log.Infof("Updating connection with ID %v:", connection.ID)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connections[connection.ID] = connection
	updateEvent := event.Event{
		Key:   connection.ID,
		Value: connection,
		Type:  Updated,
	}

	c.watchers.Send(updateEvent)
	return nil
}

// Watch watch connection events
func (c *Connections) Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error {
	log.Debug("Watching E2 node connection changes")
	replay := len(options) > 0 && options[0].Replay
	id := uuid.New()
	err := c.watchers.AddWatcher(id, ch)
	if err != nil {
		log.Error(err)
		close(ch)
		return err
	}
	go func() {
		<-ctx.Done()
		err = c.watchers.RemoveWatcher(id)
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
			for id, connection := range c.connections {
				ch <- event.Event{
					Key:   id,
					Value: connection,
					Type:  None,
				}
			}
		}()
	}
	return nil
}

// Store connection store interface
type Store interface {
	Add(ctx context.Context, id ConnectionID, connection *Connection) error

	Remove(ctx context.Context, id ConnectionID) error

	Get(ctx context.Context, id ConnectionID) (*Connection, error)

	List(ctx context.Context) []*Connection

	Update(ctx context.Context, connection *Connection) error

	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error
}

var _ Store = &Connections{}
