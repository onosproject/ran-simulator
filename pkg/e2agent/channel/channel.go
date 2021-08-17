// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package channel

import (
	"context"
	"fmt"
	"time"

	connectionsetupfaileditem "github.com/onosproject/ran-simulator/pkg/utils/e2ap/connectionupdate/connectionSetupFailedItemie"

	"github.com/onosproject/ran-simulator/pkg/e2agent/addressing"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/ran-simulator/pkg/store/channels"

	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/connectionupdate/connectionUpdateitemie"

	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/connectionupdate"

	"github.com/cenkalti/backoff"

	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm2"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/mho"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/rc"
	controlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/setup"

	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"

	"github.com/onosproject/ran-simulator/pkg/model"

	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap101"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

var log = logging.GetLogger("e2agent", "channel")

// E2Channel a new instance of E2 agent
type E2Channel interface {
	e2.ClientInterface

	Start() error

	Stop() error

	Connect() error

	GetClient() e2.ClientChannel
}

type e2Channel struct {
	node         model.Node
	model        *model.Model
	client       e2.ClientChannel
	registry     *registry.ServiceModelRegistry
	subStore     *subscriptions.Subscriptions
	channelStore channels.Store
	ricAddress   addressing.RICAddress
}

func (e *e2Channel) GetClient() e2.ClientChannel {
	return e.client
}

// NewE2Channel creates new E2 channel instance
func NewE2Channel(opts ...InstanceOption) E2Channel {
	log.Info("Creating a new E2 Instance")
	instanceOptions := &InstanceOptions{}
	for _, option := range opts {
		option(instanceOptions)
	}

	return &e2Channel{
		model:        instanceOptions.model,
		node:         instanceOptions.node,
		registry:     instanceOptions.registry,
		subStore:     instanceOptions.subStore,
		ricAddress:   instanceOptions.ricAddress,
		channelStore: instanceOptions.channelStore,
	}

}

// E2ConnectionUpdate implements E2 connection update procedure
func (e *e2Channel) E2ConnectionUpdate(ctx context.Context, request *e2appducontents.E2ConnectionUpdate) (response *e2appducontents.E2ConnectionUpdateAcknowledge, failure *e2appducontents.E2ConnectionUpdateFailure, err error) {
	log.Info("Connection Update request is received %v", request)
	connectionUpdateItemIes := make([]*e2appducontents.E2ConnectionUpdateItemIes, 0)
	connectionSetupFailedItemIes := make([]*e2appducontents.E2ConnectionSetupFailedItemIes, 0)
	ies44 := request.GetProtocolIes().GetE2ApProtocolIes44()
	ies45 := request.GetProtocolIes().GetE2ApProtocolIes45()
	ies46 := request.GetProtocolIes().GetE2ApProtocolIes46()

	// In case the E2 Node receives a E2 CONNECTION UPDATE message without any
	// IE except for Message Type IE and Transaction ID IE, it shall reply with the E2 CONNECTION
	//ACKNOWLEDGE message without performing any updates to the existing connections.
	if ies44 == nil && ies45 == nil && ies46 == nil {
		ack := connectionupdate.NewConnectionUpdate().
			BuildConnectionUpdateAcknowledge()
		return ack, nil, nil

	}

	var ricAddress addressing.RICAddress
	// If E2 Connection To Add List IE is contained in the E2 CONNECTION UPDATE message,
	//  then the E2 Node shall, if supported, use it to establish additional TNL Association(s) and configure
	// for use for RIC services and/or E2 support functions according to the TNL Association Usage IE in the message.
	if ies44 != nil {
		connectionUpdateList := ies44.GetConnectionAdd()
		if connectionUpdateList != nil {
			log.Debug("Adding new connections")
			connectionUpdateItems := connectionUpdateList.Value
			for _, connectionUpdateItem := range connectionUpdateItems {
				tnlInfo := connectionUpdateItem.GetValue().GetTnlInformation()
				tnlUsage := connectionUpdateItem.GetValue().GetTnlUsage()
				// TODO handle tnlUsage

				ricAddress = e.getRICAddress(tnlInfo)

				if ricAddress.IPAddress == nil {
					cause := &e2apies.Cause{
						Cause: &e2apies.Cause_Protocol{
							Protocol: e2apies.CauseProtocol_CAUSE_PROTOCOL_ABSTRACT_SYNTAX_ERROR_FALSELY_CONSTRUCTED_MESSAGE,
						},
					}
					connectionUpdateFailure := connectionupdate.NewConnectionUpdate(connectionupdate.WithCause(cause)).
						BuildConnectionUpdateFailure()
					return nil, connectionUpdateFailure, nil

				}

				// Adds a new channel to the channel store
				channelID := channels.NewChannelID(ricAddress.IPAddress.String(), ricAddress.Port)
				channel := &channels.Channel{
					ID: channelID,
					Status: channels.ChannelStatus{
						Phase: channels.Open,
						State: channels.Pending,
					},
				}

				err := e.channelStore.Add(ctx, channelID, channel)
				if err != nil {
					// If connection is not established then creates a connection setup failed item IE
					// to be reported in ACK
					connSetupFailedItemIe := connectionsetupfaileditem.NewConnectionSetupFailedItemIe(
						connectionsetupfaileditem.WithTnlInfo(tnlInfo)).
						BuildConnectionSetupFailedItemIes()
					connectionSetupFailedItemIes = append(connectionSetupFailedItemIes, connSetupFailedItemIe)
				}

				// If connection is established successfully, creates a connection update item IE
				// to be used in ACK
				connUpdateItemIe := connectionupdateitem.NewConnectionUpdateItemIe(
					connectionupdateitem.WithTnlInfo(tnlInfo),
					connectionupdateitem.WithTnlUsage(tnlUsage)).
					BuildConnectionUpdateItemIes()
				connectionUpdateItemIes = append(connectionUpdateItems, connUpdateItemIe)

			}

		}
	}

	// remove connections
	if ies46 != nil {
		connectionRemoveList := ies46.GetConnectionRemove()
		if connectionRemoveList != nil {
			log.Debug("Removing connections")
			connectionUpdateRemoveItems := connectionRemoveList.GetValue()
			for _, connectionUpdateRemoveItem := range connectionUpdateRemoveItems {
				tnlInfo := connectionUpdateRemoveItem.GetValue().GetTnlInformation()
				ricAddress = e.getRICAddress(tnlInfo)
				if ricAddress.IPAddress == nil {
					cause := &e2apies.Cause{
						Cause: &e2apies.Cause_Protocol{
							Protocol: e2apies.CauseProtocol_CAUSE_PROTOCOL_ABSTRACT_SYNTAX_ERROR_FALSELY_CONSTRUCTED_MESSAGE,
						},
					}
					connectionUpdateFailure := connectionupdate.NewConnectionUpdate(connectionupdate.WithCause(cause)).
						BuildConnectionUpdateFailure()
					return nil, connectionUpdateFailure, nil

				}

				channelID := channels.NewChannelID(ricAddress.IPAddress.String(), ricAddress.Port)
				channel, err := e.channelStore.Get(ctx, channelID)

				if err != nil {
					log.Warn(err)
					cause := &e2apies.Cause{
						Cause: &e2apies.Cause_Protocol{
							Protocol: e2apies.CauseProtocol_CAUSE_PROTOCOL_UNSPECIFIED,
						},
					}
					connectionUpdateFailure := connectionupdate.NewConnectionUpdate(connectionupdate.WithCause(cause)).
						BuildConnectionUpdateFailure()
					return nil, connectionUpdateFailure, nil
				}

				channel.Status.Phase = channels.Closed
				channel.Status.State = channels.Pending

				err = e.channelStore.Update(ctx, channel)
				if err != nil {
					log.Warn(err)
					cause := &e2apies.Cause{
						Cause: &e2apies.Cause_Protocol{
							Protocol: e2apies.CauseProtocol_CAUSE_PROTOCOL_UNSPECIFIED,
						},
					}
					connectionUpdateFailure := connectionupdate.NewConnectionUpdate(connectionupdate.WithCause(cause)).
						BuildConnectionUpdateFailure()
					return nil, connectionUpdateFailure, nil
				}
			}
		}

	}
	// TODO modifying connections
	if ies45 != nil {
		log.Debug("Modifying connections")
	}

	ack := connectionupdate.NewConnectionUpdate(
		connectionupdate.WithConnectionUpdateItemIes(connectionUpdateItemIes),
		connectionupdate.WithConnectionSetupFailedItemIes(connectionSetupFailedItemIes)).
		BuildConnectionUpdateAcknowledge()
	log.Info("Sending Connection Update Ack:", ack)
	return ack, nil, nil
}

func (e *e2Channel) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	ranFuncID := registry.RanFunctionID(controlutils.GetRanFunctionID(request))
	log.Debugf("Received Control Request %+v for ran function %d", request, ranFuncID)
	sm, err := e.registry.GetServiceModel(ranFuncID)
	if err != nil {
		log.Warn(err)
		// TODO If the target E2 Node receives a RIC CONTROL REQUEST message
		//  which contains a RAN Function ID IE that was not previously announced as a s
		//  supported RAN function in the E2 Setup procedure or the RIC Service Update procedure,
		//  or the E2 Node does not support the specific RIC Control procedure action, then
		//  the target E2 Node shall ignore message and send an ERROR INDICATION message to the Near-RT RIC.

		return nil, nil, err
	}
	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		response, failure, err = client.RICControl(ctx, request)
	case registry.Rcpre2:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICControl(ctx, request)
	case registry.Mho:
		client := sm.Client.(*mho.Mho)
		response, failure, err = client.RICControl(ctx, request)
	}
	if err != nil {
		return nil, nil, err
	}

	return response, failure, err
}

