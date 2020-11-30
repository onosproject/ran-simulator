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
	"context"
	"github.com/atomix/api/proto/atomix/headers"
	primitiveapi "github.com/atomix/api/proto/atomix/primitive"
	api "github.com/atomix/api/proto/atomix/session"
	"github.com/atomix/go-client/pkg/client/errors"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"math"
	"sync"
	"time"
)

// SessionOption implements a session option
type SessionOption interface {
	prepare(options *sessionOptions)
}

// WithSessionTimeout returns a session SessionOption to configure the session timeout
func WithSessionTimeout(timeout time.Duration) SessionOption {
	return sessionTimeoutOption{timeout: timeout}
}

type sessionTimeoutOption struct {
	timeout time.Duration
}

func (o sessionTimeoutOption) prepare(options *sessionOptions) {
	options.timeout = o.timeout
}

type sessionOptions struct {
	id      string
	timeout time.Duration
}

// MetadataOption implements a session metadata option
type MetadataOption interface {
	apply(options *metadataOptions)
}

// WithNamespace returns a metadata option limiting a query by namespace
func WithNamespace(namespace string) MetadataOption {
	return &metadataNamespaceOption{
		namespace: namespace,
	}
}

type metadataNamespaceOption struct {
	namespace string
}

func (o *metadataNamespaceOption) apply(options *metadataOptions) {
	options.namespace = o.namespace
}

// WithPrimitiveType returns a metadata option limiting a query by primitive type
func WithPrimitiveType(primitiveType Type) MetadataOption {
	return &metadataPrimitiveTypeOption{
		primitiveType: primitiveType,
	}
}

type metadataPrimitiveTypeOption struct {
	primitiveType Type
}

func (o *metadataPrimitiveTypeOption) apply(options *metadataOptions) {
	options.primitiveType = o.primitiveType
}

type metadataOptions struct {
	namespace     string
	primitiveType Type
}

// NewSession creates a new Session for the given partition
// name is the name of the primitive
// handler is the primitive's session handler
func NewSession(ctx context.Context, partition Partition, opts ...SessionOption) (*Session, error) {
	options := &sessionOptions{
		id:      uuid.New().String(),
		timeout: 30 * time.Second,
	}
	for i := range opts {
		opts[i].prepare(options)
	}
	session := &Session{
		Partition: partition.ID,
		conns:     net.NewConns(partition.Address),
		Timeout:   options.timeout,
		streams:   make(map[uint64]*Stream),
		mu:        sync.RWMutex{},
		ticker:    time.NewTicker(options.timeout / 2),
	}
	if err := session.open(ctx); err != nil {
		return nil, err
	}
	return session, nil
}

// Session maintains the session for a primitive
type Session struct {
	Partition  int
	Timeout    time.Duration
	SessionID  uint64
	conns      *net.Conns
	lastIndex  uint64
	requestID  uint64
	responseID uint64
	streams    map[uint64]*Stream
	mu         sync.RWMutex
	ticker     *time.Ticker
}

