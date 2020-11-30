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

package election

import (
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/golang/protobuf/proto"
	"io"
	"time"
)

// Service is a state machine for an election primitive
type Service struct {
	primitive.Service
	leader     *ElectionRegistration
	term       uint64
	timestamp  *time.Time
	candidates []*ElectionRegistration
}

// init initializes the election service
func (e *Service) init() {
	e.RegisterUnaryOperation(opEnter, e.Enter)
	e.RegisterUnaryOperation(opWithdraw, e.Withdraw)
	e.RegisterUnaryOperation(opAnoint, e.Anoint)
	e.RegisterUnaryOperation(opPromote, e.Promote)
	e.RegisterUnaryOperation(opEvict, e.Evict)
	e.RegisterUnaryOperation(opGetTerm, e.GetTerm)
	e.RegisterStreamOperation(opEvents, e.Events)
}

// Backup takes a snapshot of the service
func (e *Service) Backup(writer io.Writer) error {
	snapshot := &ElectionSnapshot{
		Term:       e.term,
		Timestamp:  e.timestamp,
		Leader:     e.leader,
		Candidates: e.candidates,
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

	snapshot := &ElectionSnapshot{}
	if err := proto.Unmarshal(bytes, snapshot); err != nil {
		return err
	}
	e.term = snapshot.Term
	e.timestamp = snapshot.Timestamp
	e.leader = snapshot.Leader
	e.candidates = snapshot.Candidates
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
	candidates := make([]*ElectionRegistration, 0, len(e.candidates))
	for _, candidate := range e.candidates {
		if primitive.SessionID(candidate.SessionID) != session.ID() {
			candidates = append(candidates, candidate)
		}
	}

	if len(candidates) != len(e.candidates) {
		e.candidates = candidates

		if primitive.SessionID(e.leader.SessionID) == session.ID() {
			e.leader = nil
			if len(e.candidates) > 0 {
				e.leader = e.candidates[0]
				e.term++
				timestamp := e.Timestamp()
				e.timestamp = &timestamp
			}
		}

		e.sendEvent(&ListenResponse{
			Type: ListenResponse_CHANGED,
			Term: e.getTerm(),
		})
	}
}

// getTerm returns the current election term
func (e *Service) getTerm() *Term {
	var leader string
	if e.leader != nil {
		leader = e.leader.ID
	}
	return &Term{
		ID:         e.term,
		Timestamp:  e.timestamp,
		Leader:     leader,
		Candidates: e.getCandidates(),
	}
}

// getCandidates returns a slice of candidate IDs
func (e *Service) getCandidates() []string {
	candidates := make([]string, len(e.candidates))
	for i, candidate := range e.candidates {
		candidates[i] = candidate.ID
	}
	return candidates
}

// Enter enters a candidate in the election
func (e *Service) Enter(bytes []byte) ([]byte, error) {
	request := &EnterRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	updated := false
	for _, candidate := range e.candidates {
		if candidate.ID == request.ID {
			candidate.SessionID = uint64(e.CurrentSession().ID())
			updated = true
			break
		}
	}

	if !updated {
		reg := &ElectionRegistration{
			ID:        request.ID,
			SessionID: uint64(e.CurrentSession().ID()),
		}

		e.candidates = append(e.candidates, reg)
		if e.leader == nil {
			e.leader = reg
			e.term++
			timestamp := e.Timestamp()
			e.timestamp = &timestamp
		}

		e.sendEvent(&ListenResponse{
			Type: ListenResponse_CHANGED,
			Term: e.getTerm(),
		})
	}

	return proto.Marshal(&EnterResponse{
		Term: e.getTerm(),
	})
}

// Withdraw withdraws a candidate from the election
func (e *Service) Withdraw(bytes []byte) ([]byte, error) {
	request := &WithdrawRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	candidates := make([]*ElectionRegistration, 0, len(e.candidates))
	for _, candidate := range e.candidates {
		if candidate.ID != request.ID {
			candidates = append(candidates, candidate)
		}
	}

	if len(candidates) != len(e.candidates) {
		e.candidates = candidates

		if e.leader.ID == request.ID {
			e.leader = nil
			if len(e.candidates) > 0 {
				e.leader = e.candidates[0]
				e.term++
				timestamp := e.Timestamp()
				e.timestamp = &timestamp
			}
		}

		e.sendEvent(&ListenResponse{
			Type: ListenResponse_CHANGED,
			Term: e.getTerm(),
		})
	}

	return proto.Marshal(&WithdrawResponse{
		Term: e.getTerm(),
	})
}

// Anoint assigns leadership to a candidate
func (e *Service) Anoint(bytes []byte) ([]byte, error) {
	request := &AnointRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if e.leader != nil && e.leader.ID == request.ID {
		return proto.Marshal(&AnointResponse{
			Term: e.getTerm(),
		})
	}

	var leader *ElectionRegistration
	for _, candidate := range e.candidates {
		if candidate.ID == request.ID {
			leader = candidate
			break
		}
	}

	if leader == nil {
		return proto.Marshal(&AnointResponse{
			Term: e.getTerm(),
		})
	}

	candidates := make([]*ElectionRegistration, 0, len(e.candidates))
	candidates = append(candidates, leader)
	for _, candidate := range e.candidates {
		if candidate.ID != request.ID {
			candidates = append(candidates, candidate)
		}
	}

	e.leader = leader
	e.term++
	timestamp := e.Timestamp()
	e.timestamp = &timestamp
	e.candidates = candidates

	e.sendEvent(&ListenResponse{
		Type: ListenResponse_CHANGED,
		Term: e.getTerm(),
	})

	return proto.Marshal(&AnointResponse{
		Term: e.getTerm(),
	})
}

// Promote increases the priority of a candidate
func (e *Service) Promote(bytes []byte) ([]byte, error) {
	request := &PromoteRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if e.leader != nil && e.leader.ID == request.ID {
		return proto.Marshal(&PromoteResponse{
			Term: e.getTerm(),
		})
	}

	var index int
	var promote *ElectionRegistration
	for i, candidate := range e.candidates {
		if candidate.ID == request.ID {
			index = i
			promote = candidate
			break
		}
	}

	if promote == nil {
		return proto.Marshal(&PromoteResponse{
			Term: e.getTerm(),
		})
	}

	candidates := make([]*ElectionRegistration, len(e.candidates))
	for i, candidate := range e.candidates {
		if i < index-1 {
			candidates[i] = candidate
		} else if i == index-1 {
			candidates[i] = promote
		} else if i == index {
			candidates[i] = e.candidates[i-1]
		} else {
			candidates[i] = candidate
		}
	}

	leader := candidates[0]
	if e.leader.ID != leader.ID {
		e.leader = leader
		e.term++
		timestamp := e.Timestamp()
		e.timestamp = &timestamp
	}
	e.candidates = candidates

	e.sendEvent(&ListenResponse{
		Type: ListenResponse_CHANGED,
		Term: e.getTerm(),
	})

	return proto.Marshal(&AnointResponse{
		Term: e.getTerm(),
	})
}

// Evict removes a candidate from the election
func (e *Service) Evict(bytes []byte) ([]byte, error) {
	request := &EvictRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	candidates := make([]*ElectionRegistration, 0, len(e.candidates))
	for _, candidate := range e.candidates {
		if candidate.ID != request.ID {
			candidates = append(candidates, candidate)
		}
	}

	if len(candidates) != len(e.candidates) {
		e.candidates = candidates

		if e.leader.ID == request.ID {
			e.leader = nil
			if len(e.candidates) > 0 {
				e.leader = e.candidates[0]
				e.term++
				timestamp := e.Timestamp()
				e.timestamp = &timestamp
			}
		}

		e.sendEvent(&ListenResponse{
			Type: ListenResponse_CHANGED,
			Term: e.getTerm(),
		})
	}

	return proto.Marshal(&WithdrawResponse{
		Term: e.getTerm(),
	})
}

// GetTerm gets the current election term
func (e *Service) GetTerm(bytes []byte) ([]byte, error) {
	return proto.Marshal(&GetTermResponse{
		Term: e.getTerm(),
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
