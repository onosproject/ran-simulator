// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package ues

import (
	"context"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/onosproject/ran-simulator/pkg/store/ues"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger()

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
	return &modelapi.GetUECountResponse{Count: uint32(s.ueStore.Len(ctx))}, nil
}

// SetUECount sets the number of UEs
func (s *Server) SetUECount(ctx context.Context, request *modelapi.SetUECountRequest) (*modelapi.SetUECountResponse, error) {
	s.ueStore.SetUECount(ctx, uint(request.Count))
	return &modelapi.SetUECountResponse{}, nil
}

func ueToAPI(ue *model.UE) *types.Ue {
	r := &types.Ue{
		IMSI:     ue.IMSI,
		Type:     string(ue.Type),
		Position: &types.Point{Lat: ue.Location.Lat, Lng: ue.Location.Lng},
		Rotation: ue.Heading,
		CRNTI:    ue.CRNTI,
		Admitted: ue.IsAdmitted,
		RrcState: uint32(ue.RrcState),
		Metrics:  nil,
		FiveQi:   int32(ue.FiveQi),
	}
	if ue.Cell != nil {
		r.ServingTower = ue.Cell.NCGI
		r.ServingTowerStrength = ue.Cell.Strength
	}
	if len(ue.Cells) > 0 {
		r.Tower1 = ue.Cells[0].NCGI
		r.Tower1Strength = ue.Cells[0].Strength
	}
	if len(ue.Cells) > 1 {
		r.Tower2 = ue.Cells[1].NCGI
		r.Tower2Strength = ue.Cells[1].Strength
	}
	if len(ue.Cells) > 2 {
		r.Tower3 = ue.Cells[2].NCGI
		r.Tower3Strength = ue.Cells[2].Strength
	}
	return r
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
	log.Infof("Received MoveToCell request: %+v", request)
	err := s.ueStore.MoveToCell(ctx, request.IMSI, request.NCGI, 0)
	if err != nil {
		return nil, err
	}
	return &modelapi.MoveToCellResponse{}, nil
}

// MoveToLocation moves the specified UE to the given location
func (s *Server) MoveToLocation(ctx context.Context, request *modelapi.MoveToLocationRequest) (*modelapi.MoveToLocationResponse, error) {
	log.Debugf("Received MoveToLocation request: %+v", request)
	return &modelapi.MoveToLocationResponse{}, s.ueStore.MoveToCoordinate(ctx, request.IMSI, model.Coordinate(*request.Location), request.Heading)
}

// DeleteUE removes the specified UE
func (s *Server) DeleteUE(ctx context.Context, request *modelapi.DeleteUERequest) (*modelapi.DeleteUEResponse, error) {
	log.Debugf("Received Delete request: %+v", request)
	_, err := s.ueStore.Delete(ctx, request.IMSI)
	return &modelapi.DeleteUEResponse{}, err
}

func eventType(ueEvent ues.UeEvent) modelapi.EventType {
	if ueEvent == ues.Created {
		return modelapi.EventType_CREATED
	} else if ueEvent == ues.Updated {
		return modelapi.EventType_UPDATED
	} else if ueEvent == ues.Deleted {
		return modelapi.EventType_DELETED
	} else {
		return modelapi.EventType_NONE
	}
}

// WatchUEs returns events pertaining to changes in the UE state.
func (s *Server) WatchUEs(request *modelapi.WatchUEsRequest, server modelapi.UEModel_WatchUEsServer) error {
	log.Debugf("Received WatchUEs request: %+v", request)
	ch := make(chan event.Event)
	err := s.ueStore.Watch(server.Context(), ch, ues.WatchOptions{Replay: !request.NoReplay, Monitor: !request.NoSubscribe})

	if err != nil {
		return err
	}

	for ueEvent := range ch {
		response := &modelapi.WatchUEsResponse{
			Ue:   ueToAPI(ueEvent.Value.(*model.UE)),
			Type: eventType(ueEvent.Type.(ues.UeEvent)),
		}
		err := server.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListUEs returns list of simulated UEs.
func (s *Server) ListUEs(request *modelapi.ListUEsRequest, server modelapi.UEModel_ListUEsServer) error {
	log.Debugf("Received listing UEs request: %v", request)
	ueList := s.ueStore.ListAllUEs(server.Context())
	for _, ue := range ueList {
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
