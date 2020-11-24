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
	"fmt"
	"github.com/atomix/go-client/pkg/client/util/net"
)

// Type is the type of a primitive
type Type string

// NewName returns a qualified primitive name with the given namespace, group, application, and name
func NewName(namespace string, group string, scope string, name string) Name {
	return Name{
		Namespace: namespace,
		Database:  group,
		Scope:     scope,
		Name:      name,
	}
}

// Name is a qualified primitive name consisting of Namespace, Database, Application, and Name
type Name struct {
	// Namespace is the namespace within which the database is stored
	Namespace string
	// Database is the database in which the primitive is stored
	Database string
	// Scope is the application scope in which the primitive is stored
	Scope string
	// Name is the simple name of the primitive
	Name string
}

func (n Name) String() string {
	return fmt.Sprintf("%s.%s.%s.%s", n.Namespace, n.Database, n.Scope, n.Name)
}

// Primitive is the base interface for primitives
type Primitive interface {
	// Name returns the fully namespaced primitive name
	Name() Name

	// Close closes the primitive
	Close(ctx context.Context) error

	// Delete deletes the primitive state from the cluster
	Delete(ctx context.Context) error
}

// Partition is the ID and address for a partition
type Partition struct {
	// ID is the partition identifier
	ID int

	// Address is the partition address
	Address net.Address
}

// Metadata is primitive metadata
type Metadata struct {
	// Type is the primitive type
	Type Type

	// Name is the primitive name
	Name Name
}
