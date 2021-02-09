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
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("trafficsim")

// NewService returns a new trafficsim Service
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
	simapi.RegisterTrafficServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	model *model.Model
}

func coordToAPI(coord model.Coordinate) *simtypes.Point {
	return &types.Point{Lat: coord.Lat, Lng: coord.Lng}
}

// GetMapLayout :
func (s *Server) GetMapLayout(ctx context.Context, req *simapi.MapLayoutRequest) (*types.MapLayout, error) {
	return &types.MapLayout{
		Center:         coordToAPI(s.model.MapLayout.Center),
		Zoom:           s.model.MapLayout.Zoom,
		Fade:           s.model.MapLayout.FadeMap,
		ShowRoutes:     s.model.MapLayout.ShowRoutes,
		ShowPower:      s.model.MapLayout.ShowPower,
		LocationsScale: s.model.MapLayout.LocationsScale,
		MinUes:         0,
		MaxUes:         0,
		CurrentRoutes:  0,
	}, nil
}

// ListRoutes :
func (s *Server) ListRoutes(req *simapi.ListRoutesRequest, stream simapi.Traffic_ListRoutesServer) error {
	return nil
}

func cellToAPI(cell model.Cell) *simtypes.Cell {
	r := &simtypes.Cell{
		Ecgi:       simtypes.Ecgi(cell.Ecgi),
		Location:   nil,
		Sector:     nil,
		Color:      "",
		MaxUEs:     0,
		Neighbors:  nil,
		TxPowerdB:  0,
		CrntiMap:   nil,
		CrntiIndex: 0,
		Port:       0,
	}
	return r
}

// ListCells :
func (s *Server) ListCells(req *simapi.ListCellsRequest, stream simapi.Traffic_ListCellsServer) error {
	for _, node := range s.model.Nodes {
		for _, cell := range node.Cells {
			resp := &simapi.ListCellsResponse{
				Cell: cellToAPI(cell),
				Type: simapi.Type_NONE,
			}
			err := stream.Send(resp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ueToAPI(ue *model.UE) *simtypes.Ue {
	r := &simtypes.Ue{
		Imsi:     simtypes.Imsi(ue.IMSI),
		Type:     string(ue.Type),
		Position: nil,
		Rotation: ue.Rotation,
		Crnti:    simtypes.Crnti(ue.Crnti),
		Admitted: ue.IsAdmitted,
	}
	if ue.Cell != nil {
		r.ServingTower = simtypes.Ecgi(ue.Cell.ID)
		r.ServingTowerStrength = ue.Cell.Strength
	}
	if len(ue.Cells) > 0 {
		r.Tower1 = simtypes.Ecgi(ue.Cells[0].ID)
		r.Tower1Strength = ue.Cells[0].Strength
	}
	if len(ue.Cells) > 1 {
		r.Tower2 = simtypes.Ecgi(ue.Cells[1].ID)
		r.Tower2Strength = ue.Cells[1].Strength
	}
	if len(ue.Cells) > 2 {
		r.Tower3 = simtypes.Ecgi(ue.Cells[2].ID)
		r.Tower3Strength = ue.Cells[2].Strength
	}
	return r
}

// ListUes :
func (s *Server) ListUes(req *simapi.ListUesRequest, stream simapi.Traffic_ListUesServer) error {
	if !req.WithoutReplay {
		for _, ue := range s.model.UEs.ListAllUEs() {
			resp := &simapi.ListUesResponse{
				Ue:   ueToAPI(ue),
				Type: simapi.Type_NONE,
			}
			err := stream.Send(resp)
			if err != nil {
				return err
			}
		}
	}

	// TODO: add subscription flag processing
	return nil
}

// SetNumberUEs changes the number of UEs in the simulation
func (s *Server) SetNumberUEs(ctx context.Context, req *simapi.SetNumberUEsRequest) (*simapi.SetNumberUEsResponse, error) {
	ueCount := req.GetNumber()
	log.Infof("Number of simulated UEs changed to %d", ueCount)
	s.model.UEs.SetUECount(uint(ueCount))
	return &simapi.SetNumberUEsResponse{Number: ueCount}, nil
}

// ResetMetrics resets the metrics on demand
func (s *Server) ResetMetrics(ctx context.Context, req *simapi.ResetMetricsMsg) (*simapi.ResetMetricsMsg, error) {
	return nil, nil
}
