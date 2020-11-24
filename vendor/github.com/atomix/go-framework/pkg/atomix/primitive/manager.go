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
	"encoding/binary"
	"fmt"
	"github.com/atomix/api/proto/atomix/primitive"
	streams "github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/gogo/protobuf/proto"
	"io"
	"time"
)

// NewManager returns an initialized Manager
func NewManager(registry Registry, context PartitionContext) *Manager {
	return &Manager{
		registry:  registry,
		context:   context,
		scheduler: newScheduler(),
		sessions:  make(map[SessionID]*sessionManager),
		services:  make(map[ServiceID]Service),
	}
}

// Manager is a Manager implementation for primitives that support sessions
type Manager struct {
	registry  Registry
	context   PartitionContext
	sessions  map[SessionID]*sessionManager
	services  map[ServiceID]Service
	scheduler *scheduler
}

// Snapshot takes a snapshot of the service
func (m *Manager) Snapshot(writer io.Writer) error {
	if err := m.snapshotSessions(writer); err != nil {
		return err
	}
	if err := m.snapshotServices(writer); err != nil {
		return err
	}
	return nil
}

func (m *Manager) snapshotSessions(writer io.Writer) error {
	return util.WriteMap(writer, m.sessions, func(id SessionID, session *sessionManager) ([]byte, error) {
		services := make([]*SessionServiceSnapshot, 0, len(session.services))
		for _, service := range session.services {
			streams := make([]*SessionStreamSnapshot, 0, len(service.streams))
			for _, stream := range service.streams {
				streams = append(streams, &SessionStreamSnapshot{
					StreamId:       uint64(stream.id),
					Type:           string(stream.op),
					SequenceNumber: stream.responseID,
					LastCompleted:  stream.completeID,
				})
			}
			services = append(services, &SessionServiceSnapshot{
				ServiceId: ServiceId(service.service),
				Streams:   streams,
			})
		}
		snapshot := &SessionSnapshot{
			SessionID:       uint64(session.id),
			Timeout:         session.timeout,
			Timestamp:       session.lastUpdated,
			CommandSequence: session.commandSequence,
			Services:        services,
		}
		return proto.Marshal(snapshot)
	})
}

func (m *Manager) snapshotServices(writer io.Writer) error {
	count := make([]byte, 4)
	binary.BigEndian.PutUint32(count, uint32(len(m.services)))
	_, err := writer.Write(count)
	if err != nil {
		return err
	}

	for id, service := range m.services {
		serviceID := ServiceId(id)
		bytes, err := proto.Marshal(&serviceID)
		if err != nil {
			return err
		}

		length := make([]byte, 4)
		binary.BigEndian.PutUint32(length, uint32(len(bytes)))

		_, err = writer.Write(length)
		if err != nil {
			return err
		}

		_, err = writer.Write(bytes)
		if err != nil {
			return err
		}

		err = service.Backup(writer)
		if err != nil {
			return err
		}
	}
	return nil
}

// Install installs a snapshot of the service
func (m *Manager) Install(reader io.Reader) error {
	if err := m.installSessions(reader); err != nil {
		return err
	}
	if err := m.installServices(reader); err != nil {
		return err
	}
	return nil
}

func (m *Manager) installSessions(reader io.Reader) error {
	m.sessions = make(map[SessionID]*sessionManager)
	return util.ReadMap(reader, m.sessions, func(data []byte) (SessionID, *sessionManager, error) {
		snapshot := &SessionSnapshot{}
		if err := proto.Unmarshal(data, snapshot); err != nil {
			return 0, nil, err
		}

		sessionManager := &sessionManager{
			id:               SessionID(snapshot.SessionID),
			timeout:          time.Duration(snapshot.Timeout),
			lastUpdated:      snapshot.Timestamp,
			ctx:              m.context,
			commandSequence:  snapshot.CommandSequence,
			commandCallbacks: make(map[uint64]func()),
			queryCallbacks:   make(map[uint64]*list.List),
			results:          make(map[uint64]streams.Result),
			services:         make(map[ServiceID]*serviceSession),
		}

		for _, service := range snapshot.Services {
			session := &serviceSession{
				sessionManager: sessionManager,
				service:        ServiceID(service.ServiceId),
				streams:        make(map[StreamID]*sessionStream),
			}

			for _, stream := range service.Streams {
				session.streams[StreamID(stream.StreamId)] = &sessionStream{
					id:         StreamID(stream.StreamId),
					op:         OperationID(stream.Type),
					session:    session,
					responseID: stream.SequenceNumber,
					completeID: stream.LastCompleted,
					ctx:        m.context,
					stream:     streams.NewNilStream(),
					results:    list.New(),
				}
			}
			sessionManager.services[ServiceID(service.ServiceId)] = session
		}
		return sessionManager.id, sessionManager, nil
	})
}

