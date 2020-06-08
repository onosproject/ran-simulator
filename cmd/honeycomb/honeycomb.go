// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"fmt"
	"github.com/onosproject/onos-topo/pkg/bulk"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Short: "ran-simulator config generation tool",
	}
	cmd.AddCommand(getHoneycombTopoCommand())
	cmd.AddCommand(getHoneycombConfigCommand())

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
	cmd.Flags().String("plmnid", "315010", "PlmnID")
	cmd.Flags().Uint16P("ecidstart", "e", 5152, "Ecid start")
	cmd.Flags().Uint16P("portstart", "p", 5152, "Port start")
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
	ecidStart, _ := cmd.Flags().GetUint16("ecidstart")
	portStart, _ := cmd.Flags().GetUint16("portstart")
	pitch, _ := cmd.Flags().GetFloat32("pitch")

	fmt.Printf("Creating honeycomb array of towers. Towers %d. Sectors: %d\n", numTowers, sectorsPerTower)

	newTopo, err := config.HoneycombTopoGenerator(numTowers, sectorsPerTower, latitude, longitude,
		types.PlmnID(plmnid), ecidStart, portStart, pitch)
	if err != nil {
		return err
	}
	err = bulk.Checker(newTopo)
	if err != nil {
		return err
	}

	viper.Set("topodevices", newTopo.TopoDevices)
	// Set the file name of the configurations file
	viper.SetConfigName("onos")
	viper.SetConfigType("yaml")
	if err := viper.WriteConfigAs(args[0]); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("Config YAML file written to %s\n", args[0])
	return nil
}

func getHoneycombConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "config outfile",
		Short:         "ran-simulator config generation tool for onos-config",
		SilenceUsage:  false,
		SilenceErrors: false,
		Args:          cobra.ExactArgs(1),
		RunE:          runHoneycombConfigCommand,
	}
	cmd.Flags().UintP("towers", "t", 0, "number of towers")
	_ = cmd.MarkFlagRequired("towers")
	cmd.Flags().UintP("sectors-per-tower", "s", 3, "sectors per tower (3 or 6)")
	cmd.Flags().String("plmnid", "315010", "PlmnID")
	cmd.Flags().Uint16P("ecidstart", "e", 5152, "Ecid start")
	return cmd
}

func runHoneycombConfigCommand(cmd *cobra.Command, args []string) error {
	numTowers, _ := cmd.Flags().GetUint("towers")
	sectorsPerTower, _ := cmd.Flags().GetUint("sectors-per-tower")
	if sectorsPerTower != 3 && sectorsPerTower != 6 {
		return fmt.Errorf("only 3 or 6 are allowed for 'sectors-per-tower'")
	}
	plmnid, _ := cmd.Flags().GetString("plmnid")
	ecidStart, _ := cmd.Flags().GetUint16("ecidstart")

	fmt.Printf("Creating honeycomb array of towers. Towers %d. Sectors: %d\n", numTowers, sectorsPerTower)

	newConfig, err := config.HoneycombConfigGenerator(numTowers, sectorsPerTower,
		types.PlmnID(plmnid), ecidStart)
	if err != nil {
		return err
	}

	viper.Set("setrequest", newConfig.SetRequest)
	// Set the file name of the configurations file
	viper.SetConfigName("onos")
	viper.SetConfigType("yaml")
	if err := viper.WriteConfigAs(args[0]); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("Config YAML file written to %s\n", args[0])
	return nil
}
