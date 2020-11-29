// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2

import (
	"context"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2/channels"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2/procedures"
	"github.com/onosproject/onos-e2t/pkg/protocols/sctp"
)

// ClientInterface is an interface for E2 client procedures
type ClientInterface procedures.E2NodeProcedures

// ClientChannel is an interface for initiating client procedures
type ClientChannel channels.E2NodeChannel

// NewClient creates a new E2 client
func NewClient(procs ClientInterface) *Client {
	return &Client{
		procs: procs,
	}
}

// Client is an E2 client
type Client struct {
	procs ClientInterface
}

func (c *Client) Connect(ctx context.Context, address string) (ClientChannel, error) {
	conn, err := sctp.Dial(ctx, address)
	if err != nil {
		return nil, err
	}
	channel := channels.NewE2NodeChannel(conn, c.procs)
	return channel, nil
}
