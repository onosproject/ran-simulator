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
	"github.com/onosproject/ran-simulator/pkg/utils"
	"gotest.tools/assert"
	"math"
	"testing"
	"time"
)

const (
	mapCenterLat           = 52.0
	mapCenterLng           = -8.0
	towerSpacingVert       = 0.01
	towerSpacingHoriz      = 0.02
	decimalDegreeTolerance = 0.0001
)

func Test_NewTowers(t *testing.T) {
	cells := NewCells(
		types.TowersParams{
			TowerRows:         2,
			TowerCols:         2,
			TowerSpacingVert:  towerSpacingVert,
			TowerSpacingHoriz: towerSpacingHoriz,
			AvgCellsPerTower:  3.5,
		},
		types.MapLayout{
			Center:     &types.Point{Lat: mapCenterLat, Lng: mapCenterLng},
			Zoom:       12,
			Fade:       false,
			ShowRoutes: false,
		})

	assert.Equal(t, 14, len(cells), "Expected 14 cells to have been created")
	for _, cell := range cells {
		testLatitude(cell, t)
		testLongitude(cell, t)
		assert.Assert(t, cell.Sector.Azimuth >= 0 && cell.Sector.Azimuth <= 270, cell.Sector)
		assert.Assert(t, cell.Sector.Arc >= 90 && cell.Sector.Arc <= 120, cell.Sector)
	}
}

// gocyclo complained when they were all in the same function
func testLatitude(cell *types.Cell, t *testing.T) {
	switch cell.Ecgi.EcID {
	case "0001420":
	case "0001421":
	case "0001422":
	case "0001423":
	case "0001424":
	case "0001425":
	case "0001426":
		assert.Assert(t, cell.Location.GetLat()-mapCenterLat-towerSpacingVert/2 < decimalDegreeTolerance,
			"cell %s off V", cell.Ecgi.EcID)
	case "0001427":
	case "0001428":
	case "0001429":
	case "000142A":
	case "000142B":
	case "000142C":
	case "000142D":
		assert.Assert(t, cell.Location.GetLat()-mapCenterLat+towerSpacingVert/2 < decimalDegreeTolerance,
			"cell %s off V", cell.Ecgi.EcID)
	default:
		t.Errorf("Unexpected cell %s V", cell.Ecgi)
	}
}

func testLongitude(cell *types.Cell, t *testing.T) {
	switch cell.Ecgi.EcID {
	case "0001420":
	case "0001421":
	case "0001422":
	case "0001427":
	case "0001428":
	case "0001429":
	case "000142A":
		assert.Assert(t, cell.Location.GetLng()+mapCenterLng-towerSpacingHoriz/2 < decimalDegreeTolerance,
			"cell %s off H", cell.Ecgi.EcID)
	case "0001423":
	case "0001424":
	case "0001425":
	case "0001426":
	case "000142B":
	case "000142C":
	case "000142D":
		assert.Assert(t, cell.Location.GetLng()+mapCenterLng+towerSpacingHoriz/2 < decimalDegreeTolerance,
			"cell %s off H", cell.Ecgi.EcID)
	default:
		t.Errorf("Unexpected cell %s H", cell.Ecgi)
	}
}

