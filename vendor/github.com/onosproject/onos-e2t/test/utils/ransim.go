// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package utils

import (
	"context"
	"io"
	"testing"

	"github.com/onosproject/helmit/pkg/kubernetes"

	"github.com/onosproject/helmit/pkg/helm"
	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/stretchr/testify/assert"

	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ConnectRansimServiceHost connects to ransim service
func ConnectRansimServiceHost(release *helm.HelmRelease) (*grpc.ClientConn, error) {
	client := kubernetes.NewForReleaseOrDie(release)
	services, err := client.CoreV1().Services().List()
	if err != nil {
		return nil, err
	}
	tlsConfig, err := creds.GetClientCredentials()
	if err != nil {
		return nil, err
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	ransimServiceAddress := getRansimServiceAddress(services[0].Name)
	return grpc.DialContext(context.Background(), ransimServiceAddress, opts...)
}

func GetNodes(t *testing.T, nodeClient modelapi.NodeModelClient) []*ransimtypes.Node {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := nodeClient.ListNodes(ctx, &modelapi.ListNodesRequest{})
	assert.NoError(t, err)
	var nodes []*ransimtypes.Node
	for {
		e2node, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			return []*ransimtypes.Node{}
		}
		nodes = append(nodes, e2node.Node)
	}
	return nodes
}

func GetCells(t *testing.T, cellClient modelapi.CellModelClient) []*ransimtypes.Cell {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := cellClient.ListCells(ctx, &modelapi.ListCellsRequest{})
	assert.NoError(t, err)
	var cellsList []*ransimtypes.Cell
	for {
		cell, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return []*ransimtypes.Cell{}
		}

		cellsList = append(cellsList, cell.Cell)
	}
	return cellsList
}

func GetNumCells(t *testing.T, cellClient modelapi.CellModelClient) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := cellClient.ListCells(ctx, &modelapi.ListCellsRequest{})
	assert.NoError(t, err)
	numCells := 0
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0
		}
		numCells++
	}
	return numCells
}

func GetNumNodes(t *testing.T, nodeClient modelapi.NodeModelClient) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := nodeClient.ListNodes(ctx, &modelapi.ListNodesRequest{})
	assert.NoError(t, err)
	numNodes := 0
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0
		}
		numNodes++
	}
	return numNodes
}

func GetRansimCellClient(t *testing.T, release *helm.HelmRelease) modelapi.CellModelClient {
	conn, err := ConnectRansimServiceHost(release)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	return modelapi.NewCellModelClient(conn)
}

func GetRansimNodeClient(t *testing.T, release *helm.HelmRelease) modelapi.NodeModelClient {
	conn, err := ConnectRansimServiceHost(release)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	return modelapi.NewNodeModelClient(conn)
}
