// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package dispatcher

import "github.com/onosproject/ran-simulator/api/trafficsim"

// Event is a stream event
type Event struct {
	// Type is the stream event type
	Type trafficsim.Type

	// UpdateType is a qualification on the type of update
	UpdateType trafficsim.UpdateType

	// Object is the event object
	Object interface{}
}
