// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package metrics

// MetricEvent is a type of event
type MetricEvent int

const (
	// None none cell event
	None MetricEvent = iota
	// Updated updated metric event
	Updated
	// Deleted deleted metric event
	Deleted
)

func (e MetricEvent) String() string {
	return [...]string{"None", "Updated", "Deleted"}[e]
}
