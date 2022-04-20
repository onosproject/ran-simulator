// SPDX-FileCopyrightText: 2022-present Intel Corporation
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
	// ORAN-E2SM-KPM version 2
	Kpm2
	// MHO
	Mho
	// O-RAN-E2SM-RC
	Rc
)
