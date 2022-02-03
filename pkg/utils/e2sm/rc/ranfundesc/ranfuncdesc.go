// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ranfundesc

import (
	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/v2/e2sm-rc-pre-v2-go"
)

// RANFunctionDescription ran function description
type RANFunctionDescription struct {
	ranFunctionShortName     string
	ranFunctionE2SmOID       string
	ranFunctionDescription   string
	ranFunctionInstance      int32
	ricEventTriggerStyleList []*e2smrcpreies.RicEventTriggerStyleList
	ricReportStyleList       []*e2smrcpreies.RicReportStyleList
}

// NewRANFunctionDescription creates a  new RAN function description
func NewRANFunctionDescription(options ...func(description *RANFunctionDescription)) *RANFunctionDescription {
	desc := &RANFunctionDescription{}
	for _, option := range options {
		option(desc)
	}
	return desc
}

// WithRANFunctionShortName sets ran parameter ID
func WithRANFunctionShortName(shortName string) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ranFunctionShortName = shortName

	}
}

// WithRANFunctionE2SmOID sets OID
func WithRANFunctionE2SmOID(oid string) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ranFunctionE2SmOID = oid

	}
}

// WithRANFunctionDescription sets ran function description
func WithRANFunctionDescription(ranFuncDesc string) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ranFunctionDescription = ranFuncDesc

	}
}

// WithRANFunctionInstance sets ran function instance
func WithRANFunctionInstance(instance int32) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ranFunctionInstance = instance
	}
}

// WithRICEventTriggerStyleList sets RIC event trigger style list
func WithRICEventTriggerStyleList(ricEventTriggerStyleList []*e2smrcpreies.RicEventTriggerStyleList) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ricEventTriggerStyleList = ricEventTriggerStyleList

	}
}

// WithRICReportStyleList sets RIC report style list
func WithRICReportStyleList(ricReportStyleList []*e2smrcpreies.RicReportStyleList) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ricReportStyleList = ricReportStyleList

	}
}

// Build builds RAN function description
func (desc *RANFunctionDescription) Build() (*e2smrcpreies.E2SmRcPreRanfunctionDescription, error) {
	ranfunctionItem := e2smrcpreies.E2SmRcPreRanfunctionDescription_E2SmRcPreRanfunctionItem001{
		RicEventTriggerStyleList: desc.ricEventTriggerStyleList,
		RicReportStyleList:       desc.ricReportStyleList,
	}

	e2smRcPrePdu := e2smrcpreies.E2SmRcPreRanfunctionDescription{
		RanFunctionName: &e2smrcpreies.RanfunctionName{
			RanFunctionShortName:   desc.ranFunctionShortName,
			RanFunctionE2SmOid:     desc.ranFunctionE2SmOID,
			RanFunctionDescription: desc.ranFunctionDescription,
			RanFunctionInstance:    &desc.ranFunctionInstance,
		},
		E2SmRcPreRanfunctionItem: &ranfunctionItem,
	}

	//ToDo - return it back once the Validation is functional again
	//if err := e2smRcPrePdu.Validate(); err != nil {
	//	return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	//}
	return &e2smRcPrePdu, nil
}
