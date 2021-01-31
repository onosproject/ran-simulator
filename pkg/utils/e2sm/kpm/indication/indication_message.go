// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package indication

import (
	e2sm_kpm_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/v1beta1/e2sm-kpm-ies"
)

// Message indication message fields for kpm service model
type Message struct {
	numberOfActiveUes int32
}

// NewIndicationMessage creates a new indication message
func NewIndicationMessage(options ...func(header *Message)) (*Message, error) {
	msg := &Message{}
	for _, option := range options {
		option(msg)
	}

	return msg, nil
}

// WithNumberOfActiveUes sets number of active UEs
func WithNumberOfActiveUes(numOfActiveUes int32) func(msg *Message) {
	return func(msg *Message) {
		msg.numberOfActiveUes = numOfActiveUes
	}
}

// CreateIndicationMessage creates indication message
func CreateIndicationMessage(indicationMessage *Message) (*e2sm_kpm_ies.E2SmKpmIndicationMessage, error) {
	e2SmIindicationMsg := e2sm_kpm_ies.E2SmKpmIndicationMessage_IndicationMessageFormat1{
		IndicationMessageFormat1: &e2sm_kpm_ies.E2SmKpmIndicationMessageFormat1{
			PmContainers: make([]*e2sm_kpm_ies.PmContainersList, 0),
		},
	}

	ocucpContainer := e2sm_kpm_ies.OcucpPfContainer{
		GNbCuCpName: &e2sm_kpm_ies.GnbCuCpName{
			Value: "test", //string
		},
		CuCpResourceStatus: &e2sm_kpm_ies.OcucpPfContainer_CuCpResourceStatus001{
			NumberOfActiveUes: indicationMessage.numberOfActiveUes, //int32
		},
	}

	containerOcuCp1 := e2sm_kpm_ies.PmContainersList{
		PerformanceContainer: &e2sm_kpm_ies.PfContainer{
			PfContainer: &e2sm_kpm_ies.PfContainer_OCuCp{
				OCuCp: &ocucpContainer,
			},
		},
		TheRancontainer: &e2sm_kpm_ies.RanContainer{
			Value: []byte("rancontainer"),
		},
	}
	e2SmIindicationMsg.IndicationMessageFormat1.PmContainers = append(e2SmIindicationMsg.IndicationMessageFormat1.PmContainers, &containerOcuCp1)

	e2smKpmPdu := e2sm_kpm_ies.E2SmKpmIndicationMessage{
		E2SmKpmIndicationMessage: &e2SmIindicationMsg,
	}

	if err := e2smKpmPdu.Validate(); err != nil {
		return nil, err
	}
	return &e2smKpmPdu, nil
}
