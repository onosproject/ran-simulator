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
	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"gotest.tools/assert"
	"testing"
	"time"
)

func Test_CrntiToName(t *testing.T) {
	assert.Equal(t, "Ue-0019", crntiToName("0019"))
	assert.Equal(t, "Ue-1234567890ABCDE", crntiToName("1234567890ABCDE"))
}

func Test_EciToName(t *testing.T) {
	assert.Equal(t, "Tower-1", eciToName("1"))
	assert.Equal(t, "Tower-0", eciToName("1234567890ABCDE"))
}

func Test_HandleRrmConfig(t *testing.T) {
	mgr, err := setUpManager()
	assert.NilError(t, err, "Unexpected error setting up manager")
	assert.Assert(t, mgr != nil, "Unexpectedly Manager is nil!")

	testReq := e2.RRMConfig{
		Ecgi: &e2.ECGI{
			PlmnId: manager.TestPlmnID,
			Ecid:   "1",
		},
		PA: []e2.XICICPA{e2.XICICPA_XICIC_PA_DB_MINUS3},
	}

	go func() {
		select {
		case updateEvent := <-mgr.TowerChannel:
			assert.Equal(t, trafficsim.Type_UPDATED, updateEvent.Type)
			assert.Equal(t, trafficsim.UpdateType_NOUPDATETYPE, updateEvent.UpdateType)
			tower, ok := updateEvent.Object.(*types.Tower)
			assert.Assert(t, ok, "Problem converting event object to Tower")
			assert.Equal(t, "0000001", tower.EcID)
			assert.Equal(t, 7, tower.GetTxPowerdB())
		case <-time.After(3 * time.Second):
			t.Errorf("Timed out on Test_HandleRrmConfig")
		}
	}()

	handleRRMConfig(&testReq)

	time.Sleep(time.Millisecond * 5)
}
