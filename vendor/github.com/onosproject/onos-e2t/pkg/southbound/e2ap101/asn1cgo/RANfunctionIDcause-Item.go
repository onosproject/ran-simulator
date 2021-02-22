// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RANfunctionIDcause-Item.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newRanFunctionIDCauseItem(rfIDi *e2appducontents.RanfunctionIdcauseItem) (*C.RANfunctionIDcause_Item_t, error) {
	cause, err := newCause(rfIDi.GetCause())
	if err != nil {
		return nil, fmt.Errorf("newCause() error %s", err.Error())
	}

	rfIDiC := C.RANfunctionIDcause_Item_t{
		ranFunctionID: newRanFunctionID(rfIDi.GetRanFunctionId()),
		cause:         *cause,
	}

	return &rfIDiC, nil
}

func decodeRanFunctionIDcauseItemBytes(rfic [72]byte) (*e2appducontents.RanfunctionIdcauseItem, error) {
	rficC := C.RANfunctionIDcause_Item_t{
		ranFunctionID: C.long(binary.LittleEndian.Uint64(rfic[:8])),
		cause: C.Cause_t{
			present: C.Cause_PR(binary.LittleEndian.Uint64(rfic[8:16])),
		},
	}
	copy(rficC.cause.choice[:], rfic[16:24])

	return decodeRanFunctionIDCauseItem(&rficC)
}

func decodeRanFunctionIDCauseItem(rfiC *C.RANfunctionIDcause_Item_t) (*e2appducontents.RanfunctionIdcauseItem, error) {
	cause, err := decodeCause(&rfiC.cause)
	if err != nil {
		return nil, fmt.Errorf("decodeCause() error %s", err.Error())
	}
	rfi := e2appducontents.RanfunctionIdcauseItem{
		RanFunctionId: &e2apies.RanfunctionId{
			Value: decodeRanFunctionID(&rfiC.ranFunctionID).Value,
		},
		Cause: cause,
	}

	return &rfi, nil
}
