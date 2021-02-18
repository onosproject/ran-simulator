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

	client := simapi.NewTrafficClient(conn)
	stream, err := client.ListCells(context.Background(), &simapi.ListCellsRequest{})
	if err != nil {
		return err
	}

	Output("%16s %7s %7s %5s %7s %7s %8s\n", "ECGI", "#UEs", "Max UEs", "TxDB", "Lat", "Lng", "Color")
	for {
		r, err := stream.Recv()
		if err != nil {
			break
		}
		cell := r.Cell
		Output("%16d %7d %7d %5.2f %7.3f %7.3f %8s\n", cell.ECGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB,
			cell.Location.Lat, cell.Location.Lng, cell.Color)
	}
	return nil
}
