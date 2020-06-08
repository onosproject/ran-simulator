// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package e2

import (
	"context"
	"fmt"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/ran-simulator/api/types"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("northbound", "e2")

// Service is an implementation of e2 service.
type Service struct {
	service.Service
	Port      int
	TowerEcID types.EcID
	PlmnID    types.PlmnID
}

// Register registers the e2 Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{port: s.Port, towerEcID: s.TowerEcID, plmnID: s.PlmnID,
		l2MeasConfig: e2.L2MeasConfig{RadioMeasReportPerUe: e2.L2MeasReportInterval_MS_500}}
	e2ap.RegisterE2APServer(r, server)
}

// Server implements the E2 gRPC service for administrative facilities.
type Server struct {
	port         int
	towerEcID    types.EcID
	plmnID       types.PlmnID
	l2MeasConfig e2.L2MeasConfig
}

// RicSubscribe - add tot he list of subscriptions
func (s *Server) RicSubscribe(ctx context.Context, req *e2ap.RicSubscriptionRequest) (*e2ap.RicSubscriptionResponse, error) {

	return nil, fmt.Errorf("not yet implemented")
}

// GetPort - expose the Port number
func (s *Server) GetPort() int {
	return s.port
}

// GetEcID - expose the tower EcID
func (s *Server) GetEcID() types.EcID {
	return s.towerEcID
}

// GetPlmnID - expose the tower PlmnID
func (s *Server) GetPlmnID() types.PlmnID {
	return s.plmnID
}

// GetECGI - expose the tower Ecgi
func (s *Server) GetECGI() types.ECGI {
	return newEcgi(s.GetEcID(), s.GetPlmnID())
}

func toTypesEcgi(e2Ecgi *e2.ECGI) types.ECGI {
	return types.ECGI{
		EcID:   types.EcID(e2Ecgi.Ecid),
		PlmnID: types.PlmnID(e2Ecgi.PlmnId),
	}
}

func toE2Ecgi(e2Ecgi *types.ECGI) e2.ECGI {
	return e2.ECGI{
		Ecid:   string(e2Ecgi.EcID),
		PlmnId: string(e2Ecgi.PlmnID),
	}
}

func newEcgi(id types.EcID, plmnID types.PlmnID) types.ECGI {
	return types.ECGI{EcID: id, PlmnID: plmnID}
}
