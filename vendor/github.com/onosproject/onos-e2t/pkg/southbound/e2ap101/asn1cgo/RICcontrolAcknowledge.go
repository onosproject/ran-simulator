// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICcontrolAcknowledge.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func xerEncodeRICcontrolAcknowledge(rca *e2appducontents.RiccontrolAcknowledge) ([]byte, error) {
	rcaC, err := newRicControlAcknowledge(rca)
	if err != nil {
		return nil, err
	}

	bytes, err := encodeXer(&C.asn_DEF_RICcontrolAcknowledge, unsafe.Pointer(rcaC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func xerDecodeRICcontrolAcknowledge(bytes []byte) (*e2appducontents.RiccontrolAcknowledge, error) {
	unsafePtr, err := decodeXer(bytes, &C.asn_DEF_RICcontrolAcknowledge)
	if err != nil {
		return nil, err
	}
	if unsafePtr == nil {
		return nil, fmt.Errorf("pointer decoded from PER is nil")
	}
	return decodeRicControlAcknowledge((*C.RICcontrolAcknowledge_t)(unsafePtr))
}

func perEncodeRICcontrolAcknowledge(rcr *e2appducontents.RiccontrolAcknowledge) ([]byte, error) {
	rcaC, err := newRicControlAcknowledge(rcr)
	if err != nil {
		return nil, err
	}

	bytes, err := encodePerBuffer(&C.asn_DEF_RICcontrolAcknowledge, unsafe.Pointer(rcaC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func perDecodeRICcontrolAcknowledge(bytes []byte) (*e2appducontents.RiccontrolAcknowledge, error) {
	unsafePtr, err := decodePer(bytes, len(bytes), &C.asn_DEF_RICcontrolAcknowledge)
	if err != nil {
		return nil, err
	}
	if unsafePtr == nil {
		return nil, fmt.Errorf("pointer decoded from PER is nil")
	}
	return decodeRicControlAcknowledge((*C.RICcontrolAcknowledge_t)(unsafePtr))
}

func newRicControlAcknowledge(rca *e2appducontents.RiccontrolAcknowledge) (*C.RICcontrolAcknowledge_t, error) {
	pIeC1710P8, err := newRicControlAcknowledgeIEs(rca.ProtocolIes)
	if err != nil {
		return nil, err
	}
	rcaC := C.RICcontrolAcknowledge_t{
		protocolIEs: *pIeC1710P8,
	}

	return &rcaC, nil
}

func decodeRicControlAcknowledge(rcaC *C.RICcontrolAcknowledge_t) (*e2appducontents.RiccontrolAcknowledge, error) {
	pIEs, err := decodeRicControlAcknowledgeIes(&rcaC.protocolIEs)
	if err != nil {
		return nil, err
	}

	ricControlAcknowledge := e2appducontents.RiccontrolAcknowledge{
		ProtocolIes: pIEs,
	}
	return &ricControlAcknowledge, nil
}
