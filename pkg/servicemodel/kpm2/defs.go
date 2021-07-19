// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm2

// MeasTypeName name of measurement type
type MeasTypeName int

const (
	// RRCConnEstabAttSum total number of RRC connection establishment attempts
	RRCConnEstabAttSum MeasTypeName = iota
	// RRCConnEstabSuccSum  total number of successful RRC Connection establishments
	RRCConnEstabSuccSum
	// RRCConnReEstabAttSum total number of RRC connection re-establishment attempts
	RRCConnReEstabAttSum
	// RRCConnReEstabAttreconfigFail  total number of RRC connection re-establishment attempts due to reconfiguration failure
	RRCConnReEstabAttreconfigFail
	// RRCConnReEstabAttHOFail total number of RRC connection re-establishment attempts due to Handover failure
	RRCConnReEstabAttHOFail
	// RRCConnReEstabAttOther total number of RRC connection re-establishment attempts due to Other reasons
	RRCConnReEstabAttOther
	// RRCConnAvg the mean number of users in RRC connected mode during each granularity period.
	RRCConnAvg
	// RRCConnMax  the max number of users in RRC connected mode during each granularity period.
	RRCConnMax
)

func (m MeasTypeName) String() string {
	return [...]string{"RRC.ConnEstabAtt.Sum",
		"RRC.ConnEstabSucc.Sum",
		"RRC.ConnReEstabAtt.Sum",
		"RRC.ConnReEstabAtt.reconfigFail",
		"RRC.ConnReEstabAtt.HOFail",
		"RRC.ConnReEstabAtt.Other",
		"RRC.Conn.Avg",
		"RRC.Conn.Max"}[m]
}

// MeasType meas type
type MeasType struct {
	measTypeName MeasTypeName
	measTypeID   int32
}

var measTypes = []MeasType{
	{
		measTypeName: RRCConnEstabAttSum,
		measTypeID:   1,
	},
	{
		measTypeName: RRCConnEstabSuccSum,
		measTypeID:   2,
	},
	{
		measTypeName: RRCConnReEstabAttSum,
		measTypeID:   3,
	},
	{
		measTypeName: RRCConnReEstabAttreconfigFail,
		measTypeID:   4,
	},
	{
		measTypeName: RRCConnReEstabAttHOFail,
		measTypeID:   5,
	},
	{
		measTypeName: RRCConnReEstabAttOther,
		measTypeID:   6,
	},
	{
		measTypeName: RRCConnAvg,
		measTypeID:   7,
	},
	{
		measTypeName: RRCConnMax,
		measTypeID:   8,
	},
}
