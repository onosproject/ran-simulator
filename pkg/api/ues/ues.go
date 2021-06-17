// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package ues

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/ues"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("api", "ues")

// NewService returns a new model Service
func NewService(ueStore ues.Store) service.Service {
	return &Service{
		ueStore: ueStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	ueStore ues.Store
}

// Register registers the TrafficSim Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		ueStore: s.ueStore,
	}
	modelapi.RegisterUEModelServer(r, server)
}

// Server implements the Routes gRPC service for administrative facilities.
type Server struct {
	ueStore ues.Store
}

// GetUECount gets the number of UEs
func (s *Server) GetUECount(ctx context.Context, request *modelapi.GetUECountRequest) (*modelapi.GetUECountResponse, error) {
	panic("implement me")
}

// SetUECount sets the number of UEs
func (s *Server) SetUECount(ctx context.Context, request *modelapi.SetUECountRequest) (*modelapi.SetUECountResponse, error) {
	panic("implement me")
}

func ueToAPI(ue *model.UE) *types.Ue {
	return &types.Ue{
		IMSI:                 ue.IMSI,
		ServingTower:         ue.Cell.NCGI,
		ServingTowerStrength: ue.Cell.Strength,
		CRNTI:                ue.CRNTI,
		Admitted:             false,
		Metrics:              nil,
		RrcState:             uint32(ue.RrcState),
	}
}

// GetUE returns information on the specified UE
func (s *Server) GetUE(ctx context.Context, request *modelapi.GetUERequest) (*modelapi.GetUEResponse, error) {
	log.Debugf("Received get UE request: %+v", request)
	ue, err := s.ueStore.Get(ctx, request.IMSI)
	if err != nil {
		return nil, err
	}
	return &modelapi.GetUEResponse{Ue: ueToAPI(ue)}, nil
}

// MoveToCell moves the specified UE to the given cell
func (s *Server) MoveToCell(ctx context.Context, request *modelapi.MoveToCellRequest) (*modelapi.MoveToCellResponse, error) {
	panic("implement me")
}

// MoveToLocation moves the specified UE to the given location
func (s *Server) MoveToLocation(ctx context.Context, request *modelapi.MoveToLocationRequest) (*modelapi.MoveToLocationResponse, error) {
	panic("implement me")
}

// DeleteUE removes the specified UE
func (s *Server) DeleteUE(ctx context.Context, request *modelapi.DeleteUERequest) (*modelapi.DeleteUEResponse, error) {
	panic("implement me")
}

// WatchUEs returns events pertaining to changes in the UE state.
func (s *Server) WatchUEs(request *modelapi.WatchUEsRequest, server modelapi.UEModel_WatchUEsServer) error {
	panic("implement me")
}

// ListUEs returns list of simulated UEs.
func (s *Server) ListUEs(request *modelapi.ListUEsRequest, server modelapi.UEModel_ListUEsServer) error {
	log.Debugf("Received listing UEs request: %v", request)
	ueList := s.ueStore.ListAllUEs(server.Context())
	for _, ue := range ueList {
		log.Info("UE:", ue)
		resp := &modelapi.ListUEsResponse{
			Ue: ueToAPI(ue),
		}
		err := server.Send(resp)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
