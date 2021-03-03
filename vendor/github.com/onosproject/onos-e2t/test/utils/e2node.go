// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package utils

import (
	"context"
	"io"

	"github.com/onosproject/onos-api/go/onos/e2t/admin"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	OnosE2TAddress = "onos-e2t:5150"
)

// GetNodeIDs get list of E2 node IDs
// TODO this function should be replaced with topology API
func GetNodeIDs() ([]string, error) {
	tlsConfig, err := creds.GetClientCredentials()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var nodeIDs []string
	if err != nil {
		return []string{}, err
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	conn, err := grpc.DialContext(ctx, OnosE2TAddress, opts...)
	if err != nil {
		return []string{}, err
	}
	adminClient := admin.NewE2TAdminServiceClient(conn)
	connections, err := adminClient.ListE2NodeConnections(ctx, &admin.ListE2NodeConnectionsRequest{})

	if err != nil {
		return []string{}, err
	}

	for {
		connection, err := connections.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return []string{}, err
		}
		if connection != nil {
			nodeID := connection.Id
			nodeIDs = append(nodeIDs, nodeID)
		}
	}
	return nodeIDs, nil
}
