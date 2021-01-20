// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Model holds the information describing the simulated RAN environment.
type Model struct {
	Controllers []*Controller
	Nodes       *SimNodes
	Locations   *SimLocations
	// MapLayout   *types.MapLayout
	// AspectRatio float64 // fold into the map layout?
	// UEs      *SimUserEquipment
	// Routes   *SimRoutes
}

// Controller represents an E2T node of the RAN controller platform
type Controller struct {
	Address string
	Port    int
}

// NewModel creates a new model of the simulated environment
func NewModel() *Model {
	return &Model{
		Nodes:     NewSimNodes(),
		Locations: nil,
	}
}

const (
	controllersKey = "controllers"
	nodesKey       = "nodes"
)

// Load from the specified YAML file
func (m *Model) Load(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	// Parse controllers
	m.Controllers = make([]*Controller, 0)
	if data[controllersKey] != nil {
		for _, val := range data[controllersKey].([]interface{}) {
			ctl, err := parseController(val)
			if err == nil {
				m.Controllers = append(m.Controllers, ctl)
			}
		}
	}

	// Parse the nodes next
	if data[nodesKey] != nil {
		for _, val := range data[nodesKey].([]interface{}) {
			node, err := parseNode(val)
			if err == nil {
				m.Nodes.Add(node)
			}
		}
	}

	// TODO: parse map layout, locations, UE constraints, etc.
	return nil
}

// Parse YAML for a controller entry
func parseController(val interface{}) (*Controller, error) {
	r := val.(map[interface{}]interface{})
	return &Controller{
		Address: r["address"].(string),
		Port:    r["port"].(int),
	}, nil
}

// Parse YAML for a simulated node entry
func parseNode(val interface{}) (*SimNode, error) {
	r := val.(map[interface{}]interface{})
	ecgi := ECGI(r["ecgi"].(string))
	smr := r["sms"].([]interface{})

	sms := make([]SMID, 0)
	for _, id := range smr {
		sms = append(sms, SMID(id.(int)))
	}

	return &SimNode{
		ECGI:          ecgi,
		ServiceModels: sms,
	}, nil
}
