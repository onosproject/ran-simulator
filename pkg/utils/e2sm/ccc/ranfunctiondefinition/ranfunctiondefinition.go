// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package ranfunctiondefinition

import (
	e2smccc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/v1/e2sm-ccc-ies"
	e2smcommon "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/v1/e2sm-common-ies"
)

// RANFunctionDefinition ran function definition fields
type RANFunctionDefinition struct {
	ranFunctionShortName                      string
	ranFunctionE2SmOID                        string
	ranFunctionDescription                    string
	ranFuncInstance                           int32
	listOfSupportedRANConfigurationStructures *e2smccc.ListOfSupportedRanconfigurationStructures
	listOfCellsForRANFunctionDefinition       *e2smccc.ListOfCellsForRanfunctionDefinition
}

// NewRANFunctionDefinition create new RAN function definition
func NewRANFunctionDefinition(options ...func(description *RANFunctionDefinition)) *RANFunctionDefinition {
	ranFunctionDefinition := &RANFunctionDefinition{}
	for _, option := range options {
		option(ranFunctionDefinition)
	}

	return ranFunctionDefinition
}

// WithRANFunctionShortName sets RAN function short name
func WithRANFunctionShortName(ranFunctionShortName string) func(description *RANFunctionDefinition) {
	return func(description *RANFunctionDefinition) {
		description.ranFunctionShortName = ranFunctionShortName
	}
}

// WithRANFunctionE2SmOID sets service model OID
func WithRANFunctionE2SmOID(ranFunctionE2SmOID string) func(description *RANFunctionDefinition) {
	return func(description *RANFunctionDefinition) {
		description.ranFunctionE2SmOID = ranFunctionE2SmOID
	}
}

// WithRANFunctionDefinition sets RAN function description
func WithRANFunctionDefinition(ranFunctionDescription string) func(description *RANFunctionDefinition) {
	return func(description *RANFunctionDefinition) {
		description.ranFunctionDescription = ranFunctionDescription
	}
}

// WithRANFunctionInstance sets RAN function instance
func WithRANFunctionInstance(ranFuncInstance int32) func(description *RANFunctionDefinition) {
	return func(description *RANFunctionDefinition) {
		description.ranFuncInstance = ranFuncInstance
	}
}

// WithListOfSupportedRanconfigurationStructures sets CCC node list
func WithListOfSupportedRanconfigurationStructures(listOfSupportedRANConfigurationStructures *e2smccc.ListOfSupportedRanconfigurationStructures) func(description *RANFunctionDefinition) {
	return func(description *RANFunctionDefinition) {
		description.listOfSupportedRANConfigurationStructures = listOfSupportedRANConfigurationStructures
	}
}

// WithListOfCellsForRanfunctionDefinition sets CCC cells list
func WithListOfCellsForRanfunctionDefinition(listOfCellsForRANFunctionDefinition *e2smccc.ListOfCellsForRanfunctionDefinition) func(description *RANFunctionDefinition) {
	return func(description *RANFunctionDefinition) {
		description.listOfCellsForRANFunctionDefinition = listOfCellsForRANFunctionDefinition
	}
}

// Build builds RAN function definition
func (r *RANFunctionDefinition) Build() (*e2smccc.E2SmCCcRAnfunctionDefinition, error) {
	e2SmCccPdu := e2smccc.E2SmCCcRAnfunctionDefinition{
		RanFunctionName: &e2smcommon.RanfunctionName{
			RanFunctionShortName:   r.ranFunctionShortName,
			RanFunctionE2SmOid:     r.ranFunctionE2SmOID,
			RanFunctionDescription: r.ranFunctionDescription,
			RanFunctionInstance:    &r.ranFuncInstance,
		},
		ListOfSupportedNodeLevelConfigurationStructures: r.listOfSupportedRANConfigurationStructures,
		ListOfCellsForRanfunctionDefinition:             r.listOfCellsForRANFunctionDefinition,
	}

	if err := e2SmCccPdu.Validate(); err != nil {
		return nil, err
	}
	return &e2SmCccPdu, nil
}
