// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package cli

import (
	"context"
	"time"

	clilib "github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/spf13/cobra"
)

func getResetMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resetmetrics",
		Short: "Reset the metrics counters",
		Args:  cobra.MaximumNArgs(0),
		RunE:  runResetMetricsCommand,
	}
	return cmd
}

func runResetMetricsCommand(cmd *cobra.Command, args []string) error {
	conn, err := clilib.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := trafficsim.NewTrafficClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err = client.ResetMetrics(ctx, &trafficsim.ResetMetricsMsg{})
	if err != nil {
		return err
	}
	return nil
}
