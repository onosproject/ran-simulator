// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package kpm2

// MeasTypeName name of measurement type
type MeasTypeName int

const (
	PDUSessionSetupReq MeasTypeName = iota

	PDUSessionSetupSucc

	PDUSessionSetupFail

	PrbUsedDL

	PrbUsedUL

	PdcpPduVolumeDL

	PdcpPduVolumeUL

	PdcpRatePerPRBDL

	PdcpRatePerPRBUL
	// RRCConnEstabAttSum total number of RRC connection establishment attempts
	RRCConnEstabAttSum
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
	return [...]string{
		"PDUSessionSetupReq",
		"PDUSessionSetupSucc",
		"PDUSessionSetupFail",
		"PrbUsedDL",
		"PrbUsedUL",
		"PdcpPduVolumeDL",
		"PdcpPduVolumeUL",
		"PdcpRatePerPRBDL",
		"PdcpRatePerPRBUL",
		"RRC.ConnEstabAtt.Sum",
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
		measTypeName: PDUSessionSetupReq,
		measTypeID:   1,
	},
	{
		measTypeName: PDUSessionSetupSucc,
		measTypeID:   2,
	},
	{
		measTypeName: PDUSessionSetupFail,
		measTypeID:   3,
	},
	{
		measTypeName: PrbUsedDL,
		measTypeID:   4,
	},
	{
		measTypeName: PrbUsedUL,
		measTypeID:   5,
	},
	{
		measTypeName: PdcpPduVolumeDL,
		measTypeID:   6,
	},
	{
		measTypeName: PdcpPduVolumeUL,
		measTypeID:   7,
	},
	{
		measTypeName: PdcpRatePerPRBDL,
		measTypeID:   8,
	},
	{
		measTypeName: PdcpRatePerPRBUL,
		measTypeID:   9,
	},
	{
		measTypeName: RRCConnEstabAttSum,
		measTypeID:   10,
	},
	{
		measTypeName: RRCConnEstabSuccSum,
		measTypeID:   11,
	},
	{
		measTypeName: RRCConnReEstabAttSum,
		measTypeID:   12,
	},
	{
		measTypeName: RRCConnReEstabAttreconfigFail,
		measTypeID:   13,
	},
	{
		measTypeName: RRCConnReEstabAttHOFail,
		measTypeID:   14,
	},
	{
		measTypeName: RRCConnReEstabAttOther,
		measTypeID:   15,
	},
	{
		measTypeName: RRCConnAvg,
		measTypeID:   16,
	},
	{
		measTypeName: RRCConnMax,
		measTypeID:   17,
	},
}
