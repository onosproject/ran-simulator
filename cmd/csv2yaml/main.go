// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	arc     = 120
	maxUEs  = 99999
	plmnID  = 62831
	stopped = "stopped"
	txPower = 11.0
)

type importData struct {
	model *model.Model
	pci   *pciData
}

type pciData struct {
	Pcis map[types.ECGI]int32 `yaml:"pcis"`
}

func main() {
	// TODO: Add usage output and better argument validation
	cellsInput := os.Args[1]
	neighborsInput := os.Args[2]
	modelOutput := os.Args[3]
	pciOutput := os.Args[4]

	cellRecords := readCsvFile(cellsInput)
	neighborRecords := readCsvFile(neighborsInput)

	data := &importData{
		model: &model.Model{PlmnID: plmnID},
		pci:   &pciData{Pcis: make(map[types.ECGI]int32)},
	}

	loadModel(data, cellRecords, neighborRecords)
	outputModelYAML(data.model, modelOutput)
	outputPciYAML(data.pci, pciOutput)
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func loadModel(data *importData, cellRecords [][]string, neighborRecords [][]string) {
	processCellRecords(data, cellRecords)
	processNeighborRecords(data.model, neighborRecords)
}

func outputModelYAML(m *model.Model, path string) {
	bytes, err := yaml.Marshal(m)
	if err != nil {
		log.Fatal("Unable to output model as YAML", err)
	}
	err = ioutil.WriteFile(path, bytes, os.FileMode(0644))
	if err != nil {
		log.Fatal("Unable to write YAML file", err)
	}
}

func outputPciYAML(pci *pciData, path string) {
	bytes, err := yaml.Marshal(pci)
	if err != nil {
		log.Fatal("Unable to output PCI as YAML", err)
	}
	err = ioutil.WriteFile(path, bytes, os.FileMode(0644))
	if err != nil {
		log.Fatal("Unable to write YAML file", err)
	}
}

func processCellRecords(data *importData, cellRecords [][]string) {
	m := data.model
	m.Nodes = make(map[string]model.Node)
	m.Cells = make(map[string]model.Cell)

	controllers := []string{"controller1"}
	serviceModels := []string{"kpm", "rc"}

	lastLoc := model.Coordinate{Lat: 0, Lng: 0}

	ni := int64(0)
	node := model.Node{}
	for _, cr := range cellRecords {
		ci, err := strconv.ParseInt(cr[0], 10, 32)
		if err != nil {
			log.Fatal("Invalid cell ID", err)
		}

		lat, err := strconv.ParseFloat(cr[1], 64)
		if err != nil {
			log.Fatal("Invalid latitude", err)
		}
		lng, err := strconv.ParseFloat(cr[2], 64)
		if err != nil {
			log.Fatal("Invalid longitude", err)
		}
		azimuth, err := strconv.ParseInt(cr[3], 10, 32)
		if err != nil {
			log.Fatal("Invalid azimuth", err)
		}

		pci, err := strconv.ParseInt(cr[4], 10, 32)
		if err != nil {
			log.Fatal("Invalid PCI", err)
		}

		loc := model.Coordinate{Lat: lat, Lng: lng}

		// Create and register a new node each time location changes
		if loc.Lat != lastLoc.Lat || loc.Lng != lastLoc.Lng {
			ni = ni + 1
			node = model.Node{
				EnbID:         genEnbID(ni),
				Controllers:   controllers,
				ServiceModels: serviceModels,
				Cells:         make([]types.ECGI, 0, 3),
				Status:        stopped,
			}
		}

		// Create and register a new cell
		cell := model.Cell{
			ECGI:      genECGI(node.EnbID, ci),
			Sector:    model.Sector{Center: loc, Arc: arc, Azimuth: int32(azimuth)},
			Color:     "none",
			MaxUEs:    maxUEs,
			TxPowerDB: txPower,
		}
		m.Cells[cellName(ci)] = cell

		// Update the PCI structure
		data.pci.Pcis[cell.ECGI] = int32(pci)

		// Associate the new cell with the current node and update the node in the model Nodes
		node.Cells = append(node.Cells, cell.ECGI)
		m.Nodes[nodeName(ni)] = node

		lastLoc = loc
	}
}

func processNeighborRecords(m *model.Model, neighborRecords [][]string) {
	for i, nr := range neighborRecords {
		// FIXME: input data should contain explicit index; not assume relative line position
		ci := int64(i) + 1
		//ci, err := strconv.ParseInt(nr[0], 10, 32)
		//if err != nil {
		//	log.Fatal("Invalid cell ID", err)
		//}
		name := cellName(ci)
		cell := m.Cells[name]
		cell.Neighbors = make([]types.ECGI, 0, len(nr))

		for _, nid := range nr {
			// Lookup neighbor cells to validate data coherence; even though it's more expensive
			nci, err := strconv.ParseInt(nid, 10, 32)
			if err != nil {
				log.Fatal("Invalid neighbor cell ID", err)
			}
			cell.Neighbors = append(cell.Neighbors, m.Cells[cellName(nci)].ECGI)
		}

		m.Cells[name] = cell
	}
}

func genEnbID(ci int64) types.EnbID {
	return types.EnbID(ci)
}
func genECGI(enbid types.EnbID, ci int64) types.ECGI {
	return types.ToECGI(plmnID, types.ToECI(enbid, types.CellID(ci)))
}

func nodeName(i int64) string {
	return fmt.Sprintf("node%d", i)
}

func cellName(i int64) string {
	return fmt.Sprintf("cell%d", i)
}
