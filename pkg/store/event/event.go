// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package event

// Event store event data structure
type Event struct {
	Key   interface{}
	Value interface{}
	Type  interface{}
}
