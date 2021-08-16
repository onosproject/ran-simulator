// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package channels

// ChannelEvent a channel event
type ChannelEvent int

const (
	// None none channel event
	None ChannelEvent = iota
	// Created created  event
	Created
	// Updated updated channel event
	Updated
	// Deleted deleted  channel event
	Deleted
)

// String converts channel event to string
func (e ChannelEvent) String() string {
	return [...]string{"None", "Created", "Updated", "Deleted"}[e]
}

// ChannelID channel ID consists of IP and port number of E2T instance
type ChannelID struct {
	ricAddress string
	ricPort    uint64
}
