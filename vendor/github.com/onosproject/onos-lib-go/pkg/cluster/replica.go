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
)

// newReplica creates a new replica
func newReplica(id ReplicaID, local bool, connector func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error)) *Replica {
	return &Replica{
		Node:      newNode(NodeID(id)),
		ID:        id,
		local:     local,
		connector: connector,
	}
}

// ReplicaID is a replica identifier
type ReplicaID NodeID

// Replica is a cluster replica
type Replica struct {
	Node
	ID        ReplicaID
	local     bool
	connector func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

// IsLocal returns a bool indicating whether the replica is local
func (r *Replica) IsLocal() bool {
	return r.local
}

// Connect connects to the replica
func (r *Replica) Connect(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	options := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	options = append(options, opts...)
	return r.connector(ctx, options...)
}

// ReplicaSet is a set of replicas
type ReplicaSet map[ReplicaID]*Replica