func (m *Manager) installServices(reader io.Reader) error {
	services := make(map[ServiceID]Service)

	countBytes := make([]byte, 4)
	n, err := reader.Read(countBytes)
	if err != nil {
		return err
	} else if n <= 0 {
		return nil
	}

	lengthBytes := make([]byte, 4)
	count := int(binary.BigEndian.Uint32(countBytes))
	for i := 0; i < count; i++ {
		n, err = reader.Read(lengthBytes)
		if err != nil {
			return err
		}
		if n > 0 {
			length := binary.BigEndian.Uint32(lengthBytes)
			bytes := make([]byte, length)
			_, err = reader.Read(bytes)
			if err != nil {
				return err
			}

			serviceID := ServiceId{}
			if err = proto.Unmarshal(bytes, &serviceID); err != nil {
				return err
			}
			primitive := m.registry.GetPrimitive(primitive.PrimitiveType(serviceID.Type))
			service := primitive.NewService(m.scheduler, newServiceContext(m.context, ServiceID(serviceID)))
			services[ServiceID(serviceID)] = service
			if err := service.Restore(reader); err != nil {
				return err
			}
		}
	}
	m.services = services

	for _, sessionManager := range m.sessions {
		for serviceID, session := range sessionManager.services {
			service, ok := m.services[serviceID]
			if ok {
				service.addSession(session)
			}
		}
	}
	return nil
}

// Command handles a service command
func (m *Manager) Command(bytes []byte, stream streams.WriteStream) {
	request := &SessionRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		stream.Error(err)
		stream.Close()
	} else {
		m.scheduler.runScheduledTasks(m.context.Timestamp())

		switch r := request.Request.(type) {
		case *SessionRequest_Command:
			m.applyCommand(r.Command, stream)
		case *SessionRequest_OpenSession:
			m.applyOpenSession(r.OpenSession, stream)
		case *SessionRequest_KeepAlive:
			m.applyKeepAlive(r.KeepAlive, stream)
		case *SessionRequest_CloseSession:
			m.applyCloseSession(r.CloseSession, stream)
		}

		m.scheduler.runImmediateTasks()
		m.scheduler.runIndex(m.context.Index())
	}
}

func (m *Manager) applyCommand(request *SessionCommandRequest, stream streams.WriteStream) {
	sessionManager, ok := m.sessions[SessionID(request.Context.SessionID)]
	if !ok {
		util.SessionEntry(m.context.NodeID(), request.Context.SessionID).
			Warn("Unknown session")
		stream.Error(fmt.Errorf("unknown session %d", request.Context.SessionID))
		stream.Close()
	} else {
		sequenceNumber := request.Context.SequenceNumber
		if sequenceNumber != 0 && sequenceNumber <= sessionManager.commandSequence {
			serviceID := ServiceID(*request.Command.Service)

			session := sessionManager.getService(serviceID)
			if session == nil {
				stream.Error(fmt.Errorf("no open session for service %s", serviceID))
				stream.Close()
				return
			}

			result, ok := session.getUnaryResult(sequenceNumber)
			if ok {
				stream.Send(result)
				stream.Close()
			} else {
				streamCtx := session.getStream(StreamID(sequenceNumber))
				if streamCtx != nil {
					streamCtx.replay(stream)
				} else {
					stream.Error(fmt.Errorf("sequence number %d has already been acknowledged", sequenceNumber))
					stream.Close()
				}
			}
		} else if sequenceNumber > sessionManager.nextCommandSequence() {
			sessionManager.scheduleCommand(sequenceNumber, func() {
				util.SessionEntry(m.context.NodeID(), request.Context.SessionID).
					Tracef("Executing command %d", sequenceNumber)
				m.applySessionCommand(request, sessionManager, stream)
			})
		} else {
			util.SessionEntry(m.context.NodeID(), request.Context.SessionID).
				Tracef("Executing command %d", sequenceNumber)
			m.applySessionCommand(request, sessionManager, stream)
		}
	}
}

