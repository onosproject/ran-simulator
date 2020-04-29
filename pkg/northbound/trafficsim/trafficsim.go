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
	"context"
	"fmt"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log = liblog.GetLogger("northbound", "trafficsim")

// NewService returns a new trafficsim Service
func NewService() (service.Service, error) {
	return &Service{}, nil
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
}

// Register registers the TrafficSim Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{}
	trafficsim.RegisterTrafficServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
}

// GetMapLayout :
func (s *Server) GetMapLayout(ctx context.Context, req *trafficsim.MapLayoutRequest) (*types.MapLayout, error) {
	return &manager.GetManager().MapLayout, nil
}

// ListRoutes :
func (s *Server) ListRoutes(req *trafficsim.ListRoutesRequest, stream trafficsim.Traffic_ListRoutesServer) error {
	if !req.WithoutReplay {
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
	}

	if req.Subscribe {
		streamID := fmt.Sprintf("route-%p", stream)
		listener, err := manager.GetManager().Dispatcher.RegisterRouteListener(streamID)
		defer manager.GetManager().Dispatcher.UnregisterRouteListener(streamID)
		if err != nil {
			log.Info("Failed setting up a listener for Route events")
			return err
		}
		log.Infof("NBI Route updates started on %s", streamID)

		for {
			select {
			case routeEvent := <-listener:
				route, objOk := routeEvent.Object.(*types.Route)
				if !objOk {
					return fmt.Errorf("could not cast object from event to Route %v", routeEvent)
				}
				msg := &trafficsim.ListRoutesResponse{
					Route: route,
					Type:  routeEvent.Type,
				}
				err := stream.SendMsg(msg)
				if err != nil {
					log.Warnf("Error sending message on stream. Closing. %v", msg)
					return err
				}
			case <-stream.Context().Done():
				log.Infof("Client has disconnected ListRoutes on %s", streamID)
				return nil
			}
		}
	}

	return nil
}

// ListCells :
func (s *Server) ListCells(req *trafficsim.ListCellsRequest, stream trafficsim.Traffic_ListCellsServer) error {
	if !req.WithoutReplay {
		manager.GetManager().CellsLock.RLock()
		for _, cell := range manager.GetManager().Cells {
			resp := &trafficsim.ListCellsResponse{
				Cell: cell,
				Type: trafficsim.Type_NONE,
			}
			err := stream.Send(resp)
			if err != nil {
				manager.GetManager().CellsLock.RUnlock()
				return err
			}
		}
		manager.GetManager().CellsLock.RUnlock()
	}

	if req.Subscribe {
		streamID := fmt.Sprintf("cell-%p", stream)
		listener, err := manager.GetManager().Dispatcher.RegisterCellListener(streamID)
		if err != nil {
			log.Info("Failed setting up a listener for Ue events")
			return err
		}
		defer manager.GetManager().Dispatcher.UnregisterCellListener(streamID)
		log.Infof("NBI Cell updates started on %s", streamID)

		for {
			select {
			case cellEvent := <-listener:
				cell, objOk := cellEvent.Object.(*types.Cell)
				if !objOk {
					return fmt.Errorf("could not cast object from event to Cell %v", cellEvent)
				}
				msg := &trafficsim.ListCellsResponse{
					Cell: cell,
					Type: cellEvent.Type,
				}
				err := stream.SendMsg(msg)
				if err != nil {
					log.Warnf("Error sending message on stream. Closing. %v", msg)
					return err
				}
			case <-stream.Context().Done():
				log.Infof("Client has disconnected ListCells on %s", streamID)
				return nil
			}
		}
	}

	return nil
}

// ListUes :
func (s *Server) ListUes(req *trafficsim.ListUesRequest, stream trafficsim.Traffic_ListUesServer) error {
	if !req.WithoutReplay {
		manager.GetManager().UserEquipmentsLock.RLock()
		for _, ue := range manager.GetManager().UserEquipments {
			resp := &trafficsim.ListUesResponse{
				Ue:   ue,
				Type: trafficsim.Type_NONE,
			}
			err := stream.Send(resp)
			if err != nil {
				manager.GetManager().UserEquipmentsLock.RUnlock()
				return err
			}
		}
		manager.GetManager().UserEquipmentsLock.RUnlock()
	}

	if req.Subscribe {
		streamID := fmt.Sprintf("ue-%p", stream)
		listener, err := manager.GetManager().Dispatcher.RegisterUeListener(streamID)
		if err != nil {
			log.Info("Failed setting up a listener for Ue events")
			return err
		}
		defer manager.GetManager().Dispatcher.UnregisterUeListener(streamID)
		log.Infof("NBI Ue updates started on %s", streamID)

		for {
			select {
			case ueEvent := <-listener:
				ue, objOk := ueEvent.Object.(*types.Ue)
				if !objOk {
					return fmt.Errorf("could not cast object from event to UE %v", ueEvent)
				}
				msg := &trafficsim.ListUesResponse{
					Ue:         manager.UeDeepCopy(ue),
					Type:       ueEvent.Type,
					UpdateType: ueEvent.UpdateType,
				}
				err := stream.SendMsg(msg)
				if err != nil {
					log.Warnf("Error sending message on stream. Closing. %v", msg)
					return err
				}
			case <-stream.Context().Done():
				log.Infof("Client has disconnected ListUes on %s", streamID)
				return nil
			}
		}
	}

	return nil
}

// SetNumberUEs - change the number of UEs in the simulation
// Cannot be set below the minimum or above the maximum
func (s *Server) SetNumberUEs(ctx context.Context, req *trafficsim.SetNumberUEsRequest) (*trafficsim.SetNumberUEsResponse, error) {
	numUes := req.GetNumber()
	minUes := manager.GetManager().MapLayout.MinUes
	maxUes := manager.GetManager().MapLayout.MaxUes
	if numUes < minUes {
		return nil, status.Errorf(codes.OutOfRange,
			"number of UEs requested %d is below minimum %d", numUes, minUes)
	} else if numUes > maxUes {
		return nil, status.Errorf(codes.OutOfRange,
			"number of UEs requested %d is above maximum %d", numUes, maxUes)
	}

	err := manager.GetManager().SetNumberUes(int(numUes))
	if err != nil {
		return nil, status.Error(codes.OutOfRange, err.Error())
	}
	return &trafficsim.SetNumberUEsResponse{Number: numUes}, nil
}

// ResetMetrics resets the metrics on demand
func (s *Server) ResetMetrics(ctx context.Context, req *trafficsim.ResetMetricsMsg) (*trafficsim.ResetMetricsMsg, error) {
	manager.GetManager().ResetMetricsChannel <- true
	log.Warn("Metrics reset")

	return &trafficsim.ResetMetricsMsg{}, nil
}
