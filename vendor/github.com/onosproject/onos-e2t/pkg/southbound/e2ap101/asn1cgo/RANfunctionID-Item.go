// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RANfunctionID-Item.h"
import "C"
import (
	"encoding/binary"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newRanFunctionIDItem(rfIDi *e2appducontents.RanfunctionIdItem) *C.RANfunctionID_Item_t {
	rfIDiC := C.RANfunctionID_Item_t{
		ranFunctionID:       newRanFunctionID(rfIDi.GetRanFunctionId()),
		ranFunctionRevision: *newRanFunctionRevision(rfIDi.GetRanFunctionRevision()),
	}

	return &rfIDiC
}

func decodeRanFunctionIDItemBytes(bytes [40]byte) (*e2appducontents.RanfunctionIdItem, error) {
	rfiC := C.RANfunctionID_Item_t{
		ranFunctionID:       C.long(binary.LittleEndian.Uint64(bytes[:8])),
		ranFunctionRevision: C.long(binary.LittleEndian.Uint64(bytes[8:16])),
	}

	return decodeRanFunctionIDItem(&rfiC)
}

func decodeRanFunctionIDItem(rfiC *C.RANfunctionID_Item_t) (*e2appducontents.RanfunctionIdItem, error) {
	rfi := e2appducontents.RanfunctionIdItem{
		RanFunctionId: &e2apies.RanfunctionId{
			Value: decodeRanFunctionID(&rfiC.ranFunctionID).Value,
		},
		RanFunctionRevision: &e2apies.RanfunctionRevision{
			Value: decodeRanFunctionRevision(&rfiC.ranFunctionRevision).Value,
		},
	}

	return &rfi, nil
}
