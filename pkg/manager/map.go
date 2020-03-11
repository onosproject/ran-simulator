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

import "github.com/onosproject/ran-simulator/api/types"

func getTowerPosition(row uint32, col uint32, params types.TowersParams, mapLayout types.MapLayout) *types.Point {
	topLeft := getMapTopLeft(params, mapLayout)

	return &types.Point{
		Lat: topLeft.Lat - params.TowerSpacingVert*float32(row),
		Lng: topLeft.Lng + params.TowerSpacingHoriz*float32(col),
	}
}

func getMapTopLeft(params types.TowersParams, mapLayout types.MapLayout) *types.Point {
	return &types.Point{
		Lat: mapLayout.GetCenter().GetLat() + params.TowerSpacingVert*float32(params.TowerRows-1)/2,
		Lng: mapLayout.GetCenter().GetLng() - params.TowerSpacingHoriz*float32(params.TowerCols-1)/2,
	}
}
