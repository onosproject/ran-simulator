// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RANfunctionOID.h"
import "C"
import (
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
)

func newRanFunctionOID(oid *e2ap_commondatatypes.RanfunctionOid) *C.RANfunctionOID_t {

	return newPrintableString(string(oid.Value))
}

func decodeRanFunctionOID(ranFunctionOidC *C.RANfunctionOID_t) *e2ap_commondatatypes.RanfunctionOid {
	rfoPs := decodePrintableString(ranFunctionOidC)
	result := e2ap_commondatatypes.RanfunctionOid{
		Value: []byte(rfoPs),
	}

	return &result
}
