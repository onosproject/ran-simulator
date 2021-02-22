// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICaction-NotAdmitted-Item.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newRicActionNotAdmittedItem(raai *e2appducontents.RicactionNotAdmittedItem) (*C.RICaction_NotAdmitted_Item_t, error) {
	cause, err := newCause(raai.GetCause())
	if err != nil {
		return nil, fmt.Errorf("newCause() %s", err.Error())
	}
	raaiC := C.RICaction_NotAdmitted_Item_t{
		ricActionID: *newRicActionID(raai.RicActionId),
		cause:       *cause,
	}

	return &raaiC, nil
}

func decodeRicActionNotAdmittedItemBytes(ranaiBytes []byte) (*e2appducontents.RicactionNotAdmittedItem, error) {
	raaiC := C.RICaction_NotAdmitted_Item_t{
		ricActionID: C.long(binary.LittleEndian.Uint64(ranaiBytes[0:8])),
		cause: C.Cause_t{
			present: C.Cause_PR(binary.LittleEndian.Uint64(ranaiBytes[8:])),
		},
	}
	copy(raaiC.cause.choice[:8], ranaiBytes[16:24])

	return decodeRicActionNotAdmittedItem(raaiC)
}

func decodeRicActionNotAdmittedItem(ranaiC C.RICaction_NotAdmitted_Item_t) (*e2appducontents.RicactionNotAdmittedItem, error) {
	cause, err := decodeCause(&ranaiC.cause)
	if err != nil {
		return nil, fmt.Errorf("decodeCause() %s", err.Error())
	}

	return &e2appducontents.RicactionNotAdmittedItem{
		RicActionId: decodeRicActionID(&ranaiC.ricActionID),
		Cause:       cause,
	}, nil
}
