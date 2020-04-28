// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
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
		Use:           "honeycomb outfile",
		Short:         "ran-simulator config generation tool",
		SilenceUsage:  false,
		SilenceErrors: false,
		Args:          cobra.ExactArgs(1),
		RunE:          runHoneycombCommand,
	}
	cmd.Flags().UintP("towers", "t", 0, "number of towers")
	_ = cmd.MarkFlagRequired("towers")
	cmd.Flags().UintP("sectors-per-tower", "s", 3, "sectors per tower (3 or 6)")
	cmd.Flags().Float32P("latitude", "a", 52.5200, "Map centre latitude in degrees")
	cmd.Flags().Float32P("longitude", "g", 13.4050, "Map centre longitude in degrees")
	cmd.Flags().String("plmnid", "315010", "PlmnID")
	cmd.Flags().Uint16P("ecidstart", "e", 5152, "Ecid start")
	cmd.Flags().Uint16P("portstart", "p", 5152, "Port start")
	cmd.Flags().Float32P("pitch", "i", 0.02, "pitch between cells in degrees")

	return cmd
}

func runHoneycombCommand(cmd *cobra.Command, args []string) error {
	numTowers, _ := cmd.Flags().GetUint("towers")
	sectorsPerTower, _ := cmd.Flags().GetUint("sectors-per-tower")
	if sectorsPerTower != 3 && sectorsPerTower != 6 {
		return fmt.Errorf("only 3 or 6 are allowed for 'sectors-per-tower'")
	}
	latitude, _ := cmd.Flags().GetFloat32("latitude")
	longitude, _ := cmd.Flags().GetFloat32("longitude")
	plmnid, _ := cmd.Flags().GetString("plmnid")
	ecidStart, _ := cmd.Flags().GetUint16("ecidstart")
	portStart, _ := cmd.Flags().GetUint16("portstart")
	pitch, _ := cmd.Flags().GetFloat32("pitch")

	fmt.Printf("Creating honeycomb array of towers. Towers %d. Sectors: %d\n", numTowers, sectorsPerTower)

	newConfig, err := config.HoneycombGenerator(numTowers, sectorsPerTower, latitude, longitude,
		types.PlmnID(plmnid), ecidStart, portStart, pitch)
	if err != nil {
		return err
	}
	err = config.Checker(newConfig)
	if err != nil {
		return err
	}

	viper.Set("mapcentre", newConfig.MapCentre)
	viper.Set("towerslayout", newConfig.TowersLayout)
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