func (e *e2Channel) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	ranFuncID := registry.RanFunctionID(subutils.GetRanFunctionID(request))
	log.Debugf("Received Subscription Request %v for ran function %d", request, ranFuncID)
	sm, err := e.registry.GetServiceModel(ranFuncID)
	id := subscriptions.NewID(subutils.GetRicInstanceID(request),
		subutils.GetRequesterID(request),
		subutils.GetRanFunctionID(request))

	if err != nil {
		log.Warn(err)
		// If the target E2 Node receives a RIC SUBSCRIPTION REQUEST
		//  message which contains a RAN Function ID IE that was not previously
		//  announced as a supported RAN function in the E2 Setup procedure or
		//  the RIC Service Update procedure, the target E2 Node shall send the RIC SUBSCRIPTION FAILURE message
		//  to the Near-RT RIC with an appropriate cause value.
		var ricActionsAccepted []*e2aptypes.RicActionID
		ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
		actionList := subutils.GetRicActionToBeSetupList(request)
		reqID := subutils.GetRequesterID(request)
		ranFuncID := subutils.GetRanFunctionID(request)
		ricInstanceID := subutils.GetRicInstanceID(request)

		for _, action := range actionList {
			actionID := e2aptypes.RicActionID(action.Value.RicActionId.Value)
			cause := &e2apies.Cause{
				Cause: &e2apies.Cause_RicRequest{
					RicRequest: e2apies.CauseRic_CAUSE_RIC_RAN_FUNCTION_ID_INVALID,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(reqID),
			subutils.WithRanFuncID(ranFuncID),
			subutils.WithRicInstanceID(ricInstanceID),
			subutils.WithActionsAccepted(ricActionsAccepted),
			subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))
		failure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}
	subscription, err := subscriptions.NewSubscription(id, request, e.client)
	if err != nil {
		return response, failure, err
	}
	err = e.subStore.Add(subscription)
	if err != nil {
		return response, failure, err
	}

	// TODO - Assumes ono-to-one mapping between ran function and server model
	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Rcpre2:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Kpm2:
		client := sm.Client.(*kpm2.Client)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Mho:
		client := sm.Client.(*mho.Mho)
		response, failure, err = client.RICSubscription(ctx, request)

	}
	// Ric subscription is failed
	if err != nil {
		return response, failure, err
	}

	return response, failure, err
}

