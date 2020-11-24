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

package indexedmap

import (
	"bytes"
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/golang/protobuf/proto"
	"io"
)

// Service is a state machine for a map primitive
type Service struct {
	primitive.Service
	lastIndex  uint64
	entries    map[string]*LinkedMapEntryValue
	indexes    map[uint64]*LinkedMapEntryValue
	firstEntry *LinkedMapEntryValue
	lastEntry  *LinkedMapEntryValue
	timers     map[string]primitive.Timer
	listeners  map[primitive.SessionID]map[primitive.StreamID]listener
}

// init initializes the map service
func (m *Service) init() {
	m.RegisterUnaryOperation(opPut, m.Put)
	m.RegisterUnaryOperation(opReplace, m.Replace)
	m.RegisterUnaryOperation(opRemove, m.Remove)
	m.RegisterUnaryOperation(opGet, m.Get)
	m.RegisterUnaryOperation(opFirstEntry, m.FirstEntry)
	m.RegisterUnaryOperation(opLastEntry, m.LastEntry)
	m.RegisterUnaryOperation(opPrevEntry, m.PrevEntry)
	m.RegisterUnaryOperation(opNextEntry, m.NextEntry)
	m.RegisterUnaryOperation(opExists, m.Exists)
	m.RegisterUnaryOperation(opSize, m.Size)
	m.RegisterUnaryOperation(opClear, m.Clear)
	m.RegisterStreamOperation(opEvents, m.Events)
	m.RegisterStreamOperation(opEntries, m.Entries)
}

// LinkedMapEntryValue is a doubly linked MapEntryValue
type LinkedMapEntryValue struct {
	*MapEntryValue
	Prev *LinkedMapEntryValue
	Next *LinkedMapEntryValue
}

// Backup takes a snapshot of the service
func (m *Service) Backup(writer io.Writer) error {
	listeners := make([]*Listener, 0)
	for sessionID, sessionListeners := range m.listeners {
		for streamID, sessionListener := range sessionListeners {
			listeners = append(listeners, &Listener{
				SessionId: uint64(sessionID),
				StreamId:  uint64(streamID),
				Key:       sessionListener.key,
			})
		}
	}

	if err := util.WriteVarInt(writer, len(listeners)); err != nil {
		return err
	}
	if err := util.WriteSlice(writer, listeners, proto.Marshal); err != nil {
		return err
	}
	if err := util.WriteVarUint64(writer, m.lastIndex); err != nil {
		return err
	}
	if err := util.WriteVarInt(writer, len(m.entries)); err != nil {
		return err
	}
	entry := m.firstEntry
	for entry != nil {
		if err := util.WriteValue(writer, entry.MapEntryValue, proto.Marshal); err != nil {
			return err
		}
		entry = entry.Next
	}
	return nil
}

// Restore restores the service from a snapshot
func (m *Service) Restore(reader io.Reader) error {
	length, err := util.ReadVarInt(reader)
	if err != nil {
		return err
	}

	listeners := make([]*Listener, length)
	err = util.ReadSlice(reader, listeners, func(data []byte) (*Listener, error) {
		listener := &Listener{}
		if err := proto.Unmarshal(data, listener); err != nil {
			return nil, err
		}
		return listener, nil
	})
	if err != nil {
		return err
	}

	m.listeners = make(map[primitive.SessionID]map[primitive.StreamID]listener)
	for _, snapshotListener := range listeners {
		sessionListeners, ok := m.listeners[primitive.SessionID(snapshotListener.SessionId)]
		if !ok {
			sessionListeners = make(map[primitive.StreamID]listener)
			m.listeners[primitive.SessionID(snapshotListener.SessionId)] = sessionListeners
		}
		sessionListeners[primitive.StreamID(snapshotListener.StreamId)] = listener{
			key:    snapshotListener.Key,
			stream: m.Session(primitive.SessionID(snapshotListener.SessionId)).Stream(primitive.StreamID(snapshotListener.StreamId)),
		}
	}

	lastIndex, err := util.ReadVarUint64(reader)
	if err != nil {
		return err
	}
	m.lastIndex = lastIndex

	entryCount, err := util.ReadVarInt(reader)
	if err != nil {
		return err
	}

	var prevEntry *LinkedMapEntryValue
	m.firstEntry = nil
	m.lastEntry = nil
	m.entries = make(map[string]*LinkedMapEntryValue)
	m.indexes = make(map[uint64]*LinkedMapEntryValue)
	for i := 0; i < entryCount; i++ {
		value, err := util.ReadValue(reader, func(data []byte) (*MapEntryValue, error) {
			entry := &MapEntryValue{}
			if err := proto.Unmarshal(data, entry); err != nil {
				return nil, err
			}
			return entry, nil
		})
		if err != nil {
			return err
		}

		linkedEntry := &LinkedMapEntryValue{
			MapEntryValue: value.(*MapEntryValue),
		}
		if m.firstEntry == nil {
			m.firstEntry = linkedEntry
		}
		m.entries[linkedEntry.Key] = linkedEntry
		m.indexes[linkedEntry.Index] = linkedEntry
		if prevEntry != nil {
			prevEntry.Next = linkedEntry
			linkedEntry.Prev = prevEntry
		}
		prevEntry = linkedEntry
		m.lastEntry = linkedEntry
	}
	return nil
}

