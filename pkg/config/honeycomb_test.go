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
