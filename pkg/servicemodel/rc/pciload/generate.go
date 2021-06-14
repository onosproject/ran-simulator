// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package pciload

import (
	"fmt"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"math/rand"
)

// GeneratePCIMetrics generates semi-random PCI metrics for the specified model
func GeneratePCIMetrics(m *model.Model, minPCI uint, maxPCI uint, maxCollisions uint, earfcnStart uint32, sizeTypes []string) *PciMetrics {
	metrics := &PciMetrics{
		Cells: make(map[types.NCGI]PciCell),
	}

	// Generate PCI pools and shuffle them
	pools := generatePools(uint(2*len(m.Cells)), minPCI, maxPCI)
	indexes := rand.Perm(len(pools))

	// Create PCI metrics for each cell using unique PCI pools; so without collisions
	earfcn := earfcnStart
	pi := 0

	// Prepare to index by NCGI and cell name alike
	ecgis := make([]types.NCGI, 0, len(m.Cells))
	names := make([]string, 0, len(m.Cells))

	for name, cell := range m.Cells {
		// Assign each cell up to two PCI pools
		ranges := make([]PciRange, 0, 2)
		for i := rand.Intn(2); i >= 0; i-- {
			ranges = append(ranges, pools[indexes[pi]])
			pi = pi + 1
		}

		// Create metrics for each cell
		metrics.Cells[cell.NCGI] = PciCell{
			CellSize: sizeTypes[rand.Intn(len(sizeTypes))],
			Earfcn:   earfcn,
			Pci:      pickPCI(ranges),
			PciPool:  ranges,
		}
		earfcn = earfcn + 1

		ecgis = append(ecgis, cell.NCGI)
		names = append(names, name)
	}

	// Now inject requested number of collisions; between neighbour cells
	// Shuffle the cells so that conflict assignment is somewhat random
	conflicts := make(map[types.NCGI]PciCell)
	cellIndexes := rand.Perm(len(ecgis))
	collisions := uint(0)
	for i := 0; i < len(cellIndexes) && collisions < maxCollisions; i++ {
		ncgi := ecgis[cellIndexes[i]]

		if _, conflicted := conflicts[ncgi]; !conflicted {
			pciCell := metrics.Cells[ncgi]
			cellRanges := pciCell.PciPool
			cell := m.Cells[names[i]]

			necgi := cell.Neighbors[rand.Intn(len(cell.Neighbors))]
			if _, conflicted := conflicts[necgi]; !conflicted {
				neighborPciCell := metrics.Cells[necgi]

				// Replace the first PCI pool of the neighbor with a randomly chosen one from the cell
				neighborPciCell.PciPool[0] = cellRanges[rand.Intn(len(cellRanges))]
				neighborPciCell.Pci = pciCell.Pci

				metrics.Cells[necgi] = neighborPciCell
				collisions = collisions + 1

				conflicts[ncgi] = pciCell
				conflicts[necgi] = neighborPciCell

				fmt.Printf("Injected conflict between %d and %d\n", ncgi, necgi)
			}
		}
	}

	return metrics
}

func pickPCI(ranges []PciRange) uint32 {
	pi := rand.Intn(len(ranges))
	return ranges[pi].Min + uint32(rand.Intn(int(ranges[pi].Max-ranges[pi].Min)))
}

// Generate 2 x cell count number of pools evenly split between min and max
func generatePools(poolCount uint, minPCI uint, maxPCI uint) []PciRange {
	pools := make([]PciRange, 0, poolCount)
	poolSize := uint32((maxPCI - minPCI) / poolCount)

	pci := uint32(minPCI)
	for i := uint(0); i < poolCount; i++ {
		pools = append(pools, PciRange{Min: pci, Max: pci + poolSize - 1})
		pci = pci + poolSize
	}
	return pools
}
