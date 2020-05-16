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

package types

// EcID is a tower ID
type EcID string

// PlmnID is a network ID
type PlmnID string

// Crnti is a UE ID relative to a tower
type Crnti string

// Imsi is a UE unique identifier
type Imsi uint64

const (
	// AzimuthKey - used in topo device attributes
	AzimuthKey = "azimuth"

	// ArcKey - used in topo device attributes
	ArcKey = "arc"

	// LatitudeKey - used in topo device attributes
	LatitudeKey = "latitude"

	// LongitudeKey - used in topo device attributes
	LongitudeKey = "longitude"

	// EcidKey - used in topo device attributes
	EcidKey = "ecid"

	// PlmnIDKey - used in topo device attributes
	PlmnIDKey = "plmnid"

	// GrpcPortKey - used in topo device attributes
	GrpcPortKey = "grpcport"
)

const (
	// E2NodeType - used in topo device type
	E2NodeType = "E2Node"

	// E2NodeVersion100 - used in topo device version
	E2NodeVersion100 = "1.0.0"
)
