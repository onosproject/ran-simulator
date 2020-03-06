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
	"math"
	"regexp"
	"strconv"

	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
)

// TestPlmnID - https://en.wikipedia.org/wiki/Mobile_country_code#Test_networks
const TestPlmnID = "001001"

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
func NewTowers(params types.TowersParams, mapLayout types.MapLayout) map[string]*types.Tower {
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
				Name:       towerName,
				Location:   &pos,
				Color:      randomColor(),
				PlmnID:     TestPlmnID,
				EcID:       makeEci(towerName),
				MaxUEs:     params.MaxUEsPerTower,
				Neighbors:  makeNeighbors(towerName, params),
				TxPowerdB:  DefaultTxPower,
				CrntiMap:   make(map[string]string),
				CrntiIndex: 0,
			}
		}
	}

	return towers
}

// Find the closest tower to any point - return closest, candidate1 and candidate2
// in order of distance
// Note this does not take any account of serving - it's just about distance
func (m *Manager) findClosestTowers(point *types.Point) ([]string, []float32) {
	var (
		closest    string
		candidate1 string
		candidate2 string
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
			closest = tower.Name
			closestDist = distance
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
	m.TowersLock.RUnlock()

	return []string{closest, candidate1, candidate2}, []float32{closestDist, candidate1Dist, candidate2Dist}
}

// GetTower returns tower based on its name
func (m *Manager) GetTower(name string) *types.Tower {
	m.TowersLock.RLock()
	defer m.TowersLock.RUnlock()
	return m.Towers[name]
}

// UpdateTower Update a tower's properties - usually power level
func (m *Manager) UpdateTower(tower string, powerAdjust float32) error {
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
func (m *Manager) NewCrnti(servingTower string, ueName string) (string, error) {
	m.TowersLock.Lock()
	defer m.TowersLock.Unlock()
	tower, ok := m.Towers[servingTower]
	if !ok {
		return "", fmt.Errorf("unknown tower %s", servingTower)
	}
	tower.CrntiIndex++
	crnti := fmt.Sprintf("%04X", tower.CrntiIndex%MaxCrnti)
	tower.CrntiMap[crnti] = ueName
	return crnti, nil
}

// DelCrnti deletes a crnti
func (m *Manager) DelCrnti(servingTower string, crnti string) error {
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
func (m *Manager) CrntiToName(crnti string, ecid string) (string, error) {
	tower, ok := m.Towers[EciToName(ecid)]
	if !ok {
		return "", fmt.Errorf("tower %s not found", ecid)
	}
	ueName, ok := tower.CrntiMap[crnti]
	if !ok {
		return "", fmt.Errorf("crnti %s not found", crnti)
	}
	return ueName, nil
}

// EciToName - TODO return the error OR refactor altogether
func EciToName(eci string) string {
	id, _ := strconv.Atoi(eci)
	return fmt.Sprintf("Tower-%d", id)
}

// Measure the distance between a point and a tower and return an answer in decimal degrees
// Simple arithmetic is used, do not use for >= 180 degrees
func distanceToTower(tower *types.Tower, point *types.Point) float32 {
	return float32(math.Hypot(
		float64(tower.GetLocation().GetLat()-point.GetLat()),
		float64(tower.GetLocation().GetLng()-point.GetLng()),
	))
}

func makeEci(towerName string) string {
	re := regexp.MustCompile("[0-9]+")
	id, _ := strconv.Atoi(re.FindAllString(towerName, 1)[0])
	return fmt.Sprintf("%07X", id)
}

func makeNeighbors(towerName string, towerParams types.TowersParams) []string {
	neighbors := make([]string, 0, 8)
	re := regexp.MustCompile("[0-9]+")
	id, _ := strconv.Atoi(re.FindAllString(towerName, 1)[0])
	id--

	nrows := int(towerParams.TowerRows)
	ncols := int(towerParams.TowerCols)

	i := id / nrows
	j := id % ncols

	for x := max(0, i-1); x <= min(i+1, nrows-1); x++ {
		for y := max(0, j-1); y <= min(j+1, ncols-1); y++ {
			if (x == i && y == j-1) || (x == i && y == j+1) || (x == i-1 && y == j) || (x == i+1 && y == j) {
				towerNum := x*nrows + y + 1
				towerName := fmt.Sprintf("Tower-%d", towerNum)
				neighbors = append(neighbors, towerName)
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

// GetTowerByName ...
func (m *Manager) GetTowerByName(name string) *types.Tower {
	for _, tower := range m.Towers {
		if tower.Name == name {
			return tower
		}
	}
	return nil
}
