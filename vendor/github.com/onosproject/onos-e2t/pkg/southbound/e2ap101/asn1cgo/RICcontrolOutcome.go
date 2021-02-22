// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICcontrolOutcome.h"
import "C"
import (
	"encoding/binary"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	"unsafe"
)

func newRicControlOutcome(rco *e2ap_commondatatypes.RiccontrolOutcome) *C.RICcontrolOutcome_t {
	return newOctetString(string(rco.GetValue()))
}

func decodeRicControlOutcomeBytes(rcoBytes []byte) *e2ap_commondatatypes.RiccontrolOutcome {
	rcmC := C.OCTET_STRING_t{
		buf:  (*C.uchar)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(rcoBytes[:8])))),
		size: C.ulong(binary.LittleEndian.Uint64(rcoBytes[8:])),
	}
	return decodeRicControlOutcome(&rcmC)
}

func decodeRicControlOutcome(rcoC *C.RICcontrolOutcome_t) *e2ap_commondatatypes.RiccontrolOutcome {
	result := e2ap_commondatatypes.RiccontrolOutcome{
		Value: []byte(decodeOctetString(rcoC)),
	}

	return &result
}
