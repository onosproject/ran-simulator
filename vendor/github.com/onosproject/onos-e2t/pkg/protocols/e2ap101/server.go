// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2

import (
	"github.com/onosproject/onos-e2t/pkg/protocols/e2ap101/channels"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2ap101/procedures"
	"github.com/onosproject/onos-e2t/pkg/protocols/sctp"
	"net"
)

// ServerHandler is a server channel handler
type ServerHandler func(channel ServerChannel) ServerInterface

// ServerInterface is an E2 server interface
type ServerInterface procedures.RICProcedures

// ServerChannel is an interface for initiating E2 server procedures
type ServerChannel channels.RICChannel

// NewServer creates a new E2 server
func NewServer(opts ...sctp.ServerOption) *Server {
	return &Server{
		server: sctp.NewServer(opts...),
	}
}

// Server is an E2 server
type Server struct {
	server *sctp.Server
}

// Serve starts the server
func (s *Server) Serve(handler ServerHandler) error {
	return s.server.Serve(func(conn net.Conn) {
		channels.NewRICChannel(conn, func(channel channels.RICChannel) procedures.RICProcedures {
			return handler(channel)
		})
	})
}

// Stop stops the server serving
func (s *Server) Stop() error {
	return s.server.Stop()
}
