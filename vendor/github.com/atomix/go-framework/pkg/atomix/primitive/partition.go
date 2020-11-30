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

package primitive

import (
	"context"
	"github.com/atomix/go-framework/pkg/atomix/stream"
	"time"
)

// Index is a partition log index
type Index uint64

// PartitionContext provides the current state of the partition
type PartitionContext interface {
	ProtocolContext

	// PartitionID is the partition identifier
	PartitionID() PartitionID

	// Index returns the current index of the service
	Index() Index

	// Timestamp returns a deterministic, monotonically increasing timestamp
	Timestamp() time.Time
}

// PartitionID is a partition identifier
type PartitionID int

// Partition is the interface for a partition client
type Partition interface {
	// MustLeader returns whether the client can only be used on the leader
	MustLeader() bool

	// IsLeader returns whether the client is the leader
	IsLeader() bool

	// Leader returns the current leader
	Leader() string

	// Write sends a write request
	Write(ctx context.Context, input []byte, stream stream.WriteStream) error

	// Read sends a read request
	Read(ctx context.Context, input []byte, stream stream.WriteStream) error
}
