// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
