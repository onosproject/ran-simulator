// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package types

// EcID is a cell ID
type EcID uint64

// PlmnID is a network ID
type PlmnID uint16

// Ecgi is a globally unique cell ID
type Ecgi uint64

// Crnti is a UE ID relative to a tower
type Crnti uint32

// Imsi is a globally unique UE identifier
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

	// AddressKey ...
	AddressKey = "address"
)

const (
	// E2NodeType - used in topo device type
	E2NodeType = "E2Node"

	// E2NodeVersion100 - used in topo device version
	E2NodeVersion100 = "1.0.0"
)
