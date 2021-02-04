// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package trafficsim

import (
	"context"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("trafficsim")

// NewService returns a new trafficsim Service
func NewService(model *model.Model) service.Service {
	return &Service{
		model: model,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	model *model.Model
}

// Register registers the TrafficSim Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		model: s.model,
	}
	trafficsim.RegisterTrafficServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	model *model.Model
}

// GetMapLayout :
func (s *Server) GetMapLayout(ctx context.Context, req *trafficsim.MapLayoutRequest) (*types.MapLayout, error) {
	return nil, nil
}

// ListRoutes :
func (s *Server) ListRoutes(req *trafficsim.ListRoutesRequest, stream trafficsim.Traffic_ListRoutesServer) error {
	return nil
}

// ListCells :
func (s *Server) ListCells(req *trafficsim.ListCellsRequest, stream trafficsim.Traffic_ListCellsServer) error {
	return nil
}

// ListUes :
func (s *Server) ListUes(req *trafficsim.ListUesRequest, stream trafficsim.Traffic_ListUesServer) error {
	return nil
}

// SetNumberUEs changes the number of UEs in the simulation
func (s *Server) SetNumberUEs(ctx context.Context, req *trafficsim.SetNumberUEsRequest) (*trafficsim.SetNumberUEsResponse, error) {
	ueCount := req.GetNumber()
	log.Infof("Number of simulated UEs changed to %d", ueCount)
	s.model.UEs.SetUECount(uint(ueCount))
	return &trafficsim.SetNumberUEsResponse{Number: ueCount}, nil
}

// ResetMetrics resets the metrics on demand
func (s *Server) ResetMetrics(ctx context.Context, req *trafficsim.ResetMetricsMsg) (*trafficsim.ResetMetricsMsg, error) {
	return nil, nil
}
