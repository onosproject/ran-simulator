// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

// RanFunctionID Ran function ID
type RanFunctionID int32

// TODO define them using standard Ran function IDs
const (

	// Internal
	Internal RanFunctionID = iota
	// Kpm
	Kpm
	// Ni
	Ni
	// Rc
	Rc
	// Kpm2
	Kpm2
)
