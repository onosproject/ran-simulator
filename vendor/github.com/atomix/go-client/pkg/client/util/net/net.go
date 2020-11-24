// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"google.golang.org/grpc"
	"sync"
)

// Address is the address of a partition
type Address string

// Connect creates a gRPC client connection to the given address
func Connect(address Address) (*grpc.ClientConn, error) {
	return grpc.Dial(
		string(address),
		grpc.WithInsecure())
}

// NewConns returns a new gRPC client connection manager
func NewConns(address Address) *Conns {
	return &Conns{
		Address: address,
		leader:  address,
	}
}

// Conns is a gRPC client connection manager
type Conns struct {
	Address Address
	leader  Address
	conn    *grpc.ClientConn
	mu      sync.RWMutex
}

// Connect gets the connection to the service
func (c *Conns) Connect() (*grpc.ClientConn, error) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn != nil {
		return conn, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	conn = c.conn
	if conn != nil {
		return conn, nil
	}

	conn, err := Connect(c.leader)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return conn, nil
}

// Reconnect reconnects the client to the given leader if necessary
func (c *Conns) Reconnect(leader Address) {
	if leader == "" {
		return
	}

	c.mu.RLock()
	connLeader := c.leader
	c.mu.RUnlock()
	if connLeader == leader {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.leader = leader
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

// Close closes the connections
func (c *Conns) Close() error {
	c.mu.Lock()
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()
	if conn != nil {
		return conn.Close()
	}
	return nil
}
