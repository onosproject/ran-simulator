// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package utils

type GNbIDBytes [4]byte

const (
	MaxGNbID22 = 1<<23 - 1
	MaxGNbID32 = 1<<33 - 1
)

func NewGNbIDWithUint64(val uint64, length int) *GNbIDBytes {
	id := new(GNbIDBytes)
	id.Set(val, length)
	return id
}

func (b *GNbIDBytes) Set(val uint64, length int) {
	if length != 32 && length != 22 {
		return
	} else if length == 22 && val > MaxGNbID22 {
		return
	} else if length == 32 && val > MaxGNbID32 {
		return
	}

	(*b)[0] = byte(val & 0xFF)
	(*b)[1] = byte(val >> 8 & 0xFF)
	(*b)[2] = byte(val >> 16 & 0xFF)
	(*b)[3] = byte(val >> 24 & 0xFF)
}

func (b *GNbIDBytes) Value() GNbIDBytes {
	return *b
}

func (b *GNbIDBytes) Bytes(length int) []byte {
	if length == 22 {
		val := make([]byte, 3)
		val[0] = (*b)[0]
		val[1] = (*b)[1]
		val[2] = (*b)[2]
		return val
	} else if length == 32 {
		val := make([]byte, 4)
		val[0] = (*b)[0]
		val[1] = (*b)[1]
		val[2] = (*b)[2]
		val[3] = (*b)[3]
		return val
	}
	return nil
}

func (b *GNbIDBytes) Uint64() uint64 {
	return uint64((*b)[0]) + (uint64((*b)[1]) << 8) + (uint64((*b)[2]) << 16) + (uint64((*b)[3]) << 24)
}

type GNbID struct {
	IDByte *GNbIDBytes
	Length int
}

func NewGNbID(val uint64, length int) *GNbID {
	return &GNbID{
		IDByte: NewGNbIDWithUint64(val, length),
		Length: length,
	}
}
