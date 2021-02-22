// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "GlobalgNB-ID.h"
import "C"
import (
	"fmt"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	"unsafe"
)

func xerEncodegNBID(gnbID *e2apies.GlobalgNbId) ([]byte, error) {
	gnbIDC, err := newGlobalgNBID(gnbID)
	if err != nil {
		return nil, err
	}

	bytes, err := encodeXer(&C.asn_DEF_GlobalgNB_ID, unsafe.Pointer(gnbIDC))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func perEncodegNBID(gnbID *e2apies.GlobalgNbId) ([]byte, error) {
	gnbIDC, err := newGlobalgNBID(gnbID)
	if err != nil {
		return nil, err
	}

	bytes, err := encodePerBuffer(&C.asn_DEF_GlobalgNB_ID, unsafe.Pointer(gnbIDC))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func newGlobalgNBID(id *e2apies.GlobalgNbId) (*C.GlobalgNB_ID_t, error) {
	if len(id.PlmnId.Value) > 3 {
		return nil, fmt.Errorf("plmnID max length is 3 - e2ap-v01.00.00.asn line 1105")
	}

	gnbChoiceC, err := newGnbIDChoice(id.GnbId)
	if err != nil {
		return nil, err
	}

	idC := C.GlobalgNB_ID_t{
		plmn_id: *newOctetString(string(id.PlmnId.Value)),
		gnb_id:  *gnbChoiceC,
	}

	return &idC, nil
}

func decodeGlobalGnbID(globalGnbID *C.GlobalgNB_ID_t) (*e2apies.GlobalgNbId, error) {
	result := new(e2apies.GlobalgNbId)
	result.PlmnId = new(e2ap_commondatatypes.PlmnIdentity)
	var err error
	result.PlmnId.Value = []byte(decodeOctetString(&globalGnbID.plmn_id))
	result.GnbId, err = decodeGnbIDChoice(&globalGnbID.gnb_id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
