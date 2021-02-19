// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package nodes

import (
	"context"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"google.golang.org/grpc"
)

// NewService returns a new model Service
func NewService(nodeStore nodes.NodeRegistry) service.Service {
	return &Service{
		nodeStore: nodeStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	nodeStore nodes.NodeRegistry
}

// Register registers the TrafficSim Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		s.nodeStore,
	}
	modelapi.RegisterNodeModelServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	nodeStore nodes.NodeRegistry
}

// CreateNode creates a new simulated E2 node
func (s *Server) CreateNode(ctx context.Context, request *modelapi.CreateNodeRequest) (*modelapi.CreateNodeResponse, error) {
	panic("implement me")
}

// GetNode retrieves the specified simulated E2 node
func (s *Server) GetNode(ctx context.Context, request *modelapi.GetNodeRequest) (*modelapi.GetNodeResponse, error) {
	node, err := s.nodeStore.GetNode(request.EnbID)
	if err != nil {
		return nil, err
	}
	return &modelapi.GetNodeResponse{Node: nodeToAPI(node)}, nil
}

func nodeToAPI(node *model.Node) *types.Node {
	return &types.Node{
		EnbID:         node.EnbID,
		Controllers:   node.Controllers,
		ServiceModels: node.ServiceModels,
		CellECGIs:     node.Cells,
	}
}

func nodeToModel(node *types.Node) *model.Node {
	return &model.Node{
		EnbID:         node.EnbID,
		Controllers:   node.Controllers,
		ServiceModels: node.ServiceModels,
		Cells:         node.CellECGIs,
	}
}

// UpdateNode updates the specified simulated E2 node
func (s *Server) UpdateNode(ctx context.Context, request *modelapi.UpdateNodeRequest) (*modelapi.UpdateNodeResponse, error) {
	err := s.nodeStore.UpdateNode(nodeToModel(request.Node))
	if err != nil {
		return nil, err
	}
	return &modelapi.UpdateNodeResponse{}, nil
}

// DeleteNode deletes the specified simulated E2 node
func (s *Server) DeleteNode(ctx context.Context, request *modelapi.DeleteNodeRequest) (*modelapi.DeleteNodeResponse, error) {
	_, err := s.nodeStore.DeleteNode(request.EnbID)
	if err != nil {
		return nil, err
	}
	return &modelapi.DeleteNodeResponse{}, nil
}

// WatchNodes monitors changes to the inventory of E2 nodes
func (s *Server) WatchNodes(request *modelapi.WatchNodesRequest, server modelapi.NodeModel_WatchNodesServer) error {
	panic("implement me")
}
