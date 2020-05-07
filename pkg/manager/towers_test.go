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
	"github.com/onosproject/ran-simulator/pkg/config"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"gotest.tools/assert"
	"math"
	"testing"
	"time"
)

func Test_NewTowers(t *testing.T) {
	config.Clear()
	towerConfig, err := config.GetTowerConfig("berlin-honeycomb-4-3.yaml")
	assert.NilError(t, err)
	cells := NewCells(towerConfig)

	assert.Equal(t, 12, len(cells), "Expected 12 cells to have been created")
	for _, cell := range cells {
		assert.Assert(t, cell.Sector.Azimuth >= 0 && cell.Sector.Azimuth <= 270, cell.Sector)
		assert.Assert(t, cell.Sector.Arc >= 90 && cell.Sector.Arc <= 120, cell.Sector)
	}
}

func Test_findClosestTowers(t *testing.T) {
	m, err := NewManager()
	assert.NilError(t, err, "Unexpected error creating manager")
	config.Clear()
	towerConfig, err := config.GetTowerConfig("berlin-rectangular-9-1.yaml")
	assert.NilError(t, err)
	m.Cells = NewCells(towerConfig)

	assert.Equal(t, 9, len(m.Cells), "Expected 9 towers to have been created")

	// Test a point outside the towers north-west
	testPointA := &types.Point{Lat: 52.12345, Lng: -8.123}
	towersA, distancesA, err := m.findStrongestCells(testPointA)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(towersA), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesA), "Expected 3 tower distancesA in findClosest")
	assert.Assert(t, distancesA[2] > distancesA[1], "Expected distance to be greater")
	assert.Assert(t, distancesA[1] > distancesA[0], "Expected distance to be greater")

	// Test a point outside the towers south-east
	testPointB := &types.Point{Lat: 51.7654, Lng: -7.9876}
	towersB, distancesB, err := m.findStrongestCells(testPointB)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(towersB), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesB), "Expected 3 tower distancesA in findClosest")
	assert.Assert(t, distancesB[2] > distancesB[1], "Expected distance to be greater")
	assert.Assert(t, distancesB[1] > distancesB[0], "Expected distance to be greater")

	// Test a point within the towers south-east of centre
	testPointC := &types.Point{Lat: 51.980, Lng: -7.950}
	towersC, distancesC, err := m.findStrongestCells(testPointC)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(towersC), "Expected 3 tower names in findClosest")
	assert.Equal(t, 3, len(distancesC), "Expected 3 tower distancesA in findClosest")
	assert.Assert(t, distancesC[2] > distancesC[1], "Expected distance to be greater")
	assert.Assert(t, distancesC[1] > distancesC[0], "Expected distance to be greater")
}

