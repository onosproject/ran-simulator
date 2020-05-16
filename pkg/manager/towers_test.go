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
	"github.com/onosproject/onos-topo/pkg/bulk"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"gotest.tools/assert"
	"math"
	"testing"
	"time"
)

func Test_NewTowers(t *testing.T) {
	topoDeviceConfig, err := bulk.GetDeviceConfig("berlin-honeycomb-4-3-topo.yaml")
	assert.NilError(t, err)

	cells := make(map[types.ECGI]*types.Cell)

	for _, td := range topoDeviceConfig.TopoDevices {
		td := td //pin
		cell, err := NewCell(&td)
		assert.NilError(t, err)
		cells[*cell.Ecgi] = cell
	}

	assert.Equal(t, 12, len(cells), "Expected 12 cells to have been created")
	for _, cell := range cells {
		assert.Assert(t, cell.Sector.Azimuth >= 0 && cell.Sector.Azimuth <= 270, cell.Sector)
		assert.Assert(t, cell.Sector.Arc >= 90 && cell.Sector.Arc <= 120, cell.Sector)
	}
}

func Test_findStrongestTowers(t *testing.T) {
	m, err := NewManager()
	assert.NilError(t, err)
	bulk.Clear()
	topoDeviceConfig, err := bulk.GetDeviceConfig("berlin-honeycomb-4-3-topo.yaml")
	assert.NilError(t, err)

	cells := make(map[types.ECGI]*types.Cell)

	for _, td := range topoDeviceConfig.TopoDevices {
		td := td //pin
		cell, err := NewCell(&td)
		assert.NilError(t, err)
		cells[*cell.Ecgi] = cell
	}
	m.Cells = cells

	assert.Equal(t, 12, len(m.Cells), "Expected 12 towers to have been created")

	// Test a point outside the towers north-west
	testPointA := &types.Point{Lat: 52.54, Lng: 13.38}
	towersA, strengthsA, err := m.findStrongestCells(testPointA)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(towersA), "Expected 3 tower names in findStrongestCells")
	assert.Equal(t, 3, len(strengthsA), "Expected 3 tower strengths in findStrongestCells")
	assert.Assert(t, strengthsA[2] < strengthsA[1], "Expected strength to be less")
	assert.Assert(t, strengthsA[1] < strengthsA[0], "Expected strength to be less")
	assert.Assert(t, towersA[0].String() != towersA[1].String())
	assert.Assert(t, towersA[0].String() != towersA[2].String())
	assert.Assert(t, towersA[1].String() != towersA[2].String())
	t.Logf("Cell %s (%f dB), Cell %s (%f dB), Cell %s (%f dB)",
		towersA[0].EcID, math.Floor(strengthsA[0]*100)/100,
		towersA[1].EcID, math.Floor(strengthsA[1]*100)/100,
		towersA[2].EcID, math.Floor(strengthsA[2]*100)/100)

	// Test a point outside the towers south-east
	testPointB := &types.Point{Lat: 52.50, Lng: 13.37}
	towersB, strengthsB, err := m.findStrongestCells(testPointB)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(towersB), "Expected 3 tower names in findStrongestCells")
	assert.Equal(t, 3, len(strengthsB), "Expected 3 tower strengths in findStrongestCells")
	assert.Assert(t, strengthsB[2] < strengthsB[1], "Expected strength to be less")
	assert.Assert(t, strengthsB[1] < strengthsB[0], "Expected strength to be less")
	assert.Assert(t, towersB[0].String() != towersB[1].String())
	assert.Assert(t, towersB[0].String() != towersB[2].String())
	assert.Assert(t, towersB[1].String() != towersB[2].String())
	t.Logf("Cell %s (%f dB), Cell %s (%f dB), Cell %s (%f dB)",
		towersB[0].EcID, math.Floor(strengthsB[0]*100)/100,
		towersB[1].EcID, math.Floor(strengthsB[1]*100)/100,
		towersB[2].EcID, math.Floor(strengthsB[2]*100)/100)

	// Test a point within the towers south-east of centre
	testPointC := &types.Point{Lat: 52.50, Lng: 13.43}
	towersC, strengthsC, err := m.findStrongestCells(testPointC)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(towersC), "Expected 3 tower names in findStrongestCells")
	assert.Equal(t, 3, len(strengthsC), "Expected 3 tower strengths in findStrongestCells")
	assert.Assert(t, strengthsC[2] < strengthsC[1], "Expected strength to be less")
	assert.Assert(t, strengthsC[1] < strengthsC[0], "Expected strength to be less")
	assert.Assert(t, towersC[0].String() != towersC[1].String())
	assert.Assert(t, towersC[0].String() != towersC[2].String())
	assert.Assert(t, towersC[1].String() != towersC[2].String())
	t.Logf("Cell %s (%f dB), Cell %s (%f dB), Cell %s (%f dB)",
		towersC[0].EcID, math.Floor(strengthsC[0]*100)/100,
		towersC[1].EcID, math.Floor(strengthsC[1]*100)/100,
		towersC[2].EcID, math.Floor(strengthsC[2]*100)/100)
}

