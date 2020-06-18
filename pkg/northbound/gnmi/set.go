// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package gnmi

import (
	"context"
	"fmt"
	"github.com/onosproject/config-models/modelplugin/e2node-1.0.0/e2node_1_0_0"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// Set implements gNMI Set
func (s *Server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	log.Infof("Cell %s. gNMI Set Request", s.GetECGI(), req)

	// Build the responses
	updateResults := make([]*gnmi.UpdateResult, 0)
	prefixPath := utils.StrPath(req.GetPrefix())
	if prefixPath == "/" {
		prefixPath = ""
	}

	cellConfig, ok := manager.GetManager().CellConfigs[s.GetECGI()]
	if !ok {
		return nil, status.Error(codes.Internal, fmt.Sprintf("cell %s not found", s.GetECGI()))
	}

	// Update - target is ignored - we already know it - the cell ecgi
	for _, update := range req.GetUpdate() {
		updatePath := utils.StrPath(update.Path)
		uintVal := uint32(update.GetVal().GetUintVal())
		if err := updateOrReplace(prefixPath+updatePath, cellConfig, &uintVal); err != nil {
			return nil, err
		}

		result := gnmi.UpdateResult{Path: update.GetPath(), Op: gnmi.UpdateResult_UPDATE}
		updateResults = append(updateResults, &result)
	}

	// Replace - target is ignored - we already know it - the cell ecgi
	for _, replace := range req.GetReplace() {
		replacePath := utils.StrPath(replace.Path)
		uintVal := uint32(replace.GetVal().GetUintVal())
		if err := updateOrReplace(prefixPath+replacePath, cellConfig, &uintVal); err != nil {
			return nil, err
		}

		result := gnmi.UpdateResult{Path: replace.GetPath(), Op: gnmi.UpdateResult_REPLACE}
		updateResults = append(updateResults, &result)
	}

	// Delete - target is ignored - we already know it - the cell ecgi
	for _, deleteGnmiPath := range req.GetDelete() {
		deletePath := utils.StrPath(deleteGnmiPath)
		if err := delete(prefixPath+deletePath, cellConfig); err != nil {
			return nil, err
		}

		result := gnmi.UpdateResult{Path: deleteGnmiPath, Op: gnmi.UpdateResult_DELETE}
		updateResults = append(updateResults, &result)
	}

	setResponse := &gnmi.SetResponse{
		Prefix:    req.GetPrefix(),
		Response:  updateResults,
		Timestamp: time.Now().Unix(),
	}

	return setResponse, nil
}

func updateOrReplace(path string, cellConfig *e2node_1_0_0.Device, value interface{}) error {
	log.Infof("Set %s to %v", path, value)
	switch path {
	case e2nodeIntervalsPdcpMeasReportPerUe:
		valUint32, ok := value.(*uint32)
		if !ok {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid update value %s %v", path, value))
		}
		cellConfig.E2Node.Intervals.PdcpMeasReportPerUe = valUint32
	case e2nodeIntervalsRadioMeasReportPerCell:
		valUint32, ok := value.(*uint32)
		if !ok {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid update value %s %v", path, value))
		}
		cellConfig.E2Node.Intervals.RadioMeasReportPerCell = valUint32
	case e2nodeIntervalsRadioMeasReportPerUe:
		valUint32, ok := value.(*uint32)
		if !ok {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid update value %s %v", path, value))
		}
		cellConfig.E2Node.Intervals.RadioMeasReportPerUe = valUint32
	case e2nodeIntervalsSchedMeasReportPerCell:
		valUint32, ok := value.(*uint32)
		if !ok {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid update value %s %v", path, value))
		}
		cellConfig.E2Node.Intervals.SchedMeasReportPerCell = valUint32
	case e2nodeIntervalsSchedMeasReportPerUe:
		valUint32, ok := value.(*uint32)
		if !ok {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid update value %s %v", path, value))
		}
		cellConfig.E2Node.Intervals.SchedMeasReportPerUe = valUint32
	default:
		return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid update path %s %v", path, value))
	}
	log.Infof("Cell Config %v", cellConfig)

	return nil
}

func delete(path string, cellConfig *e2node_1_0_0.Device) error {
	switch path {
	case e2nodeIntervalsPdcpMeasReportPerUe:
		cellConfig.E2Node.Intervals.PdcpMeasReportPerUe = nil
	case e2nodeIntervalsRadioMeasReportPerCell:
		cellConfig.E2Node.Intervals.RadioMeasReportPerCell = nil
	case e2nodeIntervalsRadioMeasReportPerUe:
		cellConfig.E2Node.Intervals.RadioMeasReportPerUe = nil
	case e2nodeIntervalsSchedMeasReportPerCell:
		cellConfig.E2Node.Intervals.SchedMeasReportPerCell = nil
	case e2nodeIntervalsSchedMeasReportPerUe:
		cellConfig.E2Node.Intervals.SchedMeasReportPerUe = nil
	default:
		return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid delete %s", path))
	}
	return nil
}
