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
	"github.com/OpenNetworkingFoundation/gmap-ran/api/types"
	"math"
	"math/rand"
)

/**
 * Generates a random latlng value in 1000 meter radius of loc
 */
func randomLatLng(mapCenterLat float32, mapCenterLng float32) types.Point {
	const r = 5000 / float64(111300) // = 100 meters
	y0 := float64(mapCenterLat)
	x0 := float64(mapCenterLng)

	u := rand.Float64()
	v := rand.Float64()

	w := r * math.Sqrt(u);
	t := 2 * math.Pi * v;
	x := w * math.Cos(t);
	y1 := w * math.Sin(t);
	x1 := x / math.Cos(y0);

	newY := roundToDecimal(y0+y1, 6)
	newX := roundToDecimal(x0+x1, 6)
	return types.Point{
		Lat: newY,
		Lng: newX,
	}
}

/**
 * Rounds number to decimals
 */
func roundToDecimal(value float64, decimals int) float32 {
	intValue := value * math.Pow10(decimals)
	return float32(math.Round(intValue) / math.Pow10(decimals))
}

func getRotationDegrees(pointA *types.Point, pointB *types.Point) float64 {
	deltaX := float64(pointB.GetLng() - pointA.GetLng())
	deltaY := float64(pointB.GetLat() - pointA.GetLat())

	return math.Atan2(deltaY, deltaX) * 180 / math.Pi
}

func randomColor() string {
	const letters = "0123456789ABCDEF";
	color := make([]uint8, 7)
	color[0] = '#'
	for i := range color {
		if i == 0 {
			continue
		}
		color[i] = letters[rand.Intn(15)]
	}
	return string(color)
}
