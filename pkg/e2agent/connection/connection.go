// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connection

import (
	"context"
	"fmt"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/pdubuilder"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	asn1libgo "github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/onosproject/ran-simulator/pkg/utils/f1ap"
	"github.com/onosproject/ran-simulator/pkg/utils/xnap"
	"net"
	"sync/atomic"
	"time"

	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2apcommondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-commondatatypes"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm2"

	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"

	connectionsetupfaileditem "github.com/onosproject/ran-simulator/pkg/utils/e2ap/connectionupdate/connectionSetupFailedItemie"

	"github.com/onosproject/ran-simulator/pkg/e2agent/addressing"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/ran-simulator/pkg/store/connections"

	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/connectionupdate/connectionUpdateitemie"

	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/connectionupdate"

	"github.com/cenkalti/backoff"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/mho"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/rc"
	rcv1 "github.com/onosproject/ran-simulator/pkg/servicemodel/rc/v1"
	controlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/setup"

	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"

	"github.com/onosproject/ran-simulator/pkg/model"

	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

var log = logging.GetLogger()

// E2Connection a client interface for of E2 connection
type E2Connection interface {
	e2.ClientInterface

	Setup() error

	Close() error

	GetClient() e2.ClientConn

	SetClient(e2.ClientConn)
}

type e2Connection struct {
	node            model.Node
	model           *model.Model
	client          e2.ClientConn
	registry        *registry.ServiceModelRegistry
	subStore        *subscriptions.Subscriptions
	connectionStore connections.Store
	ricAddress      addressing.RICAddress
	transactionID   uint64
	cellStore       cells.Store
}

// SetClient sets E2 client
func (e *e2Connection) SetClient(client e2.ClientConn) {
	e.client = client
}

// GetClient returns E2 client
func (e *e2Connection) GetClient() e2.ClientConn {
	return e.client
}

// NewE2Connection creates new E2 connection
func NewE2Connection(opts ...InstanceOption) E2Connection {
	log.Info("Creating a new E2 Connection")
	instanceOptions := &InstanceOptions{}
	for _, option := range opts {
		option(instanceOptions)
	}
	return &e2Connection{
		model:           instanceOptions.model,
		node:            instanceOptions.node,
		registry:        instanceOptions.registry,
		subStore:        instanceOptions.subStore,
		ricAddress:      instanceOptions.ricAddress,
		connectionStore: instanceOptions.connectionStore,
		client:          instanceOptions.e2Client,
		cellStore:       instanceOptions.cellStore,
	}

}

// E2ConnectionUpdate implements E2 connection update procedure
func (e *e2Connection) E2ConnectionUpdate(ctx context.Context, request *e2appducontents.E2ConnectionUpdate) (response *e2appducontents.E2ConnectionUpdateAcknowledge, failure *e2appducontents.E2ConnectionUpdateFailure, err error) {
	log.Info("Received Connection Update request %v", request)
	connectionUpdateItemIes := make([]*e2appducontents.E2ConnectionUpdateItemIes, 0)
	connectionSetupFailedItemIes := make([]*e2appducontents.E2ConnectionSetupFailedItemIes, 0)

	var ies44 *e2appducontents.E2ConnectionUpdateList
	var ies45 *e2appducontents.E2ConnectionUpdateList
	var ies46 *e2appducontents.E2ConnectionUpdateRemoveList
	var trID int32
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDE2connectionUpdateAdd) {
			// E2 Connection To Add list IE
			ies44 = v.GetValue().GetE2ConnectionUpdateAdd()
		}
		if v.Id == int32(v2.ProtocolIeIDE2connectionUpdateModify) {
			// E2 Connection To Modify list IE
			ies45 = v.GetValue().GetE2ConnectionUpdateModify()
		}
		if v.Id == int32(v2.ProtocolIeIDE2connectionUpdateRemove) {
			// E2 Connection Remove list IE
			ies46 = v.GetValue().GetE2ConnectionUpdateRemove()
		}
		if v.Id == int32(v2.ProtocolIeIDTransactionID) {
			// Transaction ID IE
			trID = v.GetValue().GetTransactionId().Value
		}
	}

	// In case the E2 Node receives a E2 CONNECTION UPDATE message without any
	// IE except for Message Type IE and Transaction ID IE, it shall reply with the E2 CONNECTION
	//ACKNOWLEDGE message without performing any updates to the existing connections.
	if ies44 == nil && ies45 == nil && ies46 == nil {
		ack := connectionupdate.NewConnectionUpdate(
			connectionupdate.WithTransactionID(trID)).
			BuildConnectionUpdateAcknowledge()
		return ack, nil, nil

	}

	var ricAddress addressing.RICAddress
	// If E2 Connection To Add List IE is contained in the E2 CONNECTION UPDATE message,
	//  then the E2 Node shall, if supported, use it to establish additional TNL Association(s) and configure
	// for use for RIC services and/or E2 support functions according to the TNL Association Usage IE in the message.
	if ies44 != nil {
		log.Debugf("Adding new connections: %+v", ies44.GetValue())
		connectionUpdateItems := ies44.GetValue()
		for _, connectionUpdateItem := range connectionUpdateItems {
			tnlInfo := connectionUpdateItem.GetValue().GetE2ConnectionUpdateItem().GetTnlInformation()
			tnlUsage := connectionUpdateItem.GetValue().GetE2ConnectionUpdateItem().GetTnlUsage()
			// TODO handle tnlUsage

			ricAddress = e.getRICAddress(tnlInfo)
			log.Debugf("RIC and IP and Port information: %v:%v", ricAddress.IPAddress, ricAddress.Port)

			if ricAddress.IPAddress == nil {
				cause := &e2apies.Cause{
					Cause: &e2apies.Cause_Protocol{
						Protocol: e2apies.CauseProtocol_CAUSE_PROTOCOL_ABSTRACT_SYNTAX_ERROR_FALSELY_CONSTRUCTED_MESSAGE,
					},
				}
				connectionUpdateFailure := connectionupdate.NewConnectionUpdate(
					connectionupdate.WithCause(cause),
					connectionupdate.WithTransactionID(trID),
					connectionupdate.WithTimeToWait(nil)).
					BuildConnectionUpdateFailure()
				return nil, connectionUpdateFailure, nil

			}

			// Adds a new connection in Connecting state
			// to the connection store to trigger reconciliation of a connection
			connectionID := connections.NewConnectionID(ricAddress.IPAddress.String(), ricAddress.Port)
			_, err := e.connectionStore.Get(ctx, connectionID)
			if err == nil {
				log.Debugf("Connection %s does exist", connectionID)
				continue
			}

			connection := &connections.Connection{
				ID: connectionID,
				Status: connections.ConnectionStatus{
					Phase: connections.Open,
					State: connections.Connecting,
				},
			}

			err = e.connectionStore.Add(ctx, connectionID, connection)
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

	// remove connections
	if ies46 != nil {
		log.Debugf("Removing connections: %+v", ies46)
		connectionUpdateRemoveItems := ies46.GetValue()
		for _, connectionUpdateRemoveItem := range connectionUpdateRemoveItems {
			tnlInfo := connectionUpdateRemoveItem.GetValue().GetE2ConnectionUpdateRemoveItem().GetTnlInformation()
			ricAddress = e.getRICAddress(tnlInfo)
			if ricAddress.IPAddress == nil {
				cause := &e2apies.Cause{
					Cause: &e2apies.Cause_Protocol{
						Protocol: e2apies.CauseProtocol_CAUSE_PROTOCOL_ABSTRACT_SYNTAX_ERROR_FALSELY_CONSTRUCTED_MESSAGE,
					},
				}
				connectionUpdateFailure := connectionupdate.NewConnectionUpdate(
					connectionupdate.WithCause(cause),
					connectionupdate.WithTransactionID(trID)).
					BuildConnectionUpdateFailure()
				return nil, connectionUpdateFailure, nil

			}

			connectionID := connections.NewConnectionID(ricAddress.IPAddress.String(), ricAddress.Port)
			connection, err := e.connectionStore.Get(ctx, connectionID)

			if err != nil {
				log.Warn(err)
				if !errors.IsNotFound(err) {
					cause := &e2apies.Cause{
						Cause: &e2apies.Cause_Protocol{
							Protocol: e2apies.CauseProtocol_CAUSE_PROTOCOL_UNSPECIFIED,
						},
					}
					connectionUpdateFailure := connectionupdate.NewConnectionUpdate(
						connectionupdate.WithCause(cause),
						connectionupdate.WithTransactionID(trID)).
						BuildConnectionUpdateFailure()
					return nil, connectionUpdateFailure, nil
				}
				connUpdateItemIe := connectionupdateitem.NewConnectionUpdateItemIe(
					connectionupdateitem.WithTnlInfo(tnlInfo)).
					BuildConnectionUpdateItemIes()
				connectionUpdateItemIes = append(connectionUpdateItemIes, connUpdateItemIe)

			} else {
				connection.Status.Phase = connections.Closed
				connection.Status.State = connections.Disconnecting
				err = e.connectionStore.Update(ctx, connection)
				if err != nil {
					log.Warn(err)
					cause := &e2apies.Cause{
						Cause: &e2apies.Cause_Protocol{
							Protocol: e2apies.CauseProtocol_CAUSE_PROTOCOL_UNSPECIFIED,
						},
					}
					connectionUpdateFailure := connectionupdate.NewConnectionUpdate(
						connectionupdate.WithCause(cause),
						connectionupdate.WithTransactionID(trID)).
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

	// After successful update of E2 interface connection(s), the E2 Node shall reply with the E2 CONNECTION UPDATE ACKNOWLEDGE message to inform
	//  the initiating Near-RT RIC that the requested E2 connection update was performed successfully.
	ack := connectionupdate.NewConnectionUpdate(
		connectionupdate.WithConnectionUpdateItemIes(connectionUpdateItemIes),
		connectionupdate.WithConnectionSetupFailedItemIes(connectionSetupFailedItemIes),
		connectionupdate.WithTransactionID(trID)).
		BuildConnectionUpdateAcknowledge()
	log.Infof("Sending Connection Update Ack: %+v", ack)
	return ack, nil, nil
}

func (e *e2Connection) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	rfID, err := controlutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID := registry.RanFunctionID(*rfID)

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
	case registry.Rcpre2:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICControl(ctx, request)
	case registry.Mho:
		client := sm.Client.(*mho.Mho)
		response, failure, err = client.RICControl(ctx, request)
	case registry.Rc:
		client := sm.Client.(*rcv1.Client)
		response, failure, err = client.RICControl(ctx, request)
	}
	if err != nil {
		return nil, nil, err
	}

	return response, failure, err
}

func (e *e2Connection) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	rfID, err := subutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	registeredRanFuncID := registry.RanFunctionID(*rfID)
	log.Debugf("Received Subscription Request %v for ran function %d", request, registeredRanFuncID)
	sm, err := e.registry.GetServiceModel(registeredRanFuncID)
	if err != nil {
		return nil, nil, err
	}
	rrID, err := subutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	riID, err := subutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}

	id := subscriptions.NewID(*riID, *rrID, *rfID)

	reqID, err := subutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := subutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := subutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}
	if err != nil {
		log.Warn(err)
		// If the target E2 Node receives a RIC SUBSCRIPTION REQUEST
		//  message which contains a RAN Function ID IE that was not previously
		//  announced as a supported RAN function in the E2 Setup procedure or
		//  the RIC Service Update procedure, the target E2 Node shall send the RIC SUBSCRIPTION FAILURE message
		//  to the Near-RT RIC with an appropriate cause value.

		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_RAN_FUNCTION_ID_INVALID,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		failure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}
	subscription, err := subscriptions.NewSubscription(id, request, e.client)
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		failure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}
	err = e.subStore.Add(subscription)
	if err != nil {
		log.Warn(err)
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		failure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	// TODO - Assumes ono-to-one mapping between ran function and server model
	switch sm.RanFunctionID {
	case registry.Rcpre2:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Kpm2:
		client := sm.Client.(*kpm2.Client)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Mho:
		client := sm.Client.(*mho.Mho)
		response, failure, err = client.RICSubscription(ctx, request)
	case registry.Rc:
		client := sm.Client.(*rcv1.Client)
		response, failure, err = client.RICSubscription(ctx, request)

	}
	// Ric subscription is failed
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		failure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	return response, failure, err
}

