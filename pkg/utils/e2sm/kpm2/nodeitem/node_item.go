// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package nodeitem

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
)

// NodeItem kpm node it
type NodeItem struct {
	globalKPMNodeID     *e2smkpmv2.GlobalKpmnodeId
	cellMeasObjectItems []*e2smkpmv2.CellMeasurementObjectItem
}

// NewNodeItem creates new kpm node item
func NewNodeItem(options ...func(item *NodeItem)) *NodeItem {
	nodeItem := &NodeItem{}
	for _, option := range options {
		option(nodeItem)
	}

	return nodeItem
}

// WithGlobalKpmNodeID sets global KPM Node ID
func WithGlobalKpmNodeID(globalKPMNodeID *e2smkpmv2.GlobalKpmnodeId) func(item *NodeItem) {
	return func(item *NodeItem) {
		item.globalKPMNodeID = globalKPMNodeID
	}
}

// WithCellMeasurementObjectItems sets cell measurement object items
func WithCellMeasurementObjectItems(cellMeasObjectItems []*e2smkpmv2.CellMeasurementObjectItem) func(item *NodeItem) {
	return func(item *NodeItem) {
		item.cellMeasObjectItems = cellMeasObjectItems
	}
}

// Build builds global kpm node item
func (item *NodeItem) Build() *e2smkpmv2.RicKpmnodeItem {
	return &e2smkpmv2.RicKpmnodeItem{
		RicKpmnodeType:            item.globalKPMNodeID,
		CellMeasurementObjectList: item.cellMeasObjectItems,
	}

}
