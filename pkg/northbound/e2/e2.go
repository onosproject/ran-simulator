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
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"google.golang.org/grpc"
)

// ServerParams - params to start a new server
type ServerParams struct {
	CaPath       string
	KeyPath      string
	CertPath     string
	TopoEndpoint string
}

var log = liblog.GetLogger("northbound", "e2")

// NewTowerServer - start a new gRPC server per tower
func NewTowerServer(towerIndex int, serverParams ServerParams) error {
	port := int16(utils.GrpcBasePort + towerIndex + 1) // Start at 5152 so this translates to 1420 in Hex
	ecID := utils.EcIDForPort(port)
	s := service.NewServer(service.NewServerConfig(serverParams.CaPath, serverParams.KeyPath, serverParams.CertPath, port, true))
	s.AddService(Service{
		port:      port,
		towerEcID: ecID,
	})

	return s.Serve(func(started string) {
		log.Info("Started E2 server on ", started)
	})
}

// Service is an implementation of e2 service.
type Service struct {
	service.Service
	port      int16
	towerEcID types.EcID
}

// Register registers the e2 Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{port: s.port, towerEcID: s.towerEcID}
	e2.RegisterInterfaceServiceServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	port      int16
	towerEcID types.EcID
}

// GetPort - expose the port number
func (s *Server) GetPort() int16 {
	return s.port
}

// GetEcID - expose the tower EcID
func (s *Server) GetEcID() types.EcID {
	return s.towerEcID
}
