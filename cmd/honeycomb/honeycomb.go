// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"fmt"
	"os"
	"strconv"

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
	cmd.Flags().Float64P("max-neighbor-distance", "d", 8000.0, "Maximum 'distance' between neighbor cells; see docs")
	cmd.Flags().Int("max-neighbors", 5, "Maximum number of neighbors a cell will have; -1 no limit")
	cmd.Flags().StringSlice("service-models", []string{"kpm/1", "rcpre2/3", "kpm2/4", "mho/5"}, "List of service models supported by the nodes")
	cmd.Flags().StringSlice("controller-addresses", []string{"onos-e2t"}, "List of E2T controller addresses or service names")
	cmd.Flags().String("plmnid", "315010", "PlmnID in MCC-MNC format, e.g. CCCNNN or CCCNN")
	cmd.Flags().Uint("ue-count", 0, "User Equipment count")
	cmd.Flags().Uint("ue-count-per-cell", 15, "Desired UE count per cell")
	cmd.Flags().String("gnbid-start", "5152", "GnbID start in hex")
	cmd.Flags().Float32P("pitch", "i", 0.02, "pitch between cells in degrees")
	cmd.Flags().Bool("single-node", false, "generate a single node for all cells")
	cmd.Flags().String("controller-yaml", "", "if specified, location of yaml file for controller")
	cmd.Flags().Uint("min-pci", 0, "minimum PCI value")
	cmd.Flags().Uint("max-pci", 503, "maximum PCI value")
	cmd.Flags().Uint("max-collisions", 8, "maximum number of collisions")
	cmd.Flags().Uint32("earfcn-start", 42, "start point for EARFCN generation")
	cmd.Flags().StringSlice("cell-types", []string{"FEMTO", "ENTERPRISE", "OUTDOOR_SMALL", "MACRO"}, "List of cell size types")
	cmd.Flags().Float64("deform-scale", .01, "scale factor for perturbation")
	return cmd
}

func runHoneycombTopoCommand(cmd *cobra.Command, args []string) error {
	numTowers, _ := cmd.Flags().GetUint("towers")
	sectorsPerTower, _ := cmd.Flags().GetUint("sectors-per-tower")
	latitude, _ := cmd.Flags().GetFloat64("latitude")
	longitude, _ := cmd.Flags().GetFloat64("longitude")
	plmnid, _ := cmd.Flags().GetString("plmnid")
	ueCount, _ := cmd.Flags().GetUint("ue-count")
	ueCountPerCell, _ := cmd.Flags().GetUint("ue-count-per-cell")
	gnbidStartS, _ := cmd.Flags().GetString("gnbid-start")
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
	deformScale, _ := cmd.Flags().GetFloat64("deform-scale")

	gnbidStart, err := strconv.ParseUint(gnbidStartS, 16, 32)
	if err != nil {
		return err
	}

	fmt.Printf("Creating honeycomb array of %d towers with %d cells each.\n", numTowers, sectorsPerTower)

	mapCenter := model.Coordinate{Lat: latitude, Lng: longitude}

	m, err := honeycomb.GenerateHoneycombTopology(mapCenter, numTowers, sectorsPerTower,
		types.PlmnIDFromString(plmnid), uint32(gnbidStart), pitch, maxDistance, maxNeighbors,
		controllerAddresses, serviceModels, singleNode, minPci, maxPci, maxCollisions, earfcnStart, cellTypes, deformScale)
	if err != nil {
		return err
	}

	m.Plmn = plmnid // we want the MCC-MNC format in our YAML
	m.UECount = ueCount

	m.UECountPerCell = ueCountPerCell

	d, err := yaml.Marshal(&m)
	if err != nil {
		fmt.Printf("Unable to marshal model data: %v", err)
		return err
	}
	if controllerFile != "" {
		err = honeycomb.WriteControllerYaml(*m, controllerFile)
		if err != nil {
			fmt.Printf("Unable to output topology operator file: %v", err)
			return err
		}
	}
	return os.WriteFile(args[0], d, 0644)
}
