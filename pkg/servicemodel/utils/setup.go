// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package utils

import (
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2ap-commondatatypes"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("servicemodel", "utils", "setup")

// SetupRequest setup request
type SetupRequest struct {
	ranFunctions types.RanFunctions
	plmnID       string
}

// NewSetupRequest creates a new setup request
func NewSetupRequest(options ...func(*SetupRequest)) (*SetupRequest, error) {
	setup := &SetupRequest{}

	for _, option := range options {
		option(setup)
	}

	return setup, nil
}

// WithRanFunctions sets ran functions
func WithRanFunctions(ranFunctions types.RanFunctions) func(*SetupRequest) {
	return func(request *SetupRequest) {
		request.ranFunctions = ranFunctions
	}
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID string) func(*SetupRequest) {
	return func(request *SetupRequest) {
		request.plmnID = plmnID

	}
}

// WithE2NodeID sets E2 node ID
func WithE2NodeID() func(*SetupRequest) {
	return func(request *SetupRequest) {

	}
}

// CreateSetupRequest creates e2 setup request
func CreateSetupRequest(request *SetupRequest) (setupRequest *e2appducontents.E2SetupRequest) {
	ranFunctionList := e2appducontents.E2SetupRequestIes_E2SetupRequestIes10{
		Id:          int32(v1beta1.ProtocolIeIDRanfunctionsAdded),
		Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2appducontents.RanfunctionsList{
			Value: make([]*e2appducontents.RanfunctionItemIes, 0),
		},
	}

	for id, ranFunctionID := range request.ranFunctions {
		ranFunction := e2appducontents.RanfunctionItemIes{
			E2ApProtocolIes10: &e2appducontents.RanfunctionItemIes_RanfunctionItemIes8{
				Id:          int32(v1beta1.ProtocolIeIDRanfunctionItem),
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
				},
			},
		}
		ranFunctionList.Value.Value = append(ranFunctionList.Value.Value, &ranFunction)
	}

	e2SetupRequest := &e2appducontents.E2SetupRequest{
		ProtocolIes: &e2appducontents.E2SetupRequestIes{
			E2ApProtocolIes3: &e2appducontents.E2SetupRequestIes_E2SetupRequestIes3{
				Id:          int32(v1beta1.ProtocolIeIDGlobalE2nodeID),
				Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
				Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
				Value: &e2apies.GlobalE2NodeId{
					GlobalE2NodeId: &e2apies.GlobalE2NodeId_GNb{
						GNb: &e2apies.GlobalE2NodeGnbId{
							GlobalGNbId: &e2apies.GlobalgNbId{
								PlmnId: &e2ap_commondatatypes.PlmnIdentity{
									Value: []byte(request.plmnID),
								},
								GnbId: &e2apies.GnbIdChoice{
									GnbIdChoice: &e2apies.GnbIdChoice_GnbId{
										GnbId: &e2ap_commondatatypes.BitString{
											Value: 0x9bcd4,
											Len:   22,
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
	err := e2SetupRequest.Validate()
	if err != nil {
		log.Warnf("Validation error %s", err.Error())
	}
	log.Debugf("Created E2SetupRequest %v", e2SetupRequest)
	return e2SetupRequest
}
