// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package utils

import (
	e2tapi "github.com/onosproject/onos-api/go/onos/e2t/e2"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/pdubuilder"
	e2sm_rc_pre_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
	"google.golang.org/protobuf/proto"
)

// Control control request fields
type Control struct {
	NodeID              string
	ServiceModelName    e2tapi.ServiceModelName
	ServiceModelVersion e2tapi.ServiceModelVersion
	ControlMessage      []byte
	ControlHeader       []byte
	ControlAckRequest   e2tapi.ControlAckRequest
	EncodingType        e2tapi.EncodingType
}

// RcControlHeader required fields for creating RC service model control header
type RcControlHeader struct {
	Priority int32
	PlmnID   []byte
	CellID   uint64
}

// RcControlMessage required fields for creating RC service model control message
type RcControlMessage struct {
	RanParameterID    int32
	RanParameterName  string
	RanParameterValue int32
}

// CreatRcControlHeader creates rc control header
func (ch *RcControlHeader) CreateRcControlHeader() ([]byte, error) {
	cellID := &e2sm_rc_pre_ies.BitString{
		Value: ch.CellID,
		Len:   28,
	}
	newE2SmRcPrePdu, err := pdubuilder.CreateE2SmRcPreControlHeader(ch.Priority, ch.PlmnID, cellID)
	if err != nil {
		return []byte{}, err
	}

	err = newE2SmRcPrePdu.Validate()
	if err != nil {
		return []byte{}, err
	}

	protoBytes, err := proto.Marshal(newE2SmRcPrePdu)
	if err != nil {
		return []byte{}, err
	}

	return protoBytes, nil
}

// CreateRcControlMessage creates rc control message
func (cm *RcControlMessage) CreateRcControlMessage() ([]byte, error) {
	newE2SmRcPrePdu, err := pdubuilder.CreateE2SmRcPreControlMessage(cm.RanParameterID, cm.RanParameterName, cm.RanParameterValue)
	if err != nil {
		return []byte{}, err
	}
	err = newE2SmRcPrePdu.Validate()
	if err != nil {
		return []byte{}, err
	}

	protoBytes, err := proto.Marshal(newE2SmRcPrePdu)
	if err != nil {
		return []byte{}, err
	}

	return protoBytes, nil
}

// Create creates a control request using SDK
func (cr *Control) Create() (*e2tapi.ControlRequest, error) {
	request := &e2tapi.ControlRequest{
		E2NodeID: e2tapi.E2NodeID(cr.NodeID),
		Header: &e2tapi.RequestHeader{
			EncodingType: cr.EncodingType,
			ServiceModel: &e2tapi.ServiceModel{
				Name:    cr.ServiceModelName,
				Version: cr.ServiceModelVersion,
			},
		},
		ControlAckRequest: cr.ControlAckRequest,
		ControlHeader:     cr.ControlHeader,
		ControlMessage:    cr.ControlMessage,
	}

	return request, nil

}
