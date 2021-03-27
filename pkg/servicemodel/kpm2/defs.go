// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm2

type MeasTypeName int

const (
	RRCConnEstabAttTot MeasTypeName = iota
	RRCConnEstabSuccTot
	RRCConnReEstabAttTot
	RRCConnReEstabAttreconfigFail
	RRCConnReEstabAttHOFail
	RRCConnReEstabAttOther
	RRCConnAvg
	RRCConnMax
)

func (m MeasTypeName) String() string {
	return [...]string{"RRC.ConnEstabAtt.Tot",
		"RRC.ConnEstabSucc.Tot",
		"RRC.ConnReEstabAtt.Tot",
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
		measTypeName: RRCConnEstabAttTot,
		measTypeID:   1,
	},
	{
		measTypeName: RRCConnEstabSuccTot,
		measTypeID:   2,
	},
	{
		measTypeName: RRCConnReEstabAttTot,
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
