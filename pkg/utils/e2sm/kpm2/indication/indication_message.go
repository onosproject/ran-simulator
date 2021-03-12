// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package indication

import (
	e2smkpmies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-ies"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	// "google.golang.org/protobuf/proto"
)

// Message indication message fields for kpm service model
type Message struct {
	numberOfActiveUes int32
	// TODO add remaining files like cu-cp name and rancontainer
}

// NewIndicationMessage creates a new indication message
func NewIndicationMessage(options ...func(message *Message)) *Message {
	msg := &Message{}
	for _, option := range options {
		option(msg)
	}

	return msg
}

// WithNumberOfActiveUes sets number of active UEs
func WithNumberOfActiveUes(numOfActiveUes int32) func(msg *Message) {
	return func(msg *Message) {
		msg.numberOfActiveUes = numOfActiveUes
	}
}

// ToAsn1Bytes converts to Asn1 bytes
func (message *Message) ToAsn1Bytes(modelPlugin modelplugins.ModelPlugin) ([]byte, error) {
	// indicationMessage, err := message.Build()
	// if err != nil {
	// 	return nil, err
	// }
	// indicationMessageProtoBytes, err := proto.Marshal(indicationMessage)
	// if err != nil {
	// 	return nil, err
	// }

	// indicationMessageAsn1Bytes, err := modelPlugin.IndicationMessageProtoToASN1(indicationMessageProtoBytes)
	// if err != nil {
	// 	return nil, err
	// }

	// return indicationMessageAsn1Bytes, nil

	return nil, nil
}

// Build builds indication message for kpm service model
func (message *Message) Build() (*e2smkpmies.E2SmKpmIndicationMessage, error) {
	// e2SmIindicationMsg := e2smkpmies.E2SmKpmIndicationMessage_IndicationMessageFormat1{
	// 	IndicationMessageFormat1: &e2smkpmies.E2SmKpmIndicationMessageFormat1{
	// 		PmContainers: make([]*e2smkpmies.PmContainersList, 0),
	// 	},
	// }

	// ocucpContainer := e2smkpmies.OcucpPfContainer{
	// 	GNbCuCpName: &e2smkpmies.GnbCuCpName{
	// 		Value: "test", //string
	// 	},
	// 	CuCpResourceStatus: &e2smkpmies.OcucpPfContainer_CuCpResourceStatus001{
	// 		NumberOfActiveUes: message.numberOfActiveUes, //int32
	// 	},
	// }

	// containerOcuCp1 := e2smkpmies.PmContainersList{
	// 	PerformanceContainer: &e2smkpmies.PfContainer{
	// 		PfContainer: &e2smkpmies.PfContainer_OCuCp{
	// 			OCuCp: &ocucpContainer,
	// 		},
	// 	},
	// 	TheRancontainer: &e2smkpmies.RanContainer{
	// 		Value: []byte("rancontainer"),
	// 	},
	// }
	// e2SmIindicationMsg.IndicationMessageFormat1.PmContainers = append(e2SmIindicationMsg.IndicationMessageFormat1.PmContainers, &containerOcuCp1)

	// e2smKpmPdu := e2smkpmies.E2SmKpmIndicationMessage{
	// 	E2SmKpmIndicationMessage: &e2SmIindicationMsg,
	// }

	// if err := e2smKpmPdu.Validate(); err != nil {
	// 	return nil, err
	// }
	// return &e2smKpmPdu, nil

	return nil, nil
}
