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

package bulk

import (
	"fmt"

	configlib "github.com/onosproject/onos-lib-go/pkg/config"
	"github.com/onosproject/onos-topo/api/device"
	deviceservice "github.com/onosproject/onos-topo/pkg/northbound/device"
)

var deviceConfig *DeviceConfig

// DeviceConfig - a wrapper around multiple devices
type DeviceConfig struct {
	TopoDevices []device.Device
}

// Clear - reset the config - needed for tests
func Clear() {
	deviceConfig = nil
	topoConfig = nil
}

// GetDeviceConfig gets the onos-topo configuration
func GetDeviceConfig(location string) (DeviceConfig, error) {
	if deviceConfig == nil {
		deviceConfig = &DeviceConfig{}
		if err := configlib.LoadNamedConfig(location, deviceConfig); err != nil {
			return DeviceConfig{}, err
		}
		if err := Checker(deviceConfig); err != nil {
			return DeviceConfig{}, err
		}
	}
	return *deviceConfig, nil
}

// Checker - check everything is within bounds
func Checker(config *DeviceConfig) error {
	if len(config.TopoDevices) == 0 {
		return fmt.Errorf("no devices found")
	}

	for _, dev := range config.TopoDevices {
		dev := dev // pin
		err := deviceservice.ValidateDevice(&dev)
		if err != nil {
			return err
		}
	}

	return nil
}
