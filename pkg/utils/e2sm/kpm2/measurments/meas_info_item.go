// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measurments

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
)

// MeasurementInfoItem measurement info item
type MeasurementInfoItem struct {
	measType      *e2smkpmv2.MeasurementType
	labelInfoList *e2smkpmv2.LabelInfoList
}

// NewMeasurementInfoItem creates a new measurement info item
func NewMeasurementInfoItem(options ...func(item *MeasurementInfoItem)) *MeasurementInfoItem {
	measInfoItem := &MeasurementInfoItem{}
	for _, option := range options {
		option(measInfoItem)
	}

	return measInfoItem
}

// WithMeasType sets measurement type
func WithMeasType(measType *e2smkpmv2.MeasurementType) func(item *MeasurementInfoItem) {
	return func(item *MeasurementInfoItem) {
		item.measType = measType
	}
}

// WithLabelInfoList sets label info list
func WithLabelInfoList(labelInfoList *e2smkpmv2.LabelInfoList) func(item *MeasurementInfoItem) {
	return func(item *MeasurementInfoItem) {
		item.labelInfoList = labelInfoList
	}
}

// Build builds measurement info item
func (m *MeasurementInfoItem) Build() (*e2smkpmv2.MeasurementInfoItem, error) {
	item := e2smkpmv2.MeasurementInfoItem{
		MeasType:      m.measType,
		LabelInfoList: m.labelInfoList,
	}

	// FIXME: Add back when ready
	//if err := item.Validate(); err != nil {
	//	return nil, errors.New(errors.Invalid, err.Error())
	//}

	return &item, nil
}
