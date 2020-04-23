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
	"fmt"
	"github.com/onosproject/ran-simulator/api/types"
	"math"
	"math/rand"
)

// ServerParams - params to start a new server
type ServerParams struct {
	CaPath       string
	KeyPath      string
	CertPath     string
	TopoEndpoint string
}

// GrpcBasePort - the base port for trafficsim - other e2 ports are stepped from this
const GrpcBasePort = 5150

// ServiceName is the default name of this Kubernetes service
const ServiceName = "ran-simulator"

// TestPlmnID - https://en.wikipedia.org/wiki/Mobile_country_code#Test_networks
const TestPlmnID = "315010"

// ImsiBaseCbrs - from https://imsiadmin.com/cbrs-assignments
const ImsiBaseCbrs = types.Imsi(315010999900000)

// RandomLatLng - Generates a random latlng value in 1000 meter radius of loc
func RandomLatLng(mapCenterLat float32, mapCenterLng float32, radius float64, aspectRatio float64) types.Point {
	var r = float64(radius)
	y0 := float64(mapCenterLat)
	x0 := float64(mapCenterLng)

	u := rand.Float64()
	v := rand.Float64()

	w := r * math.Sqrt(u)
	t := 2 * math.Pi * v
	x1 := w * math.Cos(t) * float64(aspectRatio)
	y1 := w * math.Sin(t)

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

// GetRotationDegrees - get the rotation of the car
func GetRotationDegrees(pointA *types.Point, pointB *types.Point) float64 {
	deltaX := float64(pointB.GetLng() - pointA.GetLng())
	deltaY := float64(pointB.GetLat() - pointA.GetLat())

	return math.Atan2(deltaY, deltaX) * 180 / math.Pi
}

// RandomColor from https://htmlcolorcodes.com/
func RandomColor() string {
	colorPalette := []string{
		"#641E16",
		"#78281F",
		"#512E5F",
		"#4A235A",
		"#154360",
		"#1B4F72",
		"#0E6251",
		"#0B5345",
		"#145A32",
		"#186A3B",
		"#7D6608",
		"#7E5109",
		"#784212",
		"#6E2C00",
		"#7B7D7D",
		"#626567",
		"#4D5656",
		"#424949",
		"#1B2631",
		"#17202A",

		"#C0392B",
		"#E74C3C",
		"#9B59B6",
		"#8E44AD",
		"#2980B9",
		"#3498DB",
		"#1ABC9C",
		"#16A085",
		"#27AE60",
		"#2ECC71",
		"#F1C40F",
		"#F39C12",
		"#E67E22",
		"#D35400",
		"#B3B6B7",
		"#BDC3C7",
		"#95A5A6",
		"#7F8C8D",
		"#34495E",
		"#2C3E50",
	}
	return colorPalette[rand.Intn(39)]
}

// EcIDForPort gives a consistent naming convention
func EcIDForPort(cellPort int) types.EcID {
	return types.EcID(fmt.Sprintf("%07X", cellPort))
}

// ImsiGenerator -- generate an Imsi from an index
func ImsiGenerator(ueIdx int) types.Imsi {
	return ImsiBaseCbrs + types.Imsi(ueIdx) + 1
}
