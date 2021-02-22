// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICsubsequentActionType.h"
import "C"
import (
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newRicSubsequentActionType(rsat e2apies.RicsubsequentActionType) (*C.RICsubsequentActionType_t, error) {
	var ret C.RICsubsequentActionType_t
	switch rsat {
	case e2apies.RicsubsequentActionType_RICSUBSEQUENT_ACTION_TYPE_CONTINUE:
		ret = C.RICsubsequentActionType_continue
	case e2apies.RicsubsequentActionType_RICSUBSEQUENT_ACTION_TYPE_WAIT:
		ret = C.RICsubsequentActionType_wait
	default:
		return nil, fmt.Errorf("unexpected RicsubsequentActionType %v", rsat)
	}
	return &ret, nil
}

func decodeRicSubsequentActionType(rsatC *C.RICsubsequentActionType_t) e2apies.RicsubsequentActionType {
	return e2apies.RicsubsequentActionType(*rsatC)
}
