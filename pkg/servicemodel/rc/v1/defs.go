// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package v1

const (
	modelFullName = "ORAN-E2SM-RC"
	version       = "v1"
	modelOID      = "1.3.6.1.4.1.53148.1.1.2.3"
)

const (
	eventTriggerStyle1  = "Message Event"
	eventTriggerStyle2  = "Call Process Breakpoint"
	eventTriggerStyle3  = "E2 Node Information"
	controlStyleType3   = 3
	controlStyleType200 = 200 // for PCI use-case: since there is no style for PCI use-case, define a new style
	controlActionID1    = 1

	ricInsertIndicationIDForMHO = 1
	ricInsertStyleType3         = 3

	ricPolicyStyleType3 = 3
	ricPolicyStyleName  = "Connected Mode Mobility Control"

	ricPolicyActionIDForMLB             = 1
	ricPolicyActionNameForMLB           = "Policy for Handover Control"
	ricActionDefinitionFormatTypeForMLB = 2
)

// RAN parameter IDs
const (
	// PCIRANParameterID PCI RAN parameter ID
	PCIRANParameterID = 1
	// NCGIRANParameterID NCGI RAN parameter ID
	NCGIRANParameterID = 2
	// NS xApp Id
	NSRANParameterID = 3

	// TargetPrimaryCellIDRANParameterID Target Primary Cell ID RAN parameter ID
	TargetPrimaryCellIDRANParameterID = 1
	// TargetPrimaryCellIDRANParameterName Target Primary Cell ID RAN parameter Name
	TargetPrimaryCellIDRANParameterName = "Target Primary Cell ID"
	// TargetCellRANParameterID Choice of Target Cell RAN parameter ID
	TargetCellRANParameterID = 2
	// NRCellRANParameterID NR Cell RAN parameter ID
	NRCellRANParameterID = 3
	// NRCGIRANParameterID NR CGI RAN parameter ID
	NRCGIRANParameterID = 4
	// EUTRACellRANParameterID E-UTRA Cell RAN parameter ID
	EUTRACellRANParameterID = 5
	// EUTRACGIRANParameterID E-UTRA CGI RAN parameter ID
	EUTRACGIRANParameterID = 6

	// CellSpecificOffsetRANParameterID Ocn RAN parameter ID
	CellSpecificOffsetRANParameterID = 10201
	// CellSpecificOffsetRANParameterName Ocn RAN parameter name
	CellSpecificOffsetRANParameterName = "Cell Specific Offset"
)

// UE Event IDs
const (
	A3MeasurementReportUEEventID = 2
)

// Call Process Breakpoint
const (
	CallProcessTypeIDMobilityManagement = 3
	CallBreakpointIDHandoverPreparation = 1
)
