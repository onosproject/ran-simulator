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
//extern int consumeBytesCb(void* p0, uint32_t p1, void* p2);
import "C"
import (
	"fmt"
	"math/rand"
	"sync"
	"unsafe"
)

var cbChan = make(chan []byte)
var response []byte
var chanMutex = sync.RWMutex{}

// Callback function of type C.asn_app_consume_bytes_f
//export consumeBytesCb
func consumeBytesCb(buf unsafe.Pointer, size C.uint32_t, key unsafe.Pointer) C.int {
	bytes := C.GoBytes(buf, C.int(size))
	chanMutex.Lock()
	cbChan <- bytes
	return C.int(size)
}

func encodeXer(valueType *C.asn_TYPE_descriptor_t,
	value unsafe.Pointer) ([]byte, error) {
	key := int(rand.Int31n(1e3))
	keyCint := C.int(key)
	response = nil
	// bytes get pushed back through this channel
	go func() {
		for d := range cbChan {
			//chanMutex.Lock()
			response = append(response, d...)
			chanMutex.Unlock()
		}
	}()

	xerCbF := (*C.asn_app_consume_bytes_f)(C.consumeBytesCb)
	encRetVal, err := C.xer_encode(valueType, value,
		C.XER_F_BASIC, xerCbF, unsafe.Pointer(&keyCint))
	if err != nil {
		return nil, err
	}
	if encRetVal.encoded == -1 {
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
	chanMutex.RLock()
	// Don't exit until unlocked in go routine
	defer chanMutex.RUnlock()
	return response[:], err
}
