// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"sync"

	"github.com/onosproject/ran-simulator/pkg/servicemodel"

	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
)

// ServiceModelRegistry stores list of registered service models
type ServiceModelRegistry struct {
	mu            sync.Mutex
	serviceModels map[ID]servicemodel.ServiceModel
	ranFunctions  types.RanFunctions
}

type ServiceModel struct {
	ID           ID
	ServiceModel servicemodel.ServiceModel
	Description  string
	Revision     int
}

// NewServiceModelRegistry creates a service model registry
func NewServiceModelRegistry() *ServiceModelRegistry {
	return &ServiceModelRegistry{
		serviceModels: make(map[ID]servicemodel.ServiceModel),
	}
}

// RegisterServiceModel registers a service model
func (s *ServiceModelRegistry) RegisterServiceModel(sm ServiceModel) error {
	if _, exists := s.serviceModels[sm.ID]; exists {
		return errors.New(errors.AlreadyExists, "the service model already registered")
	}

	/*ranFuncID := types.RanFunctionID(sm.ID)
	s.ranFunctions[ranFuncID] = types.RanFunctionItem{
		Description: types.RanFunctionDescription(sm.Description),
		Revision:    types.RanFunctionRevision(sm.Revision),
	}*/

	s.mu.Lock()
	s.serviceModels[sm.ID] = sm.ServiceModel
	s.mu.Unlock()
	return nil
}

// GetServiceModel finds and initialize service model interface pointer
func (s *ServiceModelRegistry) GetServiceModel(id ID, sm interface{}) error {
	if serviceModel, ok := s.serviceModels[id]; ok {
		sm = serviceModel
		return nil
	}

	return errors.New(errors.Unknown, "no service model implementation exists for ran function ID:", id)
}

// GetRanFunctions returns the list of registered ran functions
func (s *ServiceModelRegistry) GetRanFunctions() types.RanFunctions {
	return s.ranFunctions
}
