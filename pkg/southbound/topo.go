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

package southbound

import (
	"fmt"
	topodevice "github.com/onosproject/onos-topo/api/device"
	"github.com/onosproject/ran-simulator/api/types"
	"time"
)

const (
	ranSimVersion    = "1.0.0"
	ranSimType       = "Tower"
	ranSimRole       = "Simulator"
	ranSimTimeoutSec = 5
)

// CreateTowerOnTopo -- prepare the tower to be added to onos-topo
func CreateTowerOnTopo(towerID types.EcID, serviceEndpoint string, tower *types.Tower) *topodevice.Device {
	timeOut := time.Second * ranSimTimeoutSec

	towerAttributes := make(map[string]string)
	towerAttributes["longitude"] = fmt.Sprintf("%f", tower.GetLocation().GetLng())
	towerAttributes["latitude"] = fmt.Sprintf("%f", tower.GetLocation().GetLat())

	return &topodevice.Device{
		ID:          topodevice.ID(towerID),
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