func (e *e2Channel) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	ranFuncID := registry.RanFunctionID(request.ProtocolIes.E2ApProtocolIes5.Value.Value)
	log.Debugf("Received Subscription Delete Request %v for ran function ID %d", request, ranFuncID)
	subID := subscriptions.NewID(subdeleteutils.GetRicInstanceID(request),
		subdeleteutils.GetRequesterID(request),
		subdeleteutils.GetRanFunctionID(request))
	_, err = e.subStore.Get(subID)
	if err != nil {
		log.Warn(err)
		//  If the target E2 Node receives a RIC SUBSCRIPTION DELETE REQUEST
		//  message containing RIC Request ID IE that is not known, the target
		//  E2 Node shall send the RIC SUBSCRIPTION DELETE FAILURE message
		//  to the Near-RT RIC. The message shall contain the Cause IE with an appropriate value.
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRic_CAUSE_RIC_REQUEST_ID_UNKNOWN,
			},
		}
		subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
			subdeleteutils.WithRanFuncID(subdeleteutils.GetRanFunctionID(request)),
			subdeleteutils.WithRequestID(subdeleteutils.GetRequesterID(request)),
			subdeleteutils.WithRicInstanceID(subdeleteutils.GetRicInstanceID(request)),
			subdeleteutils.WithCause(cause))
		failure, err := subscriptionDelete.BuildSubscriptionDeleteFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil

	}

	sm, err := e.registry.GetServiceModel(ranFuncID)
	if err != nil {
		log.Warn(err)
		//  If the target E2 Node receives a RIC SUBSCRIPTION DELETE REQUEST message contains a
		//  RAN Function ID IE that was not previously announced as a supported RAN function
		//  in the E2 Setup procedure or the RIC Service Update procedure, the target E2 Node
		//  shall send the RIC SUBSCRIPTION DELETE FAILURE message to the Near-RT RIC.
		//  The message shall contain with an appropriate cause value.
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRic_CAUSE_RIC_RAN_FUNCTION_ID_INVALID,
			},
		}
		subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
			subdeleteutils.WithRanFuncID(subdeleteutils.GetRanFunctionID(request)),
			subdeleteutils.WithRequestID(subdeleteutils.GetRequesterID(request)),
			subdeleteutils.WithRicInstanceID(subdeleteutils.GetRicInstanceID(request)),
			subdeleteutils.WithCause(cause))
		failure, err := subscriptionDelete.BuildSubscriptionDeleteFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	switch sm.RanFunctionID {
	case registry.Kpm:
		client := sm.Client.(*kpm.Client)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	case registry.Rcpre2:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	case registry.Kpm2:
		client := sm.Client.(*kpm2.Client)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)

	}
	// Ric subscription delete procedure is failed so we are not going to update subscriptions store
	if err != nil {
		log.Warn(err)
		return response, failure, err
	}

	err = e.subStore.Remove(subID)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	return response, failure, err
}

