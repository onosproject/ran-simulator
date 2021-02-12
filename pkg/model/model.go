// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/onosproject/onos-lib-go/pkg/errors"
)

// Model simulation model
type Model struct {
	MapLayout     MapLayout               `mapstructure:"layout"`
	Nodes         map[string]Node         `mapstructure:"nodes"`
	Controllers   map[string]Controller   `mapstructure:"controllers"`
	ServiceModels map[string]ServiceModel `mapstructure:"servicemodels"`
	UECount       uint                    `mapstructure:"ueCount"`
	PlmnID        PlmnID                  `mapstructure:"plmnID"`
	UEs           UERegistry              // Not intended to be loaded from the YAML file; created separately
	// Routes   *SimRoutes
}

// Node e2 node
type Node struct {
	EnbID         EnbID           `mapstructure:"enbID"`
	Controllers   []string        `mapstructure:"controllers"`
	ServiceModels []string        `mapstructure:"servicemodels"`
	Cells         map[string]Cell `mapstructure:"cells"`
}

// Controller E2T endpoint information
type Controller struct {
	ID      string `mapstructure:"id"`
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

// Cell represents a section of coverage
type Cell struct {
	Ecgi      ECGI    `mapstructure:"ecgi"`
	Sector    Sector  `mapstructure:"sector"`
	Color     string  `mapstructure:"color"`
	MaxUEs    uint32  `mapstructure:"maxUEs"`
	Neighbors []ECGI  `mapstructure:"neighbors"`
	TxPowerDB float64 `mapstructure:"txPower"`

	// TODO: should not be needed as it coincides with sector center.
	//Location  Coordinate `mapstructure:"location"`

	// TODO: add the following later or track them differently
	//Crntis map
	//CrntiIndex uint32     `mapstructure:"crntiIndex"`
	//Port       uint32     `mapstructure:"port"`
}

// ServiceModel service model information
type ServiceModel struct {
	ID          int    `mapstructure:"id"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
}

// GetNode gets a an e2 node
func (m *Model) GetNode(name string) (Node, error) {
	if node, ok := m.Nodes[name]; ok {
		return node, nil
	}

	return Node{}, errors.New(errors.NotFound, "node not found")
}

// GetServiceModel gets a service model  based on a given name
func (m *Model) GetServiceModel(name string) (ServiceModel, error) {
	if sm, ok := m.ServiceModels[name]; ok {
		return sm, nil
	}

	return ServiceModel{}, errors.New(errors.NotFound, "the service model not found")
}

// GetNodes gets all of the simulated nodes
func (m *Model) GetNodes() []Node {
	nodes := make([]Node, 0, len(m.Nodes))
	for _, node := range m.Nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// GetController gets a controller by a given name
func (m *Model) GetController(name string) (Controller, error) {
	if controller, ok := m.Controllers[name]; ok {
		return controller, nil
	}

	return Controller{}, errors.New(errors.NotFound, "controller not found")
}

// GetControllers gets all of the controllers
func (m *Model) GetControllers() []Controller {
	var controllers []Controller
	for _, controller := range m.Controllers {
		controllers = append(controllers, controller)
	}

	return controllers
}
