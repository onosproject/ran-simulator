// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package controller

import (
	"context"
	"fmt"
	"time"

	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap101"
	e2channel "github.com/onosproject/ran-simulator/pkg/e2agent/channel"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/ran-simulator/pkg/store/channels"

	"github.com/onosproject/onos-lib-go/pkg/controller"
)

var log = logging.GetLogger("e2agent", "controller")

const defaultTimeout = 30 * time.Second
const queueSize = 100

// NewController returns a new network controller
func NewController(channelStore channels.Store) *controller.Controller {
	c := controller.NewController("Connections")
	c.Watch(&ChannelWatcher{
		channels: channelStore,
	})

	c.Reconcile(&Reconciler{
		channels: channelStore,
	})
	return c
}

// Reconciler is a E2 channel reconciler
type Reconciler struct {
	channels channels.Store
}

// Reconcile reconciles the state of a device change
func (r *Reconciler) Reconcile(id controller.ID) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	log.Info("Reconciling Channel:", id.Value)
	channelID := id.Value.(channels.ChannelID)

	channel, err := r.channels.Get(ctx, channelID)
	if err != nil {
		return controller.Result{}, err
	}

	switch channel.Status.Phase {
	case channels.Open:
		return r.reconcileOpenChannel(channel)
	case channels.Closed:
		return r.reconcileClosedChannel(channel)
	}

	return controller.Result{}, nil
}

func (r *Reconciler) reconcileOpenChannel(channel *channels.Channel) (controller.Result, error) {

	if channel.Status.State != channels.Pending {
		return controller.Result{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	addr := fmt.Sprintf("%s:%d", channel.ID.GetRICAddress(), channel.ID.GetRICPort())

	e2Channel := e2channel.NewE2Channel()
	client, err := e2.Connect(context.TODO(), addr,
		func(channel e2.ClientChannel) e2.ClientInterface {
			return e2Channel
		},
	)

	if err != nil {
		log.Warnf("Failed to reconcile opening channel %+v: %s", channel, err)
		channel.Status.State = channels.Failed
		return controller.Result{}, err
	}

	channel.Client = client
	channel.Status.State = channels.Completed

	err = r.channels.Update(ctx, channel)
	if err != nil {
		log.Warnf("Failed to reconcile opening channel %+v: %s", channel, err)
		channel.Status.State = channels.Failed
		return controller.Result{}, err
	}

	// TODO use configuration update to inform E2T about new connection

	return controller.Result{}, nil

}

func (r *Reconciler) reconcileClosedChannel(channel *channels.Channel) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if channel.Status.State == channels.Completed {
		err := r.channels.Remove(ctx, channel.ID)
		if err != nil {
			log.Warnf("Failed to reconcile closing channel %+v: %s", channel, err)
			return controller.Result{}, err
		}
	}

	if channel.Status.State == channels.Pending {
		// TODO use configuration update to inform E2T that E2 node is intended to close the connection
		//      (i.e. before calling close function)
		err := channel.Client.Close()
		if err != nil {
			log.Warnf("Failed to reconcile closing channel %+v: %s", channel, err)
			return controller.Result{}, err
		}
		channel.Status.State = channels.Completed
		err = r.channels.Update(ctx, channel)
		if err != nil {
			log.Warnf("Failed to reconcile closing channel %+v: %s", channel, err)
			return controller.Result{}, err
		}

	}

	return controller.Result{}, nil
}
