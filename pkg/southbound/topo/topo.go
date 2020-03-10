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
	"github.com/atomix/go-client/pkg/client/util"
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
	ranSimType       = "Tower"
	ranSimRole       = "Simulator"
	ranSimTimeoutSec = 5
)

var log = logging.GetLogger("southbound", "topo")

// ConnectToTopo is a go function that listens for the connection of onos-topo and
// updates the list of Tower instances on it
func ConnectToTopo(ctx context.Context, topoEndpoint string, serverParams utils.ServerParams) topodevice.DeviceServiceClient {
	log.Infof("Connecting to ONOS Topo...%s", topoEndpoint)
	// Attempt to create connection to the Topo
	opts, err := certs.HandleCertArgs(&serverParams.KeyPath, &serverParams.CertPath)
	opts = append(opts, grpc.WithStreamInterceptor(util.RetryingStreamClientInterceptor(100*time.Millisecond)))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := southbound.Connect(ctx, topoEndpoint, "", "", opts...)
	if err != nil {
		log.Fatal("Failed to connect to %s. Retry. %s", topoEndpoint, err)
	}
	return topodevice.NewDeviceServiceClient(conn)
}

// SyncToTopo updates the list of Tower instances on it
func SyncToTopo(ctx context.Context, topoClient *topodevice.DeviceServiceClient, towers map[types.EcID]*types.Tower) {

	for _, t := range towers {
		topoDevice := createTowerForTopo(t)
		resp, err := (*topoClient).Add(ctx, &topodevice.AddRequest{Device: topoDevice})
		if err != nil {
			log.Warnf("Could not add %s to onos-topo", topoTowerID(t.GetEcID()))
			continue
		}
		if resp.GetDevice().ID != topodevice.ID(topoTowerID(t.GetEcID())) {
			log.Errorf("Unexpected response from topo when adding %s. %v",
				topoTowerID(t.GetEcID()), resp)
		}
	}
	log.Infof("%d tower devices created on onos-topo", len(towers))
}

// createTowerForTopo -- prepare the tower to be added to onos-topo
func createTowerForTopo(tower *types.Tower) *topodevice.Device {
	timeOut := time.Second * ranSimTimeoutSec
	serviceEndpoint := fmt.Sprintf("%s:%d", utils.ServiceName, tower.GetPort())

	towerAttributes := make(map[string]string)
	towerAttributes["longitude"] = fmt.Sprintf("%f", tower.GetLocation().GetLng())
	towerAttributes["latitude"] = fmt.Sprintf("%f", tower.GetLocation().GetLat())
	towerAttributes["createdby"] = utils.ServiceName

	return &topodevice.Device{
		ID:          topodevice.ID(topoTowerID(tower.GetEcID())),
		Address:     serviceEndpoint,
		Version:     ranSimVersion,
		Timeout:     &timeOut,
		Credentials: topodevice.Credentials{},
		TLS:         topodevice.TlsConfig{},
		Type:        ranSimType,
		Role:        ranSimRole,
		Attributes:  towerAttributes,
	}
}

func topoTowerID(towerID types.EcID) string {
	return fmt.Sprintf("%s-%s", utils.TestPlmnID, towerID)
}
