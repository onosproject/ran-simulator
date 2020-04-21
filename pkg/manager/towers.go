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
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"math"
	"math/rand"
)

const (
	// DefaultTxPower - all cells start with this power level
	DefaultTxPower = 10

	// PowerFactor - relate power to distance in decimal degrees
	PowerFactor = 0.001
)

const defaultColor = "#000000"

const (
	maxPowerdB = 30.0
	minPowerdB = -15.0
)

// MaxCrnti - Maximum value of CRNTI
const MaxCrnti = 65523

// InvalidCrnti ...
const InvalidCrnti = "0000"

// CellIf :
type CellIf interface {
	GetPosition() types.Point
}

// NewCells - create a set of new Cells
func NewCells(params types.TowersParams, mapLayout types.MapLayout) map[types.ECGI]*types.Cell {
	cells := make(map[types.ECGI]*types.Cell)

	totalCells := uint32(math.Floor(float64(params.TowerRows*params.TowerCols) * float64(params.AvgCellsPerTower)))
	minCellsPerTower := totalCells / (params.TowerRows * params.TowerCols)
	remainder := uint32(math.Mod(float64(totalCells), float64(minCellsPerTower)))
	extraCells := make([]uint32, remainder)
	for i := range extraCells {
		extraCells[i] = uint32(rand.Int31n(int32(params.TowerRows * params.TowerCols)))
	}

	var r, c, cellNum uint32
	var cellPort = uint32(utils.GrpcBasePort + 1) // Start at 5152 so it appears as 1420 in Hex
	for r = 0; r < params.TowerRows; r++ {
		for c = 0; c < params.TowerCols; c++ {
			towerNum := r*params.TowerCols + c
			pos := getTowerPosition(r, c, params, mapLayout)
			numCells := minCellsPerTower
			for _, e := range extraCells {
				if e == towerNum {
					numCells++
				}
			}
			for cellNum = 0; cellNum < numCells; cellNum++ {
				cellPort++
				ecgi := types.ECGI{
					PlmnID: utils.TestPlmnID,
					EcID:   utils.EcIDForPort(int(cellPort)),
				}
				cells[ecgi] = &types.Cell{
					Location:   pos,
					Color:      utils.RandomColor(),
					Ecgi:       &ecgi,
					MaxUEs:     params.MaxUEsPerCell,
					Neighbors:  makeNeighbors(int(cellNum), params),
					TxPowerdB:  DefaultTxPower,
					Port:       cellPort,
					CrntiMap:   make(map[types.Crnti]types.Imsi),
					CrntiIndex: 0,
					Sector: &types.Sector{
						Azimuth: int32(float64(cellNum) / float64(numCells) * 360),
						Arc:     int32(360 / numCells),
					},
				}
			}
		}
	}

	return cells
}

// Find the closest cell to any point - return closest, candidate1 and candidate2
// in order of distance
// Note this does not take any account of serving - it's just about distance
func (m *Manager) findClosestCells(point *types.Point) ([]*types.ECGI, []float32) {
	var (
		closest    *types.ECGI
		candidate1 *types.ECGI
		candidate2 *types.ECGI
	)

	var (
		closestDist    float32 = math.MaxFloat32
		candidate1Dist float32 = math.MaxFloat32
		candidate2Dist float32 = math.MaxFloat32
	)

	m.CellsLock.RLock()
	for _, cell := range m.Cells {
		distance := distanceToCellCentroid(cell, point)
		if distance < closestDist {
			candidate2 = candidate1
			candidate2Dist = candidate1Dist
			candidate1 = closest
			candidate1Dist = closestDist
			closest = cell.Ecgi
			closestDist = distance
		} else if distance < candidate1Dist {
			candidate2 = candidate1
			candidate2Dist = candidate1Dist
			candidate1 = cell.Ecgi
			candidate1Dist = distance
		} else if distance < candidate2Dist {
			candidate2 = cell.Ecgi
			candidate2Dist = distance
		}
	}
	m.CellsLock.RUnlock()

	return []*types.ECGI{closest, candidate1, candidate2}, []float32{closestDist, candidate1Dist, candidate2Dist}
}

// GetCell returns tower based on its name
func (m *Manager) GetCell(name types.ECGI) *types.Cell {
	m.CellsLock.RLock()
	defer m.CellsLock.RUnlock()
	return m.Cells[name]
}

