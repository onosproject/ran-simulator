// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	simapi "github.com/onosproject/ran-simulator/api/trafficsim"

	"github.com/spf13/cobra"
)

func getUECountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ueCount",
		Short: "Get UE count",
		RunE:  runGetUECountCommand,
	}
	return cmd
}

func runGetUECountCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := simapi.NewTrafficClient(conn)
	stream, err := client.ListUes(context.Background(), &simapi.ListUesRequest{})
	if err != nil {
		return err
	}

	Output("%d\n", countUEs(stream))
	return nil
}

func countUEs(stream simapi.Traffic_ListUesClient) int {
	count := 0
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		count = count + 1
	}
	return count
}

func getUEsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues",
		Short: "Get UEs",
		RunE:  runGetUEsCommand,
	}
	return cmd
}

func runGetUEsCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := simapi.NewTrafficClient(conn)
	stream, err := client.ListUes(context.Background(), &simapi.ListUesRequest{})
	if err != nil {
		return err
	}

	Output("%-16s %-16s %-5s\n", "IMSI", "Serving Cell", "Admitted")
	for {
		r, err := stream.Recv()
		if err != nil {
			break
		}
		ue := r.Ue
		Output("%-16d %-16d %-5t\n", ue.IMSI, ue.ServingTower, ue.Admitted)
	}
	return nil
}
