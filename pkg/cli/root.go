// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package cli

import (
	clilib "github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "ransim"
	defaultAddress = "ran-simulator:5150"
)

// init initializes the command line
func init() {
	clilib.InitConfig(configName)
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the RAN service
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ransim {get|set|watch|setnumues|resetmetrics|log} [args]",
		Short: "ONOS RAN Simulator subsystem commands",
	}

	clilib.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(clilib.GetConfigCommand())
	cmd.AddCommand(getSetNumUEsCommand())
	cmd.AddCommand(getResetMetricsCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}
