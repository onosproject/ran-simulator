// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

// NCellID NCellID type
type NCellID [5]byte

const (
	// MaxNCellID maximum value of NCellID variable
	MaxNCellID = 1<<40 - 1
)

// NewNCellIDWithUint64 creates a new NCellID type with UInt64
func NewNCellIDWithUint64(val uint64) *NCellID {
	id := new(NCellID)
	id.Set(val)
	return id
}

// NewNCellIDWithBytes creates a new NCellID type with bytes
func NewNCellIDWithBytes(val []byte) *NCellID {
	id := new(NCellID)
	id[0] = val[0]
	id[1] = val[1]
	id[2] = val[2]
	id[3] = val[3]
	id[4] = val[4]
	return id
}

// Set sets NCellID
func (n *NCellID) Set(val uint64) {
	if val > MaxNCellID {
		return
	}
	(*n)[0] = byte(val & 0xFF)
	(*n)[1] = byte(val >> 8 & 0xFF)
	(*n)[2] = byte(val >> 16 & 0xFF)
	(*n)[3] = byte(val >> 24 & 0xFF)
	(*n)[4] = byte(val >> 32 & 0xFF)
}

// Value returns NCellID value
func (n *NCellID) Value() NCellID {
	return *n
}

// Bytes converts NCellID to byte array
func (n *NCellID) Bytes() []byte {
	val := make([]byte, 5)
	val[0] = (*n)[0]
	val[1] = (*n)[1]
	val[2] = (*n)[2]
	val[3] = (*n)[3]
	val[4] = (*n)[4]
	return val
}

// Uint64 converts NCellID to uint64
func (n *NCellID) Uint64() uint64 {
	return uint64((*n)[0]) + (uint64((*n)[1]) << 8) + (uint64((*n)[2]) << 16) + (uint64((*n)[3]) << 24) + (uint64((*n)[4]) << 32)
}
