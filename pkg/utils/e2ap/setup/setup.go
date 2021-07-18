// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package setup

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/utils"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("servicemodel", "utils", "setup")

// Setup setup request
type Setup struct {
	ranFunctions e2aptypes.RanFunctions
	plmnID       ransimtypes.Uint24
	e2NodeID     uint64
}

// NewSetupRequest creates a new setup request
func NewSetupRequest(options ...func(*Setup)) *Setup {
	setup := &Setup{}

	for _, option := range options {
		option(setup)
	}

	return setup
}

// WithRanFunctions sets ran functions
func WithRanFunctions(ranFunctions e2aptypes.RanFunctions) func(*Setup) {
	return func(request *Setup) {
		request.ranFunctions = ranFunctions
	}
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID ransimtypes.Uint24) func(*Setup) {
	return func(request *Setup) {
		request.plmnID = plmnID

	}
}

// WithE2NodeID sets E2 node ID
func WithE2NodeID(e2NodeID uint64) func(*Setup) {
	return func(request *Setup) {
		request.e2NodeID = e2NodeID
	}
}

// Build builds e2ap setup request
func (request *Setup) Build() (setupRequest *e2appducontents.E2SetupRequest, err error) {
	//plmnID := types.NewUint24(request.plmnID)
	ranFunctionList := e2appducontents.E2SetupRequestIes_E2SetupRequestIes10{
		Id:          int32(v1beta2.ProtocolIeIDRanfunctionsAdded),
		Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2appducontents.RanfunctionsList{
			Value: make([]*e2appducontents.RanfunctionItemIes, 0),
		},
	}

	for id, ranFunctionID := range request.ranFunctions {
		ranFunction := e2appducontents.RanfunctionItemIes{
			E2ApProtocolIes10: &e2appducontents.RanfunctionItemIes_RanfunctionItemIes8{
				Id:          int32(v1beta2.ProtocolIeIDRanfunctionItem),
				Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
				Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
				Value: &e2appducontents.RanfunctionItem{
					RanFunctionId: &e2apies.RanfunctionId{
						Value: int32(id),
					},
					RanFunctionDefinition: &e2ap_commondatatypes.RanfunctionDefinition{
						Value: []byte(ranFunctionID.Description),
					},
					RanFunctionRevision: &e2apies.RanfunctionRevision{
						Value: int32(ranFunctionID.Revision),
					},
					RanFunctionOid: &e2ap_commondatatypes.RanfunctionOid{
						Value: []byte(ranFunctionID.OID),
					},
				},
			},
		}
		ranFunctionList.Value.Value = append(ranFunctionList.Value.Value, &ranFunction)
	}

	e2SetupRequest := &e2appducontents.E2SetupRequest{
		ProtocolIes: &e2appducontents.E2SetupRequestIes{
			E2ApProtocolIes3: &e2appducontents.E2SetupRequestIes_E2SetupRequestIes3{
				Id:          int32(v1beta2.ProtocolIeIDGlobalE2nodeID),
				Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
				Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
				Value: &e2apies.GlobalE2NodeId{
					GlobalE2NodeId: &e2apies.GlobalE2NodeId_GNb{
						GNb: &e2apies.GlobalE2NodeGnbId{
							GlobalGNbId: &e2apies.GlobalgNbId{
								PlmnId: &e2ap_commondatatypes.PlmnIdentity{
									Value: request.plmnID.ToBytes(),
								},
								GnbId: &e2apies.GnbIdChoice{
									GnbIdChoice: &e2apies.GnbIdChoice_GnbId{
										GnbId: &e2ap_commondatatypes.BitString{
											Value: utils.Uint64ToBitString(request.e2NodeID, 28),
											Len:   28,
										}},
								},
							},
						},
					},
				},
			},
			E2ApProtocolIes10: &ranFunctionList,
		},
	}

	err = e2SetupRequest.Validate()
	if err != nil {
		log.Warnf("Validation error %s", err.Error())
		return nil, err
	}
	log.Debugf("Created E2SetupRequest %v", e2SetupRequest)
	return e2SetupRequest, nil
}
