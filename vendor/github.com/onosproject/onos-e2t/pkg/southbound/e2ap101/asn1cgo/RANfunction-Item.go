// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RANfunction-Item.h"
import "C"
import (
	"encoding/binary"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func newRanFunctionItem(rfItem *e2appducontents.RanfunctionItem) *C.RANfunction_Item_t {
	rfItemC := C.RANfunction_Item_t{
		ranFunctionID:         newRanFunctionID(rfItem.GetRanFunctionId()),
		ranFunctionRevision:   newRanFunctionRevision(rfItem.GetRanFunctionRevision()),
		ranFunctionDefinition: *newOctetString(string(rfItem.GetRanFunctionDefinition().GetValue())),
	}
	return &rfItemC
}

func decodeRanFunctionItemBytes(bytes [88]byte) (*e2appducontents.RanfunctionItem, error) {
	size := binary.LittleEndian.Uint64(bytes[16:24])
	gobytes := C.GoBytes(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[8:16]))), C.int(size))

	rfiC := C.RANfunction_Item_t{
		ranFunctionID: C.long(binary.LittleEndian.Uint64(bytes[:8])),
		ranFunctionDefinition: C.OCTET_STRING_t{
			buf:  (*C.uchar)(C.CBytes(gobytes)),
			size: C.ulong(size),
		},
		ranFunctionRevision: C.long(binary.LittleEndian.Uint64(bytes[24:32])),
	}

	return decodeRanFunctionItem(&rfiC)
}

func decodeRanFunctionItem(rfiC *C.RANfunction_Item_t) (*e2appducontents.RanfunctionItem, error) {
	rfi := e2appducontents.RanfunctionItem{
		RanFunctionId: &e2apies.RanfunctionId{
			Value: decodeRanFunctionID(&rfiC.ranFunctionID).Value,
		},
		RanFunctionRevision: &e2apies.RanfunctionRevision{
			Value: decodeRanFunctionRevision(&rfiC.ranFunctionRevision).Value,
		},
		RanFunctionDefinition: &e2ap_commondatatypes.RanfunctionDefinition{
			Value: []byte(decodeOctetString(&rfiC.ranFunctionDefinition)),
		},
	}

	return &rfi, nil
}
