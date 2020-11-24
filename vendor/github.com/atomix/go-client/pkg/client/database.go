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

package client

import (
	"context"
	databaseapi "github.com/atomix/api/proto/atomix/database"
	primitiveapi "github.com/atomix/api/proto/atomix/primitive"
	"github.com/atomix/go-client/pkg/client/counter"
	"github.com/atomix/go-client/pkg/client/election"
	"github.com/atomix/go-client/pkg/client/indexedmap"
	"github.com/atomix/go-client/pkg/client/leader"
	"github.com/atomix/go-client/pkg/client/list"
	"github.com/atomix/go-client/pkg/client/lock"
	"github.com/atomix/go-client/pkg/client/log"
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/set"
	"github.com/atomix/go-client/pkg/client/value"
	"google.golang.org/grpc"
)

// Database manages the primitives in a set of partitions
type Database struct {
	Namespace string
	Name      string

	scope    string
	conn     *grpc.ClientConn
	sessions []*primitive.Session
}

// GetPrimitives gets a list of primitives in the database
func (d *Database) GetPrimitives(ctx context.Context, opts ...primitive.MetadataOption) ([]primitive.Metadata, error) {
	client := primitiveapi.NewPrimitiveServiceClient(d.conn)

	request := &primitiveapi.GetPrimitivesRequest{
		Database: &databaseapi.DatabaseId{
			Namespace: d.Namespace,
			Name:      d.Name,
		},
		Primitive: &primitiveapi.PrimitiveId{
			Namespace: d.scope,
		},
	}

	response, err := client.GetPrimitives(ctx, request)
	if err != nil {
		return nil, err
	}

	primitives := make([]primitive.Metadata, len(response.Primitives))
	for i, p := range response.Primitives {
		var primitiveType primitive.Type
		switch p.Type {
		case primitiveapi.PrimitiveType_COUNTER:
			primitiveType = "Counter"
		case primitiveapi.PrimitiveType_ELECTION:
			primitiveType = "Election"
		case primitiveapi.PrimitiveType_INDEXED_MAP:
			primitiveType = "IndexedMap"
		case primitiveapi.PrimitiveType_LEADER_LATCH:
			primitiveType = "LeaderLatch"
		case primitiveapi.PrimitiveType_LIST:
			primitiveType = "List"
		case primitiveapi.PrimitiveType_LOCK:
			primitiveType = "Lock"
		case primitiveapi.PrimitiveType_LOG:
			primitiveType = "Log"
		case primitiveapi.PrimitiveType_MAP:
			primitiveType = "Map"
		case primitiveapi.PrimitiveType_SET:
			primitiveType = "Set"
		case primitiveapi.PrimitiveType_VALUE:
			primitiveType = "Value"
		default:
			primitiveType = "Unknown"
		}
		primitives[i] = primitive.Metadata{
			Type: primitiveType,
			Name: primitive.Name{
				Scope: p.Primitive.Namespace,
				Name:  p.Primitive.Name,
			},
		}
	}
	return primitives, nil
}

// GetCounter gets or creates a Counter with the given name
func (d *Database) GetCounter(ctx context.Context, name string) (counter.Counter, error) {
	return counter.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions)
}

// GetElection gets or creates an Election with the given name
func (d *Database) GetElection(ctx context.Context, name string, opts ...election.Option) (election.Election, error) {
	return election.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions, opts...)
}

// GetIndexedMap gets or creates a Map with the given name
func (d *Database) GetIndexedMap(ctx context.Context, name string) (indexedmap.IndexedMap, error) {
	return indexedmap.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions)
}

// GetLeaderLatch gets or creates a LeaderLatch with the given name
func (d *Database) GetLeaderLatch(ctx context.Context, name string, opts ...leader.Option) (leader.Latch, error) {
	return leader.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions, opts...)
}

// GetList gets or creates a List with the given name
func (d *Database) GetList(ctx context.Context, name string) (list.List, error) {
	return list.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions)
}

// GetLock gets or creates a Lock with the given name
func (d *Database) GetLock(ctx context.Context, name string) (lock.Lock, error) {
	return lock.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions)
}

// GetLog gets or creates a Log with the given name
func (d *Database) GetLog(ctx context.Context, name string) (log.Log, error) {
	return log.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions)
}

// GetMap gets or creates a Map with the given name
func (d *Database) GetMap(ctx context.Context, name string, opts ..._map.Option) (_map.Map, error) {
	return _map.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions, opts...)
}

// GetSet gets or creates a Set with the given name
func (d *Database) GetSet(ctx context.Context, name string) (set.Set, error) {
	return set.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions)
}

// GetValue gets or creates a Value with the given name
func (d *Database) GetValue(ctx context.Context, name string) (value.Value, error) {
	return value.New(ctx, primitive.NewName(d.Namespace, d.Name, d.scope, name), d.sessions)
}
