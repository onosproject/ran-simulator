// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package nodes

import (
	"context"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger()

// NewService returns a new model Service
func NewService(nodeStore nodes.Store, plmnID types.PlmnID) service.Service {
	return &Service{
		plmnID:    plmnID,
		nodeStore: nodeStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	plmnID    types.PlmnID
	nodeStore nodes.Store
}

// Register registers the TrafficSim Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		plmnID:    s.plmnID,
		nodeStore: s.nodeStore,
	}
	modelapi.RegisterNodeModelServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	plmnID    types.PlmnID
	nodeStore nodes.Store
}

func nodeToAPI(node *model.Node) *types.Node {
	return &types.Node{
		GnbID:         node.GnbID,
		Controllers:   node.Controllers,
		ServiceModels: node.ServiceModels,
		CellNCGIs:     node.Cells,
		Status:        node.Status,
	}
}

func nodeToModel(node *types.Node) *model.Node {
	return &model.Node{
		GnbID:         node.GnbID,
		Controllers:   node.Controllers,
		ServiceModels: node.ServiceModels,
		Cells:         node.CellNCGIs,
		Status:        node.Status,
	}
}

// GetPlmnID returns the PLMNID used by the RAN simulator for all simulated entities
func (s *Server) GetPlmnID(ctx context.Context, request *modelapi.PlmnIDRequest) (*modelapi.PlmnIDResponse, error) {
	return &modelapi.PlmnIDResponse{PlmnID: s.plmnID}, nil
}

// CreateNode creates a new simulated E2 node
func (s *Server) CreateNode(ctx context.Context, request *modelapi.CreateNodeRequest) (*modelapi.CreateNodeResponse, error) {
	log.Debugf("Received create node request: %+v", request)
	err := s.nodeStore.Add(ctx, nodeToModel(request.Node))
	if err != nil {
		return nil, err
	}
	return &modelapi.CreateNodeResponse{}, nil
}

// GetNode retrieves the specified simulated E2 node
func (s *Server) GetNode(ctx context.Context, request *modelapi.GetNodeRequest) (*modelapi.GetNodeResponse, error) {
	log.Debugf("Received get node request: %+v", request)
	node, err := s.nodeStore.Get(ctx, request.GnbID)
	if err != nil {
		return nil, err
	}
	return &modelapi.GetNodeResponse{Node: nodeToAPI(node)}, nil
}

// UpdateNode updates the specified simulated E2 node
func (s *Server) UpdateNode(ctx context.Context, request *modelapi.UpdateNodeRequest) (*modelapi.UpdateNodeResponse, error) {
	log.Debugf("Received update node request: %+v", request)
	err := s.nodeStore.Update(ctx, nodeToModel(request.Node))
	if err != nil {
		return nil, err
	}
	return &modelapi.UpdateNodeResponse{}, nil
}

// DeleteNode deletes the specified simulated E2 node
func (s *Server) DeleteNode(ctx context.Context, request *modelapi.DeleteNodeRequest) (*modelapi.DeleteNodeResponse, error) {
	log.Debugf("Received delete node request: %v", request)
	_, err := s.nodeStore.Delete(ctx, request.GnbID)
	if err != nil {
		return nil, err
	}
	return &modelapi.DeleteNodeResponse{}, nil
}

// ListNodes list of e2 nodes
func (s *Server) ListNodes(request *modelapi.ListNodesRequest, server modelapi.NodeModel_ListNodesServer) error {
	log.Debugf("Received listing nodes request: %v", request)
	nodeList, _ := s.nodeStore.List(server.Context())
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
	log.Debugf("Received watching node changes Request: %v", request)
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
	node, err := s.nodeStore.Get(ctx, request.GnbID)
	if err != nil {
		return nil, err
	}
	log.Infof("Requested '%s' of agent %d", request.Command, node.GnbID)
	// TODO: implement agent stop|start, implement connection drop|reconnect, etc.
	// For now, just put the command into the status
	err = s.nodeStore.SetStatus(ctx, node.GnbID, request.Command)
	if err != nil {
		return nil, err
	}
	return &modelapi.AgentControlResponse{Node: nodeToAPI(node)}, nil
}
