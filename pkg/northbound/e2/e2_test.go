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

package e2

import (
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/config"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

func setUpManager() (*manager.Manager, error) {
	routesParams := manager.RoutesParams{
		APIKey:    "",
		StepDelay: 1000,
	}

	towersConfig, err := config.GetTowerConfig("berlin-rectangular-4-1.yaml")
	if err != nil {
		return nil, err
	}
	mapLayout := types.MapLayout{
		Center: &towersConfig.MapCentre,
		Zoom:   12,
		MinUes: 3,
	}

	towers := manager.NewCells(towersConfig)

	locations := manager.NewLocations(towersConfig, 5, 1)

	mgr, err := manager.NewManager()
	if err != nil {
		return nil, err
	}
	mgr.MapLayout = mapLayout
	mgr.CellsLock.Lock()
	mgr.Cells = towers
	mgr.CellsLock.Unlock()
	mgr.Locations = locations

	mgr.Routes, err = mgr.NewRoutes(mapLayout, routesParams)
	if err != nil {
		return nil, err
	}
	mgr.UserEquipments, err = mgr.NewUserEquipments(mapLayout, routesParams)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}
