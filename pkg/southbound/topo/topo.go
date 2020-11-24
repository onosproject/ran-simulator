// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package topo

import (
	"context"
	"io"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"
	"github.com/onosproject/onos-topo/api/topo"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("southbound", "topo")

// CellCreationHandler a call back function to avoid import cycle
type CellCreationHandler func(device *topo.Object) error

// CellDeletionHandler a call back function to avoid import cycle
type CellDeletionHandler func(device *topo.Object) error

// ConnectToTopo is a go function that listens for the connection of onos-topo and
// listens out for Cell instances on it
func ConnectToTopo(ctx context.Context, topoEndpoint string,
	serverParams utils.ServerParams, createHandler CellCreationHandler,
	deleteHandler CellDeletionHandler) (topo.TopoClient, error) {

	log.Infof("Connecting to ONOS Topo...%s", topoEndpoint)
	// Attempt to create connection to the Topo
	opts, err := certs.HandleCertPaths(serverParams.CaPath, serverParams.KeyPath, serverParams.CertPath, true)
	if err != nil {
		log.Fatal(err)
	}
	opts = append(opts, grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(time.Second)))
	conn, err := southbound.Connect(ctx, topoEndpoint, "", "", opts...)
	if err != nil {
		log.Fatal("Failed to connect to %s. Retry. %s", topoEndpoint, err)
	}

	topoClient := topo.NewTopoClient(conn)
	stream, err := topoClient.Watch(context.Background(), &topo.WatchRequest{})
	if err != nil {
		return nil, err
	}
	for {
		in, err := stream.Recv() // Block here and wait for events from topo
		if err == io.EOF {
			// read done.
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		if in.Event.Object.Type != topo.Object_ENTITY {
			continue
		}
		/* TODO
		if in.GetDevice().GetVersion() != types.E2NodeVersion100 {
			log.Warnf("Only version %s of %s is supported", types.E2NodeVersion100, types.E2NodeType)
			continue
		}
		*/
		switch in.Event.Type {
		case topo.EventType_NONE:
			err := createHandler(&in.Event.Object)
			if err != nil {
				log.Warnf("Unable to create cell from %s. %s", in.Event.Object.ID, err.Error())
				continue
			}
		case topo.EventType_ADDED:
			err := createHandler(&in.Event.Object)
			if err != nil {
				log.Warnf("Unable to create cell from %s. %s", in.Event.Object.ID, err.Error())
				continue
			}
		case topo.EventType_UPDATED:
			// TODO
		case topo.EventType_REMOVED:
			err := deleteHandler(&in.Event.Object)
			if err != nil {
				log.Warnf("Unable to delete cell from %s. %s", in.Event.Object.ID, err.Error())
				continue
			}
		default:
			log.Warnf("topo event type %s not yet handled for %s", in.Event.Type, in.Event.Object.ID)
		}
	}
}
