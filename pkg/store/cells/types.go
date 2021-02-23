// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cells

// CellEvent a cell event
type CellEvent int

const (
	// None none cell event
	None CellEvent = iota
	// Created created cell event
	Created
	// Updated updated cell event
	Updated
	// Deleted deleted cell event
	Deleted
)

func (e CellEvent) String() string {
	return [...]string{"None", "Created", "Updated", "Deleted"}[e]
}
