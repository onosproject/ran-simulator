// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

// #cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
// #cgo LDFLAGS: -lm
// #include <stdio.h>
// #include <stdlib.h>
// #include <assert.h>
// #include "RICactions-ToBeSetup-List.h"
//#include "ProtocolIE-SingleContainer.h"
import "C"
import (
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func xerEncodeRicActionsToBeSetupList(ratbsl *e2appducontents.RicactionsToBeSetupList) ([]byte, error) {
	ratbslC, err := newRicActionToBeSetupList(ratbsl)

	if err != nil {
		return nil, err
	}

	bytes, err := encodeXer(&C.asn_DEF_RICactions_ToBeSetup_List, unsafe.Pointer(ratbslC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func perEncodeRicActionsToBeSetupList(ratbsl *e2appducontents.RicactionsToBeSetupList) ([]byte, error) {
	ratbslC, err := newRicActionToBeSetupList(ratbsl)

	if err != nil {
		return nil, err
	}

	bytes, err := encodePerBuffer(&C.asn_DEF_RICactions_ToBeSetup_List, unsafe.Pointer(ratbslC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func xerDecodeRicActionsToBeSetupList(xer []byte) (*e2appducontents.RicactionsToBeSetupList, error) {
	unsafePtr, err := decodeXer(xer, &C.asn_DEF_RICactions_ToBeSetup_List)
	if err != nil {
		return nil, err
	}
	if unsafePtr == nil {
		return nil, fmt.Errorf("pointer decoded from XER is nil")
	}
	ratsslC := (*C.RICactions_ToBeSetup_List_t)(unsafePtr)
	return decodeRicActionToBeSetupList(ratsslC)
}

func newRicActionToBeSetupList(ratbsL *e2appducontents.RicactionsToBeSetupList) (*C.RICactions_ToBeSetup_List_t, error) {
	ratbsLC := new(C.RICactions_ToBeSetup_List_t)
	for _, ricActionToBeSetupItemIe := range ratbsL.GetValue() {
		ricActionToBeItemIesScC, err := newRicActionToBeSetupItemIesSingleContainer(ricActionToBeSetupItemIe)
		if err != nil {
			return nil, fmt.Errorf("newRicActionToBeSetupItemIesSingleContainer() %s", err.Error())
		}

		if _, err = C.asn_sequence_add(unsafe.Pointer(ratbsLC), unsafe.Pointer(ricActionToBeItemIesScC)); err != nil {
			return nil, err
		}
	}

	return ratbsLC, nil
}

func decodeRicActionToBeSetupList(ratbsLC *C.RICactions_ToBeSetup_List_t) (*e2appducontents.RicactionsToBeSetupList, error) {
	ratbsL := e2appducontents.RicactionsToBeSetupList{
		Value: make([]*e2appducontents.RicactionToBeSetupItemIes, 0),
	}

	ieCount := int(ratbsLC.list.count)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*ratbsLC.list.array)) * uintptr(i)
		ratbsIeC := *(**C.ProtocolIE_SingleContainer_1713P0_t)(unsafe.Pointer(uintptr(unsafe.Pointer(ratbsLC.list.array)) + offset))
		ratbsIe, err := decodeRicActionToBeSetupItemIesSingleContainer(ratbsIeC)
		if err != nil {
			return nil, fmt.Errorf("decodeRicActionToBeSetupItemIesSingleContainer() %s", err.Error())
		}
		ratbsL.Value = append(ratbsL.Value, ratbsIe)
	}

	return &ratbsL, nil
}
