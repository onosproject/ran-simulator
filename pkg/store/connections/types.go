// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connections

import (
	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap"
)

// ConnectionEvent a connection event
type ConnectionEvent int

const (
	// None none connection event
	None ConnectionEvent = iota
	// Created created  event
	Created
	// Updated updated connection event
	Updated
	// Deleted deleted  connection event
	Deleted
)

// String converts connection event to string
func (e ConnectionEvent) String() string {
	return [...]string{"None", "Created", "Updated", "Deleted"}[e]
}

// ConnectionID consists of IP and port number of E2T instance
type ConnectionID struct {
	ricIPAddress string
	ricPort      uint64
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

// ConnectionStatus connection status
type ConnectionStatus struct {
	Phase Phase
	State State
}

// State channel state
type State int

const (

	// Connecting connecting state
	Connecting State = iota

	// Connected connected state
	Connected

	// Configuring configuring state
	Configuring

	// Configured configured state
	Configured

	// Disconnecting disconnecting state
	Disconnecting

	// Disconnected disconnected state
	Disconnected
)

// String return state in string format
func (s State) String() string {
	return [...]string{"Connecting", "Connected", "Configuring", "Configured", "Disconnecting", "Disconnected"}[s]
}

// Connection connection data for storing in connection store
type Connection struct {
	ID     ConnectionID
	Client e2.ClientConn
	Status ConnectionStatus
}
