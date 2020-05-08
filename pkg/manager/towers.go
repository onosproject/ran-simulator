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
	"sort"

	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/config"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"github.com/onosproject/ran-simulator/pkg/utils"
)

const (
	// DefaultTxPower - all cells start with this power level
	DefaultTxPower = 10

	// PowerFactor - relate power to distance in decimal degrees
	PowerFactor = 0.001

	// PowerBase - baseline for power to distance in decimal degrees
	PowerBase = 0.013

	// MaxNeighbours to find - useful limit for hex layouts
	MaxNeighbours = 6
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
func NewCells(towersConfig config.TowerConfig) map[types.ECGI]*types.Cell {
	cells := make(map[types.ECGI]*types.Cell)

	for _, tower := range towersConfig.TowersLayout {
		towerLoc := types.Point{
			Lat: tower.Latitude,
			Lng: tower.Longitude,
		}
		for _, sector := range tower.Sectors {
			ecgi := types.ECGI{
				PlmnID: tower.PlmnID,
				EcID:   sector.EcID,
			}
			cell := &types.Cell{
				Location:   &towerLoc,
				Color:      utils.RandomColor(),
				Ecgi:       &ecgi,
				MaxUEs:     uint32(sector.MaxUEs),
				TxPowerdB:  sector.InitPowerDb,
				Port:       uint32(sector.GrpcPort),
				CrntiMap:   make(map[types.Crnti]types.Imsi),
				CrntiIndex: 0,
				Sector: &types.Sector{
					Azimuth: int32(sector.Azimuth),
					Arc:     int32(sector.Arc),
				},
				ConfigAttributes: make(map[types.ConfigKey]types.ConfigValue),
			}
			cell.Sector.Centroid = centroidPosition(cell)
			cells[ecgi] = cell
		}
	}

	// go through again and update neighbours, now that all centroids have been calculated
	for _, cell := range cells {
		cell.Neighbors = makeNeighbors(cell, cells)
	}

	return cells
}

// Find the strongest power signal cell to any point - return strongest, candidate1 and candidate2
// in order of power. This is derived from
// 1. the distance of the point from cell
// 2. the arc and the azimuth of the cell
// 3. the power setting of the antenna
// Note this does not take any account of who's serving - it's just about power
// values are in dB
func (m *Manager) findStrongestCells(point *types.Point) ([]*types.ECGI, []float64, error) {
	var (
		strongest  *types.ECGI
		candidate1 *types.ECGI
		candidate2 *types.ECGI
	)

	strongestStr := -math.MaxFloat64
	candidate1Str := -math.MaxFloat64
	candidate2Str := -math.MaxFloat64

	m.CellsLock.RLock()
	for _, cell := range m.Cells {
		strength := strengthAtPoint(point, cell)

		if strength > strongestStr {
			candidate2 = candidate1
			candidate2Str = candidate1Str
			candidate1 = strongest
			candidate1Str = strongestStr
			strongest = cell.Ecgi
			strongestStr = strength
		} else if strength > candidate1Str {
			candidate2 = candidate1
			candidate2Str = candidate1Str
			candidate1 = cell.Ecgi
			candidate1Str = strength
		} else if strength > candidate2Str {
			candidate2 = cell.Ecgi
			candidate2Str = strength
		}
	}
	m.CellsLock.RUnlock()

	return []*types.ECGI{strongest, candidate1, candidate2},
		[]float64{strongestStr, candidate1Str, candidate2Str}, nil
}

func strengthAtPoint(point *types.Point, cell *types.Cell) float64 {
	distAtt := distanceAttenuation(point, cell)
	angleAtt := angleAttenuation(point, cell)

	return cell.TxPowerdB + distAtt + angleAtt
}

// distanceAttenuation is the antenna Gain as a function of the dist
// a very rough approximation to take in to account the width of
// the antenna beam. A 120° wide beam with 30° height will span ≅ 2x0.5 = 1 steradians
// A 60° wide beam will be half that and so will have double the gain
// https://en.wikipedia.org/wiki/Sector_antenna
// https://en.wikipedia.org/wiki/Steradian
func distanceAttenuation(point *types.Point, cell *types.Cell) float64 {
	latDist := point.GetLat() - cell.GetLocation().GetLat()
	realLngDist := (point.GetLng() - cell.GetLocation().GetLng()) / utils.AspectRatio(cell.GetLocation())
	r := math.Hypot(latDist, realLngDist)
	gain := 120.0 / float64(cell.GetSector().GetArc())
	return 10 * math.Log10(gain*math.Sqrt(PowerFactor/r))
}

// angleAttenuation is the attenuation of power reaching a UE due to its
// position off the centre of the beam in dB
// It is an approximation of the directivity of the antenna
// https://en.wikipedia.org/wiki/Radiation_pattern
// https://en.wikipedia.org/wiki/Sector_antenna
func angleAttenuation(point *types.Point, cell *types.Cell) float64 {

	azRads := utils.AzimuthToRads(float64(cell.Sector.Azimuth))
	pointRads := math.Atan2(point.Lat-cell.Location.Lat, point.Lng-cell.Location.Lng)
	angularOffset := math.Abs(azRads - pointRads)
	angleScaling := float64(cell.Sector.Arc) / 120.0 // Compensate for narrower beams

	// We just use a simple linear formula 0 => no loss
	// 33° => -3dB for a 120° sector according to [2]
	// assume this is 1:1 rads:attenuation e.g. 0.50 rads = 0.5 = -3dB attenuation
	return 10 * math.Log10(1-(angularOffset/math.Pi/angleScaling))
}

// GetCell returns tower based on its name
func (m *Manager) GetCell(name types.ECGI) *types.Cell {
	m.CellsLock.RLock()
	defer m.CellsLock.RUnlock()
	return m.Cells[name]
}

// UpdateCell Update a cell's properties - usually power level
func (m *Manager) UpdateCell(cell types.ECGI, powerAdjust float32) error {
	// Only the power can be updated at present
	m.CellsLock.Lock()
	c, ok := m.Cells[cell]
	if !ok {
		m.CellsLock.Unlock()
		return fmt.Errorf("unknown cell %s", cell)
	}
	currentPower := c.TxPowerdB
	if currentPower+float64(powerAdjust) < minPowerdB {
		c.TxPowerdB = minPowerdB
	} else if currentPower+float64(powerAdjust) > maxPowerdB {
		c.TxPowerdB = maxPowerdB
	} else {
		c.TxPowerdB += float64(powerAdjust)
	}
	c.GetSector().Centroid = centroidPosition(c)
	m.CellsLock.Unlock()
	m.CellsChannel <- dispatcher.Event{
		Type:   trafficsim.Type_UPDATED,
		Object: c,
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
func (m *Manager) CrntiToName(crnti types.Crnti, ecid types.ECGI) (types.Imsi, error) {
	m.CellsLock.RLock()
	defer m.CellsLock.RUnlock()
	tower, ok := m.Cells[ecid]
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
// Centroid is used **only** for the display of the beam on the GUI and for
// calculating Neighbours once at startup
// Simple arithmetic is used, do not use for lat or long diff >= 100 degrees
func centroidPosition(cell *types.Cell) *types.Point {
	if cell.Sector.Arc == 360 || cell.Sector.Arc == 0 {
		return cell.Location
	}
	// Work out the location of the centroid of the cell - ref https://en.wikipedia.org/wiki/List_of_centroids
	alpha := 2 * math.Pi * float64(cell.Sector.Arc) / 360 / 2
	dist := 2 * PowerToDist(cell.TxPowerdB) * math.Sin(alpha) / alpha / 3
	var azRads float64 = 0
	if cell.Sector.Azimuth != 90 {
		azRads = math.Pi * 2 * float64(90-cell.Sector.Azimuth) / 360
	}
	aspectRatio := utils.AspectRatio(cell.Location)
	return &types.Point{
		Lat: math.Sin(azRads)*dist + cell.Location.GetLat(),
		Lng: math.Cos(azRads)*dist/aspectRatio + cell.Location.GetLng(),
	}
}

// PowerToDist - convert power in dB to distance in decimal degrees
// Like centroid this is now used only for calculating centroid, which is
// only for the GUI and the neighbours
func PowerToDist(power float64) float64 {
	return math.Sqrt(math.Pow(10, power/10))*PowerFactor + PowerBase
}

// find the neighbours of a cell - not distance from towers, but from centroids
func makeNeighbors(self *types.Cell, allCells map[types.ECGI]*types.Cell) []*types.ECGI {
	type distance struct {
		id   *types.ECGI
		dist float64
	}
	distances := make([]distance, 0)

	selfCentroid := self.Sector.Centroid
	for _, otherCell := range allCells {
		if otherCell.Ecgi == self.Ecgi {
			continue
		}
		dist := math.Hypot(
			selfCentroid.Lng-otherCell.Sector.Centroid.Lng,
			selfCentroid.Lat-otherCell.Sector.Centroid.Lat,
		)
		distances = append(distances, distance{id: otherCell.Ecgi, dist: dist})
	}
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].dist < distances[j].dist
	})

	limit := len(allCells) - 1
	if limit > MaxNeighbours {
		limit = MaxNeighbours
	}
	neighbours := make([]*types.ECGI, limit)
	for i := 0; i < limit; i++ {
		neighbours[i] = distances[i].id
	}

	return neighbours
}

func newEcgi(id types.EcID, plmnID types.PlmnID) types.ECGI {
	return types.ECGI{EcID: id, PlmnID: plmnID}
}
