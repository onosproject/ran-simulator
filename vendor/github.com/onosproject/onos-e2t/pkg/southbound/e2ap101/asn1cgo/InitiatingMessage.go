// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "InitiatingMessage.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2ap_constants "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-constants"
	e2appdudescriptions "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-descriptions"
	"unsafe"
)

func newInitiatingMessage(im *e2appdudescriptions.InitiatingMessage) (*C.struct_InitiatingMessage, error) {

	var presentC C.InitiatingMessage__value_PR
	var pcC C.ProcedureCode_t
	var critC C.Criticality_t
	choiceC := [72]byte{} // The size of the InitiatingMessage__value_u union

	if pc := im.GetProcedureCode().GetE2Setup(); pc != nil &&
		pc.GetInitiatingMessage() != nil {

		presentC = C.InitiatingMessage__value_PR_E2setupRequest
		pcC = C.ProcedureCode_id_E2setup
		critC = C.long(C.Criticality_reject)
		e2sC, err := newE2SetupRequest(pc.GetInitiatingMessage())
		if err != nil {
			return nil, err
		}
		//	//fmt.Printf("Protocol IEs %v %v %v\n", rsrC.protocolIEs.list.array, rsrC.protocolIEs.list.count, rsrC.protocolIEs.list.size)
		//	// Now copy the rsrC over in to the choice byte by byte - the union is [72]byte
		//	// It's A_SET_OF, so has <address(8), count(4), size(4)>
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(e2sC.protocolIEs.list.array))))
		binary.LittleEndian.PutUint32(choiceC[8:], uint32(e2sC.protocolIEs.list.count))
		binary.LittleEndian.PutUint32(choiceC[12:], uint32(e2sC.protocolIEs.list.size))
	} else if pc := im.GetProcedureCode().GetRicSubscription(); pc != nil &&
		pc.GetInitiatingMessage() != nil {

		presentC = C.InitiatingMessage__value_PR_RICsubscriptionRequest
		pcC = C.ProcedureCode_id_RICsubscription
		critC = C.long(C.Criticality_reject)
		rsC, err := newRICsubscriptionRequest(pc.GetInitiatingMessage())
		if err != nil {
			return nil, err
		}
		//	//fmt.Printf("Protocol IEs %v %v %v\n", rsrC.protocolIEs.list.array, rsrC.protocolIEs.list.count, rsrC.protocolIEs.list.size)
		//	// Now copy the rsrC over in to the choice byte by byte - the union is [72]byte
		//	// It's A_SET_OF, so has <address(8), count(4), size(4)>
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(rsC.protocolIEs.list.array))))
		binary.LittleEndian.PutUint32(choiceC[8:], uint32(rsC.protocolIEs.list.count))
		binary.LittleEndian.PutUint32(choiceC[12:], uint32(rsC.protocolIEs.list.size))

	} else if pc := im.GetProcedureCode().GetRicSubscriptionDelete(); pc != nil &&
		pc.GetInitiatingMessage() != nil {

		presentC = C.InitiatingMessage__value_PR_RICsubscriptionDeleteRequest
		pcC = C.ProcedureCode_id_RICsubscriptionDelete
		critC = C.long(C.Criticality_reject)
		rsC, err := newRICsubscriptionDeleteRequest(pc.GetInitiatingMessage())
		if err != nil {
			return nil, err
		}
		//	//fmt.Printf("Protocol IEs %v %v %v\n", rsrC.protocolIEs.list.array, rsrC.protocolIEs.list.count, rsrC.protocolIEs.list.size)
		//	// Now copy the rsrC over in to the choice byte by byte - the union is [72]byte
		//	// It's A_SET_OF, so has <address(8), count(4), size(4)>
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(rsC.protocolIEs.list.array))))
		binary.LittleEndian.PutUint32(choiceC[8:], uint32(rsC.protocolIEs.list.count))
		binary.LittleEndian.PutUint32(choiceC[12:], uint32(rsC.protocolIEs.list.size))
	} else if pc := im.GetProcedureCode().GetRicIndication(); pc != nil &&
		pc.GetInitiatingMessage() != nil {

		presentC = C.InitiatingMessage__value_PR_RICindication
		pcC = C.ProcedureCode_id_RICindication
		critC = C.long(C.Criticality_reject)
		e2sC, err := newRicIndication(pc.GetInitiatingMessage())
		if err != nil {
			return nil, err
		}
		//	//fmt.Printf("Protocol IEs %v %v %v\n", rsrC.protocolIEs.list.array, rsrC.protocolIEs.list.count, rsrC.protocolIEs.list.size)
		//	// Now copy the rsrC over in to the choice byte by byte - the union is [72]byte
		//	// It's A_SET_OF, so has <address(8), count(4), size(4)>
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(e2sC.protocolIEs.list.array))))
		binary.LittleEndian.PutUint32(choiceC[8:], uint32(e2sC.protocolIEs.list.count))
		binary.LittleEndian.PutUint32(choiceC[12:], uint32(e2sC.protocolIEs.list.size))

	} else if pc := im.GetProcedureCode().GetRicControl(); pc != nil &&
		pc.GetInitiatingMessage() != nil {

		presentC = C.InitiatingMessage__value_PR_RICcontrolRequest
		pcC = C.ProcedureCode_id_RICcontrol
		critC = C.long(C.Criticality_reject)
		e2sC, err := newRicControlRequest(pc.GetInitiatingMessage())
		if err != nil {
			return nil, err
		}
		//	//fmt.Printf("Protocol IEs %v %v %v\n", rsrC.protocolIEs.list.array, rsrC.protocolIEs.list.count, rsrC.protocolIEs.list.size)
		//	// Now copy the rsrC over in to the choice byte by byte - the union is [72]byte
		//	// It's A_SET_OF, so has <address(8), count(4), size(4)>
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(e2sC.protocolIEs.list.array))))
		binary.LittleEndian.PutUint32(choiceC[8:], uint32(e2sC.protocolIEs.list.count))
		binary.LittleEndian.PutUint32(choiceC[12:], uint32(e2sC.protocolIEs.list.size))

	} else if pc := im.GetProcedureCode().GetErrorIndication(); pc != nil &&
		pc.GetInitiatingMessage() != nil {

		presentC = C.InitiatingMessage__value_PR_ErrorIndication
		pcC = C.ProcedureCode_id_ErrorIndication
		critC = C.long(C.Criticality_ignore)
		e2sC, err := newErrorIndication(pc.GetInitiatingMessage())
		if err != nil {
			return nil, err
		}
		//	//fmt.Printf("Protocol IEs %v %v %v\n", rsrC.protocolIEs.list.array, rsrC.protocolIEs.list.count, rsrC.protocolIEs.list.size)
		//	// Now copy the rsrC over in to the choice byte by byte - the union is [72]byte
		//	// It's A_SET_OF, so has <address(8), count(4), size(4)>
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(e2sC.protocolIEs.list.array))))
		binary.LittleEndian.PutUint32(choiceC[8:], uint32(e2sC.protocolIEs.list.count))
		binary.LittleEndian.PutUint32(choiceC[12:], uint32(e2sC.protocolIEs.list.size))

	} else {
		return nil, fmt.Errorf("newInitiatingMessageValue type not yet implemented")
	}

	imC := C.InitiatingMessage_t{
		procedureCode: pcC,
		criticality:   critC,
		value: C.struct_InitiatingMessage__value{
			present: presentC,
			choice:  choiceC,
		},
	}

	return &imC, nil
}

