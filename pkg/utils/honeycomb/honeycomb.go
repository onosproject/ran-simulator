// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package honeycomb

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/pmcxs/hexgrid"
)

// used for generating the pci
type auxPCI struct {
	pci     uint32
	pciPool []pciRange
}

type pciRange struct {
	min uint32
	max uint32
}

// GenerateHoneycombTopology generates a set of simulated nodes and cells organized in a honeycomb
// outward from the specified center.
func GenerateHoneycombTopology(mapCenter model.Coordinate, numTowers uint, sectorsPerTower uint, plmnID types.PlmnID,
	enbStart uint32, pitch float32, maxDistance float64, maxNeighbors int,
	controllerAddresses []string, serviceModels []string, singleNode bool, minPci uint, maxPci uint, maxCollisions uint, earfcnStart uint32, cellTypes []string, deformScale float64) (*model.Model, error) {

	earfcn := earfcnStart

	m := &model.Model{
		PlmnID:        plmnID,
		MapLayout:     model.MapLayout{Center: mapCenter, LocationsScale: 1.25},
		Cells:         make(map[string]model.Cell),
		Nodes:         make(map[string]model.Node),
		Controllers:   generateControllers(controllerAddresses),
		ServiceModels: generateServiceModels(serviceModels),
	}

	points := hexMesh(float64(pitch), numTowers, m.MapLayout.Center, deformScale)
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
	var gnbID types.GnbID
	var nodeName string
	var node model.Node
	for t = 0; t < numTowers; t++ {
		var azOffset int32
		if sectorsPerTower == 6 {
			azOffset = int32(math.Mod(float64(t), 2) * 30)
		}

		if !singleNode || t == 0 {
			gnbID = types.GnbID(enbStart + uint32(t+1))
			nodeName = fmt.Sprintf("node%d", t+1)

			node = model.Node{
				GnbID:         gnbID,
				Controllers:   controllers,
				ServiceModels: models,
				Cells:         make([]types.NCGI, 0, sectorsPerTower),
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
				NCGI: types.ToNCGI(plmnID, types.ToNCI(gnbID, cellID)),
				Sector: model.Sector{
					Center:  *points[t],
					Azimuth: azimuth,
					Arc:     arc,
					Height:  int32(rand.Intn(31) + 20),
					Tilt:    int32(rand.Intn(31) - 15)},
				Color:     "green",
				MaxUEs:    99999,
				Neighbors: make([]types.NCGI, 0, sectorsPerTower),
				TxPowerDB: 11,
				Earfcn:    earfcn,
			}
			earfcn++

			m.Cells[cellName] = cell
			node.Cells = append(node.Cells, cell.NCGI)
		}
		m.Nodes[nodeName] = node
	}

	// Add cells neighbors
	for cellName, cell := range m.Cells {
		for _, other := range m.Cells {
			if cell.NCGI != other.NCGI && isNeighbor(cell, other, maxDistance, sectorsPerTower == 1) && len(cell.Neighbors) < maxNeighbors {
				cell.Neighbors = append(cell.Neighbors, other.NCGI)
			}
		}
		m.Cells[cellName] = cell
	}
	// Add pci values
	generatePCI(m.Cells, minPci, maxPci, maxCollisions)
	// Add random cell type
	validCellTypes := make([]int32, 0)
	// validCellTypes := [4]uint32{0, 1, 2, 3}
	for cellType := range cellTypes {
		validCellTypes = append(validCellTypes, types.CellType_value[cellTypes[cellType]])
	}
	for name, cell := range m.Cells {
		tempCell := cell
		tempCell.CellType = types.CellType(rand.Intn(len(validCellTypes)))
		m.Cells[name] = tempCell
	}
	return m, nil
}

