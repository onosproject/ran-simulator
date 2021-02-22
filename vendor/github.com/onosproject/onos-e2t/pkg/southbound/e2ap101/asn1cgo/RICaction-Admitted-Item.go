// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICaction-Admitted-Item.h"
import "C"
import (
	"encoding/binary"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newRicActionAdmittedItem(raai *e2appducontents.RicactionAdmittedItem) *C.RICaction_Admitted_Item_t {
	raaiC := C.RICaction_Admitted_Item_t{
		ricActionID: *newRicActionID(raai.RicActionId),
	}

	return &raaiC
}

func decodeRicActionAdmittedItemBytes(raaiBytes [32]byte) *e2appducontents.RicactionAdmittedItem {
	raaiC := C.RICaction_Admitted_Item_t{
		ricActionID: C.long(binary.LittleEndian.Uint64(raaiBytes[0:8])),
	}

	return decodeRicActionAdmittedItem(raaiC)
}

func decodeRicActionAdmittedItem(raaiC C.RICaction_Admitted_Item_t) *e2appducontents.RicactionAdmittedItem {
	return &e2appducontents.RicactionAdmittedItem{
		RicActionId: decodeRicActionID(&raaiC.ricActionID),
	}
}
