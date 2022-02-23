// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package trafficsim

import (
	"context"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	simapi "github.com/onosproject/onos-api/go/onos/ransim/trafficsim"

	simtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger()

// NewService returns a new trafficsim Service
func NewService(model *model.Model, cellStore cells.Store, ueStore ues.Store) service.Service {
	return &Service{
		model:     model,
		cellStore: cellStore,
		ueStore:   ueStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	model     *model.Model
	cellStore cells.Store
	ueStore   ues.Store
}

// Register registers the TrafficSim Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		model:     s.model,
		cellStore: s.cellStore,
		ueStore:   s.ueStore,
	}
	simapi.RegisterTrafficServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	model     *model.Model
	cellStore cells.Store
	ueStore   ues.Store
}

// GetMapLayout :
func (s *Server) GetMapLayout(ctx context.Context, req *simapi.MapLayoutRequest) (*simtypes.MapLayout, error) {
	return &simtypes.MapLayout{
		Center:         &simtypes.Point{Lat: s.model.MapLayout.Center.Lat, Lng: s.model.MapLayout.Center.Lng},
		Zoom:           s.model.MapLayout.Zoom,
		Fade:           s.model.MapLayout.FadeMap,
		ShowRoutes:     s.model.MapLayout.ShowRoutes,
		ShowPower:      s.model.MapLayout.ShowPower,
		LocationsScale: s.model.MapLayout.LocationsScale,
	}, nil
}

func ueToAPI(ue *model.UE) *simtypes.Ue {
	r := &simtypes.Ue{
		IMSI:     ue.IMSI,
		Type:     string(ue.Type),
		Position: nil,
		Rotation: ue.Heading,
		CRNTI:    ue.CRNTI,
		Admitted: ue.IsAdmitted,
		RrcState: uint32(ue.RrcState),
	}
	if ue.Cell != nil {
		r.ServingTower = simtypes.NCGI(ue.Cell.ID)
		r.ServingTowerStrength = ue.Cell.Strength
	}
	if len(ue.Cells) > 0 {
		r.Tower1 = simtypes.NCGI(ue.Cells[0].ID)
		r.Tower1Strength = ue.Cells[0].Strength
	}
	if len(ue.Cells) > 1 {
		r.Tower2 = simtypes.NCGI(ue.Cells[1].ID)
		r.Tower2Strength = ue.Cells[1].Strength
	}
	if len(ue.Cells) > 2 {
		r.Tower3 = simtypes.NCGI(ue.Cells[2].ID)
		r.Tower3Strength = ue.Cells[2].Strength
	}
	return r
}

// ListRoutes provides means to list (and optionally monitor) simulated routes
func (s *Server) ListRoutes(req *simapi.ListRoutesRequest, stream simapi.Traffic_ListRoutesServer) error {
	// TODO: reimplement list
	// TODO: add watch capability
	return nil
}

// ListUes provides means to list (and optionally monitor) simulated UEs
func (s *Server) ListUes(request *simapi.ListUesRequest, stream simapi.Traffic_ListUesServer) error {
	log.Debugf("Received listing ues request: %v", request)
	ueList := s.ueStore.ListAllUEs(stream.Context())
	for _, ue := range ueList {
		resp := &simapi.ListUesResponse{
			Ue: ueToAPI(ue),
		}
		log.Infof("UE: %v", ue)
		err := stream.Send(resp)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// WatchUes watch ue changes
func (s *Server) WatchUes(request *simapi.WatchUesRequest, server simapi.Traffic_WatchUesServer) error {
	log.Debugf("Received watching ue changes request: %v", request)
	ch := make(chan event.Event)
	err := s.ueStore.Watch(server.Context(), ch, ues.WatchOptions{Replay: !request.NoReplay})
	if err != nil {
		return err
	}
	for ueEvent := range ch {
		response := &simapi.WatchUesResponse{
			Ue: ueToAPI(ueEvent.Value.(*model.UE)),
		}
		err := server.Send(response)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetNumberUEs changes the number of UEs in the simulation
func (s *Server) SetNumberUEs(ctx context.Context, req *simapi.SetNumberUEsRequest) (*simapi.SetNumberUEsResponse, error) {
	ueCount := req.GetNumber()
	log.Infof("Number of simulated UEs changed to %d", ueCount)
	s.ueStore.SetUECount(ctx, uint(ueCount))
	return &simapi.SetNumberUEsResponse{Number: ueCount}, nil
}

// ResetMetrics resets the metrics on demand
func (s *Server) ResetMetrics(ctx context.Context, req *simapi.ResetMetricsMsg) (*simapi.ResetMetricsMsg, error) {
	// TODO: Reimplement
	return nil, nil
}
