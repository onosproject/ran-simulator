// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package nodes

// NodeEvent a node event
type NodeEvent int

const (
	// None non node event
	None NodeEvent = iota
	// Created created node event
	Created
	// Updated updated node event
	Updated
	// Deleted deleted  node event
	Deleted
)

// String converts node event to string
func (e NodeEvent) String() string {
	return [...]string{"None", "Created", "Updated", "Deleted"}[e]
}
