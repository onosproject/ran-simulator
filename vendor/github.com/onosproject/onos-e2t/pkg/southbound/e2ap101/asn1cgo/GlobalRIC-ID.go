// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "GlobalRIC-ID.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	"unsafe"
)

func newGlobalRicID(gr *e2apies.GlobalRicId) (*C.GlobalRIC_ID_t, error) {
	if len(gr.PLmnIdentity.Value) > 3 {
		return nil, fmt.Errorf("PLmnIdentity max length is 3 - e2ap-v01.00.00.asn line 1105")
	}
	if gr.RicId.Len != 20 {
		return nil, fmt.Errorf("ric-ID has to be 20 bits exactly - e2ap-v01.00.00.asn line 1076")
	}

	idC := C.GlobalRIC_ID_t{
		pLMN_Identity: *newOctetString(string(gr.PLmnIdentity.Value)),
		ric_ID:        *newBitString(gr.RicId),
	}

	return &idC, nil
}

func decodeGlobalRicIDBytes(bytes [112]byte) (*e2apies.GlobalRicId, error) {
	grIDC := C.GlobalRIC_ID_t{
		pLMN_Identity: C.OCTET_STRING_t{
			buf:  (*C.uchar)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[0:8])))),
			size: C.ulong(binary.LittleEndian.Uint64(bytes[8:16])),
		},
		ric_ID: C.BIT_STRING_t{
			buf:         (*C.uchar)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[40:48])))),
			size:        C.ulong(binary.LittleEndian.Uint64(bytes[48:56])),
			bits_unused: C.int(binary.LittleEndian.Uint32(bytes[56:60])),
		},
	}
	return decodeGlobalRicID(&grIDC)
}

func decodeGlobalRicID(grID *C.GlobalRIC_ID_t) (*e2apies.GlobalRicId, error) {
	result := new(e2apies.GlobalRicId)
	result.PLmnIdentity = new(e2ap_commondatatypes.PlmnIdentity)
	var err error
	result.PLmnIdentity.Value = []byte(decodeOctetString(&grID.pLMN_Identity))
	result.RicId, err = decodeBitString(&grID.ric_ID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
