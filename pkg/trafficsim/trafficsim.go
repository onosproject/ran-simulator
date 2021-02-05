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

// GetMapLayout :
func (s *Server) GetMapLayout(ctx context.Context, req *simapi.MapLayoutRequest) (*types.MapLayout, error) {
	return nil, nil
}

// ListRoutes :
func (s *Server) ListRoutes(req *simapi.ListRoutesRequest, stream simapi.Traffic_ListRoutesServer) error {
	return nil
}

// ListCells :
func (s *Server) ListCells(req *simapi.ListCellsRequest, stream simapi.Traffic_ListCellsServer) error {
	return nil
}

func genbToAPI(id model.GEnbID) *simtypes.ECGI {
	return &simtypes.ECGI{
		EcID:   simtypes.EcID(id.EnbID),
		PlmnID: simtypes.PlmnID(id.PlmnID),
	}
}
func ueToAPI(ue *model.UE) *simtypes.Ue {
	r := &simtypes.Ue{
		Imsi:     simtypes.Imsi(ue.Imsi),
		Type:     string(ue.Type),
		Position: nil,
		Rotation: ue.Rotation,
		Crnti:    simtypes.Crnti(ue.Crnti),
		Admitted: ue.IsAdmitted,
	}
	if ue.Tower != nil {
		r.ServingTower = genbToAPI(ue.Tower.ID)
		r.ServingTowerStrength = ue.Tower.Strength
	}
	if len(ue.Towers) > 0 {
		r.Tower1 = genbToAPI(ue.Towers[0].ID)
		r.Tower1Strength = ue.Towers[0].Strength
	}
	if len(ue.Towers) > 1 {
		r.Tower1 = genbToAPI(ue.Towers[1].ID)
		r.Tower1Strength = ue.Towers[1].Strength
	}
	if len(ue.Towers) > 2 {
		r.Tower1 = genbToAPI(ue.Towers[2].ID)
		r.Tower1Strength = ue.Towers[2].Strength
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
