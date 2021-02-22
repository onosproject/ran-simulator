// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICsubscriptionFailure.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newRicSubscriptionFailure(rsf *e2appducontents.RicsubscriptionFailure) (*C.RICsubscriptionFailure_t, error) {
	pIeC1710P2, err := newRicSubscriptionFailureIe(rsf.ProtocolIes)
	if err != nil {
		return nil, err
	}
	rsfC := C.RICsubscriptionFailure_t{
		protocolIEs: *pIeC1710P2,
	}

	return &rsfC, nil
}

func decodeRicSubscriptionFailure(rsfC *C.RICsubscriptionFailure_t) (*e2appducontents.RicsubscriptionFailure, error) {
	pIEs, err := decodeRicSubscriptionFailureIes(&rsfC.protocolIEs)
	if err != nil {
		return nil, err
	}

	rsf := e2appducontents.RicsubscriptionFailure{
		ProtocolIes: pIEs,
	}

	return &rsf, nil
}
