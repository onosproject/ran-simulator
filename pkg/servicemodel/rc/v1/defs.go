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
	eventTriggerStyle2  = "Call Process Breakpoint"
	eventTriggerStyle3  = "E2 Node Information"
	controlStyleType200 = 200 // for PCI use-case: since there is no style for PCI use-case, define a new style
	controlActionID1    = 1
)

// RAN parameter IDs
const (
	// PCIRANParameterID PCI RAN parameter ID
	PCIRANParameterID = 1
	// NCGIRANParameterID NCGI RAN parameter ID
	NCGIRANParameterID = 2
)
