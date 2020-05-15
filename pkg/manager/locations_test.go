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
	"github.com/onosproject/onos-topo/pkg/bulk"
	"github.com/onosproject/ran-simulator/api/types"
	"gotest.tools/assert"
	"math"
	"testing"
)

func Test_NewLocations2(t *testing.T) {
	topoDeviceConfig, err := bulk.GetDeviceConfig("berlin-honeycomb-4-3-topo.yaml")
	assert.NilError(t, err)

	cells := make(map[types.ECGI]*types.Cell)

	for _, td := range topoDeviceConfig.TopoDevices {
		td := td //pin
		cell, err := NewCell(&td)
		assert.NilError(t, err)
		cells[*cell.Ecgi] = cell
	}

	centre, locations := NewLocations(cells, 30, 0.99)
	assert.Equal(t, 5250268.0, math.Round(centre.GetLat()*1e5))
	assert.Equal(t, 1340500.0, math.Round(centre.GetLng()*1e5))
	assert.Equal(t, 60, len(locations), "Unexpected number of locations")

	minLat := centre.GetLat()
	maxLat := centre.GetLat()
	minLng := centre.GetLng()
	maxLng := centre.GetLng()
	for _, c := range cells {
		if c.GetLocation().GetLat() < minLat {
			minLat = c.GetLocation().GetLat()
		}
		if c.GetLocation().GetLat() > maxLat {
			maxLat = c.GetLocation().GetLat()
		}
		if c.GetLocation().GetLng() < minLng {
			minLng = c.GetLocation().GetLng()
		}
		if c.GetLocation().GetLng() > maxLng {
			maxLng = c.GetLocation().GetLng()
		}
	}

	for k, l := range locations {
		assert.Assert(t, l.Position.GetLng() > minLng-0.1, "%s expected lng %f to be > than minLng %f", k, l.Position.GetLng(), minLng-0.1)
		assert.Assert(t, l.Position.GetLng() < maxLng+0.1, "%s expected lng %f to be < than maxLng %f", k, l.Position.GetLng(), maxLng+0.1)
	}

}
