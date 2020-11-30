// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "CriticalityDiagnostics-IE-List.h"
//#include "CriticalityDiagnostics-IE-Item.h"
import "C"

//func newCriticalityDiagnosticsIeList(cdIel *e2apies.CriticalityDiagnosticsIeList) (*C.CriticalityDiagnostics_IE_List_t, error) {
//	cdIelC := new(C.CriticalityDiagnostics_IE_List_t)
//	for _, cdIe := range cdIel.GetValue() {
//		cdIeC, err := newCriticalityDiagnosticsIEItem(cdIe)
//		if err != nil {
//			return nil, fmt.Errorf("error on newCriticalityDiagnosticsIEItem() %s", err.Error())
//		}
//		if _, err = C.asn_sequence_add(unsafe.Pointer(cdIelC), unsafe.Pointer(cdIeC)); err != nil {
//			return nil, err
//		}
//	}
//
//	return cdIelC, nil
//}
