// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package control

import (
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

// GetRequesterID gets requester ID
func GetRequesterID(request *e2appducontents.RiccontrolRequest) int32 {
	return request.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId
}

// GetRanFunctionID gets ran function ID
func GetRanFunctionID(request *e2appducontents.RiccontrolRequest) int32 {
	return request.ProtocolIes.E2ApProtocolIes5.Value.Value
}

// GetRicInstanceID gets ric instance ID
func GetRicInstanceID(request *e2appducontents.RiccontrolRequest) int32 {
	return request.ProtocolIes.E2ApProtocolIes29.Value.RicInstanceId
}
