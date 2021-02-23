// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	modelapi "github.com/onosproject/ran-simulator/api/model"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func getNodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Get all E2 nodes",
		RunE:  runGetNodesCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("watch", "w", false, "watch node changes")
	return cmd
}

func createNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <enbid> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Create an E2 node",
		RunE:  runCreateNodeCommand,
	}
	cmd.Flags().UintSlice("cells", []uint{}, "cell ECGIs")
	cmd.Flags().StringSlice("service-models", []string{}, "supported service models")
	cmd.Flags().StringSlice("controllers", []string{}, "E2T controller")
	return cmd
}

func getNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Get an E2 node",
		RunE:  runGetNodeCommand,
	}
	return cmd
}

func deleteNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete an E2 node",
		RunE:  runDeleteNodeCommand,
	}
	return cmd
}

func startNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Start E2 node agent",
		RunE:  runStartNodeCommand,
	}
	return cmd
}

func stopNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Stop E2 node agent",
		RunE:  runStopNodeCommand,
	}
	return cmd
}

func getNodeClient(cmd *cobra.Command) (modelapi.NodeModelClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return modelapi.NewNodeModelClient(conn), conn, nil
}

func runGetNodesCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		Output("%-16s %-8s %-32s %-16s %-20s\n", "EnbID", "Status", "Cell ECGIs", "Service Models", "E2T Controllers")
	}

	if watch, _ := cmd.Flags().GetBool("watch"); watch {
		stream, err := client.WatchNodes(context.Background(), &modelapi.WatchNodesRequest{NoReplay: false})
		if err != nil {
			return err
		}
		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			node := r.Node
			Output("%-16d %-8s %-32s %-16s %-20s\n", node.EnbID, node.Status, catECGIs(node.CellECGIs), catStrings(node.ServiceModels), catStrings(node.Controllers))
		}

	} else {

		stream, err := client.ListNodes(context.Background(), &modelapi.ListNodesRequest{})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			node := r.Node
			Output("%-16d %-8s %-32s %-16s %-20s\n", node.EnbID, node.Status, catECGIs(node.CellECGIs), catStrings(node.ServiceModels), catStrings(node.Controllers))
		}
	}

	return nil
}

func optionsToNode(cmd *cobra.Command, node *types.Node) (*types.Node, error) {
	cells, _ := cmd.Flags().GetUintSlice("cells")
	models, _ := cmd.Flags().GetStringSlice("service-models")
	controllers, _ := cmd.Flags().GetStringSlice("controllers")

	node.CellECGIs = toECGIs(cells)
	node.ServiceModels = models
	node.Controllers = controllers
	return node, nil
}

func runCreateNodeCommand(cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	node, err := optionsToNode(cmd, &types.Node{EnbID: types.EnbID(enbid)})
	if err != nil {
		return err
	}

	_, err = client.CreateNode(context.Background(), &modelapi.CreateNodeRequest{Node: node})
	if err != nil {
		return err
	}
	Output("Node %d created\n", enbid)
	return nil
}

func outputNode(node *types.Node) {
	Output("EnbID: %-16d\nStatus: %s\nCell EGGIs: %s\nService Models: %s\nControllers: %s\n",
		node.EnbID, node.Status, catECGIs(node.CellECGIs), catStrings(node.ServiceModels), catStrings(node.Controllers))
}

func runGetNodeCommand(cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	res, err := client.GetNode(context.Background(), &modelapi.GetNodeRequest{EnbID: types.EnbID(enbid)})
	if err != nil {
		return err
	}

	outputNode(res.Node)
	return nil
}

func runDeleteNodeCommand(cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.DeleteNode(context.Background(), &modelapi.DeleteNodeRequest{EnbID: types.EnbID(enbid)})
	if err != nil {
		return err
	}

	Output("Node %d deleted\n", enbid)
	return nil
}

func runControlCommand(command string, cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := &modelapi.AgentControlRequest{EnbID: types.EnbID(enbid), Command: command}
	res, err := client.AgentControl(context.Background(), request)
	if err != nil {
		return err
	}
	outputNode(res.Node)
	return nil
}

func runStartNodeCommand(cmd *cobra.Command, args []string) error {
	return runControlCommand("start", cmd, args)
}

func runStopNodeCommand(cmd *cobra.Command, args []string) error {
	return runControlCommand("stop", cmd, args)
}

func toECGIs(ids []uint) []types.ECGI {
	ecgis := make([]types.ECGI, 0, len(ids))
	for _, id := range ids {
		ecgis = append(ecgis, types.ECGI(id))
	}
	return ecgis
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
