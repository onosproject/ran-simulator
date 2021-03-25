// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package modelplugins

import (
	"fmt"
	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"
	kpmv2ctypes "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/kpmctypes"
	e2sm_kpm_v2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
	ransimmps "github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"
)

const smName = "e2sm_kpm"
const smVersion = "v2"
const moduleName = "e2sm_kpm_v2.so.2.0"
const smOIDKpmV2 = "1.3.6.1.4.1.53148.1.2.2.2"

type serviceModelTest string

var serviceModelTestInst serviceModelTest = ""

func (sm serviceModelTest) ServiceModelData() e2smtypes.ServiceModelData {
	smData := e2smtypes.ServiceModelData{
		Name:       smName,
		Version:    smVersion,
		ModuleName: moduleName,
		OID:        smOIDKpmV2,
	}
	return smData
}

func (sm serviceModelTest) IndicationHeaderASN1toProto(asn1Bytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) IndicationHeaderProtoToASN1(protoBytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) IndicationMessageASN1toProto(asn1Bytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) IndicationMessageProtoToASN1(protoBytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) RanFuncDescriptionASN1toProto(asn1Bytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) RanFuncDescriptionProtoToASN1(protoBytes []byte) ([]byte, error) {
	protoObj := new(e2sm_kpm_v2.E2SmKpmRanfunctionDescription)
	if err := proto.Unmarshal(protoBytes, protoObj); err != nil {
		return nil, fmt.Errorf("error unmarshalling protoBytes to E2SmKpmRanfunctionDescription %s", err)
	}

	perBytes, err := kpmv2ctypes.PerEncodeE2SmKpmRanfunctionDescription(protoObj)
	if err != nil {
		return nil, fmt.Errorf("error encoding E2SmKpmRanfunctionDescription to PER %s", err)
	}

	return perBytes, nil
}

func (sm serviceModelTest) EventTriggerDefinitionASN1toProto(asn1Bytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) EventTriggerDefinitionProtoToASN1(protoBytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) ActionDefinitionASN1toProto(asn1Bytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) ActionDefinitionProtoToASN1(protoBytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on serviceModelTest")
}

//It is redundant so far - could be reused for future, if you need to extract something specific from RanFunctionDescription message
func (sm serviceModelTest) DecodeRanFunctionDescription(asn1bytes []byte) (*e2smtypes.RanfunctionNameDef, *e2smtypes.RicEventTriggerList, *e2smtypes.RicReportList, error) {
	return nil, nil, nil, fmt.Errorf("not implemented on serviceModelTest")
}

func (sm serviceModelTest) ControlHeaderASN1toProto(asn1Bytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on KPM")
}

func (sm serviceModelTest) ControlHeaderProtoToASN1(protoBytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on KPM")
}

func (sm serviceModelTest) ControlMessageASN1toProto(asn1Bytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on KPM")
}

func (sm serviceModelTest) ControlMessageProtoToASN1(protoBytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on KPM")
}

func (sm serviceModelTest) ControlOutcomeASN1toProto(asn1Bytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on KPM")
}

func (sm serviceModelTest) ControlOutcomeProtoToASN1(protoBytes []byte) ([]byte, error) {
	return nil, fmt.Errorf("not implemented on KPM")
}

type mockRegistry struct {
	plugins map[e2smtypes.OID]ransimmps.ServiceModel
}

// NewModelRegistry create an instance of model registry
func NewMockModelRegistry() ransimmps.ModelRegistry {
	return &mockRegistry{
		plugins: make(map[e2smtypes.OID]ransimmps.ServiceModel),
	}
}

func (r mockRegistry) GetPlugins() map[e2smtypes.OID]ransimmps.ServiceModel {
	return r.plugins
}

func (r mockRegistry) GetPlugin(oid e2smtypes.OID) (ransimmps.ServiceModel, error) {
	p, ok := r.plugins[oid]
	if !ok {
		return nil, fmt.Errorf("unable to find plugin for OID %s", oid)
	}
	return p, nil
}

func (r mockRegistry) RegisterModelPlugin(moduleName string) (e2smtypes.ShortName, e2smtypes.Version, error) {

	r.plugins[smOIDKpmV2] = serviceModelTestInst

	return smName, smVersion, nil
}
