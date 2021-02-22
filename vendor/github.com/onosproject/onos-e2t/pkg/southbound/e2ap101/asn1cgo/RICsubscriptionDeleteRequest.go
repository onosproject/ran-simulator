// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICsubscriptionDeleteRequest.h"
//#include "ProtocolIE-Field.h"
import "C"
import e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"

func newRICsubscriptionDeleteRequest(rsr *e2appducontents.RicsubscriptionDeleteRequest) (*C.RICsubscriptionDeleteRequest_t, error) {
	pIeC1710P3, err := newRicSubscriptionDeleteRequestIes(rsr.GetProtocolIes())
	if err != nil {
		return nil, err
	}
	rsrC := C.RICsubscriptionDeleteRequest_t{
		protocolIEs: *pIeC1710P3,
	}

	return &rsrC, nil
}

func decodeRicSubscriptionDeleteRequest(rsdrC *C.RICsubscriptionDeleteRequest_t) (*e2appducontents.RicsubscriptionDeleteRequest, error) {
	pIEs, err := decodeRicSubscriptionDeleteRequestIes(&rsdrC.protocolIEs)
	if err != nil {
		return nil, err
	}

	rsdr := e2appducontents.RicsubscriptionDeleteRequest{
		ProtocolIes: pIEs,
	}
	return &rsdr, nil
}
