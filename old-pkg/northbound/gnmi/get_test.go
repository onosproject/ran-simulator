// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package gnmi

import (
	"testing"

	"github.com/onosproject/config-models/modelplugin/e2node-1.0.0/e2node_1_0_0"
	"github.com/onosproject/onos-topo/pkg/bulk"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"gotest.tools/assert"
)

func setUpCell() (*manager.Manager, error) {
	topoConfig, err := bulk.GetTopoConfig("berlin-honeycomb-4-3-topo.yaml")
	if err != nil {
		return nil, err
	}
	mgr, err := manager.NewManager()
	if err != nil {
		return nil, err
	}

	for _, td := range topoConfig.TopoEntities {
		td := td //pin
		cell, err := manager.NewCell(bulk.TopoEntityToTopoObject(&td))
		if err != nil {
			return nil, err
		}
		mgr.Cells[*cell.Ecgi] = cell
	}
	mgr.CellConfigs = make(map[types.ECGI]*e2node_1_0_0.Device)

	var (
		PdcpMeasReportPerUe    uint32 = 21
		RadioMeasReportPerCell uint32 = 22
		RadioMeasReportPerUe   uint32 = 23
		SchedMeasReportPerCell uint32 = 24
	)

	for ecgi := range mgr.Cells {
		mgr.CellConfigs[ecgi] = &e2node_1_0_0.Device{
			E2Node: &e2node_1_0_0.E2Node_E2Node{
				Intervals: &e2node_1_0_0.E2Node_E2Node_Intervals{
					PdcpMeasReportPerUe:    &PdcpMeasReportPerUe,
					RadioMeasReportPerCell: &RadioMeasReportPerCell,
					RadioMeasReportPerUe:   &RadioMeasReportPerUe,
					SchedMeasReportPerCell: &SchedMeasReportPerCell,
				},
			},
		}
	}

	return mgr, nil
}

func Test_getE2nodeIntervalsPdcpMeasReportPerUe(t *testing.T) {
	t.Skip()
	mgr, err := setUpCell()
	assert.NilError(t, err)
	ecgi1420 := types.ECGI{
		EcID:   "0001420",
		PlmnID: "315010",
	}
	cellConfig1420, ok := mgr.CellConfigs[ecgi1420]
	assert.Assert(t, ok, "cannot find config for", ecgi1420)

	notif, err := getE2nodeIntervalsPdcpMeasReportPerUe(ecgi1420, cellConfig1420)
	assert.NilError(t, err)
	assert.Equal(t, `path:<elem:<name:"e2node" > elem:<name:"intervals" > elem:<name:"PdcpMeasReportPerUe" > > val:<uint_val:21 > `, notif.Update[0].String())
}

func Test_getE2nodeIntervalsPdcpMeasReportPerCell(t *testing.T) {
	t.Skip()
	mgr, err := setUpCell()
	assert.NilError(t, err)
	ecgi1420 := types.ECGI{
		EcID:   "0001420",
		PlmnID: "315010",
	}
	cellConfig1420, ok := mgr.CellConfigs[ecgi1420]
	assert.Assert(t, ok, "cannot find config for %s", ecgi1420)

	notif, err := getE2nodeIntervalsRadioMeasReportPerCell(ecgi1420, cellConfig1420)
	assert.NilError(t, err)
	assert.Equal(t, `path:<elem:<name:"e2node" > elem:<name:"intervals" > elem:<name:"RadioMeasReportPerCell" > > val:<uint_val:22 > `, notif.Update[0].String())
}

func Test_Get_simple(t *testing.T) {
	_, err := setUpCell()
	assert.NilError(t, err)

	cellServer := Server{
		plmnID:    "315010",
		towerEcID: "0001420",
		port:      5152,
	}

	getRequest := &gnmi.GetRequest{
		Prefix: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "e2node"},
				{Name: "intervals"},
			},
		},
		Path: []*gnmi.Path{
			{
				Elem: []*gnmi.PathElem{
					{Name: "RadioMeasReportPerCell"},
				},
			},
			{
				Elem: []*gnmi.PathElem{
					{Name: "RadioMeasReportPerUe"},
				},
			},
		},
	}

	getResponse, err := cellServer.Get(context.Background(), getRequest)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(getResponse.Notification))
}
