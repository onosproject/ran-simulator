// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package channels

import e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap101"

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

// Phase channel phase
type Phase int

const (
	// Open open phase
	Open Phase = iota

	// Closed closed state
	Closed
)

// String return phase
func (p Phase) String() string {
	return [...]string{"Open", "Closed"}[p]
}

// ChannelStatus channel status
type ChannelStatus struct {
	Phase Phase
	State State
}

// State channel state
type State int

const (
	// Completed completed state
	Completed State = iota

	// Pending pending state
	Pending

	// Failed failed state
	Failed
)

// String return state in string format
func (s State) String() string {
	return [...]string{"Completed", "Pending", "Failed"}[s]
}

// Channel channel data structure for storing in channel store
type Channel struct {
	ID     ChannelID
	Client e2.ClientChannel
	Status ChannelStatus
}
