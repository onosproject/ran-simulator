// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measurments

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
)

// MeasurementRecordItemInteger measurement record item integer
type MeasurementRecordItemInteger struct {
	value int64
}

// NewMeasurementRecordItemInteger creates a new measurement record item integer
func NewMeasurementRecordItemInteger(options ...func(integer *MeasurementRecordItemInteger)) *MeasurementRecordItemInteger {
	measRecordItemInteger := &MeasurementRecordItemInteger{}
	for _, option := range options {
		option(measRecordItemInteger)
	}

	return measRecordItemInteger
}

// WithIntegerValue sets record item integer value
func WithIntegerValue(value int64) func(integer *MeasurementRecordItemInteger) {
	return func(recordItem *MeasurementRecordItemInteger) {
		recordItem.value = value
	}
}

// Build builds a measurement record item integer
func (m *MeasurementRecordItemInteger) Build() *e2smkpmv2.MeasurementRecordItem {
	return &e2smkpmv2.MeasurementRecordItem{
		MeasurementRecordItem: &e2smkpmv2.MeasurementRecordItem_Integer{
			Integer: m.value,
		},
	}
}

// MeasurementRecordItemReal measurement record item real
type MeasurementRecordItemReal struct {
	value float64
}

// NewMeasurementRecordItemReal creates a new measurement record item real
func NewMeasurementRecordItemReal(options ...func(integer *MeasurementRecordItemReal)) *MeasurementRecordItemReal {
	measRecordItemReal := &MeasurementRecordItemReal{}
	for _, option := range options {
		option(measRecordItemReal)
	}

	return measRecordItemReal
}

// WithRealValue sets record item integer value
func WithRealValue(value float64) func(integer *MeasurementRecordItemReal) {
	return func(recordItem *MeasurementRecordItemReal) {
		recordItem.value = value
	}
}

// Build builds measurement record item real
func (m *MeasurementRecordItemReal) Build() *e2smkpmv2.MeasurementRecordItem {
	return &e2smkpmv2.MeasurementRecordItem{
		MeasurementRecordItem: &e2smkpmv2.MeasurementRecordItem_Real{
			Real: m.value,
		},
	}
}

// NewMeasurementRecordItemNoValue create new measurement  record item no value
func NewMeasurementRecordItemNoValue() *e2smkpmv2.MeasurementRecordItem {
	return &e2smkpmv2.MeasurementRecordItem{
		MeasurementRecordItem: &e2smkpmv2.MeasurementRecordItem_NoValue{
			NoValue: 0,
		},
	}
}
