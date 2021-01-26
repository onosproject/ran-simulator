// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2smkpmies

type BitStringBuilder interface {
	NewBitString(value uint64, len uint32)
	GetValue()
	GetLen()
	GetBitString()
}

func NewBitString(value uint64, len uint32) *BitString {
	return &BitString{
		Value: value,
		Len:   len,
	}
}

func (b *BitString) GetBitString() *BitString {
	return b
}
