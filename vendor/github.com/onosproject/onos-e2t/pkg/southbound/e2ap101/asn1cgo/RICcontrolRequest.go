// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICcontrolRequest.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func xerEncodeRICcontrolRequest(rcr *e2appducontents.RiccontrolRequest) ([]byte, error) {
	rcrC, err := newRicControlRequest(rcr)
	if err != nil {
		return nil, err
	}

	bytes, err := encodeXer(&C.asn_DEF_RICcontrolRequest, unsafe.Pointer(rcrC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func xerDecodeRICcontrolRequest(bytes []byte) (*e2appducontents.RiccontrolRequest, error) {
	unsafePtr, err := decodeXer(bytes, &C.asn_DEF_RICcontrolRequest)
	if err != nil {
		return nil, err
	}
	if unsafePtr == nil {
		return nil, fmt.Errorf("pointer decoded from PER is nil")
	}
	return decodeRicControlRequest((*C.RICcontrolRequest_t)(unsafePtr))
}

func perEncodeRICcontrolRequest(rcr *e2appducontents.RiccontrolRequest) ([]byte, error) {
	rcrC, err := newRicControlRequest(rcr)
	if err != nil {
		return nil, err
	}

	bytes, err := encodePerBuffer(&C.asn_DEF_RICcontrolRequest, unsafe.Pointer(rcrC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func perDecodeRICcontrolRequest(bytes []byte) (*e2appducontents.RiccontrolRequest, error) {
	unsafePtr, err := decodePer(bytes, len(bytes), &C.asn_DEF_RICcontrolRequest)
	if err != nil {
		return nil, err
	}
	if unsafePtr == nil {
		return nil, fmt.Errorf("pointer decoded from PER is nil")
	}
	return decodeRicControlRequest((*C.RICcontrolRequest_t)(unsafePtr))
}

func newRicControlRequest(rcr *e2appducontents.RiccontrolRequest) (*C.RICcontrolRequest_t, error) {
	pIeC1710P7, err := newRicControlRequestIEs(rcr.ProtocolIes)
	if err != nil {
		return nil, err
	}
	rcrC := C.RICcontrolRequest_t{
		protocolIEs: *pIeC1710P7,
	}

	return &rcrC, nil
}

func decodeRicControlRequest(rcrC *C.RICcontrolRequest_t) (*e2appducontents.RiccontrolRequest, error) {
	pIEs, err := decodeRicControlRequestIes(&rcrC.protocolIEs)
	if err != nil {
		return nil, err
	}

	ricControlRequest := e2appducontents.RiccontrolRequest{
		ProtocolIes: pIEs,
	}
	return &ricControlRequest, nil
}
