// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

// #cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
// #cgo LDFLAGS: -lm
// #include <stdio.h>
// #include <stdlib.h>
// #include <assert.h>
// #include "E2AP-PDU.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func decodeXer(bytes []byte, valueType *C.asn_TYPE_descriptor_t) (unsafe.Pointer, error) {

	var result unsafe.Pointer
	fmt.Printf("Decode XER %d\n", len(bytes))
	decRetVal, err := C.xer_decode(nil, valueType, &result, C.CBytes(bytes), C.ulong(len(bytes)))
	if err != nil {
		return nil, err
	}
	switch decRetVal.code {
	case C.RC_OK, C.RC_WMORE:
		return result, nil
	//case C.RC_WMORE:
	//	return nil, fmt.Errorf("unhandled - want more. Consumed %v", decRetVal.consumed)
	case C.RC_FAIL:
		return nil, fmt.Errorf("failed to decode. Consumed %v", decRetVal.consumed)
	default:
		return nil, fmt.Errorf("unexpected return code %v", decRetVal.code)
	}
}
