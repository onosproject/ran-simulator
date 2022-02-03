// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package routes

// RouteEvent a node event
type RouteEvent int

const (
	// None non route event
	None RouteEvent = iota
	// Created created ue event
	Created
	// Updated updated ue event
	Updated
	// Deleted deleted  ue event
	Deleted
)

// String converts node event to string
func (e RouteEvent) String() string {
	return [...]string{"None", "Created", "Updated", "Deleted"}[e]
}
