// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"io"

	e2tadmin "github.com/onosproject/onos-api/go/onos/e2t/admin"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	e2tAddress = "onos-e2t:5150"
)

// GetE2Connections returns the list of current connections obtained via E2T admin API
func GetE2Connections() ([]*e2tadmin.ListE2NodeConnectionsResponse, error) {
	client, err := NewE2TAdminClient()
	if err != nil {
		return nil, err
	}

	stream, err := client.ListE2NodeConnections(context.Background(), &e2tadmin.ListE2NodeConnectionsRequest{})
	if err != nil {
		return nil, err
	}

	connections := make([]*e2tadmin.ListE2NodeConnectionsResponse, 0)
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

// NewE2TAdminClient returns a client for engaging with the E2T admin API
func NewE2TAdminClient() (e2tadmin.E2TAdminServiceClient, error) {
	tlsConfig, err := creds.GetClientCredentials()
	if err != nil {
		return nil, err
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	conn, err := grpc.DialContext(context.Background(), e2tAddress, opts...)
	if err != nil {
		return nil, err
	}
	return e2tadmin.NewE2TAdminServiceClient(conn), nil
}
