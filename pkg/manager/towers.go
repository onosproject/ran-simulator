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
	"github.com/OpenNetworkingFoundation/gmap-ran/api/types"
)

type TowersParams struct {
	MapCenterLat float32;
	MapCenterLng float32;
	TowerRows int;
	TowerCols int;
	TowerSpacingVert float32;
	TowerSpacingHoriz float32;
}

type TowerIf interface {
	GetPosition() types.Point
}

func NewTowers(params TowersParams) (map[string]*types.Tower) {
	topLeft := types.Point{
		Lat: params.MapCenterLat + params.TowerSpacingVert * float32(params.TowerRows) / 2,
		Lng: params.MapCenterLng - params.TowerSpacingHoriz * float32(params.TowerCols) / 2,
	}
	var towerNum = 0
	towers := make(map[string]*types.Tower)

	for r := 0; r < params.TowerRows; r++ {
		for c := 0; c < params.TowerCols; c++ {
			pos := types.Point{
				Lat: topLeft.Lat - 0.03 * float32(r),
				Lng: topLeft.Lng + 0.05 * float32(c),
			}
			towerNum = towerNum + 1
			towerName := fmt.Sprintf("Tower-%d", towerNum)
			towers[towerName] = &types.Tower{
				Name:     towerName,
				Location: &pos,
			}
		}
	}

	return towers
}
