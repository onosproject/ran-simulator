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

func Test_RandomRoute(t *testing.T) {
	startLoc := Location{Position: types.Point{Lat: 52.0, Lng: -8.0}}
	endLoc := Location{Position: types.Point{Lat: 52.2, Lng: -8.1}}

	points, err := randomRoute(&startLoc, &endLoc)
	assert.NilError(t, err, "Unexpected error generating random route")

	assert.Equal(t, len(points), 112)
	prevLat := startLoc.Position.GetLat()
	prevLng := startLoc.Position.GetLng()
	tolerance := float32(1) / stepsPerDecimalDegree * 2
	for i, p := range points {
		//t.Logf("Point %d: %f, %f", i, p.GetLat(), p.GetLng())
		if i > 0 {
			assert.Assert(t, p.GetLat() > prevLat-tolerance, "Expected Lat #%d: %f to exceed prev Lat %f", i, p.GetLat(), prevLat)
			prevLat = p.GetLat()
			assert.Assert(t, p.GetLng() < prevLng+tolerance, "Expected Lng #%d: %f to be less than prev Lng %f", i, p.GetLng(), prevLng)
			prevLng = p.GetLng()
		}
	}
}
