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

// Package device :
package device

import (
	"context"
	devicestore "github.com/onosproject/onos-topo/pkg/store/device"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-lib-go/pkg/northbound"
	deviceapi "github.com/onosproject/onos-topo/api/device"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"regexp"
	"time"
)

var log = logging.GetLogger("northbound", "device")

const (
	defaultTimeout       = 5 * time.Second
	deviceNamePattern    = `^[a-zA-Z0-9\-:_]{4,40}$`
	deviceAddressPattern = `^[a-zA-Z0-9\-_\.]+:[0-9]+$`
	deviceVersionPattern = `^(\d+(\.\d+){2,3})$`
	deviceAttrKeyPattern = `^[a-zA-Z0-9\-_\.]{1,40}$`
	displayNameMaxLength = 80
)

// NewService returns a new device Service
func NewService() (northbound.Service, error) {
	deviceStore, err := devicestore.NewAtomixStore()
	if err != nil {
		return nil, err
	}
	return &Service{
		store: deviceStore,
	}, nil
}

// Service is a Service implementation for administration.
type Service struct {
	northbound.Service
	store devicestore.Store
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{
		deviceStore: s.store,
	}
	deviceapi.RegisterDeviceServiceServer(r, server)
}

// Server implements the gRPC service for administrative facilities.
type Server struct {
	deviceStore devicestore.Store
}

// DeviceServiceClientFactory : Default DeviceServiceClient creation.
var DeviceServiceClientFactory = func(cc *grpc.ClientConn) deviceapi.DeviceServiceClient {
	return deviceapi.NewDeviceServiceClient(cc)
}

// CreateDeviceServiceClient creates and returns a new topo device client
func CreateDeviceServiceClient(cc *grpc.ClientConn) deviceapi.DeviceServiceClient {
	return DeviceServiceClientFactory(cc)
}

// ValidateDevice validates the given device
func ValidateDevice(device *deviceapi.Device) error {
	nameRegex := regexp.MustCompile(deviceNamePattern)
	if device.ID == "" {
		return status.Error(codes.InvalidArgument, "device ID is required")
	}
	if !nameRegex.MatchString(string(device.ID)) {
		return status.Errorf(codes.InvalidArgument, "device ID '%s' is invalid", device.ID)
	}

	addressRegex := regexp.MustCompile(deviceAddressPattern)
	if device.Address == "" {
		return status.Error(codes.InvalidArgument, "device address is required")
	}
	if !addressRegex.MatchString(device.Address) {
		return status.Errorf(codes.InvalidArgument, "device address '%s' is invalid", device.Address)
	}

	if device.Type == "" {
		return status.Error(codes.InvalidArgument, "device type is required")
	}
	if !nameRegex.MatchString(string(device.Type)) {
		return status.Errorf(codes.InvalidArgument, "device type '%s' is invalid", device.ID)
	}

	versionRegex := regexp.MustCompile(deviceVersionPattern)
	if device.Version == "" {
		return status.Error(codes.InvalidArgument, "device version is required")
	}
	if !versionRegex.MatchString(device.Version) {
		return status.Errorf(codes.InvalidArgument, "device version '%s' is invalid", device.Version)
	}

	if len(device.Displayname) > displayNameMaxLength {
		return status.Errorf(codes.InvalidArgument,
			"device displayname '%s' is too long. (>%d)", device.Displayname, displayNameMaxLength)
	}

	attrKeyRegex := regexp.MustCompile(deviceAttrKeyPattern)
	if device.Attributes != nil {
		for key := range device.Attributes {
			if !attrKeyRegex.MatchString(key) {
				return status.Errorf(codes.InvalidArgument, "attribute name '%s' is invalid", key)
			}
		}
	}

	if device.Timeout == nil {
		timeout := defaultTimeout
		device.Timeout = &timeout
	}
	return nil
}

// Add :
func (s *Server) Add(ctx context.Context, request *deviceapi.AddRequest) (*deviceapi.AddResponse, error) {
	device := request.Device
	if device == nil {
		return nil, status.Error(codes.InvalidArgument, "no device specified")
	} else if device.Revision > 0 {
		return nil, status.Error(codes.InvalidArgument, "device revision is already set")
	} else if err := ValidateDevice(device); err != nil {
		return nil, err
	}
	if err := s.deviceStore.Store(device); err != nil {
		return nil, err
	}
	return &deviceapi.AddResponse{
		Device: device,
	}, nil
}

// Update :
func (s *Server) Update(ctx context.Context, request *deviceapi.UpdateRequest) (*deviceapi.UpdateResponse, error) {
	device := request.Device
	if device == nil {
		return nil, status.Error(codes.InvalidArgument, "no device specified")
	} else if device.Revision == 0 {
		return nil, status.Error(codes.InvalidArgument, "device revision not set")
	} else if err := ValidateDevice(device); err != nil {
		return nil, err
	}
	if err := s.deviceStore.Store(device); err != nil {
		return nil, err
	}
	log.Infof("Updated Device %v", device)
	return &deviceapi.UpdateResponse{
		Device: device,
	}, nil
}

// Get :
func (s *Server) Get(ctx context.Context, request *deviceapi.GetRequest) (*deviceapi.GetResponse, error) {
	device, err := s.deviceStore.Load(request.ID)
	if err != nil {
		return nil, err
	} else if device == nil {
		return nil, status.Error(codes.NotFound, "device not found")
	}
	return &deviceapi.GetResponse{
		Device: device,
	}, nil
}

// List :
func (s *Server) List(request *deviceapi.ListRequest, server deviceapi.DeviceService_ListServer) error {
	if request.Subscribe {
		ch := make(chan *devicestore.Event)
		if err := s.deviceStore.Watch(ch); err != nil {
			return err
		}

		for event := range ch {
			var t deviceapi.ListResponse_Type
			switch event.Type {
			case devicestore.EventNone:
				t = deviceapi.ListResponse_NONE
			case devicestore.EventInserted:
				t = deviceapi.ListResponse_ADDED
			case devicestore.EventUpdated:
				t = deviceapi.ListResponse_UPDATED
			case devicestore.EventRemoved:
				t = deviceapi.ListResponse_REMOVED
			}
			err := server.Send(&deviceapi.ListResponse{
				Type:   t,
				Device: event.Device,
			})
			if err != nil {
				return err
			}
		}
	} else {
		ch := make(chan *deviceapi.Device)
		if err := s.deviceStore.List(ch); err != nil {
			return err
		}

		for device := range ch {
			err := server.Send(&deviceapi.ListResponse{
				Type:   deviceapi.ListResponse_NONE,
				Device: device,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Remove :
func (s *Server) Remove(ctx context.Context, request *deviceapi.RemoveRequest) (*deviceapi.RemoveResponse, error) {
	device := request.Device
	err := s.deviceStore.Delete(device)
	if err != nil {
		return nil, err
	}
	return &deviceapi.RemoveResponse{}, nil
}
