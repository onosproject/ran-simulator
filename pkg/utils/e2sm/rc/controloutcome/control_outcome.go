// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package controloutcome

import (
	e2smrcpresm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/servicemodel"
	"google.golang.org/protobuf/proto"

	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/v2/e2sm-rc-pre-v2-go"
)

// ControlOutcome required fields for control outcome
type ControlOutcome struct {
	ranParameterID int32
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

// Build builds rc control outcome message
func (co *ControlOutcome) Build() (*e2smrcpreies.E2SmRcPreControlOutcome, error) {
	e2smRcPreOutcomeFormat1 := e2smrcpreies.E2SmRcPreControlOutcomeFormat1{
		OutcomeElementList: make([]*e2smrcpreies.RanparameterItem, 0),
	}
	outcomeElementList := &e2smrcpreies.RanparameterItem{
		RanParameterId: &e2smrcpreies.RanparameterId{
			Value: co.ranParameterID,
		},
	}
	e2smRcPreOutcomeFormat1.OutcomeElementList = append(e2smRcPreOutcomeFormat1.OutcomeElementList, outcomeElementList)
	e2smRcPrePdu := e2smrcpreies.E2SmRcPreControlOutcome{
		E2SmRcPreControlOutcome: &e2smrcpreies.E2SmRcPreControlOutcome_ControlOutcomeFormat1{
			ControlOutcomeFormat1: &e2smRcPreOutcomeFormat1,
		},
	}

	//ToDo - return it back once the Validation is functional again
	//if err := e2smRcPrePdu.Validate(); err != nil {
	//	return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	//}
	return &e2smRcPrePdu, nil

}

// ToAsn1Bytes converts to Asn1 bytes
func (co *ControlOutcome) ToAsn1Bytes() ([]byte, error) {
	outcomeRcMessage, err := co.Build()
	if err != nil {
		return nil, err
	}
	outcomeProtoBytes, err := proto.Marshal(outcomeRcMessage)
	if err != nil {
		return nil, err
	}

	var rcPreServiceModel e2smrcpresm.RcPreServiceModel
	outcomeAsn1Bytes, err := rcPreServiceModel.ControlOutcomeProtoToASN1(outcomeProtoBytes)
	if err != nil {
		return nil, err
	}

	return outcomeAsn1Bytes, nil
}