func (e *e2Channel) Start() error {
	log.Infof("E2 node %d is starting; attempting to connect", e.node.GnbID)
	b := newExpBackoff()

	// Attempt to connect to the E2T controller; use exponential back-off retry
	count := 0
	connectNotify := func(err error, t time.Duration) {
		count++
		log.Infof("E2 node %d failed to connect; retry after %v; attempt %d", e.node.GnbID, b.GetElapsedTime(), count)
	}

	err := backoff.RetryNotify(e.Connect, b, connectNotify)
	if err != nil {
		return err
	}
	log.Infof("E2 node %d connected; attempting setup", e.node.GnbID)

	// Attempt to negotiate E2 setup procedure; use exponential back-off retry
	count = 0
	setupNotify := func(err error, t time.Duration) {
		count++
		log.Infof("E2 node %d failed setup procedure; retry after %v; attempt %d", e.node.GnbID, b.GetElapsedTime(), count)
	}

	err = backoff.RetryNotify(e.setup, b, setupNotify)
	log.Infof("E2 node %d completed connection setup", e.node.GnbID)
	return err
}

func (e *e2Channel) Connect() error {
	addr := fmt.Sprintf("%s:%d", e.ricAddress.IPAddress.String(), e.ricAddress.Port)
	log.Info("Connecting to E2T with IP address:", addr)
	client, err := e2.Connect(context.TODO(), addr,
		func(channel e2.ClientChannel) e2.ClientInterface {
			return e
		},
	)

	if err != nil {
		return err
	}
	// Add channels to the channel store
	channelID := channels.NewChannelID(e.ricAddress.IPAddress.String(), e.ricAddress.Port)

	channel := &channels.Channel{
		ID: channelID,
		Status: channels.ChannelStatus{
			Phase: channels.Open,
			State: channels.Completed,
		},
	}

	err = e.channelStore.Add(context.Background(),
		channelID, channel)
	if err != nil {
		return err
	}
	e.client = client

	return nil
}

func (e *e2Channel) setup() error {
	plmnID := ransimtypes.NewUint24(uint32(e.model.PlmnID))
	setupRequest := setup.NewSetupRequest(
		setup.WithRanFunctions(e.registry.GetRanFunctions()),
		setup.WithPlmnID(plmnID.Value()),
		setup.WithE2NodeID(uint64(e.node.GnbID)))

	e2SetupRequest, err := setupRequest.Build()

	if err != nil {
		log.Error(err)
		return err
	}
	_, e2SetupFailure, err := e.client.E2Setup(context.Background(), e2SetupRequest)
	if err != nil {
		log.Error(err)
		return errors.NewUnknown("E2 setup failed: %v", err)
	} else if e2SetupFailure != nil {
		err := errors.NewInvalid("E2 setup failed")
		log.Error(err)
		return err
	}
	return nil
}

func (e *e2Channel) Stop() error {
	channelID := channels.NewChannelID(e.ricAddress.IPAddress.String(), e.ricAddress.Port)
	log.Debugf("Closing E2 channel with ID %d:", channelID)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if e.client != nil {
		err := e.client.Close()
		if err != nil {
			return err
		}
		err = e.channelStore.Remove(ctx, channelID)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ E2Channel = &e2Channel{}
