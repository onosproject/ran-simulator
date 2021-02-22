// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "ENB-ID.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	"unsafe"
)

func newEnbID(enbID *e2apies.EnbId) (*C.ENB_ID_t, error) {
	var pr C.ENB_ID_PR
	var bsC C.BIT_STRING_t

	switch enbt := enbID.EnbId.(type) {
	case *e2apies.EnbId_MacroENbId:
		if enbt.MacroENbId.Len != 20 {
			return nil, fmt.Errorf("MacroENbId must be exactly 20 bits")
		}
		pr = C.ENB_ID_PR_macro_eNB_ID
		bsC = *newBitString(enbt.MacroENbId)
	case *e2apies.EnbId_HomeENbId:
		if enbt.HomeENbId.Len != 28 {
			return nil, fmt.Errorf("MacroENbId must be exactly 20 bits")
		}
		pr = C.ENB_ID_PR_home_eNB_ID
		bsC = *newBitString(enbt.HomeENbId)
	case *e2apies.EnbId_ShortMacroENbId:
		if enbt.ShortMacroENbId.Len != 18 {
			return nil, fmt.Errorf("MacroENbId must be exactly 20 bits")
		}
		pr = C.ENB_ID_PR_short_Macro_eNB_ID
		bsC = *newBitString(enbt.ShortMacroENbId)
	case *e2apies.EnbId_LongMacroENbId:
		if enbt.LongMacroENbId.Len != 21 {
			return nil, fmt.Errorf("MacroENbId must be exactly 20 bits")
		}
		pr = C.ENB_ID_PR_long_Macro_eNB_ID
		bsC = *newBitString(enbt.LongMacroENbId)
	default:
		return nil, fmt.Errorf("unexpected type for eNB ID %v", enbt)
	}

	choiceC := [48]byte{}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(bsC.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(bsC.size))
	binary.LittleEndian.PutUint32(choiceC[16:], uint32(bsC.bits_unused))

	enbIDC := C.ENB_ID_t{
		present: pr,
		choice:  choiceC,
	}

	return &enbIDC, nil
}

func decodeEnbID(enbIDC *C.ENB_ID_t) (*e2apies.EnbId, error) {
	result := new(e2apies.EnbId)

	enbIDstructC := newBitStringFromArray(enbIDC.choice)
	bitString, err := decodeBitString(enbIDstructC)
	if err != nil {
		return nil, fmt.Errorf("decodeBitString() %s", err.Error())
	}

	switch enbIDC.present {
	case C.ENB_ID_PR_macro_eNB_ID:
		result.EnbId = &e2apies.EnbId_MacroENbId{
			MacroENbId: bitString,
		}
	case C.ENB_ID_PR_home_eNB_ID:
		result.EnbId = &e2apies.EnbId_HomeENbId{
			HomeENbId: bitString,
		}
	case C.ENB_ID_PR_short_Macro_eNB_ID:
		result.EnbId = &e2apies.EnbId_ShortMacroENbId{
			ShortMacroENbId: bitString,
		}
	case C.ENB_ID_PR_long_Macro_eNB_ID:
		result.EnbId = &e2apies.EnbId_LongMacroENbId{
			LongMacroENbId: bitString,
		}
	default:
		return nil, fmt.Errorf("decodeEnbID() unexpected %v", enbIDC.present)
	}

	return result, nil
}