func Test_PowerAdjust(t *testing.T) {
	m, err := NewManager()
	assert.NilError(t, err, "Unexpected error creating manager")

	config.Clear()
	towerConfig, err := config.GetTowerConfig("berlin-rectangular-1-1.yaml")
	assert.NilError(t, err)
	m.Cells = NewCells(towerConfig)

	go func() {
		for event := range m.CellsChannel {
			assert.Equal(t, trafficsim.Type_UPDATED, event.Type)
		}
	}()

	assert.Equal(t, 1, len(m.Cells), "Expected 1 tower to have been created")
	towerID1420 := newEcgi("0001420", utils.TestPlmnID)
	err = m.UpdateCell(towerID1420, -6) // subtracted from initial 10dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	tower1, ok := m.Cells[towerID1420]
	assert.Assert(t, ok)
	assert.Equal(t, 4.0, tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with value too low - capped at -15dB /////////////////////
	err = m.UpdateCell(towerID1420, -30) // subtracted from prev 4dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	assert.Equal(t, -15.0, tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with value too high - capped at 30dB /////////////////////
	err = m.UpdateCell(towerID1420, 50) // Added to prev -15dB
	assert.NilError(t, err, "Unexpected response from adjusting power")
	assert.Equal(t, 30.0, tower1.TxPowerdB, "unexpected value for tower power")

	///////// Try with wrong name /////////////////////
	towerID1421 := newEcgi("0001421", utils.TestPlmnID)
	err = m.UpdateCell(towerID1421, -3)
	assert.Error(t, err, "unknown cell {0001421 315010}", "Expected an error for wrong name when adjusting power")

	time.Sleep(time.Millisecond * 100)
}

func Test_MakeNeighbors(t *testing.T) {
	config.Clear()
	towerConfig, err := config.GetTowerConfig("berlin-rectangular-9-1.yaml")
	assert.NilError(t, err)

	// 1420 --- 1421 --- 1422
	//   |        |        |
	// 1423 --- 1424 --- 1425
	//   |        |        |
	// 1426 --- 1427 --- 1428
	// tower num 2 is the top right - it's id is "0001422"
	cell2Ecgi := types.ECGI{
		EcID:   towerConfig.TowersLayout[2].Sectors[0].EcID,
		PlmnID: towerConfig.TowersLayout[2].PlmnID,
	}
	cells := NewCells(towerConfig)
	// The neighbors are already calculated in the above, but we do it
	// explicitly here
	neighborIDs := makeNeighbors(cells[cell2Ecgi], cells)
	assert.Equal(t, 6, len(neighborIDs), "Unexpected number of neighbors for 1422")
	assert.Equal(t, types.EcID("0001425"), neighborIDs[0].EcID)
	assert.Equal(t, types.EcID("0001421"), neighborIDs[1].EcID)

	// tower num 4 is the middle - it's id is "0001424"
	cell4Ecgi := types.ECGI{
		EcID:   towerConfig.TowersLayout[4].Sectors[0].EcID,
		PlmnID: towerConfig.TowersLayout[4].PlmnID,
	}
	neighborIDs4 := makeNeighbors(cells[cell4Ecgi], cells)
	assert.Equal(t, 6, len(neighborIDs4), "Unexpected number of neighbors for 1424")
	for idx, n := range neighborIDs4 {
		switch n.EcID {
		case "0001421":
		case "0001423":
		case "0001425":
		case "0001427":
			assert.Assert(t, idx < 4, "Expected named cells to be in the closest 4")
		}
	}
}

func Test_StrengthToTower1Sector(t *testing.T) {
	cell := &types.Cell{
		Location: &types.Point{
			Lat: 52,
			Lng: -8,
		},
		Sector: &types.Sector{
			Azimuth: 0,
			Arc:     360,
		},
		TxPowerdB: 10.0, // Does not matter in this case 360
	}
	strength := strengthAtPoint(
		&types.Point{
			Lat: 52.01,
			Lng: -8.01,
		}, cell)
	assert.Equal(t, -2423.0, math.Round(strength*1e3), "Unexpected strength for single sector tower")
}

func Test_StrengthToTower2Sectors(t *testing.T) {
	cell := &types.Cell{
		Location: &types.Point{
			Lat: 52,
			Lng: -8,
		},
		Sector: &types.Sector{
			Azimuth: 90,
			Arc:     180,
		},
		TxPowerdB: 10.0,
	}
	strength := strengthAtPoint(&types.Point{
		Lat: 52.01,
		Lng: -8.01,
	}, cell)
	assert.Equal(t, 1045.0, math.Round(strength*1e3), "Unexpected strength for 2 sector tower")
}

func Test_StrengthToTower3Sectors(t *testing.T) {
	cell := &types.Cell{
		Location: &types.Point{
			Lat: 52,
			Lng: -8,
		},
		Sector: &types.Sector{
			Azimuth: 120,
			Arc:     120,
		},
		TxPowerdB: 10.0,
	}
	strength := strengthAtPoint(
		&types.Point{
			Lat: 52.001, // Much closer
			Lng: -8.001,
		}, cell)
	assert.Equal(t, 8219.0, math.Floor(strength*1e3), "Unexpected strength for 3 sector tower")
}

func Test_PowerToDist(t *testing.T) {
	distm20 := PowerToDist(-20) // -20dB
	assert.Equal(t, 1309, int(math.Floor(distm20*1e5)))

	distm10 := PowerToDist(-10) // -10dB
	assert.Equal(t, 1331, int(math.Floor(distm10*1e5)))

	dist0 := PowerToDist(0) // 0dB
	assert.Equal(t, 1399, int(math.Floor(dist0*1e5)))

	dist10 := PowerToDist(10)
	assert.Equal(t, 1616, int(math.Floor(dist10*1e5)))

	dist20 := PowerToDist(20)
	assert.Equal(t, 2300, int(math.Floor(dist20*1e5)))

	dist30 := PowerToDist(30)
	assert.Equal(t, 4462, int(math.Floor(dist30*1e5)))
}

func Test_AngularAttenuationWideBeam(t *testing.T) {
	cell := &types.Cell{
		Ecgi: &types.ECGI{
			EcID:   "test1",
			PlmnID: utils.TestPlmnID,
		},
		Location: &types.Point{
			Lat: 0,
			Lng: 0,
		},
		Sector: &types.Sector{
			Azimuth:  120, // Beam is facing south east
			Arc:      120, // wide beam
			Centroid: nil,
		},
	}

	// try at 45° above 3 o'clock
	point1 := &types.Point{
		Lat: 0.02, // Further away
		Lng: 0.02,
	}
	aAtt1 := angleAttenuation(point1, cell)
	assert.Equal(t, -2341.0, math.Round(aAtt1*1e3))

	dAtt1 := distanceAttenuation(point1, cell)
	assert.Equal(t, -7258.0, math.Round(dAtt1*1e3))

	// try at 0° - 3 o'clock
	point2 := &types.Point{
		Lat: 0.0,   // to the east
		Lng: 0.001, //Closer
	}
	att2 := angleAttenuation(point2, cell)
	assert.Equal(t, -792.0, math.Round(att2*1e3))

	dAtt2 := distanceAttenuation(point2, cell)
	assert.Equal(t, -0.0, math.Round(dAtt2*1e3)) // Zero as dist = 0.001 = PowerFactor = reference dist

	// try at -30° - middle of beam at 120° azimuth - south east
	point3 := &types.Point{
		Lat: -0.0057735,
		Lng: 0.01,
	}
	att3 := angleAttenuation(point3, cell)
	assert.Equal(t, 0.0, math.Round(att3*1e3))

	dAtt3 := distanceAttenuation(point3, cell)
	assert.Equal(t, -5312.0, math.Round(dAtt3*1e3))
}

func Test_AngularAttenuationNarrowBeam(t *testing.T) {
	cell := &types.Cell{
		Ecgi: &types.ECGI{
			EcID:   "test1",
			PlmnID: utils.TestPlmnID,
		},
		Location: &types.Point{
			Lat: 0,
			Lng: 0,
		},
		Sector: &types.Sector{
			Azimuth:  120, // Beam is facing south east
			Arc:      60,  // narrow beam
			Centroid: nil,
		},
	}

	// try at 45°
	point := &types.Point{
		Lat: 0.02, // Further away
		Lng: 0.02,
	}
	aAtt1 := angleAttenuation(point, cell)
	assert.Equal(t, -7782.0, math.Round(aAtt1*1e3))

	dAtt1 := distanceAttenuation(point, cell)
	assert.Equal(t, -4247.0, math.Round(dAtt1*1e3))

	// try at 0°
	point2 := &types.Point{
		Lat: 0.0,
		Lng: 0.01,
	}
	aAtt2 := angleAttenuation(point2, cell)
	assert.Equal(t, -1761.0, math.Round(aAtt2*1e3))

	dAtt2 := distanceAttenuation(point2, cell)
	assert.Equal(t, -1990.0, math.Round(dAtt2*1e3))

	// try at -30° - middle of beam
	point3 := &types.Point{
		Lat: -0.0057735,
		Lng: 0.01,
	}
	aAtt3 := angleAttenuation(point3, cell)
	assert.Equal(t, 0.0, math.Round(aAtt3*1e3))

	dAtt3 := distanceAttenuation(point3, cell)
	assert.Equal(t, -2302.0, math.Round(dAtt3*1e3))
}
