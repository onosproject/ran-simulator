// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/spf13/cobra"
)

func getCellsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cells",
		Short: "Get Cells",
		RunE:  runGetCellsCommand,
	}
	return cmd
}

func runGetCellsCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := modelapi.NewCellModelClient(conn)
	stream, err := client.WatchCells(context.Background(), &modelapi.WatchCellsRequest{NoReplay: false, NoSubscribe: true})
	if err != nil {
		return err
	}

	Output("%-16s %7s %7s %7s %9s %9s %-8s\n", "ECGI", "#UEs", "Max UEs", "TxDB", "Lat", "Lng", "Color")
	for {
		r, err := stream.Recv()
		if err != nil {
			break
		}
		cell := r.Cell
		Output("%-16d %7d %7d %7.2f %9.3f %9.3f %-8s\n", cell.ECGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB,
			cell.Location.Lat, cell.Location.Lng, cell.Color)
	}
	return nil
}
