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
