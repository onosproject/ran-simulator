// Copyright 2020-present Open Networking Foundation.
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

package peer

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sync"
)

// ID is a peer identifier
type ID string

// NewPeer returns a new group peer
func NewPeer(id ID, host string, port int) *Peer {
	return &Peer{
		ID:   id,
		Host: host,
		Port: port,
	}
}

// Peer is a peers group peer
type Peer struct {
	ID   ID
	Host string
	Port int
	conn *grpc.ClientConn
	mu   sync.RWMutex
}

// Connect connects to the member
func (m *Peer) Connect(ctx context.Context, opts ...ConnectOption) (*grpc.ClientConn, error) {
	options := applyConnectOptions(opts...)

	m.mu.RLock()
	conn := m.conn
	m.mu.RUnlock()
	if conn != nil {
		return conn, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conn != nil {
		return m.conn, nil
	}

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", m.Host, m.Port), append([]grpc.DialOption{grpc.WithBlock()}, options.dialOptions...)...)
	if err != nil {
		return nil, err
	}
	m.conn = conn
	return conn, err
}
