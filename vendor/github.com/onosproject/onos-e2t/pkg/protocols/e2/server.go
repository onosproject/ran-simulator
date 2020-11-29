// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2

import (
	"github.com/onosproject/onos-e2t/pkg/protocols/e2/channels"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2/procedures"
	"github.com/onosproject/onos-e2t/pkg/protocols/sctp"
	"net"
)

// Handler is a server channel handler
type Handler func(channel ServerChannel)

// ServerInterface is an E2 server interface
type ServerInterface procedures.RICProcedures

// ServerChannel is an interface for initiating E2 server procedures
type ServerChannel channels.RICChannel

// NewServer creates a new E2 server
func NewServer(procs ServerInterface) *Server {
	return &Server{
		procs: procs,
	}
}

// Server is an E2 server
type Server struct {
	server *sctp.Server
	procs  ServerInterface
}

// Serve starts the server
func (s *Server) Serve(handler Handler) error {
	return s.server.Serve(func(conn net.Conn) {
		handler(channels.NewRICChannel(conn, s.procs))
	})
}

// Stop stops the server serving
func (s *Server) Stop() error {
	return s.server.Stop()
}