func Test_PowerAdjust(t *testing.T) {
	m, err := NewManager()
	assert.NilError(t, err, "Unexpected error creating manager")
	bulk.Clear()
	topoDeviceConfig, err := bulk.GetDeviceConfig("berlin-rectangular-4-1-topo.yaml")
	assert.NilError(t, err)

	cells := make(map[types.ECGI]*types.Cell)

	for _, td := range topoDeviceConfig.TopoDevices {
		td := td //pin
		cell, err := NewCell(&td)
		assert.NilError(t, err)
		cells[*cell.Ecgi] = cell
	}
	m.Cells = cells

	go func() {
		for event := range m.CellsChannel {
			assert.Equal(t, trafficsim.Type_UPDATED, event.Type)
		}
	}()

	assert.Equal(t, 4, len(m.Cells), "Expected 4 towers to have been created")
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
	towerID1431 := newEcgi("0001431", utils.TestPlmnID)
	err = m.UpdateCell(towerID1431, -3)
	assert.Error(t, err, "unknown cell {0001431 315010}", "Expected an error for wrong name when adjusting power")

	time.Sleep(time.Millisecond * 100)
}

func Test_MakeNeighbors(t *testing.T) {
	bulk.Clear()
	topoDeviceConfig, err := bulk.GetDeviceConfig("berlin-rectangular-9-1-topo.yaml")
	assert.NilError(t, err)

	cells := make(map[types.ECGI]*types.Cell)

	for _, td := range topoDeviceConfig.TopoDevices {
		td := td //pin
		cell, err := NewCell(&td)
		assert.NilError(t, err)
		cells[*cell.Ecgi] = cell
	}

	// 1420 --- 1421 --- 1422
	//   |        |        |
	// 1423 --- 1424 --- 1425
	//   |        |        |
	// 1426 --- 1427 --- 1428
	// tower num 2 is the top right - it's id is "0001422"
	cell2Ecgi := types.ECGI{
		EcID:   types.EcID(topoDeviceConfig.TopoDevices[2].Attributes[types.EcidKey]),
		PlmnID: types.PlmnID(topoDeviceConfig.TopoDevices[2].Attributes[types.PlmnIDKey]),
	}
	// The neighbors are already calculated in the above, but we do it
	// explicitly here
	neighborIDs := makeNeighbors(cells[cell2Ecgi], cells)
	assert.Equal(t, 6, len(neighborIDs), "Unexpected number of neighbors for 1422")
	assert.Equal(t, types.EcID("0001425"), neighborIDs[0].EcID)
	assert.Equal(t, types.EcID("0001421"), neighborIDs[1].EcID)

	// tower num 4 is the middle - it's id is "0001424"
	cell4Ecgi := types.ECGI{
		EcID:   types.EcID(topoDeviceConfig.TopoDevices[4].Attributes[types.EcidKey]),
		PlmnID: types.PlmnID(topoDeviceConfig.TopoDevices[4].Attributes[types.PlmnIDKey]),
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
	assert.Equal(t, -1551.0, math.Round(strength*1e3), "Unexpected strength for single sector tower")
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
	assert.Equal(t, -1173.0, math.Round(strength*1e3), "Unexpected strength for 2 sector tower")
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
	assert.Equal(t, -2195.0, math.Floor(strength*1e3), "Unexpected strength for 3 sector tower")
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
			Azimuth: 120, // Beam is facing south east
			Arc:     120, // wide beam
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

	// try at 135° - north west - 10:30 o'clock
	point4 := &types.Point{
		Lat: 0.02,
		Lng: -0.02,
	}
	att4 := angleAttenuation(point4, cell)
	assert.Equal(t, -10792.0, math.Round(att4*1e3))

	dAtt4 := distanceAttenuation(point4, cell)
	assert.Equal(t, -7258.0, math.Round(dAtt4*1e3))
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
			Azimuth: 120, // Beam is facing south east
			Arc:     60,  // narrow beam
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

func Test_centroidPosition(t *testing.T) {
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
			Azimuth: 120, // Beam is facing south east
			Arc:     60,  // narrow beam
		},
	}

	point := centroidPosition(cell)

	assert.Equal(t, -0.0044563384065730675, point.GetLat())
	assert.Equal(t, 0.007718604535905087, point.GetLng())

}
