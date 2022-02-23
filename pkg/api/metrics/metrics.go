// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package nodes

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"

	metricsapi "github.com/onosproject/onos-api/go/onos/ransim/metrics"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger()

// NewService returns a new metrics Service
func NewService(metricsStore metrics.Store) service.Service {
	return &Service{
		store: metricsStore,
	}
}

// Service is a Service implementation for administration.
type Service struct {
	service.Service
	store metrics.Store
}

// Register registers the TrafficSim Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		s.store,
	}
	metricsapi.RegisterMetricsServiceServer(r, server)
}

// Server implements the metrics gRPC service
type Server struct {
	store metrics.Store
}

func bitWidth(t string) int {
	if strings.HasSuffix(t, "64") {
		return 64
	} else if strings.HasSuffix(t, "32") {
		return 32
	} else if strings.HasSuffix(t, "16") {
		return 16
	} else if strings.HasSuffix(t, "8") {
		return 8
	} else {
		return 64
	}
}

func value(v string, t string) interface{} {
	if strings.HasPrefix(t, "uint") {
		i, err := strconv.ParseUint(v, 10, bitWidth(t))
		if err == nil {
			return i
		}
	} else if strings.HasPrefix(t, "int") {
		i, err := strconv.ParseInt(v, 10, bitWidth(t))
		if err == nil {
			return i
		}
	} else if strings.HasPrefix(t, "float") {
		i, err := strconv.ParseFloat(v, bitWidth(t))
		if err == nil {
			return i
		}
	} else if strings.HasPrefix(t, "bool") {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return v
}

func valueType(v interface{}) (string, string) {
	return fmt.Sprintf("%v", v), reflect.TypeOf(v).Kind().String()
}

// List lists all metrics of the specified entity
func (s *Server) List(ctx context.Context, request *metricsapi.ListRequest) (*metricsapi.ListResponse, error) {
	log.Debugf("Received list metrics request: %+v", request)
	mmap, err := s.store.List(ctx, request.EntityID)
	if err != nil {
		return nil, err
	}

	mlist := make([]*metricsapi.Metric, 0, len(mmap))
	for k, v := range mmap {
		mlist = append(mlist, metricToAPI(request.EntityID, k, v))
	}

	return &metricsapi.ListResponse{
		EntityID: request.EntityID,
		Metrics:  mlist,
	}, nil
}

func metricToAPI(entityID uint64, name string, v interface{}) *metricsapi.Metric {
	vv, vt := valueType(v)
	return &metricsapi.Metric{EntityID: entityID, Key: name, Value: vv, Type: vt}
}

// Set sets the value of the specified metric
func (s *Server) Set(ctx context.Context, request *metricsapi.SetRequest) (*metricsapi.SetResponse, error) {
	log.Debugf("Received set metric request: %+v", request)
	m := request.Metric
	err := s.store.Set(ctx, m.EntityID, m.Key, value(m.Value, m.Type))
	if err != nil {
		return nil, err
	}
	return &metricsapi.SetResponse{}, nil
}

// Get retrieves the value of the specified metric
func (s *Server) Get(ctx context.Context, request *metricsapi.GetRequest) (*metricsapi.GetResponse, error) {
	log.Debugf("Received get metric request: %+v", request)
	if m, ok := s.store.Get(ctx, request.EntityID, request.Name); ok {
		return &metricsapi.GetResponse{Metric: metricToAPI(request.EntityID, request.Name, m)}, nil
	}
	return nil, errors.New(errors.NotFound, "metric not found")
}

// Delete deletes the value of the specified metric
func (s *Server) Delete(ctx context.Context, request *metricsapi.DeleteRequest) (*metricsapi.DeleteResponse, error) {
	log.Debugf("Received delete metric request: %+v", request)
	err := s.store.Delete(ctx, request.EntityID, request.Name)
	if err != nil {
		return nil, err
	}
	return &metricsapi.DeleteResponse{}, nil
}

// DeleteAll deletes all metrics of the specified entity
func (s *Server) DeleteAll(ctx context.Context, request *metricsapi.DeleteAllRequest) (*metricsapi.DeleteAllResponse, error) {
	log.Debugf("Received delete all metric request: %+v", request)
	err := s.store.DeleteAll(ctx, request.EntityID)
	if err != nil {
		return nil, err
	}
	return &metricsapi.DeleteAllResponse{}, nil
}

func eventType(metricsEvent metrics.MetricEvent) metricsapi.EventType {
	if metricsEvent == metrics.Deleted {
		return metricsapi.EventType_DELETED
	}
	return metricsapi.EventType_UPDATED
}

// Watch watches for all metric updates and deletes
func (s *Server) Watch(request *metricsapi.WatchRequest, server metricsapi.MetricsService_WatchServer) error {
	log.Debugf("Received watch metrics request: %+v", request)
	ch := make(chan event.Event)
	err := s.store.Watch(server.Context(), ch)
	if err != nil {
		return err
	}

	for metricsEvent := range ch {
		key := metricsEvent.Key.(metrics.Key)
		response := &metricsapi.WatchResponse{Type: eventType(metricsEvent.Type.(metrics.MetricEvent))}
		if metricsEvent.Type == metrics.Deleted {
			response.Metric = metricToAPI(key.EntityID, key.Name, "")
		} else {
			response.Metric = metricToAPI(key.EntityID, key.Name, metricsEvent.Value)
		}
		err := server.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}
