// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connection

import (
	"context"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"sync/atomic"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	"fmt"
	"time"

	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"

	"github.com/onosproject/ran-simulator/pkg/model"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/configupdate"

	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap"
	e2connection "github.com/onosproject/ran-simulator/pkg/e2agent/connection"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/ran-simulator/pkg/store/connections"

	"github.com/onosproject/onos-lib-go/pkg/controller"
)

var log = logging.GetLogger()

const defaultTimeout = 30 * time.Second
const queueSize = 100

// NewController returns a new connection controller. This controller is responsible to open and close
// E2 connections that are the result of the E2 Connection Update procedure or E2 Configuration update procedure
func NewController(connections connections.Store, node model.Node, model *model.Model,
	registry *registry.ServiceModelRegistry, subStore *subscriptions.Subscriptions, cellStore cells.Store) *controller.Controller {
	c := controller.NewController("E2Connections")
	c.Watch(&Watcher{
		connections: connections,
	})

	c.Reconcile(&Reconciler{
		connections: connections,
		node:        node,
		model:       model,
		registry:    registry,
		subStore:    subStore,
		cellStore:   cellStore,
	})
	return c
}

// Reconciler is a E2 connection reconciler
type Reconciler struct {
	connections   connections.Store
	node          model.Node
	model         *model.Model
	registry      *registry.ServiceModelRegistry
	subStore      *subscriptions.Subscriptions
	transactionID uint64
	cellStore     cells.Store
}

// Reconcile reconciles the state of a device change
func (r *Reconciler) Reconcile(id controller.ID) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	log.Info("Reconciling Connection:", id.Value)
	connectionID := id.Value.(connections.ConnectionID)

	connection, err := r.connections.Get(ctx, connectionID)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Warn(err)
			return controller.Result{}, err
		}
		log.Info("Connection %s not found", id)
		return controller.Result{}, nil
	}

	switch connection.Status.Phase {
	case connections.Open:
		log.Infof("Reconcile opening connection: %s", connection.ID)
		return r.reconcileOpenConnection(connection)
	case connections.Closed:
		log.Infof("Reconcile closing connection: %s", connection.ID)
		return r.reconcileClosedConnection(connection)
	}

	return controller.Result{}, nil
}

func (r *Reconciler) configureDataConn(ctx context.Context, connection *connections.Connection) (controller.Result, error) {
	log.Infof("Configuring data connection %s", connection.ID)
	plmnID := ransimtypes.NewUint24(uint32(r.model.PlmnID))
	transactionID := atomic.AddUint64(&r.transactionID, 1) % 255

	configUpdate, err := configupdate.NewConfigurationUpdate(
		configupdate.WithTransactionID(int32(transactionID)),
		configupdate.WithE2NodeID(uint64(r.node.GnbID)),
		configupdate.WithPlmnID(plmnID.Value())).
		Build()
	if err != nil {
		log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
		return controller.Result{}, err
	}
	log.Infof("Sending Configuration update request:%+v", configUpdate)
	configUpdateAck, configUpdateFailure, err := connection.Client.E2ConfigurationUpdate(ctx, configUpdate)
	if err != nil {
		log.Warnf("Failed to reconcile configuring connection %+v: %s", connection, err)
		return controller.Result{}, err
	}
	if configUpdateFailure != nil {
		err = errors.NewUnknown("Failed to reconcile configuring connection %+v: %s", connection, err)
		log.Warn(err)
		return controller.Result{}, err
	}

	// Update the state of connection to Configured after receiving config update ack
	if configUpdateAck != nil {
		log.Infof("Config update ack is received:%+v", configUpdateAck)
		connection.Status.State = connections.Configured
		err = r.connections.Update(ctx, connection)
		if err != nil {
			log.Warnf("Failed to reconcile configuring connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
	}

	return controller.Result{}, nil

}

func (r *Reconciler) reconcileOpenConnection(connection *connections.Connection) (controller.Result, error) {

	// If the connection state is in configured  state returns with nil error
	if connection.Status.State == connections.Configured {
		log.Infof("Connection %+v is configured", connection)
		return controller.Result{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	// If the connection state is in configuring state then configure the connection
	if connection.Status.State == connections.Configuring {
		log.Infof("Reconcile Configuring connection: %+v", connection)
		return r.configureDataConn(ctx, connection)
	}

	addr := fmt.Sprintf("%s:%d", connection.ID.GetRICIPAddress(), connection.ID.GetRICPort())
	// If the connection state is in Connecting state then opens a connection to RIC
	// and update connectivity status to Connected
	if connection.Status.State == connections.Connecting {
		log.Infof("Reconcile Connecting connection %+v", connection)
		e2Connection := e2connection.NewE2Connection(
			e2connection.WithNode(r.node),
			e2connection.WithModel(r.model),
			e2connection.WithSMRegistry(r.registry),
			e2connection.WithSubStore(r.subStore),
			e2connection.WithConnectionStore(r.connections),
			e2connection.WithCellStore(r.cellStore))

		client, err := e2.Connect(ctx, addr, func(channel e2.ClientConn) e2.ClientInterface {
			return e2Connection
		})

		if err != nil {
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			return controller.Result{}, err
		}

		e2Connection.SetClient(client)
		connection.Client = client
		connection.Status.State = connections.Connected
		err = r.connections.Update(ctx, connection)
		if err != nil {
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
	}

	// Since the Connection is already established, then E2 NODE CONFIGURATION UPDATE procedure shall be the first E2AP procedure triggered on an
	//  additional TNLA of an already setup E2 interface instance after the TNL association has become operational, and the Near-RT RIC shall
	//  associate the TNLA to the E2 interface instance using the included Global E2 Node ID.
	if connection.Status.State == connections.Connected {
		log.Infof("Reconcile Connected connection %+v", connection)
		connection.Status.State = connections.Configuring
		err := r.connections.Update(ctx, connection)
		if err != nil {
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
		return controller.Result{}, nil
	}

	return controller.Result{}, nil

}

func (r *Reconciler) reconcileClosedConnection(connection *connections.Connection) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if connection.Status.State == connections.Disconnected {
		log.Infof("Reconcile disconnected connection %+v", connection)
		err := r.connections.Remove(ctx, connection.ID)
		if err != nil {
			log.Warnf("Failed to reconcile closing connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
		return controller.Result{}, nil
	}

	if connection.Status.State == connections.Disconnecting {
		log.Infof("Reconcile disconnecting connection %+v", connection)
		// TODO use configuration update to inform E2T that E2 node is intended to close the connection
		//      (i.e. before calling close function)
		err := connection.Client.Close()
		if err != nil {
			log.Warnf("Failed to reconcile closing connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
		connection.Status.State = connections.Disconnected
		err = r.connections.Update(ctx, connection)
		if err != nil {
			log.Warnf("Failed to reconcile closing connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
		return controller.Result{}, nil
	}

	return controller.Result{}, nil
}