func (m *Manager) applySessionCommand(request *SessionCommandRequest, session *sessionManager, stream streams.WriteStream) {
	m.applyServiceCommand(request.Command, request.Context, session, stream)
	session.completeCommand(request.Context.SequenceNumber)
}

func (m *Manager) applyServiceCommand(request *ServiceCommandRequest, context *SessionCommandContext, sessionManager *sessionManager, stream streams.WriteStream) {
	switch request.Request.(type) {
	case *ServiceCommandRequest_Operation:
		m.applyServiceCommandOperation(request, context, sessionManager, stream)
	case *ServiceCommandRequest_Create:
		m.applyServiceCommandCreate(request, context, sessionManager, stream)
	case *ServiceCommandRequest_Close:
		m.applyServiceCommandClose(request, context, sessionManager, stream)
	case *ServiceCommandRequest_Delete:
		m.applyServiceCommandDelete(request, context, sessionManager, stream)
	default:
		stream.Error(fmt.Errorf("unknown service command"))
		stream.Close()
	}
}

func (m *Manager) applyServiceCommandOperation(request *ServiceCommandRequest, context *SessionCommandContext, sessionManager *sessionManager, stream streams.WriteStream) {
	serviceID := ServiceID(*request.Service)

	service, ok := m.services[serviceID]
	if !ok {
		stream.Error(fmt.Errorf("unknown service %s", serviceID))
		stream.Close()
		return
	}

	session := sessionManager.getService(serviceID)
	if session == nil {
		stream.Error(fmt.Errorf("no open session for service %s", serviceID))
		stream.Close()
		return
	}

	operationID := OperationID(request.GetOperation().Method)
	service.setCurrentSession(session)
	service.setCurrentOperation(operationID)

	operation := service.GetOperation(operationID)
	if unaryOp, ok := operation.(UnaryOperation); ok {
		output, err := unaryOp.Execute(request.GetOperation().Value)
		result := session.addUnaryResult(context.SequenceNumber, streams.Result{
			Value: output,
			Error: err,
		})
		stream.Send(result)
		stream.Close()
	} else if streamOp, ok := operation.(StreamingOperation); ok {
		streamCtx := session.addStream(StreamID(context.SequenceNumber), operationID, stream)
		streamOp.Execute(request.GetOperation().Value, streamCtx)
	} else {
		stream.Close()
	}
}

func (m *Manager) applyServiceCommandCreate(request *ServiceCommandRequest, context *SessionCommandContext, sessionManager *sessionManager, stream streams.WriteStream) {
	defer stream.Close()

	serviceID := ServiceID(*request.Service)

	service, ok := m.services[serviceID]
	if !ok {
		primitive := m.registry.GetPrimitive(primitive.PrimitiveType(request.Service.Type))
		if primitive == nil {
			stream.Result(proto.Marshal(&SessionResponse{
				Response: &SessionResponse_Command{
					Command: &SessionCommandResponse{
						Context: &SessionResponseContext{
							Index:    uint64(m.context.Index()),
							Sequence: context.SequenceNumber,
							Type:     SessionResponseType_RESPONSE,
							Status:   SessionResponseStatus_INVALID,
							Message:  fmt.Sprintf("unknown primitive type '%s'", request.Service.Type),
						},
						Response: &ServiceCommandResponse{
							Response: &ServiceCommandResponse_Create{
								Create: &ServiceCreateResponse{},
							},
						},
					},
				},
			}))
			return
		}
		service = primitive.NewService(m.scheduler, newServiceContext(m.context, serviceID))
		m.services[serviceID] = service
	}

	session := sessionManager.getService(serviceID)
	if session == nil {
		session = sessionManager.addService(serviceID)
		service.addSession(session)
		service.setCurrentSession(session)
		if open, ok := service.(SessionOpenService); ok {
			open.SessionOpen(session)
		}
	}

	stream.Result(proto.Marshal(&SessionResponse{
		Response: &SessionResponse_Command{
			Command: &SessionCommandResponse{
				Context: &SessionResponseContext{
					Index:    uint64(m.context.Index()),
					Sequence: context.SequenceNumber,
					Type:     SessionResponseType_RESPONSE,
					Status:   SessionResponseStatus_OK,
				},
				Response: &ServiceCommandResponse{
					Response: &ServiceCommandResponse_Create{
						Create: &ServiceCreateResponse{},
					},
				},
			},
		},
	}))
}

