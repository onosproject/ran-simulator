// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "ProtocolIE-Container.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	"unsafe"

	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
)

func newE2SetupRequestIes(esv *e2appducontents.E2SetupRequestIes) (*C.ProtocolIE_Container_1710P11_t, error) {
	pIeC1710P11 := new(C.ProtocolIE_Container_1710P11_t)

	if esv.GetE2ApProtocolIes3() != nil {
		ie3C, err := newE2setupRequestIe3GlobalE2NodeID(esv.GetE2ApProtocolIes3())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P11), unsafe.Pointer(ie3C)); err != nil {
			return nil, err
		}
	}

	if esv.GetE2ApProtocolIes10() != nil {
		ie10C, err := newE2setupRequestIe10RanFunctionList(esv.GetE2ApProtocolIes10())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P11), unsafe.Pointer(ie10C)); err != nil {
			return nil, err
		}
	}

	if esv.GetE2ApProtocolIes33() != nil {
		ie33C, err := newE2setupRequestIe33E2nodeComponentConfigUpdateList(esv.GetE2ApProtocolIes33())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P11), unsafe.Pointer(ie33C)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P11, nil
}

func decodeE2SetupRequestIes(protocolIEsC *C.ProtocolIE_Container_1710P11_t) (*e2appducontents.E2SetupRequestIes, error) {
	pIEs := new(e2appducontents.E2SetupRequestIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P11 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i) // Forget the rest - this works - 7Nov20
		e2srIeC := *(**C.E2setupRequestIEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeE2setupRequestIE(e2srIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes3 != nil {
			pIEs.E2ApProtocolIes3 = ie.E2ApProtocolIes3
		}
		if ie.E2ApProtocolIes10 != nil {
			pIEs.E2ApProtocolIes10 = ie.E2ApProtocolIes10
		}
		if ie.E2ApProtocolIes33 != nil {
			pIEs.E2ApProtocolIes33 = ie.E2ApProtocolIes33
		}
	}

	return pIEs, nil
}

func newE2SetupResponseIes(e2srIEs *e2appducontents.E2SetupResponseIes) (*C.ProtocolIE_Container_1710P12_t, error) {
	pIeC1710P12 := new(C.ProtocolIE_Container_1710P12_t)

	if e2srIEs.GetE2ApProtocolIes4() != nil {
		ie4C, err := newE2setupResponseIe4GlobalRicID(e2srIEs.GetE2ApProtocolIes4())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P12), unsafe.Pointer(ie4C)); err != nil {
			return nil, err
		}
	}

	if e2srIEs.GetE2ApProtocolIes9() != nil {
		ie9C, err := newE2setupResponseIe9RanFunctionsAccepted(e2srIEs.GetE2ApProtocolIes9())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P12), unsafe.Pointer(ie9C)); err != nil {
			return nil, err
		}
	}

	if e2srIEs.GetE2ApProtocolIes13() != nil {
		ie13C, err := newE2setupResponseIe13RanFunctionsRejected(e2srIEs.GetE2ApProtocolIes13())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P12), unsafe.Pointer(ie13C)); err != nil {
			return nil, err
		}
	}

	if e2srIEs.GetE2ApProtocolIes35() != nil {
		ie35C, err := newE2setupResponseIe35E2nodeComponentConfigUpdateAckList(e2srIEs.GetE2ApProtocolIes35())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P12), unsafe.Pointer(ie35C)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P12, nil
}

