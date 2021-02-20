// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cells

import (
	"context"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("api", "cells")

// NewService returns a new model Service
func NewService(cellStore cells.CellRegistry) service.Service {
	return &Service{
		cellStore: cellStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	cellStore cells.CellRegistry
}

// Register registers the TrafficSim Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		cellStore: s.cellStore,
	}
	modelapi.RegisterCellModelServer(r, server)
}

// Server implements the TrafficSim gRPC service for administrative facilities.
type Server struct {
	cellStore cells.CellRegistry
}

func cellToAPI(cell *model.Cell) *types.Cell {
	sector := sectorToAPI(cell.Sector)
	return &types.Cell{
		ECGI:      cell.ECGI,
		Location:  sector.Centroid,
		Sector:    sector,
		Color:     cell.Color,
		MaxUEs:    cell.MaxUEs,
		Neighbors: cell.Neighbors,
		TxPowerdB: cell.TxPowerDB,
	}
}

func cellToModel(cell *types.Cell) *model.Cell {
	return &model.Cell{
		ECGI: cell.ECGI,
		Sector: model.Sector{
			Center: model.Coordinate{Lat: cell.Sector.Centroid.Lat, Lng: cell.Sector.Centroid.Lng},
			Arc:    cell.Sector.Arc, Azimuth: cell.Sector.Azimuth,
		},
		Color:     cell.Color,
		MaxUEs:    cell.MaxUEs,
		Neighbors: cell.Neighbors,
		TxPowerDB: cell.TxPowerdB,
	}
}

func sectorToAPI(sector model.Sector) *types.Sector {
	return &types.Sector{
		Azimuth:  sector.Azimuth,
		Arc:      sector.Arc,
		Centroid: &types.Point{Lat: sector.Center.Lat, Lng: sector.Center.Lng},
	}
}

// CreateCell creates a new simulated cell
func (s *Server) CreateCell(ctx context.Context, request *modelapi.CreateCellRequest) (*modelapi.CreateCellResponse, error) {
	err := s.cellStore.AddCell(cellToModel(request.Cell))
	if err != nil {
		return nil, err
	}
	return &modelapi.CreateCellResponse{}, nil
}

// GetCell retrieves the specified simulated cell
func (s *Server) GetCell(ctx context.Context, request *modelapi.GetCellRequest) (*modelapi.GetCellResponse, error) {
	node, err := s.cellStore.GetCell(request.ECGI)
	if err != nil {
		return nil, err
	}
	return &modelapi.GetCellResponse{Cell: cellToAPI(node)}, nil
}

// UpdateCell updates the specified simulated cell
func (s *Server) UpdateCell(ctx context.Context, request *modelapi.UpdateCellRequest) (*modelapi.UpdateCellResponse, error) {
	err := s.cellStore.UpdateCell(cellToModel(request.Cell))
	if err != nil {
		return nil, err
	}
	return &modelapi.UpdateCellResponse{}, nil
}

// DeleteCell deletes the specified simulated cell
func (s *Server) DeleteCell(ctx context.Context, request *modelapi.DeleteCellRequest) (*modelapi.DeleteCellResponse, error) {
	_, err := s.cellStore.DeleteCell(request.ECGI)
	if err != nil {
		return nil, err
	}
	return &modelapi.DeleteCellResponse{}, nil
}

func eventType(t uint8) modelapi.EventType {
	if t == cells.ADDED {
		return modelapi.EventType_CREATED
	} else if t == cells.UPDATED {
		return modelapi.EventType_UPDATED
	} else if t == cells.DELETED {
		return modelapi.EventType_DELETED
	} else {
		return modelapi.EventType_NONE
	}
}

// WatchCells monitors changes to the inventory of cells
func (s *Server) WatchCells(request *modelapi.WatchCellsRequest, server modelapi.CellModel_WatchCellsServer) error {
	log.Infof("Watching cells [%v]...", request)
	ch := make(chan cells.CellEvent)
	s.cellStore.WatchCells(ch, cells.WatchOptions{Replay: !request.NoReplay, Monitor: !request.NoSubscribe})

	for event := range ch {
		response := &modelapi.WatchCellsResponse{
			Cell: cellToAPI(event.Cell),
			Type: eventType(event.Type),
		}
		err := server.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}
