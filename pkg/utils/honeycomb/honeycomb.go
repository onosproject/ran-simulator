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
	"strconv"
	"strings"
)

// GenerateHoneycombTopology generates a set of simulated nodes and cells organized in a honeycomb
// outward from the specified center.
func GenerateHoneycombTopology(mapCenter model.Coordinate, numTowers uint, sectorsPerTower uint, plmnID types.PlmnID,
	enbStart uint32, pitch float32, maxDistance float64, maxNeighbors int,
	controllerAddresses []string, serviceModels []string, singleNode bool) (*model.Model, error) {

	m := &model.Model{
		PlmnID:        plmnID,
		MapLayout:     model.MapLayout{Center: mapCenter, LocationsScale: 1.25},
		Cells:         make(map[string]model.Cell),
		Nodes:         make(map[string]model.Node),
		Controllers:   generateControllers(controllerAddresses),
		ServiceModels: generateServiceModels(serviceModels),
	}

	aspectRatio := utils.AspectRatio(mapCenter.Lat)
	points := hexMesh(float64(pitch), numTowers)
	arc := int32(360.0 / sectorsPerTower)

	controllers := make([]string, 0, len(controllerAddresses))
	for name := range m.Controllers {
		controllers = append(controllers, name)
	}

	models := make([]string, 0, len(serviceModels))
	for name := range m.ServiceModels {
		models = append(models, name)
	}

	var t, s uint
	var enbID types.EnbID
	var nodeName string
	var node model.Node
	for t = 0; t < numTowers; t++ {
		var azOffset int32 = 0
		if sectorsPerTower == 6 {
			azOffset = int32(math.Mod(float64(t), 2) * 30)
		}

		if !singleNode || t == 0 {
			enbID = types.EnbID(enbStart + uint32(t+1))
			nodeName = fmt.Sprintf("node%d", t+1)

			node = model.Node{
				EnbID:         enbID,
				Controllers:   controllers,
				ServiceModels: models,
				Cells:         make([]types.ECGI, 0, sectorsPerTower),
				Status:        "stopped",
			}
		}

		for s = 0; s < sectorsPerTower; s++ {
			cellID := types.CellID(s + 1)
			if singleNode && sectorsPerTower == 1 {
				cellID = types.CellID(t + 1)
			}
			cellName := fmt.Sprintf("cell%d", (t*sectorsPerTower)+s+1)

			azimuth := azOffset
			if s > 0 {
				azimuth = int32(360.0*s/sectorsPerTower + uint(azOffset))
			}

			cell := model.Cell{
				ECGI: types.ToECGI(plmnID, types.ToECI(enbID, cellID)),
				Sector: model.Sector{
					Center: model.Coordinate{
						Lat: mapCenter.Lat + points[t].Lat,
						Lng: mapCenter.Lng + points[t].Lng/aspectRatio},
					Azimuth: azimuth,
					Arc:     arc},
				Color:     "green",
				MaxUEs:    99999,
				Neighbors: make([]types.ECGI, 0, sectorsPerTower),
				TxPowerDB: 11,
			}

			m.Cells[cellName] = cell
			node.Cells = append(node.Cells, cell.ECGI)
		}

		m.Nodes[nodeName] = node
	}

	// Add cells neighbors
	for cellName, cell := range m.Cells {
		for _, other := range m.Cells {
			if cell.ECGI != other.ECGI && isNeighbor(cell, other, maxDistance, sectorsPerTower == 1) && len(cell.Neighbors) < maxNeighbors {
				cell.Neighbors = append(cell.Neighbors, other.ECGI)
			}
		}
		m.Cells[cellName] = cell
	}

	return m, nil
}

func generateControllers(addresses []string) map[string]model.Controller {
	controllers := make(map[string]model.Controller)
	for i, address := range addresses {
		name := fmt.Sprintf("e2t-%d", i+1)
		controllers[name] = model.Controller{ID: name, Address: address, Port: 36421}
	}
	return controllers
}

func generateServiceModels(namesAndIDs []string) map[string]model.ServiceModel {
	models := make(map[string]model.ServiceModel)
	for i, nameAndID := range namesAndIDs {
		fields := strings.Split(nameAndID, "/")
		id := int64(i)
		if len(fields) > 1 {
			id, _ = strconv.ParseInt(fields[1], 10, 32)
		}
		models[fields[0]] = model.ServiceModel{ID: int(id), Version: "1.0.0", Description: fields[0] + " service model"}
	}
	return models
}

// Cells are neighbors if their sectors have the same coordinates or if their center arc vectors fall within a distance/2
func isNeighbor(cell model.Cell, other model.Cell, maxDistance float64, onlyDistance bool) bool {
	return (cell.Sector.Center.Lat == other.Sector.Center.Lat && cell.Sector.Center.Lng == other.Sector.Center.Lng) ||
		(onlyDistance && utils.Distance(cell.Sector.Center, other.Sector.Center) <= maxDistance) ||
		utils.Distance(reachPoint(cell.Sector, maxDistance), reachPoint(other.Sector, maxDistance)) <= maxDistance/2
}

// Calculate the end-point of the center arc vector a distance from the sector center
func reachPoint(sector model.Sector, distance float64) model.Coordinate {
	return utils.TargetPoint(sector.Center, float64((sector.Azimuth+sector.Arc/2)%360), distance)
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
