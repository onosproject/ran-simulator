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

package primitive

import (
	"container/list"
	"context"
	api "github.com/atomix/api/proto/atomix/session"
	streams "github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"time"
)

// SessionID is a session identifier
type SessionID uint64

// Session is a service session
type Session interface {
	// ID returns the session identifier
	ID() SessionID

	// Streams returns all open streams
	Streams() []Stream

	// Stream returns a stream by ID
	Stream(id StreamID) Stream

	// StreamsOf returns all open streams for the given operation
	StreamsOf(op OperationID) []Stream
}

// newSessionManager creates a new session manager
func newSessionManager(ctx PartitionContext, timeout *time.Duration) *sessionManager {
	if timeout == nil {
		defaultTimeout := 30 * time.Second
		timeout = &defaultTimeout
	}
	session := &sessionManager{
		id:               SessionID(ctx.Index()),
		timeout:          *timeout,
		lastUpdated:      ctx.Timestamp(),
		ctx:              ctx,
		commandCallbacks: make(map[uint64]func()),
		queryCallbacks:   make(map[uint64]*list.List),
		results:          make(map[uint64]streams.Result),
		services:         make(map[ServiceID]*serviceSession),
	}
	util.SessionEntry(ctx.NodeID(), uint64(session.id)).
		Debug("Session open")
	return session
}

// sessionManager manages the ordering of request and response streams for a single client
type sessionManager struct {
	id               SessionID
	timeout          time.Duration
	lastUpdated      time.Time
	ctx              PartitionContext
	commandSequence  uint64
	commandCallbacks map[uint64]func()
	queryCallbacks   map[uint64]*list.List
	results          map[uint64]streams.Result
	services         map[ServiceID]*serviceSession
}

// getService gets the service session
func (s *sessionManager) getService(id ServiceID) *serviceSession {
	return s.services[id]
}

// addService adds a service session
func (s *sessionManager) addService(id ServiceID) *serviceSession {
	session := newServiceSession(s, id)
	s.services[id] = session
	return session
}

// removeService removes a service session
func (s *sessionManager) removeService(id ServiceID) *serviceSession {
	session := s.services[id]
	delete(s.services, id)
	return session
}

// timedOut returns a boolean indicating whether the session is timed out
func (s *sessionManager) timedOut(time time.Time) bool {
	return s.lastUpdated.UnixNano() > 0 && time.Sub(s.lastUpdated) > s.timeout
}

// getResult gets a unary result
func (s *sessionManager) getUnaryResult(id uint64) (streams.Result, bool) {
	result, ok := s.results[id]
	return result, ok
}

// addUnaryResult adds a unary result
func (s *sessionManager) addUnaryResult(id uint64, result streams.Result) streams.Result {
	bytes, err := proto.Marshal(&SessionResponse{
		Response: &SessionResponse_Command{
			Command: &SessionCommandResponse{
				Context: &SessionResponseContext{
					StreamID: uint64(id),
					Index:    uint64(s.ctx.Index()),
					Sequence: 1,
					Status:   getStatus(result.Error),
					Message:  getMessage(result.Error),
				},
				Response: &ServiceCommandResponse{
					Response: &ServiceCommandResponse_Operation{
						Operation: &ServiceOperationResponse{
							result.Value.([]byte),
						},
					},
				},
			},
		},
	})
	result = streams.Result{
		Value: bytes,
		Error: err,
	}
	s.results[id] = result
	return result
}

// scheduleQuery schedules a query to be executed after the given sequence number
func (s *sessionManager) scheduleQuery(sequenceNumber uint64, f func()) {
	queries, ok := s.queryCallbacks[sequenceNumber]
	if !ok {
		queries = list.New()
		s.queryCallbacks[sequenceNumber] = queries
	}
	queries.PushBack(f)
}

// scheduleCommand schedules a command to be executed at the given sequence number
func (s *sessionManager) scheduleCommand(sequenceNumber uint64, f func()) {
	s.commandCallbacks[sequenceNumber] = f
}

// nextCommandSequence returns the next command sequence number for the session
func (s *sessionManager) nextCommandSequence() uint64 {
	return s.commandSequence + 1
}

// completeCommand completes operations up to the given sequence number and executes commands and
// queries pending for the sequence number to be completed
func (s *sessionManager) completeCommand(sequenceNumber uint64) {
	for i := s.commandSequence + 1; i <= sequenceNumber; i++ {
		s.commandSequence = i
		queries, ok := s.queryCallbacks[i]
		if ok {
			query := queries.Front()
			for query != nil {
				query.Value.(func())()
				query = query.Next()
			}
			delete(s.queryCallbacks, i)
		}

		command, ok := s.commandCallbacks[s.nextCommandSequence()]
		if ok {
			command()
			delete(s.commandCallbacks, i)
		}
	}
}

