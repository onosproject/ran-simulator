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
	Nodes     *SimNodes
	Locations *SimLocations
	// MapLayout   *types.MapLayout
	// AspectRatio float64 // fold into the map layout?
	// UEs      *SimUserEquipment
	// Routes   *SimRoutes
}

// NewModel creates a new model of the simulated environment
func NewModel() *Model {
	return &Model{
		Nodes:     NewSimNodes(),
		Locations: nil,
	}
}

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

	// Parse the nodes first
	for _, val := range data["nodes"].([]interface{}) {
		node, err := parseNode(val)
		if err == nil {
			m.Nodes.Add(node)
		}
	}

	// TODO: parse map layout, locations, UE constraints, etc.
	return nil
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
