// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mobility

import (
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"math"
)

// powerFactor relates power to distance in decimal degrees
const powerFactor = 0.001

// StrengthAtLocation returns the signal strength at location relative to the specified cell.
func StrengthAtLocation(coord model.Coordinate, cell model.Cell) float64 {
	distAtt := distanceAttenuation(coord, cell)
	angleAtt := angleAttenuation(coord, cell)
	pathLoss := getPathLoss(coord, cell)
	return cell.TxPowerDB + distAtt + angleAtt - pathLoss
}

// distanceAttenuation is the antenna Gain as a function of the dist
// a very rough approximation to take in to account the width of
// the antenna beam. A 120° wide beam with 30° height will span ≅ 2x0.5 = 1 steradians
// A 60° wide beam will be half that and so will have double the gain
// https://en.wikipedia.org/wiki/Sector_antenna
// https://en.wikipedia.org/wiki/Steradian
func distanceAttenuation(coord model.Coordinate, cell model.Cell) float64 {
	latDist := coord.Lat - cell.Sector.Center.Lat
	realLngDist := (coord.Lng - cell.Sector.Center.Lng) / utils.AspectRatio(cell.Sector.Center.Lat)
	r := math.Hypot(latDist, realLngDist)
	gain := 120.0 / float64(cell.Sector.Arc)
	return 10 * math.Log10(gain*math.Sqrt(powerFactor/r))
}

// angleAttenuation is the attenuation of power reaching a UE due to its
// position off the centre of the beam in dB
// It is an approximation of the directivity of the antenna
// https://en.wikipedia.org/wiki/Radiation_pattern
// https://en.wikipedia.org/wiki/Sector_antenna
func angleAttenuation(coord model.Coordinate, cell model.Cell) float64 {
	azRads := utils.AzimuthToRads(float64(cell.Sector.Azimuth))
	pointRads := math.Atan2(coord.Lat-cell.Sector.Center.Lat, coord.Lng-cell.Sector.Center.Lng)
	angularOffset := math.Abs(azRads - pointRads)
	angleScaling := float64(cell.Sector.Arc) / 120.0 // Compensate for narrower beams

	// We just use a simple linear formula 0 => no loss
	// 33° => -3dB for a 120° sector according to [2]
	// assume this is 1:1 rads:attenuation e.g. 0.50 rads = 0.5 = -3dB attenuation
	//return 10 * math.Log10(1-(angularOffset/math.Pi/angleScaling))
	return -math.Min(12*math.Pow((angularOffset/(math.Pi*2/3)/angleScaling), 2), 30)
}

func getPathLoss(coord model.Coordinate, cell model.Cell) float64 {
	return getFreeSpacePathLoss(coord, cell)
}

func getFreeSpacePathLoss(coord model.Coordinate, cell model.Cell) float64 {
	distanceKM := getEuclianDistanceFromGPS(coord, cell)
	// Assuming we're using CBRS frequency 3.6 GHz
	// 92.45 is the constant value of 20 * log10(4*pi / c) in Kilometer scale
	pathLoss := 20*math.Log10(distanceKM) + 20*math.Log10(3.6) + 92.45
	return pathLoss
}

func getEuclianDistanceFromGPS(coord model.Coordinate, cell model.Cell) float64 {
	earthRadius := 6378.137
	dLat := coord.Lat*math.Pi/180 - cell.Sector.Center.Lat*math.Pi/180
	dLng := coord.Lng*math.Pi/180 - cell.Sector.Center.Lng*math.Pi/180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(coord.Lat*math.Pi/180)*math.Cos(cell.Sector.Center.Lat*math.Pi/180)*
		math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
