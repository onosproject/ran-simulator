// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2smkpmies

type RanfunctionNameBuilder interface {
	NewRanfunctionName(shortName string, oid string, description string, instance int32)
	GetRanFunctionShortName()
	GetRanFunctionE2SmOid()
	GetRanFunctionDescription()
	GetRanFunctionInstance()
	GetRanfunctionName()
}

func NewRanfunctionName(shortName string, oid string, description string, instance int32) *RanfunctionName {
	return &RanfunctionName{
		RanFunctionShortName:   shortName,
		RanFunctionE2SmOid:     oid,
		RanFunctionDescription: description,
		RanFunctionInstance:    instance,
	}
}

func (b *RanfunctionName) GetRanfunctionName() *RanfunctionName {
	return b
}
