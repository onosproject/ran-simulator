// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package cells

import (
	"context"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger()

// ManagementDelegate provides means to clear and load the model and resume the simulation
type ManagementDelegate interface {
	// PauseAndClear pauses simulation and clears the model
	PauseAndClear(ctx context.Context)

	// LoadModel loads the new model into the simulator
	LoadModel(ctx context.Context, modelData []byte) error

	// LoadMetrics loads new metrics into the simulator
	LoadMetrics(ctx context.Context, name string, metricsData []byte) error

	// Resume resume the simulation
	Resume(ctx context.Context)
}

// NewService returns a new model Service
func NewService(delegate ManagementDelegate) service.Service {
	return &Service{
		delegate: delegate,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	delegate ManagementDelegate
}

// Register registers the ModelService with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		delegate: s.delegate,
	}
	modelapi.RegisterModelServiceServer(r, server)
}

// Server implements the ModelService gRPC service
type Server struct {
	delegate ManagementDelegate
}

// Load loads new data sets into the simulator
func (s *Server) Load(ctx context.Context, request *modelapi.LoadRequest) (*modelapi.LoadResponse, error) {
	log.Debugf("Received model load request: %v", request)

	// Stop simulation and clear model
	s.delegate.PauseAndClear(ctx)

	// Load all new data sets
	for _, ds := range request.DataSet {
		if ds.Type == "model" {
			if err := s.delegate.LoadModel(ctx, ds.Data); err != nil {
				return nil, err
			}
		} else {
			if err := s.delegate.LoadMetrics(ctx, ds.Type, ds.Data); err != nil {
				return nil, err
			}
		}
	}

	if request.Resume {
		s.delegate.Resume(ctx)
	}

	return &modelapi.LoadResponse{}, nil
}

// Clear clears model data
func (s *Server) Clear(ctx context.Context, request *modelapi.ClearRequest) (*modelapi.ClearResponse, error) {
	log.Debugf("Received model clear request: %v", request)
	s.delegate.PauseAndClear(ctx)
	if request.Resume {
		s.delegate.Resume(ctx)
	}
	return &modelapi.ClearResponse{}, nil
}
