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

	m.Towers = newTowers(
		types.TowersParams{
			TowerRows:         3,
			TowerCols:         3,
			TowerSpacingVert:  0.01,
			TowerSpacingHoriz: 0.02,
		},
		types.MapLayout{
			Center:     &types.Point{Lat: 52.0, Lng: -8.0},
			Zoom:       12,
			Fade:       false,
			ShowRoutes: false,
		})

	assert.Equal(t, 9, len(m.Towers), "Expected 9 towersA to have been created")
	for _, tower := range m.Towers {
		switch tower.Name {
		case "Tower-1":
			assert.Equal(t, float32(52.014999), tower.Location.GetLat())
			assert.Equal(t, float32(-8.030000), tower.Location.GetLng())
		case "Tower-2":
			assert.Equal(t, float32(52.014999), tower.Location.GetLat())
			assert.Equal(t, float32(-7.9799995), tower.Location.GetLng())
		case "Tower-3":
			assert.Equal(t, float32(52.014999), tower.Location.GetLat())
			assert.Equal(t, float32(-7.930000), tower.Location.GetLng())
		case "Tower-4":
			assert.Equal(t, float32(51.985001), tower.Location.GetLat())
			assert.Equal(t, float32(-8.030000), tower.Location.GetLng())
		case "Tower-5":
			assert.Equal(t, float32(51.985001), tower.Location.GetLat())
			assert.Equal(t, float32(-7.9799995), tower.Location.GetLng())
		case "Tower-6":
			assert.Equal(t, float32(51.985), tower.Location.GetLat())
			assert.Equal(t, float32(-7.93), tower.Location.GetLng())
		case "Tower-7":
			assert.Equal(t, float32(51.954998), tower.Location.GetLat())
			assert.Equal(t, float32(-8.030000), tower.Location.GetLng())
		case "Tower-8":
			assert.Equal(t, float32(51.954998), tower.Location.GetLat())
			assert.Equal(t, float32(-7.9799995), tower.Location.GetLng())
		case "Tower-9":
			assert.Equal(t, float32(51.954998), tower.Location.GetLat())
			assert.Equal(t, float32(-7.93), tower.Location.GetLng())
		default:
			t.Errorf("Unexpected tower %s", tower.Name)
		}
	}

	// Test a point outside the towers north-west
	testPointA := &types.Point{Lat: 52.12345, Lng: -8.123}
	towersA, distancesA := m.findClosestTowers(testPointA)
	assert.Equal(t, 3, len(towersA), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesA), "Expected 3 tower distancesA in findClosest")

	assert.Equal(t, "Tower-1", towersA[0], "Unexpected name for closest tower")
	assert.Equal(t, float32(0.14286664), distancesA[0], "Unexpected dist for closest tower")

	assert.Equal(t, "Tower-4", towersA[1], "Unexpected name for 2nd closest tower")
	assert.Equal(t, float32(0.16678624), distancesA[1], "Unexpected dist for 2nd closest tower")

	assert.Equal(t, "Tower-2", towersA[2], "Unexpected name for 3rd closest tower")
	assert.Equal(t, float32(0.17947416), distancesA[2], "Unexpected dist for 3rd closest tower")

	assert.Assert(t, distancesA[2] > distancesA[1], "Expected distance to be greater")
	assert.Assert(t, distancesA[1] > distancesA[0], "Expected distance to be greater")


	// Test a point outside the towers south-east
	testPointB := &types.Point{Lat: 51.7654, Lng: -7.9876}
	towersB, distancesB := m.findClosestTowers(testPointB)
	assert.Equal(t, 3, len(towersB), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesB), "Expected 3 tower distancesA in findClosest")

	assert.Equal(t, "Tower-8", towersB[0], "Unexpected name for closest tower")
	assert.Equal(t, float32(0.18975036), distancesB[0], "Unexpected dist for closest tower")

	assert.Equal(t, "Tower-7", towersB[1], "Unexpected name for 2nd closest tower")
	assert.Equal(t, float32(0.19428119), distancesB[1], "Unexpected dist for 2nd closest tower")

	assert.Equal(t, "Tower-9", towersB[2], "Unexpected name for 3rd closest tower")
	assert.Equal(t, float32(0.19815448), distancesB[2], "Unexpected dist for 3rd closest tower")

	assert.Assert(t, distancesB[2] > distancesB[1], "Expected distance to be greater")
	assert.Assert(t, distancesB[1] > distancesB[0], "Expected distance to be greater")

	// Test a point within the towers south-east of centre
	testPointC := &types.Point{Lat: 51.980, Lng: -7.950}
	towersC, distancesC := m.findClosestTowers(testPointC)
	assert.Equal(t, 3, len(towersC), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesC), "Expected 3 tower distancesA in findClosest")

	assert.Equal(t, "Tower-6", towersC[0], "Unexpected name for closest tower")
	assert.Equal(t, float32(0.02061577), distancesC[0], "Unexpected dist for closest tower")

	assert.Equal(t, "Tower-5", towersC[1], "Unexpected name for 2nd closest tower")
	assert.Equal(t, float32(0.030413724), distancesC[1], "Unexpected dist for 2nd closest tower")

	assert.Equal(t, "Tower-9", towersC[2], "Unexpected name for 3rd closest tower")
	assert.Equal(t, float32(0.032016803), distancesC[2], "Unexpected dist for 3rd closest tower")

	assert.Assert(t, distancesC[2] > distancesC[1], "Expected distance to be greater")
	assert.Assert(t, distancesC[1] > distancesC[0], "Expected distance to be greater")
}