// Put puts a key/value pair in the map
func (m *Service) Put(value []byte) ([]byte, error) {
	request := &PutRequest{}
	if err := proto.Unmarshal(value, request); err != nil {
		return nil, err
	}

	var oldEntry *LinkedMapEntryValue
	if request.Index > 0 {
		oldEntry = m.indexes[request.Index]
		if oldEntry != nil && oldEntry.Key != request.Key {
			return proto.Marshal(&PutResponse{
				Status: UpdateStatus_PRECONDITION_FAILED,
			})
		}
	} else {
		oldEntry = m.entries[request.Key]
	}

	if oldEntry == nil {
		// If the version is positive then reject the request.
		if !request.IfEmpty && request.Version > 0 {
			return proto.Marshal(&PutResponse{
				Status: UpdateStatus_PRECONDITION_FAILED,
			})
		}

		// Increment the index for a new entry
		var index uint64
		if request.Index > 0 {
			if request.Index > m.lastIndex {
				m.lastIndex = request.Index
			}
			index = request.Index
		} else {
			m.lastIndex++
			index = m.lastIndex
		}

		// Create a new entry value and set it in the map.
		newEntry := &LinkedMapEntryValue{
			MapEntryValue: &MapEntryValue{
				Index:   index,
				Key:     request.Key,
				Value:   request.Value,
				Version: uint64(m.Index()),
				TTL:     request.TTL,
				Created: m.Timestamp(),
				Updated: m.Timestamp(),
			},
		}
		m.entries[newEntry.Key] = newEntry
		m.indexes[newEntry.Index] = newEntry

		// Set the first entry if not set
		if m.firstEntry == nil {
			m.firstEntry = newEntry
		}

		// If the last entry is set, link it to the new entry
		if request.Index > 0 {
			if m.lastIndex == request.Index {
				if m.lastEntry != nil {
					m.lastEntry.Next = newEntry
					newEntry.Prev = m.lastEntry
				}
				m.lastEntry = newEntry
			}
		} else {
			if m.lastEntry != nil {
				m.lastEntry.Next = newEntry
				newEntry.Prev = m.lastEntry
			}
		}

		// Update the last entry
		m.lastEntry = newEntry

		// Schedule the timeout for the value if necessary.
		m.scheduleTTL(request.Key, newEntry)

		// Publish an event to listener streams.
		m.sendEvent(&ListenResponse{
			Type:    ListenResponse_INSERTED,
			Key:     request.Key,
			Index:   newEntry.Index,
			Value:   newEntry.Value,
			Version: newEntry.Version,
			Created: newEntry.Created,
			Updated: newEntry.Updated,
		})

		return proto.Marshal(&PutResponse{
			Status: UpdateStatus_OK,
			Index:  newEntry.Index,
			Key:    newEntry.Key,
		})
	}

	// If the version is -1 then reject the request.
	// If the version is positive then compare the version to the current version.
	if request.IfEmpty || (!request.IfEmpty && request.Version > 0 && request.Version != oldEntry.Version) {
		return proto.Marshal(&PutResponse{
			Status:          UpdateStatus_PRECONDITION_FAILED,
			Index:           oldEntry.Index,
			Key:             oldEntry.Key,
			PreviousValue:   oldEntry.Value,
			PreviousVersion: oldEntry.Version,
		})
	}

	// If the value is equal to the current value, return a no-op.
	if bytes.Equal(oldEntry.Value, request.Value) {
		return proto.Marshal(&PutResponse{
			Status:          UpdateStatus_NOOP,
			Index:           oldEntry.Index,
			Key:             oldEntry.Key,
			PreviousValue:   oldEntry.Value,
			PreviousVersion: oldEntry.Version,
		})
	}

	// Create a new entry value and set it in the map.
	newEntry := &LinkedMapEntryValue{
		MapEntryValue: &MapEntryValue{
			Index:   oldEntry.Index,
			Key:     oldEntry.Key,
			Value:   request.Value,
			Version: uint64(m.Index()),
			TTL:     request.TTL,
			Created: oldEntry.Created,
			Updated: m.Timestamp(),
		},
		Prev: oldEntry.Prev,
		Next: oldEntry.Next,
	}
	m.entries[newEntry.Key] = newEntry
	m.indexes[newEntry.Index] = newEntry

	// Update links for previous and next entries
	if newEntry.Prev != nil {
		oldEntry.Prev.Next = newEntry
	} else {
		m.firstEntry = newEntry
	}
	if newEntry.Next != nil {
		oldEntry.Next.Prev = newEntry
	} else {
		m.lastEntry = newEntry
	}

	// Schedule the timeout for the value if necessary.
	m.scheduleTTL(request.Key, newEntry)

	// Publish an event to listener streams.
	m.sendEvent(&ListenResponse{
		Type:    ListenResponse_UPDATED,
		Key:     request.Key,
		Index:   newEntry.Index,
		Value:   newEntry.Value,
		Version: newEntry.Version,
		Created: newEntry.Created,
		Updated: newEntry.Updated,
	})

	return proto.Marshal(&PutResponse{
		Status:          UpdateStatus_OK,
		Index:           newEntry.Index,
		Key:             newEntry.Key,
		PreviousValue:   oldEntry.Value,
		PreviousVersion: oldEntry.Version,
		Created:         newEntry.Created,
		Updated:         newEntry.Updated,
	})
}

