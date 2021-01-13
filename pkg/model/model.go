// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package model

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
