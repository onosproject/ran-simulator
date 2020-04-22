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
	"github.com/onosproject/ran-simulator/api/types"
	"gotest.tools/assert"
	"testing"
)

func Test_NewLocations(t *testing.T) {
	mapLayout := types.MapLayout{
		Center:        &types.Point{Lat: 52.8, Lng: -8.2},
		Zoom:          12,
		Fade:          false,
		ShowRoutes:    false,
		ShowPower:     false,
		MinUes:        3,
		MaxUes:        30,
		CurrentRoutes: 0,
	}

	towersParams := types.TowersParams{
		TowerRows:         2,
		TowerCols:         2,
		TowerSpacingVert:  0.02,
		TowerSpacingHoriz: 0.02,
		LocationsScale:    1.0,
		MaxUEsPerCell:     4,
		AvgCellsPerTower:  3.0,
	}

	locations := NewLocations(towersParams, mapLayout)

	assert.Equal(t, 60, len(locations), "Unexpected number of locations")
	gridHalf := float32(towersParams.TowerCols-1) * towersParams.TowerSpacingHoriz / 2
	minLng := mapLayout.GetCenter().GetLng() - gridHalf
	maxLng := mapLayout.GetCenter().GetLng() + gridHalf
	for k, l := range locations {
		assert.Assert(t, l.Position.GetLng() > minLng, "%s expected lng %f to be < than maxLng %f", k, l.Position.GetLng(), maxLng)
		assert.Assert(t, l.Position.GetLng() < maxLng, "%s expected lng %f to be > than minLng %f", k, l.Position.GetLat(), minLng)
	}

}
