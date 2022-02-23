// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package cells

import (
	"context"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger()

// NewService returns a new model Service
func NewService(cellStore cells.Store) service.Service {
	return &Service{
		cellStore: cellStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	cellStore cells.Store
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
	cellStore cells.Store
}

func cellToAPI(cell *model.Cell) *types.Cell {
	sector := sectorToAPI(cell.Sector)
	measurementParams := measurementParamsToAPI(cell.MeasurementParams)
	return &types.Cell{
		NCGI:              cell.NCGI,
		Location:          sector.Centroid,
		Sector:            sector,
		Color:             cell.Color,
		MaxUEs:            cell.MaxUEs,
		Neighbors:         cell.Neighbors,
		TxPowerdB:         cell.TxPowerDB,
		MeasurementParams: measurementParams,
		RrcIdleCount:      cell.RrcIdleCount,
		RrcConnectedCount: cell.RrcConnectedCount,
		Pci:               cell.PCI,
	}
}

func cellToModel(cell *types.Cell) *model.Cell {
	return &model.Cell{
		NCGI: cell.NCGI,
		Sector: model.Sector{
			Center:  model.Coordinate{Lat: cell.Sector.Centroid.Lat, Lng: cell.Sector.Centroid.Lng},
			Arc:     cell.Sector.Arc,
			Azimuth: cell.Sector.Azimuth,
			Tilt:    cell.Sector.Tilt,
			Height:  cell.Sector.Height,
		},
		Color:     cell.Color,
		MaxUEs:    cell.MaxUEs,
		Neighbors: cell.Neighbors,
		TxPowerDB: cell.TxPowerdB,
		MeasurementParams: model.MeasurementParams{
			TimeToTrigger:          cell.MeasurementParams.TimeToTrigger,
			FrequencyOffset:        cell.MeasurementParams.FrequencyOffset,
			PCellIndividualOffset:  cell.MeasurementParams.PcellIndividualOffset,
			NCellIndividualOffsets: cell.MeasurementParams.NcellIndividualOffsets,
			Hysteresis:             cell.MeasurementParams.Hysteresis,
			EventA3Params: model.EventA3Params{
				A3Offset:      cell.MeasurementParams.EventA3Params.A3Offset,
				ReportOnLeave: cell.MeasurementParams.EventA3Params.ReportOnLeave,
			},
		},
		PCI: cell.Pci,
	}
}

func sectorToAPI(sector model.Sector) *types.Sector {
	return &types.Sector{
		Azimuth:  sector.Azimuth,
		Arc:      sector.Arc,
		Centroid: &types.Point{Lat: sector.Center.Lat, Lng: sector.Center.Lng},
		Tilt:     sector.Tilt,
		Height:   sector.Height,
	}
}

func measurementParamsToAPI(params model.MeasurementParams) *types.MeasurementParams {
	return &types.MeasurementParams{
		TimeToTrigger:          params.TimeToTrigger,
		FrequencyOffset:        params.FrequencyOffset,
		PcellIndividualOffset:  params.PCellIndividualOffset,
		NcellIndividualOffsets: params.NCellIndividualOffsets,
		Hysteresis:             params.Hysteresis,
		EventA3Params:          eventA3ParamsToAPI(params.EventA3Params),
	}
}

func eventA3ParamsToAPI(params model.EventA3Params) *types.EventA3Params {
	return &types.EventA3Params{
		A3Offset:      params.A3Offset,
		ReportOnLeave: params.ReportOnLeave,
	}
}

// CreateCell creates a new simulated cell
func (s *Server) CreateCell(ctx context.Context, request *modelapi.CreateCellRequest) (*modelapi.CreateCellResponse, error) {
	log.Debugf("Received create cell request: %v", request)
	err := s.cellStore.Add(ctx, cellToModel(request.Cell))
	if err != nil {
		return nil, err
	}
	return &modelapi.CreateCellResponse{}, nil
}

// GetCell retrieves the specified simulated cell
func (s *Server) GetCell(ctx context.Context, request *modelapi.GetCellRequest) (*modelapi.GetCellResponse, error) {
	log.Debugf("Received get cell request: %v", request)
	node, err := s.cellStore.Get(ctx, request.NCGI)
	if err != nil {
		return nil, err
	}
	return &modelapi.GetCellResponse{Cell: cellToAPI(node)}, nil
}

// UpdateCell updates the specified simulated cell
func (s *Server) UpdateCell(ctx context.Context, request *modelapi.UpdateCellRequest) (*modelapi.UpdateCellResponse, error) {
	log.Debugf("Received update cell request: %v", request)
	err := s.cellStore.Update(ctx, cellToModel(request.Cell))
	if err != nil {
		return nil, err
	}
	return &modelapi.UpdateCellResponse{}, nil
}

// DeleteCell deletes the specified simulated cell
func (s *Server) DeleteCell(ctx context.Context, request *modelapi.DeleteCellRequest) (*modelapi.DeleteCellResponse, error) {
	log.Debugf("Received delete cell request: %v", request)
	_, err := s.cellStore.Delete(ctx, request.NCGI)
	if err != nil {
		return nil, err
	}
	return &modelapi.DeleteCellResponse{}, nil
}

func eventType(cellEvent cells.CellEvent) modelapi.EventType {
	if cellEvent == cells.Created {
		return modelapi.EventType_CREATED
	} else if cellEvent == cells.Updated {
		return modelapi.EventType_UPDATED
	} else if cellEvent == cells.Deleted {
		return modelapi.EventType_DELETED
	} else {
		return modelapi.EventType_NONE
	}
}

// ListCells list all of the cells
func (s *Server) ListCells(request *modelapi.ListCellsRequest, server modelapi.CellModel_ListCellsServer) error {
	log.Debugf("Received listing cells request: %v", request)
	cellList, err := s.cellStore.List(server.Context())
	if err != nil {
		return err
	}
	for _, cell := range cellList {
		resp := &modelapi.ListCellsResponse{
			Cell: cellToAPI(cell),
		}
		err = server.Send(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

// WatchCells monitors changes to the inventory of cells
func (s *Server) WatchCells(request *modelapi.WatchCellsRequest, server modelapi.CellModel_WatchCellsServer) error {
	log.Debugf("Received watching cell changes request: %v", request)
	ch := make(chan event.Event)
	err := s.cellStore.Watch(server.Context(), ch, cells.WatchOptions{Replay: !request.NoReplay, Monitor: !request.NoSubscribe})
	if err != nil {
		return err
	}

	for cellEvent := range ch {
		response := &modelapi.WatchCellsResponse{
			Cell: cellToAPI(cellEvent.Value.(*model.Cell)),
			Type: eventType(cellEvent.Type.(cells.CellEvent)),
		}
		err := server.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}
