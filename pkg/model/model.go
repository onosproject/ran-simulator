// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("pkg", "model")

type Ecgi string

type Model struct {
	Nodes         map[string]Node         `yaml:"nodes"`
	Controllers   map[string]Controller   `yaml:"controllers"`
	ServiceModels map[string]ServiceModel `yaml:"servicemodels"`
}

type Node struct {
	Ecgi          Ecgi     `yaml:"ecgi"`
	Controllers   []string `yaml:"controllers"`
	ServiceModels []string `yaml:"servicemodels"`
}

type Controller struct {
	ID      string `yaml:"id"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type ServiceModel struct {
	ID          int    `yaml:"id"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

func (m *Model) GetNode(name string) Node {
	if node, ok := m.Nodes[name]; ok {
		return node
	}

	return Node{}
}

func (m *Model) GetServiceModel(name string) ServiceModel {
	if sm, ok := m.ServiceModels[name]; ok {
		return sm
	}

	return ServiceModel{}
}

// GetAllNodes gets all of the simulated nodes
func (m *Model) GetNodes() []Node {
	var nodes []Node
	for _, node := range m.Nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// GetController gets a controller by a given name
func (m *Model) GetController(name string) Controller {
	if controller, ok := m.Controllers[name]; ok {
		return controller
	}

	return Controller{}
}

// GetAllControllers gets all of the controllers
func (m *Model) GetControllers() []Controller {
	var controllers []Controller
	for _, controller := range m.Controllers {
		controllers = append(controllers, controller)
	}

	return controllers
}
