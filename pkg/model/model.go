// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
)

// Model simulation model
type Model struct {
	MapLayout     MapLayout               `mapstructure:"layout" yaml:"layout"`
	Nodes         map[string]Node         `mapstructure:"nodes" yaml:"nodes"`
	Cells         map[string]Cell         `mapstructure:"cells" yaml:"cells"`
	Controllers   map[string]Controller   `mapstructure:"controllers" yaml:"controllers"`
	ServiceModels map[string]ServiceModel `mapstructure:"servicemodels" yaml:"servicemodels"`
	UECount       uint                    `mapstructure:"ueCount" yaml:"ueCount"`
	Plmn          string                  `mapstructure:"plmnID" yaml:"plmnID"`
	PlmnID        types.PlmnID            `mapstructure:"plmnNumber" yaml:"plmnNumber"` // overridden and derived post-load from "Plmn" field
	APIKey        string                  `mapstructure:"apiKey" yaml:"apiKey"`         // Google Maps API key (optional)
}

// Coordinate represents a geographical location
type Coordinate struct {
	Lat float64 `mapstructure:"lat"`
	Lng float64 `mapstructure:"lng"`
}

// Sector represents a 2D arc emanating from a location
type Sector struct {
	Center  Coordinate `mapstructure:"center"`
	Azimuth int32      `mapstructure:"azimuth"`
	Arc     int32      `mapstructure:"arc"`
}

// Route represents a series of points for tracking movement of user-equipment
type Route struct {
	IMSI        types.IMSI
	Points      []*Coordinate
	Color       string
	SpeedAvg    uint32
	SpeedStdDev uint32
	Reverse     bool
	NextPoint   uint32
}

// Node e2 node
type Node struct {
	EnbID         types.EnbID  `mapstructure:"enbID"`
	Controllers   []string     `mapstructure:"controllers"`
	ServiceModels []string     `mapstructure:"servicemodels"`
	Cells         []types.ECGI `mapstructure:"cells"`
	Status        string       `mapstructure:"status"`
}

// Controller E2T endpoint information
type Controller struct {
	ID      string `mapstructure:"id"`
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

// EventA3Params represents a set of A3 handover event parameters
type EventA3Params struct {
	A3Offset          int32 `mapstructure:"a3Offset"`          // Integer: -30 to 30
	A3TimeToTrigger   int32 `mapstructure:"a3TimeToTrigger"`   // timetotrigger: 0ms, 40ms, 64ms, 80ms, 100ms, 128ms, 160ms, 256ms, 320ms, 480ms, 512ms, 640ms, 1024ms, 1280ms, 2560ms, 5120ms
	A3Hysteresis      int32 `mapstructure:"a3Hysteresis"`      // Integer: 0 to 30
	A3CellOffset      int32 `mapstructure:"a3CellOffset"`      // q-offsetrange: -24, -22, -20, -18, -16, -14, -12, -10, -8, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24
	A3FrequencyOffset int32 `mapstructure:"a3FrequencyOffset"` // q-offsetrange: -24, -22, -20, -18, -16, -14, -12, -10, -8, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24
}

// Cell represents a section of coverage
type Cell struct {
	ECGI          types.ECGI    `mapstructure:"ecgi"`
	Sector        Sector        `mapstructure:"sector"`
	Color         string        `mapstructure:"color"`
	MaxUEs        uint32        `mapstructure:"maxUEs"`
	Neighbors     []types.ECGI  `mapstructure:"neighbors"`
	TxPowerDB     float64       `mapstructure:"txPower"`
	EventA3Params EventA3Params `mapstructure:"eventA3Params"`
}

// UEType represents type of user-equipment
type UEType string

// UECell represents UE-cell relationship
type UECell struct {
	ID       types.GEnbID
	ECGI     types.ECGI // Auxiliary form of association
	Strength float64
}

// UE represents user-equipment, i.e. phone, IoT device, etc.
type UE struct {
	IMSI     types.IMSI
	Type     UEType
	Location Coordinate
	Heading  uint32

	Cell  *UECell
	CRNTI types.CRNTI
	Cells []*UECell

	IsAdmitted bool
}

// ServiceModel service model information
type ServiceModel struct {
	ID          int    `mapstructure:"id"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
}

// GetServiceModel gets a service model based on a given name.
func (m *Model) GetServiceModel(name string) (ServiceModel, error) {
	if sm, ok := m.ServiceModels[name]; ok {
		return sm, nil
	}
	return ServiceModel{}, errors.New(errors.NotFound, "the service model not found")
}

// GetController gets a controller by a given name
func (m *Model) GetController(name string) (Controller, error) {
	if controller, ok := m.Controllers[name]; ok {
		return controller, nil
	}
	return Controller{}, errors.New(errors.NotFound, "controller not found")
}
