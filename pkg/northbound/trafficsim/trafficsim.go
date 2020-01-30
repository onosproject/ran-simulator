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
	"fmt"
	"context"
	"github.com/OpenNetworkingFoundation/gmap-ran/api/trafficsim"
	"github.com/OpenNetworkingFoundation/gmap-ran/api/types"
	"github.com/OpenNetworkingFoundation/gmap-ran/pkg/manager"
	"github.com/OpenNetworkingFoundation/gmap-ran/pkg/service"
	"google.golang.org/grpc"
	log "k8s.io/klog"
)

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

func (s *Server) GetMapLayout(ctx context.Context, req *trafficsim.MapLayoutRequest) (*types.MapLayout, error) {
	return &manager.GetManager().MapLayout, nil
}

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
		defer manager.GetManager().Dispatcher.UnregisterRouteListener(streamID);
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
					log.Warningf("Error sending message on stream. Closing. %v", msg)
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
	if !req.WithoutReplay {
		for _, ue := range manager.GetManager().UserEquipments {
			resp := &trafficsim.ListUesResponse{
				Ue:   ue,
				Type: trafficsim.Type_NONE,
			}
			err := stream.Send(resp)
			if err != nil {
				return err
			}
		}
	}

	if req.Subscribe {
		streamID := fmt.Sprintf("ue-%p", stream)
		listener, err := manager.GetManager().Dispatcher.RegisterUeListener(streamID)
		if err != nil {
			log.Info("Failed setting up a listener for Ue events")
			return err
		}
		defer manager.GetManager().Dispatcher.UnregisterUeListener(streamID);
		log.Infof("NBI Ue updates started on %s", streamID)

		for {
			select {
			case ueEvent := <-listener:
				ue, objOk := ueEvent.Object.(*types.Ue)
				if !objOk {
					return fmt.Errorf("could not cast object from event to UE %v", ueEvent)
				}
				msg := &trafficsim.ListUesResponse{
					Ue:   ue,
					Type: ueEvent.Type,
				}
				err := stream.SendMsg(msg)
				if err != nil {
					log.Warningf("Error sending message on stream. Closing. %v", msg)
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
