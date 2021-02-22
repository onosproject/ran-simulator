// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICsubscriptionDeleteFailure.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newRicSubscriptionDeleteFailure(rsdf *e2appducontents.RicsubscriptionDeleteFailure) (*C.RICsubscriptionDeleteFailure_t, error) {
	pIeC1710P5, err := newRicSubscriptionDeleteFailureIe(rsdf.ProtocolIes)
	if err != nil {
		return nil, err
	}
	rsdfC := C.RICsubscriptionDeleteFailure_t{
		protocolIEs: *pIeC1710P5,
	}

	return &rsdfC, nil
}

func decodeRicSubscriptionDeleteFailure(rsdfC *C.RICsubscriptionDeleteFailure_t) (*e2appducontents.RicsubscriptionDeleteFailure, error) {
	pIEs, err := decodeRicSubscriptionDeleteFailureIes(&rsdfC.protocolIEs)
	if err != nil {
		return nil, err
	}

	rsdf := e2appducontents.RicsubscriptionDeleteFailure{
		ProtocolIes: pIEs,
	}

	return &rsdf, nil
}
