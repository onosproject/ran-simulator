// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/api/types"
)

// Model simulation model
type Model struct {
	MapLayout     MapLayout               `yaml:"layout"`
	Nodes         map[string]Node         `yaml:"nodes"`
	Controllers   map[string]Controller   `yaml:"controllers"`
	ServiceModels map[string]ServiceModel `yaml:"servicemodels"`
	UECount       uint                    `yaml:"ueCount"`
	PlmnID        types.PlmnID            `yaml:"plmnID"`

	// Not intended to be loaded from the YAML file; created separately
	UEs UERegistry
}

// Coordinate represents a geographical location
type Coordinate struct {
	Lat float64 `yaml:"lat"`
	Lng float64 `yaml:"lng"`
}

// Sector represents a 2D arc emanating from a location
type Sector struct {
	Center  Coordinate `yaml:"center"`
	Azimuth int32      `yaml:"azimuth"`
	Arc     int32      `yaml:"arc"`
}

// Route represents a named series of points for tracking movement of user-equipment
type Route struct {
	Name   string
	Points []*Coordinate
	Color  string
}

// Node e2 node
type Node struct {
	EnbID         types.EnbID     `yaml:"enbID"`
	Controllers   []string        `yaml:"controllers"`
	ServiceModels []string        `yaml:"servicemodels"`
	Cells         map[string]Cell `yaml:"cells"`
}

// Controller E2T endpoint information
type Controller struct {
	ID      string `yaml:"id"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

// Cell represents a section of coverage
type Cell struct {
	Ecgi      types.ECGI   `yaml:"ecgi"`
	Sector    Sector       `yaml:"sector"`
	Color     string       `yaml:"color"`
	MaxUEs    uint32       `yaml:"maxUEs"`
	Neighbors []types.ECGI `yaml:"neighbors"`
	TxPowerDB float64      `yaml:"txPower"`

	// TODO: should not be needed as it coincides with sector center.
	//Location  Coordinate `yaml:"location"`

	// TODO: add the following later or track them differently
	//Crntis map
	//CrntiIndex uint32     `yaml:"crntiIndex"`
	//Port       uint32     `yaml:"port"`
}

// ServiceModel service model information
type ServiceModel struct {
	ID          int    `yaml:"id"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
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
