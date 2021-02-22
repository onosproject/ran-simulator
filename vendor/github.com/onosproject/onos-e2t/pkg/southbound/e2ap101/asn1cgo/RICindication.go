// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICindication.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newRicIndication(ri *e2appducontents.Ricindication) (*C.RICindication_t, error) {
	pIeC1710P6, err := newRicIndicationIEs(ri.ProtocolIes)
	if err != nil {
		return nil, err
	}
	riC := C.RICindication_t{
		protocolIEs: *pIeC1710P6,
	}

	return &riC, nil
}

func decodeRicIndication(riC *C.RICindication_t) (*e2appducontents.Ricindication, error) {
	pIEs, err := decodeRicIndicationIes(&riC.protocolIEs)
	if err != nil {
		return nil, err
	}

	ricIndication := e2appducontents.Ricindication{
		ProtocolIes: pIEs,
	}
	return &ricIndication, nil
}