func decodeE2SetupResponseIes(protocolIEsC *C.ProtocolIE_Container_1710P12_t) (*e2appducontents.E2SetupResponseIes, error) {
	pIEs := new(e2appducontents.E2SetupResponseIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P11 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		e2srIeC := *(**C.E2setupResponseIEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeE2setupResponseIE(e2srIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes4 != nil {
			pIEs.E2ApProtocolIes4 = ie.E2ApProtocolIes4
		}
		if ie.E2ApProtocolIes9 != nil {
			pIEs.E2ApProtocolIes9 = ie.E2ApProtocolIes9
		}
		if ie.E2ApProtocolIes13 != nil {
			pIEs.E2ApProtocolIes13 = ie.E2ApProtocolIes13
		}
		if ie.E2ApProtocolIes35 != nil {
			pIEs.E2ApProtocolIes35 = ie.E2ApProtocolIes35
		}
	}

	return pIEs, nil
}

func newRicSubscriptionResponseIe(rsrIEs *e2appducontents.RicsubscriptionResponseIes) (*C.ProtocolIE_Container_1710P1_t, error) {
	pIeC1710P1 := new(C.ProtocolIE_Container_1710P1_t)

	if rsrIEs.GetE2ApProtocolIes5() != nil {
		ie5C, err := newRicSubscriptionResponseIe5RanFunctionID(rsrIEs.GetE2ApProtocolIes5())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P1), unsafe.Pointer(ie5C)); err != nil {
			return nil, err
		}
	}
	if rsrIEs.GetE2ApProtocolIes17() != nil {
		ie17C, err := newRicSubscriptionResponseIe17RactionAdmittedList(rsrIEs.GetE2ApProtocolIes17())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P1), unsafe.Pointer(ie17C)); err != nil {
			return nil, err
		}
	}
	// TODO: Comment back in when RICactionRejected is handled
	//if rsrIEs.GetE2ApProtocolIes18() != nil {
	//	ie18C, err := newE2setupResponseIe4GlobalRicID(rsrIEs.GetE2ApProtocolIes18())
	//	if err != nil {
	//		return nil, err
	//	}
	//	if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P1), unsafe.Pointer(ie18C)); err != nil {
	//		return nil, err
	//	}
	//}
	if rsrIEs.GetE2ApProtocolIes29() != nil {
		ie29C, err := newRicSubscriptionResponseIe29RicRequestID(rsrIEs.GetE2ApProtocolIes29())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P1), unsafe.Pointer(ie29C)); err != nil {
			return nil, err
		}
	}
	return pIeC1710P1, nil
}

func decodeRicSubscriptionResponseIes(protocolIEsC *C.ProtocolIE_Container_1710P1_t) (*e2appducontents.RicsubscriptionResponseIes, error) {
	pIEs := new(e2appducontents.RicsubscriptionResponseIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P1 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		rsrIeC := *(**C.RICsubscriptionResponse_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicSubscriptionResponseIE(rsrIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes17 != nil {
			pIEs.E2ApProtocolIes17 = ie.E2ApProtocolIes17
		}
		if ie.E2ApProtocolIes18 != nil {
			pIEs.E2ApProtocolIes18 = ie.E2ApProtocolIes18
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
	}

	return pIEs, nil
}

func newRicSubscriptionRequestIes(rsrIEs *e2appducontents.RicsubscriptionRequestIes) (*C.ProtocolIE_Container_1710P0_t, error) {
	pIeC1710P0 := new(C.ProtocolIE_Container_1710P0_t)

	if rsrIEs.GetE2ApProtocolIes5() != nil {
		ie5C, err := newRicSubscriptionRequestIe5RanFunctionID(rsrIEs.E2ApProtocolIes5)
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P0), unsafe.Pointer(ie5C)); err != nil {
			return nil, err
		}
	}

	if rsrIEs.GetE2ApProtocolIes29() != nil {
		ie29C, err := newRicSubscriptionRequestIe29RicRequestID(rsrIEs.E2ApProtocolIes29)
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P0), unsafe.Pointer(ie29C)); err != nil {
			return nil, err
		}
	}

	if rsrIEs.GetE2ApProtocolIes30() != nil {
		ie30C, err := newRicSubscriptionRequestIe30RicSubscriptionDetails(rsrIEs.E2ApProtocolIes30)
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P0), unsafe.Pointer(ie30C)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P0, nil
}

func decodeRicSubscriptionRequestIes(protocolIEsC *C.ProtocolIE_Container_1710P0_t) (*e2appducontents.RicsubscriptionRequestIes, error) {
	pIEs := new(e2appducontents.RicsubscriptionRequestIes)

	ieCount := int(protocolIEsC.list.count)
	//	fmt.Printf("1544P0 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		rsrIeC := *(**C.RICsubscriptionRequest_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicSubscriptionRequestIE(rsrIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
		if ie.E2ApProtocolIes30 != nil {
			pIEs.E2ApProtocolIes30 = ie.E2ApProtocolIes30
		}
	}

	return pIEs, nil
}

