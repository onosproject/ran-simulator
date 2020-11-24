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

package log

import (
	api "github.com/atomix/api/proto/atomix/log"
	primitiveapi "github.com/atomix/api/proto/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix"
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"google.golang.org/grpc"
)

// Type is the log primitive type
const Type = primitiveapi.PrimitiveType_LOG

// RegisterPrimitive registers the primitive on the given node
func RegisterPrimitive(node *atomix.Node) {
	node.RegisterPrimitive(Type, &Primitive{})
}

// Primitive is the counter primitive
type Primitive struct{}

// RegisterServer registers the log server with the protocol
func (p *Primitive) RegisterServer(server *grpc.Server, protocol primitive.Protocol) {
	api.RegisterLogServiceServer(server, &Server{
		Server: &primitive.Server{
			Type:     primitive.ServiceType_LOG,
			Protocol: protocol,
		},
	})
}

// NewService creates a new log service
func (p *Primitive) NewService(scheduler primitive.Scheduler, context primitive.ServiceContext) primitive.Service {
	service := &Service{
		Service:   primitive.NewService(scheduler, context),
		indexes:   make(map[uint64]*LinkedLogEntryValue),
		listeners: make(map[primitive.SessionID]map[primitive.StreamID]listener),
	}
	service.init()
	return service
}

var _ primitive.Primitive = &Primitive{}
