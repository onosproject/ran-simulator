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
	"fmt"
	"github.com/onosproject/ran-simulator/pkg/config"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"math"
	"math/rand"

	"github.com/onosproject/ran-simulator/api/types"
)

// Location :
type Location struct {
	Name     string
	Position types.Point
}

// NewLocations - create a new set of locations
func NewLocations(towersConfig config.TowerConfig, maxUEs int, locationsScale float32) map[string]*Location {
	locations := make(map[string]*Location)

	minLat := towersConfig.MapCentre.GetLat()
	maxLat := towersConfig.MapCentre.GetLat()
	minLng := towersConfig.MapCentre.GetLng()
	maxLng := towersConfig.MapCentre.GetLng()
	for _, tower := range towersConfig.TowersLayout {
		if tower.Latitude < minLat {
			minLat = tower.Latitude
		} else if tower.Latitude > maxLat {
			maxLat = tower.Latitude
		}
		if tower.Longitude < minLng {
			minLng = tower.Longitude
		} else if tower.Longitude > maxLng {
			maxLng = tower.Longitude
		}
	}
	radius := float64(locationsScale) * math.Hypot(float64(maxLat-minLat), float64(maxLng-minLng)) / 2
	aspectRatio := utils.AspectRatio(&towersConfig.MapCentre)
	for l := 0; l < (maxUEs * 2); l++ {
		pos := utils.RandomLatLng(towersConfig.MapCentre.GetLat(), towersConfig.MapCentre.GetLng(),
			radius, aspectRatio)
		name := fmt.Sprintf("Location-%d", l)
		loc := Location{
			Name:     name,
			Position: pos,
		}
		locations[name] = &loc
	}

	return locations
}

func (m *Manager) getRandomLocation(exclude string) (*Location, error) {
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
