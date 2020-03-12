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

// DefaultTxPower - all base-stations start with this power level
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

// TowerIf :
type TowerIf interface {
	GetPosition() types.Point
}

// NewTowers - create a set of new towers
func NewTowers(params types.TowersParams, mapLayout types.MapLayout) map[types.EcID]*types.Tower {
	towers := make(map[types.EcID]*types.Tower)

	var r, c uint32
	for r = 0; r < params.TowerRows; r++ {
		for c = 0; c < params.TowerCols; c++ {
			pos := getTowerPosition(r, c, params, mapLayout)
			towerNum := r*params.TowerCols + c
			towerPort := utils.GrpcBasePort + towerNum + 2 // Start at 5152 so it appears as 1420 in Hex
			ecid := utils.EcIDForPort(int(towerPort))
			towers[ecid] = &types.Tower{
				Location:   pos,
				Color:      utils.RandomColor(),
				PlmnID:     utils.TestPlmnID,
				EcID:       ecid,
				MaxUEs:     params.MaxUEsPerTower,
				Neighbors:  makeNeighbors(int(towerNum), params),
				TxPowerdB:  DefaultTxPower,
				Port:       towerPort,
				CrntiMap:   make(map[types.Crnti]types.UEName),
				CrntiIndex: 0,
			}
		}
	}

	return towers
}

// Find the closest tower to any point - return closest, candidate1 and candidate2
// in order of distance
// Note this does not take any account of serving - it's just about distance
func (m *Manager) findClosestTowers(point *types.Point) ([]types.EcID, []float32) {
	var (
		closest    types.EcID
		candidate1 types.EcID
		candidate2 types.EcID
	)

	var (
		closestDist    float32 = math.MaxFloat32
		candidate1Dist float32 = math.MaxFloat32
		candidate2Dist float32 = math.MaxFloat32
	)

	m.TowersLock.RLock()
	for _, tower := range m.Towers {
		distance := distanceToTower(tower, point)
		if distance < closestDist {
			candidate2 = candidate1
			candidate2Dist = candidate1Dist
			candidate1 = closest
			candidate1Dist = closestDist
			closest = tower.EcID
			closestDist = distance
		} else if distance < candidate1Dist {
			candidate2 = candidate1
			candidate2Dist = candidate1Dist
			candidate1 = tower.EcID
			candidate1Dist = distance
		} else if distance < candidate2Dist {
			candidate2 = tower.EcID
			candidate2Dist = distance
		}
	}
	m.TowersLock.RUnlock()

	return []types.EcID{closest, candidate1, candidate2}, []float32{closestDist, candidate1Dist, candidate2Dist}
}

// GetTower returns tower based on its name
func (m *Manager) GetTower(name types.EcID) *types.Tower {
	m.TowersLock.RLock()
	defer m.TowersLock.RUnlock()
	return m.Towers[name]
}

// UpdateTower Update a tower's properties - usually power level
func (m *Manager) UpdateTower(tower types.EcID, powerAdjust float32) error {
	// Only the power can be updated at present
	m.TowersLock.Lock()
	t, ok := m.Towers[tower]
	if !ok {
		m.TowersLock.Unlock()
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
	m.TowersLock.Unlock()
	m.TowerChannel <- dispatcher.Event{
		Type:   trafficsim.Type_UPDATED,
		Object: t,
	}
	return nil
}

// NewCrnti allocs a new crnti
func (m *Manager) NewCrnti(servingTower types.EcID, ueName types.UEName) (types.Crnti, error) {
	m.TowersLock.Lock()
	defer m.TowersLock.Unlock()
	tower, ok := m.Towers[servingTower]
	if !ok {
		return "", fmt.Errorf("unknown tower %s", servingTower)
	}
	tower.CrntiIndex++
	crnti := types.Crnti(fmt.Sprintf("%04X", tower.CrntiIndex%MaxCrnti))
	tower.CrntiMap[crnti] = ueName
	return crnti, nil
}

// DelCrnti deletes a crnti
func (m *Manager) DelCrnti(servingTower types.EcID, crnti types.Crnti) error {
	m.TowersLock.Lock()
	defer m.TowersLock.Unlock()
	tower, ok := m.Towers[servingTower]
	if !ok {
		return fmt.Errorf("unknown tower %s", servingTower)
	}
	crntiMap := tower.CrntiMap
	delete(crntiMap, crnti)
	return nil
}

// CrntiToName ...
func (m *Manager) CrntiToName(crnti types.Crnti, ecid types.EcID) (types.UEName, error) {
	tower, ok := m.Towers[ecid]
	if !ok {
		return "", fmt.Errorf("tower %s not found", ecid)
	}
	ueName, ok := tower.CrntiMap[crnti]
	if !ok {
		return "", fmt.Errorf("crnti %s/%s not found", ecid, crnti)
	}
	return ueName, nil
}

// Measure the distance between a point and a tower and return an answer in decimal degrees
// Simple arithmetic is used, do not use for >= 180 degrees
func distanceToTower(tower *types.Tower, point *types.Point) float32 {
	return float32(math.Hypot(
		float64(tower.GetLocation().GetLat()-point.GetLat()),
		float64(tower.GetLocation().GetLng()-point.GetLng()),
	))
}

func makeNeighbors(id int, towerParams types.TowersParams) []types.EcID {
	neighbors := make([]types.EcID, 0, 8)

	nrows := int(towerParams.TowerRows)
	ncols := int(towerParams.TowerCols)

	i := id / nrows
	j := id % ncols

	for x := max(0, i-1); x <= min(i+1, nrows-1); x++ {
		for y := max(0, j-1); y <= min(j+1, ncols-1); y++ {
			if (x == i && y == j-1) || (x == i && y == j+1) || (x == i-1 && y == j) || (x == i+1 && y == j) {
				towerID := x*nrows + y + 2 + utils.GrpcBasePort
				towerEcID := utils.EcIDForPort(towerID)
				neighbors = append(neighbors, towerEcID)
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