// Replace replaces a key/value pair in the map
func (m *Service) Replace(value []byte) ([]byte, error) {
	request := &ReplaceRequest{}
	if err := proto.Unmarshal(value, request); err != nil {
		return nil, err
	}

	var oldEntry *LinkedMapEntryValue
	if request.Index > 0 {
		oldEntry = m.indexes[request.Index]
	} else {
		oldEntry = m.entries[request.Key]
	}

	if oldEntry == nil {
		return proto.Marshal(&ReplaceResponse{
			Status: UpdateStatus_PRECONDITION_FAILED,
		})
	}

	// If the version was specified and does not match the entry version, fail the replace.
	if request.PreviousVersion != 0 && request.PreviousVersion != oldEntry.Version {
		return proto.Marshal(&ReplaceResponse{
			Status: UpdateStatus_PRECONDITION_FAILED,
		})
	}

	// If the value was specified and does not match the entry value, fail the replace.
	if len(request.PreviousValue) != 0 && bytes.Equal(request.PreviousValue, oldEntry.Value) {
		return proto.Marshal(&ReplaceResponse{
			Status: UpdateStatus_PRECONDITION_FAILED,
		})
	}

	// If we've made it this far, update the entry.
	// Create a new entry value and set it in the map.
	newEntry := &LinkedMapEntryValue{
		MapEntryValue: &MapEntryValue{
			Index:   oldEntry.Index,
			Key:     oldEntry.Key,
			Value:   request.NewValue,
			Version: uint64(m.Index()),
			TTL:     request.TTL,
			Created: oldEntry.Created,
			Updated: m.Timestamp(),
		},
		Prev: oldEntry.Prev,
		Next: oldEntry.Next,
	}

	m.entries[newEntry.Key] = newEntry
	m.indexes[newEntry.Index] = newEntry

	// Update links for previous and next entries
	if newEntry.Prev != nil {
		oldEntry.Prev.Next = newEntry
	} else {
		m.firstEntry = newEntry
	}
	if newEntry.Next != nil {
		oldEntry.Next.Prev = newEntry
	} else {
		m.lastEntry = newEntry
	}

	// Schedule the timeout for the value if necessary.
	m.scheduleTTL(request.Key, newEntry)

	// Publish an event to listener streams.
	m.sendEvent(&ListenResponse{
		Type:    ListenResponse_UPDATED,
		Key:     request.Key,
		Index:   newEntry.Index,
		Value:   newEntry.Value,
		Version: newEntry.Version,
		Created: newEntry.Created,
		Updated: newEntry.Updated,
	})

	return proto.Marshal(&ReplaceResponse{
		Status:          UpdateStatus_OK,
		Index:           newEntry.Index,
		Key:             newEntry.Key,
		PreviousValue:   oldEntry.Value,
		PreviousVersion: oldEntry.Version,
		NewVersion:      newEntry.Version,
		Created:         newEntry.Created,
		Updated:         newEntry.Updated,
	})
}

