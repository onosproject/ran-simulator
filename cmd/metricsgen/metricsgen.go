// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package main

import (
	"fmt"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/rc/pciload"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// A simple tool to generate metrics for cells from the specified model.
func main() {
	rootCmd := getRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metricsgen",
		Short: "PCI cell assignment based on specified RAN model",
	}
	cmd.AddCommand(getPCIMetricsCommand())
	return cmd
}

func getPCIMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pci outfile",
		Short:         "ran-simulator config generation tool for onos-topo",
		SilenceUsage:  false,
		SilenceErrors: false,
		Args:          cobra.ExactArgs(1),
		RunE:          runPCIMetricsCommand,
	}
	cmd.Flags().String("model", "model.yaml", "path of the model.yaml file")
	_ = cmd.MarkFlagRequired("model")
	cmd.Flags().Uint("min-pci", 0, "minimum PCI value")
	cmd.Flags().Uint("max-pci", 503, "maximum PCI value")
	cmd.Flags().Uint("max-collisions", 8, "maximum number of collisions")
	cmd.Flags().Uint32("earfcn-start", 42, "start point for EARFCN generation")
	cmd.Flags().StringSlice("cell-types", []string{"ENTERPRISE", "FEMTO", "OUTDOOR_SMALL"}, "List of cell size types")
	return cmd
}

type auxTop struct {
	Rc *auxRc `yaml:"rc"`
}

type auxRc struct {
	Pci *pciload.PciMetrics `yaml:"pci"`
}

func runPCIMetricsCommand(cmd *cobra.Command, args []string) error {
	modelPath, _ := cmd.Flags().GetString("model")
	minPCI, _ := cmd.Flags().GetUint("min-pci")
	maxPCI, _ := cmd.Flags().GetUint("max-pci")
	maxCollisions, _ := cmd.Flags().GetUint("max-collisions")
	earfcnStart, _ := cmd.Flags().GetUint32("earfcn-start")
	cellSizeTypes, _ := cmd.Flags().GetStringSlice("cell-types")

	m := &model.Model{}
	err := model.LoadConfig(m, modelPath)
	if err != nil {
		fmt.Printf("Unable to read model data: %v", err)
		return err
	}

	pciMetrics := pciload.GeneratePCIMetrics(m, minPCI, maxPCI, maxCollisions, earfcnStart, cellSizeTypes)

	d, err := yaml.Marshal(&auxTop{Rc: &auxRc{Pci: pciMetrics}})
	if err != nil {
		fmt.Printf("Unable to marshal PCI metrics: %v", err)
		return err
	}

	return ioutil.WriteFile(args[0], d, 0644)
}
