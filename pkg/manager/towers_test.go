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
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"gotest.tools/assert"
	"testing"
	"time"
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

	assert.Equal(t, 9, len(m.Towers), "Expected 9 towers to have been created")
	for _, tower := range m.Towers {
		switch tower.EcID {
		case "0001420":
		case "0001421":
		case "0001422":
			assert.Assert(t, tower.Location.GetLat()-mapCenterLat-towerSpacingVert < decimalDegreeTolerance)
		case "0001423":
		case "0001424":
		case "0001425":
			assert.Assert(t, tower.Location.GetLat()-mapCenterLat < decimalDegreeTolerance)
		case "0001426":
		case "0001427":
		case "0001428":
			assert.Assert(t, tower.Location.GetLat()-mapCenterLat+towerSpacingVert < decimalDegreeTolerance)
		default:
			t.Errorf("Unexpected tower %s", tower.EcID)
		}
		switch tower.EcID {
		case "0001420":
		case "0001423":
		case "0001426":
			assert.Assert(t, tower.Location.GetLng()+mapCenterLng+towerSpacingHoriz < decimalDegreeTolerance)
		case "0001421":
		case "0001424":
		case "0001427":
			assert.Assert(t, tower.Location.GetLng()-mapCenterLng < decimalDegreeTolerance)
		case "0001422":
		case "0001425":
		case "0001428":
			assert.Assert(t, tower.Location.GetLng()-mapCenterLng-towerSpacingHoriz < decimalDegreeTolerance)
		default:
			t.Errorf("Unexpected tower %s", tower.EcID)
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

func Test_PowerAdjust(t *testing.T) {
	m, err := NewManager()
	assert.NilError(t, err, "Unexpected error creating manager")

	const mapCenterLat = 52.0
	const mapCenterLng = -8.0
	const towerSpacingVert = 0.01
	const towerSpacingHoriz = 0.02
	m.Towers = NewTowers(
		types.TowersParams{
			TowerRows:         1,
			TowerCols:         1,
			TowerSpacingVert:  towerSpacingVert,
			TowerSpacingHoriz: towerSpacingHoriz,
		},
		types.MapLayout{
			Center:     &types.Point{Lat: mapCenterLat, Lng: mapCenterLng},
			Zoom:       12,
			Fade:       false,
			ShowRoutes: false,
		})
	go func() {
		for event := range m.TowerChannel {
			assert.Equal(t, trafficsim.Type_UPDATED, event.Type)
		}
	}()

	assert.Equal(t, 1, len(m.Towers), "Expected 1 tower to have been created")

	err = m.UpdateTower("0001420", -6) // subtracted from initial 10dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	tower1, ok := m.Towers["0001420"]
	assert.Assert(t, ok)
	assert.Equal(t, float32(4.0), tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with value too low - capped at -15dB /////////////////////
	err = m.UpdateTower("0001420", -30) // subtracted from prev 4dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	assert.Equal(t, float32(-15.0), tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with value too high - capped at 30dB /////////////////////
	err = m.UpdateTower("0001420", 50) // Added to prev -15dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	assert.Equal(t, float32(30.0), tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with wrong name /////////////////////
	err = m.UpdateTower("0001421", -3)
	assert.Error(t, err, "unknown tower 0001421", "Expected an error for wrong name when adjusting power")

	time.Sleep(time.Millisecond * 100)
}
