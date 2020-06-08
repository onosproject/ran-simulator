// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

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
	s := Server{}

	testReq := e2.RRMConfig{
		Ecgi: &e2.ECGI{
			PlmnId: utils.TestPlmnID,
			Ecid:   "0001420",
		},
		PA: []e2.XICICPA{e2.XICICPA_XICIC_PA_DB_MINUS3},
	}

	go func() {
		select {
		case updateEvent := <-mgr.CellsChannel:
			assert.Equal(t, trafficsim.Type_UPDATED, updateEvent.Type)
			assert.Equal(t, trafficsim.UpdateType_NOUPDATETYPE, updateEvent.UpdateType)
			cell, ok := updateEvent.Object.(*types.Cell)
			assert.Assert(t, ok, "Problem converting event object to Tower %v", updateEvent.Object)
			assert.Equal(t, types.EcID("0001420"), cell.Ecgi.EcID)
			mgr.CellsLock.RLock()
			assert.Equal(t, 7.0, cell.GetTxPowerdB())
			mgr.CellsLock.RUnlock()
		case <-time.After(time.Millisecond * 100):
			t.Errorf("Timed out on Test_HandleRrmConfig")
		}
	}()

	s.handleRRMConfig(&testReq)
	time.Sleep(time.Millisecond * 110)
	stopManager(mgr)
}
