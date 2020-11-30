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

package lock

import (
	"context"
	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/lock"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"google.golang.org/grpc"
)

// Type is the lock type
const Type primitive.Type = "Lock"

// Client provides an API for creating Locks
type Client interface {
	// GetLock gets the Lock instance of the given name
	GetLock(ctx context.Context, name string) (Lock, error)
}

// Lock provides distributed concurrency control
type Lock interface {
	primitive.Primitive

	// Lock acquires the lock
	Lock(ctx context.Context, opts ...LockOption) (uint64, error)

	// Unlock releases the lock
	Unlock(ctx context.Context, opts ...UnlockOption) (bool, error)

	// IsLocked returns a bool indicating whether the lock is held
	IsLocked(ctx context.Context, opts ...IsLockedOption) (bool, error)
}

// New creates a new Lock primitive for the given partitions
// The lock will be created in one of the given partitions.
func New(ctx context.Context, name primitive.Name, partitions []*primitive.Session) (Lock, error) {
	i, err := util.GetPartitionIndex(name.Name, len(partitions))
	if err != nil {
		return nil, err
	}
	return newLock(ctx, name, partitions[i])
}

// newLock creates a new Lock primitive for the given partition
func newLock(ctx context.Context, name primitive.Name, partition *primitive.Session) (*lock, error) {
	instance, err := primitive.NewInstance(ctx, name, partition, &primitiveHandler{})
	if err != nil {
		return nil, err
	}
	return &lock{
		name:     name,
		instance: instance,
	}, nil
}

// lock is the single partition implementation of Lock
type lock struct {
	name     primitive.Name
	instance *primitive.Instance
}

func (l *lock) Name() primitive.Name {
	return l.name
}

func (l *lock) Lock(ctx context.Context, opts ...LockOption) (uint64, error) {
	response, err := l.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLockServiceClient(conn)
		request := &api.LockRequest{
			Header: header,
		}
		for i := range opts {
			opts[i].beforeLock(request)
		}
		response, err := client.Lock(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterLock(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return 0, err
	}
	return response.(*api.LockResponse).Version, nil
}

func (l *lock) Unlock(ctx context.Context, opts ...UnlockOption) (bool, error) {
	response, err := l.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLockServiceClient(conn)
		request := &api.UnlockRequest{
			Header: header,
		}
		for i := range opts {
			opts[i].beforeUnlock(request)
		}
		response, err := client.Unlock(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterUnlock(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return false, err
	}
	return response.(*api.UnlockResponse).Unlocked, nil
}

func (l *lock) IsLocked(ctx context.Context, opts ...IsLockedOption) (bool, error) {
	response, err := l.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLockServiceClient(conn)
		request := &api.IsLockedRequest{
			Header: header,
		}
		for i := range opts {
			opts[i].beforeIsLocked(request)
		}
		response, err := client.IsLocked(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		for i := range opts {
			opts[i].afterIsLocked(response)
		}
		return response.Header, response, nil
	})
	if err != nil {
		return false, err
	}
	return response.(*api.IsLockedResponse).IsLocked, nil
}

func (l *lock) Close(ctx context.Context) error {
	return l.instance.Close(ctx)
}

func (l *lock) Delete(ctx context.Context) error {
	return l.instance.Delete(ctx)
}
