// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package types

import (
	"fmt"
	"strconv"
)

// PlmnID is a globally unique network identifier (Public Land Mobile Network)
type PlmnID uint32

// EnbID is an eNodeB Identifier
type EnbID uint32

// CellID is a node-local cell identifier
type CellID uint8

// ECI is a E-UTRAN Cell Identifier
type ECI uint32

// GEnbID is a Globally eNodeB identifier
type GEnbID uint64

// ECGI is E-UTRAN Cell Global Identifier
type ECGI uint64

// CRNTI is a cell-specific UE identifier
type CRNTI uint32

// MSIN is Mobile Subscriber Identification Number
type MSIN uint32

// IMSI is International Mobile Subscriber Identity
type IMSI uint64

const (
	mask28               = 0xfffffff
	mask20               = 0xfffff00
	lowest24             = 0x0ffffff
	maskSecondNibble     = 0x00000f0
	maskSeventhNibble    = 0xf000000
	maskThirteenthNibble = 0xf000000000000
)

// EncodePlmnID encodes MCC and MNC strings into a PLMNID hex string
func EncodePlmnID(mcc string, mnc string) string {
	if len(mnc) == 2 {
		return string(mcc[1]) + string(mcc[0]) + "F" + string(mcc[2]) + string(mnc[1]) + string(mnc[0])
	} else {
		return string(mcc[1]) + string(mcc[0]) + string(mnc[2]) + string(mcc[2]) + string(mnc[1]) + string(mnc[0])
	}
}

// DecodePlmnID decodes MCC and MNC strings from PLMNID hex string
func DecodePlmnID(plmnID string) (mcc string, mnc string) {
	if plmnID[2] == 'f' || plmnID[2] == 'F' {
		return string(plmnID[1]) + string(plmnID[0]) + string(plmnID[3]),
			string(plmnID[5]) + string(plmnID[4])
	} else {
		return string(plmnID[1]) + string(plmnID[0]) + string(plmnID[3]),
			string(plmnID[5]) + string(plmnID[4]) + string(plmnID[2])
	}
}

// ToPlmnID encodes the specified MCC and MNC strings into a numeric PLMNID
func ToPlmnID(mcc string, mnc string) PlmnID {
	s := EncodePlmnID(mcc, mnc)
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0
	}
	return PlmnID(n)
}

// PlmnIDFromHexString converts string form of PLMNID in its hex form into a numeric one suitable for APIs
func PlmnIDFromHexString(plmnID string) PlmnID {
	n, err := strconv.ParseUint(plmnID, 16, 32)
	if err != nil {
		return 0
	}
	return PlmnID(n)
}

// PlmnIDFromString converts string form of PLMNID given as a simple MCC-MCN catenation into a numeric one suitable for APIs
func PlmnIDFromString(plmnID string) PlmnID {
	return ToPlmnID(plmnID[0:3], plmnID[3:])
}

// PlmnIDToString generates the MCC-MCN catenation format from the specified numeric PLMNID
func PlmnIDToString(plmnID PlmnID) string {
	hexString := fmt.Sprintf("%x", plmnID)
	mcc, mnc := DecodePlmnID(hexString)
	return mcc + mnc
}

// ToECI produces ECI from the specified components
func ToECI(enbID EnbID, cid CellID) ECI {
	if cid&maskSecondNibble == 0 {
		return ECI(uint(enbID)<<4 | uint(cid)) // Unclear whether this clause is needed
	}
	return ECI(uint(enbID)<<8 | uint(cid))
}

// ToECGI produces ECGI from the specified components
func ToECGI(plmnID PlmnID, eci ECI) ECGI {
	if uint(eci)&maskSeventhNibble == 0 {
		return ECGI(uint(plmnID)<<24 | (uint(eci) & mask28)) // Unclear whether this clause is needed
	}
	return ECGI(uint(plmnID)<<28 | (uint(eci) & mask28))
}

// ToGEnbID produces GEnbID from the specified components
func ToGEnbID(plmnID PlmnID, enbID EnbID) GEnbID {
	return GEnbID(uint(plmnID)<<28 | (uint(enbID) << 8 & mask20))
}

// GetPlmnID extracts PLMNID from the specified ECGI, GEnbID or IMSI
func GetPlmnID(id uint64) PlmnID {
	if id&maskThirteenthNibble == 0 {
		return PlmnID(id >> 24)
	}
	return PlmnID(id >> 28)
}

// GetCellID extracts Cell ID from the specified ECGI or GEnbID
func GetCellID(id uint64) CellID {
	if id&maskThirteenthNibble == 0 {
		return CellID(id & 0xf)
	}
	return CellID(id & 0xff)
}

// GetEnbID extracts Enb ID from the specified ECGI or GEnbID
func GetEnbID(id uint64) EnbID {
	if id&maskThirteenthNibble == 0 {
		return EnbID((id & mask20) >> 4)
	}
	return EnbID((id & mask20) >> 8)
}

// GetECI extracts ECI from the specified ECGI or GEnbID
func GetECI(id uint64) ECI {
	if id&maskThirteenthNibble == 0 {
		return ECI(id & lowest24)
	}
	return ECI(id & mask28)
}

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
