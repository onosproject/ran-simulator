// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package v1beta1

// Driven from e2ap_constants.proto
// TODO: Automate the generation of this file

type ProcedureCodeT int32

const (
	ProcedureCodeIDdummy                 ProcedureCodeT = 0
	ProcedureCodeIDE2setup               ProcedureCodeT = 1
	ProcedureCodeIDErrorIndication       ProcedureCodeT = 2
	ProcedureCodeIDReset                 ProcedureCodeT = 3
	ProcedureCodeIDRICcontrol            ProcedureCodeT = 4
	ProcedureCodeIDRICindication         ProcedureCodeT = 5
	ProcedureCodeIDRICserviceQuery       ProcedureCodeT = 6
	ProcedureCodeIDRICserviceUpdate      ProcedureCodeT = 7
	ProcedureCodeIDRICsubscription       ProcedureCodeT = 8
	ProcedureCodeIDRICsubscriptionDelete ProcedureCodeT = 9
)

type ProtocolIeID int32

const (
	ProtocolIeIDCause                    ProtocolIeID = 1
	ProtocolIeIDCriticalityDiagnostics   ProtocolIeID = 2
	ProtocolIeIDGlobalE2nodeID           ProtocolIeID = 3
	ProtocolIeIDGlobalRicID              ProtocolIeID = 4
	ProtocolIeIDRanfunctionID            ProtocolIeID = 5
	ProtocolIeIDRanfunctionIDItem        ProtocolIeID = 6
	ProtocolIeIDRanfunctionIeCauseItem   ProtocolIeID = 7
	ProtocolIeIDRanfunctionItem          ProtocolIeID = 8
	ProtocolIeIDRanfunctionsAccepted     ProtocolIeID = 9
	ProtocolIeIDRanfunctionsAdded        ProtocolIeID = 10
	ProtocolIeIDRanfunctionsDeleted      ProtocolIeID = 11
	ProtocolIeIDRanfunctionsModified     ProtocolIeID = 12
	ProtocolIeIDRanfunctionsRejected     ProtocolIeID = 13
	ProtocolIeIDRicactionAdmittedItem    ProtocolIeID = 14
	ProtocolIeIDRicactionID              ProtocolIeID = 15
	ProtocolIeIDRicactionNotAdmittedItem ProtocolIeID = 16
	ProtocolIeIDRicactionsAdmitted       ProtocolIeID = 17
	ProtocolIeIDRicactionsNotAdmitted    ProtocolIeID = 18
	ProtocolIeIDRicactionToBeSetupItem   ProtocolIeID = 19
	ProtocolIeIDRiccallProcessID         ProtocolIeID = 20
	ProtocolIeIDRiccontrolAckRequest     ProtocolIeID = 21
	ProtocolIeIDRiccontrolHeader         ProtocolIeID = 22
	ProtocolIeIDRiccontrolMessage        ProtocolIeID = 23
	ProtocolIeIDRiccontrolStatus         ProtocolIeID = 24
	ProtocolIeIDRicindicationHeader      ProtocolIeID = 25
	ProtocolIeIDRicindicationMessage     ProtocolIeID = 26
	ProtocolIeIDRicindicationSn          ProtocolIeID = 27
	ProtocolIeIDRicindicationType        ProtocolIeID = 28
	ProtocolIeIDRicrequestID             ProtocolIeID = 29
	ProtocolIeIDRicsubscriptionDetails   ProtocolIeID = 30
	ProtocolIeIDTimeToWait               ProtocolIeID = 31
	ProtocolIeIDRiccontrolOutcome        ProtocolIeID = 32
)
