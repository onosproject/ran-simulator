// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package trafficsim

import (
	"context"
	simapi "github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/pkg/model"
	"io"
	"net"
	"testing"

	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func newTestService() (northbound.Service, error) {
	return &Service{
		model: &model.Model{
			UECount: 10,
			UEs:     model.NewUERegistry(10),
		},
	}, nil
}

func createServerConnection(t *testing.T) *grpc.ClientConn {
	lis = bufconn.Listen(1024 * 1024)
	s, err := newTestService()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	server := grpc.NewServer()
	s.Register(server)

	go func() {
		if err := server.Serve(lis); err != nil {
			assert.NoError(t, err, "Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}

func TestServiceBasics(t *testing.T) {
	conn := createServerConnection(t)
	client := simapi.NewTrafficClient(conn)
	assert.NotNil(t, client, "unable to create gRPC client")

	stream, err := client.ListUes(context.TODO(), &simapi.ListUesRequest{WithoutReplay: false})
	assert.NoError(t, err, "unable to list UEs")

	assert.Equal(t, 10, countUEs(t, stream), "incorrect UE count")
	_, err = client.SetNumberUEs(context.TODO(), &simapi.SetNumberUEsRequest{
		Number: 16,
	})
	assert.NoError(t, err, "unable to set UE count")

	stream, err = client.ListUes(context.TODO(), &simapi.ListUesRequest{WithoutReplay: false})
	assert.NoError(t, err, "unable to list UEs")
	assert.Equal(t, 16, countUEs(t, stream), "incorrect revised UE count")
}

func countUEs(t *testing.T, stream simapi.Traffic_ListUesClient) int {
	count := 0
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err, "unable to read UE stream")
		count = count + 1
	}
	return count
}
