// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package types

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

const (
	// MaxUint24 maximum value of uint24 variable
	MaxUint24 = 1<<24 - 1
)

// Uint24 unti24 type
type Uint24 [3]uint8

// NewUint24 creates a new uint24 type
func NewUint24(val uint32) *Uint24 {
	var u = new(Uint24)
	u.Set(val)
	return u
}

// Value returns value
func (u *Uint24) Value() Uint24 {
	return *u
}

// Set sets uint24 value
func (u *Uint24) Set(val uint32) {
	if val > MaxUint24 {
		return
	}
	(*u)[0] = uint8(val & 0xFF)
	(*u)[1] = uint8((val >> 8) & 0xFF)
	(*u)[2] = uint8((val >> 16) & 0xFF)
}

// Uint32 converts uint24 to uint32
func (u Uint24) Uint32() uint32 {
	return uint32(u[0]) | uint32(u[1])<<8 | uint32(u[2])<<16
}

// String converts uint24 to string
func (u Uint24) String() string {
	return strconv.Itoa(int(u.Uint32()))
}

// ToBytes converts uint24 to bytes array
func (u Uint24) ToBytes() []byte {
	var buf = &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, u); err != nil {
		return nil
	}
	return buf.Bytes()
}

// Uint24ToUint32 converts uint24 uint32
func Uint24ToUint32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 3; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}
