// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package cli

import (
	"context"
	"fmt"
	"strconv"
	"time"

	clilib "github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/spf13/cobra"
)

func getSetNumUEsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setnumues <numues>",
		Short: "Change the number of UEs in the RAN simulation",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetNumUEsCommand,
	}
	return cmd
}

func runSetNumUEsCommand(cmd *cobra.Command, args []string) error {
	conn, err := clilib.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	numUEs, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	request := trafficsim.SetNumberUEsRequest{
		Number: uint32(numUEs),
	}
	client := trafficsim.NewTrafficClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	response, err := client.SetNumberUEs(ctx, &request)
	if err != nil {
		return err
	}
	if response.Number != uint32(numUEs) {
		return fmt.Errorf("Unexpected # UEs. Expected %d Got %d", numUEs, response.Number)
	}

	return nil
}
