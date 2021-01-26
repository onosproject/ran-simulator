// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"sync"

	"github.com/onosproject/ran-simulator/pkg/modelplugins"

	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
)

var log = logging.GetLogger("registry")

// ServiceModelRegistry stores list of registered service models
type ServiceModelRegistry struct {
	mu            sync.RWMutex
	serviceModels map[RanFunctionID]ServiceModel
	ranFunctions  e2aptypes.RanFunctions
}

// ServiceModel service model
type ServiceModel struct {
	RanFunctionID       RanFunctionID
	ModelName           string
	Version             string
	Description         []byte // ASN1 bytes from Service Model
	Revision            int
	Client              servicemodel.Client
	ModelPluginRegistry *modelplugins.ModelPluginRegistry
}

// NewServiceModelRegistry creates a service model registry
func NewServiceModelRegistry() *ServiceModelRegistry {
	return &ServiceModelRegistry{
		serviceModels: make(map[RanFunctionID]ServiceModel),
		ranFunctions:  make(map[e2aptypes.RanFunctionID]e2aptypes.RanFunctionItem),
	}
}

// RegisterServiceModel registers a service model
func (s *ServiceModelRegistry) RegisterServiceModel(sm ServiceModel) error {
	log.Info("Register Service Model:", sm.ModelName)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.serviceModels[sm.RanFunctionID]; exists {
		return errors.New(errors.AlreadyExists, "the service model already registered")
	}
	ranFuncID := e2aptypes.RanFunctionID(sm.RanFunctionID)
	s.ranFunctions[ranFuncID] = e2aptypes.RanFunctionItem{
		Description: sm.Description,
		Revision:    e2aptypes.RanFunctionRevision(sm.Revision),
	}
	s.serviceModels[sm.RanFunctionID] = sm

	return nil
}

// GetServiceModel finds and initialize service model interface pointer
func (s *ServiceModelRegistry) GetServiceModel(id RanFunctionID) (ServiceModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sm, ok := s.serviceModels[id]
	if ok {
		return sm, nil
	}
	return ServiceModel{}, errors.New(errors.Unknown, "no service model implementation exists for ran function ID:", id)
}

// GetRanFunctions returns the list of registered ran functions
func (s *ServiceModelRegistry) GetRanFunctions() e2aptypes.RanFunctions {
	return s.ranFunctions
}
