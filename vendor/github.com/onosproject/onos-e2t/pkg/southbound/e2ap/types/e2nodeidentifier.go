// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package types

import "fmt"

// PlmnID - a 3 byte representation of the Plmn ID
//- digits 0 to 9, encoded 0000 to 1001,
//- 1111 used as filler digit,
//two digits per octet,
//- bits 4 to 1 of octet n encoding digit 2n-1
//- bits 8 to 5 of octet n encoding digit 2n
//
//-The PLMN identity consists of 3 digits from MCC followed by either
//-a filler digit plus 2 digits from MNC (in case of 2 digit MNC) or
//-3 digits from MNC (in case of a 3 digit MNC).
type PlmnID [3]byte

func PlmnIDFromSlice(plmnIDBytes []byte) (PlmnID, error) {
	if len(plmnIDBytes) != 3 {
		return [3]byte{}, fmt.Errorf("plmnID must be 3 bytes long")
	}
	plmnID := [3]byte{plmnIDBytes[0], plmnIDBytes[1], plmnIDBytes[2]}
	return plmnID, nil
}

type E2NodeType int32

const (
	E2NodeTypeGNB E2NodeType = iota
	E2NodeTypeEnGNB
	E2NodeTypeNgENB
	E2NodeTypeENB
)

// E2NodeIdentity a simplified model of the E2 Node
type E2NodeIdentity struct {
	Plmn           PlmnID
	NodeType       E2NodeType
	NodeIdentifier []byte
}

func NewE2NodeIdentity(plmnIDSlice []byte) (*E2NodeIdentity, error) {
	plmnID, err := PlmnIDFromSlice(plmnIDSlice)
	if err != nil {
		return nil, err
	}

	return &E2NodeIdentity{
		Plmn: plmnID,
	}, nil
}