func newRicIndicationIEs(riIes *e2appducontents.RicindicationIes) (*C.ProtocolIE_Container_1710P6_t, error) {
	pIeC1710P6 := new(C.ProtocolIE_Container_1710P6_t)

	if riIes.GetE2ApProtocolIes5() != nil {
		ie5c, err := newRicIndicationIe5RanFunctionID(riIes.GetE2ApProtocolIes5())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P6), unsafe.Pointer(ie5c)); err != nil {
			return nil, err
		}
	}

	if riIes.GetE2ApProtocolIes15() != nil {
		ie15c, err := newRicIndicationIe15RicActionID(riIes.GetE2ApProtocolIes15())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P6), unsafe.Pointer(ie15c)); err != nil {
			return nil, err
		}
	}

	if riIes.GetE2ApProtocolIes20() != nil {
		ie20c, err := newRicIndicationIe20RiccallProcessID(riIes.GetE2ApProtocolIes20())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P6), unsafe.Pointer(ie20c)); err != nil {
			return nil, err
		}
	}

	if riIes.GetE2ApProtocolIes25() != nil {
		ie25c, err := newRicIndicationIe25RicIndicationHeader(riIes.GetE2ApProtocolIes25())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P6), unsafe.Pointer(ie25c)); err != nil {
			return nil, err
		}
	}

	if riIes.GetE2ApProtocolIes26() != nil {
		ie26c, err := newRicIndicationIe26RicIndicationMessage(riIes.GetE2ApProtocolIes26())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P6), unsafe.Pointer(ie26c)); err != nil {
			return nil, err
		}
	}

	if riIes.GetE2ApProtocolIes27() != nil {
		ie27c, err := newRicIndicationIe27RicIndicationSn(riIes.GetE2ApProtocolIes27())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P6), unsafe.Pointer(ie27c)); err != nil {
			return nil, err
		}
	}

	if riIes.GetE2ApProtocolIes28() != nil {
		ie28c, err := newRicIndicationIe28RicIndicationType(riIes.GetE2ApProtocolIes28())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P6), unsafe.Pointer(ie28c)); err != nil {
			return nil, err
		}
	}

	if riIes.GetE2ApProtocolIes29() != nil {
		ie29c, err := newRicIndicationIe29RicRequestID(riIes.GetE2ApProtocolIes29())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P6), unsafe.Pointer(ie29c)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P6, nil
}

