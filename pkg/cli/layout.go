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

func getLayoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "layout",
		Short: "Get Layout",
		RunE:  runGetLayoutCommand,
	}
	return cmd
}

func runGetLayoutCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := simapi.NewTrafficClient(conn)
	ml, err := client.GetMapLayout(context.Background(), &simapi.MapLayoutRequest{})
	if err != nil {
		return err
	}

	Output("Center: %7.3f,%7.3f\nZoom: %5.2f\nFade: %v\nShowRoutes: %v\nShowPower: %v\nLocationsScale: %5.2f\n",
		ml.Center.Lat, ml.Center.Lng, ml.Zoom, ml.Fade, ml.ShowRoutes, ml.ShowPower, ml.LocationsScale)
	return nil
}
