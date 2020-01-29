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

package trafficsim

import (
	"github.com/OpenNetworkingFoundation/gmap-ran/api/trafficsim"
	"github.com/OpenNetworkingFoundation/gmap-ran/pkg/manager"
	"github.com/OpenNetworkingFoundation/gmap-ran/pkg/service"
	"github.com/onosproject/onos-config/pkg/utils/logging"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("northbound", "trafficsim")

// NewService returns a new trafficsim Service
func NewService() (service.Service, error) {
	return &Service{}, nil
}

// Service is an implementation of TrafficSim service.
type Service struct {
	service.Service
}

// Register registers the TrafficSim Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{}
	trafficsim.RegisterTrafficServer(r, server);
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
}

func (s *Server) ListRoutes(req *trafficsim.ListRoutesRequest, stream trafficsim.Traffic_ListRoutesServer) error {
	for _, route := range manager.GetManager().Routes {
		resp := &trafficsim.ListRoutesResponse{
			Route: route,
			Type:  trafficsim.Type_NONE,
		}

		err := stream.Send(resp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) ListTowers(req *trafficsim.ListTowersRequest, stream trafficsim.Traffic_ListTowersServer) error {
	for _, tower := range manager.GetManager().Towers {
		err := stream.Send(tower)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) ListUes(req *trafficsim.ListUesRequest, stream trafficsim.Traffic_ListUesServer) error {
	return nil
}
