// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

import (
	"sync"

	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
)

var log = logging.GetLogger("registry")

// ServiceModelRegistry stores list of registered service models
type ServiceModelRegistry struct {
	mu            sync.RWMutex
	serviceModels map[ID]servicemodel.ServiceModel
	ranFunctions  types.RanFunctions
}

// ServiceModelConfig service model configuration
type ServiceModelConfig struct {
	ID           ID
	ServiceModel servicemodel.ServiceModel
	Description  []byte // ASN1 bytes from Service Model
	Revision     int
}

// NewServiceModelRegistry creates a service model registry
func NewServiceModelRegistry() *ServiceModelRegistry {
	return &ServiceModelRegistry{
		serviceModels: make(map[ID]servicemodel.ServiceModel),
		ranFunctions:  make(map[types.RanFunctionID]types.RanFunctionItem),
	}
}

// RegisterServiceModel registers a service model
func (s *ServiceModelRegistry) RegisterServiceModel(sm ServiceModelConfig) error {
	log.Info("Register Service Model:", sm)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.serviceModels[sm.ID]; exists {
		return errors.New(errors.AlreadyExists, "the service model already registered")
	}
	ranFuncID := types.RanFunctionID(sm.ID)
	s.ranFunctions[ranFuncID] = types.RanFunctionItem{
		Description: sm.Description,
		Revision:    types.RanFunctionRevision(sm.Revision),
	}
	s.serviceModels[sm.ID] = sm.ServiceModel
	return nil
}

// GetServiceModel finds and initialize service model interface pointer
func (s *ServiceModelRegistry) GetServiceModel(id ID) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sm, ok := s.serviceModels[id]
	if ok {
		return sm, nil
	}
	return nil, errors.New(errors.Unknown, "no service model implementation exists for ran function ID:", id)
}

// GetRanFunctions returns the list of registered ran functions
func (s *ServiceModelRegistry) GetRanFunctions() types.RanFunctions {
	return s.ranFunctions
}
