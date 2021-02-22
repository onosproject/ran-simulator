// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "ErrorIndication.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newErrorIndication(ei *e2appducontents.ErrorIndication) (*C.ErrorIndication_t, error) {
	pIeC1710P10, err := newErrorIndicationIe(ei.ProtocolIes)
	if err != nil {
		return nil, err
	}
	eiC := C.ErrorIndication_t{
		protocolIEs: *pIeC1710P10,
	}

	return &eiC, nil
}

func decodeErrorIndication(eiC *C.ErrorIndication_t) (*e2appducontents.ErrorIndication, error) {
	pIEs, err := decodeErrorIndicationIes(&eiC.protocolIEs)
	if err != nil {
		return nil, err
	}

	ei := e2appducontents.ErrorIndication{
		ProtocolIes: pIEs,
	}

	return &ei, nil
}
