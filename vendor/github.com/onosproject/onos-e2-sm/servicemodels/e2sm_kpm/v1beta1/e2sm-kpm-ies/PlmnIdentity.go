// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2smkpmies

import "fmt"

type PlmnIdentityBuilder interface {
	NewPlmnID(value []byte)
	GetValue()
	GetPlmnID()
}

func NewPlmnID(value []byte) (*PlmnIdentity, error) {
	if len(value) != 3 {
		return nil, fmt.Errorf("Size of the PlmnID must be 3 bytes")
	}
	return &PlmnIdentity{
		Value: value,
	}, nil
}

func (b *PlmnIdentity) GetPlmnID() *PlmnIdentity {
	return b
}
