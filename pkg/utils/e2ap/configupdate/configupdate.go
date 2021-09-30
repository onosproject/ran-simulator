// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package configupdate

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/utils"

	"github.com/onosproject/onos-e2t/api/e2ap/v2beta1"
	e2apcommondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2beta1/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2beta1/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2beta1/e2ap-pdu-contents"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
)

// ConfigurationUpdate configuration update procedure data structure
type ConfigurationUpdate struct {
	transactionID int32
	plmnID        ransimtypes.Uint24
	e2NodeID      uint64
}

// NewConfigurationUpdate creates a new instance of configuration update
func NewConfigurationUpdate(options ...func(update *ConfigurationUpdate)) *ConfigurationUpdate {
	configUpdate := &ConfigurationUpdate{}

	for _, option := range options {
		option(configUpdate)
	}
	return configUpdate
}

// WithTransactionID sets transaction ID
func WithTransactionID(transID int32) func(update *ConfigurationUpdate) {
	return func(configUpdate *ConfigurationUpdate) {
		configUpdate.transactionID = transID
	}
}

// WithE2NodeID sets E2 node ID
func WithE2NodeID(e2NodeID uint64) func(update *ConfigurationUpdate) {
	return func(configUpdate *ConfigurationUpdate) {
		configUpdate.e2NodeID = e2NodeID
	}
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID ransimtypes.Uint24) func(update *ConfigurationUpdate) {
	return func(configUpdate *ConfigurationUpdate) {
		configUpdate.plmnID = plmnID

	}
}

// Build builds a configuration update request
func (c *ConfigurationUpdate) Build() (*e2appducontents.E2NodeConfigurationUpdate, error) {
	gE2NodeID := &e2apies.GlobalE2NodeId{
		GlobalE2NodeId: &e2apies.GlobalE2NodeId_GNb{
			GNb: &e2apies.GlobalE2NodeGnbId{
				GlobalGNbId: &e2apies.GlobalgNbId{
					PlmnId: &e2apcommondatatypes.PlmnIdentity{
						Value: c.plmnID.ToBytes(),
					},
					GnbId: &e2apies.GnbIdChoice{
						GnbIdChoice: &e2apies.GnbIdChoice_GnbId{
							GnbId: &asn1.BitString{
								Value: utils.Uint64ToBitString(c.e2NodeID, 28),
								Len:   28,
							}},
					},
				},
			},
		},
	}

	configUpdateNodeID := &e2appducontents.E2NodeConfigurationUpdateIes_E2NodeConfigurationUpdateIes3{
		Id:          int32(v2beta1.ProtocolIeIDGlobalE2nodeID),
		Presence:    int32(e2apcommondatatypes.Presence_PRESENCE_MANDATORY),
		Criticality: int32(e2apcommondatatypes.Criticality_CRITICALITY_REJECT),
		Value:       gE2NodeID,
	}

	configUpdate := &e2appducontents.E2NodeConfigurationUpdate{
		ProtocolIes: &e2appducontents.E2NodeConfigurationUpdateIes{
			E2ApProtocolIes49: &e2appducontents.E2NodeConfigurationUpdateIes_E2NodeConfigurationUpdateIes49{
				Id:          int32(v2beta1.ProtocolIeIDTransactionID),
				Criticality: int32(e2apcommondatatypes.Criticality_CRITICALITY_REJECT),
				Value: &e2apies.TransactionId{
					Value: c.transactionID,
				},
				Presence: int32(e2apcommondatatypes.Presence_PRESENCE_MANDATORY),
			},
			E2ApProtocolIes3: configUpdateNodeID,
		},
	}
	return configUpdate, nil
}
