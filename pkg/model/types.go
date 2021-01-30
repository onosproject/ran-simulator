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
