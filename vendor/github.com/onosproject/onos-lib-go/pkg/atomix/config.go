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
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/env"
)

const defaultPort = 5678

// DatabaseType is the type of a database
type DatabaseType string

const (
	// DatabaseTypeConsensus indicates a consensus database
	DatabaseTypeConsensus DatabaseType = "consensus"
	// DatabaseTypeCache indicates a cache database
	DatabaseTypeCache DatabaseType = "cache"
	// DatabaseTypeEvent indicates an event database
	DatabaseTypeEvent DatabaseType = "event"

	// DatabaseTypeConfig indicates a configuration database
	// Deprecated: Use DatabaseTypeConsensus instead
	DatabaseTypeConfig DatabaseType = "config"
	// DatabaseTypeTimeSeries indicates a time series database
	// Deprecated: Use DatabaseTypeEvent instead
	DatabaseTypeTimeSeries DatabaseType = "timeSeries"
	// DatabaseTypeRelational indicates a relational database
	// Deprecated: Use DatabaseTypeConfig or DatabaseTypeCache instead
	DatabaseTypeRelational DatabaseType = "relational"
)

var databaseTypeCompat = map[DatabaseType]DatabaseType{
	DatabaseTypeConfig:     DatabaseTypeConsensus,
	DatabaseTypeEvent:      DatabaseTypeTimeSeries,
	DatabaseTypeConsensus:  DatabaseTypeConfig,
	DatabaseTypeTimeSeries: DatabaseTypeEvent,
}

// Config is the Atomix configuration
type Config struct {
	// Controller is the Atomix controller address
	Controller string `yaml:"controller,omitempty"`
	// Member is the Atomix member name
	Member string `yaml:"member,omitempty"`
	// Host is the Atomix member hostname
	Host string `yaml:"host,omitempty"`
	// Port is the Atomix member port
	Port int `yaml:"port,omitempty"`
	// Namespace is the Atomix namespace
	Namespace string `yaml:"namespace,omitempty"`
	// Scope is the Atomix client/application scope
	Scope string `yaml:"scope,omitempty"`
	// Databases is a mapping of database types to databases
	Databases map[DatabaseType]string `yaml:"databases"`
}

// GetController gets the Atomix controller address
func (c Config) GetController() string {
	if c.Controller == "" {
		namespace := c.GetNamespace()
		if namespace != "" {
			c.Controller = fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", namespace)
		}
	}
	return c.Controller
}

// GetMember gets the Atomix member name
func (c Config) GetMember() string {
	if c.Member == "" {
		c.Member = env.GetPodName()
	}
	return c.Member
}

// GetHost gets the Atomix peer host
func (c Config) GetHost() string {
	if c.Host == "" {
		c.Host = env.GetPodID()
	}
	return c.Host
}

// GetPort gets the Atomix peer port
func (c Config) GetPort() int {
	if c.Port == 0 {
		c.Port = defaultPort
	}
	return c.Port
}

// GetNamespace gets the Atomix client namespace
func (c Config) GetNamespace() string {
	if c.Namespace == "" {
		c.Namespace = env.GetServiceNamespace()
	}
	return c.Namespace
}

// GetScope gets the Atomix client scope
func (c Config) GetScope() string {
	if c.Scope == "" {
		c.Scope = env.GetServiceName()
	}
	if c.Scope == "" {
		c.Scope = c.GetNamespace()
	}
	return c.Scope
}

// GetDatabase gets the database name for the given database type
func (c Config) GetDatabase(databaseType DatabaseType) string {
	db, ok := c.Databases[databaseType]
	if ok {
		return db
	}
	dbType, ok := databaseTypeCompat[databaseType]
	if ok {
		return c.Databases[dbType]
	}
	return ""
}
