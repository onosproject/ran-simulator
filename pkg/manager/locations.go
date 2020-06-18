// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package manager

import (
	"fmt"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"math"
	"math/rand"

	"github.com/onosproject/ran-simulator/api/types"
)

// LocationID type - an alias for string
type LocationID string

// Location :
type Location struct {
	Name     LocationID
	Position types.Point
}

// NewLocations - create a new set of locations
func NewLocations(cells map[types.ECGI]*types.Cell, maxUEs int, locationsScale float32) (*types.Point, map[LocationID]*Location) {
	locations := make(map[LocationID]*Location)

	minLat := 90.0
	maxLat := -90.0
	minLng := 180.0
	maxLng := -180.0

	for _, cell := range cells {
		if cell.GetLocation().GetLat() < minLat {
			minLat = cell.GetLocation().GetLat()
		}
		if cell.GetLocation().GetLat() > maxLat {
			maxLat = cell.GetLocation().GetLat()
		}
		if cell.GetLocation().GetLng() < minLng {
			minLng = cell.GetLocation().GetLng()
		}
		if cell.GetLocation().GetLng() > maxLng {
			maxLng = cell.GetLocation().GetLng()
		}
	}
	centre := types.Point{Lat: minLat + (maxLat-minLat)/2, Lng: minLng + (maxLng-minLng)/2}
	radius := float64(locationsScale) * math.Hypot(maxLat-minLat, maxLng-minLng) / 2
	aspectRatio := utils.AspectRatio(&centre)
	for l := 0; l < (maxUEs * 2); l++ {
		pos := utils.RandomLatLng(centre.GetLat(), centre.GetLng(),
			radius, aspectRatio)
		name := LocationID(fmt.Sprintf("Location-%d", l))
		loc := Location{
			Name:     name,
			Position: pos,
		}
		locations[name] = &loc
	}

	return &centre, locations
}

func (m *Manager) getRandomLocation(exclude LocationID) (*Location, error) {
	randomIndex := rand.Intn(len(m.Locations) - 1)
	idx := 0
	for name, loc := range m.Locations {
		if idx == randomIndex {
			if name == exclude {
				randomIndex = randomIndex + 1
				idx = idx + 1
				continue
			}
			return loc, nil
		}
		idx = idx + 1
	}
	return nil, fmt.Errorf("no location found")
}
