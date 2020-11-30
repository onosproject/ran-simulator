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
	"google.golang.org/grpc"
)

// NewInstance creates a new primitive instance
func NewInstance(ctx context.Context, name Name, session *Session, handler Handler) (*Instance, error) {
	instance := &Instance{
		Name:    name,
		Session: session,
		handler: handler,
	}
	if err := instance.create(ctx); err != nil {
		return nil, err
	}
	return instance, nil
}

// Instance is a primitive instance
type Instance struct {
	Name    Name
	Session *Session
	handler Handler
}

// DoCreate sends a create session request
func (i *Instance) DoCreate(ctx context.Context, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) error {
	return i.Session.doCreate(ctx, i.Name, f)
}

// DoClose sends a session close request
func (i *Instance) DoClose(ctx context.Context, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) error {
	return i.Session.doClose(ctx, i.Name, f)
}

// DoQuery sends a session query request
func (i *Instance) DoQuery(ctx context.Context, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) (interface{}, error) {
	return i.Session.doQuery(ctx, i.Name, f)
}

// DoCommand sends a session command request
func (i *Instance) DoCommand(ctx context.Context, f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error)) (interface{}, error) {
	return i.Session.doCommand(ctx, i.Name, f)
}

// DoQueryStream sends a session query stream request
func (i *Instance) DoQueryStream(
	ctx context.Context,
	f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error),
	responseFunc func(interface{}) (*headers.ResponseHeader, interface{}, error)) (<-chan interface{}, error) {
	return i.Session.doQueryStream(ctx, i.Name, f, responseFunc)
}

// DoCommandStream sends a session command stream request
func (i *Instance) DoCommandStream(
	ctx context.Context,
	f func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error),
	responseFunc func(interface{}) (*headers.ResponseHeader, interface{}, error)) (<-chan interface{}, error) {
	return i.Session.doCommandStream(ctx, i.Name, f, responseFunc)
}

// create creates the instance
func (i *Instance) create(ctx context.Context) error {
	return i.handler.Create(ctx, i)
}

// Close closes the instance
func (i *Instance) Close(ctx context.Context) error {
	return i.handler.Close(ctx, i)
}

// Delete deletes the instance
func (i *Instance) Delete(ctx context.Context) error {
	return i.handler.Delete(ctx, i)
}
