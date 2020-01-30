// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dispatcher

import (
	"fmt"
	log "k8s.io/klog"
	"sync"
)

type Dispatcher struct {
	nbiUeListenersLock sync.RWMutex
	nbiUeListeners map[string]chan Event
	nbiRouteListenersLock sync.RWMutex
	nbiRouteListeners map[string]chan Event
}

// NewDispatcher creates and initializes a new event dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		nbiUeListeners: make(map[string]chan Event),
		nbiRouteListeners: make(map[string]chan Event),
	}
}

func (d *Dispatcher) ListenUeEvents(ueEventChannel <-chan Event) {
	log.Info("User Equipment Event listener initialized")

	for ueEvent := range ueEventChannel {
		d.nbiUeListenersLock.RLock()
		for _, nbiChan := range d.nbiUeListeners {
			nbiChan <- ueEvent
		}
		d.nbiUeListenersLock.RUnlock()
	}
}

func (d *Dispatcher) RegisterUeListener(subscriber string) (chan Event, error) {
	d.nbiUeListenersLock.Lock()
	defer d.nbiUeListenersLock.Unlock()
	if _, ok := d.nbiUeListeners[subscriber]; ok {
		return nil, fmt.Errorf("NBI UE %s is already registered", subscriber)
	}
	channel := make(chan Event);
	d.nbiUeListeners[subscriber] = channel
	return channel, nil
}

func (d *Dispatcher) UnregisterUeListener(subscriber string) {
	d.nbiUeListenersLock.Lock()
	defer d.nbiUeListenersLock.Unlock()
	channel, ok := d.nbiUeListeners[subscriber]
	if !ok {
		log.Infof("Subscriber %s had not been registered", subscriber)
		return
	}
	delete(d.nbiUeListeners, subscriber)
	close(channel)
}

func (d *Dispatcher) ListenRouteEvents(routeEventChannel <-chan Event) {
	log.Info("Route Event listener initialized")

	for routeEvent := range routeEventChannel {
		d.nbiRouteListenersLock.RLock()
		for _, nbiChan := range d.nbiRouteListeners {
			nbiChan <- routeEvent
		}
		d.nbiRouteListenersLock.RUnlock()
	}
}

func (d *Dispatcher) RegisterRouteListener(subscriber string) (chan Event, error) {
	d.nbiRouteListenersLock.Lock()
	defer d.nbiRouteListenersLock.Unlock()
	if _, ok := d.nbiRouteListeners[subscriber]; ok {
		return nil, fmt.Errorf("NBI Route %s is already registered", subscriber)
	}
	channel := make(chan Event);
	d.nbiRouteListeners[subscriber] = channel
	return channel, nil
}

func (d *Dispatcher) UnregisterRouteListener(subscriber string) {
	d.nbiRouteListenersLock.Lock()
	defer d.nbiRouteListenersLock.Unlock()
	channel, ok := d.nbiRouteListeners[subscriber]
	if !ok {
		log.Infof("Subscriber %s had not been registered", subscriber)
		return
	}
	delete(d.nbiRouteListeners, subscriber)
	close(channel)
}

// GetListeners returns a list of registered listeners names
func (d *Dispatcher) GetListeners() []string {
	listenerKeys := make([]string, 0)
	d.nbiUeListenersLock.RLock()
	defer d.nbiUeListenersLock.RUnlock()
	for k := range d.nbiUeListeners {
		listenerKeys = append(listenerKeys, k)
	}
	d.nbiRouteListenersLock.RLock()
	defer d.nbiRouteListenersLock.RUnlock()
	for k := range d.nbiRouteListeners {
		listenerKeys = append(listenerKeys, k)
	}
	return listenerKeys
}
