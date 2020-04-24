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

package config

import (
	"fmt"
	configlib "github.com/onosproject/onos-lib-go/pkg/config"
	"github.com/onosproject/ran-simulator/api/types"
)

var towerConfig *TowerConfig

// Sector - one side of the tower
type Sector struct {
	EcID        types.EcID `yaml:"ecid"`
	GrpcPort    uint16     `yaml:"grpcport"`
	Azimuth     uint16     `yaml:"azimuth"`
	Arc         uint16     `yaml:"arc"`
	MaxUEs      uint16     `yaml:"maxues"`
	InitPowerDb float32    `yaml:"initpowerdb"`
}

// TowersLayout an individual tower with sectors
type TowersLayout struct {
	TowerID   string       `yaml:"towerid"`
	PlmnID    types.PlmnID `yaml:"plmnid"`
	Latitude  float32      `yaml:"latitude"`
	Longitude float32      `yaml:"longitude"`
	Sectors   []Sector
}

// TowerConfig is the ran-simulator configuration
type TowerConfig struct {
	MapCentre    types.Point
	TowersLayout []TowersLayout
}

// Clear - reset the config - needed for tests
func Clear() {
	towerConfig = nil
}

// GetTowerConfig gets the onos-towerConfig configuration
func GetTowerConfig(location string) (TowerConfig, error) {
	if towerConfig == nil {
		towerConfig = &TowerConfig{}
		if err := configlib.LoadNamedConfig(location, towerConfig); err != nil {
			return TowerConfig{}, err
		}
		if err := Checker(towerConfig); err != nil {
			return TowerConfig{}, err
		}
	}
	return *towerConfig, nil
}

// Checker - check everything is within bounds
func Checker(config *TowerConfig) error {
	if config.MapCentre.Lat < -90 || config.MapCentre.Lat > 90 {
		return fmt.Errorf("map centre latitude outside range -90, 90")
	}
	if config.MapCentre.Lng < -180 || config.MapCentre.Lng > 180 {
		return fmt.Errorf("map centre longitude outside range -180, 180")
	}

	// Highly unlikely - no cell towers in mid atlantic
	if config.MapCentre.Lat == 0 && config.MapCentre.Lng == 0 {
		return fmt.Errorf("map centre invalid: 0,0")
	}

	ecgis := make(map[string]interface{})
	grpcports := make(map[uint16]interface{})
	latlngs := make(map[string]interface{})
	for _, tower := range config.TowersLayout {
		latlng := fmt.Sprintf("%f %f", tower.Latitude, tower.Longitude)
		if _, ok := latlngs[latlng]; ok {
			return fmt.Errorf("%s lat lng repeated in %s", latlng, tower.TowerID)
		}
		latlngs[latlng] = struct{}{}

		if len(tower.PlmnID) != 6 {
			return fmt.Errorf("the PlmnID must be 6 chars: %s", tower.PlmnID)
		}

		if len(tower.Sectors) == 0 {
			return fmt.Errorf("every tower must have at least 1 sector: %s", tower.TowerID)
		}

		for _, sector := range tower.Sectors {
			if len(sector.EcID) != 7 {
				return fmt.Errorf("the Ecid must be 7 chars: %s", sector.EcID)
			}

			ecgi := fmt.Sprintf("%s-%s", tower.PlmnID, sector.EcID)
			if _, ok := ecgis[ecgi]; ok {
				return fmt.Errorf("%s ecid repeated in %s", sector.EcID, tower.TowerID)
			}
			ecgis[ecgi] = struct{}{}

			if sector.GrpcPort < 1024 || sector.GrpcPort == 5150 {
				return fmt.Errorf("invalid value for grpc port %d", sector.GrpcPort)
			}
			if _, ok := grpcports[sector.GrpcPort]; ok {
				return fmt.Errorf("%d grpcport repeated in %s", sector.GrpcPort, tower.TowerID)
			}
			grpcports[sector.GrpcPort] = struct{}{}

			if sector.Arc < 1 || sector.Arc > 360 {
				return fmt.Errorf("arc must be 1-360° %d", sector.Arc)
			}

			if sector.Azimuth > 359 {
				return fmt.Errorf("azimuth must be 0-359° %d", sector.Azimuth)
			}

			if sector.MaxUEs < 1 || sector.MaxUEs > 65535 {
				return fmt.Errorf("MaxUEs must be 1-65535° %d", sector.MaxUEs)
			}
		}
	}
	return nil
}
