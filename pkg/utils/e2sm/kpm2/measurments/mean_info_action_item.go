// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measurments

import e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"

// MeasurementInfoActionItem measurement info action item
type MeasurementInfoActionItem struct {
	measTypeName string
	measTypeID   int32
}

// NewMeasurementInfoActionItem creates a new measurement info action item
func NewMeasurementInfoActionItem(options ...func(item *MeasurementInfoActionItem)) *MeasurementInfoActionItem {
	measInfoActionItem := &MeasurementInfoActionItem{}
	for _, option := range options {
		option(measInfoActionItem)
	}

	return measInfoActionItem
}

// WithMeasTypeName sets measurement type name
func WithMeasTypeName(measTypeName string) func(item *MeasurementInfoActionItem) {
	return func(item *MeasurementInfoActionItem) {
		item.measTypeName = measTypeName
	}
}

// WithMeasTypeID sets measurement type ID
func WithMeasTypeID(measTypeID int32) func(item *MeasurementInfoActionItem) {
	return func(item *MeasurementInfoActionItem) {
		item.measTypeID = measTypeID
	}
}

// Build builds measurement info action item
func (m *MeasurementInfoActionItem) Build() (*e2smkpmv2.MeasurementInfoActionItem, error) {
	return &e2smkpmv2.MeasurementInfoActionItem{
		MeasName: &e2smkpmv2.MeasurementTypeName{
			Value: m.measTypeName,
		},
		MeasId: &e2smkpmv2.MeasurementTypeId{
			Value: m.measTypeID,
		},
	}, nil
}
