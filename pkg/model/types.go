// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

// EnbID  E2 node ID
type EnbID string

// PlmnID plmnID
type PlmnID string

// Ecgi Ecgi
type Ecgi string

// GEnbID global E2 node ID
type GEnbID struct {
	PlmnID PlmnID
	EnbID  EnbID
}

// Imsi represents a unique UE identifier
type Imsi uint64

// Crnti is a tower-specific UE identifier
type Crnti string

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
