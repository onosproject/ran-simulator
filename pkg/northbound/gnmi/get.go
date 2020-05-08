// Copyright 2020-present Open Networking Foundation.
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

package gnmi

import (
	"context"
	"fmt"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get implements gNMI Get
func (s *Server) Get(ctx context.Context, req *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	cell := manager.GetManager().GetCell(s.GetECGI())
	if cell == nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("cell %s not found", s.GetECGI()))
	}

	notifications := make([]*gnmi.Notification, 0)

	if req.Path == nil || len(req.Path) == 0 {
		// return everything
		for key, value := range cell.GetConfigAttributes() {
			path, err := utils.ParseGNMIElements(utils.SplitPath(string(key)))
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("path error on %s %v", s.GetECGI(), err))
			}
			typedValue, err := bytesToGnmiValue(value)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("value error on %s %v", s.GetECGI(), err))
			}
			notification := gnmi.Notification{
				Timestamp: 0,
				Update: []*gnmi.Update{
					{Path: path, Val: typedValue},
				},
				Delete: nil,
			}
			notifications = append(notifications, &notification)
		}
	} else {
		// TODO: Implement the rest of Get
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("gNMI Get not yet supported on Port %d", s.port))
	}

	response := gnmi.GetResponse{
		Notification: notifications,
	}

	return &response, nil
}

func gnmiValueToBytes(val *gnmi.TypedValue) ([]byte, error) {
	test := make([]byte, 0)
	valBytes, err := val.XXX_Marshal(test, true)
	return valBytes, err
}

func bytesToGnmiValue(val []byte) (*gnmi.TypedValue, error) {
	typedVal := gnmi.TypedValue{}
	if err := typedVal.XXX_Unmarshal(val); err != nil {
		return nil, err
	}
	return &typedVal, nil
}
