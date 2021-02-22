// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "GlobalE2node-ID.h"
//#include "GlobalE2node-gNB-ID.h"
//#include "GlobalE2node-eNB-ID.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	"unsafe"
)

func newGlobalE2nodeID(gnID *e2apies.GlobalE2NodeId) (*C.GlobalE2node_ID_t, error) {
	var prC C.GlobalE2node_ID_PR

	choiceC := [8]byte{} // The size of the GlobalE2node_ID_u
	switch choice := gnID.GetGlobalE2NodeId().(type) {
	case *e2apies.GlobalE2NodeId_GNb:
		prC = C.GlobalE2node_ID_PR_gNB

		globalgNBIDC, err := newGlobalE2nodegNBID(choice.GNb)
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(globalgNBIDC))))
	case *e2apies.GlobalE2NodeId_ENb:
		prC = C.GlobalE2node_ID_PR_eNB

		globalEnbIDC, err := newGlobalE2nodeeNBID(choice.ENb)
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(globalEnbIDC))))
	default:
		return nil, fmt.Errorf("handling of %v not yet implemented", choice)
	}

	gnIDC := C.GlobalE2node_ID_t{
		present: prC,
		choice:  choiceC,
	}

	return &gnIDC, nil
}

func decodeGlobalE2NodeID(globalE2nodeIDchoice [48]byte) (*e2apies.GlobalE2NodeId, error) {

	present := C.long(binary.LittleEndian.Uint64(globalE2nodeIDchoice[0:8]))
	result := new(e2apies.GlobalE2NodeId)

	switch present {
	case C.GlobalE2node_ID_PR_gNB:
		bufC := globalE2nodeIDchoice[8:16]
		gNbC := *(**C.GlobalE2node_gNB_ID_t)(unsafe.Pointer(&bufC[0]))
		gNB, err := decodeGlobalE2nodegNBID(gNbC)
		if err != nil {
			return nil, fmt.Errorf("decodeGlobalE2NodeID() %v", err)
		}

		result.GlobalE2NodeId = &e2apies.GlobalE2NodeId_GNb{
			GNb: gNB,
		}
	case C.GlobalE2node_ID_PR_eNB:
		bufC := globalE2nodeIDchoice[8:16]
		eNbC := *(**C.GlobalE2node_eNB_ID_t)(unsafe.Pointer(&bufC[0]))
		eNB, err := decodeGlobalE2nodeeNBID(eNbC)
		if err != nil {
			return nil, fmt.Errorf("decodeGlobalE2nodeeNBID() %v", err)
		}

		result.GlobalE2NodeId = &e2apies.GlobalE2NodeId_ENb{
			ENb: eNB,
		}
	default:
		return nil, fmt.Errorf("decodeGlobalE2NodeID(). %v not yet implemneted", present)
	}

	return result, nil
}
