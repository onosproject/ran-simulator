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

	"github.com/onosproject/ran-simulator/pkg/utils"

	e2 "github.com/onosproject/onos-ric/api/sb"
	"gotest.tools/assert"
)

func Test_GenerateReport(t *testing.T) {
	mgr, err := setUpManager()
	assert.NilError(t, err, "Unexpected error setting up manager")
	assert.Assert(t, mgr != nil, "Unexpectedly Manager is nil!")

	mgr.UserEquipmentsLock.RLock()
	ue1, ok := mgr.UserEquipments[utils.ImsiGenerator(0)]
	assert.Assert(t, ok, "Expected to find Ue-0001")
	mgr.UserEquipmentsLock.RUnlock()

	msg, err := generateReport(ue1)
	assert.NilError(t, err)
	assert.Equal(t, e2.MessageType_RADIO_MEAS_REPORT_PER_UE, msg.GetHdr().MessageType)
	rmrpu := msg.GetMsg().GetRadioMeasReportPerUE()
	assert.Assert(t, ok, "Expected msg.S to convert to RadioMeasReportPerUE")
	assert.Equal(t, "0001", rmrpu.GetCrnti())
	assert.Equal(t, 4, len(rmrpu.GetRadioReportServCells()))
	for _, rr := range rmrpu.GetRadioReportServCells() {
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
	stopManager(mgr)
}

func Test_MakeCqi(t *testing.T) {

	testCases := []struct {
		strength float64
		cqi      uint32
	}{
		{-5, 2},
		{-4, 3},
		{-1, 6},
		{0, 7},
		{1.0, 8},
		{1.5, 9},
		{3.0, 10},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.cqi, makeCqi(tc.strength), "unexpected cqi %d for strength %f", tc.cqi, tc.strength)
	}
}
