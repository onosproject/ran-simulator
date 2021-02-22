// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2

import (
	"context"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2ap101/channels"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2ap101/procedures"
	"github.com/onosproject/onos-e2t/pkg/protocols/sctp"
)

// ClientHandler is a client handler function
type ClientHandler func(channel ClientChannel) ClientInterface

// ClientInterface is an interface for E2 client procedures
type ClientInterface procedures.E2NodeProcedures

// ClientChannel is an interface for initiating client procedures
type ClientChannel channels.E2NodeChannel

// Connect connects to the given address
func Connect(ctx context.Context, address string, handler ClientHandler) (ClientChannel, error) {
	conn, err := sctp.Dial(ctx, address)
	if err != nil {
		return nil, err
	}
	channel := channels.NewE2NodeChannel(conn, func(channel channels.E2NodeChannel) procedures.E2NodeProcedures {
		return handler(channel)
	})
	return channel, nil
}