func (m *Manager) applyServiceCommandClose(request *ServiceCommandRequest, context *SessionCommandContext, sessionManager *sessionManager, stream streams.WriteStream) {
	serviceID := ServiceID(*request.Service)

	service, ok := m.services[serviceID]
	if ok {
		session := sessionManager.removeService(serviceID)
		if session != nil {
			service.setCurrentSession(session)
			service.removeSession(session)
			if closed, ok := service.(SessionClosedService); ok {
				closed.SessionClosed(session)
			}
		}
	}

	stream.Result(proto.Marshal(&SessionResponse{
		Response: &SessionResponse_Command{
			Command: &SessionCommandResponse{
				Context: &SessionResponseContext{
					Index:    uint64(m.context.Index()),
					Sequence: context.SequenceNumber,
					Type:     SessionResponseType_RESPONSE,
					Status:   SessionResponseStatus_OK,
				},
				Response: &ServiceCommandResponse{
					Response: &ServiceCommandResponse_Close{
						Close: &ServiceCloseResponse{},
					},
				},
			},
		},
	}))
	stream.Close()
}

func (m *Manager) applyServiceCommandDelete(request *ServiceCommandRequest, context *SessionCommandContext, session *sessionManager, stream streams.WriteStream) {
	defer stream.Close()

	serviceID := ServiceID(*request.Service)

	_, ok := m.services[serviceID]
	if !ok {
		stream.Result(proto.Marshal(&SessionResponse{
			Response: &SessionResponse_Command{
				Command: &SessionCommandResponse{
					Context: &SessionResponseContext{
						Index:    uint64(m.context.Index()),
						Sequence: context.SequenceNumber,
						Type:     SessionResponseType_RESPONSE,
						Status:   SessionResponseStatus_NOT_FOUND,
						Message:  fmt.Sprintf("unknown service '%s.%s'", serviceID.Namespace, serviceID.Name),
					},
					Response: &ServiceCommandResponse{
						Response: &ServiceCommandResponse_Delete{
							Delete: &ServiceDeleteResponse{},
						},
					},
				},
			},
		}))
	} else {
		delete(m.services, serviceID)
		for _, session := range m.sessions {
			session.removeService(serviceID)
		}

		stream.Result(proto.Marshal(&SessionResponse{
			Response: &SessionResponse_Command{
				Command: &SessionCommandResponse{
					Context: &SessionResponseContext{
						Index:    uint64(m.context.Index()),
						Sequence: context.SequenceNumber,
						Type:     SessionResponseType_RESPONSE,
						Status:   SessionResponseStatus_OK,
					},
					Response: &ServiceCommandResponse{
						Response: &ServiceCommandResponse_Delete{
							Delete: &ServiceDeleteResponse{},
						},
					},
				},
			},
		}))
	}
}

func (m *Manager) applyOpenSession(request *OpenSessionRequest, stream streams.WriteStream) {
	session := newSessionManager(m.context, request.Timeout)
	m.sessions[session.id] = session
	stream.Result(proto.Marshal(&SessionResponse{
		Response: &SessionResponse_OpenSession{
			OpenSession: &OpenSessionResponse{
				SessionID: uint64(session.id),
			},
		},
	}))
	stream.Close()
}

// applyKeepAlive applies a KeepAliveRequest to the service
func (m *Manager) applyKeepAlive(request *KeepAliveRequest, stream streams.WriteStream) {
	session, ok := m.sessions[SessionID(request.SessionID)]
	if !ok {
		util.SessionEntry(m.context.NodeID(), request.SessionID).
			Warn("Unknown session")
		stream.Error(fmt.Errorf("unknown session %d", request.SessionID))
	} else {
		util.SessionEntry(m.context.NodeID(), request.SessionID).
			Tracef("Recording keep-alive %v", request)

		// Update the session's last updated timestamp to prevent it from expiring
		session.lastUpdated = m.context.Timestamp()

		// Clear the results up to the given command sequence number
		for _, service := range session.services {
			service.ack(request.CommandSequence, request.Streams)
		}

		// Expire sessions that have not been kept alive
		m.expireSessions()

		// Send the response
		stream.Result(proto.Marshal(&SessionResponse{
			Response: &SessionResponse_KeepAlive{
				KeepAlive: &KeepAliveResponse{},
			},
		}))
	}
	stream.Close()
}

