// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package nodes

import (
	"context"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("api", "nodes")

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

// CreateNode creates a new simulated E2 node
func (s *Server) CreateNode(ctx context.Context, request *modelapi.CreateNodeRequest) (*modelapi.CreateNodeResponse, error) {
	err := s.nodeStore.AddNode(nodeToModel(request.Node))
	if err != nil {
		return nil, err
	}
	return &modelapi.CreateNodeResponse{}, nil
}

// GetNode retrieves the specified simulated E2 node
func (s *Server) GetNode(ctx context.Context, request *modelapi.GetNodeRequest) (*modelapi.GetNodeResponse, error) {
	node, err := s.nodeStore.GetNode(request.EnbID)
	if err != nil {
		return nil, err
	}
	return &modelapi.GetNodeResponse{Node: nodeToAPI(node)}, nil
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

func eventType(t uint8) modelapi.EventType {
	if t == nodes.ADDED {
		return modelapi.EventType_CREATED
	} else if t == nodes.UPDATED {
		return modelapi.EventType_UPDATED
	} else if t == nodes.DELETED {
		return modelapi.EventType_DELETED
	} else {
		return modelapi.EventType_NONE
	}
}

// WatchNodes monitors changes to the inventory of E2 nodes
func (s *Server) WatchNodes(request *modelapi.WatchNodesRequest, server modelapi.NodeModel_WatchNodesServer) error {
	log.Infof("Watching nodes [%v]...", request)
	ch := make(chan nodes.NodeEvent)
	s.nodeStore.WatchNodes(ch, nodes.WatchOptions{Replay: !request.NoReplay, Monitor: !request.NoSubscribe})

	for event := range ch {
		response := &modelapi.WatchNodesResponse{
			Node: nodeToAPI(event.Node),
			Type: eventType(event.Type),
		}
		err := server.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}
