// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package modelplugins

import (
	"fmt"
	"plugin"
	"sync"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	types "github.com/onosproject/onos-api/go/onos/e2t/e2sm"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("modelregistry")

// ModelFullName service model name
type ModelFullName string

// ModelVersion service model version
type ModelVersion string

// ModelOid service model OID
type ModelOid string

// ModelRegistry is the object for the saving information about device models
type ModelRegistry interface {
	GetPlugins() map[ModelOid]ServiceModel
	GetPlugin(oid ModelOid) (ServiceModel, error)
	RegisterModelPlugin(moduleName string) (ModelFullName, ModelVersion, error)
}

type modelRegistry struct {
	plugins map[ModelOid]ServiceModel
	mu      sync.RWMutex
}

// NewModelRegistry create an instance of model registry
func NewModelRegistry() ModelRegistry {
	return &modelRegistry{
		plugins: make(map[ModelOid]ServiceModel),
	}
}

// ServiceModel is a set of methods that each model plugin should implement
type ServiceModel interface {
	ServiceModelData2() (string, string, string, string)
	ServiceModelData() (string, string, string)
	IndicationHeaderASN1toProto(asn1Bytes []byte) ([]byte, error)
	IndicationHeaderProtoToASN1(protoBytes []byte) ([]byte, error)
	IndicationMessageASN1toProto(asn1Bytes []byte) ([]byte, error)
	IndicationMessageProtoToASN1(protoBytes []byte) ([]byte, error)
	RanFuncDescriptionASN1toProto(asn1Bytes []byte) ([]byte, error)
	RanFuncDescriptionProtoToASN1(protoBytes []byte) ([]byte, error)
	EventTriggerDefinitionASN1toProto(asn1Bytes []byte) ([]byte, error)
	EventTriggerDefinitionProtoToASN1(protoBytes []byte) ([]byte, error)
	ActionDefinitionASN1toProto(asn1Bytes []byte) ([]byte, error)
	ActionDefinitionProtoToASN1(protoBytes []byte) ([]byte, error)
	DecodeRanFunctionDescription(asn1bytes []byte) (*types.RanfunctionNameDef, *types.RicEventTriggerList, *types.RicReportList, error)
	ControlHeaderASN1toProto(asn1Bytes []byte) ([]byte, error)
	ControlHeaderProtoToASN1(protoBytes []byte) ([]byte, error)
	ControlMessageASN1toProto(asn1Bytes []byte) ([]byte, error)
	ControlMessageProtoToASN1(protoBytes []byte) ([]byte, error)
	ControlOutcomeASN1toProto(asn1Bytes []byte) ([]byte, error)
	ControlOutcomeProtoToASN1(protoBytes []byte) ([]byte, error)
}

// GetModelPlugins get model plugins
func (r *modelRegistry) GetPlugins() map[ModelOid]ServiceModel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	plugins := make(map[ModelOid]ServiceModel, len(r.plugins))
	for id, plugin := range r.plugins {
		plugins[id] = plugin
	}
	return plugins
}

// GetPlugin returns the model plugin interface
func (r *modelRegistry) GetPlugin(oid ModelOid) (ServiceModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	serviceModel, ok := r.plugins[oid]
	if !ok {
		err := errors.NewNotFound("Model plugin '%s' not found", oid)
		return nil, err
	}
	return serviceModel, nil

}

// RegisterModelPlugin adds an external model plugin to the model registry at startup
// or through the 'admin' gRPC interface. Once plugins are loaded they cannot be unloaded
func (r *modelRegistry) RegisterModelPlugin(moduleName string) (ModelFullName, ModelVersion, error) {
	log.Info("Loading module ", moduleName)
	modelPluginModule, err := plugin.Open(moduleName)
	if err != nil {
		log.Warnf("Unable to load module %s %s", moduleName, err)
		return "", "", err
	}
	symbolMP, err := modelPluginModule.Lookup("ServiceModel")
	if err != nil {
		log.Warn("Unable to find ServiceModel in module ", moduleName, err)
		return "", "", err
	}
	serviceModelPlugin, ok := symbolMP.(ServiceModel)
	if !ok {
		log.Warnf("Unable to use ServiceModelPlugin in %s", moduleName)
		return "", "", fmt.Errorf("symbol loaded from module %s is not a ServiceModel",
			moduleName)
	}
	name, version, _, oid := serviceModelPlugin.ServiceModelData2()
	modelOid := ModelOid(oid)
	log.Infof("Loaded %s %s from %s", name, version, moduleName)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[modelOid] = serviceModelPlugin

	return ModelFullName(name), ModelVersion(version), nil
}
