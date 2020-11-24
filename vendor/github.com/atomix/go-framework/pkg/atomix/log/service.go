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

package log

import (
	"bytes"
	"io"

	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/golang/protobuf/proto"
)

// Service is a state machine for a log primitive
type Service struct {
	primitive.Service
	lastIndex  uint64
	indexes    map[uint64]*LinkedLogEntryValue
	firstEntry *LinkedLogEntryValue
	lastEntry  *LinkedLogEntryValue
	listeners  map[primitive.SessionID]map[primitive.StreamID]listener
}

// init initializes the log service
func (m *Service) init() {
	m.RegisterUnaryOperation(opAppend, m.Append)
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

// LinkedLogEntryValue is a doubly linked LogEntryValue
type LinkedLogEntryValue struct {
	*LogEntryValue
	Prev *LinkedLogEntryValue
	Next *LinkedLogEntryValue
}

// Backup takes a snapshot of the service
func (m *Service) Backup(writer io.Writer) error {
	listeners := make([]*Listener, 0)
	for sessionID, sessionListeners := range m.listeners {
		for streamID, sessionListener := range sessionListeners {
			listeners = append(listeners, &Listener{
				SessionId: uint64(sessionID),
				StreamId:  uint64(streamID),
				Index:     sessionListener.index,
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
	entry := m.firstEntry
	for entry != nil {
		if err := util.WriteValue(writer, entry.LogEntryValue, proto.Marshal); err != nil {
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
			index:  snapshotListener.Index,
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

	var prevEntry *LinkedLogEntryValue
	m.firstEntry = nil
	m.lastEntry = nil
	m.indexes = make(map[uint64]*LinkedLogEntryValue)
	for i := 0; i < entryCount; i++ {
		value, err := util.ReadValue(reader, func(data []byte) (*LogEntryValue, error) {
			entry := &LogEntryValue{}
			if err := proto.Unmarshal(data, entry); err != nil {
				return nil, err
			}
			return entry, nil
		})
		if err != nil {
			return err
		}

		linkedEntry := &LinkedLogEntryValue{
			LogEntryValue: value.(*LogEntryValue),
		}
		if m.firstEntry == nil {
			m.firstEntry = linkedEntry
		}
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

// Append appends a value to the end of the log
func (m *Service) Append(value []byte) ([]byte, error) {
	request := &AppendRequest{}
	if err := proto.Unmarshal(value, request); err != nil {
		return nil, err
	}

	oldEntry := m.indexes[request.Index]

	if oldEntry == nil {

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

		// Create a new entry value and set it in the log.
		newEntry := &LinkedLogEntryValue{
			LogEntryValue: &LogEntryValue{
				Index:     index,
				Value:     request.Value,
				Timestamp: m.Timestamp(),
			},
		}
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

		m.sendEvent(&ListenResponse{
			Type:      ListenResponse_APPENDED,
			Index:     newEntry.Index,
			Value:     newEntry.Value,
			Timestamp: newEntry.Timestamp,
		})
		return proto.Marshal(&AppendResponse{
			Status: UpdateStatus_OK,
			Index:  newEntry.Index,
		})
	}

	// If the value is equal to the current value, return a no-op.
	if bytes.Equal(oldEntry.Value, request.Value) {
		return proto.Marshal(&AppendResponse{
			Status: UpdateStatus_NOOP,
			Index:  oldEntry.Index,
		})
	}

	// Create a new entry value and set it in the log.
	newEntry := &LinkedLogEntryValue{
		LogEntryValue: &LogEntryValue{
			Index:     oldEntry.Index,
			Value:     request.Value,
			Timestamp: oldEntry.Timestamp,
		},
		Prev: oldEntry.Prev,
		Next: oldEntry.Next,
	}
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

	// Publish an event to listener streams.
	m.sendEvent(&ListenResponse{
		Type:      ListenResponse_APPENDED,
		Index:     newEntry.Index,
		Value:     newEntry.Value,
		Timestamp: newEntry.Timestamp,
	})

	return proto.Marshal(&AppendResponse{
		Status:    UpdateStatus_OK,
		Index:     newEntry.Index,
		Timestamp: newEntry.Timestamp,
	})
}

// Remove removes a key/value pair from the log
func (m *Service) Remove(bytes []byte) ([]byte, error) {
	request := &RemoveRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	var entry *LinkedLogEntryValue
	if request.Index > 0 {
		entry = m.indexes[request.Index]
	}

	if entry == nil {
		return proto.Marshal(&RemoveResponse{
			Status: UpdateStatus_NOOP,
		})
	}

	// Delete the entry from the log.
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
		Type:      ListenResponse_REMOVED,
		Index:     entry.Index,
		Value:     entry.Value,
		Timestamp: entry.Timestamp,
	})

	return proto.Marshal(&RemoveResponse{
		Status:        UpdateStatus_OK,
		Index:         entry.Index,
		PreviousValue: entry.Value,
		Timestamp:     entry.Timestamp,
	})
}

// Get gets a value from the log
func (m *Service) Get(bytes []byte) ([]byte, error) {
	request := &GetRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	var entry *LinkedLogEntryValue
	var ok bool
	if request.Index > 0 {
		entry, ok = m.indexes[request.Index]
	}

	if !ok {
		return proto.Marshal(&GetResponse{})
	}
	return proto.Marshal(&GetResponse{
		Index:     entry.Index,
		Value:     entry.Value,
		Timestamp: entry.Timestamp,
	})
}

// FirstEntry gets the first entry from the log
func (m *Service) FirstEntry(bytes []byte) ([]byte, error) {
	request := &FirstEntryRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if m.firstEntry == nil {
		return proto.Marshal(&FirstEntryResponse{})
	}
	return proto.Marshal(&FirstEntryResponse{
		Index:     m.firstEntry.Index,
		Value:     m.firstEntry.Value,
		Timestamp: m.firstEntry.Timestamp,
	})
}

// LastEntry gets the last entry from the log
func (m *Service) LastEntry(bytes []byte) ([]byte, error) {
	request := &LastEntryRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if m.lastEntry == nil {
		return proto.Marshal(&LastEntryResponse{})
	}
	return proto.Marshal(&LastEntryResponse{
		Index:     m.lastEntry.Index,
		Value:     m.lastEntry.Value,
		Timestamp: m.lastEntry.Timestamp,
	})
}

// PrevEntry gets the previous entry from the log
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
		Index:     entry.Index,
		Value:     entry.Value,
		Timestamp: entry.Timestamp,
	})
}

// NextEntry gets the next entry from the log
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
		Index:     entry.Index,
		Value:     entry.Value,
		Timestamp: entry.Timestamp,
	})
}

// Exists checks if the log contains an index
func (m *Service) Exists(bytes []byte) ([]byte, error) {
	request := &ContainsIndexRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	_, ok := m.indexes[request.Index]
	return proto.Marshal(&ContainsIndexResponse{
		ContainsIndex: ok,
	})
}

// Size returns the size of the log
func (m *Service) Size(bytes []byte) ([]byte, error) {
	return proto.Marshal(&SizeResponse{
		Size_: int32(len(m.indexes)),
	})
}

// Clear removes all entries from the log
func (m *Service) Clear(value []byte) ([]byte, error) {
	m.indexes = make(map[uint64]*LinkedLogEntryValue)
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
				Type:      ListenResponse_NONE,
				Index:     entry.Index,
				Value:     entry.Value,
				Timestamp: entry.Timestamp,
			})
			if err != nil {
				stream.Error(err)
				continue
			}
			if lis.index > 0 {
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

func (m *Service) sendEvent(event *ListenResponse) {
	bytes, _ := proto.Marshal(event)
	for sessionID, listeners := range m.listeners {

		session := m.Session(sessionID)
		if session != nil {
			for _, listener := range listeners {
				if listener.index > 0 {
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

// Entries returns a stream of entries to the client
func (m *Service) Entries(value []byte, stream primitive.Stream) {
	defer stream.Close()
	entry := m.firstEntry
	for entry != nil {
		stream.Result(proto.Marshal(&EntriesResponse{
			Index:     entry.Index,
			Value:     entry.Value,
			Timestamp: entry.Timestamp,
		}))
		entry = entry.Next
	}
}

type listener struct {
	index  uint64
	stream stream.WriteStream
}
