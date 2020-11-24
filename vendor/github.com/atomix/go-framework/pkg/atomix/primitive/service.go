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
	"io"
)

// ServiceID is a service identifier
type ServiceID ServiceId

// ServiceContext provides information about the context within which a service is running
type ServiceContext interface {
	PartitionContext

	// ServiceID is the service identifier
	ServiceID() ServiceID

	// ServiceType returns the service type
	ServiceType() ServiceType

	// CurrentOperation returns the current operation identifier
	CurrentOperation() OperationID

	// CurrentSession returns the current session
	CurrentSession() Session

	// Session returns the session with the given identifier
	Session(id SessionID) Session

	// Sessions returns a list of open sessions
	Sessions() []Session
}

// internalContext provides setters for the service context
type internalContext interface {
	ServiceContext
	setCurrentOperation(op OperationID)
	setCurrentSession(session Session)
	addSession(session Session)
	removeSession(session Session)
}

func newServiceContext(ctx PartitionContext, id ServiceID) ServiceContext {
	return &serviceContext{
		PartitionContext: ctx,
		serviceID:        id,
		sessions:         make(map[SessionID]Session),
	}
}

// serviceContext is a default implementation of the service context
type serviceContext struct {
	PartitionContext
	serviceID      ServiceID
	sessions       map[SessionID]Session
	currentSession Session
	currentOp      OperationID
}

func (c *serviceContext) ServiceID() ServiceID {
	return c.serviceID
}

func (c *serviceContext) ServiceType() ServiceType {
	return c.serviceID.Type
}

func (c *serviceContext) CurrentOperation() OperationID {
	return c.currentOp
}

// setOperation sets the current operation
func (c *serviceContext) setCurrentOperation(op OperationID) {
	c.currentOp = op
}

func (c *serviceContext) CurrentSession() Session {
	return c.currentSession
}

// setCurrentSession sets the current session
func (c *serviceContext) setCurrentSession(session Session) {
	c.currentSession = session
}

func (c *serviceContext) Session(id SessionID) Session {
	return c.sessions[id]
}

func (c *serviceContext) Sessions() []Session {
	sessions := make([]Session, 0, len(c.sessions))
	for _, session := range c.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// addSession adds a session to the service
func (c *serviceContext) addSession(session Session) {
	c.sessions[session.ID()] = session
}

// removeSession removes a session from the service
func (c *serviceContext) removeSession(session Session) {
	delete(c.sessions, session.ID())
}

var _ ServiceContext = &serviceContext{}

// SessionOpenService is an interface for listening to session open events
type SessionOpenService interface {
	// SessionOpen is called when a session is opened for a service
	SessionOpen(Session)
}

// SessionClosedService is an interface for listening to session closed events
type SessionClosedService interface {
	// SessionClosed is called when a session is closed for a service
	SessionClosed(Session)
}

// SessionExpiredService is an interface for listening to session expired events
type SessionExpiredService interface {
	// SessionExpired is called when a session is expired for a service
	SessionExpired(Session)
}

// BackupService is an interface for backing up a service
type BackupService interface {
	// Backup is called to take a snapshot of the service state
	Backup(writer io.Writer) error
}

// RestoreService is an interface for restoring up a service
type RestoreService interface {
	// Restore is called to restore the service state from a snapshot
	Restore(reader io.Reader) error
}

// Service is a primitive service
type Service interface {
	BackupService
	RestoreService
	Executor
	Scheduler
	internalContext
}

// NewService creates a new primitive service
func NewService(scheduler Scheduler, context ServiceContext) Service {
	return &managedService{
		Executor:        newExecutor(),
		Scheduler:       scheduler,
		internalContext: context.(internalContext),
	}
}

// managedService is a primitive service
type managedService struct {
	BackupService
	RestoreService
	Executor
	Scheduler
	internalContext
}
