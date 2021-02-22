// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

// #cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
// #cgo LDFLAGS: -lm
// #include <stdio.h>
// #include <stdlib.h>
// #include <assert.h>
// #include "RICrequestID.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func xerDecodeRicRequestID(bytes []byte) (*e2apies.RicrequestId, error) {
	unsafePtr, err := decodeXer(bytes, &C.asn_DEF_RICrequestID)
	if err != nil {
		return nil, err
	}
	if unsafePtr == nil {
		return nil, fmt.Errorf("pointer decoded from XER is nil")
	}
	ricIndicationC := (*C.RICrequestID_t)(unsafePtr)
	ricIndication := decodeRicRequestID(ricIndicationC)

	return ricIndication, nil
}

func newRicRequestID(rrID *e2apies.RicrequestId) *C.RICrequestID_t {
	rrIDC := C.RICrequestID_t{
		ricRequestorID: C.long(rrID.RicRequestorId),
		ricInstanceID:  C.long(rrID.RicInstanceId),
	}
	return &rrIDC
}

func decodeRicRequestIDBytes(ricRequestIDCchoice []byte) *e2apies.RicrequestId {
	ricRequestorID := binary.LittleEndian.Uint64(ricRequestIDCchoice[0:8])
	ricInstanceID := binary.LittleEndian.Uint64(ricRequestIDCchoice[8:16])

	rrID := C.RICrequestID_t{
		ricRequestorID: C.long(ricRequestorID),
		ricInstanceID:  C.long(ricInstanceID),
	}

	return decodeRicRequestID(&rrID)
}

func decodeRicRequestID(ricRequestIDCchoice *C.RICrequestID_t) *e2apies.RicrequestId {

	result := e2apies.RicrequestId{
		RicRequestorId: int32(ricRequestIDCchoice.ricRequestorID),
		RicInstanceId:  int32(ricRequestIDCchoice.ricInstanceID),
	}

	return &result
}
