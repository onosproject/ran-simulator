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

package local

import (
	"context"
	"github.com/atomix/api/proto/atomix/database"
	"github.com/atomix/go-framework/pkg/atomix"
	"github.com/atomix/go-framework/pkg/atomix/cluster"
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/stream"
	"net"
	"time"
)

// NewNode returns a new Atomix Node with a local protocol implementation
func NewNode(lis net.Listener, partitions []primitive.PartitionID) *atomix.Node {
	return atomix.NewNode("local", &database.DatabaseConfig{}, NewProtocol(partitions), atomix.WithLocal(lis))
}

// NewProtocol returns an Atomix Protocol instance
func NewProtocol(partitions []primitive.PartitionID) primitive.Protocol {
	return &Protocol{
		partitions: partitions,
	}
}

// Protocol implements the Atomix protocol in process
type Protocol struct {
	partitions []primitive.PartitionID
	clients    map[primitive.PartitionID]*localClient
}

func (p *Protocol) Start(cluster cluster.Cluster, registry primitive.Registry) error {
	clients := make(map[primitive.PartitionID]*localClient)
	for _, partitionID := range p.partitions {
		context := &localContext{
			partition: partitionID,
		}
		client := &localClient{
			state:   primitive.NewManager(registry, context),
			context: context,
			ch:      make(chan localRequest),
		}
		client.start()
		clients[partitionID] = client
	}
	p.clients = clients
	return nil
}

func (p *Protocol) Partition(partitionID primitive.PartitionID) primitive.Partition {
	return p.clients[partitionID]
}

func (p *Protocol) Partitions() []primitive.Partition {
	partitions := make([]primitive.Partition, 0, len(p.clients))
	for _, partition := range p.clients {
		partitions = append(partitions, partition)
	}
	return partitions
}

func (p *Protocol) Stop() error {
	for _, partition := range p.clients {
		partition.stop()
	}
	return nil
}

type localContext struct {
	partition primitive.PartitionID
	index     primitive.Index
	timestamp time.Time
}

func (c *localContext) NodeID() string {
	return "local"
}

func (c *localContext) PartitionID() primitive.PartitionID {
	return c.partition
}

func (c *localContext) Index() primitive.Index {
	return c.index
}

func (c *localContext) Timestamp() time.Time {
	return c.timestamp
}

type operationType string

const (
	command operationType = "command"
	query   operationType = "query"
)

type localRequest struct {
	op     operationType
	input  []byte
	stream stream.WriteStream
}

type localClient struct {
	state   *primitive.Manager
	context *localContext
	ch      chan localRequest
}

func (c *localClient) MustLeader() bool {
	return false
}

func (c *localClient) IsLeader() bool {
	return false
}

func (c *localClient) Leader() string {
	return ""
}

func (c *localClient) start() {
	go c.processRequests()
}

func (c *localClient) stop() {
	close(c.ch)
}

func (c *localClient) processRequests() {
	for request := range c.ch {
		if request.op == command {
			c.context.index++
			c.context.timestamp = time.Now()
			c.state.Command(request.input, request.stream)
		} else {
			c.state.Query(request.input, request.stream)
		}
	}
}

func (c *localClient) Write(ctx context.Context, input []byte, stream stream.WriteStream) error {
	c.ch <- localRequest{
		op:     command,
		input:  input,
		stream: stream,
	}
	return nil
}

func (c *localClient) Read(ctx context.Context, input []byte, stream stream.WriteStream) error {
	c.ch <- localRequest{
		op:     query,
		input:  input,
		stream: stream,
	}
	return nil
}
