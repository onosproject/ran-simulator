// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

// Package gnmi implements the northbound gNMI service for the configuration subsystem.
package gnmi

import (
	"context"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
)

var log = liblog.GetLogger("northbound", "gnmi")

// Service implements Service for GNMI
type Service struct {
	northbound.Service
	Port      int
	TowerEcID types.EcID
	PlmnID    types.PlmnID
}

// Register registers the GNMI server with grpc
func (s Service) Register(r *grpc.Server) {
	gnmi.RegisterGNMIServer(r, &Server{port: s.Port, towerEcID: s.TowerEcID, plmnID: s.PlmnID})
}

// Server implements the grpc GNMI service
type Server struct {
	port      int
	towerEcID types.EcID
	plmnID    types.PlmnID
}

// Capabilities implements gNMI Capabilities
func (s *Server) Capabilities(ctx context.Context, req *gnmi.CapabilityRequest) (*gnmi.CapabilityResponse, error) {
	log.Infof("gNMI Capabilities requested for %s-%s", s.GetPlmnID(), s.GetEcID())
	return &gnmi.CapabilityResponse{
		SupportedModels: []*gnmi.ModelData{
			{
				Name:         "e2node",
				Organization: "Open Networking Foundation",
				Version:      "2020-05-01",
			},
		},
		SupportedEncodings: []gnmi.Encoding{gnmi.Encoding_PROTO},
		GNMIVersion:        "0.6.0",
	}, nil
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

func newEcgi(id types.EcID, plmnID types.PlmnID) types.ECGI {
	return types.ECGI{EcID: id, PlmnID: plmnID}
}
