// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICactionType.h"
import "C"
import (
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newRicActionType(rat e2apies.RicactionType) (*C.RICactionType_t, error) {
	var ret C.RICactionType_t
	switch rat {
	case e2apies.RicactionType_RICACTION_TYPE_REPORT:
		ret = C.RICactionType_report
	case e2apies.RicactionType_RICACTION_TYPE_INSERT:
		ret = C.RICactionType_insert
	case e2apies.RicactionType_RICACTION_TYPE_POLICY:
		ret = C.RICactionType_policy
	default:
		return nil, fmt.Errorf("unexpected RicActionType %v", rat)
	}
	return &ret, nil
}

//func decodeRicActionTypeBytes(bytes []byte) e2apies.RicactionType {
//	raIDC := C.long(binary.LittleEndian.Uint64(bytes[:8]))
//	return decodeRicActionType(&raIDC)
//}

func decodeRicActionType(ratC *C.RICactionType_t) e2apies.RicactionType {
	return e2apies.RicactionType(*ratC)
}