func Test_findClosestTowers(t *testing.T) {
	m, err := NewManager()
	assert.NilError(t, err, "Unexpected error creating manager")

	m.Cells = NewCells(
		types.TowersParams{
			TowerRows:         3,
			TowerCols:         3,
			TowerSpacingVert:  towerSpacingVert,
			TowerSpacingHoriz: towerSpacingHoriz,
			AvgCellsPerTower:  1,
		},
		types.MapLayout{
			Center:     &types.Point{Lat: mapCenterLat, Lng: mapCenterLng},
			Zoom:       12,
			Fade:       false,
			ShowRoutes: false,
		})

	assert.Equal(t, 9, len(m.Cells), "Expected 9 towers to have been created")

	// Test a point outside the towers north-west
	testPointA := &types.Point{Lat: 52.12345, Lng: -8.123}
	towersA, distancesA := m.findClosestCells(testPointA)
	assert.Equal(t, 3, len(towersA), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesA), "Expected 3 tower distancesA in findClosest")
	assert.Assert(t, distancesA[2] > distancesA[1], "Expected distance to be greater")
	assert.Assert(t, distancesA[1] > distancesA[0], "Expected distance to be greater")

	// Test a point outside the towers south-east
	testPointB := &types.Point{Lat: 51.7654, Lng: -7.9876}
	towersB, distancesB := m.findClosestCells(testPointB)
	assert.Equal(t, 3, len(towersB), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesB), "Expected 3 tower distancesA in findClosest")
	assert.Assert(t, distancesB[2] > distancesB[1], "Expected distance to be greater")
	assert.Assert(t, distancesB[1] > distancesB[0], "Expected distance to be greater")

	// Test a point within the towers south-east of centre
	testPointC := &types.Point{Lat: 51.980, Lng: -7.950}
	towersC, distancesC := m.findClosestCells(testPointC)
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
	m.Cells = NewCells(
		types.TowersParams{
			TowerRows:         1,
			TowerCols:         1,
			TowerSpacingVert:  towerSpacingVert,
			TowerSpacingHoriz: towerSpacingHoriz,
			AvgCellsPerTower:  1.0,
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

	assert.Equal(t, 1, len(m.Cells), "Expected 1 tower to have been created")
	towerID1420 := newEcgi("0001420", utils.TestPlmnID)
	err = m.UpdateCell(towerID1420, -6) // subtracted from initial 10dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	tower1, ok := m.Cells[towerID1420]
	assert.Assert(t, ok)
	assert.Equal(t, float32(4.0), tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with value too low - capped at -15dB /////////////////////
	err = m.UpdateCell(towerID1420, -30) // subtracted from prev 4dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	assert.Equal(t, float32(-15.0), tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with value too high - capped at 30dB /////////////////////
	err = m.UpdateCell(towerID1420, 50) // Added to prev -15dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	assert.Equal(t, float32(30.0), tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with wrong name /////////////////////
	towerID1421 := newEcgi("0001421", utils.TestPlmnID)
	err = m.UpdateCell(towerID1421, -3)
	assert.Error(t, err, "unknown tower {0001421 315010}", "Expected an error for wrong name when adjusting power")

	time.Sleep(time.Millisecond * 100)
}

func Test_MakeNeighbors(t *testing.T) {
	towerParams := types.TowersParams{
		TowerRows:         3,
		TowerCols:         3,
		TowerSpacingVert:  towerSpacingVert,
		TowerSpacingHoriz: towerSpacingHoriz,
		AvgCellsPerTower:  1.0,
	}

	// 1420 --- 1421 --- 1422
	//   |        |        |
	// 1423 --- 1424 --- 1425
	//   |        |        |
	// 1426 --- 1427 --- 1428
	// tower num 2 is the top right - it's id is "0001422"
	neighborIDs := makeNeighbors(2, towerParams)
	assert.Equal(t, 2, len(neighborIDs), "Unexpected number of neighbors for 2")
	assert.Equal(t, types.EcID("0001421"), neighborIDs[0].EcID)
	assert.Equal(t, types.EcID("0001425"), neighborIDs[1].EcID)

	// tower num 4 is the middle - it's id is "0001424"
	neighborIDs4 := makeNeighbors(4, towerParams)
	assert.Equal(t, 4, len(neighborIDs4), "Unexpected number of neighbors for 5")
	assert.Equal(t, types.EcID("0001421"), neighborIDs4[0].EcID)
	assert.Equal(t, types.EcID("0001423"), neighborIDs4[1].EcID)
	assert.Equal(t, types.EcID("0001425"), neighborIDs4[2].EcID)
	assert.Equal(t, types.EcID("0001427"), neighborIDs4[3].EcID)
}

func Test_distToTower1Sector(t *testing.T) {
	dist := distanceToCellCentroid(
		&types.Cell{
			Location: &types.Point{
				Lat: 52,
				Lng: -8,
			},
			Sector: &types.Sector{
				Azimuth: 0,
				Arc:     360,
			},
			TxPowerdB: 10.0, // Does not matter in this case 360
		},
		&types.Point{
			Lat: 52.01,
			Lng: -8.01,
		})
	assert.Equal(t, 1414, int(math.Floor(float64(dist*1e5))), "Unexpected distance for single sector tower")
}

func Test_distToTower2Sectors(t *testing.T) {
	dist := distanceToCellCentroid(
		&types.Cell{
			Location: &types.Point{
				Lat: 52,
				Lng: -8,
			},
			Sector: &types.Sector{
				Azimuth: 90,
				Arc:     180,
			},
			TxPowerdB: 10.0,
		},
		&types.Point{
			Lat: 52.01,
			Lng: -8.01,
		})
	assert.Equal(t, 1512, int(math.Floor(float64(dist*1e5))), "Unexpected distance for single sector tower")
}

func Test_distToTower3Sectors(t *testing.T) {
	dist := distanceToCellCentroid(
		&types.Cell{
			Location: &types.Point{
				Lat: 52,
				Lng: -8,
			},
			Sector: &types.Sector{
				Azimuth: 120,
				Arc:     120,
			},
			TxPowerdB: 10.0,
		},
		&types.Point{
			Lat: 52.01,
			Lng: -8.01,
		})
	assert.Equal(t, 1542, int(math.Floor(float64(dist*1e5))), "Unexpected distance for single sector tower")
}

func Test_PowerToDist(t *testing.T) {
	distm20 := PowerToDist(-20) // -20dB
	assert.Equal(t, 10, int(math.Floor(distm20*1e5)))

	distm10 := PowerToDist(-10) // -10dB
	assert.Equal(t, 31, int(math.Floor(distm10*1e5)))

	dist0 := PowerToDist(0) // 0dB
	assert.Equal(t, 100, int(math.Floor(dist0*1e5)))

	dist10 := PowerToDist(10)
	assert.Equal(t, 316, int(math.Floor(dist10*1e5)))

	dist20 := PowerToDist(20)
	assert.Equal(t, 1000, int(math.Floor(dist20*1e5)))
}