func generatePCI(cells map[string]model.Cell, minPCI uint, maxPCI uint, maxCollisions uint) {
	if len(cells) > int(maxPCI-minPCI) {
		panic("Too little space in between the minimum and maximum PCI values. Try setting --min-pci lower or --max-pci higher")
	}
	pciCells := make(map[types.NCGI]auxPCI)

	// Generate PCI pools and shuffle them
	pools := generatePools(uint(2*len(cells)), minPCI, maxPCI)
	indexes := rand.Perm(len(pools))

	// Create PCI metrics for each cell using unique PCI pools; so without collisions
	pi := 0

	// Prepare to index by NCGI and cell name alike
	ecgis := make([]types.NCGI, 0, len(cells))
	names := make([]string, 0, len(cells))

	for name, cell := range cells {
		// Assign each cell up to two PCI pools
		ranges := make([]pciRange, 0, 2)
		for i := rand.Intn(2); i >= 0; i-- {
			ranges = append(ranges, pools[indexes[pi]])
			pi = pi + 1
		}

		// Create metrics for each cell
		pciCells[cell.NCGI] = auxPCI{
			pci:     pickPCI(ranges),
			pciPool: ranges,
		}

		ecgis = append(ecgis, cell.NCGI)
		names = append(names, name)
	}

	// Now inject requested number of collisions; between neighbour cells
	// Shuffle the cells so that conflict assignment is somewhat random
	conflicts := make(map[types.NCGI]auxPCI)
	cellIndexes := rand.Perm(len(ecgis))
	collisions := uint(0)
	for i := 0; i < len(cellIndexes) && collisions < maxCollisions; i++ {
		ncgi := ecgis[cellIndexes[i]]

		if _, conflicted := conflicts[ncgi]; !conflicted {
			pciCell := pciCells[ncgi]
			cellRanges := pciCell.pciPool
			cell := cells[names[i]]

			necgi := cell.Neighbors[rand.Intn(len(cell.Neighbors))]
			if _, conflicted := conflicts[necgi]; !conflicted {
				neighborPciCell := pciCells[necgi]

				// Replace the first PCI pool of the neighbor with a randomly chosen one from the cell
				neighborPciCell.pciPool[0] = cellRanges[rand.Intn(len(cellRanges))]
				neighborPciCell.pci = pciCell.pci
				pciCells[necgi] = neighborPciCell
				collisions = collisions + 1

				conflicts[ncgi] = pciCell
				conflicts[necgi] = neighborPciCell

				fmt.Printf("Injected conflict between %d and %d\n", ncgi, necgi)
			}
		}
	}
	for i := 0; i < len(cellIndexes); i++ {
		tempCell := cells[names[i]]
		tempCell.PCI = pciCells[ecgis[i]].pci
		cells[names[i]] = tempCell
	}
}

func pickPCI(ranges []pciRange) uint32 {
	pi := rand.Intn(len(ranges))
	if ranges[pi].max-ranges[pi].min == 0 {
		return ranges[pi].min
	}
	return ranges[pi].min + uint32(rand.Intn(int(ranges[pi].max-ranges[pi].min)))
}

// Generate 2 x cell count number of pools evenly split between min and max
func generatePools(poolCount uint, minPCI uint, maxPCI uint) []pciRange {
	pools := make([]pciRange, 0, poolCount)
	poolSize := uint32((maxPCI - minPCI) / poolCount)

	pci := uint32(minPCI)
	for i := uint(0); i < poolCount; i++ {
		pools = append(pools, pciRange{min: pci, max: pci + poolSize - 1})
		pci = pci + poolSize
	}
	return pools
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

func hexMesh(pitch float64, numTowers uint, center model.Coordinate, deformScale float64) []*model.Coordinate {
	rings, _ := numRings(numTowers)
	points := make([]*model.Coordinate, 0)
	hexArray := hexgrid.HexRange(hexgrid.NewHex(0, 0), int(rings))
	// randomly generate a center point (will be biased towards poles). this is deterministic since go rand is deterministic
	theta := utils.DegreesToRads(center.Lat)
	phi := utils.DegreesToRads(center.Lng)
	// deform mesh
	for _, h := range hexArray {
		x, y := hexgrid.Point(hexgrid.HexToPixel(hexgrid.LayoutPointY00(pitch, pitch), h))
		// angle offset in radians
		x = utils.DegreesToRads(x + (rand.Float64()-0.5)*deformScale)
		y = utils.DegreesToRads(y + (rand.Float64()-0.5)*deformScale)
		// perturb each individual point
		lat := (math.Asin(math.Cos(theta)*math.Sin(x) + math.Cos(y)*math.Sin(theta)*math.Cos(x))) * 180 / math.Pi
		lon := (math.Atan2(math.Sin(y), -math.Tan(x)*math.Sin(theta)+math.Cos(y)*math.Cos(theta)) + phi) * 180 / math.Pi
		points = append(points, &model.Coordinate{Lat: lat, Lng: lon})
		// logging location
		// fmt.Printf("%f, %f\n", lat, lon)
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
