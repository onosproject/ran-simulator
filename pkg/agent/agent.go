// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package agent

import (
	"context"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2ap-commondatatypes"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

// Agent is an E2 agent
type Agent interface {
	// Start starts the agent
	Start() error

	// Stop stops the agent
	Stop() error
}

// NewE2Agent creates a new E2 agent
func NewE2Agent(registry *registry.ServiceModelRegistry, address string) Agent {
	return &e2Agent{
		address:  address,
		registry: registry,
	}
}

// e2Agent is an E2 agent
type e2Agent struct {
	address  string
	channel  e2.ClientChannel
	registry *registry.ServiceModelRegistry
}

func (a *e2Agent) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	ranFuncID := registry.ID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	switch ranFuncID {
	case registry.KPM:
		var kpmService kpm.ServiceModel
		err = a.registry.GetServiceModel(ranFuncID, &kpmService)
		if err != nil {
			return nil, nil, err
		}
		return kpmService.RICControl(ctx, request)

	}
	return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

}

func (a *e2Agent) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	ranFuncID := registry.ID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	switch ranFuncID {
	case registry.KPM:
		var kpmService kpm.ServiceModel
		err = a.registry.GetServiceModel(ranFuncID, &kpmService)
		if err != nil {
			return nil, nil, err
		}
		return kpmService.RICSubscription(ctx, request)

	}
	return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

}

func (a *e2Agent) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	ranFuncID := registry.ID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	switch ranFuncID {
	case registry.KPM:
		var kpmService kpm.ServiceModel
		err = a.registry.GetServiceModel(ranFuncID, &kpmService)
		if err != nil {
			return nil, nil, err
		}
		return kpmService.RICSubscriptionDelete(ctx, request)

	}
	return nil, nil, errors.New(errors.NotSupported, "ran function id %v is not supported", ranFuncID)

}

func (a *e2Agent) Start() error {
	client := e2.NewClient(a)
	channel, err := client.Connect(context.Background(), a.address)
	if err != nil {
		return err
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
									Value: []byte{'o', 'n', 'f'},
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
			E2ApProtocolIes10: &e2appducontents.E2SetupRequestIes_E2SetupRequestIes10{
				Id:          int32(v1beta1.ProtocolIeIDRanfunctionsAdded),
				Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
				Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
				Value: &e2appducontents.RanfunctionsList{
					Value: []*e2appducontents.RanfunctionItemIes{}, // TODO: Add RAN functions
				},
			},
		},
	}
	_, e2SetupFailure, err := channel.E2Setup(context.Background(), e2SetupRequest)
	if err != nil {
		return errors.NewUnknown("E2 setup failed: %v", err)
	} else if e2SetupFailure != nil {
		return errors.NewInvalid("E2 setup failed")
	}

	a.channel = channel
	return nil
}

func (a *e2Agent) Stop() error {
	if a.channel != nil {
		return a.channel.Close()
	}
	return nil
}

var _ Agent = &e2Agent{}

var _ e2.ClientInterface = &e2Agent{}
