// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package utils

import (
	"github.com/onosproject/ran-simulator/pkg/model"
	"math"
)

// Earth radius in meters
const earthRadius = 6378100

// See: http://en.wikipedia.org/wiki/Haversine_formula

// Distance returns the distance in meters between two geo coordinates
func Distance(c1 model.Coordinate, c2 model.Coordinate) float64 {
	var la1, lo1, la2, lo2 float64
	la1 = c1.Lat * math.Pi / 180
	lo1 = c1.Lng * math.Pi / 180
	la2 = c2.Lat * math.Pi / 180
	lo2 = c2.Lng * math.Pi / 180

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}

// TargetPoint returns the target coordinate specified distance and heading from the starting coordinate
func TargetPoint(c model.Coordinate, bearing float64, dist float64) model.Coordinate {
	var la1, lo1, la2, lo2, azimuth, d float64
	la1 = c.Lat * math.Pi / 180
	lo1 = c.Lng * math.Pi / 180
	azimuth = bearing * math.Pi / 180
	d = dist / earthRadius

	la2 = math.Asin(math.Sin(la1)*math.Cos(d) + math.Cos(la1)*math.Sin(d)*math.Cos(azimuth))
	lo2 = lo1 + math.Atan2(math.Sin(azimuth)*math.Sin(d)*math.Cos(la1), math.Cos(d)-math.Sin(la1)*math.Sin(la2))

	return model.Coordinate{Lat: la2 * 180 / math.Pi, Lng: lo2 * 180 / math.Pi}
}

// InitialBearing returns initial bearing from c1 to c2
func InitialBearing(c1 model.Coordinate, c2 model.Coordinate) float64 {
	var la1, lo1, la2, lo2 float64
	la1 = c1.Lat * math.Pi / 180
	lo1 = c1.Lng * math.Pi / 180
	la2 = c2.Lat * math.Pi / 180
	lo2 = c2.Lng * math.Pi / 180

	y := math.Sin(lo2-lo1) * math.Cos(la2)
	x := math.Cos(la1)*math.Sin(la2) - math.Sin(la1)*math.Cos(la2)*math.Cos(lo2-lo1)
	theta := math.Atan2(y, x)
	return math.Mod(theta*180/math.Pi+360, 360.0) // in degrees
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