// expireSessions expires sessions that have not been kept alive within their timeout
func (m *Manager) expireSessions() {
	for id, sessionManager := range m.sessions {
		if sessionManager.timedOut(m.context.Timestamp()) {
			sessionManager.close()
			delete(m.sessions, id)
			for _, session := range sessionManager.services {
				service, ok := m.services[session.service]
				if ok {
					service.removeSession(session)
					if expired, ok := service.(SessionExpiredService); ok {
						expired.SessionExpired(session)
					}
				}
			}
		}
	}
}

func (m *Manager) applyCloseSession(request *CloseSessionRequest, stream streams.WriteStream) {
	sessionManager, ok := m.sessions[SessionID(request.SessionID)]
	if !ok {
		util.SessionEntry(m.context.NodeID(), request.SessionID).
			Warn("Unknown session")
		stream.Error(fmt.Errorf("unknown session %d", request.SessionID))
	} else {
		// Close the session and notify the service.
		delete(m.sessions, sessionManager.id)
		sessionManager.close()
		for _, session := range sessionManager.services {
			service, ok := m.services[session.service]
			if ok {
				service.removeSession(session)
				if expired, ok := service.(SessionExpiredService); ok {
					expired.SessionExpired(session)
				}
			}
		}

		// Send the response
		stream.Result(proto.Marshal(&SessionResponse{
			Response: &SessionResponse_CloseSession{
				CloseSession: &CloseSessionResponse{},
			},
		}))
	}
	stream.Close()
}

// Query handles a service query
func (m *Manager) Query(bytes []byte, stream streams.WriteStream) {
	request := &SessionRequest{}
	err := proto.Unmarshal(bytes, request)
	if err != nil {
		stream.Error(err)
		stream.Close()
	} else {
		query := request.GetQuery()
		if Index(query.Context.LastIndex) > m.context.Index() {
			util.SessionEntry(m.context.NodeID(), query.Context.SessionID).
				Tracef("Query index %d greater than last index %d", query.Context.LastIndex, m.context.Index())
			m.scheduler.ScheduleIndex(Index(query.Context.LastIndex), func() {
				m.sequenceQuery(query, stream)
			})
		} else {
			util.SessionEntry(m.context.NodeID(), query.Context.SessionID).
				Tracef("Sequencing query %d <= %d", query.Context.LastIndex, m.context.Index())
			m.sequenceQuery(query, stream)
		}
	}
}

func (m *Manager) sequenceQuery(request *SessionQueryRequest, stream streams.WriteStream) {
	sessionManager, ok := m.sessions[SessionID(request.Context.SessionID)]
	if !ok {
		util.SessionEntry(m.context.NodeID(), request.Context.SessionID).
			Warn("Unknown session")
		stream.Error(fmt.Errorf("unknown session %d", request.Context.SessionID))
		stream.Close()
	} else {
		sequenceNumber := request.Context.LastSequenceNumber
		if sequenceNumber > sessionManager.commandSequence {
			util.SessionEntry(m.context.NodeID(), request.Context.SessionID).
				Tracef("Query ID %d greater than last ID %d", sequenceNumber, sessionManager.commandSequence)
			sessionManager.scheduleQuery(sequenceNumber, func() {
				util.SessionEntry(m.context.NodeID(), request.Context.SessionID).
					Tracef("Executing query %d", sequenceNumber)
				m.applyServiceQuery(request.Query, request.Context, sessionManager, stream)
			})
		} else {
			util.SessionEntry(m.context.NodeID(), request.Context.SessionID).
				Tracef("Executing query %d", sequenceNumber)
			m.applyServiceQuery(request.Query, request.Context, sessionManager, stream)
		}
	}
}

func (m *Manager) applyServiceQuery(request *ServiceQueryRequest, context *SessionQueryContext, sessionManager *sessionManager, stream streams.WriteStream) {
	switch request.Request.(type) {
	case *ServiceQueryRequest_Operation:
		m.applyServiceQueryOperation(request, context, sessionManager, stream)
	case *ServiceQueryRequest_Metadata:
		m.applyServiceQueryMetadata(request, context, sessionManager, stream)
	default:
		stream.Error(fmt.Errorf("unknown service query"))
		stream.Close()
	}
}

