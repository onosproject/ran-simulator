// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ranfuncdescription

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
)

// RANFunctionDescription ran function description fields
type RANFunctionDescription struct {
	ranFunctionShortName     string
	ranFunctionE2SmOID       string
	ranFunctionDescription   string
	ranFuncInstance          int32
	ricKpmNodeList           []*e2smkpmv2.RicKpmnodeItem
	ricEventTriggerStyleList []*e2smkpmv2.RicEventTriggerStyleItem
	ricReportStyleList       []*e2smkpmv2.RicReportStyleItem
}

// NewRANFunctionDescription create new RAN function description
func NewRANFunctionDescription(options ...func(description *RANFunctionDescription)) *RANFunctionDescription {
	ranFunctionDescription := &RANFunctionDescription{}
	for _, option := range options {
		option(ranFunctionDescription)
	}

	return ranFunctionDescription
}

// WithRANFunctionShortName sets RAN function short name
func WithRANFunctionShortName(ranFunctionShortName string) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ranFunctionShortName = ranFunctionShortName
	}
}

// WithRANFunctionE2SmOID sets service model OID
func WithRANFunctionE2SmOID(ranFunctionE2SmOID string) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ranFunctionE2SmOID = ranFunctionE2SmOID
	}
}

// WithRANFunctionDescription sets RAN function description
func WithRANFunctionDescription(ranFunctionDescription string) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ranFunctionDescription = ranFunctionDescription
	}
}

// WithRANFunctionInstance sets RAN function instance
func WithRANFunctionInstance(ranFuncInstance int32) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ranFuncInstance = ranFuncInstance
	}
}

// WithRICKPMNodeList sets KPM node list
func WithRICKPMNodeList(ricKpmNodeList []*e2smkpmv2.RicKpmnodeItem) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ricKpmNodeList = ricKpmNodeList
	}
}

// WithRICEventTriggerStyleList sets event trigger style list
func WithRICEventTriggerStyleList(ricEventTriggerStyleList []*e2smkpmv2.RicEventTriggerStyleItem) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ricEventTriggerStyleList = ricEventTriggerStyleList
	}
}

// WithRICReportStyleList sets RIC report style list
func WithRICReportStyleList(ricReportStyleList []*e2smkpmv2.RicReportStyleItem) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ricReportStyleList = ricReportStyleList
	}
}

// Build builds RAN function description
func (r *RANFunctionDescription) Build() (*e2smkpmv2.E2SmKpmRanfunctionDescription, error) {
	e2SmKpmPdu := e2smkpmv2.E2SmKpmRanfunctionDescription{
		RanFunctionName: &e2smkpmv2.RanfunctionName{
			RanFunctionShortName:   r.ranFunctionShortName,
			RanFunctionE2SmOid:     r.ranFunctionE2SmOID,
			RanFunctionDescription: r.ranFunctionDescription,
			RanFunctionInstance:    &r.ranFuncInstance,
		},
		RicKpmNodeList:           r.ricKpmNodeList,
		RicEventTriggerStyleList: r.ricEventTriggerStyleList,
		RicReportStyleList:       r.ricReportStyleList,
	}

	// FIXME: Add back when ready
	//if err := e2SmKpmPdu.Validate(); err != nil {
	//	return nil, errors.New(errors.Invalid, err.Error())
	//}
	return &e2SmKpmPdu, nil
}
