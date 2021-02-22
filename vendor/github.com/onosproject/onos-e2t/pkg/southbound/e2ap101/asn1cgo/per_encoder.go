// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "E2AP-PDU.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func encodePerBuffer(valueType *C.asn_TYPE_descriptor_t,
	value unsafe.Pointer) ([]byte, error) {

	perBuf := C.malloc(C.sizeof_uchar * 1024) // C allocated pointer
	defer C.free(perBuf)

	encRetVal, err := C.asn_encode_to_buffer(nil, C.ATS_ALIGNED_BASIC_PER, valueType, value, perBuf, C.ulong(1024))
	if err != nil {
		return nil, err
	}
	if encRetVal.encoded == -1 {
		//fmt.Printf("error on %v\n", *encRetVal.failed_type)
		var i C.uint
		for i = 0; i < encRetVal.failed_type.elements_count; i++ {
			step := C.uint(unsafe.Sizeof(C.asn_TYPE_member_t{})) * i
			element := (*C.asn_TYPE_member_t)(unsafe.Pointer(uintptr(unsafe.Pointer(encRetVal.failed_type.elements)) + uintptr(step)))
			fmt.Printf("Element %v %v\n", i, C.GoString(element.name))
		}

		return nil, fmt.Errorf("error encoding. Name: %v Tag: %v #Tags: %v Alltags: %v, Elements: %v",
			C.GoString(encRetVal.failed_type.name),
			C.GoString(encRetVal.failed_type.xml_tag),
			encRetVal.failed_type.tags_count,
			encRetVal.failed_type.all_tags_count,
			encRetVal.failed_type.elements_count)
	}
	bytes := make([]byte, encRetVal.encoded)
	for i := 0; i < int(encRetVal.encoded); i++ {
		b := *(*C.uchar)(unsafe.Pointer(uintptr(perBuf) + uintptr(i)))
		bytes[i] = byte(b)
	}
	return bytes, nil
}
