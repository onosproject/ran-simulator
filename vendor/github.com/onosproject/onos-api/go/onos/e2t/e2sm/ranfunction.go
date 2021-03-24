// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

type ShortName string

type OID string

type Description string

type Instance int32

type StyleType int32

type StyleName string

type FormatType int32

type Version string

type ModuleName string

type RanfunctionNameDef struct {
	RanFunctionShortName   ShortName
	RanFunctionE2SmOid     OID
	RanFunctionDescription Description
	RanFunctionInstance    Instance
}

type RicReportStyleDef struct {
	RicReportStyleType             StyleType
	RicReportStyleName             StyleName
	RicIndicationHeaderFormatType  FormatType
	RicIndicationMessageFormatType FormatType
}

type RicEventTriggerDef struct {
	RicEventStyleType  StyleType
	RicEventStyleName  StyleName
	RicEventFormatType FormatType
}

type RicReportList map[StyleType]RicReportStyleDef
type RicEventTriggerList map[StyleType]RicEventTriggerDef
