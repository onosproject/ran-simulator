// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RANfunctions-List.h"
//#include "ProtocolIE-SingleContainer.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func xerEncodeRanFunctionsList(rfl *e2appducontents.RanfunctionsList) ([]byte, error) {
	rflC, err := newRanFunctionsList(rfl)

	if err != nil {
		return nil, err
	}

	bytes, err := encodeXer(&C.asn_DEF_RANfunctions_List, unsafe.Pointer(rflC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func perEncodeRanFunctionsList(rfl *e2appducontents.RanfunctionsList) ([]byte, error) {
	rflC, err := newRanFunctionsList(rfl)

	if err != nil {
		return nil, err
	}

	bytes, err := encodePerBuffer(&C.asn_DEF_RANfunctions_List, unsafe.Pointer(rflC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func xerDecodeRanFunctionList(xer []byte) (*e2appducontents.RanfunctionsList, error) {
	unsafePtr, err := decodeXer(xer, &C.asn_DEF_RANfunctions_List)
	if err != nil {
		return nil, err
	}
	if unsafePtr == nil {
		return nil, fmt.Errorf("pointer decoded from XER is nil")
	}
	rflC := (*C.RANfunctions_List_t)(unsafePtr)
	return decodeRanFunctionsList(rflC)
}

func perDecodeRanFunctionList(per []byte) (*e2appducontents.RanfunctionsList, error) {
	unsafePtr, err := decodePer(per, len(per), &C.asn_DEF_RANfunctions_List)
	if err != nil {
		return nil, err
	}
	if unsafePtr == nil {
		return nil, fmt.Errorf("pointer decoded from XER is nil")
	}
	rflC := (*C.RANfunctions_List_t)(unsafePtr)
	return decodeRanFunctionsList(rflC)
}

func newRanFunctionsList(rfl *e2appducontents.RanfunctionsList) (*C.RANfunctions_List_t, error) {
	rflC := new(C.RANfunctions_List_t)
	for _, ranfunctionItemIe := range rfl.GetValue() {
		ranfunctionItemIesScC, err := newRanFunctionItemIesSingleContainer(ranfunctionItemIe)
		if err != nil {
			return nil, fmt.Errorf("newRanFunctionsList() %s", err.Error())
		}

		if _, err = C.asn_sequence_add(unsafe.Pointer(rflC), unsafe.Pointer(ranfunctionItemIesScC)); err != nil {
			return nil, err
		}
	}

	return rflC, nil
}

func decodeRanFunctionsListBytes(ranFunctionListChoice [48]byte) (*e2appducontents.RanfunctionsList, error) {
	array := (**C.struct_ProtocolIE_SingleContainer)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(ranFunctionListChoice[0:8]))))
	count := C.int(binary.LittleEndian.Uint32(ranFunctionListChoice[8:12]))
	size := C.int(binary.LittleEndian.Uint32(ranFunctionListChoice[12:16]))

	ranFunctionListChoiceC := C.RANfunctions_List_t{
		list: C.struct___69{
			array: array,
			size:  size,
			count: count,
		},
	}

	return decodeRanFunctionsList(&ranFunctionListChoiceC)
}

func decodeRanFunctionsList(rflC *C.RANfunctions_List_t) (*e2appducontents.RanfunctionsList, error) {
	rfl := e2appducontents.RanfunctionsList{
		Value: make([]*e2appducontents.RanfunctionItemIes, 0),
	}

	ieCount := int(rflC.list.count)
	//fmt.Printf("RanFunctionListC %T List %T %v Array %T %v Deref %v\n", rflC, rflC.list, rflC.list, rflC.list.array, *rflC.list.array, *(rflC.list.array))
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*rflC.list.array)) * uintptr(i)
		rfiIeC := *(**C.ProtocolIE_SingleContainer_1713P8_t)(unsafe.Pointer(uintptr(unsafe.Pointer(rflC.list.array)) + offset))
		//fmt.Printf("Value %T %p %v\n", rfiIeC, rfiIeC, rfiIeC)
		rfiIe, err := decodeRanFunctionItemIesSingleContainer(rfiIeC)
		if err != nil {
			return nil, fmt.Errorf("decodeRanFunctionsList() %s", err.Error())
		}
		rfl.Value = append(rfl.Value, rfiIe)
	}

	return &rfl, nil
}
