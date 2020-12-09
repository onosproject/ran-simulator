// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package sctp

import (
	"context"
	"github.com/ishidawataru/sctp"
	"net"
	"strconv"
)

// DialOptions is SCTP options
type DialOptions struct {
	WriteBufferSize int
	ReadBufferSize  int
}

// DialOption is an SCTP option function
type DialOption func(*DialOptions)

// WithWriteBuffer sets the write buffer size
func WithWriteBuffer(size int) DialOption {
	return func(options *DialOptions) {
		options.WriteBufferSize = size
	}
}

// WithReadBuffer sets the read buffer size
func WithReadBuffer(size int) DialOption {
	return func(options *DialOptions) {
		options.ReadBufferSize = size
	}
}

// Dial connects to an SCTP server
func Dial(ctx context.Context, address string, opts ...DialOption) (net.Conn, error) {
	options := DialOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return nil, err
	}

	addr := &sctp.SCTPAddr{
		IPAddrs: []net.IPAddr{*ip},
		Port:    port,
	}
	conn, err := sctp.DialSCTP("sctp", nil, addr)
	if err != nil {
		return nil, err
	}

	if options.WriteBufferSize != 0 {
		if err := conn.SetWriteBuffer(options.WriteBufferSize); err != nil {
			return nil, err
		}
	}
	if options.ReadBufferSize != 0 {
		if err := conn.SetReadBuffer(options.ReadBufferSize); err != nil {
			return nil, err
		}
	}
	return conn, nil
}