// close closes the session and completes all its streams
func (s *sessionManager) close() {
	util.SessionEntry(s.ctx.NodeID(), uint64(s.id)).
		Debug("Session closed")
	for _, service := range s.services {
		service.close()
	}
}

// newServiceSession creates a new service session
func newServiceSession(session *sessionManager, service ServiceID) *serviceSession {
	return &serviceSession{
		sessionManager: session,
		service:        service,
		streams:        make(map[StreamID]*sessionStream),
	}
}

// serviceSession manages sessions within a service
type serviceSession struct {
	*sessionManager
	service ServiceID
	streams map[StreamID]*sessionStream
}

func (s *serviceSession) ID() SessionID {
	return s.id
}

func (s *serviceSession) Streams() []Stream {
	streams := make([]Stream, 0, len(s.streams))
	for _, stream := range s.streams {
		streams = append(streams, stream)
	}
	return streams
}

func (s *serviceSession) Stream(id StreamID) Stream {
	return s.streams[id]
}

func (s *serviceSession) StreamsOf(op OperationID) []Stream {
	streams := make([]Stream, 0, len(s.streams))
	for _, stream := range s.streams {
		if stream.OperationID() == op {
			streams = append(streams, stream)
		}
	}
	return streams
}

// addStream adds a stream at the given sequence number
func (s *serviceSession) addStream(id StreamID, op OperationID, outStream streams.WriteStream) Stream {
	stream := &sessionStream{
		id:      id,
		op:      op,
		session: s,
		ctx:     s.ctx,
		stream:  outStream,
		results: list.New(),
	}
	s.streams[id] = stream
	stream.open()
	util.StreamEntry(s.ctx.NodeID(), uint64(s.ID()), uint64(id)).
		Trace("Stream open")
	return stream
}

// getStream returns a stream by the request sequence number
func (s *serviceSession) getStream(id StreamID) *sessionStream {
	return s.streams[id]
}

// ack acknowledges response streams up to the given request sequence number
func (s *serviceSession) ack(id uint64, streams map[uint64]uint64) {
	for responseID := range s.results {
		if responseID > id {
			continue
		}
		delete(s.results, responseID)
	}
	for streamID, stream := range s.streams {
		// If the stream ID is greater than the acknowledged sequence number, skip it
		if uint64(stream.ID()) > id {
			continue
		}

		// If the stream is still held by the client, ack the stream.
		// Otherwise, close the stream.
		streamAck, open := streams[uint64(stream.ID())]
		if open {
			stream.ack(streamAck)
		} else {
			stream.Close()
			delete(s.streams, streamID)
		}
	}
}

// close closes the session and completes all its streams
func (s *serviceSession) close() {
	for _, stream := range s.streams {
		stream.Close()
	}
}

var _ Session = &serviceSession{}

// SessionServer is an implementation of SessionServiceServer for session management
type SessionServer struct {
	*Server
}

// OpenSession opens a new session
func (s *SessionServer) OpenSession(ctx context.Context, request *api.OpenSessionRequest) (*api.OpenSessionResponse, error) {
	log.Tracef("Received OpenSessionRequest %+v", request)
	header, err := s.DoOpenSession(ctx, request.Header, request.Timeout)
	if err != nil {
		return nil, err
	}
	response := &api.OpenSessionResponse{
		Header: header,
	}
	log.Tracef("Sending OpenSessionResponse %+v", response)
	return response, nil
}

// KeepAlive keeps a session alive
func (s *SessionServer) KeepAlive(ctx context.Context, request *api.KeepAliveRequest) (*api.KeepAliveResponse, error) {
	log.Tracef("Received KeepAliveRequest %+v", request)
	header, err := s.DoKeepAliveSession(ctx, request.Header)
	if err != nil {
		return nil, err
	}
	response := &api.KeepAliveResponse{
		Header: header,
	}
	log.Tracef("Sending KeepAliveResponse %+v", response)
	return response, nil
}

// CloseSession closes a session
func (s *SessionServer) CloseSession(ctx context.Context, request *api.CloseSessionRequest) (*api.CloseSessionResponse, error) {
	log.Tracef("Received CloseSessionRequest %+v", request)
	header, err := s.DoCloseSession(ctx, request.Header)
	if err != nil {
		return nil, err
	}
	response := &api.CloseSessionResponse{
		Header: header,
	}
	log.Tracef("Sending CloseSessionResponse %+v", response)
	return response, nil
}
