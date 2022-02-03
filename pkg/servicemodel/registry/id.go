// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package registry

// RanFunctionID Ran function ID
type RanFunctionID int32

// ModelOid service model OID
type ModelOid string

// TODO define them using standard Ran function IDs
const (

	// Internal
	Internal RanFunctionID = iota
	// Kpm
	Kpm
	// Ni
	Ni
	// Rcpre2
	Rcpre2
	// Kpm2
	Kpm2
	// MHO
	Mho
)
