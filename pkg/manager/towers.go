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
)

// DefaultTxPower - all cells start with this power level
const DefaultTxPower = 10

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

// NewTowers - create a set of new towers
func NewTowers(params types.TowersParams, mapLayout types.MapLayout) map[types.ECGI]*types.Cell {
	towers := make(map[types.ECGI]*types.Cell)

	var r, c uint32
	for r = 0; r < params.TowerRows; r++ {
		for c = 0; c < params.TowerCols; c++ {
			pos := getTowerPosition(r, c, params, mapLayout)
			towerNum := r*params.TowerCols + c
			towerPort := utils.GrpcBasePort + towerNum + 2 // Start at 5152 so it appears as 1420 in Hex
			ecgi := types.ECGI{
				PlmnID: utils.TestPlmnID,
				EcID:   utils.EcIDForPort(int(towerPort)),
			}
			towers[ecgi] = &types.Cell{
				Location:   pos,
				Color:      utils.RandomColor(),
				Ecgi:       &ecgi,
				MaxUEs:     params.MaxUEsPerCell,
				Neighbors:  makeNeighbors(int(towerNum), params),
				TxPowerdB:  DefaultTxPower,
				Port:       towerPort,
				CrntiMap:   make(map[types.Crnti]types.Imsi),
				CrntiIndex: 0,
			}
		}
	}

	return towers
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
	for _, tower := range m.Cells {
		distance := distanceToCell(tower, point)
		if distance < closestDist {
			candidate2 = candidate1
			candidate2Dist = candidate1Dist
			candidate1 = closest
			candidate1Dist = closestDist
			closest = tower.Ecgi
			closestDist = distance
		} else if distance < candidate1Dist {
			candidate2 = candidate1
			candidate2Dist = candidate1Dist
			candidate1 = tower.Ecgi
			candidate1Dist = distance
		} else if distance < candidate2Dist {
			candidate2 = tower.Ecgi
			candidate2Dist = distance
		}
	}
	m.CellsLock.RUnlock()

	return []*types.ECGI{closest, candidate1, candidate2}, []float32{closestDist, candidate1Dist, candidate2Dist}
}

// GetTower returns tower based on its name
func (m *Manager) GetTower(name types.ECGI) *types.Cell {
	m.CellsLock.RLock()
	defer m.CellsLock.RUnlock()
	return m.Cells[name]
}

// UpdateTower Update a tower's properties - usually power level
func (m *Manager) UpdateTower(tower types.ECGI, powerAdjust float32) error {
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

// Measure the distance between a point and a tower and return an answer in decimal degrees
// Simple arithmetic is used, do not use for >= 180 degrees
func distanceToCell(tower *types.Cell, point *types.Point) float32 {
	return float32(math.Hypot(
		float64(tower.GetLocation().GetLat()-point.GetLat()),
		float64(tower.GetLocation().GetLng()-point.GetLng()),
	))
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
				towerID := x*nrows + y + 2 + utils.GrpcBasePort
				towerEcgi := newEcgi(utils.EcIDForPort(towerID), utils.TestPlmnID)
				neighbors = append(neighbors, &towerEcgi)
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
