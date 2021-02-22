// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package trafficsim

import (
	"context"

	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	simapi "github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	simtypes "github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("trafficsim")

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
func (s *Server) GetMapLayout(ctx context.Context, req *simapi.MapLayoutRequest) (*types.MapLayout, error) {
	return &types.MapLayout{
		Center:         &types.Point{Lat: s.model.MapLayout.Center.Lat, Lng: s.model.MapLayout.Center.Lng},
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
		Rotation: ue.Rotation,
		CRNTI:    ue.CRNTI,
		Admitted: ue.IsAdmitted,
	}
	if ue.Cell != nil {
		r.ServingTower = simtypes.ECGI(ue.Cell.ID)
		r.ServingTowerStrength = ue.Cell.Strength
	}
	if len(ue.Cells) > 0 {
		r.Tower1 = simtypes.ECGI(ue.Cells[0].ID)
		r.Tower1Strength = ue.Cells[0].Strength
	}
	if len(ue.Cells) > 1 {
		r.Tower2 = simtypes.ECGI(ue.Cells[1].ID)
		r.Tower2Strength = ue.Cells[1].Strength
	}
	if len(ue.Cells) > 2 {
		r.Tower3 = simtypes.ECGI(ue.Cells[2].ID)
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
func (s *Server) ListUes(req *simapi.ListUesRequest, stream simapi.Traffic_ListUesServer) error {
	ch := make(chan ues.UEEvent)
	log.Infof("UE Store: %v", s.ueStore)
	s.ueStore.WatchUEs(ch, ues.WatchOptions{Replay: !req.WithoutReplay, Monitor: req.Subscribe})

	for event := range ch {
		resp := &simapi.ListUesResponse{
			Ue:   ueToAPI(event.UE),
			Type: eventType(event.Type),
		}
		err := stream.Send(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func eventType(t uint8) simapi.Type {
	if t == ues.ADDED {
		return simapi.Type_ADDED
	} else if t == ues.UPDATED {
		return simapi.Type_ADDED
	} else if t == ues.DELETED {
		return simapi.Type_ADDED
	} else {
		return simapi.Type_NONE
	}
}

// SetNumberUEs changes the number of UEs in the simulation
func (s *Server) SetNumberUEs(ctx context.Context, req *simapi.SetNumberUEsRequest) (*simapi.SetNumberUEsResponse, error) {
	ueCount := req.GetNumber()
	log.Infof("Number of simulated UEs changed to %d", ueCount)
	s.ueStore.SetUECount(uint(ueCount))
	return &simapi.SetNumberUEsResponse{Number: ueCount}, nil
}

// ResetMetrics resets the metrics on demand
func (s *Server) ResetMetrics(ctx context.Context, req *simapi.ResetMetricsMsg) (*simapi.ResetMetricsMsg, error) {
	// TODO: Reimplement
	return nil, nil
}
