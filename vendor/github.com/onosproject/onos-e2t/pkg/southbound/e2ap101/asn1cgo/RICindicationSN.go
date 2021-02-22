// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICindicationSN.h"
import "C"
import (
	"encoding/binary"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newRicIndicationSn(rih *e2apies.RicindicationSn) *C.RICindicationSN_t {
	snC := C.long(rih.GetValue())
	return &snC
}

func decodeRicIndicationSnBytes(bytes []byte) *e2apies.RicindicationSn {
	raIDC := C.long(binary.LittleEndian.Uint64(bytes[:8]))
	return decodeRicIndicationSn(&raIDC)
}

func decodeRicIndicationSn(raIDC *C.RICindicationSN_t) *e2apies.RicindicationSn {
	return &e2apies.RicindicationSn{
		Value: int32(*raIDC),
	}
}
