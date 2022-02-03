// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package utils

import (
	"math"
	"math/rand"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
)

// ServerParams - params to start a new server
type ServerParams struct {
	CaPath         string
	KeyPath        string
	CertPath       string
	TopoEndpoint   string
	AddK8sSvcPorts bool
}

// GrpcBasePort - the base port for trafficsim - other e2 ports are stepped from this
const GrpcBasePort = 5150

// ServiceName is the default name of this Kubernetes service
const ServiceName = "ran-simulator"

// TestPlmnID - https://en.wikipedia.org/wiki/Mobile_country_code#Test_networks
const TestPlmnID = "315010"

// ImsiBaseCbrs - from https://imsiadmin.com/cbrs-assignments
const ImsiBaseCbrs = types.IMSI(315010999900000)

// RandomLatLng - Generates a random latlng value in 1000 meter radius of loc
func RandomLatLng(mapCenterLat float64, mapCenterLng float64, radius float64, aspectRatio float64) types.Point {
	var r = radius
	y0 := mapCenterLat
	x0 := mapCenterLng

	u := rand.Float64()
	v := rand.Float64()

	w := r * math.Sqrt(u)
	t := 2 * math.Pi * v
	x1 := w * math.Cos(t) / aspectRatio
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
func roundToDecimal(value float64, decimals int) float64 {
	intValue := value * math.Pow10(decimals)
	return math.Round(intValue) / math.Pow10(decimals)
}

// GetRotationDegrees - get the rotation of the car
func GetRotationDegrees(pointA *types.Point, pointB *types.Point) float64 {
	deltaX := pointB.GetLng() - pointA.GetLng()
	deltaY := pointB.GetLat() - pointA.GetLat()

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

// ImsiGenerator -- generate an Imsi from an index
func ImsiGenerator(ueIdx int) types.IMSI {
	return ImsiBaseCbrs + types.IMSI(ueIdx) + 1
}

// AzimuthToRads - angle measured in degrees clockwise from north, expressed in rads from 3 o'clock anticlockwise
func AzimuthToRads(azimuth float64) float64 {
	if azimuth == 90 {
		return 0
	}
	return DegreesToRads(90 - azimuth)
}

// DegreesToRads - general conversion of degrees to rads, both starting at 3 o'clock going anticlockwise
func DegreesToRads(degrees float64) float64 {
	return 2 * math.Pi * degrees / 360
}

// AspectRatio - Compensate for the narrowing of meridians at higher latitudes
func AspectRatio(latitude float64) float64 {
	return math.Cos(DegreesToRads(latitude))
}

// Uint64ToBitString converts uint64 to a bit string byte array
func Uint64ToBitString(value uint64, bitCount int) []byte {
	result := make([]byte, bitCount/8+1)
	if bitCount%8 > 0 {
		value = value << (8 - bitCount%8)
	}

	for i := 0; i <= (bitCount / 8); i++ {
		result[i] = byte(value >> (((bitCount / 8) - i) * 8) & 0xFF)
	}

	return result
}

// BitStringToUint64 converts bit string to uint 64
func BitStringToUint64(bitString []byte, bitCount int) uint64 {
	var result uint64
	for i, b := range bitString {
		result += uint64(b) << ((len(bitString) - i - 1) * 8)
	}
	if bitCount%8 != 0 {
		return result >> (8 - bitCount%8)
	}
	return result
}
