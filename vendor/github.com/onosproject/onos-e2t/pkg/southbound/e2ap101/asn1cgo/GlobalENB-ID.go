// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "GlobalENB-ID.h"
import "C"
import (
	"fmt"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	"unsafe"
)

func xerEncodeeNBID(gnbID *e2apies.GlobalEnbId) ([]byte, error) {
	gnbIDC, err := newGlobaleNBID(gnbID)
	if err != nil {
		return nil, err
	}

	bytes, err := encodeXer(&C.asn_DEF_GlobalENB_ID, unsafe.Pointer(gnbIDC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func perEncodeeNBID(gnbID *e2apies.GlobalEnbId) ([]byte, error) {
	gnbIDC, err := newGlobaleNBID(gnbID)
	if err != nil {
		return nil, err
	}

	bytes, err := encodePerBuffer(&C.asn_DEF_GlobalENB_ID, unsafe.Pointer(gnbIDC))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func newGlobaleNBID(id *e2apies.GlobalEnbId) (*C.GlobalENB_ID_t, error) {
	if len(id.PLmnIdentity.Value) > 3 {
		return nil, fmt.Errorf("plmnID max length is 3 - e2ap-v01.00.00.asn line 1105")
	}

	enbC, err := newEnbID(id.ENbId)
	if err != nil {
		return nil, err
	}

	idC := C.GlobalENB_ID_t{
		pLMN_Identity: *newOctetString(string(id.PLmnIdentity.Value)),
		eNB_ID:        *enbC,
	}

	return &idC, nil
}

func decodeGlobalEnbID(globalEnbID *C.GlobalENB_ID_t) (*e2apies.GlobalEnbId, error) {
	result := new(e2apies.GlobalEnbId)
	result.PLmnIdentity = new(e2ap_commondatatypes.PlmnIdentity)
	var err error
	result.PLmnIdentity.Value = []byte(decodeOctetString(&globalEnbID.pLMN_Identity))
	result.ENbId, err = decodeEnbID(&globalEnbID.eNB_ID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
