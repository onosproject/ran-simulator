// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RANfunctionID.h"
import "C"
import (
	"encoding/binary"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newRanFunctionID(rfID *e2apies.RanfunctionId) C.long {
	return C.long(rfID.Value)
}

func decodeRanFunctionIDBytes(ranFunctionIDCbytes []byte) *e2apies.RanfunctionId {
	rfC := (C.RANfunctionID_t)(binary.LittleEndian.Uint64(ranFunctionIDCbytes[0:8]))

	return decodeRanFunctionID(&rfC)
}

func decodeRanFunctionID(ranFunctionIDC *C.RANfunctionID_t) *e2apies.RanfunctionId {
	result := e2apies.RanfunctionId{
		Value: int32(*ranFunctionIDC),
	}

	return &result
}
