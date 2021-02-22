// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICaction-NotAdmitted-List.h"
//#include "ProtocolIE-SingleContainer.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func newRicActionNotAdmittedList(ranaL *e2appducontents.RicactionNotAdmittedList) (*C.RICaction_NotAdmitted_List_t, error) {
	ranaLC := new(C.RICaction_NotAdmitted_List_t)

	for _, rana := range ranaL.GetValue() {
		ranaC, err := newRicActionNotAdmittedItemIEItemIesSingleContainer(rana)
		if err != nil {
			return nil, fmt.Errorf("error on newRicActionNotAdmittedItemIEItemIesSingleContainer() %s", err.Error())
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(ranaLC), unsafe.Pointer(ranaC)); err != nil {
			return nil, err
		}
	}

	return ranaLC, nil
}

func decodeRicActionNotAdmittedListBytes(ranaLBytes []byte) (*e2appducontents.RicactionNotAdmittedList, error) {
	array := (**C.struct_ProtocolIE_SingleContainer)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(ranaLBytes[0:8]))))
	count := C.int(binary.LittleEndian.Uint32(ranaLBytes[8:12]))
	size := C.int(binary.LittleEndian.Uint32(ranaLBytes[12:16]))

	ranaLC := C.RICaction_NotAdmitted_List_t{
		list: C.struct___108{
			array: array,
			size:  size,
			count: count,
		},
	}

	return decodeRicActionNotAdmittedList(&ranaLC)
}

func decodeRicActionNotAdmittedList(ranaLC *C.RICaction_NotAdmitted_List_t) (*e2appducontents.RicactionNotAdmittedList, error) {

	ranaL := e2appducontents.RicactionNotAdmittedList{
		Value: make([]*e2appducontents.RicactionNotAdmittedItemIes, 0),
	}

	ieCount := int(ranaLC.list.count)
	//fmt.Printf("RicactionAdmittedList %T List %T %v Array %T %v Deref %v\n", rflC, rflC.list, rflC.list, rflC.list.array, *rflC.list.array, *(rflC.list.array))
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*ranaLC.list.array)) * uintptr(i)
		rfIDiIeC := *(**C.ProtocolIE_SingleContainer_1713P2_t)(unsafe.Pointer(uintptr(unsafe.Pointer(ranaLC.list.array)) + offset))
		//fmt.Printf("Value %T %p %v\n", rfIDiIeC, rfIDiIeC, rfIDiIeC)
		rfIDiIe, err := decodeRicActionNotAdmittedItemIesSingleContainer(rfIDiIeC)
		if err != nil {
			return nil, fmt.Errorf("decodeRicActionNotAdmittedItemIesSingleContainer() %s", err.Error())
		}
		ranaL.Value = append(ranaL.Value, rfIDiIe)
	}

	return &ranaL, nil
}
