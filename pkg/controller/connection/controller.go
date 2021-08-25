// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package connection

import (
	"context"
	"fmt"
	"time"

	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap101"
	e2connection "github.com/onosproject/ran-simulator/pkg/e2agent/connection"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/ran-simulator/pkg/store/connections"

	"github.com/onosproject/onos-lib-go/pkg/controller"
)

var log = logging.GetLogger("e2agent", "controller")

const defaultTimeout = 30 * time.Second
const queueSize = 100

// NewController returns a new channel controller. This controller is responsible to open and close
// E2 channels that are the result of the E2 Connection Update procedure or E2 Configuration update procedure
func NewController(connections connections.Store) *controller.Controller {
	c := controller.NewController("E2Connections")
	c.Watch(&ConnectionWatcher{
		connections: connections,
	})

	c.Reconcile(&Reconciler{
		connections: connections,
	})
	return c
}

// Reconciler is a E2 channel reconciler
type Reconciler struct {
	connections connections.Store
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
		client, err := e2.Connect(ctx, addr, func(channel e2.ClientChannel) e2.ClientInterface {
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
			connection.Status.State = connections.Disconnected
			return controller.Result{}, err
		}
	}
	if connection.Status.State == connections.Connected {
		log.Debug("Sending configuration update")
		// TODO use configuration update to inform E2T about new connection
		// 		and change channel state to Initialized
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
