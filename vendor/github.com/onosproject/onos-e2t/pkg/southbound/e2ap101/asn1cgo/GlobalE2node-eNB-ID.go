// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "GlobalE2node-eNB-ID.h"
import "C"
import (
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newGlobalE2nodeeNBID(enbID *e2apies.GlobalE2NodeEnbId) (*C.GlobalE2node_eNB_ID_t, error) {

	globaleNBID, err := newGlobaleNBID(enbID.GlobalENbId)
	if err != nil {
		return nil, err
	}

	globaleNBIDC := C.GlobalE2node_eNB_ID_t{
		global_eNB_ID: *globaleNBID,
	}

	return &globaleNBIDC, nil
}

func decodeGlobalE2nodeeNBID(eNBC *C.GlobalE2node_eNB_ID_t) (*e2apies.GlobalE2NodeEnbId, error) {
	result := new(e2apies.GlobalE2NodeEnbId)
	var err error
	result.GlobalENbId, err = decodeGlobalEnbID(&eNBC.global_eNB_ID)
	if err != nil {
		return nil, fmt.Errorf("error decodeGlobalE2nodeeNBID() %v", err)
	}

	return result, nil
}
