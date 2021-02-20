// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/spf13/cobra"
	"strconv"
)

func getCellsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cells",
		Short: "Get all cells",
		RunE:  runGetCellsCommand,
	}
	return cmd
}

func getCellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Get a cell",
		RunE:  runGetCellCommand,
	}
	return cmd
}

func deleteCellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a cell",
		RunE:  runDeleteCellCommand,
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

	Output("%-16s %7s %7s %7s %9s %9s %7s %7s %-8s\n", "ECGI", "#UEs", "Max UEs", "TxDB", "Lat", "Lng", "Azimuth", "Arc", "Color")
	for {
		r, err := stream.Recv()
		if err != nil {
			break
		}
		cell := r.Cell
		Output("%-16d %7d %7d %7.2f %9.3f %9.3f %6d %5d %-8s\n",
			cell.ECGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB,
			cell.Location.Lat, cell.Location.Lng, cell.Sector.Azimuth, cell.Sector.Arc, cell.Color)
	}
	return nil
}

func runGetCellCommand(cmd *cobra.Command, args []string) error {
	ecgi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := modelapi.NewCellModelClient(conn)
	res, err := client.GetCell(context.Background(), &modelapi.GetCellRequest{ECGI: types.ECGI(ecgi)})
	if err != nil {
		return err
	}

	cell := res.Cell
	Output("ECGI: %-16d\nUE Count: %d\nMax UEs: %d\nTxPower dB: %.2f\n",
		cell.ECGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB)
	Output("Latitude: %.3f\nLongitude: %.3f\nAzimuth: %d\nArc: %d\nColor: %s\n",
		cell.Location.Lat, cell.Location.Lng, cell.Sector.Azimuth, cell.Sector.Arc, cell.Color)
	return nil
}

func runDeleteCellCommand(cmd *cobra.Command, args []string) error {
	ecgi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := modelapi.NewCellModelClient(conn)
	_, err = client.DeleteCell(context.Background(), &modelapi.DeleteCellRequest{ECGI: types.ECGI(ecgi)})
	if err != nil {
		return err
	}

	Output("Cell %d deleted\n", ecgi)
	return nil
}
