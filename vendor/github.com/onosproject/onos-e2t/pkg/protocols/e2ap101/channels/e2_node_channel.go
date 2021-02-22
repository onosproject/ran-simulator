// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package channels

import (
	"context"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2appdudescriptions "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-descriptions"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2ap101/procedures"
	"github.com/onosproject/onos-e2t/pkg/utils/async"
	"io"
	"net"
)

// E2NodeHandler is a function for wrapping an E2NodeChannel
type E2NodeHandler func(channel E2NodeChannel) procedures.E2NodeProcedures

// E2NodeChannel is a channel for an E2 node
type E2NodeChannel interface {
	Channel
	procedures.RICProcedures
}

// NewE2NodeChannel creates a new E2 node channel
func NewE2NodeChannel(conn net.Conn, handler E2NodeHandler, opts ...Option) E2NodeChannel {
	parent := newThreadSafeChannel(conn, opts...)
	channel := &e2NodeChannel{
		threadSafeChannel: parent,
	}
	procs := handler(channel)
	channel.e2Setup = procedures.NewE2SetupInitiator(parent.send)
	channel.ricControl = procedures.NewRICControlProcedure(parent.send, procs)
	channel.ricIndication = procedures.NewRICIndicationInitiator(parent.send)
	channel.ricSubscription = procedures.NewRICSubscriptionProcedure(parent.send, procs)
	channel.ricSubscriptionDelete = procedures.NewRICSubscriptionDeleteProcedure(parent.send, procs)
	channel.open()
	return channel
}

// e2NodeChannel is an E2 node channel
type e2NodeChannel struct {
	*threadSafeChannel
	e2Setup               *procedures.E2SetupInitiator
	ricControl            *procedures.RICControlProcedure
	ricIndication         *procedures.RICIndicationInitiator
	ricSubscription       *procedures.RICSubscriptionProcedure
	ricSubscriptionDelete *procedures.RICSubscriptionDeleteProcedure
}

func (c *e2NodeChannel) open() {
	go c.recvPDUs()
}

func (c *e2NodeChannel) recvPDUs() {
	for {
		pdu, err := c.recv()
		if err == io.EOF {
			c.Close()
			return
		}
		if err != nil {
			log.Error(err)
		} else {
			c.recvPDU(pdu)
		}
	}
}

func (c *e2NodeChannel) recvPDU(pdu *e2appdudescriptions.E2ApPdu) {
	if c.e2Setup.Matches(pdu) {
		go c.e2Setup.Handle(pdu)
	} else if c.ricControl.Matches(pdu) {
		go c.ricControl.Handle(pdu)
	} else if c.ricIndication.Matches(pdu) {
		c.ricIndication.Handle(pdu)
	} else if c.ricSubscription.Matches(pdu) {
		go c.ricSubscription.Handle(pdu)
	} else if c.ricSubscriptionDelete.Matches(pdu) {
		go c.ricSubscriptionDelete.Handle(pdu)
	} else {
		log.Errorf("Unsupported E2AP message: %+v", pdu)
	}
}

func (c *e2NodeChannel) E2Setup(ctx context.Context, request *e2appducontents.E2SetupRequest) (response *e2appducontents.E2SetupResponse, failure *e2appducontents.E2SetupFailure, err error) {
	return c.e2Setup.Initiate(ctx, request)
}

func (c *e2NodeChannel) RICIndication(ctx context.Context, request *e2appducontents.Ricindication) (err error) {
	return c.ricIndication.Initiate(ctx, request)
}

func (c *e2NodeChannel) Close() error {
	procedures := []procedures.ElementaryProcedure{
		c.e2Setup,
		c.ricControl,
		c.ricIndication,
		c.ricSubscription,
		c.ricSubscriptionDelete,
	}
	err := async.Apply(len(procedures), func(i int) error {
		return procedures[i].Close()
	})
	if err != nil {
		return err
	}
	return c.threadSafeChannel.Close()
}

var _ E2NodeChannel = &e2NodeChannel{}
