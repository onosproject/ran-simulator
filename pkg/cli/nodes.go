// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/onosproject/ran-simulator/api/types"
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

	Output("%-16s %-32s %-16s %-20s\n", "EnbID", "Cell ECGIs", "Service Models", "E2T Controllers")
	for {
		r, err := stream.Recv()
		if err != nil {
			break
		}
		node := r.Node
		Output("%-16d %-32s %-16s %-20s\n", node.EnbID, catECGIs(node.CellECGIs), catStrings(node.ServiceModels), catStrings(node.Controllers))
	}
	return nil
}

func catECGIs(ecgis []types.ECGI) string {
	s := ""
	for _, ecgi := range ecgis {
		s = s + fmt.Sprintf(",%d", ecgi)
	}
	if len(s) > 1 {
		return s[1:]
	}
	return s
}

func catStrings(strings []string) string {
	s := ""
	for _, string := range strings {
		s = s + "," + string
	}
	if len(s) > 1 {
		return s[1:]
	}
	return s
}
