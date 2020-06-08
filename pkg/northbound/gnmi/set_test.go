// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package gnmi

import (
	"context"
	"github.com/onosproject/config-models/modelplugin/e2node-1.0.0/e2node_1_0_0"
	"github.com/onosproject/onos-topo/pkg/bulk"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/openconfig/gnmi/proto/gnmi"
	"gotest.tools/assert"
	"testing"
)

func Test_Set_simple(t *testing.T) {
	topoDeviceConfig, err := bulk.GetDeviceConfig("berlin-honeycomb-4-3-topo.yaml")
	assert.NilError(t, err)

	mgr, err := manager.NewManager()
	assert.NilError(t, err)
	for _, td := range topoDeviceConfig.TopoDevices {
		td := td //pin
		cell, err := manager.NewCell(&td)
		assert.NilError(t, err)
		mgr.Cells[*cell.Ecgi] = cell
	}
	mgr.CellConfigs = make(map[types.ECGI]*e2node_1_0_0.Device)
	for ecgi := range mgr.Cells {
		mgr.CellConfigs[ecgi] = &e2node_1_0_0.Device{
			E2Node: &e2node_1_0_0.E2Node_E2Node{
				Intervals: &e2node_1_0_0.E2Node_E2Node_Intervals{},
			},
		}
	}

	cellServer := Server{
		plmnID:    "315010",
		towerEcID: "0001420",
		port:      5152,
	}

	setRequest := &gnmi.SetRequest{
		Prefix: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "e2node"},
				{Name: "intervals"},
			},
		},
		Update: []*gnmi.Update{
			{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "RadioMeasReportPerCell"},
					},
				},
				Val: &gnmi.TypedValue{
					Value: &gnmi.TypedValue_UintVal{
						UintVal: 20,
					},
				},
			},
			{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "RadioMeasReportPerUe"},
					},
				},
				Val: &gnmi.TypedValue{
					Value: &gnmi.TypedValue_UintVal{
						UintVal: 21,
					},
				},
			},
		},
	}

	resp, err := cellServer.Set(context.Background(), setRequest)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(resp.Response))

}
