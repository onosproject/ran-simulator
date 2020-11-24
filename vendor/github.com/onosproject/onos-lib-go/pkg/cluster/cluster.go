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

package cluster

import (
	"context"
	"github.com/atomix/go-client/pkg/client"
	"github.com/atomix/go-client/pkg/client/peer"
	"google.golang.org/grpc"
	"io"
	"sync"
)

// Service is a gRPC service
type Service func(NodeID, *grpc.Server)

// New creates a new cluster
func New(client *client.Client) (Cluster, error) {
	c := &atomixCluster{
		group:    client.Group(),
		replicas: make(ReplicaSet),
		watchers: make([]chan<- ReplicaSet, 0),
	}
	if err := c.open(); err != nil {
		return nil, err
	}
	return c, nil
}

// Cluster is an interface for interacting with the onos-ric cluster
type Cluster interface {
	io.Closer

	// Node returns the local node
	Node() Node

	// Replica returns a replica by ID
	Replica(ReplicaID) *Replica

	// Replicas returns the set of remote replicas
	Replicas() ReplicaSet

	// Watch watches the cluster for changes to the replicas
	Watch(chan<- ReplicaSet) error
}

// atomixCluster is the default Atomix based cluster implementation
type atomixCluster struct {
	group    *peer.Group
	replicas ReplicaSet
	watchers []chan<- ReplicaSet
	mu       sync.RWMutex
}

// open opens the cluster
func (c *atomixCluster) open() error {
	ch := make(chan peer.Set)
	err := c.group.Watch(context.Background(), ch)
	if err != nil {
		return err
	}
	go func() {
		for peers := range ch {
			c.mu.Lock()
			for id, p := range peers {
				_, ok := c.replicas[ReplicaID(id)]
				if !ok {
					var local bool
					member := c.group.Member()
					if member != nil {
						local = member.ID == id
					}
					c.replicas[ReplicaID(id)] = newReplica(ReplicaID(p.ID), local, func(p *peer.Peer) func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
						return func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
							return p.Connect(ctx, peer.WithDialOptions(opts...))
						}
					}(p))
				}
			}
			for id := range c.replicas {
				_, ok := peers[peer.ID(id)]
				if !ok {
					delete(c.replicas, id)
				}
			}
			c.mu.Unlock()

			replicas := c.Replicas()
			c.mu.RLock()
			for _, watcher := range c.watchers {
				watcher <- replicas
			}
			c.mu.RUnlock()
		}
	}()
	return nil
}

// Node returns the local node
func (c *atomixCluster) Node() Node {
	member := c.group.Member()
	if member == nil {
		return Node{}
	}
	return newNode(NodeID(member.ID))
}

// Replica returns a replica by ID
func (c *atomixCluster) Replica(id ReplicaID) *Replica {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.replicas[id]
}

// Replicas returns the set of all replicas
func (c *atomixCluster) Replicas() ReplicaSet {
	c.mu.RLock()
	defer c.mu.RUnlock()
	replicas := make(ReplicaSet)
	for id, replica := range c.replicas {
		replicas[id] = replica
	}
	return replicas
}

// Watch watches the cluster for changes
func (c *atomixCluster) Watch(ch chan<- ReplicaSet) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.watchers = append(c.watchers, ch)
	return nil
}

// Close closes the cluster
func (c *atomixCluster) Close() error {
	return c.group.Close()
}
