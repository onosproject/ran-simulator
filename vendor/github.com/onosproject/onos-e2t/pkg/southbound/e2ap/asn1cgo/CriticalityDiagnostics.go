// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "CriticalityDiagnostics.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"
	"unsafe"
)

func newCriticalityDiagnostics(cd *e2apies.CriticalityDiagnostics) (*C.CriticalityDiagnostics_t, error) {
	pc, err := newProcedureCode(cd.GetProcedureCode())
	if err != nil {
		return nil, fmt.Errorf("newProcedureCode() %s", err.Error())
	}
	tm, err := newTriggeringMessage(cd.GetTriggeringMessage())
	if err != nil {
		return nil, fmt.Errorf("newTriggeringMessage() %s", err.Error())
	}
	c, err := criticalityToC(cd.GetProcedureCriticality())
	if err != nil {
		return nil, fmt.Errorf("criticalityToC() %s", err.Error())
	}
	// TODO - add this back in
	//cdie, err := newCriticalityDiagnosticsIeList(cd.GetIEsCriticalityDiagnostics())
	//if err != nil {
	//	return nil, fmt.Errorf("newCriticalityDiagnosticsIeList() %s", err.Error())
	//}
	cdC := C.CriticalityDiagnostics_t{
		procedureCode:        &pc,
		triggeringMessage:    &tm,
		procedureCriticality: &c,
		ricRequestorID:       newRicRequestID(cd.GetRicRequestorId()),
		// TODO - add this back in
		//iEsCriticalityDiagnostics: cdie,
	}

	return &cdC, nil
}

func decodeCriticalityDiagnosticsBytes(bytes []byte) (*e2apies.CriticalityDiagnostics, error) {
	cdC := C.CriticalityDiagnostics_t{
		procedureCode:        (*C.long)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[:8])))),
		triggeringMessage:    (*C.long)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[8:])))),
		procedureCriticality: (*C.long)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[16:])))),
		//ricRequestorID: (*C.RICrequestID_t)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(bytes[24:])))),
	}
	return decodeCriticalityDiagnostics(&cdC)
}

func decodeCriticalityDiagnostics(cdC *C.CriticalityDiagnostics_t) (*e2apies.CriticalityDiagnostics, error) {

	ret := e2apies.CriticalityDiagnostics{
		ProcedureCode:             decodeProcedureCode(*cdC.procedureCode),
		TriggeringMessage:         decodeTriggeringMessage(*cdC.triggeringMessage),
		ProcedureCriticality:      decodeCriticality(*cdC.procedureCriticality),
		RicRequestorId:            nil,
		IEsCriticalityDiagnostics: nil,
	}

	return &ret, nil
}
