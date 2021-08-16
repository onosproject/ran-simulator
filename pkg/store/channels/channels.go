// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package channels

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/onosproject/ran-simulator/pkg/store/watcher"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap101"
)

var log = logging.GetLogger("store", "channels")

// WatchOptions allows tailoring the WatchNodes behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

// NewChannelID creates a new channel ID
func NewChannelID(ricAddress string, ricPort uint64) ChannelID {
	return ChannelID{
		ricAddress: ricAddress,
		ricPort:    ricPort,
	}
}

// GetRICAddress gets RIC address
func (ch ChannelID) GetRICAddress() string {
	return ch.ricAddress
}

// GetRICPort gets RIC port
func (ch ChannelID) GetRICPort() uint64 {
	return ch.ricPort
}

// Channels data structure for storing channels
type Channels struct {
	channels map[ChannelID]e2.ClientChannel
	mu       sync.RWMutex
	watchers *watcher.Watchers
}

// NewStore creates a new e2 agents store
func NewStore() *Channels {
	watchers := watcher.NewWatchers()
	return &Channels{
		channels: make(map[ChannelID]e2.ClientChannel),
		mu:       sync.RWMutex{},
		watchers: watchers,
	}
}

// Add adds a channel to channel store
func (c *Channels) Add(ctx context.Context, id ChannelID, channel e2.ClientChannel) error {
	log.Info("Adding a channel with channel ID: %v", id)
	c.mu.Lock()
	defer c.mu.Unlock()
	if id.ricAddress == "" || id.ricPort == 0 {
		return errors.NewInvalid("ric address or port number is invalid")
	}
	c.channels[id] = channel
	addChannelEvent := event.Event{
		Key:   id,
		Value: channel,
		Type:  Created,
	}

	c.watchers.Send(addChannelEvent)

	return nil

}

// Remove removes a channel from channel store
func (c *Channels) Remove(ctx context.Context, id ChannelID) error {
	log.Info("Removing a channel with channel ID: %v", id)
	c.mu.Lock()
	defer c.mu.Unlock()
	removeChannelEvent := event.Event{
		Key:   id,
		Value: c.channels[id],
		Type:  Deleted,
	}
	delete(c.channels, id)
	c.watchers.Send(removeChannelEvent)

	return nil
}

// Get gets channel based on a given channel ID
func (c *Channels) Get(ctx context.Context, id ChannelID) (e2.ClientChannel, error) {
	log.Debugf("Getting a channel with channel ID: %v", id)
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.channels[id]; ok {
		return val, nil
	}
	return nil, errors.NewNotFound("channel with ID %v not found", id)
}

// List list all of the available channels
func (c *Channels) List(ctx context.Context) []e2.ClientChannel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	channels := make([]e2.ClientChannel, 0)
	for _, channel := range c.channels {
		channels = append(channels, channel)
	}

	return channels
}

// Watch watch channel events
func (c *Channels) Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error {
	log.Debug("Watching E2 node channel changes")
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
			for id, channel := range c.channels {
				ch <- event.Event{
					Key:   id,
					Value: channel,
					Type:  None,
				}
			}
		}()
	}
	return nil
}

// Store channel store interface
type Store interface {
	Add(ctx context.Context, id ChannelID, channel e2.ClientChannel) error

	Remove(ctx context.Context, id ChannelID) error

	Get(ctx context.Context, id ChannelID) (e2.ClientChannel, error)

	List(ctx context.Context) []e2.ClientChannel

	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error
}

var _ Store = &Channels{}
