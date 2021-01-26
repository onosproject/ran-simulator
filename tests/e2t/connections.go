// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2t

import (
	"context"
	e2tadmin "github.com/onosproject/onos-api/go/onos/e2t/admin"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onosproject/ran-simulator/tests/utils"
)

const (
	e2tAddress = "onos-e2t:5150"
)

// TestConnections test connectivity between e2 agents and e2t
func (s *TestSuite) TestConnections(t *testing.T) {
	// Creates an instance of the simulator
	simulator := utils.CreateRanSimulatorWithName(t, "ran-simulator")
	err := simulator.Install(true)
	assert.NoError(t, err, "could not install device simulator %v", err)

	client, err := newAdminClient()
	assert.NoError(t, err, "unable to connect to E2T admin service %v", err)
	if client == nil {
		return
	}

	connections, err := client.ListE2NodeConnections(context.Background(), &e2tadmin.ListE2NodeConnectionsRequest{})
	assert.NoError(t, err, "unable to fetch connections from E2T admin service %v", err)
	if connections == nil {
		return
	}

	count := 0
	for {
		_, err := connections.Recv()
		if err == io.EOF {
			break
		} else if err == nil {
			count = count + 1
		}
	}
	assert.Equal(t, 2, count, "incorrect connection count")
}

func newAdminClient() (e2tadmin.E2TAdminServiceClient, error) {
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
