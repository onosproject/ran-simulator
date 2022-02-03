// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"

	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	ransimAddress = "ran-simulator:5150"
)

// NewRansimConnection returns a connection for engaging with the ransim API
func NewRansimConnection() (*grpc.ClientConn, error) {
	tlsConfig, err := creds.GetClientCredentials()
	if err != nil {
		return nil, err
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	conn, err := grpc.DialContext(context.Background(), ransimAddress, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// NewRansimNodeClient returns a node model client
func NewRansimNodeClient() (modelapi.NodeModelClient, error) {
	conn, err := NewRansimConnection()
	if err != nil {
		return nil, err
	}
	return modelapi.NewNodeModelClient(conn), nil
}

// NewRansimCellClient returns a cell model client
func NewRansimCellClient() (modelapi.CellModelClient, error) {
	conn, err := NewRansimConnection()
	if err != nil {
		return nil, err
	}
	return modelapi.NewCellModelClient(conn), nil
}
