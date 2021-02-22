// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "GNB-ID-Choice.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	"unsafe"
)

func newGnbIDChoice(gnbIDCh *e2apies.GnbIdChoice) (*C.GNB_ID_Choice_t, error) {
	var pr C.GNB_ID_Choice_PR

	choiceC := [48]byte{}

	switch choice := gnbIDCh.GetGnbIdChoice().(type) {
	case *e2apies.GnbIdChoice_GnbId:
		pr = C.GNB_ID_Choice_PR_gnb_ID
		bsC := newBitString(choice.GnbId)
		//fmt.Printf("gNB ID %v %v %v %v\n", bsC, unsafe.Sizeof(bsC.size), unsafe.Sizeof(bsC.bits_unused), *bsC.buf)

		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(bsC.buf))))
		binary.LittleEndian.PutUint64(choiceC[8:], uint64(bsC.size))
		binary.LittleEndian.PutUint32(choiceC[16:], uint32(bsC.bits_unused))
	default:
		return nil, fmt.Errorf("newGnbIDChoice undhandled type %v", choice)
	}

	gnbChC := C.GNB_ID_Choice_t{
		present: pr,
		choice:  choiceC,
	}
	return &gnbChC, nil
}

func decodeGnbIDChoice(gnbIDC *C.GNB_ID_Choice_t) (*e2apies.GnbIdChoice, error) {
	result := new(e2apies.GnbIdChoice)

	switch gnbIDC.present {
	case C.GNB_ID_Choice_PR_gnb_ID:
		//fmt.Printf("GNB_ID_Choice_t %+v\n", gnbIDC.choice)
		gnbIDstructC := newBitStringFromArray(gnbIDC.choice)

		bitString, err := decodeBitString(gnbIDstructC)
		if err != nil {
			return nil, err
		}
		result.GnbIdChoice = &e2apies.GnbIdChoice_GnbId{
			GnbId: bitString,
		}
	default:
		return nil, fmt.Errorf("decodeGnbIDChoice() %v not yet implemented", gnbIDC.present)
	}

	return result, nil
}
