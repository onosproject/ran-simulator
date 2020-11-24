// Copyright 2019-present Open Networking Foundation.
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

package device

import (
	"context"
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	deviceapi "github.com/onosproject/onos-topo/api/device"
	"github.com/onosproject/onos-topo/pkg/config"
	"io"
	"time"
)

// NewAtomixStore returns a new persistent Store
func NewAtomixStore() (Store, error) {
	ricConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	database, err := atomix.GetDatabase(ricConfig.Atomix, ricConfig.Atomix.GetDatabase(atomix.DatabaseTypeConsensus))
	if err != nil {
		return nil, err
	}

	devices, err := database.GetMap(context.Background(), "devices")
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		devices: devices,
	}, nil
}

// NewLocalStore returns a new local device store
func NewLocalStore() (Store, error) {
	node, address := atomix.StartLocalNode()
	name := primitive.Name{
		Namespace: "local",
		Name:      "devices",
	}

	session, err := primitive.NewSession(context.TODO(), primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}

	devices, err := _map.New(context.Background(), name, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		devices: devices,
		closer:  node.Stop,
	}, nil
}

// Store stores topology information
type Store interface {
	io.Closer

	// Load loads a device from the store
	Load(deviceID deviceapi.ID) (*deviceapi.Device, error)

	// Store stores a device in the store
	Store(*deviceapi.Device) error

	// Delete deletes a device from the store
	Delete(*deviceapi.Device) error

	// List streams devices to the given channel
	List(chan<- *deviceapi.Device) error

	// Watch streams device events to the given channel
	Watch(chan<- *Event) error
}

// atomixStore is the device implementation of the Store
type atomixStore struct {
	devices _map.Map
	closer  func() error
}

func (s *atomixStore) Load(deviceID deviceapi.ID) (*deviceapi.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	entry, err := s.devices.Get(ctx, string(deviceID))
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}
	return decodeDevice(entry)
}

func (s *atomixStore) Store(device *deviceapi.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bytes, err := proto.Marshal(device)
	if err != nil {
		return err
	}

	// Put the device in the map using an optimistic lock if this is an update
	var entry *_map.Entry
	if device.Revision == 0 {
		entry, err = s.devices.Put(ctx, string(device.ID), bytes)
	} else {
		entry, err = s.devices.Put(ctx, string(device.ID), bytes, _map.IfVersion(_map.Version(device.Revision)))
	}

	if err != nil {
		return err
	}

	// Update the device metadata
	device.Revision = deviceapi.Revision(entry.Version)
	return err
}

func (s *atomixStore) Delete(device *deviceapi.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if device.Revision > 0 {
		_, err := s.devices.Remove(ctx, string(device.ID), _map.IfVersion(_map.Version(device.Revision)))
		return err
	}
	_, err := s.devices.Remove(ctx, string(device.ID))
	return err
}

func (s *atomixStore) List(ch chan<- *deviceapi.Device) error {
	mapCh := make(chan *_map.Entry)
	if err := s.devices.Entries(context.Background(), mapCh); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for entry := range mapCh {
			if device, err := decodeDevice(entry); err == nil {
				ch <- device
			}
		}
	}()
	return nil
}

func (s *atomixStore) Watch(ch chan<- *Event) error {
	mapCh := make(chan *_map.Event)
	if err := s.devices.Watch(context.Background(), mapCh, _map.WithReplay()); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range mapCh {
			if device, err := decodeDevice(event.Entry); err == nil {
				ch <- &Event{
					Type:   EventType(event.Type),
					Device: device,
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = s.devices.Close(ctx)
	cancel()
	if s.closer != nil {
		return s.closer()
	}
	return nil
}

func decodeDevice(entry *_map.Entry) (*deviceapi.Device, error) {
	device := &deviceapi.Device{}
	if err := proto.Unmarshal(entry.Value, device); err != nil {
		return nil, err
	}
	device.ID = deviceapi.ID(entry.Key)
	device.Revision = deviceapi.Revision(entry.Version)
	return device, nil
}

// EventType provides the type for a device event
type EventType string

const (
	// EventNone is no event
	EventNone EventType = ""
	// EventInserted is inserted
	EventInserted EventType = "inserted"
	// EventUpdated is updated
	EventUpdated EventType = "updated"
	// EventRemoved is removed
	EventRemoved EventType = "removed"
)

// Event is a store event for a device
type Event struct {
	Type   EventType
	Device *deviceapi.Device
}
