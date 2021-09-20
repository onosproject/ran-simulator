// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package connectionupdate

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v2beta1"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2beta1/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2beta1/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2beta1/e2ap-pdu-contents"
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

	ie39 := &e2appducontents.E2ConnectionUpdateAckIes_E2ConnectionUpdateAckIes39{}
	ie40 := &e2appducontents.E2ConnectionUpdateAckIes_E2ConnectionUpdateAckIes40{}
	ie49 := &e2appducontents.E2ConnectionUpdateAckIes_E2ConnectionUpdateAckIes49{
		Id:          int32(v2beta1.ProtocolIeIDTransactionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.TransactionId{
			Value: c.transactionID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	if c.connectionUpdateItemIes != nil {
		ie39 = &e2appducontents.E2ConnectionUpdateAckIes_E2ConnectionUpdateAckIes39{
			Id:          int32(v2beta1.ProtocolIeIDRicrequestID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value: &e2appducontents.E2ConnectionUpdateList{
				Value: make([]*e2appducontents.E2ConnectionUpdateItemIes, 0),
			},
			Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	}

	if c.connectionSetupFailedItemIes != nil {
		ie40 = &e2appducontents.E2ConnectionUpdateAckIes_E2ConnectionUpdateAckIes40{
			Id:          int32(v2beta1.ProtocolIeIDE2connectionSetupFailed),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value: &e2appducontents.E2ConnectionSetupFailedList{
				Value: make([]*e2appducontents.E2ConnectionSetupFailedItemIes, 0),
			},
			Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}
	}

	response := &e2appducontents.E2ConnectionUpdateAcknowledge{
		ProtocolIes: &e2appducontents.E2ConnectionUpdateAckIes{},
	}

	if len(c.connectionUpdateItemIes) != 0 {
		response.GetProtocolIes().E2ApProtocolIes39 = ie39
		response.GetProtocolIes().GetE2ApProtocolIes39().Value.Value = c.connectionUpdateItemIes
	}
	if len(c.connectionSetupFailedItemIes) != 0 {
		response.GetProtocolIes().E2ApProtocolIes40 = ie40
		response.GetProtocolIes().GetE2ApProtocolIes40().Value.Value = c.connectionSetupFailedItemIes
	}
	response.GetProtocolIes().E2ApProtocolIes49 = ie49
	return response
}

// BuildConnectionUpdateFailure creates a connection update failure message
func (c *ConnectionUpdate) BuildConnectionUpdateFailure() *e2appducontents.E2ConnectionUpdateFailure {
	ie49 := &e2appducontents.E2ConnectionUpdateFailureIes_E2ConnectionUpdateFailureIes49{
		Id:          int32(v2beta1.ProtocolIeIDTransactionID),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2apies.TransactionId{
			Value: c.transactionID,
		},
		Presence: int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
	}

	failure := &e2appducontents.E2ConnectionUpdateFailure{
		ProtocolIes: &e2appducontents.E2ConnectionUpdateFailureIes{
			E2ApProtocolIes1: &e2appducontents.E2ConnectionUpdateFailureIes_E2ConnectionUpdateFailureIes1{
				Id:          int32(v2beta1.ProtocolIeIDCause),
				Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
				Value:       c.cause,
				Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
			},
			//E2ApProtocolIes31: ie31, // TODO add time to wait whenever is needed
			//E2ApProtocolIes2: &criticalityDiagnostics, // TODO
			E2ApProtocolIes49: ie49,
		},
	}
	return failure
}
