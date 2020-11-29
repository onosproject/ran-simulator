// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package sctp

import (
	"context"
	"github.com/ishidawataru/sctp"
	"net"
	"time"
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

	ip, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return nil, err
	}

	var init sctp.InitMsg
	now := time.Now()
	if deadline, ok := ctx.Deadline(); ok && deadline.After(now) {
		init.MaxInitTimeout = uint16(deadline.Sub(now).Milliseconds())
	}

	addr := &sctp.SCTPAddr{
		IPAddrs: []net.IPAddr{*ip},
	}
	conn, err := sctp.DialSCTPExt("sctp", nil, addr, init)
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
