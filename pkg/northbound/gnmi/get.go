// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

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
	prefixPath := utils.StrPath(req.GetPrefix())
	if prefixPath == "/" {
		prefixPath = ""
	}
	allRequested := true
	namedAttrs := make(map[string]interface{})

	if req.Path != nil && len(req.Path) > 0 {
		allRequested = false
		for _, p := range req.GetPath() {
			namedAttrs[prefixPath+utils.StrPath(p)] = struct{}{}
		}
	}

	// return everything that has a value
	if _, ok := namedAttrs[e2nodeIntervalsPdcpMeasReportPerUe]; ok || allRequested {
		if notif, err := getE2nodeIntervalsPdcpMeasReportPerUe(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
	}
	if _, ok := namedAttrs[e2nodeIntervalsRadioMeasReportPerCell]; ok || allRequested {
		if notif, err := getE2nodeIntervalsRadioMeasReportPerCell(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
	}
	if _, ok := namedAttrs[e2nodeIntervalsRadioMeasReportPerUe]; ok || allRequested {
		if notif, err := getE2nodeIntervalsRadioMeasReportPerUe(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
	}
	if _, ok := namedAttrs[e2nodeIntervalsSchedMeasReportPerCell]; ok || allRequested {
		if notif, err := getE2nodeIntervalsSchedMeasReportPerCell(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
	}
	if _, ok := namedAttrs[e2nodeIntervalsSchedMeasReportPerUe]; ok || allRequested {
		if notif, err := getE2nodeIntervalsSchedMeasReportPerUe(s.GetECGI(), cellConfig); err == nil && notif != nil {
			notifications = append(notifications, notif)
		}
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
