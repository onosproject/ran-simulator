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

package set

import (
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/golang/protobuf/proto"
	"io"
)

// Service is a state machine for a list primitive
type Service struct {
	primitive.Service
	values map[string]bool
}

// init initializes the list service
func (s *Service) init() {
	s.RegisterUnaryOperation(opSize, s.Size)
	s.RegisterUnaryOperation(opContains, s.Contains)
	s.RegisterUnaryOperation(opAdd, s.Add)
	s.RegisterUnaryOperation(opRemove, s.Remove)
	s.RegisterUnaryOperation(opClear, s.Clear)
	s.RegisterStreamOperation(opEvents, s.Events)
	s.RegisterStreamOperation(opIterate, s.Iterate)
}

// Backup takes a snapshot of the service
func (s *Service) Backup(writer io.Writer) error {
	return util.WriteMap(writer, s.values, func(key string, value bool) ([]byte, error) {
		return []byte(key), nil
	})
}

// Restore restores the service from a snapshot
func (s *Service) Restore(reader io.Reader) error {
	s.values = make(map[string]bool)
	return util.ReadMap(reader, s.values, func(data []byte) (string, bool, error) {
		return string(data), true, nil
	})
}

// Size gets the number of elements in the set
func (s *Service) Size(bytes []byte) ([]byte, error) {
	return proto.Marshal(&SizeResponse{
		Size_: uint32(len(s.values)),
	})
}

// Contains checks whether the set contains an element
func (s *Service) Contains(bytes []byte) ([]byte, error) {
	request := &ContainsRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	_, ok := s.values[request.Value]
	return proto.Marshal(&ContainsResponse{
		Contains: ok,
	})
}

// Add adds an element to the set
func (s *Service) Add(bytes []byte) ([]byte, error) {
	request := &AddRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if _, ok := s.values[request.Value]; !ok {
		s.values[request.Value] = true

		s.sendEvent(&ListenResponse{
			Type:  ListenResponse_ADDED,
			Value: request.Value,
		})

		return proto.Marshal(&AddResponse{
			Added: true,
		})
	}
	return proto.Marshal(&AddResponse{
		Added: false,
	})
}

// Remove removes an element from the set
func (s *Service) Remove(bytes []byte) ([]byte, error) {
	request := &RemoveRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if _, ok := s.values[request.Value]; ok {
		delete(s.values, request.Value)

		s.sendEvent(&ListenResponse{
			Type:  ListenResponse_REMOVED,
			Value: request.Value,
		})

		return proto.Marshal(&RemoveResponse{
			Removed: true,
		})
	}
	return proto.Marshal(&RemoveResponse{
		Removed: false,
	})
}

// Clear removes all elements from the set
func (s *Service) Clear(bytes []byte) ([]byte, error) {
	s.values = make(map[string]bool)
	return proto.Marshal(&ClearResponse{})
}

// Events registers a channel on which to send set change events
func (s *Service) Events(bytes []byte, stream primitive.Stream) {
	request := &ListenRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		stream.Error(err)
		stream.Close()
		return
	}

	if request.Replay {
		for value := range s.values {
			stream.Result(proto.Marshal(&ListenResponse{
				Type:  ListenResponse_NONE,
				Value: value,
			}))
		}
	}
}

// Iterate sends all current set elements on the given channel
func (s *Service) Iterate(bytes []byte, stream primitive.Stream) {
	defer stream.Close()
	for value := range s.values {
		stream.Result(proto.Marshal(&IterateResponse{
			Value: value,
		}))
	}
}

func (s *Service) sendEvent(event *ListenResponse) {
	bytes, err := proto.Marshal(event)
	for _, session := range s.Sessions() {
		for _, stream := range session.StreamsOf(opEvents) {
			stream.Result(bytes, err)
		}
	}
}
