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
	"github.com/onosproject/config-models/modelplugin/e2node-1.0.0/e2node_1_0_0"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get implements gNMI Get
func (s *Server) Get(ctx context.Context, req *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	cellConfig, ok := manager.GetManager().CellConfigs[s.GetECGI()]
	if !ok {
		return nil, status.Error(codes.Internal, fmt.Sprintf("cell %s not found", s.GetECGI()))
	}

	notifications := make([]*gnmi.Notification, 0)

	if req.Path == nil || len(req.Path) == 0 {
		// return everything that has a value
		if notif, err := getE2nodeIntervalsPdcpMeasReportPerUe(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
		if notif, err := getE2nodeIntervalsRadioMeasReportPerCell(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
		if notif, err := getE2nodeIntervalsRadioMeasReportPerUe(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
		if notif, err := getE2nodeIntervalsSchedMeasReportPerCell(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
		if notif, err := getE2nodeIntervalsSchedMeasReportPerUe(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
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

func getE2nodeIntervalsPdcpMeasReportPerUe(ecgi types.ECGI, cellConfig *e2node_1_0_0.Device) (*gnmi.Notification, error) {
	if cellConfig.E2Node.Intervals.PdcpMeasReportPerUe != nil {
		path, err := utils.ParseGNMIElements(utils.SplitPath(e2nodeIntervalsPdcpMeasReportPerUe))
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("path error on %s %v", ecgi, err))
		}
		val := valFromUint(cellConfig.E2Node.Intervals.PdcpMeasReportPerUe)
		return newNotif(path, val), nil
	}
	return nil, nil
}

func getE2nodeIntervalsRadioMeasReportPerCell(ecgi types.ECGI, cellConfig *e2node_1_0_0.Device) (*gnmi.Notification, error) {
	if cellConfig.E2Node.Intervals.RadioMeasReportPerCell != nil {
		path, err := utils.ParseGNMIElements(utils.SplitPath(e2nodeIntervalsRadioMeasReportPerCell))
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("path error on %s %v", ecgi, err))
		}
		val := valFromUint(cellConfig.E2Node.Intervals.RadioMeasReportPerCell)
		return newNotif(path, val), nil
	}
	return nil, nil
}

func getE2nodeIntervalsRadioMeasReportPerUe(ecgi types.ECGI, cellConfig *e2node_1_0_0.Device) (*gnmi.Notification, error) {
	if cellConfig.E2Node.Intervals.RadioMeasReportPerUe != nil {
		path, err := utils.ParseGNMIElements(utils.SplitPath(e2nodeIntervalsRadioMeasReportPerUe))
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("path error on %s %v", ecgi, err))
		}
		val := valFromUint(cellConfig.E2Node.Intervals.RadioMeasReportPerUe)
		return newNotif(path, val), nil
	}
	return nil, nil
}

func getE2nodeIntervalsSchedMeasReportPerCell(ecgi types.ECGI, cellConfig *e2node_1_0_0.Device) (*gnmi.Notification, error) {
	if cellConfig.E2Node.Intervals.SchedMeasReportPerCell != nil {
		path, err := utils.ParseGNMIElements(utils.SplitPath(e2nodeIntervalsSchedMeasReportPerCell))
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("path error on %s %v", ecgi, err))
		}
		val := valFromUint(cellConfig.E2Node.Intervals.SchedMeasReportPerCell)
		return newNotif(path, val), nil
	}
	return nil, nil
}

func getE2nodeIntervalsSchedMeasReportPerUe(ecgi types.ECGI, cellConfig *e2node_1_0_0.Device) (*gnmi.Notification, error) {
	if cellConfig.E2Node.Intervals.SchedMeasReportPerUe != nil {
		path, err := utils.ParseGNMIElements(utils.SplitPath(e2nodeIntervalsSchedMeasReportPerUe))
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("path error on %s %v", ecgi, err))
		}
		val := valFromUint(cellConfig.E2Node.Intervals.SchedMeasReportPerUe)
		return newNotif(path, val), nil
	}
	return nil, nil
}

func newNotif(path *gnmi.Path, val *gnmi.TypedValue) *gnmi.Notification {
	return &gnmi.Notification{
		Timestamp: 0,
		Update: []*gnmi.Update{
			{Path: path, Val: val},
		}}
}

func valFromUint(val *uint32) *gnmi.TypedValue {
	return &gnmi.TypedValue{
		Value: &gnmi.TypedValue_UintVal{
			UintVal: uint64(*val),
		}}
}
