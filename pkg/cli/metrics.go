// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	metricsapi "github.com/onosproject/ran-simulator/api/metrics"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"strconv"
)

func getMetricCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metric <entity-id> <metric-name>",
		Short: "Get metric value",
		Args:  cobra.ExactArgs(2),
		RunE:  runGetMetricCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	return cmd
}

func setMetricCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metric <entity-id> <metric-name> <value> --type {string|intX|uintX|floatX|bool; where X={8|16|32|64}}",
		Short: "Set metric value",
		Args:  cobra.ExactArgs(3),
		RunE:  runSetMetricCommand,
	}
	cmd.Flags().String("type", "string", "value type: string|intX|uintX|floatX|bool; where X={8|16|32|64}")
	return cmd
}

func deleteMetricCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metric <entity-id> <metric-name>",
		Short: "Delete a metric",
		Args:  cobra.ExactArgs(2),
		RunE:  runDeleteMetricCommand,
	}
	return cmd
}

func getMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics <entity-id>",
		Short: "Get all metrics of an entity",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetMetricsCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	return cmd
}

func deleteMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrica <entity-id>",
		Short: "Delete all metrics of an entity",
		Args:  cobra.ExactArgs(1),
		RunE:  runDeleteMetricsCommand,
	}
	return cmd
}

func getMetricsClient(cmd *cobra.Command) (metricsapi.MetricsServiceClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return metricsapi.NewMetricsServiceClient(conn), conn, nil
}

func runGetMetricCommand(cmd *cobra.Command, args []string) error {
	entityID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}
	name := args[1]

	client, conn, err := getMetricsClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := client.Get(context.Background(), &metricsapi.GetRequest{EntityID: entityID, Name: name})
	if err != nil {
		return err
	}

	m := resp.Metric

	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		Output("%s=%s (%s)\n", m.Key, m.Value, m.Type)
	} else {
		Output("%s\n", m.Value)
	}
	return nil
}

func runGetMetricsCommand(cmd *cobra.Command, args []string) error {
	entityID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getMetricsClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := client.List(context.Background(), &metricsapi.ListRequest{EntityID: entityID})
	if err != nil {
		return err
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	for _, m := range resp.Metrics {
		if verbose {
			Output("%s=%s (%s)\n", m.Key, m.Value, m.Type)
		} else {
			Output("%s=%s\n", m.Key, m.Value)
		}
	}
	return nil
}

func runSetMetricCommand(cmd *cobra.Command, args []string) error {
	entityID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}
	name := args[1]
	value := args[2]
	valueType, _ := cmd.Flags().GetString("type")

	client, conn, err := getMetricsClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	metric := &metricsapi.Metric{
		EntityID: entityID,
		Key:      name,
		Value:    value,
		Type:     valueType,
	}
	_, err = client.Set(context.Background(), &metricsapi.SetRequest{Metric: metric})
	if err != nil {
		return err
	}
	return nil
}

func runDeleteMetricCommand(cmd *cobra.Command, args []string) error {
	entityID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}
	name := args[1]

	client, conn, err := getMetricsClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.Delete(context.Background(), &metricsapi.DeleteRequest{EntityID: entityID, Name: name})
	if err != nil {
		return err
	}
	return nil
}

func runDeleteMetricsCommand(cmd *cobra.Command, args []string) error {
	entityID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getMetricsClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.DeleteAll(context.Background(), &metricsapi.DeleteAllRequest{EntityID: entityID})
	if err != nil {
		return err
	}
	return nil
}
