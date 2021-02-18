// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package nodes

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
	modelapi.RegisterNodeModelServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	model *model.Model
}

// CreateNode creates a new simulated E2 node
func (s *Server) CreateNode(ctx context.Context, request *modelapi.CreateNodeRequest) (*modelapi.CreateNodeResponse, error) {
	panic("implement me")
}

// GetNode retrieves the specified simulated E2 node
func (s *Server) GetNode(ctx context.Context, request *modelapi.GetNodeRequest) (*modelapi.GetNodeResponse, error) {
	panic("implement me")
}

// UpdateNode updates the specified simulated E2 node
func (s *Server) UpdateNode(ctx context.Context, request *modelapi.UpdateNodeRequest) (*modelapi.UpdateNodeResponse, error) {
	panic("implement me")
}

// DeleteNode deletes the specified simulated E2 node
func (s *Server) DeleteNode(ctx context.Context, request *modelapi.DeleteNodeRequest) (*modelapi.DeleteNodeResponse, error) {
	panic("implement me")
}

// WatchNodes monitors changes to the inventory of E2 nodes
func (s *Server) WatchNodes(request *modelapi.WatchNodesRequest, server modelapi.NodeModel_WatchNodesServer) error {
	panic("implement me")
}