// Remove removes a key/value pair from the map
func (m *Service) Remove(bytes []byte) ([]byte, error) {
	request := &RemoveRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	var entry *LinkedMapEntryValue
	if request.Index > 0 {
		entry = m.indexes[request.Index]
	} else {
		entry = m.entries[request.Key]
	}

	if entry == nil {
		return proto.Marshal(&RemoveResponse{
			Status: UpdateStatus_NOOP,
		})
	}

	// If the request version is set, verify that the request version matches the entry version.
	if request.Version > 0 && request.Version != entry.Version {
		return proto.Marshal(&RemoveResponse{
			Status: UpdateStatus_PRECONDITION_FAILED,
		})
	}

	// Delete the entry from the map.
	delete(m.entries, entry.Key)
	delete(m.indexes, entry.Index)

	// Cancel any TTLs.
	m.cancelTTL(request.Key)

	// Update links for previous and next entries
	if entry.Prev != nil {
		entry.Prev.Next = entry.Next
	} else {
		m.firstEntry = entry.Next
	}
	if entry.Next != nil {
		entry.Next.Prev = entry.Prev
	} else {
		m.lastEntry = entry.Prev
	}

	// Publish an event to listener streams.
	m.sendEvent(&ListenResponse{
		Type:    ListenResponse_REMOVED,
		Key:     entry.Key,
		Index:   entry.Index,
		Value:   entry.Value,
		Version: entry.Version,
		Created: entry.Created,
		Updated: entry.Updated,
	})

	return proto.Marshal(&RemoveResponse{
		Status:          UpdateStatus_OK,
		Index:           entry.Index,
		Key:             entry.Key,
		PreviousValue:   entry.Value,
		PreviousVersion: entry.Version,
		Created:         entry.Created,
		Updated:         entry.Updated,
	})
}

// Get gets a value from the map
func (m *Service) Get(bytes []byte) ([]byte, error) {
	request := &GetRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	var entry *LinkedMapEntryValue
	var ok bool
	if request.Index > 0 {
		entry, ok = m.indexes[request.Index]
	} else {
		entry, ok = m.entries[request.Key]
	}

	if !ok {
		return proto.Marshal(&GetResponse{})
	}
	return proto.Marshal(&GetResponse{
		Index:   entry.Index,
		Key:     entry.Key,
		Value:   entry.Value,
		Version: entry.Version,
		Created: entry.Created,
		Updated: entry.Updated,
	})
}

// FirstEntry gets the first entry from the map
func (m *Service) FirstEntry(bytes []byte) ([]byte, error) {
	request := &FirstEntryRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if m.firstEntry == nil {
		return proto.Marshal(&FirstEntryResponse{})
	}
	return proto.Marshal(&FirstEntryResponse{
		Index:   m.firstEntry.Index,
		Key:     m.firstEntry.Key,
		Value:   m.firstEntry.Value,
		Version: m.firstEntry.Version,
		Created: m.firstEntry.Created,
		Updated: m.firstEntry.Updated,
	})
}

// LastEntry gets the last entry from the map
func (m *Service) LastEntry(bytes []byte) ([]byte, error) {
	request := &LastEntryRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if m.lastEntry == nil {
		return proto.Marshal(&LastEntryResponse{})
	}
	return proto.Marshal(&LastEntryResponse{
		Index:   m.lastEntry.Index,
		Key:     m.lastEntry.Key,
		Value:   m.lastEntry.Value,
		Version: m.lastEntry.Version,
		Created: m.lastEntry.Created,
		Updated: m.lastEntry.Updated,
	})
}

// PrevEntry gets the previous entry from the map
func (m *Service) PrevEntry(bytes []byte) ([]byte, error) {
	request := &PrevEntryRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	entry, ok := m.indexes[request.Index]
	if !ok {
		for _, e := range m.indexes {
			if entry == nil || (e.Index < request.Index && e.Index > entry.Index) {
				entry = e
			}
		}
	} else {
		entry = entry.Prev
	}

	if entry == nil {
		return proto.Marshal(&PrevEntryResponse{})
	}
	return proto.Marshal(&PrevEntryResponse{
		Index:   entry.Index,
		Key:     entry.Key,
		Value:   entry.Value,
		Version: entry.Version,
		Created: entry.Created,
		Updated: entry.Updated,
	})
}

// NextEntry gets the next entry from the map
func (m *Service) NextEntry(bytes []byte) ([]byte, error) {
	request := &NextEntryRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	entry, ok := m.indexes[request.Index]
	if !ok {
		for _, e := range m.indexes {
			if entry == nil || (e.Index > request.Index && e.Index < entry.Index) {
				entry = e
			}
		}
	} else {
		entry = entry.Next
	}

	if entry == nil {
		return proto.Marshal(&NextEntryResponse{})
	}
	return proto.Marshal(&NextEntryResponse{
		Index:   entry.Index,
		Key:     entry.Key,
		Value:   entry.Value,
		Version: entry.Version,
		Created: entry.Created,
		Updated: entry.Updated,
	})
}

