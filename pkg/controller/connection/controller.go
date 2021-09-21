// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package connection

import (
	"context"
	"fmt"
	"time"

	"github.com/onosproject/ran-simulator/pkg/model"

	"github.com/onosproject/ran-simulator/pkg/tranidpool"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/utils/e2ap/configupdate"

	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap"
	e2connection "github.com/onosproject/ran-simulator/pkg/e2agent/connection"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/ran-simulator/pkg/store/connections"

	"github.com/onosproject/onos-lib-go/pkg/controller"
)

var log = logging.GetLogger("controller", "connection")

const defaultTimeout = 30 * time.Second
const queueSize = 100

// NewController returns a new connection controller. This controller is responsible to open and close
// E2 connections that are the result of the E2 Connection Update procedure or E2 Configuration update procedure
func NewController(connections connections.Store,
	transactionIDPool *tranidpool.TransactionIDPool, node model.Node, model *model.Model) *controller.Controller {
	c := controller.NewController("E2Connections")
	c.Watch(&Watcher{
		connections: connections,
	})

	c.Reconcile(&Reconciler{
		connections:       connections,
		transactionIDPool: transactionIDPool,
		node:              node,
		model:             model,
	})
	return c
}

// Reconciler is a E2 connection reconciler
type Reconciler struct {
	connections       connections.Store
	transactionIDPool *tranidpool.TransactionIDPool
	node              model.Node
	model             *model.Model
}

// Reconcile reconciles the state of a device change
func (r *Reconciler) Reconcile(id controller.ID) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	log.Info("Reconciling Connection:", id.Value)
	connectionID := id.Value.(connections.ConnectionID)

	connection, err := r.connections.Get(ctx, connectionID)
	if err != nil {
		return controller.Result{}, err
	}

	switch connection.Status.Phase {
	case connections.Open:
		return r.reconcileOpenConnection(connection)
	case connections.Closed:
		return r.reconcileClosedConnection(connection)
	}

	return controller.Result{}, nil
}

func (r *Reconciler) reconcileOpenConnection(connection *connections.Connection) (controller.Result, error) {

	// If the connection state is in Initialized  state returns with nil error
	if connection.Status.State == connections.Initialized {
		return controller.Result{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	addr := fmt.Sprintf("%s:%d", connection.ID.GetRICIPAddress(), connection.ID.GetRICPort())

	if connection.Status.State == connections.Disconnected {
		e2Connection := e2connection.NewE2Connection()
		client, err := e2.Connect(ctx, addr, func(channel e2.ClientConn) e2.ClientInterface {
			return e2Connection
		})

		if err != nil {
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			return controller.Result{}, err
		}

		connection.Client = client
		connection.Status.State = connections.Connected

		err = r.connections.Update(ctx, connection)
		if err != nil {
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
	}
	if connection.Status.State == connections.Connected {
		log.Debugf("Sending configuration update for connection: %+v", connection)
		plmnID := ransimtypes.NewUint24(uint32(r.model.PlmnID))
		transactionID, err := r.transactionIDPool.NewID()
		if err != nil {
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
		log.Debugf("Test Trans ID:", transactionID)
		configUpdate, err := configupdate.NewConfigurationUpdate(configupdate.WithTransactionID(int32(transactionID)),
			configupdate.WithE2NodeID(uint64(r.node.GnbID)),
			configupdate.WithPlmnID(plmnID.Value())).Build()
		if err != nil {
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			r.transactionIDPool.Release(transactionID)
			return controller.Result{}, err
		}
		configUpdateAck, configUpdateFailure, err := connection.Client.E2ConfigurationUpdate(ctx, configUpdate)
		if err != nil {
			r.transactionIDPool.Release(transactionID)
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
		if configUpdateFailure != nil {
			r.transactionIDPool.Release(transactionID)
			log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
			return controller.Result{}, err
		}

		if configUpdateAck != nil {
			connection.Status.State = connections.Initialized
			err = r.connections.Update(ctx, connection)
			if err != nil {
				log.Warnf("Failed to reconcile opening connection %+v: %s", connection, err)
				return controller.Result{}, err
			}
		}
		r.transactionIDPool.Release(transactionID)
	}

	return controller.Result{}, nil

}

func (r *Reconciler) reconcileClosedConnection(connection *connections.Connection) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if connection.Status.State == connections.Disconnected {
		err := r.connections.Remove(ctx, connection.ID)
		if err != nil {
			log.Warnf("Failed to reconcile closing connection %+v: %s", connection, err)
			return controller.Result{}, err
		}
	}

	if connection.Status.State == connections.Initialized {
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
	}

	return controller.Result{}, nil
}
