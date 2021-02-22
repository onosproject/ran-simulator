// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "E2setupFailure.h"
//#include "ProtocolIE-Field.h"
import "C"
import e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"

func newE2setupFailure(ei *e2appducontents.E2SetupFailure) (*C.E2setupFailure_t, error) {
	pIeC1710P13, err := newE2setupFailureIe(ei.ProtocolIes)
	if err != nil {
		return nil, err
	}
	e2sfC := C.E2setupFailure_t{
		protocolIEs: *pIeC1710P13,
	}

	return &e2sfC, nil
}

func decodeE2setupFailure(eiC *C.E2setupFailure_t) (*e2appducontents.E2SetupFailure, error) {
	pIEs, err := decodeE2setupFailureIes(&eiC.protocolIEs)
	if err != nil {
		return nil, err
	}

	e2sf := e2appducontents.E2SetupFailure{
		ProtocolIes: pIEs,
	}

	return &e2sf, nil
}
