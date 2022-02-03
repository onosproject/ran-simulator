// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package reportstyle

import (
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
)

// Item Report style Item fields
type Item struct {
	ricStyleType            int32
	ricStyleName            string
	ricFormatType           int32
	measInfoActionList      *e2smkpmv2.MeasurementInfoActionList
	indicationHdrFormatType int32
	indicationMsgFormatType int32
}

// NewReportStyleItem create new report style item
func NewReportStyleItem(options ...func(item *Item)) *Item {
	reportStyleItem := &Item{}
	for _, option := range options {
		option(reportStyleItem)
	}

	return reportStyleItem
}

// WithRICStyleType sets RIC style type
func WithRICStyleType(ricStyleType int32) func(item *Item) {
	return func(item *Item) {
		item.ricStyleType = ricStyleType
	}
}

// WithRICStyleName sets RIC style name
func WithRICStyleName(ricStyleName string) func(item *Item) {
	return func(item *Item) {
		item.ricStyleName = ricStyleName
	}
}

// WithRICFormatType sets RIC format type
func WithRICFormatType(ricFormatType int32) func(item *Item) {
	return func(item *Item) {
		item.ricFormatType = ricFormatType
	}
}

// WithMeasInfoActionList sets meas info list
func WithMeasInfoActionList(measInfoActionList *e2smkpmv2.MeasurementInfoActionList) func(item *Item) {
	return func(item *Item) {
		item.measInfoActionList = measInfoActionList
	}
}

// WithIndicationHdrFormatType sets indication header format type
func WithIndicationHdrFormatType(indicationHdrFormatType int32) func(item *Item) {
	return func(item *Item) {
		item.indicationHdrFormatType = indicationHdrFormatType
	}
}

// WithIndicationMsgFormatType sets indication message format type
func WithIndicationMsgFormatType(indicationMsgFormatType int32) func(item *Item) {
	return func(item *Item) {
		item.indicationMsgFormatType = indicationMsgFormatType
	}
}

// Build builds RIC report style item
func (i *Item) Build() *e2smkpmv2.RicReportStyleItem {
	return &e2smkpmv2.RicReportStyleItem{
		RicReportStyleType: &e2smkpmv2.RicStyleType{
			Value: i.ricStyleType,
		},
		RicReportStyleName: &e2smkpmv2.RicStyleName{
			Value: i.ricStyleName,
		},
		RicActionFormatType: &e2smkpmv2.RicFormatType{
			Value: i.ricFormatType,
		},
		MeasInfoActionList: i.measInfoActionList,
		RicIndicationHeaderFormatType: &e2smkpmv2.RicFormatType{
			Value: i.indicationHdrFormatType,
		},
		RicIndicationMessageFormatType: &e2smkpmv2.RicFormatType{
			Value: i.indicationMsgFormatType,
		},
	}
}
