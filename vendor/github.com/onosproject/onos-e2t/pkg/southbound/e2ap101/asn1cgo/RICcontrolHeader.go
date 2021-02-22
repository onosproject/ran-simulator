// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICcontrolHeader.h"
import "C"
import (
	"encoding/binary"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	"unsafe"
)

func newRicControlHeader(rch *e2ap_commondatatypes.RiccontrolHeader) *C.RICcontrolHeader_t {
	return newOctetString(string(rch.GetValue()))
}

func decodeRicControlHeaderBytes(rchBytes []byte) *e2ap_commondatatypes.RiccontrolHeader {
	rchC := C.OCTET_STRING_t{
		buf:  (*C.uchar)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(rchBytes[:8])))),
		size: C.ulong(binary.LittleEndian.Uint64(rchBytes[8:])),
	}
	return decodeRicControlHeader(&rchC)
}

func decodeRicControlHeader(rchC *C.RICcontrolHeader_t) *e2ap_commondatatypes.RiccontrolHeader {
	result := e2ap_commondatatypes.RiccontrolHeader{
		Value: []byte(decodeOctetString(rchC)),
	}

	return &result
}