func (m *Manager) applyServiceQueryOperation(request *ServiceQueryRequest, context *SessionQueryContext, sessionManager *sessionManager, stream streams.WriteStream) {
	serviceID := ServiceID(*request.Service)

	service, ok := m.services[serviceID]

	// If the service does not exist, reject the operation
	if !ok {
		stream.Error(fmt.Errorf("unknown service %s", serviceID))
		stream.Close()
		return
	}

	session := sessionManager.getService(serviceID)
	if session == nil {
		stream.Error(fmt.Errorf("no open session for service %s", serviceID))
		stream.Close()
		return
	}

	// Set the current session on the service
	operationID := OperationID(request.GetOperation().Method)
	service.setCurrentSession(session)
	service.setCurrentOperation(operationID)

	// Get the service operation
	operation := service.GetOperation(operationID)
	if operation == nil {
		stream.Error(fmt.Errorf("unknown operation: %s", request.GetOperation().Method))
		stream.Close()
		return
	}

	index := m.context.Index()
	responseStream := streams.NewEncodingStream(stream, func(value interface{}, err error) (interface{}, error) {
		return proto.Marshal(&SessionResponse{
			Response: &SessionResponse_Query{
				Query: &SessionQueryResponse{
					Context: &SessionResponseContext{
						Index:    uint64(index),
						Sequence: context.LastSequenceNumber,
						Status:   getStatus(err),
						Message:  getMessage(err),
					},
					Response: &ServiceQueryResponse{
						Response: &ServiceQueryResponse_Operation{
							Operation: &ServiceOperationResponse{
								Result: value.([]byte),
							},
						},
					},
				},
			},
		})
	})

	if unaryOp, ok := operation.(UnaryOperation); ok {
		responseStream.Result(unaryOp.Execute(request.GetOperation().Value))
		responseStream.Close()
	} else if streamOp, ok := operation.(StreamingOperation); ok {
		stream.Result(proto.Marshal(&SessionResponse{
			Response: &SessionResponse_Query{
				Query: &SessionQueryResponse{
					Context: &SessionResponseContext{
						Index:  uint64(index),
						Type:   SessionResponseType_OPEN_STREAM,
						Status: SessionResponseStatus_OK,
					},
				},
			},
		}))

		responseStream = streams.NewCloserStream(responseStream, func(_ streams.WriteStream) {
			stream.Result(proto.Marshal(&SessionResponse{
				Response: &SessionResponse_Query{
					Query: &SessionQueryResponse{
						Context: &SessionResponseContext{
							Index:  uint64(index),
							Type:   SessionResponseType_CLOSE_STREAM,
							Status: SessionResponseStatus_OK,
						},
					},
				},
			}))
		})

		queryStream := &queryStream{
			WriteStream: responseStream,
			id:          StreamID(m.context.Index()),
			op:          operationID,
			session:     session,
		}

		streamOp.Execute(request.GetOperation().Value, queryStream)
	} else {
		stream.Close()
	}
}

func (m *Manager) applyServiceQueryMetadata(request *ServiceQueryRequest, context *SessionQueryContext, session *sessionManager, stream streams.WriteStream) {
	defer stream.Close()

	services := []*ServiceId{}
	serviceType := request.GetMetadata().Type
	namespace := request.GetMetadata().Namespace
	for name, service := range m.services {
		if (serviceType == 0 || service.ServiceType() == serviceType) && (namespace == "" || name.Namespace == namespace) {
			services = append(services, &ServiceId{
				Type:      service.ServiceType(),
				Namespace: name.Namespace,
				Name:      name.Name,
			})
		}
	}

	stream.Result(proto.Marshal(&SessionResponse{
		Response: &SessionResponse_Query{
			Query: &SessionQueryResponse{
				Context: &SessionResponseContext{
					Index:    uint64(m.context.Index()),
					Sequence: context.LastSequenceNumber,
					Status:   SessionResponseStatus_OK,
				},
				Response: &ServiceQueryResponse{
					Response: &ServiceQueryResponse_Metadata{
						Metadata: &ServiceMetadataResponse{
							Services: services,
						},
					},
				},
			},
		},
	}))
}
