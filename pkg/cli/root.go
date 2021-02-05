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
	cmd.AddCommand(getGetCommand())
	//cmd.AddCommand(getSetCommand())

	cmd.AddCommand(getCompletionCommand())

	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {ueCount,ues} [args]",
		Short: "RAN simulator get info commands",
	}

	clilib.AddConfigFlags(cmd, defaultAddress)

	cmd.AddCommand(clilib.GetConfigCommand())
	cmd.AddCommand(getUECountCommand())
	cmd.AddCommand(getUEsCommand())
	return cmd
}
