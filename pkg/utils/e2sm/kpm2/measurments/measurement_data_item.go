// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measurments

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
)

// MeasurementDataItem measurement data item
type MeasurementDataItem struct {
	mr             *e2smkpmv2.MeasurementRecord
	incompleteFlag e2smkpmv2.IncompleteFlag
}

// NewMeasurementDataItem creates a new measurement data item
func NewMeasurementDataItem(options ...func(item *MeasurementDataItem)) *MeasurementDataItem {
	measDataItem := &MeasurementDataItem{}
	for _, option := range options {
		option(measDataItem)
	}

	return measDataItem
}

// WithMeasurementRecord sets measurement record
func WithMeasurementRecord(mr *e2smkpmv2.MeasurementRecord) func(item *MeasurementDataItem) {
	return func(item *MeasurementDataItem) {
		item.mr = mr
	}
}

// WithIncompleteFlag sets incomplete flag
func WithIncompleteFlag(incompleteFlag e2smkpmv2.IncompleteFlag) func(item *MeasurementDataItem) {
	return func(item *MeasurementDataItem) {
		item.incompleteFlag = incompleteFlag
	}
}

// Build builds a measurement data item
func (m *MeasurementDataItem) Build() (*e2smkpmv2.MeasurementDataItem, error) {
	mdi := e2smkpmv2.MeasurementDataItem{
		MeasRecord:     m.mr,
		IncompleteFlag: &m.incompleteFlag,
	}

	// FIXME: Add back when ready
	//if err := mdi.Validate(); err != nil {
	//	return nil, errors.New(errors.Invalid, err.Error())
	//}
	return &mdi, nil
}
