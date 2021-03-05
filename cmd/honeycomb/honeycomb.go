// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package main

import (
	"fmt"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/utils/honeycomb"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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
	cmd.Flags().UintP("sectors-per-tower", "s", 3, "sectors per tower (3 or 6)")
	cmd.Flags().Float64P("latitude", "a", 52.5200, "Map centre latitude in degrees")
	cmd.Flags().Float64P("longitude", "g", 13.4050, "Map centre longitude in degrees")
	cmd.Flags().Float64P("max-neighbor-distance", "d", 3860.0, "Maximum distance between neighbor cells")
	cmd.Flags().String("plmnid", "315010", "PlmnID in MCC-MNC format, e.g. CCCNNN or CCCNN")
	cmd.Flags().Uint32P("enbidstart", "e", 5152, "EnbID start")
	cmd.Flags().Float32P("pitch", "i", 0.02, "pitch between cells in degrees")
	return cmd
}

func runHoneycombTopoCommand(cmd *cobra.Command, args []string) error {
	numTowers, _ := cmd.Flags().GetUint("towers")
	sectorsPerTower, _ := cmd.Flags().GetUint("sectors-per-tower")
	if sectorsPerTower != 3 && sectorsPerTower != 6 {
		return fmt.Errorf("only 3 or 6 are allowed for 'sectors-per-tower'")
	}
	latitude, _ := cmd.Flags().GetFloat64("latitude")
	longitude, _ := cmd.Flags().GetFloat64("longitude")
	plmnid, _ := cmd.Flags().GetString("plmnid")
	enbidStart, _ := cmd.Flags().GetUint32("enbidstart")
	pitch, _ := cmd.Flags().GetFloat32("pitch")
	maxDistance, _ := cmd.Flags().GetFloat64("max-neighbor-distance")

	fmt.Printf("Creating honeycomb array of towers. Towers %d. Sectors: %d\n", numTowers, sectorsPerTower)

	mapCenter := model.Coordinate{Lat: latitude, Lng: longitude}

	m, err := honeycomb.GenerateHoneycombTopology(mapCenter, numTowers, sectorsPerTower,
		types.PlmnIDFromString(plmnid), enbidStart, pitch, maxDistance)
	if err != nil {
		return err
	}

	m.Plmn = plmnid // we want the MCC-MNC format in our YAML

	d, err := yaml.Marshal(&m)
	if err != nil {
		fmt.Printf("Unable to marshal model data: %v", err)
		return err
	}

	return ioutil.WriteFile(args[0], d, 0644)
}
