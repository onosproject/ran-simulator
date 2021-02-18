// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cells

import (
	"context"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/onosproject/ran-simulator/pkg/model"
	"google.golang.org/grpc"
)

// NewService returns a new model Service
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
	modelapi.RegisterCellModelServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	model *model.Model
}

// CreateCell creates a new simulated cell
func (s *Server) CreateCell(ctx context.Context, request *modelapi.CreateCellRequest) (*modelapi.CreateCellResponse, error) {
	return &modelapi.CreateCellResponse{}, nil
}

// GetCell retrieves the specified simulated cell
func (s *Server) GetCell(ctx context.Context, request *modelapi.GetCellRequest) (*modelapi.GetCellResponse, error) {
	panic("implement me")
}

// UpdateCell updates the specified simulated cell
func (s *Server) UpdateCell(ctx context.Context, request *modelapi.UpdateCellRequest) (*modelapi.UpdateCellResponse, error) {
	panic("implement me")
}

// DeleteCell deletes the specified simulated cell
func (s *Server) DeleteCell(ctx context.Context, request *modelapi.DeleteCellRequest) (*modelapi.DeleteCellResponse, error) {
	panic("implement me")
}

// WatchCells monitors changes to the inventory of cells
func (s *Server) WatchCells(request *modelapi.WatchCellsRequest, server modelapi.CellModel_WatchCellsServer) error {
	panic("implement me")
}
