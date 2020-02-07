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
	"math/rand"

	"github.com/onosproject/ran-simulator/api/types"
)

// LocationsParams :
type LocationsParams struct {
	NumLocations int
}

// Location :
type Location struct {
	Name     string
	Position types.Point
}

func newLocations(params LocationsParams, towersParams types.TowersParams, mapLayout types.MapLayout) map[string]*Location {
	locations := make(map[string]*Location)

	for l := 0; l < params.NumLocations; l++ {
		pos := randomLatLng(mapLayout.Center.GetLat(), mapLayout.GetCenter().GetLng())
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
