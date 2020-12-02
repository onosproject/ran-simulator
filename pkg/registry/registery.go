// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"fmt"
	"sync"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/servicemodels"
)

// ServiceModelRegistry stores list of registered service models
type ServiceModelRegistry struct {
	mu            sync.Mutex
	serviceModels map[ID]servicemodels.ServiceModel
}

// NewServiceModelRegistry creates a service model registry
func NewServiceModelRegistry() *ServiceModelRegistry {
	return &ServiceModelRegistry{
		serviceModels: make(map[ID]servicemodels.ServiceModel),
	}
}

// RegisterServiceModel registers a service model
func (s *ServiceModelRegistry) RegisterServiceModel(id ID, sm servicemodels.ServiceModel) error {
	if _, exists := s.serviceModels[id]; exists {
		return errors.New(errors.AlreadyExists, "the service model already registered")
	}

	s.mu.Lock()
	s.serviceModels[id] = sm
	s.mu.Unlock()
	return nil
}

// GetServiceModel finds and initialize service model interface pointer
func (s *ServiceModelRegistry) GetServiceModel(id ID, sm interface{}) error {
	if _, ok := s.serviceModels[id]; ok {
		return nil
	}
	return fmt.Errorf("unknown service: %T", sm)
}
