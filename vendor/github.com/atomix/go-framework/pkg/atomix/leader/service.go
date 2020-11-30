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

package leader

import (
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/golang/protobuf/proto"
	"io"
)

// Service is a state machine for an election primitive
type Service struct {
	primitive.Service
	leader       *LatchParticipant
	latch        uint64
	participants []*LatchParticipant
}

// init initializes the election service
func (e *Service) init() {
	e.RegisterUnaryOperation(opLatch, e.Latch)
	e.RegisterUnaryOperation(opGetLatch, e.GetLatch)
	e.RegisterStreamOperation(opEvents, e.Events)
}

// Backup takes a snapshot of the service
func (e *Service) Backup(writer io.Writer) error {
	snapshot := &LatchSnapshot{
		Latch:        e.latch,
		Leader:       e.leader,
		Participants: e.participants,
	}
	bytes, err := proto.Marshal(snapshot)
	if err != nil {
		return err
	}
	return util.WriteBytes(writer, bytes)
}

// Restore restores the service from a snapshot
func (e *Service) Restore(reader io.Reader) error {
	bytes, err := util.ReadBytes(reader)
	if err != nil {
		return err
	}

	snapshot := &LatchSnapshot{}
	if err := proto.Unmarshal(bytes, snapshot); err != nil {
		return err
	}
	e.latch = snapshot.Latch
	e.leader = snapshot.Leader
	e.participants = snapshot.Participants
	return nil
}

// SessionExpired is called when a session is expired by the server
func (e *Service) SessionExpired(session primitive.Session) {
	e.close(session)
}

// SessionClosed is called when a session is closed by the client
func (e *Service) SessionClosed(session primitive.Session) {
	e.close(session)
}

// close elects a new leader when a session is closed
func (e *Service) close(session primitive.Session) {
	candidates := make([]*LatchParticipant, 0, len(e.participants))
	for _, candidate := range e.participants {
		if primitive.SessionID(candidate.SessionID) != session.ID() {
			candidates = append(candidates, candidate)
		}
	}

	if len(candidates) != len(e.participants) {
		e.participants = candidates

		if primitive.SessionID(e.leader.SessionID) == session.ID() {
			e.leader = nil
			if len(e.participants) > 0 {
				e.leader = e.participants[0]
				e.latch++
			}
		}

		e.sendEvent(&ListenResponse{
			Type:  ListenResponse_CHANGED,
			Latch: e.getLatch(),
		})
	}
}

// getLatch returns the current election latch
func (e *Service) getLatch() *Latch {
	var leader string
	if e.leader != nil {
		leader = e.leader.ID
	}
	return &Latch{
		ID:           e.latch,
		Leader:       leader,
		Participants: e.getParticipants(),
	}
}

// getParticipants returns a slice of candidate IDs
func (e *Service) getParticipants() []string {
	candidates := make([]string, len(e.participants))
	for i, candidate := range e.participants {
		candidates[i] = candidate.ID
	}
	return candidates
}

// Latch attempts to acquire the latch
func (e *Service) Latch(bytes []byte) ([]byte, error) {
	request := &LatchRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	reg := &LatchParticipant{
		ID:        request.ID,
		SessionID: uint64(e.CurrentSession().ID()),
	}

	e.participants = append(e.participants, reg)
	if e.leader == nil {
		e.leader = reg
		e.latch++
	}

	e.sendEvent(&ListenResponse{
		Type:  ListenResponse_CHANGED,
		Latch: e.getLatch(),
	})

	return proto.Marshal(&LatchResponse{
		Latch: e.getLatch(),
	})
}

// GetLatch gets the current latch
func (e *Service) GetLatch(bytes []byte) ([]byte, error) {
	return proto.Marshal(&GetResponse{
		Latch: e.getLatch(),
	})
}

// Events registers the given channel to receive election events
func (e *Service) Events(bytes []byte, stream primitive.Stream) {
	// Keep the stream open for events
}

func (e *Service) sendEvent(event *ListenResponse) {
	bytes, err := proto.Marshal(event)
	for _, session := range e.Sessions() {
		for _, stream := range session.StreamsOf(opEvents) {
			stream.Result(bytes, err)
		}
	}
}
