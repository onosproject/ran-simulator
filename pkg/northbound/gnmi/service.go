// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package gnmi implements the northbound gNMI service for the configuration subsystem.
package gnmi

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"io/ioutil"
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
	log.Info("gNMI Capabilities requested for %s-%s", s.GetPlmnID(), s.GetEcID())
	v, _ := getGNMIServiceVersion()
	return &gnmi.CapabilityResponse{
		SupportedModels: []*gnmi.ModelData{
			{
				Name:         "e2node",
				Organization: "Open Networking Foundation",
				Version:      "2020-05-01",
			},
		},
		SupportedEncodings: []gnmi.Encoding{gnmi.Encoding_PROTO},
		GNMIVersion:        *v,
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

// getGNMIServiceVersion returns a pointer to the gNMI service version string.
// The method is non-trivial because of the way it is defined in the proto file.
func getGNMIServiceVersion() (*string, error) {
	gzB, _ := (&gnmi.Update{}).Descriptor()
	r, err := gzip.NewReader(bytes.NewReader(gzB))
	if err != nil {
		return nil, fmt.Errorf("error in initializing gzip reader: %v", err)
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error in reading gzip data: %v", err)
	}
	desc := &descriptor.FileDescriptorProto{}
	if err := proto.Unmarshal(b, desc); err != nil {
		return nil, fmt.Errorf("error in unmarshaling proto: %v", err)
	}
	ver, err := proto.GetExtension(desc.Options, gnmi.E_GnmiService)
	if err != nil {
		return nil, fmt.Errorf("error in getting version from proto extension: %v", err)
	}
	return ver.(*string), nil
}
