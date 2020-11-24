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
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"sync"
)

const bufSize = 1024 * 1024

// NewTestFactory creates a new cluster factory
func NewTestFactory(services ...Service) *TestFactory {
	return &TestFactory{
		clusters: make(map[NodeID]*localCluster),
		services: services,
	}
}

// TestFactory is a factory for creating test clusters
type TestFactory struct {
	clusters map[NodeID]*localCluster
	services []Service
	mu       sync.RWMutex
}

// NewCluster creates a new test cluster
func (f *TestFactory) NewCluster(nodeID NodeID) (Cluster, error) {
	f.mu.Lock()
	cluster := &localCluster{
		factory:  f,
		nodeID:   nodeID,
		replicas: make(ReplicaSet),
		watchers: make([]chan<- ReplicaSet, 0),
	}
	f.clusters[nodeID] = cluster
	f.mu.Unlock()

	if err := cluster.open(); err != nil {
		return nil, err
	}
	return cluster, nil
}

// closeCluster removes a cluster from the factory
func (f *TestFactory) closeCluster(nodeID NodeID) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.clusters, nodeID)
}

// getClusters returns a list of clusters
func (f *TestFactory) getClusters() []*localCluster {
	f.mu.RLock()
	defer f.mu.RUnlock()
	clusters := make([]*localCluster, 0, len(f.clusters))
	for _, cluster := range f.clusters {
		clusters = append(clusters, cluster)
	}
	return clusters
}

type localCluster struct {
	factory  *TestFactory
	nodeID   NodeID
	replicas ReplicaSet
	watchers []chan<- ReplicaSet
	lis      *bufconn.Listener
	mu       sync.RWMutex
}

func (c *localCluster) open() error {
	server := grpc.NewServer()
	for _, service := range c.factory.services {
		service(c.nodeID, server)
	}

	c.lis = bufconn.Listen(bufSize)
	go func() {
		_ = server.Serve(c.lis)
	}()

	wg := &sync.WaitGroup{}
	clusters := c.factory.getClusters()
	for _, cluster := range clusters {
		wg.Add(1)
		go func(cluster *localCluster) {
			cluster.addReplica(newReplica(ReplicaID(c.nodeID), c.nodeID == cluster.nodeID, func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
				opts = append(opts, grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return c.lis.Dial()
				}))
				return grpc.DialContext(ctx, "local", opts...)
			}))
			c.addReplica(newReplica(ReplicaID(cluster.nodeID), c.nodeID == cluster.nodeID, func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
				opts = append(opts, grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return cluster.lis.Dial()
				}))
				return grpc.DialContext(ctx, "local", opts...)
			}))
			wg.Done()
		}(cluster)
	}
	wg.Wait()
	return nil
}

func (c *localCluster) Node() Node {
	return newNode(c.nodeID)
}

func (c *localCluster) addReplica(replica *Replica) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.replicas[replica.ID]; !ok {
		c.replicas[replica.ID] = replica
	}
	replicas := make(ReplicaSet)
	for id, replica := range c.replicas {
		replicas[id] = replica
	}
	for _, watcher := range c.watchers {
		watcher <- replicas
	}
}

func (c *localCluster) removeReplica(replicaID ReplicaID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.replicas, replicaID)
	replicas := make(ReplicaSet)
	for id, replica := range c.replicas {
		replicas[id] = replica
	}
	for _, watcher := range c.watchers {
		watcher <- replicas
	}
}

func (c *localCluster) Replica(id ReplicaID) *Replica {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.replicas[id]
}

func (c *localCluster) Replicas() ReplicaSet {
	c.mu.RLock()
	defer c.mu.RUnlock()
	replicas := make(ReplicaSet)
	for id, replica := range c.replicas {
		replicas[id] = replica
	}
	return replicas
}

func (c *localCluster) Watch(ch chan<- ReplicaSet) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.watchers = append(c.watchers, ch)
	return nil
}

func (c *localCluster) Close() error {
	c.factory.closeCluster(c.nodeID)
	wg := &sync.WaitGroup{}
	clusters := c.factory.getClusters()
	for _, cluster := range clusters {
		wg.Add(1)
		go func(cluster *localCluster) {
			cluster.removeReplica(ReplicaID(c.nodeID))
		}(cluster)
	}
	wg.Done()
	return nil
}
