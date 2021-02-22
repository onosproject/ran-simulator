// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "ProtocolIE-Field.h"
//#include "ProtocolIE-SingleContainer.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func newRanFunctionsIDcauseList(rfIDcl *e2appducontents.RanfunctionsIdcauseList) (*C.RANfunctionsIDcause_List_t, error) {
	rfIDclC := new(C.RANfunctionsIDcause_List_t)
	for _, rfIDCause := range rfIDcl.GetValue() {
		rfIDcauseC, err := newRanFunctionIDcauseItemIesSingleContainer(rfIDCause)
		if err != nil {
			return nil, fmt.Errorf("error on newRanFunctionIDcauseItemIesSingleContainer() %s", err.Error())
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(rfIDclC), unsafe.Pointer(rfIDcauseC)); err != nil {
			return nil, err
		}
	}

	return rfIDclC, nil
}

func decodeRanFunctionsIDCauseListBytes(ranFunctionIDCauseListChoice [112]byte) (*e2appducontents.RanfunctionsIdcauseList, error) {
	array := (**C.struct_ProtocolIE_SingleContainer)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(ranFunctionIDCauseListChoice[0:8]))))
	count := C.int(binary.LittleEndian.Uint32(ranFunctionIDCauseListChoice[8:12]))
	size := C.int(binary.LittleEndian.Uint32(ranFunctionIDCauseListChoice[12:16]))

	rfIDCauselC := C.RANfunctionsIDcause_List_t{
		list: C.struct___108{
			array: array,
			size:  size,
			count: count,
		},
	}

	return decodeRanFunctionsCauseIDList(&rfIDCauselC)
}

func decodeRanFunctionsCauseIDList(rfIDCauselC *C.RANfunctionsIDcause_List_t) (*e2appducontents.RanfunctionsIdcauseList, error) {
	rfIDcausel := e2appducontents.RanfunctionsIdcauseList{
		Value: make([]*e2appducontents.RanfunctionIdcauseItemIes, 0),
	}

	ieCount := int(rfIDCauselC.list.count)
	//fmt.Printf("RanFunctionIDListC %T List %T %v Array %T %v Deref %v\n", rflC, rflC.list, rflC.list, rflC.list.array, *rflC.list.array, *(rflC.list.array))
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*rfIDCauselC.list.array)) * uintptr(i)
		rfIDciIeC := *(**C.ProtocolIE_SingleContainer_1713P10_t)(unsafe.Pointer(uintptr(unsafe.Pointer(rfIDCauselC.list.array)) + offset))
		//fmt.Printf("Value %T %p %v\n", rfIDciIeC, rfIDciIeC, rfIDciIeC)
		rfIDiIe, err := decodeRanFunctionIDCauseItemIesSingleContainer(rfIDciIeC)
		if err != nil {
			return nil, fmt.Errorf("decodeRanFunctionIDCauseItemIesSingleContainer() %s", err.Error())
		}
		rfIDcausel.Value = append(rfIDcausel.Value, rfIDiIe)
	}

	return &rfIDcausel, nil
}
