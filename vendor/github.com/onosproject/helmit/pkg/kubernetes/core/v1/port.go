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

package v1

import (
	"fmt"
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	corev1 "k8s.io/api/core/v1"
)

// ServicePort is a service port
type ServicePort struct {
	resource.Client
	corev1.ServicePort
	service *Service
}

// Address returns the address of the port
func (p *ServicePort) Address(qualified bool) string {
	return fmt.Sprintf("%s:%d", p.service.Hostname(qualified), p.Port)
}

// Ports returns a list of service ports
func (s *Service) Ports() []*ServicePort {
	ports := make([]*ServicePort, len(s.Object.Spec.Ports))
	for i, port := range s.Object.Spec.Ports {
		ports[i] = &ServicePort{
			Client:      s.Resource.Client,
			ServicePort: port,
			service:     s,
		}
	}
	return ports
}

// Port returns a service port by name
func (s *Service) Port(name string) *ServicePort {
	for _, port := range s.Object.Spec.Ports {
		if port.Name == name {
			return &ServicePort{
				Client:      s.Resource.Client,
				ServicePort: port,
				service:     s,
			}
		}
	}
	return nil
}