// open creates the session and begins keep-alives
func (s *Session) open(ctx context.Context) error {
	err := s.doSession(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		request := &api.OpenSessionRequest{
			Header:  header,
			Timeout: &s.Timeout,
		}
		client := api.NewSessionServiceClient(conn)
		response, err := client.OpenSession(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return err
	}

	go func() {
		for range s.ticker.C {
			_ = s.keepAlive(context.TODO())
		}
	}()
	return nil
}

// keepAlive keeps the session alive
func (s *Session) keepAlive(ctx context.Context) error {
	return s.doSession(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		request := &api.KeepAliveRequest{
			Header: header,
		}
		client := api.NewSessionServiceClient(conn)
		response, err := client.KeepAlive(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
}

// Close closes the session
func (s *Session) Close() error {
	err := s.close(context.TODO())
	s.ticker.Stop()
	return err
}

// close closes the session
func (s *Session) close(ctx context.Context) error {
	return s.doSession(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		request := &api.CloseSessionRequest{
			Header: header,
		}
		client := api.NewSessionServiceClient(conn)
		response, err := client.CloseSession(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
}

func getPrimitiveID(name Name) primitiveapi.PrimitiveId {
	return primitiveapi.PrimitiveId{
		Name:      name.Name,
		Namespace: name.Namespace,
	}
}

// getState gets the header for the current state of the session
func (s *Session) getState(primitive primitiveapi.PrimitiveId) *headers.RequestHeader {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &headers.RequestHeader{
		Primitive: primitive,
		Partition: uint32(s.Partition),
		SessionID: s.SessionID,
		Index:     s.lastIndex,
		RequestID: s.responseID,
		Streams:   s.getStreamHeaders(),
	}
}

// getQueryHeader gets the current read header
func (s *Session) getQueryHeader(primitive primitiveapi.PrimitiveId) *headers.RequestHeader {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &headers.RequestHeader{
		Primitive: primitive,
		Partition: uint32(s.Partition),
		SessionID: s.SessionID,
		Index:     s.lastIndex,
		RequestID: s.requestID,
	}
}

// nextCommandHeader returns the next write header
func (s *Session) nextCommandHeader(primitive primitiveapi.PrimitiveId) *headers.RequestHeader {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestID = s.requestID + 1
	header := &headers.RequestHeader{
		Primitive: primitive,
		Partition: uint32(s.Partition),
		SessionID: s.SessionID,
		Index:     s.lastIndex,
		RequestID: s.requestID,
	}
	return header
}

// nextStreamHeader returns the next write stream and header
func (s *Session) nextStreamHeader(primitive primitiveapi.PrimitiveId) (*Stream, *headers.RequestHeader) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestID = s.requestID + 1
	stream := &Stream{
		ID:      s.requestID,
		session: s,
	}
	s.streams[s.requestID] = stream
	header := &headers.RequestHeader{
		Primitive: primitive,
		Partition: uint32(s.Partition),
		SessionID: s.SessionID,
		Index:     s.lastIndex,
		RequestID: s.requestID,
	}
	return stream, header
}

func (s *Session) doSession(ctx context.Context, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) error {
	header := s.getState(primitiveapi.PrimitiveId{})
	_, err := s.doRequest(ctx, header, func(conn *grpc.ClientConn) (*headers.ResponseHeader, interface{}, error) {
		return f(ctx, conn, header)
	})
	return err
}

// doPrimitive sends a primitive request
func (s *Session) doPrimitive(ctx context.Context, name Name, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) error {
	header := s.nextCommandHeader(getPrimitiveID(name))
	_, err := s.doRequest(ctx, header, func(conn *grpc.ClientConn) (*headers.ResponseHeader, interface{}, error) {
		return f(ctx, conn, header)
	})
	return err
}

// doCreate sends a create session request
func (s *Session) doCreate(ctx context.Context, name Name, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) error {
	return s.doPrimitive(ctx, name, f)
}

// doClose sends a session close request
func (s *Session) doClose(ctx context.Context, name Name, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) error {
	return s.doPrimitive(ctx, name, f)
}

// doQuery sends a session query request
func (s *Session) doQuery(ctx context.Context, name Name, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) (interface{}, error) {
	header := s.getQueryHeader(getPrimitiveID(name))
	return s.doRequest(ctx, header, func(conn *grpc.ClientConn) (*headers.ResponseHeader, interface{}, error) {
		return f(ctx, conn, header)
	})
}

// doCommand sends a session command request
func (s *Session) doCommand(ctx context.Context, name Name, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) (interface{}, error) {
	header := s.nextCommandHeader(getPrimitiveID(name))
	return s.doRequest(ctx, header, func(conn *grpc.ClientConn) (*headers.ResponseHeader, interface{}, error) {
		return f(ctx, conn, header)
	})
}

func (s *Session) doRequest(ctx context.Context, requestHeader *headers.RequestHeader, f func(conn *grpc.ClientConn) (*headers.ResponseHeader, interface{}, error)) (interface{}, error) {
	i := 0
	for {
		conn, err := s.conns.Connect()
		if err != nil {
			return nil, err
		}
		responseHeader, response, err := f(conn)
		if err == nil {
			switch responseHeader.Status {
			case headers.ResponseStatus_OK:
				s.recordResponse(requestHeader, responseHeader)
				return response, nil
			case headers.ResponseStatus_NOT_LEADER:
				s.conns.Reconnect(net.Address(responseHeader.Leader))
				continue
			default:
				s.recordResponse(requestHeader, responseHeader)
				return response, errors.FromHeader(responseHeader)
			}
		} else if err == context.Canceled {
			return nil, errors.NewCanceled(err.Error())
		} else {
			select {
			case <-time.After(time.Duration(math.Max(math.Pow(float64(i), 2), 1000)) * time.Millisecond):
				i++
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}

// doQueryStream sends a session query stream request
func (s *Session) doQueryStream(
	ctx context.Context,
	name Name,
	f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error),
	responseFunc func(interface{}) (*headers.ResponseHeader, interface{}, error)) (<-chan interface{}, error) {
	conn, err := s.conns.Connect()
	if err != nil {
		return nil, err
	}

	requestHeader := s.getQueryHeader(getPrimitiveID(name))
	responses, err := f(ctx, conn, requestHeader)
	if err != nil {
		return nil, err
	}

	handshakeCh := make(chan struct{})
	responseCh := make(chan interface{})
	go s.queryStream(ctx, f, responseFunc, responses, requestHeader, handshakeCh, responseCh)

	select {
	case <-handshakeCh:
		return responseCh, nil
	case <-time.After(15 * time.Second):
		return nil, errors.NewTimeout("handshake timed out")
	}
}

func (s *Session) queryStream(
	ctx context.Context,
	f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error),
	responseFunc func(interface{}) (*headers.ResponseHeader, interface{}, error),
	responses interface{},
	requestHeader *headers.RequestHeader,
	handshakeCh chan<- struct{},
	responseCh chan interface{}) {
	for {
		responseHeader, response, err := responseFunc(responses)
		if err != nil {
			close(responseCh)
			return
		}

		switch responseHeader.Type {
		case headers.ResponseType_OPEN_STREAM:
			close(handshakeCh)
		case headers.ResponseType_CLOSE_STREAM:
			close(responseCh)
			return
		case headers.ResponseType_RESPONSE:
			switch responseHeader.Status {
			case headers.ResponseStatus_OK:
				// Record the response
				s.recordResponse(requestHeader, responseHeader)
				responseCh <- response
			case headers.ResponseStatus_NOT_LEADER:
				s.conns.Reconnect(net.Address(responseHeader.Leader))
				conn, err := s.conns.Connect()
				if err != nil {
					close(responseCh)
				} else {
					responses, err := f(ctx, conn, requestHeader)
					if err != nil {
						close(responseCh)
					} else {
						go s.queryStream(ctx, f, responseFunc, responses, requestHeader, nil, responseCh)
					}
				}
				return
			case headers.ResponseStatus_ERROR:
				close(responseCh)
				return
			}
		}
	}
}

// doCommandStream sends a session command stream request
func (s *Session) doCommandStream(
	ctx context.Context,
	name Name,
	f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error),
	responseFunc func(interface{}) (*headers.ResponseHeader, interface{}, error)) (<-chan interface{}, error) {
	conn, err := s.conns.Connect()
	if err != nil {
		return nil, err
	}

	stream, requestHeader := s.nextStreamHeader(getPrimitiveID(name))
	responses, err := f(ctx, conn, requestHeader)
	if err != nil {
		stream.Close()
		return nil, err
	}

	// Create a goroutine to close the stream when the context is canceled.
	// This will ensure that the server is notified the stream has been closed on the next keep-alive.
	go func() {
		<-ctx.Done()
		stream.Close()
	}()

	handshakeCh := make(chan struct{})
	responseCh := make(chan interface{})
	go s.commandStream(ctx, f, responseFunc, responses, stream, requestHeader, handshakeCh, responseCh)

	select {
	case <-handshakeCh:
		return responseCh, nil
	case <-time.After(15 * time.Second):
		return nil, errors.NewTimeout("handshake timed out")
	}
}

func (s *Session) commandStream(
	ctx context.Context,
	f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error),
	responseFunc func(interface{}) (*headers.ResponseHeader, interface{}, error),
	responses interface{},
	stream *Stream,
	requestHeader *headers.RequestHeader,
	handshakeCh chan<- struct{},
	responseCh chan<- interface{}) {
	for {
		responseHeader, response, err := responseFunc(responses)
		if err != nil {
			close(responseCh)
			stream.Close()
			return
		}

		switch responseHeader.Type {
		case headers.ResponseType_OPEN_STREAM:
			if stream.serialize(responseHeader) && handshakeCh != nil {
				close(handshakeCh)
			}
		case headers.ResponseType_CLOSE_STREAM:
			if stream.serialize(responseHeader) {
				close(responseCh)
				stream.Close()
				return
			}
		case headers.ResponseType_RESPONSE:
			switch responseHeader.Status {
			case headers.ResponseStatus_OK:
				// Record the response
				s.recordResponse(requestHeader, responseHeader)

				// Attempt to serialize the response to the stream and skip the response if serialization failed.
				if stream.serialize(responseHeader) {
					responseCh <- response
				}
			case headers.ResponseStatus_NOT_LEADER:
				s.conns.Reconnect(net.Address(responseHeader.Leader))
				conn, err := s.conns.Connect()
				if err != nil {
					close(responseCh)
					stream.Close()
				} else {
					responses, err := f(ctx, conn, requestHeader)
					if err != nil {
						close(responseCh)
						stream.Close()
					} else {
						go s.commandStream(ctx, f, responseFunc, responses, stream, requestHeader, nil, responseCh)
					}
				}
				return
			case headers.ResponseStatus_ERROR:
				close(responseCh)
				stream.Close()
				return
			}
		}
	}
}

// recordResponse records the index in a response header
func (s *Session) recordResponse(requestHeader *headers.RequestHeader, responseHeader *headers.ResponseHeader) {
	// Use a double-checked lock to avoid locking when multiple responses are received for an index.
	s.mu.RLock()
	if responseHeader.Index > s.lastIndex {
		s.mu.RUnlock()
		s.mu.Lock()

		// If the session ID is set, ensure the session is initialized
		if responseHeader.SessionID > s.SessionID {
			s.SessionID = responseHeader.SessionID
			s.lastIndex = responseHeader.SessionID
		}

		// If the request ID is greater than the highest response ID, update the response ID.
		if requestHeader.RequestID > s.responseID {
			s.responseID = requestHeader.RequestID
		}

		// If the response index has increased, update the last received index
		if responseHeader.Index > s.lastIndex {
			s.lastIndex = responseHeader.Index
		}
		s.mu.Unlock()
	} else {
		s.mu.RUnlock()
	}
}

// deleteStream deletes the given stream from the session
func (s *Session) deleteStream(streamID uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.streams, streamID)
}

// getStreamHeaders returns a slice of headers for all open streams
func (s *Session) getStreamHeaders() []headers.StreamHeader {
	result := make([]headers.StreamHeader, 0, len(s.streams))
	for _, stream := range s.streams {
		if stream.ID <= s.responseID {
			result = append(result, stream.getHeader())
		}
	}
	return result
}

// Stream manages the context for a single response stream within a session
type Stream struct {
	ID         uint64
	session    *Session
	responseID uint64
	mu         sync.RWMutex
}

// getHeader returns the current header for the stream
func (s *Stream) getHeader() headers.StreamHeader {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return headers.StreamHeader{
		StreamID:   s.ID,
		ResponseID: s.responseID,
	}
}

// serialize updates the stream response metadata and returns whether the response was received in sequential order
func (s *Stream) serialize(header *headers.ResponseHeader) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if header.ResponseID == s.responseID+1 {
		s.responseID++
		return true
	}
	return false
}

// Close closes the stream
func (s *Stream) Close() {
	s.session.deleteStream(s.ID)
}
