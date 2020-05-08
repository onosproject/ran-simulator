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

package gnmi

import (
	"context"
	"fmt"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// Set implements gNMI Set
func (s *Server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	log.Info("gNMI Set Request", req)

	// Build the responses
	updateResults := make([]*gnmi.UpdateResult, 0)

	cell := manager.GetManager().GetCell(s.GetECGI())
	if cell == nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("cell %s not found", s.GetECGI()))
	}

	// Update - target is ignored - we already know it - the cell ecgi
	for _, update := range req.GetUpdate() {
		key, value, err := s.formatUpdateOrReplace(req.GetPrefix(), update)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error handling delete value %v", update))
		}
		cell.GetConfigAttributes()[key] = value
		result := gnmi.UpdateResult{Path: update.GetPath(), Op: gnmi.UpdateResult_UPDATE}
		updateResults = append(updateResults, &result)
	}

	// Replace - target is ignored - we already know it - the cell ecgi
	for _, replace := range req.GetReplace() {
		key, value, err := s.formatUpdateOrReplace(req.GetPrefix(), replace)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error handling replace value %v", replace))
		}
		cell.GetConfigAttributes()[key] = value
		result := gnmi.UpdateResult{Path: replace.GetPath(), Op: gnmi.UpdateResult_REPLACE}
		updateResults = append(updateResults, &result)
	}

	// Delete - target is ignored - we already know it - the cell ecgi
	for _, deletePath := range req.GetDelete() {
		key, err := s.formatDelete(req.GetPrefix(), deletePath)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error handling delete value %v", deletePath))
		}
		delete(cell.GetConfigAttributes(), key)
		result := gnmi.UpdateResult{Path: deletePath, Op: gnmi.UpdateResult_DELETE}
		updateResults = append(updateResults, &result)
	}

	setResponse := &gnmi.SetResponse{
		Response:  updateResults,
		Timestamp: time.Now().Unix(),
	}

	return setResponse, nil
}

// This deals with a path and a value
func (s *Server) formatUpdateOrReplace(prefix *gnmi.Path, u *gnmi.Update) (types.ConfigKey, types.ConfigValue, error) {
	prefixPath := utils.StrPath(prefix)
	updatePath := utils.StrPath(u.Path)
	valBytes, err := gnmiValueToBytes(u.Val)
	if err != nil {
		return "", nil, nil
	}
	return types.ConfigKey(prefixPath + updatePath), valBytes, nil
}

// This deals with a path
func (s *Server) formatDelete(prefix *gnmi.Path, d *gnmi.Path) (types.ConfigKey, error) {
	prefixPath := utils.StrPath(prefix)
	deletePath := utils.StrPath(d)
	return types.ConfigKey(prefixPath + deletePath), nil
}
