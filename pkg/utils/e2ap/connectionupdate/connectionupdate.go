// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connectionupdate

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
)

// ConnectionUpdate info for building connection update ack and failure responses
type ConnectionUpdate struct {
	connectionUpdateItemIes      []*e2appducontents.E2ConnectionUpdateItemIes
	connectionSetupFailedItemIes []*e2appducontents.E2ConnectionSetupFailedItemIes
	cause                        *e2apies.Cause
	timeToWait                   *e2apies.TimeToWait
	transactionID                int32
}

// NewConnectionUpdate creates a new instance of connection update
func NewConnectionUpdate(options ...func(update *ConnectionUpdate)) *ConnectionUpdate {
	connectionUpdate := &ConnectionUpdate{}

	for _, option := range options {
		option(connectionUpdate)
	}
	return connectionUpdate
}

// WithCause sets cause of failure
func WithCause(cause *e2apies.Cause) func(*ConnectionUpdate) {
	return func(connectionUpdate *ConnectionUpdate) {
		connectionUpdate.cause = cause
	}
}

// WithTimeToWait sets time to wait
func WithTimeToWait(timeToWait *e2apies.TimeToWait) func(*ConnectionUpdate) {
	return func(connectionUpdate *ConnectionUpdate) {
		connectionUpdate.timeToWait = timeToWait
	}
}

// WithConnectionUpdateItemIes sets connection update item ies
func WithConnectionUpdateItemIes(connectionUpdateItemIes []*e2appducontents.E2ConnectionUpdateItemIes) func(*ConnectionUpdate) {
	return func(connectionUpdate *ConnectionUpdate) {
		connectionUpdate.connectionUpdateItemIes = connectionUpdateItemIes
	}
}

// WithConnectionSetupFailedItemIes sets connection setup failed item Ies
func WithConnectionSetupFailedItemIes(connectionSetupFailedItemIes []*e2appducontents.E2ConnectionSetupFailedItemIes) func(*ConnectionUpdate) {
	return func(connectionUpdate *ConnectionUpdate) {
		connectionUpdate.connectionSetupFailedItemIes = connectionSetupFailedItemIes
	}
}

// WithTransactionID sets transaction ID
func WithTransactionID(transID int32) func(*ConnectionUpdate) {
	return func(connectionUpdate *ConnectionUpdate) {
		connectionUpdate.transactionID = transID
	}
}

// BuildConnectionUpdateAcknowledge creates a connection update acknowledge
func (c *ConnectionUpdate) BuildConnectionUpdateAcknowledge() *e2appducontents.E2ConnectionUpdateAcknowledge {

	response := &e2appducontents.E2ConnectionUpdateAcknowledge{
		ProtocolIes: make([]*e2appducontents.E2ConnectionUpdateAckIes, 0),
	}
	response.SetTransactionID(c.transactionID)

	if len(c.connectionUpdateItemIes) != 0 {
		e2cul := &e2appducontents.E2ConnectionUpdateList{
			Value: make([]*e2appducontents.E2ConnectionUpdateItemIes, 0),
		}
		e2cul.Value = c.connectionUpdateItemIes

		ie := &e2appducontents.E2ConnectionUpdateAckIes{
			Id:          int32(v2.ProtocolIeIDE2connectionSetup),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value: &e2appducontents.E2ConnectionUpdateAckIe{
				E2ConnectionUpdateAckIe: &e2appducontents.E2ConnectionUpdateAckIe_E2ConnectionSetup{
					E2ConnectionSetup: e2cul,
				},
			},
		}
		response.ProtocolIes = append(response.ProtocolIes, ie)
	}

	if len(c.connectionSetupFailedItemIes) != 0 {
		e2csfl := &e2appducontents.E2ConnectionSetupFailedList{
			Value: make([]*e2appducontents.E2ConnectionSetupFailedItemIes, 0),
		}
		e2csfl.Value = c.connectionSetupFailedItemIes

		ie := &e2appducontents.E2ConnectionUpdateAckIes{
			Id:          int32(v2.ProtocolIeIDE2connectionSetupFailed),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value: &e2appducontents.E2ConnectionUpdateAckIe{
				E2ConnectionUpdateAckIe: &e2appducontents.E2ConnectionUpdateAckIe_E2ConnectionSetupFailed{
					E2ConnectionSetupFailed: e2csfl,
				},
			},
		}
		response.ProtocolIes = append(response.ProtocolIes, ie)
	}

	return response
}

// BuildConnectionUpdateFailure creates a connection update failure message
func (c *ConnectionUpdate) BuildConnectionUpdateFailure() *e2appducontents.E2ConnectionUpdateFailure {

	failure := &e2appducontents.E2ConnectionUpdateFailure{
		ProtocolIes: make([]*e2appducontents.E2ConnectionUpdateFailureIes, 0),
	}
	failure.SetTransactionID(c.transactionID).SetCause(c.cause)

	return failure
}
