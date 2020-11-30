// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package sctp

import (
	"github.com/ishidawataru/sctp"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"net"
)

var log = logging.GetLogger("southbound", "sctp")

const (
	defaultSCTPPort = 36421
)

// ServerOptions is SCTP options
type ServerOptions struct {
	Port            int
	WriteBufferSize int
	ReadBufferSize  int
}

// ServerOption is an SCTP option function
type ServerOption func(*ServerOptions)

// WriteBuffer sets the write buffer size
func WriteBuffer(size int) ServerOption {
	return func(options *ServerOptions) {
		options.WriteBufferSize = size
	}
}

// ReadBuffer sets the read buffer size
func ReadBuffer(size int) ServerOption {
	return func(options *ServerOptions) {
		options.ReadBufferSize = size
	}
}

// Handler is a handler for SCTP connections
type Handler func(conn net.Conn)

// NewServer creates a new southbound server
func NewServer(opts ...ServerOption) *Server {
	return &Server{
		options: applyServerOptions(opts...),
	}
}

// Server is a southbound server
type Server struct {
	options ServerOptions
	lis     net.Listener
}

// Serve starts the server
func (s *Server) Serve(handler Handler) error {
	addr := &sctp.SCTPAddr{
		IPAddrs: []net.IPAddr{},
		Port:    s.options.Port,
	}

	lis, err := sctp.ListenSCTP("sctp", addr)
	if err != nil {
		return err
	}
	s.lis = lis

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Errorf("Failed to accept connection: %v", err)
				continue
			}

			log.Infof("Accepted connection from %s", conn.RemoteAddr())
			sconn := conn.(*sctp.SCTPConn)

			// Configure the connection read buffer
			if s.options.ReadBufferSize != 0 {
				err := sconn.SetWriteBuffer(s.options.WriteBufferSize)
				if err != nil {
					log.Errorf("Failed to configure connection: %v", err)
					continue
				}
			}

			// Configure the connection write buffer
			if s.options.WriteBufferSize != 0 {
				err := sconn.SetWriteBuffer(s.options.WriteBufferSize)
				if err != nil {
					log.Errorf("Failed to configure connection: %v", err)
					continue
				}
			}

			go handler(conn)
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	return s.lis.Close()
}

func applyServerOptions(opts ...ServerOption) ServerOptions {
	options := ServerOptions{
		Port: defaultSCTPPort,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
