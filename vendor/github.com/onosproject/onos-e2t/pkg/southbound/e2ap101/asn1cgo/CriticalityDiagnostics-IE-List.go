// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "CriticalityDiagnostics-IE-List.h"
//#include "CriticalityDiagnostics-IE-Item.h"
import "C"
import (
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	"unsafe"
)

func newCriticalityDiagnosticsIeList(cdIel *e2apies.CriticalityDiagnosticsIeList) (*C.CriticalityDiagnostics_IE_List_t, error) {
	cdIelC := new(C.CriticalityDiagnostics_IE_List_t)
	for _, cdIe := range cdIel.GetValue() {
		cdIeC, err := newCriticalityDiagnosticsIEItem(cdIe)
		if err != nil {
			return nil, fmt.Errorf("error on newCriticalityDiagnosticsIEItem() %s", err.Error())
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(cdIelC), unsafe.Pointer(cdIeC)); err != nil {
			return nil, err
		}
	}

	return cdIelC, nil
}

func decodeCriticalityDiagnosticsIeList(cdIelC *C.CriticalityDiagnostics_IE_List_t) (*e2apies.CriticalityDiagnosticsIeList, error) {
	cdIel := e2apies.CriticalityDiagnosticsIeList{
		Value: make([]*e2apies.CriticalityDiagnosticsIeItem, 0),
	}

	ieCount := int(cdIelC.list.count)
	//fmt.Printf("RanFunctionListC %T List %T %v Array %T %v Deref %v\n", rflC, rflC.list, rflC.list, rflC.list.array, *rflC.list.array, *(rflC.list.array))
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*cdIelC.list.array)) * uintptr(i)
		cdIeC := *(**C.CriticalityDiagnostics_IE_Item_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cdIelC.list.array)) + offset))
		//fmt.Printf("Value %T %p %v\n", rfiIeC, rfiIeC, rfiIeC)
		cdIe, err := decodeCriticalityDiagnosticsIEItem(cdIeC)
		if err != nil {
			return nil, fmt.Errorf("decodeCriticalityDiagnosticsIeList() %s", err.Error())
		}
		cdIel.Value = append(cdIel.Value, cdIe)
	}

	return &cdIel, nil
}
