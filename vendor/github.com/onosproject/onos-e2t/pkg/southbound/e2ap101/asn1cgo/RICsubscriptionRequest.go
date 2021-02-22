// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICsubscriptionRequest.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func xerEncodeRICsubscriptionRequest(rsr *e2appducontents.RicsubscriptionRequest) ([]byte, error) {
	rsrC, err := newRICsubscriptionRequest(rsr)
	if err != nil {
		return nil, err
	}

	bytes, err := encodeXer(&C.asn_DEF_RICsubscriptionRequest, unsafe.Pointer(rsrC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func perEncodeRICsubscriptionRequest(rsr *e2appducontents.RicsubscriptionRequest) ([]byte, error) {
	rsrC, err := newRICsubscriptionRequest(rsr)
	if err != nil {
		return nil, err
	}

	bytes, err := encodePerBuffer(&C.asn_DEF_RICsubscriptionRequest, unsafe.Pointer(rsrC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func newRICsubscriptionRequest(rsr *e2appducontents.RicsubscriptionRequest) (*C.RICsubscriptionRequest_t, error) {
	pIeC1710P0, err := newRicSubscriptionRequestIes(rsr.GetProtocolIes())
	if err != nil {
		return nil, err
	}
	rsrC := C.RICsubscriptionRequest_t{
		protocolIEs: *pIeC1710P0,
	}

	return &rsrC, nil
}

func decodeRicSubscriptionRequest(ricSubscriptionRequestC *C.RICsubscriptionRequest_t) (*e2appducontents.RicsubscriptionRequest, error) {
	pIEs, err := decodeRicSubscriptionRequestIes(&ricSubscriptionRequestC.protocolIEs)
	if err != nil {
		return nil, err
	}

	ricSubscriptionRequest := e2appducontents.RicsubscriptionRequest{
		ProtocolIes: pIEs,
	}
	return &ricSubscriptionRequest, nil
}
