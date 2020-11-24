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
	"errors"
	"time"

	api "github.com/onosproject/onos-lib-go/api/logging"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Gets a logger attribute (e.g. level)",
	}
	cmd.AddCommand(getGetLevelCommand())
	return cmd
}

func getGetLevelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "level logger_name",
		Short: "Gets a logger level",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetLevelCommand,
	}

	return cmd
}

func runGetLevelCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
	}()

	name := args[0]
	if name == "" {
		return errors.New("the logger name should be provided")
	}

	client := api.NewLoggerClient(conn)
	req := api.GetLevelRequest{
		LoggerName: name,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	response, err := client.GetLevel(ctx, &req)

	if err != nil {
		return err
	}

	cli.Output("%s logger level is %s\n", name, response.Level.String())

	return err
}
