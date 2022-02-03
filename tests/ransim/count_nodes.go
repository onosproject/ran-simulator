// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ransim

import (
	"context"
	"io"
	"testing"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"

	"github.com/stretchr/testify/assert"

	"github.com/onosproject/ran-simulator/tests/utils"
)

// TestCountNodes tests if there is the correct number of nodes
func (s *TestSuite) TestCountNodes(t *testing.T) {
	nodes, err := getNodes()
	assert.NoError(t, err, "unable to connect to Ransim node service %v", err)
	assert.Equal(t, 2, len(nodes))
}

func getNodes() ([]*modelapi.ListNodesResponse, error) {
	client, err := utils.NewRansimNodeClient()
	if err != nil {
		return nil, err
	}
	stream, err := client.ListNodes(context.Background(), &modelapi.ListNodesRequest{})
	if err != nil {
		return nil, err
	}

	connections := make([]*modelapi.ListNodesResponse, 0)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err == nil {
			connections = append(connections, resp)
		}
	}
	return connections, err

}
