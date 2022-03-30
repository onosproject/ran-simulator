// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connectionupdateitem

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-commondatatypes"
	e2ap_ies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
)

// IEs connection item update
type IEs struct {
	tnlInfo  *e2ap_ies.Tnlinformation
	tnlUsage e2ap_ies.Tnlusage
}

// NewConnectionUpdateItemIe creates a new instance of connection update Item Ie
func NewConnectionUpdateItemIe(options ...func(update *IEs)) *IEs {
	connectionUpdateItemIe := &IEs{}

	for _, option := range options {
		option(connectionUpdateItemIe)
	}

	return connectionUpdateItemIe
}

// WithTnlInfo sets tnl info
func WithTnlInfo(tnlInfo *e2ap_ies.Tnlinformation) func(*IEs) {
	return func(connectionUpdateItemIe *IEs) {
		connectionUpdateItemIe.tnlInfo = tnlInfo
	}
}

// WithTnlUsage sets tnl usage
func WithTnlUsage(tnlUsage e2ap_ies.Tnlusage) func(*IEs) {
	return func(connectionUpdateItemIe *IEs) {
		connectionUpdateItemIe.tnlUsage = tnlUsage
	}
}

// BuildConnectionUpdateItemIes builds connection Update Item Ies
func (c *IEs) BuildConnectionUpdateItemIes() *e2appducontents.E2ConnectionUpdateItemIes {
	connectionUpdateItemIes := &e2appducontents.E2ConnectionUpdateItemIes{
		Id:          int32(v2.ProtocolIeIDE2connectionUpdateItem),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
		Value: &e2appducontents.E2ConnectionUpdateItemIe{
			E2ConnectionUpdateItemIe: &e2appducontents.E2ConnectionUpdateItemIe_E2ConnectionUpdateItem{
				E2ConnectionUpdateItem: &e2appducontents.E2ConnectionUpdateItem{
					TnlInformation: c.tnlInfo,
					TnlUsage:       c.tnlUsage,
				},
			},
		},
	}
	return connectionUpdateItemIes
}
