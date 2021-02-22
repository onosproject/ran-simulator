// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "TypeOfError.h"
import "C"
import (
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newTypeOfError(toe e2apies.TypeOfError) (C.TypeOfError_t, error) {
	var ret C.TypeOfError_t
	switch toe {
	case e2apies.TypeOfError_TYPE_OF_ERROR_NOT_UNDERSTOOD:
		ret = C.TypeOfError_not_understood
	case e2apies.TypeOfError_TYPE_OF_ERROR_MISSING:
		ret = C.TypeOfError_missing
	default:
		return 0, fmt.Errorf("unexpected TypeOfError %v", toe)
	}
	return ret, nil
}

func decodeTypeOfError(toeC C.TypeOfError_t) e2apies.TypeOfError {
	return e2apies.TypeOfError(toeC)
}