// Exists checks if the map contains a key
func (m *Service) Exists(bytes []byte) ([]byte, error) {
	request := &ContainsKeyRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	_, ok := m.entries[request.Key]
	return proto.Marshal(&ContainsKeyResponse{
		ContainsKey: ok,
	})
}

// Size returns the size of the map
func (m *Service) Size(bytes []byte) ([]byte, error) {
	return proto.Marshal(&SizeResponse{
		Size_: uint32(len(m.entries)),
	})
}

// Clear removes all entries from the map
func (m *Service) Clear(value []byte) ([]byte, error) {
	m.entries = make(map[string]*LinkedMapEntryValue)
	m.indexes = make(map[uint64]*LinkedMapEntryValue)
	m.firstEntry = nil
	m.lastEntry = nil
	return proto.Marshal(&ClearResponse{})
}

// Events sends change events to the client
func (m *Service) Events(bytes []byte, stream primitive.Stream) {
	request := &ListenRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		stream.Error(err)
		stream.Close()
		return
	}

	// Create and populate the listener
	lis := listener{
		key:    request.Key,
		index:  request.Index,
		stream: stream,
	}
	listeners, ok := m.listeners[stream.Session().ID()]
	if !ok {
		listeners = make(map[primitive.StreamID]listener)
		m.listeners[stream.Session().ID()] = listeners
	}
	listeners[stream.ID()] = lis

	if request.Replay {
		entry := m.firstEntry
		for entry != nil {
			bytes, err := proto.Marshal(&ListenResponse{
				Type:    ListenResponse_NONE,
				Key:     entry.Key,
				Index:   entry.Index,
				Value:   entry.Value,
				Version: entry.Version,
				Created: entry.Created,
				Updated: entry.Updated,
			})
			if err != nil {
				stream.Error(err)
				continue
			}

			if lis.key != "" {
				if entry.Key == lis.key {
					lis.stream.Value(bytes)
				}
			} else if lis.index > 0 {
				if entry.Index == lis.index {
					lis.stream.Value(bytes)
				}
			} else {
				lis.stream.Value(bytes)
			}
			entry = entry.Next
		}
	}
}

// Entries returns a stream of entries to the client
func (m *Service) Entries(value []byte, stream primitive.Stream) {
	defer stream.Close()
	entry := m.firstEntry
	for entry != nil {
		stream.Result(proto.Marshal(&EntriesResponse{
			Key:     entry.Key,
			Index:   entry.Index,
			Value:   entry.Value,
			Version: entry.Version,
			Created: entry.Created,
			Updated: entry.Updated,
		}))
		entry = entry.Next
	}
}

func (m *Service) scheduleTTL(key string, entry *LinkedMapEntryValue) {
	m.cancelTTL(key)
	if entry.TTL != nil {
		m.timers[key] = m.ScheduleOnce(entry.Created.Add(*entry.TTL).Sub(m.Timestamp()), func() {
			delete(m.entries, key)
			delete(m.indexes, entry.Index)

			// Update links for previous and next entries
			if entry.Prev != nil {
				entry.Prev.Next = entry.Next
			} else {
				m.firstEntry = entry.Next
			}
			if entry.Next != nil {
				entry.Next.Prev = entry.Prev
			} else {
				m.lastEntry = entry.Prev
			}

			m.sendEvent(&ListenResponse{
				Type:    ListenResponse_REMOVED,
				Key:     key,
				Index:   entry.Index,
				Value:   entry.Value,
				Version: uint64(entry.Version),
				Created: entry.Created,
				Updated: entry.Updated,
			})
		})
	}
}

func (m *Service) cancelTTL(key string) {
	timer, ok := m.timers[key]
	if ok {
		timer.Cancel()
	}
}

func (m *Service) sendEvent(event *ListenResponse) {
	bytes, _ := proto.Marshal(event)
	for sessionID, listeners := range m.listeners {
		session := m.Session(sessionID)
		if session != nil {
			for _, listener := range listeners {
				if listener.key != "" {
					if event.Key == listener.key {
						listener.stream.Value(bytes)
					}
				} else if listener.index > 0 {
					if event.Index == listener.index {
						listener.stream.Value(bytes)
					}
				} else {
					listener.stream.Value(bytes)
				}
			}
		}
	}
}

type listener struct {
	key    string
	index  uint64
	stream stream.WriteStream
}
