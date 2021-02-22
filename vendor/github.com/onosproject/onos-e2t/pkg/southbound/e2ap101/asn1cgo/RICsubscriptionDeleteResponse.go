// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICsubscriptionDeleteResponse.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newRicSubscriptionDeleteResponse(rsr *e2appducontents.RicsubscriptionDeleteResponse) (*C.RICsubscriptionDeleteResponse_t, error) {
	pIeC1710P4, err := newRicSubscriptionDeleteResponseIe(rsr.ProtocolIes)
	if err != nil {
		return nil, err
	}
	rsrC := C.RICsubscriptionDeleteResponse_t{
		protocolIEs: *pIeC1710P4,
	}

	return &rsrC, nil
}

func decodeRicSubscriptionDeleteResponse(rsrC *C.RICsubscriptionDeleteResponse_t) (*e2appducontents.RicsubscriptionDeleteResponse, error) {
	pIEs, err := decodeRicSubscriptionDeleteResponseIes(&rsrC.protocolIEs)
	if err != nil {
		return nil, err
	}

	rsr := e2appducontents.RicsubscriptionDeleteResponse{
		ProtocolIes: pIEs,
	}

	return &rsr, nil
}
