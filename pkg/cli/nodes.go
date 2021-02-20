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

func getNodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Get E2 Nodes",
		RunE:  runGetNodesCommand,
	}
	return cmd
}

func runGetNodesCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := modelapi.NewNodeModelClient(conn)
	stream, err := client.WatchNodes(context.Background(), &modelapi.WatchNodesRequest{NoReplay: false, NoSubscribe: true})
	if err != nil {
		return err
	}

	//Output("%16s %7s %7s %5s %7s %7s %8s\n", "ECGI", "#UEs", "Max UEs", "TxDB", "Lat", "Lng", "Color")
	for {
		r, err := stream.Recv()
		if err != nil {
			break
		}
		node := r.Node
		Output("%16d %v, %v, %v\n", node.EnbID, node.CellECGIs, node.ServiceModels, node.Controllers)
	}
	return nil
}
