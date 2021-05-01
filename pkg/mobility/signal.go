// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

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
	return cell.TxPowerDB + distAtt + angleAtt
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
	return 10 * math.Log10(1-(angularOffset/math.Pi/angleScaling))
}
