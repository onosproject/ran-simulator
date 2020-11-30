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

package client

import (
	"context"
	"errors"
	"fmt"
	databaseapi "github.com/atomix/api/proto/atomix/database"
	"github.com/atomix/go-client/pkg/client/peer"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"github.com/atomix/go-client/pkg/client/util/net"
	"google.golang.org/grpc"
	"sort"
	"time"
)

// New creates a new Atomix client
func New(address string, opts ...Option) (*Client, error) {
	ctx := context.Background()
	options := applyOptions(opts...)
	if options.joinTimeout != nil {
		c, cancel := context.WithTimeout(context.Background(), *options.joinTimeout)
		defer cancel()
		ctx = c
	}
	return NewWithContext(ctx, address, opts...)
}

// NewWithContext returns a new Atomix client
func NewWithContext(ctx context.Context, address string, opts ...Option) (*Client, error) {
	options := applyOptions(opts...)

	// Set up a connection to the server.
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(util.RetryingUnaryClientInterceptor()), grpc.WithStreamInterceptor(util.RetryingStreamClientInterceptor(time.Second)))
	if err != nil {
		return nil, err
	}

	clusterOpts := []peer.Option{
		peer.WithNamespace(options.namespace),
		peer.WithScope(options.scope),
		peer.WithMemberID(options.memberID),
		peer.WithHost(options.peerHost),
		peer.WithPort(options.peerPort),
		peer.WithServices(options.peerServices...),
		peer.WithServerOptions(options.peerServerOpts...),
	}
	if options.joinTimeout != nil {
		clusterOpts = append(clusterOpts, peer.WithJoinTimeout(*options.joinTimeout))
	}

	peers, err := peer.NewGroupWithContext(ctx, address, clusterOpts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		peers:   peers,
		options: *options,
	}, nil
}

// Client is an Atomix client
type Client struct {
	conn    *grpc.ClientConn
	peers   *peer.Group
	options options
}

// Group returns the peer group
func (c *Client) Group() *peer.Group {
	return c.peers
}

// GetDatabases returns a list of all databases in the client's namespace
func (c *Client) GetDatabases(ctx context.Context) ([]*Database, error) {
	client := databaseapi.NewDatabaseServiceClient(c.conn)
	request := &databaseapi.GetDatabasesRequest{
		Namespace: c.options.namespace,
	}

	response, err := client.GetDatabases(ctx, request)
	if err != nil {
		return nil, err
	}

	databases := make([]*Database, len(response.Databases))
	for i, databaseProto := range response.Databases {
		database, err := c.newDatabase(ctx, &databaseProto)
		if err != nil {
			return nil, err
		}
		databases[i] = database
	}
	return databases, nil
}

// GetDatabase gets a database client by name from the client's namespace
func (c *Client) GetDatabase(ctx context.Context, name string) (*Database, error) {
	client := databaseapi.NewDatabaseServiceClient(c.conn)
	request := &databaseapi.GetDatabaseRequest{
		ID: databaseapi.DatabaseId{
			Name:      name,
			Namespace: c.options.namespace,
		},
	}

	response, err := client.GetDatabase(ctx, request)
	if err != nil {
		return nil, err
	} else if response.Database == nil {
		return nil, errors.New("unknown database " + name)
	}
	return c.newDatabase(ctx, response.Database)
}

func (c *Client) newDatabase(ctx context.Context, databaseProto *databaseapi.Database) (*Database, error) {
	// Ensure the partitions are sorted in case the controller sent them out of order.
	partitionProtos := databaseProto.Partitions
	sort.Slice(partitionProtos, func(i, j int) bool {
		return partitionProtos[i].PartitionID.Partition < partitionProtos[j].PartitionID.Partition
	})

	// Iterate through the partitions and create gRPC client connections for each partition.
	partitions := make([]primitive.Partition, len(databaseProto.Partitions))
	for i, partitionProto := range partitionProtos {
		ep := partitionProto.Endpoints[0]
		partitions[i] = primitive.Partition{
			ID:      int(partitionProto.PartitionID.Partition),
			Address: net.Address(fmt.Sprintf("%s:%d", ep.Host, ep.Port)),
		}
	}

	// Iterate through partitions and open sessions
	sessions := make([]*primitive.Session, len(partitions))
	for i, partition := range partitions {
		session, err := primitive.NewSession(ctx, partition, primitive.WithSessionTimeout(c.options.sessionTimeout))
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return &Database{
		Namespace: databaseProto.ID.Namespace,
		Name:      databaseProto.ID.Name,
		scope:     c.options.scope,
		sessions:  sessions,
		conn:      c.conn,
	}, nil
}

// Close closes the client
func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}
