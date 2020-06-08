// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

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
	tolerance := 1.0 / stepsPerDecimalDegree * 2
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
