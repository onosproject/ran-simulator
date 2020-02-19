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
	"github.com/onosproject/ran-simulator/api/types"
	"gotest.tools/assert"
	"testing"
)

func Test_findClosestTowers(t *testing.T) {
	m, err := NewManager()
	assert.NilError(t, err, "Unexpected error creating manager")
	const mapCenterLat = 52.0
	const mapCenterLng = -8.0
	const towerSpacingVert = 0.01
	const towerSpacingHoriz = 0.02
	const decimalDegreeTolerance = 0.0001

	m.Towers = NewTowers(
		types.TowersParams{
			TowerRows:         3,
			TowerCols:         3,
			TowerSpacingVert:  towerSpacingVert,
			TowerSpacingHoriz: towerSpacingHoriz,
		},
		types.MapLayout{
			Center:     &types.Point{Lat: mapCenterLat, Lng: mapCenterLng},
			Zoom:       12,
			Fade:       false,
			ShowRoutes: false,
		})

	assert.Equal(t, 9, len(m.Towers), "Expected 9 towersA to have been created")
	for _, tower := range m.Towers {
		switch tower.Name {
		case "Tower-1":
		case "Tower-2":
		case "Tower-3":
			assert.Assert(t, tower.Location.GetLat()-mapCenterLat-towerSpacingVert < decimalDegreeTolerance)
		case "Tower-4":
		case "Tower-5":
		case "Tower-6":
			assert.Assert(t, tower.Location.GetLat()-mapCenterLat < decimalDegreeTolerance)
		case "Tower-7":
		case "Tower-8":
		case "Tower-9":
			assert.Assert(t, tower.Location.GetLat()-mapCenterLat+towerSpacingVert < decimalDegreeTolerance)
		default:
			t.Errorf("Unexpected tower %s", tower.Name)
		}
		switch tower.Name {
		case "Tower-1":
		case "Tower-4":
		case "Tower-7":
			assert.Assert(t, tower.Location.GetLng()+mapCenterLng+towerSpacingHoriz < decimalDegreeTolerance)
		case "Tower-2":
		case "Tower-5":
		case "Tower-8":
			assert.Assert(t, tower.Location.GetLng()-mapCenterLng < decimalDegreeTolerance)
		case "Tower-3":
		case "Tower-6":
		case "Tower-9":
			assert.Assert(t, tower.Location.GetLng()-mapCenterLng-towerSpacingHoriz < decimalDegreeTolerance)
		default:
			t.Errorf("Unexpected tower %s", tower.Name)
		}
	}

	// Test a point outside the towers north-west
	testPointA := &types.Point{Lat: 52.12345, Lng: -8.123}
	towersA, distancesA := m.findClosestTowers(testPointA)
	assert.Equal(t, 3, len(towersA), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesA), "Expected 3 tower distancesA in findClosest")
	assert.Assert(t, distancesA[2] > distancesA[1], "Expected distance to be greater")
	assert.Assert(t, distancesA[1] > distancesA[0], "Expected distance to be greater")

	// Test a point outside the towers south-east
	testPointB := &types.Point{Lat: 51.7654, Lng: -7.9876}
	towersB, distancesB := m.findClosestTowers(testPointB)
	assert.Equal(t, 3, len(towersB), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesB), "Expected 3 tower distancesA in findClosest")
	assert.Assert(t, distancesB[2] > distancesB[1], "Expected distance to be greater")
	assert.Assert(t, distancesB[1] > distancesB[0], "Expected distance to be greater")

	// Test a point within the towers south-east of centre
	testPointC := &types.Point{Lat: 51.980, Lng: -7.950}
	towersC, distancesC := m.findClosestTowers(testPointC)
	assert.Equal(t, 3, len(towersC), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesC), "Expected 3 tower distancesA in findClosest")
	assert.Assert(t, distancesC[2] > distancesC[1], "Expected distance to be greater")
	assert.Assert(t, distancesC[1] > distancesC[0], "Expected distance to be greater")
}