func decodeRicIndicationIes(protocolIEsC *C.ProtocolIE_Container_1710P6_t) (*e2appducontents.RicindicationIes, error) {
	pIEs := new(e2appducontents.RicindicationIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P6 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i) // Forget the rest - this works - 7Nov20
		riIeC := *(**C.RICindication_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicIndicationIE(riIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes15 != nil {
			pIEs.E2ApProtocolIes15 = ie.E2ApProtocolIes15
		}
		if ie.E2ApProtocolIes20 != nil {
			pIEs.E2ApProtocolIes20 = ie.E2ApProtocolIes20
		}
		if ie.E2ApProtocolIes25 != nil {
			pIEs.E2ApProtocolIes25 = ie.E2ApProtocolIes25
		}
		if ie.E2ApProtocolIes26 != nil {
			pIEs.E2ApProtocolIes26 = ie.E2ApProtocolIes26
		}
		if ie.E2ApProtocolIes27 != nil {
			pIEs.E2ApProtocolIes27 = ie.E2ApProtocolIes27
		}
		if ie.E2ApProtocolIes28 != nil {
			pIEs.E2ApProtocolIes28 = ie.E2ApProtocolIes28
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
	}

	return pIEs, nil
}

func newRicControlRequestIEs(rcRIes *e2appducontents.RiccontrolRequestIes) (*C.ProtocolIE_Container_1710P7_t, error) {
	pIeC1710P7 := new(C.ProtocolIE_Container_1710P7_t)

	if rcRIes.GetE2ApProtocolIes5() != nil {
		ie5c, err := newRicControlRequestIe5RanFunctionID(rcRIes.GetE2ApProtocolIes5())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P7), unsafe.Pointer(ie5c)); err != nil {
			return nil, err
		}
	}

	if rcRIes.GetE2ApProtocolIes20() != nil {
		ie20c, err := newRicControlRequestIe20RiccallProcessID(rcRIes.GetE2ApProtocolIes20())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P7), unsafe.Pointer(ie20c)); err != nil {
			return nil, err
		}
	}

	if rcRIes.GetE2ApProtocolIes22() != nil {
		ie22c, err := newRicControlRequestIe22RiccontrolHeader(rcRIes.GetE2ApProtocolIes22())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P7), unsafe.Pointer(ie22c)); err != nil {
			return nil, err
		}
	}

	if rcRIes.GetE2ApProtocolIes23() != nil {
		ie23c, err := newRicControlRequestIe23RiccontrolMessage(rcRIes.GetE2ApProtocolIes23())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P7), unsafe.Pointer(ie23c)); err != nil {
			return nil, err
		}
	}

	if rcRIes.GetE2ApProtocolIes21() != nil {
		ie21c, err := newRicControlRequestIe21RiccontrolAckRequest(rcRIes.GetE2ApProtocolIes21())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P7), unsafe.Pointer(ie21c)); err != nil {
			return nil, err
		}
	}

	if rcRIes.GetE2ApProtocolIes29() != nil {
		ie29c, err := newRicControlRequestIe29RicRequestID(rcRIes.GetE2ApProtocolIes29())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P7), unsafe.Pointer(ie29c)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P7, nil
}

func decodeRicControlRequestIes(protocolIEsC *C.ProtocolIE_Container_1710P7_t) (*e2appducontents.RiccontrolRequestIes, error) {
	pIEs := new(e2appducontents.RiccontrolRequestIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P6 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i) // Forget the rest - this works - 7Nov20
		riIeC := *(**C.RICcontrolRequest_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicControlRequestIE(riIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes20 != nil {
			pIEs.E2ApProtocolIes20 = ie.E2ApProtocolIes20
		}
		if ie.E2ApProtocolIes22 != nil {
			pIEs.E2ApProtocolIes22 = ie.E2ApProtocolIes22
		}
		if ie.E2ApProtocolIes23 != nil {
			pIEs.E2ApProtocolIes23 = ie.E2ApProtocolIes23
		}
		if ie.E2ApProtocolIes21 != nil {
			pIEs.E2ApProtocolIes21 = ie.E2ApProtocolIes21
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
	}

	return pIEs, nil
}

func newRicControlAcknowledgeIEs(rcaIes *e2appducontents.RiccontrolAcknowledgeIes) (*C.ProtocolIE_Container_1710P8_t, error) {
	pIeC1710P8 := new(C.ProtocolIE_Container_1710P8_t)

	if rcaIes.GetE2ApProtocolIes5() != nil {
		ie5c, err := newRicControlAcknowledgeIe5RanFunctionID(rcaIes.GetE2ApProtocolIes5())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P8), unsafe.Pointer(ie5c)); err != nil {
			return nil, err
		}
	}

	if rcaIes.GetE2ApProtocolIes20() != nil {
		ie20c, err := newRicControlAcknowledgeIe20RiccallProcessID(rcaIes.GetE2ApProtocolIes20())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P8), unsafe.Pointer(ie20c)); err != nil {
			return nil, err
		}
	}

	if rcaIes.GetE2ApProtocolIes24() != nil {
		ie22c, err := newRicControlAcknowledgeIe24RiccontrolStatus(rcaIes.GetE2ApProtocolIes24())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P8), unsafe.Pointer(ie22c)); err != nil {
			return nil, err
		}
	}

	if rcaIes.GetE2ApProtocolIes29() != nil {
		ie29c, err := newRicControlAcknowledgeIe29RicRequestID(rcaIes.GetE2ApProtocolIes29())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P8), unsafe.Pointer(ie29c)); err != nil {
			return nil, err
		}
	}

	if rcaIes.GetE2ApProtocolIes32() != nil {
		ie32c, err := newRicControlAcknowledgeIe32RiccontrolOutcome(rcaIes.GetE2ApProtocolIes32())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P8), unsafe.Pointer(ie32c)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P8, nil
}