func (e *e2Connection) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	var ranFunctionID int32
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRanfunctionID) {
			// E2 Connection To Add list IE
			ranFunctionID = v.GetValue().GetRanfunctionId().GetValue()
		}
	}

	ranFuncID := registry.RanFunctionID(ranFunctionID)
	log.Debugf("Received Subscription Delete Request %v for ran function ID %d", request, ranFuncID)

	rrID, err := subdeleteutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	rfID, err := subdeleteutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	riID, err := subdeleteutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}

	subID := subscriptions.NewID(*riID, *rrID, *rfID)
	_, err = e.subStore.Get(subID)
	if err != nil {
		log.Warn(err)
		//  If the target E2 Node receives a RIC SUBSCRIPTION DELETE REQUEST
		//  message containing RIC Request ID IE that is not known, the target
		//  E2 Node shall send the RIC SUBSCRIPTION DELETE FAILURE message
		//  to the Near-RT RIC. The message shall contain the Cause IE with an appropriate value.
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}

		rrID, err := subdeleteutils.GetRequesterID(request)
		if err != nil {
			return nil, nil, err
		}
		rfID, err := subdeleteutils.GetRanFunctionID(request)
		if err != nil {
			return nil, nil, err
		}
		riID, err := subdeleteutils.GetRicInstanceID(request)
		if err != nil {
			return nil, nil, err
		}

		subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
			subdeleteutils.WithRanFuncID(*rfID),
			subdeleteutils.WithRequestID(*rrID),
			subdeleteutils.WithRicInstanceID(*riID),
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
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_RAN_FUNCTION_ID_INVALID,
			},
		}
		rrID, err := subdeleteutils.GetRequesterID(request)
		if err != nil {
			return nil, nil, err
		}
		rfID, err := subdeleteutils.GetRanFunctionID(request)
		if err != nil {
			return nil, nil, err
		}
		riID, err := subdeleteutils.GetRicInstanceID(request)
		if err != nil {
			return nil, nil, err
		}

		subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
			subdeleteutils.WithRanFuncID(*rfID),
			subdeleteutils.WithRequestID(*rrID),
			subdeleteutils.WithRicInstanceID(*riID),
			subdeleteutils.WithCause(cause))
		failure, err := subscriptionDelete.BuildSubscriptionDeleteFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	switch sm.RanFunctionID {
	case registry.Rcpre2:
		client := sm.Client.(*rc.Client)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	case registry.Kpm2:
		client := sm.Client.(*kpm2.Client)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	case registry.Mho:
		client := sm.Client.(*mho.Mho)
		response, failure, err = client.RICSubscriptionDelete(ctx, request)
	case registry.Rc:
		client := sm.Client.(*rcv1.Client)
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

func (e *e2Connection) connectAndSetup() error {
	log.Infof("E2 node %d is starting; attempting to connect", e.node.GnbID)
	b := newExpBackoff()

	// Attempt to connect to the E2T controller; use exponential back-off retry
	count := 0
	connectNotify := func(err error, t time.Duration) {
		count++
		log.Infof("E2 node %d failed to connect; retry after %v; attempt %d", e.node.GnbID, b.GetElapsedTime(), count)
	}

	err := backoff.RetryNotify(e.connect, b, connectNotify)
	if err != nil {
		return err
	}
	log.Infof("E2 node %d connected; attempting setup", e.node.GnbID)

	// Attempt to negotiate E2 setup procedure; use exponential back-off retry
	count = 0
	setupNotify := func(err error, t time.Duration) {
		count++
		log.Infof("E2 node %d failed setup procedure; retry after %v; attempt %d: %+v", e.node.GnbID, b.GetElapsedTime(), count, err)
	}

	err = backoff.RetryNotify(e.setup, b, setupNotify)
	log.Infof("E2 node %d completed connection setup: %+v", e.node.GnbID, err)
	return err

}

func (e *e2Connection) Setup() error {
	err := e.connectAndSetup()
	if err != nil {
		return err
	}

	go func() {
		<-e.client.Context().Done()
		log.Warn("Context is cancelled, reconnecting...")
		controller, err := e.model.GetController(e.node.Controllers[0])
		if err != nil {
			return
		}

		controllerAddresses, err := net.LookupHost(controller.Address)
		if err != nil {
			return
		}

		ricAddress := addressing.RICAddress{
			IPAddress: net.ParseIP(controllerAddresses[0]),
			Port:      uint64(controller.Port),
		}
		e.ricAddress = ricAddress
		err = e.Setup()
		if err != nil {
			return
		}

	}()

	return err
}

func (e *e2Connection) connect() error {
	addr := fmt.Sprintf("%s:%d", e.ricAddress.IPAddress.String(), e.ricAddress.Port)
	log.Info("Connecting to E2T with IP address:", addr)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := e2.Connect(ctx, addr,
		func(channel e2.ClientConn) e2.ClientInterface {
			return e
		},
	)

	if err != nil {
		return err
	}

	e.client = client
	return nil
}

var (
	defaultGnBDUID                         = int64(21)
	defaultRRCVerBytes                     = []byte{0xE0}
	defaultRRCVerLen                       = uint32(3)
	defaultNrCellIDLen                     = uint32(36)
	defaultSulFreqBandIndicationNr         = int32(32)
	defaultFreqBandIndicatorNr             = int32(1)
	defaultMeasureTimingConfigurationBytes = []byte{0xF1, 0xF1, 0xF1}
	defaultRANAC                           = int32(255)
	defaultTacBytes                        = []byte{0x01, 0x01, 0x01}
	defaultSD                              = []byte{0x01, 0x23, 0x45}
	defaultSST                             = []byte{0x01}
	defaultAMFRegionValue                  = []byte{0xdd} // todo need to be changed
	defaultAMFRegionLen                    = uint32(8)
)

func (e *e2Connection) setup() error {
	plmnID := ransimtypes.NewUint24(uint32(e.model.PlmnID))

	configAdditionList := &e2appducontents.E2NodeComponentConfigAdditionList{
		Value: make([]*e2appducontents.E2NodeComponentConfigAdditionItemIes, 0),
	}

	// TODO initialize component interfaces properly. It is just initialized with some default values
	// 	to avoid encoding error.
	e2ncIDF1 := pdubuilder.CreateE2NodeComponentIDF1(21)
	e2ncIDXn := pdubuilder.CreateE2NodeComponentIDXn(&e2apies.GlobalNgRannodeId{
		GlobalNgRannodeId: &e2apies.GlobalNgRannodeId_GNb{
			GNb: &e2apies.GlobalgNbId{
				PlmnId: &e2apcommondatatypes.PlmnIdentity{
					Value: plmnID.ToBytes(),
				},
				GnbId: &e2apies.GnbIdChoice{
					GnbIdChoice: &e2apies.GnbIdChoice_GnbId{
						GnbId: &asn1libgo.BitString{
							Value: utils.Uint64ToBitString(uint64(e.node.GnbID), 22),
							Len:   22,
						},
					},
				},
			},
		},
	})

	sCellItemListF1 := make([]f1ap.SCellItemInfo, 0)
	sCellItemListXn := make([]xnap.XnItemCellInfo, 0)
	nCellItemMapXn := make(map[ransimtypes.NCGI][]xnap.XnItemCellInfo)
	e2NodePlmn := plmnID.ToBytes()
	for _, c := range e.node.Cells {
		m, err := e.cellStore.Get(context.Background(), c)
		if err != nil {
			log.Warnf("failed to fetch cell %+v: %+v", m, err)
		}
		nci := utils.NewNCellIDWithUint64(uint64(ransimtypes.GetNCI(m.NCGI)))
		sCellItem := f1ap.SCellItemInfo{
			PlmnIDBytes:                     e2NodePlmn,
			NrCellIDBytes:                   nci.Bytes(),
			NrCellIDLen:                     defaultNrCellIDLen,
			NrPCI:                           int32(m.PCI),
			SulFreqBandIndicationNr:         defaultSulFreqBandIndicationNr,
			FreqBandIndicatorNr:             defaultFreqBandIndicatorNr,
			NrArfcn:                         int32(m.Earfcn),
			MeasureTimingConfigurationBytes: defaultMeasureTimingConfigurationBytes,
		}
		sCellItemListF1 = append(sCellItemListF1, sCellItem)
		xnSCellItem := xnap.XnItemCellInfo{
			NCGIKey:                         m.NCGI,
			NrCellIDBytes:                   nci.Bytes(),
			NrCellIDLen:                     defaultNrCellIDLen,
			NrPCI:                           int32(m.PCI),
			NrArfcn:                         int32(m.Earfcn),
			SulFreqBand:                     defaultSulFreqBandIndicationNr,
			FreqBand:                        defaultFreqBandIndicatorNr,
			MeasureTimingConfigurationBytes: defaultMeasureTimingConfigurationBytes,
			RanAC:                           defaultRANAC,
		}
		sCellItemListXn = append(sCellItemListXn, xnSCellItem)

		nCellItemListXn := make([]xnap.XnItemCellInfo, 0)
		for _, n := range m.Neighbors {
			nCell, err := e.cellStore.Get(context.Background(), n)
			if err != nil {
				log.Warnf("failed to fetch neighbor cell %+v, err: %+v", nCell, err)
				continue
			}
			neighborNci := utils.NewNCellIDWithUint64(uint64(ransimtypes.GetNCI(nCell.NCGI)))
			xnNCellitem := xnap.XnItemCellInfo{
				NCGIKey:                         nCell.NCGI,
				NrCellIDBytes:                   neighborNci.Bytes(),
				NrCellIDLen:                     defaultNrCellIDLen,
				NrPCI:                           int32(nCell.PCI),
				NrArfcn:                         int32(nCell.Earfcn),
				SulFreqBand:                     defaultSulFreqBandIndicationNr,
				FreqBand:                        defaultFreqBandIndicatorNr,
				MeasureTimingConfigurationBytes: defaultMeasureTimingConfigurationBytes,
				RanAC:                           defaultRANAC,
			}
			nCellItemListXn = append(nCellItemListXn, xnNCellitem)
		}
		nCellItemMapXn[m.NCGI] = nCellItemListXn
	}

	f1SetupRequestBytes, err := f1ap.CreateF1SetupRequest(defaultGnBDUID, defaultRRCVerBytes, defaultRRCVerLen, sCellItemListF1)
	if err != nil {
		return err
	}
	xnSetupRequestBytes, err := xnap.CreateXnSetupRequest(e2NodePlmn, utils.Uint64ToBitString(uint64(e.node.GnbID), 22), defaultTacBytes, []xnap.XnItemSlice{
		{
			Sst: defaultSST,
			Sd:  defaultSD,
		},
	}, xnap.XnItemAMFRegion{AmfRegionID: defaultAMFRegionValue, AmfRegionIDLen: defaultAMFRegionLen}, sCellItemListXn, nCellItemMapXn)
	if err != nil {
		return err
	}
	configComponentAdditionItems := []*types.E2NodeComponentConfigAdditionItem{
		{
			E2NodeComponentType: e2apies.E2NodeComponentInterfaceType_E2NODE_COMPONENT_INTERFACE_TYPE_F1,
			E2NodeComponentID:   e2ncIDF1,
			E2NodeComponentConfiguration: e2apies.E2NodeComponentConfiguration{
				E2NodeComponentRequestPart:  f1SetupRequestBytes,
				E2NodeComponentResponsePart: []byte{0x04, 0x05, 0x06},
			},
		},
		{
			E2NodeComponentType: e2apies.E2NodeComponentInterfaceType_E2NODE_COMPONENT_INTERFACE_TYPE_XN,
			E2NodeComponentID:   e2ncIDXn,
			E2NodeComponentConfiguration: e2apies.E2NodeComponentConfiguration{
				E2NodeComponentRequestPart:  xnSetupRequestBytes,
				E2NodeComponentResponsePart: []byte{0x04, 0x05, 0x06},
			},
		},
	}
	for _, configAdditionItem := range configComponentAdditionItems {
		cui := &e2appducontents.E2NodeComponentConfigAdditionItemIes{
			Id:          int32(v2.ProtocolIeIDE2nodeComponentConfigAdditionItem),
			Criticality: int32(e2apcommondatatypes.Criticality_CRITICALITY_REJECT),
			Value: &e2appducontents.E2NodeComponentConfigAdditionItemIe{
				E2NodeComponentConfigAdditionItemIe: &e2appducontents.E2NodeComponentConfigAdditionItemIe_E2NodeComponentConfigAdditionItem{
					E2NodeComponentConfigAdditionItem: &e2appducontents.E2NodeComponentConfigAdditionItem{
						E2NodeComponentInterfaceType: configAdditionItem.E2NodeComponentType,
						E2NodeComponentId:            configAdditionItem.E2NodeComponentID,
						E2NodeComponentConfiguration: &configAdditionItem.E2NodeComponentConfiguration,
					},
				},
			},
		}
		configAdditionList.Value = append(configAdditionList.Value, cui)
	}

	transactionID := atomic.AddUint64(&e.transactionID, 1) % 255

	setupRequest := setup.NewSetupRequest(
		setup.WithRanFunctions(e.registry.GetRanFunctions()),
		setup.WithPlmnID(plmnID.Value()),
		setup.WithE2NodeID(uint64(e.node.GnbID)),
		setup.WithComponentConfigUpdateList(configAdditionList),
		setup.WithTransactionID(int32(transactionID)))

	e2SetupRequest, err := setupRequest.Build()

	if err != nil {
		log.Error(err)
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	e2SetupAck, e2SetupFailure, err := e.client.E2Setup(ctx, e2SetupRequest)
	if err != nil {
		log.Warn(err)
		return errors.NewUnknown("E2 setup failed: %v", err)
	} else if e2SetupFailure != nil {
		err := errors.NewInvalid("E2 setup failed")
		log.Warn(err)
		return err
	}
	log.Infof("E2 Setup Ack is received:%+v", e2SetupAck)
	// Add connection to the connection store
	connectionID := connections.NewConnectionID(e.ricAddress.IPAddress.String(), e.ricAddress.Port)

	connection := &connections.Connection{
		ID: connectionID,
		Status: connections.ConnectionStatus{
			Phase: connections.Open,
			State: connections.Configured,
		},
		Client: e.client,
	}

	err = e.connectionStore.Add(ctx,
		connectionID, connection)
	if err != nil {
		return err
	}
	return nil
}

func (e *e2Connection) Close() error {
	connectionID := connections.NewConnectionID(e.ricAddress.IPAddress.String(), e.ricAddress.Port)
	log.Debugf("Closing E2 connection with ID %d:", connectionID)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if e.client != nil {
		err := e.client.Close()
		if err != nil {
			return err
		}
		err = e.connectionStore.Remove(ctx, connectionID)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ E2Connection = &e2Connection{}
