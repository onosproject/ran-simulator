// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package types

// A sequence number 0-65535 e2ap-v01.00.00.asn1.asn1 line 1117
type RicIndicationSn uint16

// The E2SM Indication Header e2ap-v01.00.00.asn1.asn1 line 1110
type RicIndicationHeader []byte

// The E2SM Indication Header e2ap-v01.00.00.asn1.asn1 line 1115
type RicIndicationMessage []byte
