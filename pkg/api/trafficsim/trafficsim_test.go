// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package trafficsim

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"

	simapi "github.com/onosproject/onos-api/go/onos/ransim/trafficsim"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"

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
	m := &model.Model{}
	err := model.LoadConfig(m, "../../model/test")
	if err != nil {
		return &Service{}, err
	}
	nodeStore := nodes.NewNodeRegistry(m.Nodes)
	cellStore := cells.NewCellRegistry(m.Cells, nodeStore)
	ueStore := ues.NewUERegistry(m.UECount, cellStore, "random")
	return &Service{model: m, cellStore: cellStore, ueStore: ueStore}, nil
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
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}

func TestMapLayout(t *testing.T) {
	client := simapi.NewTrafficClient(createServerConnection(t))
	assert.NotNil(t, client, "unable to create gRPC client")

	r, err := client.GetMapLayout(context.TODO(), &simapi.MapLayoutRequest{})
	assert.NoError(t, err, "unable to get map layout")
	assert.Equal(t, 45.0, r.Center.Lat, "incorrect latitude")
	assert.Equal(t, -30.0, r.Center.Lng, "incorrect longitude")
	assert.Equal(t, float32(0.8), r.Zoom, "incorrect zoom")
	assert.Equal(t, float32(1.0), r.LocationsScale, "incorrect scale")
	assert.Equal(t, true, r.Fade)
	assert.Equal(t, true, r.ShowRoutes)
	assert.Equal(t, true, r.ShowPower)
}

func TestServiceBasics(t *testing.T) {
	client := simapi.NewTrafficClient(createServerConnection(t))
	assert.NotNil(t, client, "unable to create gRPC client")

	stream, err := client.ListUes(context.Background(), &simapi.ListUesRequest{})
	assert.NoError(t, err, "unable to list UEs")
	numUes := countUEs(t, stream)
	assert.Equal(t, 12, numUes)
	_, err = client.SetNumberUEs(context.TODO(), &simapi.SetNumberUEsRequest{
		Number: 16,
	})
	assert.NoError(t, err, "unable to set UE count")

	stream, err = client.ListUes(context.TODO(), &simapi.ListUesRequest{})
	assert.NoError(t, err, "unable to list UEs")
	numUes = countUEs(t, stream)
	assert.Equal(t, 16, numUes)

}

func countUEs(t *testing.T, stream simapi.Traffic_ListUesClient) int {
	count := 0
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		count = count + 1
		t.Log(count)
	}
	return count
}
