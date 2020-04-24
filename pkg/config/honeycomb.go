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

package config

import (
	"fmt"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/pmcxs/hexgrid"
)

// HoneycombGenerator - used by the cli tool "honeycomb"
func HoneycombGenerator(numTowers uint, sectorsPerTower uint, latitude float32,
	longitude float32, plmnid types.PlmnID, ecidStart uint16, portstart uint16, pitch float32) (*TowerConfig, error) {

	mapCentre := types.Point{
		Lat: latitude,
		Lng: longitude,
	}

	aspectRatio := float32(utils.AspectRatio(&mapCentre))
	newConfig := TowerConfig{
		MapCentre:    mapCentre,
		TowersLayout: make([]TowersLayout, numTowers),
	}
	points := hexMesh(float64(pitch), numTowers)
	var t, s uint
	for t = 0; t < numTowers; t++ {
		tower := TowersLayout{
			TowerID:   fmt.Sprintf("Tower-%d", t+1),
			PlmnID:    plmnid,
			Latitude:  latitude + points[t].Lat,
			Longitude: longitude + points[t].Lng/aspectRatio,
			Sectors:   make([]Sector, sectorsPerTower),
		}
		for s = 0; s < sectorsPerTower; s++ {
			var azimuth uint = 0
			if s > 0 {
				azimuth = 360.0 * s / sectorsPerTower
			}
			sector := Sector{
				EcID:        types.EcID(fmt.Sprintf("%07x", ecidStart+uint16(t*sectorsPerTower)+uint16(s))),
				GrpcPort:    portstart + uint16(t*sectorsPerTower) + uint16(s),
				Azimuth:     uint16(azimuth),
				Arc:         360.0 / uint16(sectorsPerTower),
				MaxUEs:      5,
				InitPowerDb: 10,
			}
			tower.Sectors[s] = sector
		}
		newConfig.TowersLayout[t] = tower
	}
	return &newConfig, nil
}

func hexMesh(pitch float64, numTowers uint) []*types.Point {
	rings, _ := numRings(numTowers)
	points := make([]*types.Point, 0)
	hexArray := hexgrid.HexRange(hexgrid.NewHex(0, 0), int(rings))

	for _, h := range hexArray {
		x, y := hexgrid.Point(hexgrid.HexToPixel(hexgrid.LayoutPointY00(pitch, pitch), h))
		points = append(points, &types.Point{
			Lat: float32(x),
			Lng: float32(y),
		})
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
