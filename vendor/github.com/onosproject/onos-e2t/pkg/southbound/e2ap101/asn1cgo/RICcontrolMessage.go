// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICcontrolMessage.h"
import "C"
import (
	"encoding/binary"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	"unsafe"
)

func newRicControlMessage(rcm *e2ap_commondatatypes.RiccontrolMessage) *C.RICcontrolMessage_t {
	return newOctetString(string(rcm.GetValue()))
}

func decodeRicControlMessageBytes(rcmBytes []byte) *e2ap_commondatatypes.RiccontrolMessage {
	rcmC := C.OCTET_STRING_t{
		buf:  (*C.uchar)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(rcmBytes[:8])))),
		size: C.ulong(binary.LittleEndian.Uint64(rcmBytes[8:])),
	}
	return decodeRicControlMessage(&rcmC)
}

func decodeRicControlMessage(rcmC *C.RICcontrolMessage_t) *e2ap_commondatatypes.RiccontrolMessage {
	result := e2ap_commondatatypes.RiccontrolMessage{
		Value: []byte(decodeOctetString(rcmC)),
	}

	return &result
}
