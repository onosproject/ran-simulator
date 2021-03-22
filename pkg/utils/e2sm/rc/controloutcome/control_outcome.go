// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package controloutcome

import (
	"fmt"

	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"

	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

// ControlOutcome required fields for control outcome
type ControlOutcome struct {
	ranParameterID    int32
	ranParameterValue int32
}

// NewControlOutcome creates a new control outcome
func NewControlOutcome(options ...func(outcome *ControlOutcome)) *ControlOutcome {
	outcome := &ControlOutcome{}
	for _, option := range options {
		option(outcome)
	}
	return outcome
}

// WithRanParameterID sets ran parameter ID
func WithRanParameterID(ranParameterID int32) func(co *ControlOutcome) {
	return func(co *ControlOutcome) {
		co.ranParameterID = ranParameterID

	}
}

// WithRanParameterValue sets ran parameter value
func WithRanParameterValue(ranParameterValue int32) func(co *ControlOutcome) {
	return func(co *ControlOutcome) {
		co.ranParameterValue = ranParameterValue
	}
}

// Build builds rc control outcome message
func (co *ControlOutcome) Build() (*e2sm_rc_pre_ies.E2SmRcPreControlOutcome, error) {
	e2smRcPreOutcomeFormat1 := e2sm_rc_pre_ies.E2SmRcPreControlOutcomeFormat1{
		OutcomeElementList: make([]*e2sm_rc_pre_ies.RanparameterItem, 0),
	}
	outcomeElementList := &e2sm_rc_pre_ies.RanparameterItem{
		RanParameterId: &e2sm_rc_pre_ies.RanparameterId{
			Value: co.ranParameterID,
		},
		RanParameterValue: &e2sm_rc_pre_ies.RanparameterValue{
			RanparameterValue: &e2sm_rc_pre_ies.RanparameterValue_ValueInt{
				ValueInt: co.ranParameterValue,
			},
		},
	}
	e2smRcPreOutcomeFormat1.OutcomeElementList = append(e2smRcPreOutcomeFormat1.OutcomeElementList, outcomeElementList)
	e2smRcPrePdu := e2sm_rc_pre_ies.E2SmRcPreControlOutcome{
		E2SmRcPreControlOutcome: &e2sm_rc_pre_ies.E2SmRcPreControlOutcome_ControlOutcomeFormat1{
			ControlOutcomeFormat1: &e2smRcPreOutcomeFormat1,
		},
	}

	if err := e2smRcPrePdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &e2smRcPrePdu, nil

}

// ToAsn1Bytes converts to Asn1 bytes
func (co *ControlOutcome) ToAsn1Bytes(modelPlugin modelplugins.ServiceModel) ([]byte, error) {
	outcomeRcMessage, err := co.Build()
	if err != nil {
		return nil, err
	}
	outcomeProtoBytes, err := proto.Marshal(outcomeRcMessage)
	if err != nil {
		return nil, err
	}

	outcomeAsn1Bytes, err := modelPlugin.ControlOutcomeProtoToASN1(outcomeProtoBytes)
	if err != nil {
		return nil, err
	}

	return outcomeAsn1Bytes, nil
}
