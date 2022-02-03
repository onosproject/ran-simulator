// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package honeycomb

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"os"
)

// WriteControllerYaml outputs YAML file that can be consumed by the onos topo operator.
func WriteControllerYaml(model model.Model, location string) error {
	f, err := os.OpenFile(location, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	// print out Nodes and their connections
	for _, node := range model.Nodes {
		printNode(w, node)
	}

	// print out Cells
	for _, cell := range model.Cells {
		printCell(w, cell)

		for _, neighbor := range cell.Neighbors {
			printCellNeighbor(w, cell, neighbor)
		}
	}

	// print the node-cell relations separately
	for _, node := range model.Nodes {
		for _, ncgi := range node.Cells {
			printNodeCellRelation(w, node, ncgi)
		}
	}

	return w.Flush()
}

func printNode(w *bufio.Writer, node model.Node) {
	_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nkind: Entity\nmetadata:\n")
	_, _ = w.WriteString(fmt.Sprintf("  name: e2.1.%x\n", node.GnbID))
	_, _ = w.WriteString("spec:\n")
	_, _ = w.WriteString(fmt.Sprintf("  uri: e2:1/%x\n", node.GnbID))
	_, _ = w.WriteString("  kind:\n    name: e2node\n")
	_, _ = w.WriteString("  aspects:\n")
	_, _ = w.WriteString("    onos.topo.E2Node:\n")
	_, _ = w.WriteString("      service_models:\n")
	_, _ = w.WriteString("---\n")
}

func printCell(w *bufio.Writer, cell model.Cell) {
	_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nkind: Entity\nmetadata:\n")
	_, _ = w.WriteString(fmt.Sprintf("  name: e2.1.%x.%x\n", types.GetGnbID(uint64(cell.NCGI)), types.GetCellID(uint64(cell.NCGI))))
	_, _ = w.WriteString("spec:\n")
	_, _ = w.WriteString(fmt.Sprintf("  uri: e2:1/%x/%x\n", types.GetGnbID(uint64(cell.NCGI)), types.GetCellID(uint64(cell.NCGI))))
	_, _ = w.WriteString("  kind:\n    name: e2cell\n")
	_, _ = w.WriteString("  aspects:\n")
	_, _ = w.WriteString("    onos.topo.Location:\n")
	_, _ = w.WriteString(fmt.Sprintf("      lat: %f\n", cell.Sector.Center.Lat))
	_, _ = w.WriteString(fmt.Sprintf("      lng: %f\n", cell.Sector.Center.Lng))
	_, _ = w.WriteString("    onos.topo.E2Cell:\n")
	_, _ = w.WriteString(fmt.Sprintf("      earfcn: %d\n", cell.Earfcn))
	_, _ = w.WriteString(fmt.Sprintf("      cell_type: %s\n", cell.CellType.String()))
	_, _ = w.WriteString("    onos.topo.Coverage:\n")
	_, _ = w.WriteString(fmt.Sprintf("      arc_width: %d\n", cell.Sector.Arc))
	_, _ = w.WriteString(fmt.Sprintf("      tilt: %d\n", cell.Sector.Tilt))
	_, _ = w.WriteString(fmt.Sprintf("      height: %d\n", cell.Sector.Height))
	_, _ = w.WriteString(fmt.Sprintf("      azimuth: %d\n", cell.Sector.Azimuth))
	_, _ = w.WriteString("---\n")
}

func printNodeCellRelation(w *bufio.Writer, node model.Node, ncgi types.NCGI) {
	rid, _ := uuid.NewRandom()
	_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nkind: Relation\nmetadata:\n")
	_, _ = w.WriteString(fmt.Sprintf("  name: rid.%s\n", rid.String()))
	_, _ = w.WriteString("spec:\n")
	_, _ = w.WriteString(fmt.Sprintf("  uri: rid:%s\n", rid.String()))
	_, _ = w.WriteString("  kind:\n    name: contains\n")
	_, _ = w.WriteString(fmt.Sprintf("  source:\n    uri: e2:1/%x\n", node.GnbID))
	_, _ = w.WriteString(fmt.Sprintf("  target:\n    uri: e2:1/%x/%x\n", node.GnbID, types.GetCellID(uint64(ncgi))))
	_, _ = w.WriteString("---\n")
}

func printCellNeighbor(w *bufio.Writer, cell model.Cell, neighbor types.NCGI) {
	rid, _ := uuid.NewRandom()
	_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nkind: Relation\nmetadata:\n")
	_, _ = w.WriteString(fmt.Sprintf("  name: rid.%s\n", rid.String()))
	_, _ = w.WriteString("spec:\n")
	_, _ = w.WriteString(fmt.Sprintf("  uri: rid:%s\n", rid.String()))
	_, _ = w.WriteString("  kind:\n    name: neighbors\n")
	_, _ = w.WriteString(fmt.Sprintf("  source:\n    uri: e2:1/%x/%x\n", types.GetGnbID(uint64(cell.NCGI)), types.GetCellID(uint64(cell.NCGI))))
	_, _ = w.WriteString(fmt.Sprintf("  target:\n    uri: e2:1/%x/%x\n", types.GetGnbID(uint64(neighbor)), types.GetCellID(uint64(neighbor))))
	_, _ = w.WriteString("---\n")
}
