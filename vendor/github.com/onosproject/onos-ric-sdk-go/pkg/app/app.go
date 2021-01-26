// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package app

// ID is an application identifier
type ID string

// InstanceID is an application instance identifier
type InstanceID string

// AppContext is the main application context; an object on which all SDK app-visible functions reside.
type AppContext struct {
	ID         ID
	InstanceID InstanceID

	// other context required by various SDK facilities
}
