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
	"github.com/atomix/api/proto/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/errors"
	streams "github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
	"time"
)

// Server is a base server for servers that support sessions
type Server struct {
	Protocol Protocol
	Type     ServiceType
}

// DoCommand submits a command to the service
func (s *Server) DoCommand(ctx context.Context, name string, input []byte, header *headers.RequestHeader) ([]byte, *headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return nil, &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_Command{
			Command: &SessionCommandRequest{
				Context: &SessionCommandContext{
					SessionID:      header.SessionID,
					SequenceNumber: header.RequestID,
				},
				Command: &ServiceCommandRequest{
					Service: &ServiceId{
						Type:      s.Type,
						Name:      header.Primitive.Name,
						Namespace: header.Primitive.Namespace,
					},
					Request: &ServiceCommandRequest_Operation{
						Operation: &ServiceOperationRequest{
							Method: name,
							Value:  input,
						},
					},
				},
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	// Write the request
	if err := partition.Write(ctx, bytes, stream); err != nil {
		return nil, nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, nil, errors.Proto(errors.NewInternal("write channel closed"))
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	commandResponse := sessionResponse.GetCommand()
	responseHeader := &headers.ResponseHeader{
		SessionID:  header.SessionID,
		StreamID:   commandResponse.Context.StreamID,
		ResponseID: commandResponse.Context.Sequence,
		Index:      commandResponse.Context.Index,
		Status:     headers.ResponseStatus(commandResponse.Context.Status),
		Message:    commandResponse.Context.Message,
	}
	return commandResponse.Response.GetOperation().Result, responseHeader, nil
}

// DoCommandStream submits a streaming command to the service
func (s *Server) DoCommandStream(ctx context.Context, name string, input []byte, header *headers.RequestHeader, stream streams.WriteStream) error {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		stream.Value(SessionOutput{
			Header: &headers.ResponseHeader{
				Status: headers.ResponseStatus_NOT_LEADER,
				Leader: partition.Leader(),
			},
		})
		stream.Close()
		return nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_Command{
			Command: &SessionCommandRequest{
				Context: &SessionCommandContext{
					SessionID:      header.SessionID,
					SequenceNumber: header.RequestID,
				},
				Command: &ServiceCommandRequest{
					Service: &ServiceId{
						Type:      s.Type,
						Name:      header.Primitive.Name,
						Namespace: header.Primitive.Namespace,
					},
					Request: &ServiceCommandRequest_Operation{
						Operation: &ServiceOperationRequest{
							Method: name,
							Value:  input,
						},
					},
				},
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return errors.Proto(errors.NewInternal(err.Error()))
	}

	stream = streams.NewEncodingStream(stream, func(value interface{}, err error) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		sessionResponse := &SessionResponse{}
		err = proto.Unmarshal(value.([]byte), sessionResponse)
		if err != nil {
			return SessionOutput{
				Result: streams.Result{
					Error: err,
				},
			}, nil
		}

		commandResponse := sessionResponse.GetCommand()
		responseHeader := &headers.ResponseHeader{
			SessionID:  header.SessionID,
			StreamID:   commandResponse.Context.StreamID,
			ResponseID: commandResponse.Context.Sequence,
			Index:      commandResponse.Context.Index,
			Type:       headers.ResponseType(commandResponse.Context.Type),
			Status:     headers.ResponseStatus(commandResponse.Context.Status),
			Message:    commandResponse.Context.Message,
		}

		var result []byte
		if commandResponse.Response != nil {
			result = commandResponse.Response.GetOperation().Result
		}
		return SessionOutput{
			Header: responseHeader,
			Result: streams.Result{
				Value: result,
			},
		}, nil
	})

	go func() {
		_ = partition.Write(ctx, bytes, stream)
	}()
	return nil
}

// DoQuery submits a query to the service
func (s *Server) DoQuery(ctx context.Context, name string, input []byte, header *headers.RequestHeader) ([]byte, *headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return nil, &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_Query{
			Query: &SessionQueryRequest{
				Context: &SessionQueryContext{
					SessionID:          header.SessionID,
					LastIndex:          header.Index,
					LastSequenceNumber: header.RequestID,
				},
				Query: &ServiceQueryRequest{
					Service: &ServiceId{
						Type:      s.Type,
						Name:      header.Primitive.Name,
						Namespace: header.Primitive.Namespace,
					},
					Request: &ServiceQueryRequest_Operation{
						Operation: &ServiceOperationRequest{
							Method: name,
							Value:  input,
						},
					},
				},
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	// Read the request
	if err := partition.Read(ctx, bytes, stream); err != nil {
		return nil, nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, nil, errors.Proto(errors.NewInternal("write channel closed"))
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	queryResponse := sessionResponse.GetQuery()
	responseHeader := &headers.ResponseHeader{
		SessionID: header.SessionID,
		Index:     queryResponse.Context.Index,
		Status:    headers.ResponseStatus(queryResponse.Context.Status),
		Message:   queryResponse.Context.Message,
	}
	return queryResponse.Response.GetOperation().Result, responseHeader, nil
}

// DoQueryStream submits a streaming query to the service
func (s *Server) DoQueryStream(ctx context.Context, name string, input []byte, header *headers.RequestHeader, stream streams.WriteStream) error {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		stream.Value(SessionOutput{
			Header: &headers.ResponseHeader{
				Status: headers.ResponseStatus_NOT_LEADER,
				Leader: partition.Leader(),
			},
		})
		stream.Close()
		return nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_Query{
			Query: &SessionQueryRequest{
				Context: &SessionQueryContext{
					SessionID:          header.SessionID,
					LastIndex:          header.Index,
					LastSequenceNumber: header.RequestID,
				},
				Query: &ServiceQueryRequest{
					Service: &ServiceId{
						Type:      s.Type,
						Name:      header.Primitive.Name,
						Namespace: header.Primitive.Namespace,
					},
					Request: &ServiceQueryRequest_Operation{
						Operation: &ServiceOperationRequest{
							Method: name,
							Value:  input,
						},
					},
				},
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return errors.Proto(errors.NewInternal(err.Error()))
	}

	stream = streams.NewDecodingStream(stream, func(value interface{}, err error) (interface{}, error) {
		sessionResponse := &SessionResponse{}
		if err := proto.Unmarshal(value.([]byte), sessionResponse); err != nil {
			return nil, errors.Proto(errors.NewInternal(err.Error()))
		}
		queryResponse := sessionResponse.GetQuery()
		responseHeader := &headers.ResponseHeader{
			SessionID: header.SessionID,
			Index:     queryResponse.Context.Index,
			Type:      headers.ResponseType(queryResponse.Context.Type),
			Status:    headers.ResponseStatus(queryResponse.Context.Status),
			Message:   queryResponse.Context.Message,
		}
		var result []byte
		if queryResponse.Response != nil {
			result = queryResponse.Response.GetOperation().Result
		}
		return SessionOutput{
			Header: responseHeader,
			Result: streams.Result{
				Value: result,
			},
		}, nil
	})
	go func() {
		_ = partition.Read(ctx, bytes, stream)
	}()
	return nil
}

// DoMetadata submits a metadata query to the service
func (s *Server) DoMetadata(ctx context.Context, serviceType primitive.PrimitiveType, namespace string, header *headers.RequestHeader) ([]*ServiceId, *headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return nil, &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_Query{
			Query: &SessionQueryRequest{
				Context: &SessionQueryContext{
					SessionID:          header.SessionID,
					LastIndex:          header.Index,
					LastSequenceNumber: header.RequestID,
				},
				Query: &ServiceQueryRequest{
					Request: &ServiceQueryRequest_Metadata{
						Metadata: &ServiceMetadataRequest{
							Type:      ServiceType(serviceType),
							Namespace: namespace,
						},
					},
				},
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	// Read the request
	if err := partition.Read(ctx, bytes, stream); err != nil {
		return nil, nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, nil, errors.NewInternal("read channel closed")
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	queryResponse := sessionResponse.GetQuery()
	responseHeader := &headers.ResponseHeader{
		SessionID: header.SessionID,
		Index:     queryResponse.Context.Index,
		Status:    headers.ResponseStatus(queryResponse.Context.Status),
		Message:   queryResponse.Context.Message,
	}
	return queryResponse.Response.GetMetadata().Services, responseHeader, nil
}

// DoOpenSession opens a new session
func (s *Server) DoOpenSession(ctx context.Context, header *headers.RequestHeader, timeout *time.Duration) (*headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_OpenSession{
			OpenSession: &OpenSessionRequest{
				Timeout: timeout,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	// Write the request
	if err := partition.Write(ctx, bytes, stream); err != nil {
		return nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, errors.Proto(errors.NewInternal("write channel closed"))
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	sessionID := sessionResponse.GetOpenSession().SessionID
	return &headers.ResponseHeader{
		SessionID: sessionID,
		Index:     sessionID,
	}, nil
}

// DoKeepAliveSession keeps a session alive
func (s *Server) DoKeepAliveSession(ctx context.Context, header *headers.RequestHeader) (*headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	streams := make(map[uint64]uint64)
	for _, stream := range header.Streams {
		streams[stream.StreamID] = stream.ResponseID
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_KeepAlive{
			KeepAlive: &KeepAliveRequest{
				SessionID:       header.SessionID,
				CommandSequence: header.RequestID,
				Streams:         streams,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Write the request
	if err := partition.Write(ctx, bytes, stream); err != nil {
		return nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, errors.Proto(errors.NewInternal("write channel closed"))
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}
	return &headers.ResponseHeader{
		SessionID: header.SessionID,
	}, nil
}

// DoCloseSession closes a session
func (s *Server) DoCloseSession(ctx context.Context, header *headers.RequestHeader) (*headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_CloseSession{
			CloseSession: &CloseSessionRequest{
				SessionID: header.SessionID,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	// Write the request
	if err := partition.Write(ctx, bytes, stream); err != nil {
		return nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, errors.Proto(errors.NewInternal("write channel closed"))
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}
	return &headers.ResponseHeader{
		SessionID: header.SessionID,
	}, nil
}

// DoCreateService creates the service
func (s *Server) DoCreateService(ctx context.Context, header *headers.RequestHeader) (*headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_Command{
			Command: &SessionCommandRequest{
				Context: &SessionCommandContext{
					SessionID:      header.SessionID,
					SequenceNumber: header.RequestID,
				},
				Command: &ServiceCommandRequest{
					Service: &ServiceId{
						Type:      s.Type,
						Name:      header.Primitive.Name,
						Namespace: header.Primitive.Namespace,
					},
					Request: &ServiceCommandRequest_Create{
						Create: &ServiceCreateRequest{},
					},
				},
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	// Write the request
	if err := partition.Write(ctx, bytes, stream); err != nil {
		return nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, errors.Proto(errors.NewInternal("write channel closed"))
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	commandResponse := sessionResponse.GetCommand()
	responseHeader := &headers.ResponseHeader{
		SessionID:  header.SessionID,
		StreamID:   commandResponse.Context.StreamID,
		ResponseID: commandResponse.Context.Sequence,
		Index:      commandResponse.Context.Index,
		Status:     headers.ResponseStatus(commandResponse.Context.Status),
		Message:    commandResponse.Context.Message,
	}
	return responseHeader, nil
}

// DoCloseService closes the service
func (s *Server) DoCloseService(ctx context.Context, header *headers.RequestHeader) (*headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_Command{
			Command: &SessionCommandRequest{
				Context: &SessionCommandContext{
					SessionID:      header.SessionID,
					SequenceNumber: header.RequestID,
				},
				Command: &ServiceCommandRequest{
					Service: &ServiceId{
						Type:      s.Type,
						Name:      header.Primitive.Name,
						Namespace: header.Primitive.Namespace,
					},
					Request: &ServiceCommandRequest_Close{
						Close: &ServiceCloseRequest{},
					},
				},
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	// Write the request
	if err := partition.Write(ctx, bytes, stream); err != nil {
		return nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, errors.Proto(errors.NewInternal("write channel closed"))
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	commandResponse := sessionResponse.GetCommand()
	responseHeader := &headers.ResponseHeader{
		SessionID:  header.SessionID,
		StreamID:   commandResponse.Context.StreamID,
		ResponseID: commandResponse.Context.Sequence,
		Index:      commandResponse.Context.Index,
		Status:     headers.ResponseStatus(commandResponse.Context.Status),
		Message:    commandResponse.Context.Message,
	}
	return responseHeader, nil
}

// DoDeleteService deletes the service
func (s *Server) DoDeleteService(ctx context.Context, header *headers.RequestHeader) (*headers.ResponseHeader, error) {
	// If the client requires a leader and is not the leader, return an error
	partition := s.Protocol.Partition(PartitionID(header.Partition))
	if partition.MustLeader() && !partition.IsLeader() {
		return &headers.ResponseHeader{
			Status: headers.ResponseStatus_NOT_LEADER,
			Leader: partition.Leader(),
		}, nil
	}

	sessionRequest := &SessionRequest{
		Request: &SessionRequest_Command{
			Command: &SessionCommandRequest{
				Context: &SessionCommandContext{
					SessionID:      header.SessionID,
					SequenceNumber: header.RequestID,
				},
				Command: &ServiceCommandRequest{
					Service: &ServiceId{
						Type:      s.Type,
						Name:      header.Primitive.Name,
						Namespace: header.Primitive.Namespace,
					},
					Request: &ServiceCommandRequest_Delete{
						Delete: &ServiceDeleteRequest{},
					},
				},
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	// Create a unary stream
	stream := streams.NewUnaryStream()

	// Write the request
	if err := partition.Write(ctx, bytes, stream); err != nil {
		return nil, errors.Proto(err)
	}

	// Wait for the result
	result, ok := stream.Receive()
	if !ok {
		return nil, errors.Proto(errors.NewInternal("write channel closed"))
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, errors.Proto(result.Error)
	}

	sessionResponse := &SessionResponse{}
	err = proto.Unmarshal(result.Value.([]byte), sessionResponse)
	if err != nil {
		return nil, errors.Proto(errors.NewInternal(err.Error()))
	}

	commandResponse := sessionResponse.GetCommand()
	responseHeader := &headers.ResponseHeader{
		SessionID:  header.SessionID,
		StreamID:   commandResponse.Context.StreamID,
		ResponseID: commandResponse.Context.Sequence,
		Index:      commandResponse.Context.Index,
		Status:     headers.ResponseStatus(commandResponse.Context.Status),
		Message:    commandResponse.Context.Message,
	}
	return responseHeader, nil
}

// SessionOutput is a result for session-supporting servers containing session header information
type SessionOutput struct {
	streams.Result
	Header *headers.ResponseHeader
}
