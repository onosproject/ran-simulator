// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package channels

import (
	"context"
	"io"
	"net"

	e2appdudescriptions "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-descriptions"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/asn1cgo"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

const defaultRecvBufSize = 1024 * 4

var log = logging.GetLogger("protocols", "e2")

// Options is E2 connection options
type Options struct {
	RecvBufferSize int
}

// Option is an E2 connection option
type Option func(*Options)

// WithRecvBuffer sets the connection receive buffer size
func WithRecvBuffer(size int) Option {
	return func(options *Options) {
		options.RecvBufferSize = size
	}
}

// Channel is the base interface for E2 channels
type Channel interface {
	io.Closer
	// Context returns the channel context
	Context() context.Context
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

// newThreadSafeChannel creates a new thread safe channel
func newThreadSafeChannel(conn net.Conn, opts ...Option) *threadSafeChannel {
	options := Options{
		RecvBufferSize: defaultRecvBufSize,
	}
	for _, opt := range opts {
		opt(&options)
	}
	ctx, cancel := context.WithCancel(context.Background())
	channel := &threadSafeChannel{
		conn:    conn,
		sendCh:  make(chan asyncMessage),
		recvCh:  make(chan e2appdudescriptions.E2ApPdu),
		options: options,
		ctx:     ctx,
		cancel:  cancel,
	}
	channel.open()
	return channel
}

// threadSafeChannel is a thread-safe Channel implementation
type threadSafeChannel struct {
	conn    net.Conn
	sendCh  chan asyncMessage
	recvCh  chan e2appdudescriptions.E2ApPdu
	options Options
	ctx     context.Context
	cancel  context.CancelFunc
}

func (c *threadSafeChannel) Context() context.Context {
	return c.ctx
}

func (c *threadSafeChannel) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *threadSafeChannel) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *threadSafeChannel) open() {
	go c.processSends()
	go c.processRecvs()
}

// send sends a message on the connection
func (c *threadSafeChannel) send(msg *e2appdudescriptions.E2ApPdu) error {
	errCh := make(chan error, 1)
	c.sendCh <- asyncMessage{
		msg:   *msg,
		errCh: errCh,
	}
	return <-errCh
}

// processSends processes the send channel
func (c *threadSafeChannel) processSends() {
	for msg := range c.sendCh {
		err := c.processSend(msg.msg)
		if err == io.EOF {
			c.Close()
		} else if err != nil {
			msg.errCh <- err
		}
		close(msg.errCh)
	}
}

// processSend processes a send
func (c *threadSafeChannel) processSend(msg e2appdudescriptions.E2ApPdu) error {
	bytes, err := asn1cgo.PerEncodeE2apPdu(&msg)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(bytes)
	return err
}

// recv receives a message on the connection
func (c *threadSafeChannel) recv() (*e2appdudescriptions.E2ApPdu, error) {
	msg, ok := <-c.recvCh
	if !ok {
		return nil, io.EOF
	}
	return &msg, nil
}

// processRecvs processes the receive channel
func (c *threadSafeChannel) processRecvs() {
	buf := make([]byte, c.options.RecvBufferSize)
	for {
		n, err := c.conn.Read(buf)
		if err == io.EOF {
			c.Close()
			return
		}
		if err != nil {
			log.Error(err)
		} else {
			err := c.processRecv(buf[:n])
			if err != nil {
				log.Error(err)
			}
		}
	}
}

// processRecvs processes the receive channel
func (c *threadSafeChannel) processRecv(bytes []byte) error {
	msg, err := asn1cgo.PerDecodeE2apPdu(bytes)
	if err != nil {
		return err
	}
	c.recvCh <- *msg
	return nil
}

func (c *threadSafeChannel) Close() error {
	defer func() {
		if err := recover(); err != nil {
			log.Debug("recovering from panic:", err)
		}
	}()
	close(c.sendCh)
	close(c.recvCh)
	c.cancel()
	return c.conn.Close()
}

type asyncMessage struct {
	msg   e2appdudescriptions.E2ApPdu
	errCh chan error
}
