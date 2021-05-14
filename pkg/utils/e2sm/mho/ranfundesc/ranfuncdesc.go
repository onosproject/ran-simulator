// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package ranfundesc

import (
	"fmt"

	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho/v1/e2sm-mho"
)

// RANFunctionDescription ran function description
type RANFunctionDescription struct {
	ranFunctionShortName     string
	ranFunctionE2SmOID       string
	ranFunctionDescription   string
	ranFunctionInstance      int32
	ricEventTriggerStyleList []*e2sm_mho.RicEventTriggerStyleList
	ricReportStyleList       []*e2sm_mho.RicReportStyleList
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
func WithRICEventTriggerStyleList(ricEventTriggerStyleList []*e2sm_mho.RicEventTriggerStyleList) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ricEventTriggerStyleList = ricEventTriggerStyleList

	}
}

// WithRICReportStyleList sets RIC report style list
func WithRICReportStyleList(ricReportStyleList []*e2sm_mho.RicReportStyleList) func(description *RANFunctionDescription) {
	return func(description *RANFunctionDescription) {
		description.ricReportStyleList = ricReportStyleList

	}
}

// Build builds RAN function description
func (desc *RANFunctionDescription) Build() (*e2sm_mho.E2SmMhoRanfunctionDescription, error) {
	e2smMhoPdu := e2sm_mho.E2SmMhoRanfunctionDescription{
		RanFunctionName: &e2sm_mho.RanfunctionName{
			RanFunctionShortName:   desc.ranFunctionShortName,
			RanFunctionE2SmOid:     desc.ranFunctionE2SmOID,
			RanFunctionDescription: desc.ranFunctionDescription,
			RanFunctionInstance:    desc.ranFunctionInstance,
		},
		RicEventTriggerStyleList: desc.ricEventTriggerStyleList,
		RicReportStyleList:       desc.ricReportStyleList,
	}

	if err := e2smMhoPdu.Validate(); err != nil {
		return nil, fmt.Errorf("error validating E2SmPDU %s", err.Error())
	}
	return &e2smMhoPdu, nil
}
