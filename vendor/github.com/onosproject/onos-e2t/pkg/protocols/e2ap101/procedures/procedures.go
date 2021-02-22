// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package procedures

import (
	e2appdudescriptions "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-descriptions"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"io"
)

var log = logging.GetLogger("protocols", "e2")

// ElementaryProcedure is anb identifier interface for E2 elementary procedure interfaces
type ElementaryProcedure interface {
	io.Closer
	// Matches returns a bool indicating whether the given PDU matches the procedure
	Matches(pdu *e2appdudescriptions.E2ApPdu) bool
	// Handle handles a matching PDU
	Handle(pdu *e2appdudescriptions.E2ApPdu)
}

// Dispatcher is a function for sending a message
type Dispatcher func(pdu *e2appdudescriptions.E2ApPdu) error

// RICProcedures implements the procedures for a RIC
type RICProcedures interface {
	E2Setup
	RICIndication
}

// E2NodeProcedures implements the procedures for an E2 node
type E2NodeProcedures interface {
	RICControl
	RICSubscription
	RICSubscriptionDelete
}
