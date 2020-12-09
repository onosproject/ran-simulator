// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package types

import "github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"

type RanFunctionDescription []byte
type RanFunctionRevision int
type RanFunctionID uint8

type RanFunctionItem struct {
	Description RanFunctionDescription
	Revision    RanFunctionRevision
}

type RanFunctions map[RanFunctionID]RanFunctionItem

type RanFunctionRevisions map[RanFunctionID]RanFunctionRevision

type RanFunctionCauses map[RanFunctionID]*e2apies.Cause
