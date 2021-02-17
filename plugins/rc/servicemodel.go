// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package main

import (
	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"
)

func (s serviceModel) ServiceModelData() (string, string, string) {
	return "e2sm_rc_pre", "v1", ""
}

func (s serviceModel) IndicationHeaderASN1toProto(asn1Bytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) IndicationHeaderProtoToASN1(protoBytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) IndicationMessageASN1toProto(asn1Bytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) IndicationMessageProtoToASN1(protoBytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) RanFuncDescriptionASN1toProto(asn1Bytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) RanFuncDescriptionProtoToASN1(protoBytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) EventTriggerDefinitionASN1toProto(asn1Bytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) EventTriggerDefinitionProtoToASN1(protoBytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) ActionDefinitionASN1toProto(asn1Bytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) ActionDefinitionProtoToASN1(protoBytes []byte) ([]byte, error) {
	// TODO
	var bytes []byte
	return bytes, nil
}

func (s serviceModel) DecodeRanFunctionDescription(asn1bytes []byte) (*e2smtypes.RanfunctionNameDef, *e2smtypes.RicEventTriggerList, *e2smtypes.RicReportList, error) {
	// TODO
	ranfunctionNameDef := &e2smtypes.RanfunctionNameDef {
		RanFunctionShortName: "a",
		RanFunctionE2SmOid: "a",
		RanFunctionDescription: "a",
		RanFunctionInstance: 1,
	}

	ricEventTriggerList := &e2smtypes.RicEventTriggerList {
	}

	ricReportList := &e2smtypes.RicReportList {
	}

	return ranfunctionNameDef, ricEventTriggerList, ricReportList, nil
}
