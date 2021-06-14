// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/utils/honeycomb"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// A simple too to generate a tower configuration in a honeycomb layout
func main() {
	rootCmd := getRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "honeycomb",
		Short: "honeycomb RAN topology generator",
	}
	cmd.AddCommand(getHoneycombTopoCommand())
	return cmd
}

func getHoneycombTopoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "topo outfile",
		Short:         "ran-simulator config generation tool for onos-topo",
		SilenceUsage:  false,
		SilenceErrors: false,
		Args:          cobra.ExactArgs(1),
		RunE:          runHoneycombTopoCommand,
	}
	cmd.Flags().UintP("towers", "t", 0, "number of towers")
	_ = cmd.MarkFlagRequired("towers")
	cmd.Flags().UintP("sectors-per-tower", "s", 3, "sectors per tower")
	cmd.Flags().Float64P("latitude", "a", 52.5200, "Map centre latitude in degrees")
	cmd.Flags().Float64P("longitude", "g", 13.4050, "Map centre longitude in degrees")
	cmd.Flags().Float64P("max-neighbor-distance", "d", 3600.0, "Maximum 'distance' between neighbor cells; see docs")
	cmd.Flags().Int("max-neighbors", 5, "Maximum number of neighbors a cell will have; -1 no limit")
	cmd.Flags().StringSlice("service-models", []string{"kpm/1", "ni/2", "rc/3"}, "List of service models supported by the nodes")
	cmd.Flags().StringSlice("controller-addresses", []string{"onos-e2t"}, "List of E2T controller addresses or service names")
	cmd.Flags().String("plmnid", "315010", "PlmnID in MCC-MNC format, e.g. CCCNNN or CCCNN")
	cmd.Flags().Uint32P("enbidstart", "e", 5152, "GnbID start")
	cmd.Flags().Float32P("pitch", "i", 0.02, "pitch between cells in degrees")
	cmd.Flags().Bool("single-node", false, "generate a single node for all cells")
	cmd.Flags().String("controller-yaml", "", "if specified, location of yaml file for controller")
	cmd.Flags().Uint("min-pci", 0, "minimum PCI value")
	cmd.Flags().Uint("max-pci", 503, "maximum PCI value")
	cmd.Flags().Uint("max-collisions", 8, "maximum number of collisions")
	cmd.Flags().Uint32("earfcn-start", 42, "start point for EARFCN generation")
	cmd.Flags().StringSlice("cell-types", []string{"FEMTO", "ENTERPRISE", "OUTDOOR_SMALL", "MACRO"}, "List of cell size types")
	return cmd
}

func runHoneycombTopoCommand(cmd *cobra.Command, args []string) error {
	numTowers, _ := cmd.Flags().GetUint("towers")
	sectorsPerTower, _ := cmd.Flags().GetUint("sectors-per-tower")
	latitude, _ := cmd.Flags().GetFloat64("latitude")
	longitude, _ := cmd.Flags().GetFloat64("longitude")
	plmnid, _ := cmd.Flags().GetString("plmnid")
	enbidStart, _ := cmd.Flags().GetUint32("enbidstart")
	pitch, _ := cmd.Flags().GetFloat32("pitch")
	maxDistance, _ := cmd.Flags().GetFloat64("max-neighbor-distance")
	maxNeighbors, _ := cmd.Flags().GetInt("max-neighbors")
	controllerAddresses, _ := cmd.Flags().GetStringSlice("controller-addresses")
	serviceModels, _ := cmd.Flags().GetStringSlice("service-models")
	singleNode, _ := cmd.Flags().GetBool("single-node")
	controllerFile, _ := cmd.Flags().GetString("controller-yaml")

	minPci, _ := cmd.Flags().GetUint("min-pci")
	maxPci, _ := cmd.Flags().GetUint("max-pci")
	maxCollisions, _ := cmd.Flags().GetUint("max-collisions")
	earfcnStart, _ := cmd.Flags().GetUint32("earfcn-start")
	cellTypes, _ := cmd.Flags().GetStringSlice("cell-types")

	fmt.Printf("Creating honeycomb array of %d towers with %d cells each.\n", numTowers, sectorsPerTower)

	mapCenter := model.Coordinate{Lat: latitude, Lng: longitude}

	m, err := honeycomb.GenerateHoneycombTopology(mapCenter, numTowers, sectorsPerTower,
		types.PlmnIDFromString(plmnid), enbidStart, pitch, maxDistance, maxNeighbors, controllerAddresses, serviceModels, singleNode, minPci, maxPci, maxCollisions, earfcnStart, cellTypes)
	if err != nil {
		return err
	}

	m.Plmn = plmnid // we want the MCC-MNC format in our YAML

	d, err := yaml.Marshal(&m)
	if err != nil {
		fmt.Printf("Unable to marshal model data: %v", err)
		return err
	}
	if controllerFile != "" {
		writeControllerYaml(*m, controllerFile)
	}
	return ioutil.WriteFile(args[0], d, 0644)
}

func writeControllerYaml(model model.Model, location string) {
	f, err := os.OpenFile(location, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	w := bufio.NewWriter(f)

	// print out Nodes and their connections
	for node := range model.Nodes {
		// print the node
		_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nmetadata:\n  kind: Entity\n")
		_, _ = w.WriteString(fmt.Sprintf("  name: %d\n", model.Nodes[node].GnbID))
		_, _ = w.WriteString("spec:\n  aspects:\n    servicemodels:\n")
		// and print service models
		for serviceModel := range model.Nodes[node].ServiceModels {
			_, _ = w.WriteString(fmt.Sprintf("      - %s\n", model.Nodes[node].ServiceModels[serviceModel]))
		}
		_, _ = w.WriteString("---\n")
		// then print the connections separately
		for cell := range model.Nodes[node].Cells {
			_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nkind: Entity\nmetadata:\n  name: e2-node-cell\n")
			_, _ = w.WriteString("spec:\n  aspects:\n")
			_, _ = w.WriteString(fmt.Sprintf("    nodeid: %d\n", model.Nodes[node].GnbID))
			_, _ = w.WriteString(fmt.Sprintf("    cellid: %d\n", model.Nodes[node].Cells[cell]))
			_, _ = w.WriteString("---\n")
		}
	}

	// print out Cells
	for cell := range model.Cells {
		// print the cell
		_, _ = w.WriteString("apiVersion: topo.onosproject.org/v1beta1\nmetadata:\n  kind: Entity\n")
		_, _ = w.WriteString(fmt.Sprintf("  name: %d\n", model.Cells[cell].NCGI))
		_, _ = w.WriteString("spec:\n  aspects:\n")
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

	w.Flush()
}
