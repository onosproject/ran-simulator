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

package manager

import (
	"fmt"
	"math"

	"github.com/onosproject/ran-simulator/api/types"
)

// TowersParams :
type TowersParams struct {
	TowerRows         int
	TowerCols         int
	TowerSpacingVert  float32
	TowerSpacingHoriz float32
}

// TowerIf :
type TowerIf interface {
	GetPosition() types.Point
}

func newTowers(params TowersParams, mapLayout types.MapLayout) map[string]*types.Tower {
	topLeft := types.Point{
		Lat: mapLayout.GetCenter().GetLat() + params.TowerSpacingVert*float32(params.TowerRows)/2,
		Lng: mapLayout.GetCenter().GetLng() - params.TowerSpacingHoriz*float32(params.TowerCols)/2,
	}
	var towerNum = 0
	towers := make(map[string]*types.Tower)

	for r := 0; r < params.TowerRows; r++ {
		for c := 0; c < params.TowerCols; c++ {
			pos := types.Point{
				Lat: topLeft.Lat - 0.03*float32(r),
				Lng: topLeft.Lng + 0.05*float32(c),
			}
			towerNum = towerNum + 1
			towerName := fmt.Sprintf("Tower-%d", towerNum)
			towers[towerName] = &types.Tower{
				Name:     towerName,
				Location: &pos,
				Color:    randomColor(),
			}
		}
	}

	return towers
}

// Find the closest tower to any point - return serving, candidate1 and candidate2
// in order of distance
func (m *Manager) findClosestTower(point *types.Point) (string, string, string) {
	var serving string
	var candidate1 string
	var candidate2 string

	var servingDist float32 = math.MaxFloat32
	var candidate1Dist float32 = math.MaxFloat32
	var candidate2Dist float32 = math.MaxFloat32

	for _, tower := range m.Towers {
		distance := distanceToTower(tower, point)
		if distance < servingDist {
			candidate2 = candidate1
			candidate2Dist = candidate1Dist
			candidate1 = serving
			candidate1Dist = servingDist
			serving = tower.Name
			servingDist = distance
		} else if distance < candidate1Dist {
			candidate2 = candidate1
			candidate2Dist = candidate1Dist
			candidate1 = tower.Name
			candidate1Dist = distance
		} else if distance < candidate2Dist {
			candidate2 = tower.Name
			candidate2Dist = distance
		}
	}

	return serving, candidate1, candidate2
}

// Measure the distance between a point and a tower and return an answer in decimal degrees
// Simple arithmetic is used, do not use for >= 180 degrees
func distanceToTower(tower *types.Tower, point *types.Point) float32 {
	return float32(math.Hypot(
		float64(tower.GetLocation().GetLat()-point.GetLat()),
		float64(tower.GetLocation().GetLng()-point.GetLng()),
	))
}

// GetTowerByName ...
func (m *Manager) GetTowerByName(name string) *types.Tower {
	for _, tower := range m.Towers {
		if tower.Name == name {
			return tower
		}
	}
	return nil
}
