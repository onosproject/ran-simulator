// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

// #cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
// #cgo LDFLAGS: -lm
// #include <stdio.h>
// #include <stdlib.h>
// #include <assert.h>
// #include "RICsubscriptionDetails.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func newRicSubscriptionDetails(rsDet *e2appducontents.RicsubscriptionDetails) (*C.RICsubscriptionDetails_t, error) {

	raTbsL, err := newRicActionToBeSetupList(rsDet.RicActionToBeSetupList)
	if err != nil {
		return nil, err
	}

	rsDetC := C.RICsubscriptionDetails_t{
		ricEventTriggerDefinition: *newRicEventTriggerDefinition(rsDet.GetRicEventTriggerDefinition()),
		ricAction_ToBeSetup_List:  *raTbsL,
	}

	return &rsDetC, nil
}

func decodeRicSubscriptionDetailsBytes(bytes []byte) (*e2appducontents.RicsubscriptionDetails, error) {
	array := (**C.struct_ProtocolIE_SingleContainer)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[40:]))))
	count := C.int(binary.LittleEndian.Uint32(bytes[48:]))
	size := C.int(binary.LittleEndian.Uint32(bytes[52:]))

	rsDetC := C.RICsubscriptionDetails_t{
		ricEventTriggerDefinition: C.RICeventTriggerDefinition_t{
			buf:  (*C.uchar)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[:8])))),
			size: C.ulong(binary.LittleEndian.Uint64(bytes[8:])),
		},
		ricAction_ToBeSetup_List: C.RICactions_ToBeSetup_List_t{
			list: C.struct___87{ // TODO: tie this down with a predictable name
				array: array,
				size:  size,
				count: count,
			},
		},
	}

	return decodeRicSubscriptionDetails(&rsDetC)
}

func decodeRicSubscriptionDetails(rsDetC *C.RICsubscriptionDetails_t) (*e2appducontents.RicsubscriptionDetails, error) {
	ratbsL, err := decodeRicActionToBeSetupList(&rsDetC.ricAction_ToBeSetup_List)
	if err != nil {
		return nil, fmt.Errorf("decodeRicActionToBeSetupList() %s", err.Error())
	}

	rsDet := e2appducontents.RicsubscriptionDetails{
		RicEventTriggerDefinition: decodeRicEventTriggerDefinition(&rsDetC.ricEventTriggerDefinition),
		RicActionToBeSetupList:    ratbsL,
	}
	return &rsDet, nil
}
