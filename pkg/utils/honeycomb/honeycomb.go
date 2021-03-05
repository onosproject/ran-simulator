// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package honeycomb

import (
	"fmt"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/pmcxs/hexgrid"
	"math"
)

// GenerateHoneycombTopology generates a set of simulated nodes and cells organized in a honeycomb
// outward from the specified center.
func GenerateHoneycombTopology(mapCenter model.Coordinate, numTowers uint, sectorsPerTower uint,
	plmnID types.PlmnID, enbStart uint32, pitch float32, maxDistance float64) (*model.Model, error) {

	m := &model.Model{
		PlmnID:    plmnID,
		MapLayout: model.MapLayout{Center: mapCenter},
		Cells:     make(map[string]model.Cell),
		Nodes:     make(map[string]model.Node),
	}

	aspectRatio := utils.AspectRatio(mapCenter.Lat)
	points := hexMesh(float64(pitch), numTowers)
	arc := int32(360.0 / sectorsPerTower)

	var t, s uint
	for t = 0; t < numTowers; t++ {
		var azOffset int32 = 0
		if sectorsPerTower == 6 {
			azOffset = int32(math.Mod(float64(t), 2) * 30)
		}

		enbID := types.EnbID(enbStart + uint32(t+1))
		nodeName := fmt.Sprintf("node%d", t+1)

		node := model.Node{
			EnbID:         enbID,
			Controllers:   []string{"e2t-1"},
			ServiceModels: []string{"kpm", "rc"},
			Cells:         make([]types.ECGI, 0, sectorsPerTower),
			Status:        "stopped",
		}

		m.Nodes[nodeName] = node

		for s = 0; s < sectorsPerTower; s++ {
			cellID := types.CellID(s + 1)
			cellName := fmt.Sprintf("cell%d", (t*sectorsPerTower)+s+1)

			azimuth := azOffset
			if s > 0 {
				azimuth = int32(360.0*s/sectorsPerTower + uint(azOffset))
			}

			cell := model.Cell{
				ECGI: types.ToECGI(plmnID, types.ToECI(node.EnbID, cellID)),
				Sector: model.Sector{
					Center: model.Coordinate{
						Lat: mapCenter.Lat + points[t].Lat,
						Lng: mapCenter.Lng + points[t].Lng/aspectRatio},
					Azimuth: azimuth,
					Arc:     arc},
				Color:     "green",
				MaxUEs:    99999,
				Neighbors: nil,
				TxPowerDB: 11,
			}

			m.Cells[cellName] = cell
			node.Cells = append(node.Cells, cell.ECGI)
		}

	}

	// Add cells neighbors
	for cellName, cell := range m.Cells {
		for _, other := range m.Cells {
			if isNeighbor(cell, other, maxDistance) {
				cell.Neighbors = append(cell.Neighbors, other.ECGI)
			}
		}
		m.Cells[cellName] = cell
	}

	return m, nil
}

func isNeighbor(cell model.Cell, other model.Cell, maxDistance float64) bool {
	return (cell.Sector.Center.Lat == other.Sector.Center.Lat && cell.Sector.Center.Lng == other.Sector.Center.Lng) ||
		distance(cell.Sector.Center, other.Sector.Center) <= maxDistance
}

// http://en.wikipedia.org/wiki/Haversine_formula
func distance(c1 model.Coordinate, c2 model.Coordinate) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = c1.Lat * math.Pi / 180
	lo1 = c1.Lng * math.Pi / 180
	la2 = c2.Lat * math.Pi / 180
	lo2 = c2.Lng * math.Pi / 180

	r = 6378100 // Earth radius in meters

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func hexMesh(pitch float64, numTowers uint) []*model.Coordinate {
	rings, _ := numRings(numTowers)
	points := make([]*model.Coordinate, 0)
	hexArray := hexgrid.HexRange(hexgrid.NewHex(0, 0), int(rings))

	for _, h := range hexArray {
		x, y := hexgrid.Point(hexgrid.HexToPixel(hexgrid.LayoutPointY00(pitch, pitch), h))
		points = append(points, &model.Coordinate{Lat: x, Lng: y})
	}
	return points
}

// Number of cells in the hexagon layout 3x^2+9x+7
func numRings(numTowers uint) (uint, error) {
	switch n := numTowers; {
	case n <= 7:
		return 1, nil
	case n <= 19:
		return 2, nil
	case n <= 37:
		return 3, nil
	case n <= 61:
		return 4, nil
	case n <= 91:
		return 5, nil
	case n <= 127:
		return 6, nil
	case n <= 169:
		return 7, nil
	case n <= 217:
		return 8, nil
	case n <= 271:
		return 9, nil
	case n <= 331:
		return 10, nil
	case n <= 469:
		return 11, nil
	default:
		return 0, fmt.Errorf(">469 not handled %d", numTowers)
	}

}
