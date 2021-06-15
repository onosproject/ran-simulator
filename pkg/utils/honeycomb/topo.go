// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package honeycomb

import (
	"bufio"
	"fmt"
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
	for node := range model.Nodes {
		// print the node
		_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nkind: Entity\nmetadata:\n")
		_, _ = w.WriteString(fmt.Sprintf("  name: \"%d\"\n", model.Nodes[node].GnbID))
		_, _ = w.WriteString("spec:\n")
		_, _ = w.WriteString("  kind:\n    name: e2-node\n")
		_, _ = w.WriteString("  aspects:\n")
		_, _ = w.WriteString("    onos.topo.E2Node:\n")
		_, _ = w.WriteString("      service_models:\n")
		_, _ = w.WriteString("---\n")
	}

	// print out Cells
	for cell := range model.Cells {
		// print the cell
		_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nkind: Entity\nmetadata:\n")
		_, _ = w.WriteString(fmt.Sprintf("  name: \"%d\"\n", model.Cells[cell].NCGI))
		_, _ = w.WriteString("spec:\n")
		_, _ = w.WriteString("  kind:\n    name: e-cell\n")
		_, _ = w.WriteString("  aspects:\n")
		_, _ = w.WriteString("    onos.topo.Location:\n")
		_, _ = w.WriteString(fmt.Sprintf("      lat: %f\n", model.Cells[cell].Sector.Center.Lat))
		_, _ = w.WriteString(fmt.Sprintf("      lng: %f\n", model.Cells[cell].Sector.Center.Lng))
		_, _ = w.WriteString("    onos.topo.E2Cell:\n")
		_, _ = w.WriteString(fmt.Sprintf("      earfcn: %d\n", model.Cells[cell].Earfcn))
		_, _ = w.WriteString(fmt.Sprintf("      cell_type: %s\n", model.Cells[cell].CellType.String()))
		_, _ = w.WriteString("    onos.topo.Coverage:\n")
		_, _ = w.WriteString(fmt.Sprintf("      arc_width: %d\n", model.Cells[cell].Sector.Arc))
		_, _ = w.WriteString(fmt.Sprintf("      tilt: %d\n", model.Cells[cell].Sector.Tilt))
		_, _ = w.WriteString(fmt.Sprintf("      height: %d\n", model.Cells[cell].Sector.Height))
		_, _ = w.WriteString(fmt.Sprintf("      azimuth: %d\n", model.Cells[cell].Sector.Azimuth))
		_, _ = w.WriteString("---\n")
	}

	// print the node-cell relations separately
	for _, node := range model.Nodes {
		for _, ncgi := range node.Cells {
			_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nkind: Relation\nmetadata:\n")
			_, _ = w.WriteString(fmt.Sprintf("  name: \"%d-%d\"\n", node.GnbID, ncgi))
			_, _ = w.WriteString("spec:\n")
			_, _ = w.WriteString("  kind:\n    name: e-node-cell\n")
			_, _ = w.WriteString(fmt.Sprintf("  source:\n    name: \"%d\"\n", node.GnbID))
			_, _ = w.WriteString(fmt.Sprintf("  target:\n    name: \"%d\"\n", ncgi))
			_, _ = w.WriteString("---\n")
		}
	}

	return w.Flush()
}
