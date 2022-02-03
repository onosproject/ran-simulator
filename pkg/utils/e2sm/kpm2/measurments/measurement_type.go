// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measurments

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
)

// MeasurementTypeMeasName measurement type meas name
type MeasurementTypeMeasName struct {
	Name string
}

// NewMeasurementTypeMeasName creates a new measurement type meas name
func NewMeasurementTypeMeasName(options ...func(*MeasurementTypeMeasName)) *MeasurementTypeMeasName {
	measTypeName := &MeasurementTypeMeasName{}
	for _, option := range options {
		option(measTypeName)
	}

	return measTypeName
}

// WithMeasurementName sets measurement name
func WithMeasurementName(name string) func(*MeasurementTypeMeasName) {
	return func(measurementName *MeasurementTypeMeasName) {
		measurementName.Name = name
	}

}

// Build builds measurement type meas name
func (m *MeasurementTypeMeasName) Build() (*e2smkpmv2.MeasurementType, error) {
	measType := e2smkpmv2.MeasurementType{
		MeasurementType: &e2smkpmv2.MeasurementType_MeasName{
			MeasName: &e2smkpmv2.MeasurementTypeName{
				Value: m.Name,
			},
		},
	}

	// FIXME: Add back when ready
	//if err := measType.Validate(); err != nil {
	//	return nil, errors.New(errors.Invalid, err.Error())
	//}

	return &measType, nil
}

// MeasurementTypeMeasID measurement type meas ID
type MeasurementTypeMeasID struct {
	ID int32
}

// NewMeasurementTypeMeasID creates a new measurement type meas ID
func NewMeasurementTypeMeasID(options ...func(id *MeasurementTypeMeasID)) *MeasurementTypeMeasID {
	measTypeID := &MeasurementTypeMeasID{}
	for _, option := range options {
		option(measTypeID)
	}

	return measTypeID
}

// Build builds a measurement type meas ID
func (m *MeasurementTypeMeasID) Build() (*e2smkpmv2.MeasurementType, error) {
	measType := e2smkpmv2.MeasurementType{
		MeasurementType: &e2smkpmv2.MeasurementType_MeasId{
			MeasId: &e2smkpmv2.MeasurementTypeId{
				Value: m.ID,
			},
		},
	}

	// FIXME: Add back when ready
	//if err := measType.Validate(); err != nil {
	//	return nil, errors.New(errors.Invalid, err.Error())
	//}

	return &measType, nil
}