func decodeRicControlAcknowledgeIes(protocolIEsC *C.ProtocolIE_Container_1710P8_t) (*e2appducontents.RiccontrolAcknowledgeIes, error) {
	pIEs := new(e2appducontents.RiccontrolAcknowledgeIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P6 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i) // Forget the rest - this works - 7Nov20
		riIeC := *(**C.RICcontrolAcknowledge_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicControlAcknowledgeIE(riIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes20 != nil {
			pIEs.E2ApProtocolIes20 = ie.E2ApProtocolIes20
		}
		if ie.E2ApProtocolIes24 != nil {
			pIEs.E2ApProtocolIes24 = ie.E2ApProtocolIes24
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
		if ie.E2ApProtocolIes32 != nil {
			pIEs.E2ApProtocolIes32 = ie.E2ApProtocolIes32
		}
	}

	return pIEs, nil
}

func newRicSubscriptionDeleteRequestIes(rsdrIEs *e2appducontents.RicsubscriptionDeleteRequestIes) (*C.ProtocolIE_Container_1710P3_t, error) {
	pIeC1710P3 := new(C.ProtocolIE_Container_1710P3_t)

	if rsdrIEs.GetE2ApProtocolIes5() != nil {
		ie5C, err := newRicSubscriptionDeleteRequestIe5RanFunctionID(rsdrIEs.E2ApProtocolIes5)
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P3), unsafe.Pointer(ie5C)); err != nil {
			return nil, err
		}
	}

	if rsdrIEs.GetE2ApProtocolIes29() != nil {
		ie29C, err := newRicSubscriptionDeleteRequestIe29RicRequestID(rsdrIEs.E2ApProtocolIes29)
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P3), unsafe.Pointer(ie29C)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P3, nil
}

func decodeRicSubscriptionDeleteRequestIes(protocolIEsC *C.ProtocolIE_Container_1710P3_t) (*e2appducontents.RicsubscriptionDeleteRequestIes, error) {
	pIEs := new(e2appducontents.RicsubscriptionDeleteRequestIes)

	ieCount := int(protocolIEsC.list.count)
	//	fmt.Printf("1544P0 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		rsrIeC := *(**C.RICsubscriptionDeleteRequest_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicSubscriptionDeleteRequestIE(rsrIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
	}

	return pIEs, nil
}

func newRicSubscriptionDeleteResponseIe(rsrIEs *e2appducontents.RicsubscriptionDeleteResponseIes) (*C.ProtocolIE_Container_1710P4_t, error) {
	pIeC1710P4 := new(C.ProtocolIE_Container_1710P4_t)

	if rsrIEs.GetE2ApProtocolIes5() != nil {
		ie5C, err := newRicSubscriptionDeleteResponseIe5RanFunctionID(rsrIEs.GetE2ApProtocolIes5())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P4), unsafe.Pointer(ie5C)); err != nil {
			return nil, err
		}
	}

	if rsrIEs.GetE2ApProtocolIes29() != nil {
		ie29C, err := newRicSubscriptionDeleteResponseIe29RicRequestID(rsrIEs.GetE2ApProtocolIes29())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P4), unsafe.Pointer(ie29C)); err != nil {
			return nil, err
		}
	}
	return pIeC1710P4, nil
}

