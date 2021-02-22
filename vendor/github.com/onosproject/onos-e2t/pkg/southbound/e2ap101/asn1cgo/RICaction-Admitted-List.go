// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICaction-Admitted-List.h"
//#include "ProtocolIE-SingleContainer.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func newRicActionAdmittedList(raal *e2appducontents.RicactionAdmittedList) (*C.RICaction_Admitted_List_t, error) {
	rfIDlC := new(C.RICaction_Admitted_List_t)

	for _, raaID := range raal.GetValue() {
		rfIDC, err := newRicActionAdmittedItemIEItemIesSingleContainer(raaID)
		if err != nil {
			return nil, fmt.Errorf("error on newRicActionAdmittedItemIEItemIesSingleContainer() %s", err.Error())
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(rfIDlC), unsafe.Pointer(rfIDC)); err != nil {
			return nil, err
		}
	}

	return rfIDlC, nil
}

func decodeRicActionAdmittedListBytes(raalBytes []byte) (*e2appducontents.RicactionAdmittedList, error) {
	array := (**C.struct_ProtocolIE_SingleContainer)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(raalBytes[0:8]))))
	count := C.int(binary.LittleEndian.Uint32(raalBytes[8:12]))
	size := C.int(binary.LittleEndian.Uint32(raalBytes[12:16]))

	raalC := C.RICaction_Admitted_List_t{
		list: C.struct___95{
			array: array,
			size:  size,
			count: count,
		},
	}

	return decodeRicActionAdmittedList(&raalC)
}

func decodeRicActionAdmittedList(raalC *C.RICaction_Admitted_List_t) (*e2appducontents.RicactionAdmittedList, error) {

	raal := e2appducontents.RicactionAdmittedList{
		Value: make([]*e2appducontents.RicactionAdmittedItemIes, 0),
	}

	ieCount := int(raalC.list.count)
	//fmt.Printf("RicactionAdmittedList %T List %T %v Array %T %v Deref %v\n", rflC, rflC.list, rflC.list, rflC.list.array, *rflC.list.array, *(rflC.list.array))
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*raalC.list.array)) * uintptr(i)
		rfIDiIeC := *(**C.ProtocolIE_SingleContainer_1713P1_t)(unsafe.Pointer(uintptr(unsafe.Pointer(raalC.list.array)) + offset))
		//fmt.Printf("Value %T %p %v\n", rfIDiIeC, rfIDiIeC, rfIDiIeC)
		rfIDiIe, err := decodeRicActionAdmittedItemIesSingleContainer(rfIDiIeC)
		if err != nil {
			return nil, fmt.Errorf("decodeRicActionAdmittedItemIesSingleContainer() %s", err.Error())
		}
		raal.Value = append(raal.Value, rfIDiIe)
	}

	return &raal, nil
}
