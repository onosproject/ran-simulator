// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"context"
	"strconv"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func getCellsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cells",
		Short: "Get all cells",
		RunE:  runGetCellsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("watch", "w", false, "watch cell changes")

	return cmd
}

func createCellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <enbid> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a cell",
		RunE:  runCreateCellCommand,
	}
	cmd.Flags().Uint32("max-ues", 10000, "maximum number of UEs connected")
	cmd.Flags().Float64("tx-power", 11.0, "transmit power (dB)")
	cmd.Flags().Float64("lat", 11.0, "geo location latitude")
	cmd.Flags().Float64("lng", 11.0, "geo location longitude")
	cmd.Flags().Int32("azimuth", 0, "azimuth of the coverage arc")
	cmd.Flags().Int32("arc", 120, "angle width of the coverage arc")
	cmd.Flags().UintSlice("neighbors", []uint{}, "neighbor cell ECGIs")
	cmd.Flags().String("color", "blue", "color label")
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

func updateCellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <enbid> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Update a cell",
		RunE:  runUpdateCellCommand,
	}
	cmd.Flags().Uint32("max-ues", 10000, "maximum number of UEs connected")
	cmd.Flags().Float64("tx-power", 11.0, "transmit power (dB)")
	cmd.Flags().Float64("lat", 11.0, "geo location latitude")
	cmd.Flags().Float64("lng", 11.0, "geo location longitude")
	cmd.Flags().Int32("azimuth", 0, "azimuth of the coverage arc")
	cmd.Flags().Int32("arc", 120, "angle width of the coverage arc")
	cmd.Flags().UintSlice("neighbors", []uint{}, "neighbor cell ECGIs")
	cmd.Flags().String("color", "blue", "color label")
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

func getCellClient(cmd *cobra.Command) (modelapi.CellModelClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}

	return modelapi.NewCellModelClient(conn), conn, nil
}

func runGetCellsCommand(cmd *cobra.Command, args []string) error {
	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		Output("%-16s %7s %7s %7s %9s %9s %7s %7s %-8s %s\n",
			"ECGI", "#UEs", "Max UEs", "TxDB", "Lat", "Lng", "Azimuth", "Arc", "Color", "Neighbors")
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	if watch, _ := cmd.Flags().GetBool("watch"); watch {
		stream, err := client.WatchCells(context.Background(), &modelapi.WatchCellsRequest{NoReplay: false})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			cell := r.Cell
			Output("%-16d %7d %7d %7.2f %9.3f %9.3f %7d %7d %-8s %s\n",
				cell.ECGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB,
				cell.Location.Lat, cell.Location.Lng, cell.Sector.Azimuth, cell.Sector.Arc, cell.Color,
				catECGIs(cell.Neighbors))
		}

	} else {

		stream, err := client.ListCells(context.Background(), &modelapi.ListCellsRequest{})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			cell := r.Cell
			Output("%-16d %7d %7d %7.2f %9.3f %9.3f %7d %7d %-8s %s\n",
				cell.ECGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB,
				cell.Location.Lat, cell.Location.Lng, cell.Sector.Azimuth, cell.Sector.Arc, cell.Color,
				catECGIs(cell.Neighbors))
		}
	}
	return nil
}

func optionsToCell(cmd *cobra.Command, cell *types.Cell, update bool) (*types.Cell, error) {
	arc, _ := cmd.Flags().GetInt32("arc")
	azimuth, _ := cmd.Flags().GetInt32("azimuth")
	lat, _ := cmd.Flags().GetFloat64("lat")
	lng, _ := cmd.Flags().GetFloat64("lng")

	if cell.Location == nil {
		cell.Location = &types.Point{Lat: lat, Lng: lng}
	} else {
		if !update || cmd.Flags().Changed("lat") {
			cell.Location.Lng = lng
		}
		if !update || cmd.Flags().Changed("lng") {
			cell.Location.Lng = lng
		}
	}

	if cell.Sector == nil {
		cell.Sector = &types.Sector{Centroid: cell.Location, Azimuth: azimuth, Arc: arc}
	} else {
		cell.Sector.Centroid = cell.Location
		if !update || cmd.Flags().Changed("arc") {
			cell.Sector.Arc = arc
		}
		if !update || cmd.Flags().Changed("azimuth") {
			cell.Sector.Azimuth = azimuth
		}
	}

	color, _ := cmd.Flags().GetString("color")
	if !update || cmd.Flags().Changed("color") {
		cell.Color = color
	}

	maxUEs, _ := cmd.Flags().GetUint32("max-ues")
	if !update || cmd.Flags().Changed("max-ues") {
		cell.MaxUEs = maxUEs
	}

	txDb, _ := cmd.Flags().GetFloat64("tx-power")
	if !update || cmd.Flags().Changed("tx-power") {
		cell.TxPowerdB = txDb
	}

	neighbors, _ := cmd.Flags().GetUintSlice("neighbors")
	if !update || cmd.Flags().Changed("neighbors") {
		cell.Neighbors = toECGIs(neighbors)
	}
	return cell, nil
}

func runCreateCellCommand(cmd *cobra.Command, args []string) error {
	ecgi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	cell, err := optionsToCell(cmd, &types.Cell{ECGI: types.ECGI(ecgi)}, false)
	if err != nil {
		return err
	}

	_, err = client.CreateCell(context.Background(), &modelapi.CreateCellRequest{Cell: cell})
	if err != nil {
		return err
	}
	Output("Cell %d created\n", ecgi)
	return nil
}

func runUpdateCellCommand(cmd *cobra.Command, args []string) error {
	ecgi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get the cell first to prime the update cell with existing values and allow sparse update
	gres, err := client.GetCell(context.Background(), &modelapi.GetCellRequest{ECGI: types.ECGI(ecgi)})
	if err != nil {
		return err
	}

	cell, err := optionsToCell(cmd, gres.Cell, true)
	if err != nil {
		return err
	}

	_, err = client.UpdateCell(context.Background(), &modelapi.UpdateCellRequest{Cell: cell})
	if err != nil {
		return err
	}
	Output("Cell %d updated\n", ecgi)
	return nil
}

func runGetCellCommand(cmd *cobra.Command, args []string) error {
	ecgi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	res, err := client.GetCell(context.Background(), &modelapi.GetCellRequest{ECGI: types.ECGI(ecgi)})
	if err != nil {
		return err
	}

	cell := res.Cell
	Output("ECGI: %-16d\nUE Count: %d\nMax UEs: %d\nTxPower dB: %.2f\n",
		cell.ECGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB)
	Output("Latitude: %.3f\nLongitude: %.3f\nAzimuth: %d\nArc: %d\nColor: %s\nNeighbors: %s\n",
		cell.Location.Lat, cell.Location.Lng, cell.Sector.Azimuth, cell.Sector.Arc, cell.Color,
		catECGIs(cell.Neighbors))
	return nil
}

func runDeleteCellCommand(cmd *cobra.Command, args []string) error {
	ecgi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.DeleteCell(context.Background(), &modelapi.DeleteCellRequest{ECGI: types.ECGI(ecgi)})
	if err != nil {
		return err
	}

	Output("Cell %d deleted\n", ecgi)
	return nil
}
