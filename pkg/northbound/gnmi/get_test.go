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

package gnmi

import (
	"github.com/onosproject/config-models/modelplugin/e2node-1.0.0/e2node_1_0_0"
	"github.com/onosproject/ran-simulator/api/types"
	"gotest.tools/assert"
	"testing"
)

func setUpCell() (*types.ECGI, *e2node_1_0_0.Device) {
	ecgi := types.ECGI{
		EcID:   "315010",
		PlmnID: "0001420",
	}

	var (
		PdcpMeasReportPerUe    uint32 = 21
		RadioMeasReportPerCell uint32 = 22
		RadioMeasReportPerUe   uint32 = 23
		SchedMeasReportPerCell uint32 = 24
		SchedMeasReportPerUe   uint32 = 25
	)

	device := e2node_1_0_0.Device{
		E2Node: &e2node_1_0_0.E2Node_E2Node{
			Intervals: &e2node_1_0_0.E2Node_E2Node_Intervals{
				PdcpMeasReportPerUe:    &PdcpMeasReportPerUe,
				RadioMeasReportPerCell: &RadioMeasReportPerCell,
				RadioMeasReportPerUe:   &RadioMeasReportPerUe,
				SchedMeasReportPerCell: &SchedMeasReportPerCell,
				SchedMeasReportPerUe:   &SchedMeasReportPerUe,
			},
		},
	}

	return &ecgi, &device
}

func Test_getE2nodeIntervalsPdcpMeasReportPerUe(t *testing.T) {
	ecgi, device := setUpCell()

	notif, err := getE2nodeIntervalsPdcpMeasReportPerUe(*ecgi, device)
	assert.NilError(t, err)
	assert.Equal(t, `path:<elem:<name:"e2node" > elem:<name:"intervals" > elem:<name:"PdcpMeasReportPerUe" > > val:<uint_val:21 > `, notif.Update[0].String())
}

func Test_getE2nodeIntervalsPdcpMeasReportPerCell(t *testing.T) {
	ecgi, device := setUpCell()

	notif, err := getE2nodeIntervalsRadioMeasReportPerCell(*ecgi, device)
	assert.NilError(t, err)
	assert.Equal(t, `path:<elem:<name:"e2node" > elem:<name:"intervals" > elem:<name:"RadioMeasReportPerCell" > > val:<uint_val:22 > `, notif.Update[0].String())
}
