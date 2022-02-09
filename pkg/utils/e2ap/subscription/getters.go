// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"fmt"
	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
)

// GetRequesterID gets requester ID
func GetRequesterID(request *e2appducontents.RicsubscriptionRequest) (*int32, error) {
	var res int32 = -1
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicrequestID) {
			res = v.GetValue().GetRrId().GetRicRequestorId()
			break
		}
	}

	if res == -1 {
		return nil, fmt.Errorf("RicRequestID was not found")
	}

	return &res, nil
}

// GetRanFunctionID gets ran function ID
func GetRanFunctionID(request *e2appducontents.RicsubscriptionRequest) (*int32, error) {
	var res int32 = -1
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRanfunctionID) {
			res = v.GetValue().GetRfId().GetValue()
			break
		}
	}

	if res == -1 {
		return nil, fmt.Errorf("RanFunctionID was not found")
	}

	return &res, nil
}

// GetRicInstanceID gets ric instance ID
func GetRicInstanceID(request *e2appducontents.RicsubscriptionRequest) (*int32, error) {
	var res int32 = -1
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicrequestID) {
			res = v.GetValue().GetRrId().GetRicInstanceId()
			break
		}
	}

	if res == -1 {
		return nil, fmt.Errorf("RicInstanceID was not found")
	}

	return &res, nil
}

// GetRicActionToBeSetupList get ric action list
func GetRicActionToBeSetupList(request *e2appducontents.RicsubscriptionRequest) []*e2appducontents.RicactionToBeSetupItemIes {
	var res []*e2appducontents.RicactionToBeSetupItemIes
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			res = v.GetValue().GetRsd().GetRicActionToBeSetupList().GetValue()
			break
		}
	}

	return res
}
