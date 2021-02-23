// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package nodes

import (
	"context"

	"github.com/onosproject/ran-simulator/pkg/store/event"

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
func NewService(nodeStore nodes.Store) service.Service {
	return &Service{
		nodeStore: nodeStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	nodeStore nodes.Store
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
	nodeStore nodes.Store
}

func nodeToAPI(node *model.Node) *types.Node {
	return &types.Node{
		EnbID:         node.EnbID,
		Controllers:   node.Controllers,
		ServiceModels: node.ServiceModels,
		CellECGIs:     node.Cells,
		Status:        node.Status,
	}
}

func nodeToModel(node *types.Node) *model.Node {
	return &model.Node{
		EnbID:         node.EnbID,
		Controllers:   node.Controllers,
		ServiceModels: node.ServiceModels,
		Cells:         node.CellECGIs,
		Status:        node.Status,
	}
}

// CreateNode creates a new simulated E2 node
func (s *Server) CreateNode(ctx context.Context, request *modelapi.CreateNodeRequest) (*modelapi.CreateNodeResponse, error) {
	err := s.nodeStore.Add(ctx, nodeToModel(request.Node))
	if err != nil {
		return nil, err
	}
	return &modelapi.CreateNodeResponse{}, nil
}

// GetNode retrieves the specified simulated E2 node
func (s *Server) GetNode(ctx context.Context, request *modelapi.GetNodeRequest) (*modelapi.GetNodeResponse, error) {
	node, err := s.nodeStore.Get(ctx, request.EnbID)
	if err != nil {
		return nil, err
	}
	return &modelapi.GetNodeResponse{Node: nodeToAPI(node)}, nil
}

// UpdateNode updates the specified simulated E2 node
func (s *Server) UpdateNode(ctx context.Context, request *modelapi.UpdateNodeRequest) (*modelapi.UpdateNodeResponse, error) {
	err := s.nodeStore.Update(ctx, nodeToModel(request.Node))
	if err != nil {
		return nil, err
	}
	return &modelapi.UpdateNodeResponse{}, nil
}

// DeleteNode deletes the specified simulated E2 node
func (s *Server) DeleteNode(ctx context.Context, request *modelapi.DeleteNodeRequest) (*modelapi.DeleteNodeResponse, error) {
	_, err := s.nodeStore.Delete(ctx, request.EnbID)
	if err != nil {
		return nil, err
	}
	return &modelapi.DeleteNodeResponse{}, nil
}

// ListNodes list of e2 nodes
func (s *Server) ListNodes(request *modelapi.ListNodesRequest, server modelapi.NodeModel_ListNodesServer) error {
	nodeList, _ := s.nodeStore.List(server.Context())
	log.Info("List of nodes:", nodeList)
	for _, node := range nodeList {
		log.Info("Node:", node)
		resp := &modelapi.ListNodesResponse{
			Node: nodeToAPI(node),
		}
		err := server.Send(resp)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func eventType(nodeEvent nodes.NodeEvent) modelapi.EventType {
	if nodeEvent == nodes.Created {
		return modelapi.EventType_CREATED
	} else if nodeEvent == nodes.Updated {
		return modelapi.EventType_UPDATED
	} else if nodeEvent == nodes.Deleted {
		return modelapi.EventType_DELETED
	} else {
		return modelapi.EventType_NONE
	}
}

// WatchNodes monitors changes to the inventory of E2 nodes
func (s *Server) WatchNodes(request *modelapi.WatchNodesRequest, server modelapi.NodeModel_WatchNodesServer) error {
	log.Infof("Watching nodes [%v]...", request)
	ch := make(chan event.Event)
	err := s.nodeStore.Watch(server.Context(), ch, nodes.WatchOptions{Replay: !request.NoReplay, Monitor: !request.NoSubscribe})

	if err != nil {
		return err
	}

	for nodeEvent := range ch {
		response := &modelapi.WatchNodesResponse{
			Node: nodeToAPI(nodeEvent.Value.(*model.Node)),
			Type: eventType(nodeEvent.Type.(nodes.NodeEvent)),
		}
		err := server.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}

// AgentControl allows control over the lifecycle of the agent running on behalf of the simulated E2 node
func (s *Server) AgentControl(ctx context.Context, request *modelapi.AgentControlRequest) (*modelapi.AgentControlResponse, error) {
	node, err := s.nodeStore.Get(ctx, request.EnbID)
	if err != nil {
		return nil, err
	}
	log.Infof("Requested '%s' of agent %d", request.Command, node.EnbID)
	// TODO: implement agent stop|start, implement connection drop|reconnect, etc.
	// For now, just put the command into the status
	err = s.nodeStore.SetStatus(ctx, node.EnbID, request.Command)
	if err != nil {
		return nil, err
	}
	return &modelapi.AgentControlResponse{Node: nodeToAPI(node)}, nil
}
