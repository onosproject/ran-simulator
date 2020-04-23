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

package manager

import (
	"github.com/onosproject/ran-simulator/pkg/config"
	"gotest.tools/assert"
	"testing"
)

func Test_NewLocations(t *testing.T) {
	towersConfig, err := config.GetTowerConfig("berlin-rectangular-4-1.yaml")
	assert.NilError(t, err)

	locations := NewLocations(towersConfig, 30)

	assert.Equal(t, 60, len(locations), "Unexpected number of locations")

	minLat := towersConfig.MapCentre.GetLat()
	maxLat := towersConfig.MapCentre.GetLat()
	minLng := towersConfig.MapCentre.GetLng()
	maxLng := towersConfig.MapCentre.GetLng()
	for _, tower := range towersConfig.TowersLayout {
		if tower.Latitude < minLat {
			minLat = tower.Latitude
		} else if tower.Latitude > maxLat {
			maxLat = tower.Latitude
		}
		if tower.Longitude < minLng {
			minLng = tower.Longitude
		} else if tower.Longitude > maxLng {
			maxLng = tower.Longitude
		}
	}

	for k, l := range locations {
		assert.Assert(t, l.Position.GetLng() > minLng, "%s expected lng %f to be < than maxLng %f", k, l.Position.GetLng(), maxLng)
		assert.Assert(t, l.Position.GetLng() < maxLng, "%s expected lng %f to be > than minLng %f", k, l.Position.GetLat(), minLng)
	}

}
