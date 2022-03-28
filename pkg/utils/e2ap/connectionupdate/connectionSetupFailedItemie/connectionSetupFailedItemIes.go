// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connectionsetupfaileditem

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
)

// IEs connection setup failed item Ies
type IEs struct {
	tnlInfo *e2apies.Tnlinformation
	cause   *e2apies.Cause
}

// NewConnectionSetupFailedItemIe creates a new instance of connection setup failed Item Ie
func NewConnectionSetupFailedItemIe(options ...func(update *IEs)) *IEs {
	connectionSetupFailedItemIe := &IEs{}

	for _, option := range options {
		option(connectionSetupFailedItemIe)
	}

	return connectionSetupFailedItemIe
}

// WithTnlInfo sets tnl info
func WithTnlInfo(tnlInfo *e2apies.Tnlinformation) func(ie *IEs) {
	return func(connectionSetupFailedItemIe *IEs) {
		connectionSetupFailedItemIe.tnlInfo = tnlInfo
	}
}

// BuildConnectionSetupFailedItemIes builds connection setup failed Item Ies
func (c *IEs) BuildConnectionSetupFailedItemIes() *e2appducontents.E2ConnectionSetupFailedItemIes {
	connectionSetupFailedItemIes := &e2appducontents.E2ConnectionSetupFailedItemIes{
		Id:          int32(v2.ProtocolIeIDE2connectionSetupFailedItem),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
		Value: &e2appducontents.E2ConnectionSetupFailedItemIe{
			E2ConnectionSetupFailedItemIe: &e2appducontents.E2ConnectionSetupFailedItemIe_E2ConnectionSetupFailedItem{
				E2ConnectionSetupFailedItem: &e2appducontents.E2ConnectionSetupFailedItem{
					TnlInformation: c.tnlInfo,
					Cause:          c.cause,
				},
			},
		},
	}

	return connectionSetupFailedItemIes
}
