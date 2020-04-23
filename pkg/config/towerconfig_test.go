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

package config

import (
	"gotest.tools/assert"
	"testing"
)

func Test_LoadConfig(t *testing.T) {
	config, err := GetTowerConfig("berlin-honeycomb-4-3.yaml")
	assert.NilError(t, err, "Unexpected error loading towerConfig")
	assert.Equal(t, 4, len(config.TowersLayout), "Unexpected number of towers")

	tower1 := config.TowersLayout[0]
	assert.Equal(t, "Tower-1", tower1.TowerID)
	assert.Equal(t, 3, len(tower1.Sectors))
}