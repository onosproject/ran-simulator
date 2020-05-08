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

package utils

import (
	"github.com/onosproject/ran-simulator/api/types"
	"gotest.tools/assert"
	"math"
	"testing"
)

const (
	PosCenLat = 52.12345
	PosCenLng = 13.12345
	Pos1Lat   = 52.12350
	Pos1Lng   = 13.12350
	Pos2Lat   = 52.12340
	Pos2Lng   = 13.12340
)

func Test_GetRotationDegrees(t *testing.T) {
	centre := types.Point{
		Lat: PosCenLat,
		Lng: PosCenLng,
	}
	p1 := types.Point{
		Lat: Pos1Lat,
		Lng: Pos1Lng,
	}
	p2 := types.Point{
		Lat: Pos2Lat,
		Lng: Pos2Lng,
	}
	p3 := types.Point{
		Lat: Pos2Lat,
		Lng: Pos1Lng,
	}
	p4 := types.Point{
		Lat: Pos1Lat,
		Lng: Pos2Lng,
	}
	r1 := GetRotationDegrees(&centre, &p1)
	assert.Equal(t, 45.0, math.Round(r1), "Unexpected r1")

	r2 := GetRotationDegrees(&centre, &p2)
	assert.Equal(t, -135.0, math.Round(r2), "Unexpected r2")

	r3 := GetRotationDegrees(&centre, &p3)
	assert.Equal(t, -45.0, math.Round(r3), "Unexpected r3")

	r4 := GetRotationDegrees(&centre, &p4)
	assert.Equal(t, 135.0, math.Round(r4), "Unexpected r4")
}

func Test_RandomColor(t *testing.T) {
	c := RandomColor()
	assert.Equal(t, 7, len(c))
	assert.Equal(t, uint8('#'), c[0])
}

func Test_GetRandomLngLat(t *testing.T) {
	const scale = 0.2
	for i := 0; i < 100; i++ {
		pt := RandomLatLng(0.0, 0.0, scale, 1)
		assert.Assert(t, pt.GetLat() < scale, "Expecting position %f to be within scale", pt.GetLat())
	}
}

func Test_AzimuthToRads(t *testing.T) {
	assert.Equal(t, math.Pi/2, AzimuthToRads(0))
	assert.Equal(t, 0.0, AzimuthToRads(90))
	assert.Equal(t, -math.Pi/2, AzimuthToRads(180))
	assert.Equal(t, -math.Pi, AzimuthToRads(270))
	assert.Equal(t, math.Round(10e6*math.Pi/3), math.Round(10e6*AzimuthToRads(30)))
}

func Test_AspectRatio(t *testing.T) {
	ar := AspectRatio(&types.Point{Lat: 52.52, Lng: 13.13})
	assert.Equal(t, 608, int(math.Round(ar*1e3)))
}
