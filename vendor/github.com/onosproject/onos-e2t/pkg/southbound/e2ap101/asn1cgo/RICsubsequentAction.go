// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

// #cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
// #cgo LDFLAGS: -lm
// #include <stdio.h>
// #include <stdlib.h>
// #include <assert.h>
// #include "RICsubsequentAction.h"
import "C"
import (
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newRicSubsequentAction(rsa *e2apies.RicsubsequentAction) (*C.RICsubsequentAction_t, error) {
	rsatC, err := newRicSubsequentActionType(rsa.RicSubsequentActionType)
	if err != nil {
		return nil, fmt.Errorf("newRicSubsequentActionType() %s", err.Error())
	}
	rttwC, err := newRicTimeToWait(rsa.RicTimeToWait)
	if err != nil {
		return nil, fmt.Errorf("newRicTimeToWait() %s", err.Error())
	}
	rsaC := C.RICsubsequentAction_t{
		ricSubsequentActionType: *rsatC,
		ricTimeToWait:           *rttwC,
	}

	return &rsaC, nil
}

func decodeRicSubsequentAction(rsaC *C.RICsubsequentAction_t) (*e2apies.RicsubsequentAction, error) {
	rsa := e2apies.RicsubsequentAction{
		RicSubsequentActionType: decodeRicSubsequentActionType(&rsaC.ricSubsequentActionType),
		RicTimeToWait:           decodeRicTimeToWait(&rsaC.ricTimeToWait),
	}

	return &rsa, nil
}
