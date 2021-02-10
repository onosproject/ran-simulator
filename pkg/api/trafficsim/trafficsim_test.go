// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package trafficsim

import (
	"context"
	simapi "github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
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
	m := &model.Model{}
	err := model.LoadConfig(m, "../../model/test")
	if err != nil {
		return &Service{}, err
	}
	nodeStore := nodes.NewNodeRegistry(m.Nodes)
	cellStore := cells.NewCellRegistry(m.Cells, nodeStore)
	ueStore := ues.NewUERegistry(m.UECount, cellStore)
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
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
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

	stream, err := client.ListUes(context.TODO(), &simapi.ListUesRequest{WithoutReplay: false})
	assert.NoError(t, err, "unable to list UEs")

	assert.Equal(t, 12, countItems(t, stream, &simapi.ListUesResponse{}), "incorrect UE count")
	_, err = client.SetNumberUEs(context.TODO(), &simapi.SetNumberUEsRequest{
		Number: 16,
	})
	assert.NoError(t, err, "unable to set UE count")

	stream, err = client.ListUes(context.TODO(), &simapi.ListUesRequest{WithoutReplay: false})
	assert.NoError(t, err, "unable to list UEs")
	assert.Equal(t, 16, countItems(t, stream, &simapi.ListUesResponse{}), "incorrect revised UE count")
}

func countItems(t *testing.T, stream grpc.ClientStream, msg interface{}) int {
	count := 0
	for {
		err := stream.RecvMsg(msg)
		if err == io.EOF {
			break
		}
		assert.NoError(t, err, "unable to read stream")
		count = count + 1
	}
	return count
}
