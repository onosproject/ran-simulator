// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package routes

import (
	"context"

	"github.com/onosproject/ran-simulator/pkg/store/routes"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/pkg/model"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger()

// NewService returns a new model Service
func NewService(routeStore routes.Store) service.Service {
	return &Service{
		routeStore: routeStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	routeStore routes.Store
}

// Register registers the TrafficSim Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		routeStore: s.routeStore,
	}
	modelapi.RegisterRouteModelServer(r, server)
}

// Server implements the Routes gRPC service for administrative facilities.
type Server struct {
	routeStore routes.Store
}

func routeToAPI(route *model.Route) *types.Route {
	points := make([]*types.Point, 0, len(route.Points))
	for _, p := range route.Points {
		points = append(points, &types.Point{Lat: p.Lat, Lng: p.Lng})
	}
	return &types.Route{
		RouteID:    route.IMSI,
		Waypoints:  points,
		Color:      route.Color,
		SpeedAvg:   route.SpeedAvg,
		SpeedStdev: route.SpeedStdDev,
		NextPoint:  route.NextPoint,
		Reverse:    route.Reverse,
	}
}

func routeToModel(route *types.Route) *model.Route {
	points := make([]*model.Coordinate, 0, len(route.Waypoints))
	for _, p := range route.Waypoints {
		points = append(points, &model.Coordinate{Lat: p.Lat, Lng: p.Lng})
	}
	return &model.Route{
		IMSI:        route.RouteID,
		Points:      points,
		Color:       route.Color,
		SpeedAvg:    route.SpeedAvg,
		SpeedStdDev: route.SpeedStdev,
		NextPoint:   route.NextPoint,
		Reverse:     route.Reverse,
	}
}

// CreateRoute creates a new simulated route of a UE
func (s *Server) CreateRoute(ctx context.Context, request *modelapi.CreateRouteRequest) (*modelapi.CreateRouteResponse, error) {
	log.Debugf("Received create route request: %+v", request)
	err := s.routeStore.Add(ctx, routeToModel(request.Route))
	if err != nil {
		return nil, err
	}
	return &modelapi.CreateRouteResponse{}, nil
}

// GetRoute retrieves the specified UE route
func (s *Server) GetRoute(ctx context.Context, request *modelapi.GetRouteRequest) (*modelapi.GetRouteResponse, error) {
	log.Debugf("Received get route request: %+v", request)
	route, err := s.routeStore.Get(ctx, request.IMSI)
	if err != nil {
		return nil, err
	}
	return &modelapi.GetRouteResponse{Route: routeToAPI(route)}, nil
}

// DeleteRoute deletes the specified simulated E2 route
func (s *Server) DeleteRoute(ctx context.Context, request *modelapi.DeleteRouteRequest) (*modelapi.DeleteRouteResponse, error) {
	log.Debugf("Received delete route request: %v", request)
	_, err := s.routeStore.Delete(ctx, request.IMSI)
	if err != nil {
		return nil, err
	}
	return &modelapi.DeleteRouteResponse{}, nil
}

// ListRoutes list of e2 routes
func (s *Server) ListRoutes(request *modelapi.ListRoutesRequest, server modelapi.RouteModel_ListRoutesServer) error {
	log.Debugf("Received listing routes request: %v", request)
	routeList := s.routeStore.List(server.Context())
	for _, route := range routeList {
		log.Info("Route:", route)
		resp := &modelapi.ListRoutesResponse{
			Route: routeToAPI(route),
		}
		err := server.Send(resp)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func eventType(routeEvent routes.RouteEvent) modelapi.EventType {
	if routeEvent == routes.Created {
		return modelapi.EventType_CREATED
	} else if routeEvent == routes.Updated {
		return modelapi.EventType_UPDATED
	} else if routeEvent == routes.Deleted {
		return modelapi.EventType_DELETED
	} else {
		return modelapi.EventType_NONE
	}
}

// WatchRoutes monitors changes to the inventory of E2 routes
func (s *Server) WatchRoutes(request *modelapi.WatchRoutesRequest, server modelapi.RouteModel_WatchRoutesServer) error {
	log.Debugf("Received watching route changes Request: %v", request)
	ch := make(chan event.Event)
	err := s.routeStore.Watch(server.Context(), ch, routes.WatchOptions{Replay: !request.NoReplay, Monitor: !request.NoSubscribe})

	if err != nil {
		return err
	}

	for routeEvent := range ch {
		response := &modelapi.WatchRoutesResponse{
			Route: routeToAPI(routeEvent.Value.(*model.Route)),
			Type:  eventType(routeEvent.Type.(routes.RouteEvent)),
		}
		err := server.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}
