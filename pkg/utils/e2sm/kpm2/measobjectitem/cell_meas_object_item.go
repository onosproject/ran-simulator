// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measobjectitem

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
)

// CellMeasObjectItem cell measurement object item
type CellMeasObjectItem struct {
	cellObjID    string
	cellGlobalID *e2smkpmv2.CellGlobalId
}

// NewCellMeasObjectItem create new cell measurement object item
func NewCellMeasObjectItem(options ...func(*CellMeasObjectItem)) *CellMeasObjectItem {
	cellMeasObjectItem := &CellMeasObjectItem{}
	for _, option := range options {
		option(cellMeasObjectItem)
	}

	return cellMeasObjectItem
}

// WithCellObjectID sets cell object ID
func WithCellObjectID(cellObjID string) func(item *CellMeasObjectItem) {
	return func(item *CellMeasObjectItem) {
		item.cellObjID = cellObjID
	}
}

// WithCellGlobalID sets cell global ID
func WithCellGlobalID(cellGlobalID *e2smkpmv2.CellGlobalId) func(item *CellMeasObjectItem) {
	return func(item *CellMeasObjectItem) {
		item.cellGlobalID = cellGlobalID
	}
}

// Build builds a cell measurement object item
func (c *CellMeasObjectItem) Build() *e2smkpmv2.CellMeasurementObjectItem {
	return &e2smkpmv2.CellMeasurementObjectItem{
		CellObjectId: &e2smkpmv2.CellObjectId{
			Value: c.cellObjID,
		},
		CellGlobalId: c.cellGlobalID,
	}
}
