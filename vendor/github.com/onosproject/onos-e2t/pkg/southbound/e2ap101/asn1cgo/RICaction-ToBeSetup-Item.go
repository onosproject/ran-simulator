// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICaction-ToBeSetup-Item.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func newRicActionToBeSetupItem(ratbsItem *e2appducontents.RicactionToBeSetupItem) (*C.RICaction_ToBeSetup_Item_t, error) {
	ratC, err := newRicActionType(ratbsItem.GetRicActionType())
	if err != nil {
		return nil, fmt.Errorf("newRicActionType() %s", err.Error())
	}

	rsaC, err := newRicSubsequentAction(ratbsItem.RicSubsequentAction)
	if err != nil {
		return nil, fmt.Errorf("newRicSubsequentAction() %s", err.Error())
	}

	ratbsItemC := C.RICaction_ToBeSetup_Item_t{
		ricActionID:         *newRicActionID(ratbsItem.GetRicActionId()),
		ricActionType:       *ratC,
		ricActionDefinition: newRicActionDefinition(ratbsItem.GetRicActionDefinition()),
		ricSubsequentAction: rsaC,
	}

	return &ratbsItemC, nil
}

func decodeRicActionToBeSetupItemBytes(bytes [56]byte) (*e2appducontents.RicactionToBeSetupItem, error) {

	rfiC := C.RICaction_ToBeSetup_Item_t{
		ricActionID:         C.long(binary.LittleEndian.Uint64(bytes[:8])),
		ricActionType:       C.long(binary.LittleEndian.Uint64(bytes[8:16])),
		ricActionDefinition: (*C.RICactionDefinition_t)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[16:24])))),
		ricSubsequentAction: (*C.struct_RICsubsequentAction)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[24:32])))),
	}

	return decodeRicActionToBeSetupItem(&rfiC)
}

func decodeRicActionToBeSetupItem(rfiC *C.RICaction_ToBeSetup_Item_t) (*e2appducontents.RicactionToBeSetupItem, error) {
	rsa, err := decodeRicSubsequentAction(rfiC.ricSubsequentAction)
	if err != nil {
		return nil, fmt.Errorf("decodeRicSubsequentAction() %s", err.Error())
	}

	rfi := e2appducontents.RicactionToBeSetupItem{
		RicActionId:         decodeRicActionID(&rfiC.ricActionID),
		RicActionType:       decodeRicActionType(&rfiC.ricActionType),
		RicActionDefinition: decodeRicActionDefinition(rfiC.ricActionDefinition),
		RicSubsequentAction: rsa,
	}

	return &rfi, nil
}
