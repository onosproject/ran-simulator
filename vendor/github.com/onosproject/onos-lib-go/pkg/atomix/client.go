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

package atomix

import (
	"context"
	"fmt"
	"github.com/atomix/go-client/pkg/client"
	"github.com/atomix/go-client/pkg/client/peer"
	netutil "github.com/atomix/go-client/pkg/client/util/net"
	"github.com/atomix/go-framework/pkg/atomix"
	"github.com/atomix/go-framework/pkg/atomix/counter"
	"github.com/atomix/go-framework/pkg/atomix/election"
	"github.com/atomix/go-framework/pkg/atomix/indexedmap"
	"github.com/atomix/go-framework/pkg/atomix/leader"
	"github.com/atomix/go-framework/pkg/atomix/list"
	"github.com/atomix/go-framework/pkg/atomix/lock"
	"github.com/atomix/go-framework/pkg/atomix/log"
	"github.com/atomix/go-framework/pkg/atomix/map"
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/set"
	"github.com/atomix/go-framework/pkg/atomix/value"
	"github.com/atomix/go-local/pkg/atomix/local"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"google.golang.org/grpc"
	"net"
)

const basePort = 45000

// StartLocalNode starts a single local Atomix node
func StartLocalNode() (*atomix.Node, netutil.Address) {
	for port := basePort; port < basePort+100; port++ {
		address := netutil.Address(fmt.Sprintf("localhost:%d", port))
		lis, err := net.Listen("tcp", string(address))
		if err != nil {
			continue
		}
		node := local.NewNode(lis, []primitive.PartitionID{1})
		counter.RegisterPrimitive(node)
		election.RegisterPrimitive(node)
		indexedmap.RegisterPrimitive(node)
		lock.RegisterPrimitive(node)
		log.RegisterPrimitive(node)
		leader.RegisterPrimitive(node)
		list.RegisterPrimitive(node)
		_map.RegisterPrimitive(node)
		set.RegisterPrimitive(node)
		value.RegisterPrimitive(node)
		_ = node.Start()
		return node, address
	}
	panic("cannot find open port")
}

// GetClient returns the Atomix client
func GetClient(config Config) (*client.Client, error) {
	opts := []client.Option{
		client.WithNamespace(config.GetNamespace()),
		client.WithScope(config.GetScope()),
	}
	member := config.GetMember()
	host := config.GetHost()
	if host != "" {
		opts = append(opts, client.WithPeerHost(config.GetHost()))
		opts = append(opts, client.WithPeerPort(config.GetPort()))
		for _, s := range serviceRegistry.services {
			service := func(service cluster.Service) func(peer.ID, *grpc.Server) {
				return func(id peer.ID, server *grpc.Server) {
					service(cluster.NodeID(id), server)
				}
			}(s)
			opts = append(opts, client.WithPeerService(service))
		}
	}
	if member != "" {
		opts = append(opts, client.WithMemberID(config.GetMember()))
	} else if host != "" {
		opts = append(opts, client.WithMemberID(config.GetHost()))
	}

	return client.New(config.GetController(), opts...)
}

// GetDatabase returns the Atomix database
func GetDatabase(config Config, database string) (*client.Database, error) {
	client, err := GetClient(config)
	if err != nil {
		return nil, err
	}
	return client.GetDatabase(context.Background(), database)
}
