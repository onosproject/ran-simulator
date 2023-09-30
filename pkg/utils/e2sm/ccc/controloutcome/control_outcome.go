// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package controloutcome

import (
	e2smcccsm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/servicemodel"
	"google.golang.org/protobuf/proto"

	e2smccc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/v1/e2sm-ccc-ies"
)

// Message control outcome message format 1 fields for ccc service model
type ControlOutcome struct {
	timestamp                              []byte
	ranConfigurationStructuresAcceptedList *e2smccc.RanConfigurationStructuresAcceptedList
	ranConfigurationStructuresFailedList   *e2smccc.RanConfigurationStructuresFailedList
}

// NewControlOutcome creates a new control outcome
func NewControlOutcome(options ...func(outcome *ControlOutcome)) *ControlOutcome {
	outcome := &ControlOutcome{}
	for _, option := range options {
		option(outcome)
	}
	return outcome
}

// WithReceivedTimestamp sets timestamp
func WithReceivedTimestamp(timestamp []byte) func(ctrl *ControlOutcome) {
	return func(ctrl *ControlOutcome) {
		ctrl.timestamp = timestamp
	}
}

// WithRanConfigurationStructuresAcceptedList sets measurement info list
func WithRanConfigurationStructuresAcceptedList(ranConfigurationStructuresAcceptedList *e2smccc.RanConfigurationStructuresAcceptedList) func(ctrl *ControlOutcome) {
	return func(ctrl *ControlOutcome) {
		ctrl.ranConfigurationStructuresAcceptedList = ranConfigurationStructuresAcceptedList
	}
}

// WithRanConfigurationStructuresFailedList sets measurement info list
func WithRanConfigurationStructuresFailedList(ranConfigurationStructuresFailedList *e2smccc.RanConfigurationStructuresFailedList) func(ctrl *ControlOutcome) {
	return func(ctrl *ControlOutcome) {
		ctrl.ranConfigurationStructuresFailedList = ranConfigurationStructuresFailedList
	}
}

// Build builds ccc control outcome message
func (ctrl *ControlOutcome) Build() (*e2smccc.E2SmCCcRIcControlOutcome, error) {
	e2smCccPdu := &e2smccc.E2SmCCcRIcControlOutcome{
		ControlOutcomeFormat: &e2smccc.ControlOutcomeFormat{
			ControlOutcomeFormat: &e2smccc.ControlOutcomeFormat_E2SmCccControlOutcomeFormat1{
				E2SmCccControlOutcomeFormat1: &e2smccc.E2SmCCcControlOutcomeFormat1{
					ReceivedTimestamp:                      ctrl.timestamp,
					RanConfigurationStructuresAcceptedList: ctrl.ranConfigurationStructuresAcceptedList,
					RanConfigurationStructuresFailedList:   ctrl.ranConfigurationStructuresFailedList,
				},
			},
		},
	}

	if err := e2smCccPdu.Validate(); err != nil {
		return nil, err
	}

	return e2smCccPdu, nil
}

// ToAsn1Bytes converts to Asn1 bytes
func (ctrl *ControlOutcome) ToAsn1Bytes() ([]byte, error) {
	// Creating a control outcome
	outcomeMessage, err := ctrl.Build()
	if err != nil {
		return nil, err
	}
	outcomeProtoBytes, err := proto.Marshal(outcomeMessage)
	if err != nil {
		return nil, err
	}

	var cccServiceModel e2smcccsm.CCCServiceModel
	outcomeAsn1Bytes, err := cccServiceModel.ControlOutcomeProtoToASN1(outcomeProtoBytes)
	if err != nil {
		return nil, err
	}

	return outcomeAsn1Bytes, nil
}
