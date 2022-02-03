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

// TestCountCells tests if there is the correct number of cells
func (s *TestSuite) TestCountCells(t *testing.T) {
	cells, err := getCells()
	assert.NoError(t, err, "unable to connect to Ransim cell service %v", err)
	assert.Equal(t, 6, len(cells))
}

func getCells() ([]*modelapi.ListCellsResponse, error) {
	client, err := utils.NewRansimCellClient()
	if err != nil {
		return nil, err
	}
	stream, err := client.ListCells(context.Background(), &modelapi.ListCellsRequest{})
	if err != nil {
		return nil, err
	}

	connections := make([]*modelapi.ListCellsResponse, 0)
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
