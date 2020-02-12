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
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"math"

	"github.com/onosproject/ran-simulator/api/types"
)

// TowerIf :
type TowerIf interface {
	GetPosition() types.Point
}

func newTowers(params types.TowersParams, mapLayout types.MapLayout) map[string]*types.Tower {
	topLeft := types.Point{
		Lat: mapLayout.GetCenter().GetLat() + params.TowerSpacingVert*float32(params.TowerRows-1)/2,
		Lng: mapLayout.GetCenter().GetLng() - params.TowerSpacingHoriz*float32(params.TowerCols-1)/2,
	}
	var towerNum = 0
	towers := make(map[string]*types.Tower)

	for r := 0; r < int(params.TowerRows); r++ {
		for c := 0; c < int(params.TowerCols); c++ {
			pos := types.Point{
				Lat: topLeft.Lat - params.TowerSpacingVert*float32(r),
				Lng: topLeft.Lng + params.TowerSpacingHoriz*float32(c),
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
func (m *Manager) findClosestTowers(point *types.Point) ([]string, []float32) {
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

	return []string{serving, candidate1, candidate2}, []float32{servingDist, candidate1Dist, candidate2Dist}
}

// GetTower returns tower based on its name
func (m *Manager) GetTower(name string) *types.Tower {
	return m.Towers[name]
}

// UpdateTower Update a tower's properties - usually power level
func (m *Manager) UpdateTower(tower *types.Tower) {
	// Only the power can be updated at present
	m.Towers[tower.GetName()].TxPower = tower.TxPower
	m.TowerChannel <- dispatcher.Event{
		Type:   trafficsim.Type_UPDATED,
		Object: tower,
	}
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
