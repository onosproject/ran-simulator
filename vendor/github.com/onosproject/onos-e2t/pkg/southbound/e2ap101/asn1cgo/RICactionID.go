// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICactionID.h"
import "C"
import (
	"encoding/binary"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newRicActionID(raID *e2apies.RicactionId) *C.RICactionID_t {
	raIDC := C.RICactionID_t(raID.GetValue())
	return &raIDC
}

func decodeRicActionIDBytes(bytes []byte) *e2apies.RicactionId {
	raIDC := C.long(binary.LittleEndian.Uint64(bytes[:8]))
	return decodeRicActionID(&raIDC)
}

func decodeRicActionID(raIDC *C.RICactionID_t) *e2apies.RicactionId {
	return &e2apies.RicactionId{
		Value: int32(*raIDC),
	}
}
