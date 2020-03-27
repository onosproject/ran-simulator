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

package e2

import (
	"testing"
	"time"

	"github.com/onosproject/ran-simulator/pkg/utils"

	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"gotest.tools/assert"
)

func Test_EciToName(t *testing.T) {
	mgr, err := setUpManager()
	assert.NilError(t, err, "Unexpected error setting up manager")
	assert.Assert(t, mgr != nil, "Unexpectedly Manager is nil!")
}

func Test_HandleRrmConfig(t *testing.T) {
	mgr, err := setUpManager()
	assert.NilError(t, err, "Unexpected error setting up manager")
	assert.Assert(t, mgr != nil, "Unexpectedly Manager is nil!")

	testReq := e2.RRMConfig{
		Ecgi: &e2.ECGI{
			PlmnId: utils.TestPlmnID,
			Ecid:   "0001420",
		},
		PA: []e2.XICICPA{e2.XICICPA_XICIC_PA_DB_MINUS3},
	}

	go func() {
		select {
		case updateEvent := <-mgr.TowerChannel:
			assert.Equal(t, trafficsim.Type_UPDATED, updateEvent.Type)
			assert.Equal(t, trafficsim.UpdateType_NOUPDATETYPE, updateEvent.UpdateType)
			tower, ok := updateEvent.Object.(*types.Tower)
			assert.Assert(t, ok, "Problem converting event object to Tower %v", updateEvent.Object)
			assert.Equal(t, types.EcID("0001420"), tower.Ecgi.EcID)
			mgr.TowersLock.RLock()
			assert.Equal(t, float32(7), tower.GetTxPowerdB())
			mgr.TowersLock.RUnlock()
		case <-time.After(time.Millisecond * 100):
			t.Errorf("Timed out on Test_HandleRrmConfig")
		}
	}()

	handleRRMConfig(&testReq)
	time.Sleep(time.Millisecond * 110)

	mgr.Close()
}
