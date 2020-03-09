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
	"gotest.tools/assert"
	"testing"
)

func Test_GenerateReport(t *testing.T) {
	mgr, err := setUpManager()
	assert.NilError(t, err, "Unexpected error setting up manager")
	assert.Assert(t, mgr != nil, "Unexpectedly Manager is nil!")

	mgr.UserEquipmentsLock.RLock()
	ue1, ok := mgr.UserEquipments["Ue-0001"]
	assert.Assert(t, ok, "Expected to find Ue-0001")
	mgr.UserEquipmentsLock.RUnlock()

	msg := generateReport(ue1)
	assert.Equal(t, e2.MessageType_RADIO_MEAS_REPORT_PER_UE, msg.MessageType)
	rmrpu, ok := msg.S.(*e2.TelemetryMessage_RadioMeasReportPerUE)
	assert.Assert(t, ok, "Expected msg.S to convert to RadioMeasReportPerUE")
	assert.Equal(t, "0001", rmrpu.RadioMeasReportPerUE.Crnti)
	assert.Equal(t, 3, len(rmrpu.RadioMeasReportPerUE.RadioReportServCells))
	for _, rr := range rmrpu.RadioMeasReportPerUE.RadioReportServCells {
		switch ecid := rr.Ecgi.Ecid; ecid {
		case "0001420":
		case "0001421":
		case "0001422":
		case "0001423":
			// ok
		default:
			t.Errorf("Unexpected Ecid %s in report", ecid)
		}
	}

	mgr.Close()
}
