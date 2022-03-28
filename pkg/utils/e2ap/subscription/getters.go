// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"github.com/onosproject/onos-lib-go/pkg/errors"

	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
)

// GetRequesterID gets requester ID
func GetRequesterID(request *e2appducontents.RicsubscriptionRequest) (*int32, error) {
	var res int32 = -1
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicrequestID) {
			res = v.GetValue().GetRicrequestId().GetRicRequestorId()
			break
		}
	}

	if res == -1 {
		return nil, errors.NewNotFound("RicRequestID was not found")
	}

	return &res, nil
}

// GetRanFunctionID gets ran function ID
func GetRanFunctionID(request *e2appducontents.RicsubscriptionRequest) (*int32, error) {
	var res int32 = -1
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRanfunctionID) {
			res = v.GetValue().GetRanfunctionId().GetValue()
			break
		}
	}

	if res == -1 {
		return nil, errors.NewNotFound("RanFunctionID was not found")
	}

	return &res, nil
}

// GetRicInstanceID gets ric instance ID
func GetRicInstanceID(request *e2appducontents.RicsubscriptionRequest) (*int32, error) {
	var res int32 = -1
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicrequestID) {
			res = v.GetValue().GetRicrequestId().GetRicInstanceId()
			break
		}
	}

	if res == -1 {
		return nil, errors.NewNotFound("RicInstanceID was not found")
	}

	return &res, nil
}

// GetRicActionToBeSetupList get ric action list
func GetRicActionToBeSetupList(request *e2appducontents.RicsubscriptionRequest) []*e2appducontents.RicactionToBeSetupItemIes {
	var res []*e2appducontents.RicactionToBeSetupItemIes
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			res = v.GetValue().GetRicsubscriptionDetails().GetRicActionToBeSetupList().GetValue()
			break
		}
	}

	return res
}
