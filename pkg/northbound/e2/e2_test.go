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
	"github.com/onosproject/ran-simulator/pkg/manager"
)

func setUpManager() (*manager.Manager, error) {
	mapLayout := types.MapLayout{
		Center:     &types.Point{Lat: 52.0, Lng: 8.0},
		Zoom:       12,
		Fade:       false,
		ShowRoutes: false,
		ShowPower:  false,
		MinUes:     5,
		MaxUes:     5,
	}
	towerParams := types.TowersParams{
		TowerRows:         2,
		TowerCols:         2,
		TowerSpacingVert:  0.01,
		TowerSpacingHoriz: 0.01,
		MaxUEsPerCell:     4,
		LocationsScale:    1.0,
		AvgCellsPerTower:  1.0,
	}
	routesParams := manager.RoutesParams{
		APIKey:    "",
		StepDelay: 1000,
	}

	towers := manager.NewCells(towerParams, mapLayout)

	locations := manager.NewLocations(towerParams, mapLayout)

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
	mgr.UserEquipments = mgr.NewUserEquipments(mapLayout, routesParams)

	return mgr, nil
}
