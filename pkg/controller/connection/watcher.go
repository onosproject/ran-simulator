// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connection

import (
	"context"
	"sync"

	"github.com/onosproject/onos-lib-go/pkg/controller"
	"github.com/onosproject/ran-simulator/pkg/store/connections"
	"github.com/onosproject/ran-simulator/pkg/store/event"
)

// Watcher is a connection watcher
type Watcher struct {
	connections  connections.Store
	connectionCh chan event.Event
	cancel       context.CancelFunc
	mu           sync.Mutex
}

// Start starts the connection watcher
func (w *Watcher) Start(ch chan<- controller.ID) error {
	log.Info("Starting Connection Watcher")
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	w.connectionCh = make(chan event.Event, queueSize)
	ctx, cancel := context.WithCancel(context.Background())
	err := w.connections.Watch(ctx, w.connectionCh, connections.WatchOptions{Replay: true})
	if err != nil {
		cancel()
		return err
	}
	w.cancel = cancel

	go func() {
		for connectionEvent := range w.connectionCh {
			log.Debugf("Received connection event: %v", connectionEvent)
			ch <- controller.NewID(connectionEvent.Key)

		}
		close(ch)
	}()

	return nil

}

// Stop stops the connection watcher
func (w *Watcher) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	w.mu.Unlock()
}