// UpdateCell Update a tower's properties - usually power level
func (m *Manager) UpdateCell(tower types.ECGI, powerAdjust float32) error {
	// Only the power can be updated at present
	m.CellsLock.Lock()
	t, ok := m.Cells[tower]
	if !ok {
		m.CellsLock.Unlock()
		return fmt.Errorf("unknown tower %s", tower)
	}
	currentPower := t.TxPowerdB
	if currentPower+powerAdjust < minPowerdB {
		t.TxPowerdB = minPowerdB
	} else if currentPower+powerAdjust > maxPowerdB {
		t.TxPowerdB = maxPowerdB
	} else {
		t.TxPowerdB += powerAdjust
	}
	m.CellsLock.Unlock()
	m.TowerChannel <- dispatcher.Event{
		Type:   trafficsim.Type_UPDATED,
		Object: t,
	}
	return nil
}

// NewCrnti allocs a new crnti
func (m *Manager) NewCrnti(servingTower *types.ECGI, imsi types.Imsi) (types.Crnti, error) {
	m.CellsLock.Lock()
	defer m.CellsLock.Unlock()
	tower, ok := m.Cells[*servingTower]
	if !ok {
		return "", fmt.Errorf("unknown tower %s", servingTower)
	}
	tower.CrntiIndex++
	crnti := types.Crnti(fmt.Sprintf("%04X", tower.CrntiIndex%MaxCrnti))
	tower.CrntiMap[crnti] = imsi
	return crnti, nil
}

// DelCrnti deletes a crnti
func (m *Manager) DelCrnti(servingTower *types.ECGI, crnti types.Crnti) error {
	m.CellsLock.Lock()
	defer m.CellsLock.Unlock()
	tower, ok := m.Cells[*servingTower]
	if !ok {
		return fmt.Errorf("unknown tower %s", servingTower)
	}
	crntiMap := tower.CrntiMap
	delete(crntiMap, crnti)
	return nil
}

// CrntiToName ...
func (m *Manager) CrntiToName(crnti types.Crnti, ecid *types.ECGI) (types.Imsi, error) {
	m.CellsLock.RLock()
	defer m.CellsLock.RUnlock()
	tower, ok := m.Cells[*ecid]
	if !ok {
		return 0, fmt.Errorf("tower %s not found", ecid)
	}
	imsi, ok := tower.CrntiMap[crnti]
	if !ok {
		return 0, fmt.Errorf("crnti %s/%s not found", ecid, crnti)
	}
	return imsi, nil
}

// Measure the distance between a point and a cell centroid and return an answer in decimal degrees
// Simple arithmetic is used, do not use for lat or long diff >= 100 degrees
func distanceToCellCentroid(cell *types.Cell, point *types.Point) float32 {
	if cell.Sector.Arc == 360 || cell.Sector.Arc == 0 {
		return float32(math.Hypot(
			float64(cell.GetLocation().GetLat()-point.GetLat()),
			float64(cell.GetLocation().GetLng()-point.GetLng()),
		))
	}
	// Work out the location of the centroid of the cell - ref https://en.wikipedia.org/wiki/Circular_sector
	alpha := 2 * math.Pi * float64(cell.Sector.Arc) / 360 / 2
	dist := 2 * PowerToDist(cell.TxPowerdB) * math.Sin(alpha) / alpha / 3
	var azRads float64 = 0
	if cell.Sector.Azimuth != 90 {
		azRads = math.Pi * 2 / float64(cell.Sector.Azimuth-90) / 360
	}
	centroidLat := math.Sin(azRads)*dist + float64(cell.Location.GetLat())
	centroidLng := math.Cos(azRads)*dist + float64(cell.Location.GetLng())
	return float32(math.Hypot(
		centroidLat-float64(point.GetLat()),
		centroidLng-float64(point.GetLng()),
	))
}

// PowerToDist - convert power in dB to distance in decimal degrees
func PowerToDist(power float32) float64 {
	return math.Sqrt(math.Pow(10, float64(power)/10)) * PowerFactor
}

func makeNeighbors(id int, towerParams types.TowersParams) []*types.ECGI {
	neighbors := make([]*types.ECGI, 0)

	nrows := int(towerParams.TowerRows)
	ncols := int(towerParams.TowerCols)

	i := id / nrows
	j := id % ncols

	for x := max(0, i-1); x <= min(i+1, nrows-1); x++ {
		for y := max(0, j-1); y <= min(j+1, ncols-1); y++ {
			if (x == i && y == j-1) || (x == i && y == j+1) || (x == i-1 && y == j) || (x == i+1 && y == j) {
				cellID := x*nrows + y + 2 + utils.GrpcBasePort
				cellEcgi := newEcgi(utils.EcIDForPort(cellID), utils.TestPlmnID)
				neighbors = append(neighbors, &cellEcgi)
			}
		}
	}
	return neighbors
}

// min ...
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// max ...
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func newEcgi(id types.EcID, plmnID types.PlmnID) types.ECGI {
	return types.ECGI{EcID: id, PlmnID: plmnID}
}
