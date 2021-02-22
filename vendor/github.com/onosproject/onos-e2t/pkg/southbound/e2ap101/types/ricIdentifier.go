// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package types

type RicIdentifierBits uint64
type RicIdentifierLen uint32

type RicIdentifier struct {
	RicIdentifierValue RicIdentifierBits
	RicIdentifierLen   RicIdentifierLen
}

type RicIdentity struct {
	RicIdentifier RicIdentifier
	PlmnID        PlmnID
}