func decodeRicSubscriptionDeleteResponseIes(protocolIEsC *C.ProtocolIE_Container_1710P4_t) (*e2appducontents.RicsubscriptionDeleteResponseIes, error) {
	pIEs := new(e2appducontents.RicsubscriptionDeleteResponseIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P1 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		rsrIeC := *(**C.RICsubscriptionDeleteResponse_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicSubscriptionDeleteResponseIE(rsrIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
	}

	return pIEs, nil
}

func newRicSubscriptionFailureIe(rsdIEs *e2appducontents.RicsubscriptionFailureIes) (*C.ProtocolIE_Container_1710P2_t, error) {
	pIeC1710P2 := new(C.ProtocolIE_Container_1710P2_t)

	if rsdIEs.GetE2ApProtocolIes2() != nil {
		ie2C, err := newRicSubscriptionFailureIe2CriticalityDiagnostics(rsdIEs.GetE2ApProtocolIes2())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P2), unsafe.Pointer(ie2C)); err != nil {
			return nil, err
		}
	}

	if rsdIEs.GetE2ApProtocolIes5() != nil {
		ie5C, err := newRicSubscriptionFailureIe5RanFunctionID(rsdIEs.GetE2ApProtocolIes5())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P2), unsafe.Pointer(ie5C)); err != nil {
			return nil, err
		}
	}

	if rsdIEs.GetE2ApProtocolIes18() != nil {
		ie2C, err := newRicSubscriptionFailureIe18RicActionNotAdmittedList(rsdIEs.GetE2ApProtocolIes18())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P2), unsafe.Pointer(ie2C)); err != nil {
			return nil, err
		}
	}

	if rsdIEs.GetE2ApProtocolIes29() != nil {
		ie29C, err := newRicSubscriptionFailureIe29RicRequestID(rsdIEs.GetE2ApProtocolIes29())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P2), unsafe.Pointer(ie29C)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P2, nil
}

func decodeRicSubscriptionFailureIes(protocolIEsC *C.ProtocolIE_Container_1710P2_t) (*e2appducontents.RicsubscriptionFailureIes, error) {
	pIEs := new(e2appducontents.RicsubscriptionFailureIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P1 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		rsfIeC := *(**C.RICsubscriptionFailure_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicSubscriptionFailureIE(rsfIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes2 != nil {
			pIEs.E2ApProtocolIes2 = ie.E2ApProtocolIes2
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes18 != nil {
			pIEs.E2ApProtocolIes18 = ie.E2ApProtocolIes18
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
	}

	return pIEs, nil
}

func newRicSubscriptionDeleteFailureIe(rsdfIEs *e2appducontents.RicsubscriptionDeleteFailureIes) (*C.ProtocolIE_Container_1710P5_t, error) {
	pIeC1710P5 := new(C.ProtocolIE_Container_1710P5_t)

	if rsdfIEs.GetE2ApProtocolIes5() != nil {
		ie5C, err := newRicSubscriptionDeleteFailureIe5RanFunctionID(rsdfIEs.GetE2ApProtocolIes5())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P5), unsafe.Pointer(ie5C)); err != nil {
			return nil, err
		}
	}

	if rsdfIEs.GetE2ApProtocolIes29() != nil {
		ie29C, err := newRicSubscriptionDeleteFailureIe29RicRequestID(rsdfIEs.GetE2ApProtocolIes29())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P5), unsafe.Pointer(ie29C)); err != nil {
			return nil, err
		}
	}

	if rsdfIEs.GetE2ApProtocolIes1() != nil {
		ie1C, err := newRicSubscriptionDeleteFailureIe1Cause(rsdfIEs.GetE2ApProtocolIes1())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P5), unsafe.Pointer(ie1C)); err != nil {
			return nil, err
		}
	}

	if rsdfIEs.GetE2ApProtocolIes2() != nil {
		ie2C, err := newRicSubscriptionDeleteFailureIe2CriticalityDiagnostics(rsdfIEs.GetE2ApProtocolIes2())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P5), unsafe.Pointer(ie2C)); err != nil {
			return nil, err
		}
	}
	return pIeC1710P5, nil
}

