// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"fmt"
	clilib "github.com/onosproject/onos-lib-go/pkg/cli"

	// Needed to keep ran-sim happy for the mo
	_ "github.com/onosproject/onos-lib-go/pkg/cli"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

const (
	configName     = "ransim"
	defaultAddress = "ran-simulator:5150"
)

// init initializes the command line
func init() {
	clilib.InitConfig(configName)
}

// Execute runs the root command and any sub-commands.
func Execute() {
	rootCmd := GetRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GenerateCliDocs generate markdown files for onos-cli commands
func GenerateCliDocs() {
	cmd := GetRootCommand()
	err := doc.GenMarkdownTree(cmd, "docs/cli")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// GetRootCommand returns the root onos command
func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "simcli",
		Short:                  "RAN simulator command line client",
		BashCompletionFunction: getBashCompletions(),
		SilenceUsage:           true,
		SilenceErrors:          true,
	}
	cmd.AddCommand(getCreateCommand())
	cmd.AddCommand(getDeleteCommand())
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(getSetCommand())
	//cmd.AddCommand(getUpdateCommand())

	cmd.AddCommand(startNodeCommand())
	cmd.AddCommand(stopNodeCommand())

	cmd.AddCommand(getCompletionCommand())

	return cmd
}

func getCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {node,cell} [args]",
		Short: "Commands for creating simulated entities",
	}

	cmd.AddCommand(createNodeCommand())
	cmd.AddCommand(createCellCommand())
	return cmd
}

func getDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete {node,cell} [args]",
		Short: "Commands for deleting simulated entities",
	}
	cmd.AddCommand(deleteNodeCommand())
	cmd.AddCommand(deleteCellCommand())
	cmd.AddCommand(deleteMetricCommand())
	cmd.AddCommand(deleteMetricsCommand())
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {layout,cells,ues,ueCount} [args]",
		Short: "Commands for retrieving RAN simulator model and other information",
	}

	clilib.AddConfigFlags(cmd, defaultAddress)

	cmd.AddCommand(clilib.GetConfigCommand())
	cmd.AddCommand(getLayoutCommand())

	cmd.AddCommand(getNodesCommand())
	cmd.AddCommand(getNodeCommand())

	cmd.AddCommand(getCellsCommand())
	cmd.AddCommand(getCellCommand())

	cmd.AddCommand(getUEsCommand())
	//cmd.AddCommand(getUECommand())

	cmd.AddCommand(getUECountCommand())

	cmd.AddCommand(getMetricCommand())
	cmd.AddCommand(getMetricsCommand())
	return cmd
}

func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {metric} [args]",
		Short: "Commands for setting RAN simulator model metrics and other information",
	}

	cmd.AddCommand(updateNodeCommand())
	cmd.AddCommand(updateCellCommand())
	cmd.AddCommand(setMetricCommand())
	return cmd
}