func decodeInitiatingMessage(initMsgC *C.InitiatingMessage_t) (*e2appdudescriptions.InitiatingMessage, error) {

	initiatingMessage := new(e2appdudescriptions.InitiatingMessage)

	listArrayAddr := initMsgC.value.choice[0:8]

	switch initMsgC.value.present {
	case C.InitiatingMessage__value_PR_E2setupRequest:
		e2srC := *(**C.E2setupRequestIEs_t)(unsafe.Pointer(&listArrayAddr[0]))
		esC := C.E2setupRequest_t{
			protocolIEs: C.ProtocolIE_Container_1710P11_t{
				list: C.struct___53{ // TODO: tie this down with a predictable name
					array: (**C.E2setupRequestIEs_t)(unsafe.Pointer(e2srC)),
					count: C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[8:12])),
					size:  C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[12:16])),
				},
			},
		}
		//fmt.Printf("E2SetupRequestC %+v\n %+v\n", initMsgC, riC)
		e2sr, err := decodeE2setupRequest(&esC)
		if err != nil {
			return nil, err
		}
		initiatingMessage.ProcedureCode = &e2appdudescriptions.E2ApElementaryProcedures{
			E2Setup: &e2appdudescriptions.E2Setup{
				InitiatingMessage: e2sr,
				ProcedureCode: &e2ap_constants.IdE2Setup{
					Value: int32(v1beta2.ProcedureCodeIDE2setup),
				},
				Criticality: &e2ap_commondatatypes.CriticalityReject{},
			},
		}
	case C.InitiatingMessage__value_PR_RICsubscriptionRequest:
		ricsrC := *(**C.RICsubscriptionRequest_IEs_t)(unsafe.Pointer(&listArrayAddr[0]))
		srC := C.RICsubscriptionRequest_t{
			protocolIEs: C.ProtocolIE_Container_1710P0_t{
				list: C.struct___86{ // TODO: tie this down with a predictable name
					array: (**C.RICsubscriptionRequest_IEs_t)(unsafe.Pointer(ricsrC)),
					count: C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[8:12])),
					size:  C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[12:16])),
				},
			},
		}
		//fmt.Printf("RICsubscriptionRequest_t %+v\n %+v\n", initMsgC, sdrC)

		sr, err := decodeRicSubscriptionRequest(&srC)
		if err != nil {
			return nil, err
		}

		// TODO: Get the value
		initiatingMessage.ProcedureCode = &e2appdudescriptions.E2ApElementaryProcedures{
			RicSubscription: &e2appdudescriptions.RicSubscription{
				InitiatingMessage: sr,
				ProcedureCode: &e2ap_constants.IdRicsubscription{
					Value: int32(v1beta2.ProcedureCodeIDRICsubscription),
				},
				Criticality: &e2ap_commondatatypes.CriticalityReject{},
			},
		}

	case C.InitiatingMessage__value_PR_RICsubscriptionDeleteRequest:
		ricsdrC := *(**C.RICsubscriptionDeleteRequest_IEs_t)(unsafe.Pointer(&listArrayAddr[0]))
		sdrC := C.RICsubscriptionDeleteRequest_t{
			protocolIEs: C.ProtocolIE_Container_1710P3_t{
				list: C.struct___80{ // TODO: tie this down with a predictable name
					array: (**C.RICsubscriptionDeleteRequest_IEs_t)(unsafe.Pointer(ricsdrC)),
					count: C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[8:12])),
					size:  C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[12:16])),
				},
			},
		}
		//fmt.Printf("RICsubscriptionRequest_t %+v\n %+v\n", initMsgC, sdrC)

		sdr, err := decodeRicSubscriptionDeleteRequest(&sdrC)
		if err != nil {
			return nil, err
		}

		// TODO: Get the value
		initiatingMessage.ProcedureCode = &e2appdudescriptions.E2ApElementaryProcedures{
			RicSubscriptionDelete: &e2appdudescriptions.RicSubscriptionDelete{
				InitiatingMessage: sdr,
				ProcedureCode: &e2ap_constants.IdRicsubscriptionDelete{
					Value: int32(v1beta2.ProcedureCodeIDRICsubscriptionDelete),
				},
				Criticality: &e2ap_commondatatypes.CriticalityReject{},
			},
		}

	case C.InitiatingMessage__value_PR_RICindication:
		riIesC := *(**C.RICindication_IEs_t)(unsafe.Pointer(&listArrayAddr[0]))
		riC := C.RICindication_t{
			protocolIEs: C.ProtocolIE_Container_1710P6_t{
				list: C.struct___77{ // TODO: tie this down with a predictable name
					array: (**C.RICindication_IEs_t)(unsafe.Pointer(riIesC)),
					count: C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[8:12])),
					size:  C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[12:16])),
				},
			},
		}
		fmt.Printf("RICindication_t %+v\n %+v\n", initMsgC, riC)

		ri, err := decodeRicIndication(&riC)
		if err != nil {
			return nil, fmt.Errorf("decodeRicIndication() %s", err.Error())
		}
		initiatingMessage.ProcedureCode = &e2appdudescriptions.E2ApElementaryProcedures{
			RicIndication: &e2appdudescriptions.RicIndication{
				InitiatingMessage: ri,
				ProcedureCode: &e2ap_constants.IdRicindication{
					Value: int32(v1beta2.ProcedureCodeIDRICindication),
				},
				Criticality: &e2ap_commondatatypes.CriticalityIgnore{
					Criticality: e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE,
				},
			},
		}

	case C.InitiatingMessage__value_PR_RICcontrolRequest:
		rcrIesC := *(**C.RICcontrolRequest_IEs_t)(unsafe.Pointer(&listArrayAddr[0]))
		rcrC := C.RICcontrolRequest_t{
			protocolIEs: C.ProtocolIE_Container_1710P7_t{
				list: C.struct___76{ // TODO: tie this down with a predictable name
					array: (**C.RICcontrolRequest_IEs_t)(unsafe.Pointer(rcrIesC)),
					count: C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[8:12])),
					size:  C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[12:16])),
				},
			},
		}
		fmt.Printf("RICcontrolRequest_t %+v\n %+v\n", initMsgC, rcrC)

		rcr, err := decodeRicControlRequest(&rcrC)
		if err != nil {
			return nil, fmt.Errorf("decodeRicControlRequest() %s", err.Error())
		}
		initiatingMessage.ProcedureCode = &e2appdudescriptions.E2ApElementaryProcedures{
			RicControl: &e2appdudescriptions.RicControl{
				InitiatingMessage: rcr,
				ProcedureCode: &e2ap_constants.IdRiccontrol{
					Value: int32(v1beta2.ProcedureCodeIDRICcontrol),
				},
				Criticality: &e2ap_commondatatypes.CriticalityReject{
					Criticality: e2ap_commondatatypes.Criticality_CRITICALITY_REJECT,
				},
			},
		}

	case C.InitiatingMessage__value_PR_ErrorIndication:
		riIesC := *(**C.ErrorIndication_IEs_t)(unsafe.Pointer(&listArrayAddr[0]))
		riC := C.ErrorIndication_t{
			protocolIEs: C.ProtocolIE_Container_1710P10_t{
				list: C.struct___52{ // TODO: tie this down with a predictable name
					array: (**C.ErrorIndication_IEs_t)(unsafe.Pointer(riIesC)),
					count: C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[8:12])),
					size:  C.int(binary.LittleEndian.Uint32(initMsgC.value.choice[12:16])),
				},
			},
		}
		fmt.Printf("ErrorIndication_t %+v\n %+v\n", initMsgC, riC)

		ri, err := decodeErrorIndication(&riC)
		if err != nil {
			return nil, fmt.Errorf("decodeErrorIndication() %s", err.Error())
		}
		initiatingMessage.ProcedureCode = &e2appdudescriptions.E2ApElementaryProcedures{
			ErrorIndication: &e2appdudescriptions.ErrorIndicationEp{
				InitiatingMessage: ri,
				ProcedureCode: &e2ap_constants.IdErrorIndication{
					Value: int32(v1beta2.ProcedureCodeIDErrorIndication),
				},
				Criticality: &e2ap_commondatatypes.CriticalityIgnore{
					Criticality: e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE,
				},
			},
		}

	default:
		return nil, fmt.Errorf("decodeInitiatingMessage() %v not yet implemented", initMsgC.value.present)
	}

	return initiatingMessage, nil
}