func decodeRicSubscriptionDeleteFailureIes(protocolIEsC *C.ProtocolIE_Container_1710P5_t) (*e2appducontents.RicsubscriptionDeleteFailureIes, error) {
	pIEs := new(e2appducontents.RicsubscriptionDeleteFailureIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P1 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		rsdfIeC := *(**C.RICsubscriptionDeleteFailure_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeRicSubscriptionDeleteFailureIE(rsdfIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
		if ie.E2ApProtocolIes1 != nil {
			pIEs.E2ApProtocolIes1 = ie.E2ApProtocolIes1
		}
		if ie.E2ApProtocolIes2 != nil {
			pIEs.E2ApProtocolIes2 = ie.E2ApProtocolIes2
		}
	}

	return pIEs, nil
}

func newErrorIndicationIe(eiIEs *e2appducontents.ErrorIndicationIes) (*C.ProtocolIE_Container_1710P10_t, error) {
	pIeC1710P10 := new(C.ProtocolIE_Container_1710P10_t)

	if eiIEs.GetE2ApProtocolIes2() != nil {
		ie2C, err := newErrorIndicationIe2CriticalityDiagnostics(eiIEs.GetE2ApProtocolIes2())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P10), unsafe.Pointer(ie2C)); err != nil {
			return nil, err
		}
	}

	if eiIEs.GetE2ApProtocolIes5() != nil {
		ie5C, err := newErrorIndicationIe5RanFunctionID(eiIEs.GetE2ApProtocolIes5())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P10), unsafe.Pointer(ie5C)); err != nil {
			return nil, err
		}
	}

	if eiIEs.GetE2ApProtocolIes1() != nil {
		ie1C, err := newErrorIndicationIe1Cause(eiIEs.GetE2ApProtocolIes1())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P10), unsafe.Pointer(ie1C)); err != nil {
			return nil, err
		}
	}

	if eiIEs.GetE2ApProtocolIes29() != nil {
		ie29C, err := newErrorIndicationIe29RicRequestID(eiIEs.GetE2ApProtocolIes29())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P10), unsafe.Pointer(ie29C)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P10, nil
}

func decodeErrorIndicationIes(protocolIEsC *C.ProtocolIE_Container_1710P10_t) (*e2appducontents.ErrorIndicationIes, error) {
	pIEs := new(e2appducontents.ErrorIndicationIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P1 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		eiIeC := *(**C.ErrorIndication_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeErrorIndicationIE(eiIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes2 != nil {
			pIEs.E2ApProtocolIes2 = ie.E2ApProtocolIes2
		}
		if ie.E2ApProtocolIes5 != nil {
			pIEs.E2ApProtocolIes5 = ie.E2ApProtocolIes5
		}
		if ie.E2ApProtocolIes1 != nil {
			pIEs.E2ApProtocolIes1 = ie.E2ApProtocolIes1
		}
		if ie.E2ApProtocolIes29 != nil {
			pIEs.E2ApProtocolIes29 = ie.E2ApProtocolIes29
		}
	}

	return pIEs, nil
}

func newE2setupFailureIe(e2sfIEs *e2appducontents.E2SetupFailureIes) (*C.ProtocolIE_Container_1710P13_t, error) {
	pIeC1710P13 := new(C.ProtocolIE_Container_1710P13_t)

	if e2sfIEs.GetE2ApProtocolIes1() != nil {
		ie1C, err := newE2setupFailureIe1Cause(e2sfIEs.GetE2ApProtocolIes1())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P13), unsafe.Pointer(ie1C)); err != nil {
			return nil, err
		}
	}

	if e2sfIEs.GetE2ApProtocolIes2() != nil {
		ie2C, err := newE2setupIe2CriticalityDiagnostics(e2sfIEs.GetE2ApProtocolIes2())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P13), unsafe.Pointer(ie2C)); err != nil {
			return nil, err
		}
	}

	if e2sfIEs.GetE2ApProtocolIes31() != nil {
		ie31C, err := newE2setupFailureIe31TimeToWait(e2sfIEs.GetE2ApProtocolIes31())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P13), unsafe.Pointer(ie31C)); err != nil {
			return nil, err
		}
	}

	if e2sfIEs.GetE2ApProtocolIes48() != nil {
		ie48C, err := newE2setupFailureIe48Tnlinformation(e2sfIEs.GetE2ApProtocolIes48())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P13), unsafe.Pointer(ie48C)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P13, nil
}

