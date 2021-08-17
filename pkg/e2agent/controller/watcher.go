// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package controller

import (
	"context"
	"sync"

	"github.com/onosproject/onos-lib-go/pkg/controller"
	"github.com/onosproject/ran-simulator/pkg/store/channels"
	"github.com/onosproject/ran-simulator/pkg/store/event"
)

// ChannelWatcher is a channel watcher
type ChannelWatcher struct {
	channels  channels.Store
	channelCh chan event.Event
	cancel    context.CancelFunc
	mu        sync.Mutex
}

// Start starts the channel watcher
func (w *ChannelWatcher) Start(ch chan<- controller.ID) error {
	log.Info("Starting Channel Watcher")
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	w.channelCh = make(chan event.Event, queueSize)
	ctx, cancel := context.WithCancel(context.Background())
	err := w.channels.Watch(ctx, w.channelCh, channels.WatchOptions{Replay: true})
	if err != nil {
		cancel()
		return err
	}
	w.cancel = cancel

	go func() {
		for channelEvent := range w.channelCh {
			ch <- controller.NewID(channelEvent.Key)

		}
		close(ch)
	}()

	return nil

}

// Stop stops the channel watcher
func (w *ChannelWatcher) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	w.mu.Unlock()
}
