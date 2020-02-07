// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
