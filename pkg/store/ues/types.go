// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ues

// UeEvent a node event
type UeEvent int

const (
	// None non ue event
	None UeEvent = iota
	// Created created ue event
	Created
	// Updated updated ue event
	Updated
	// Deleted deleted  ue event
	Deleted
)

// String converts node event to string
func (e UeEvent) String() string {
	return [...]string{"None", "Created", "Updated", "Deleted"}[e]
}
