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
	"time"

	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("northbound", "e2")

// NewTowerServer - start a new gRPC server per tower
func NewTowerServer(ecgi types.ECGI, port uint16, serverParams utils.ServerParams) error {
	s := service.NewServer(service.NewServerConfig(serverParams.CaPath, serverParams.KeyPath, serverParams.CertPath, int16(port), true))
	s.AddService(Service{
		port:      int(port),
		towerEcID: ecgi.EcID,
		plmnID:    ecgi.PlmnID,
	})

	return s.Serve(func(started string) {
		log.Info("Started E2 server on ", started)
	})
}

// Service is an implementation of e2 service.
type Service struct {
	service.Service
	port      int
	towerEcID types.EcID
	plmnID    types.PlmnID
}

// Register registers the e2 Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{port: s.port, towerEcID: s.towerEcID, plmnID: s.plmnID}
	e2ap.RegisterE2APServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	port            int
	towerEcID       types.EcID
	plmnID          types.PlmnID
	stream          e2ap.E2AP_RicChanServer
	indChan         chan e2ap.RicIndication
	l2MeasConfig    e2.L2MeasConfig
	telemetryTicker *time.Ticker
}

// GetPort - expose the port number
func (s *Server) GetPort() int {
	return s.port
}

// GetEcID - expose the tower EcID
func (s *Server) GetEcID() types.EcID {
	return s.towerEcID
}

// GetPlmnID - expose the tower PlmnID
func (s *Server) GetPlmnID() types.PlmnID {
	return s.plmnID
}

// GetECGI - expose the tower Ecgi
func (s *Server) GetECGI() types.ECGI {
	return newEcgi(s.GetEcID(), s.GetPlmnID())
}

func toTypesEcgi(e2Ecgi *e2.ECGI) types.ECGI {
	return types.ECGI{
		EcID:   types.EcID(e2Ecgi.Ecid),
		PlmnID: types.PlmnID(e2Ecgi.PlmnId),
	}
}

func toE2Ecgi(e2Ecgi *types.ECGI) e2.ECGI {
	return e2.ECGI{
		Ecid:   string(e2Ecgi.EcID),
		PlmnId: string(e2Ecgi.PlmnID),
	}
}

func newEcgi(id types.EcID, plmnID types.PlmnID) types.ECGI {
	return types.ECGI{EcID: id, PlmnID: plmnID}
}
