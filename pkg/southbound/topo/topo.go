// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topo

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"
	topodevice "github.com/onosproject/onos-topo/api/device"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"google.golang.org/grpc"
	"time"
)

const (
	ranSimVersion    = "1.0.0"
	ranSimType       = "RanSimulator"
	ranSimTimeoutSec = 5
)

var log = logging.GetLogger("southbound", "topo")

// ConnectToTopo is a go function that listens for the connection of onos-topo and
// updates the list of Cell instances on it
func ConnectToTopo(ctx context.Context, topoEndpoint string, serverParams utils.ServerParams) topodevice.DeviceServiceClient {
	log.Infof("Connecting to ONOS Topo...%s", topoEndpoint)
	// Attempt to create connection to the Topo
	opts, err := certs.HandleCertPaths(serverParams.CaPath, serverParams.KeyPath, serverParams.CertPath, true)
	if err != nil {
		log.Fatal(err)
	}
	opts = append(opts, grpc.WithUnaryInterceptor(southbound.RetryingUnaryClientInterceptor()))
	conn, err := southbound.Connect(ctx, topoEndpoint, "", "", opts...)
	if err != nil {
		log.Fatal("Failed to connect to %s. Retry. %s", topoEndpoint, err)
	}
	return topodevice.NewDeviceServiceClient(conn)
}

// SyncToTopo updates the list of Cell instances on it
// Will block until topo comes available
func SyncToTopo(ctx context.Context, topoClient *topodevice.DeviceServiceClient, cells map[types.ECGI]*types.Cell) {

	for _, t := range cells {
		topoDevice := createCellForTopo(t)
		resp, err := (*topoClient).Add(ctx, &topodevice.AddRequest{Device: topoDevice})
		if err != nil {
			log.Warnf("Could not add %s to onos-topo %s", topoCellID(t.GetEcgi()), err.Error())
			continue
		}
		if resp.GetDevice().ID != topodevice.ID(topoCellID(t.GetEcgi())) {
			log.Errorf("Unexpected response from topo when adding %s. %v",
				topoCellID(t.GetEcgi()), resp)
		}
	}
	log.Infof("%d cell devices created on onos-topo", len(cells))
}

// createCellForTopo -- prepare the cell to be added to onos-topo
func createCellForTopo(cell *types.Cell) *topodevice.Device {
	timeOut := time.Second * ranSimTimeoutSec
	serviceEndpoint := fmt.Sprintf("%s:%d", utils.ServiceName, cell.GetPort())

	cellAttributes := make(map[string]string)
	cellAttributes["longitude"] = fmt.Sprintf("%f", cell.GetLocation().GetLng())
	cellAttributes["latitude"] = fmt.Sprintf("%f", cell.GetLocation().GetLat())
	cellAttributes["azimuth"] = fmt.Sprintf("%d", cell.GetSector().GetAzimuth())
	cellAttributes["arc"] = fmt.Sprintf("%d", cell.GetSector().GetArc())
	cellAttributes["createdby"] = utils.ServiceName

	return &topodevice.Device{
		ID:          topodevice.ID(topoCellID(cell.GetEcgi())),
		Address:     serviceEndpoint,
		Version:     ranSimVersion,
		Timeout:     &timeOut,
		Credentials: topodevice.Credentials{},
		TLS:         topodevice.TlsConfig{},
		Type:        ranSimType,
		Attributes:  cellAttributes,
	}
}

func topoCellID(cellID *types.ECGI) string {
	return fmt.Sprintf("%s-%s", cellID.PlmnID, cellID.EcID)
}