func decodeE2setupFailureIes(protocolIEsC *C.ProtocolIE_Container_1710P13_t) (*e2appducontents.E2SetupFailureIes, error) {
	pIEs := new(e2appducontents.E2SetupFailureIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P1 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		eiIeC := *(**C.E2setupFailureIEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		ie, err := decodeE2setupFailureIE(eiIeC)
		if err != nil {
			return nil, err
		}
		if ie.E2ApProtocolIes1 != nil {
			pIEs.E2ApProtocolIes1 = ie.E2ApProtocolIes1
		}
		if ie.E2ApProtocolIes2 != nil {
			pIEs.E2ApProtocolIes2 = ie.E2ApProtocolIes2
		}
		if ie.E2ApProtocolIes31 != nil {
			pIEs.E2ApProtocolIes31 = ie.E2ApProtocolIes31
		}
		if ie.E2ApProtocolIes48 != nil {
			pIEs.E2ApProtocolIes48 = ie.E2ApProtocolIes48
		}
	}

	return pIEs, nil
}

func newE2ConnectionUpdateIe(e2cuIEs *e2appducontents.E2ConnectionUpdateIes) (*C.ProtocolIE_Container_1710P14_t, error) {
	pIeC1710P14 := new(C.ProtocolIE_Container_1710P14_t)

	if e2cuIEs.GetE2ApProtocolIes44() != nil {
		ie44C, err := newE2connectionUpdateIe44E2connectionUpdateList(e2cuIEs.GetE2ApProtocolIes44())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P14), unsafe.Pointer(ie44C)); err != nil {
			return nil, err
		}
	}

	if e2cuIEs.GetE2ApProtocolIes45() != nil {
		ie45C, err := newE2connectionUpdateIe45E2connectionUpdateList(e2cuIEs.GetE2ApProtocolIes45())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P14), unsafe.Pointer(ie45C)); err != nil {
			return nil, err
		}
	}

	if e2cuIEs.GetE2ApProtocolIes46() != nil {
		ie46C, err := newE2connectionUpdateIe46E2connectionUpdateRemoveList(e2cuIEs.GetE2ApProtocolIes46())
		if err != nil {
			return nil, err
		}
		if _, err = C.asn_sequence_add(unsafe.Pointer(pIeC1710P14), unsafe.Pointer(ie46C)); err != nil {
			return nil, err
		}
	}

	return pIeC1710P14, nil
}

func decodeE2connectionUpdateIes(protocolIEsC *C.ProtocolIE_Container_1710P14_t) (*e2appducontents.E2ConnectionUpdateIes, error) {
	pIEs := new(e2appducontents.E2ConnectionUpdateIes)

	ieCount := int(protocolIEsC.list.count)
	//fmt.Printf("1544P1 Type %T Count %v Size %v\n", *protocolIEsC.list.array, protocolIEsC.list.count, protocolIEsC.list.size)
	for i := 0; i < ieCount; i++ {
		//ToDo - uncomment once decodeE2connectionUpdateIe is implemented
		//offset := unsafe.Sizeof(unsafe.Pointer(*protocolIEsC.list.array)) * uintptr(i)
		//eiIeC := *(**C.E2connectionUpdate_IEs_t)(unsafe.Pointer(uintptr(unsafe.Pointer(protocolIEsC.list.array)) + offset))

		//ToDo - Implement decodeE2connectionUpdateIe function -- see analogical function in RICcontrolRequest message chain
		//ie, err := decodeE2connectionUpdateIe(eiIeC)
		//if err != nil {
		//	return nil, err
		//}
		//if ie.E2ApProtocolIes44 != nil {
		//	pIEs.E2ApProtocolIes44 = ie.E2ApProtocolIes44
		//}
		//if ie.E2ApProtocolIes45 != nil {
		//	pIEs.E2ApProtocolIes45 = ie.E2ApProtocolIes45
		//}
		//if ie.E2ApProtocolIes46 != nil {
		//	pIEs.E2ApProtocolIes46 = ie.E2ApProtocolIes46
		//}
	}

	return pIEs, nil
}
