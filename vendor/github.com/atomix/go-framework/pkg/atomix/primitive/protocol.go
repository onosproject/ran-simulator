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
	"github.com/atomix/go-framework/pkg/atomix/cluster"
)

// ProtocolContext provides the current state of the protocol
type ProtocolContext interface {
	// NodeID is the local node identifier
	NodeID() string
}

// Protocol is the interface to be implemented by replication protocols
type Protocol interface {
	// Partition returns a partition
	Partition(partitionID PartitionID) Partition

	// Partitions returns the protocol partitions
	Partitions() []Partition

	// Start starts the protocol
	Start(cluster cluster.Cluster, registry Registry) error

	// Stop stops the protocol
	Stop() error
}
