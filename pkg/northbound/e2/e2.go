// Copyright 2020-present Open Networking Foundation.
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

package e2

import (
	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/onosproject/ran-simulator/pkg/service"
	"google.golang.org/grpc"
)

// NewService returns a new trafficsim Service
func NewService() (service.Service, error) {
	return &Service{}, nil
}

// Service is an implementation of e2 service.
type Service struct {
	service.Service
}

// Register registers the e2 Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{}
	e2.RegisterInterfaceServiceServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
}

// SendTelemetry ...
func (s *Server) SendTelemetry(req *e2.L2MeasConfig, stream e2.InterfaceService_SendTelemetryServer) error {
	mgr := GetManager()
	return mgr.RunTelemetry(stream)
}

// SendControl ...
func (s *Server) SendControl(stream e2.InterfaceService_SendControlServer) error {
	mgr := GetManager()
	return mgr.RunControl(stream)
}
