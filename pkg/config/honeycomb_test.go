// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package config

import (
	"gotest.tools/assert"
	"testing"
)

func Test_HexArray7(t *testing.T) {
	points := hexMesh(0.01, 7)
	assert.Equal(t, 7, len(points))

	var testcases = []struct {
		pointidx uint
		lat      float64
		lng      float64
	}{
		{0, -0.017320508075688773, 0},
		{1, -0.008660254037844387, 0.015},
		{2, -0.008660254037844387, -0.015},
		{3, 0.0, 0},
		{4, 0.008660254037844387, 0.015},
		{5, 0.008660254037844387, -0.015},
		{6, 0.017320508075688773, 0},
	}
	for _, tc := range testcases {
		assert.Equal(t, tc.lat, points[tc.pointidx].Lat, tc)
		assert.Equal(t, tc.lng, points[tc.pointidx].Lng, tc)
	}
}

func Test_HexArrayN(t *testing.T) {
	points := hexMesh(0.01, 12)
	assert.Equal(t, 19, len(points))

	for _, p := range points {
		t.Logf("P: %v", p)
	}

	var testcases = []struct {
		pointidx uint
		lat      float64
		lng      float64
	}{
		{9, 0, 0},
	}
	for _, tc := range testcases {
		assert.Equal(t, tc.lat, points[tc.pointidx].Lat, tc)
		assert.Equal(t, tc.lng, points[tc.pointidx].Lng, tc)
	}
}
