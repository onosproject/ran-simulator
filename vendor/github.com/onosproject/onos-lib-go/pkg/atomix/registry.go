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

package atomix

import (
	"github.com/onosproject/onos-lib-go/pkg/cluster"
)

var serviceRegistry = &ServiceRegistry{
	services: make([]cluster.Service, 0),
}

// RegisterService registers a service with the Atomix cluster
func RegisterService(service cluster.Service) {
	serviceRegistry.Register(service)
}

// ServiceRegistry is a registry of Atomix services
type ServiceRegistry struct {
	services []cluster.Service
}

// Register registers a service
func (r *ServiceRegistry) Register(service cluster.Service) {
	r.services = append(r.services, service)
}
