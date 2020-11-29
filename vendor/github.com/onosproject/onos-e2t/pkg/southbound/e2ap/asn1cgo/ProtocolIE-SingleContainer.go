// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "ProtocolIE-SingleContainer.h"
//#include "RICaction-ToBeSetup-Item.h"
//#include "ProtocolIE-Field.h"
import "C"
import (
	"fmt"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
)

func newRicActionToBeSetupItemIesSingleContainer(rfItemIes *e2appducontents.RicactionToBeSetupItemIes) (*C.ProtocolIE_SingleContainer_1547P0_t, error) {
	pIeSC1547P0, err := newRicActionToBeSetupItemIEs(rfItemIes)

	return (*C.ProtocolIE_SingleContainer_1547P0_t)(pIeSC1547P0), err
}

func newRicActionAdmittedItemIEItemIesSingleContainer(raaItemIes *e2appducontents.RicactionAdmittedItemIes) (*C.ProtocolIE_SingleContainer_1547P1_t, error) {
	pIeSC1547P1, err := newRicActionAdmittedItemIEs(raaItemIes)

	return (*C.ProtocolIE_SingleContainer_1547P1_t)(pIeSC1547P1), err
}

func newRicActionNotAdmittedItemIEItemIesSingleContainer(ranaItemIes *e2appducontents.RicactionNotAdmittedItemIes) (*C.ProtocolIE_SingleContainer_1547P2_t, error) {
	pIeSC1547P2, err := newRicActionNotAdmittedItemIEs(ranaItemIes)

	return (*C.ProtocolIE_SingleContainer_1547P2_t)(pIeSC1547P2), err
}

func newRanFunctionItemIesSingleContainer(rfItemIes *e2appducontents.RanfunctionItemIes) (*C.ProtocolIE_SingleContainer_1547P3_t, error) {
	pIeSC1547P3, err := newRANfunctionItemIEs(rfItemIes)

	return (*C.ProtocolIE_SingleContainer_1547P3_t)(pIeSC1547P3), err
}

func newRanFunctionIDItemIesSingleContainer(rfIDItemIes *e2appducontents.RanfunctionIdItemIes) (*C.ProtocolIE_SingleContainer_1547P4_t, error) {
	pIeSC1547P4, err := newRANfunctionIDItemIEs(rfIDItemIes)

	return (*C.ProtocolIE_SingleContainer_1547P4_t)(pIeSC1547P4), err
}

func newRanFunctionIDcauseItemIesSingleContainer(rfIDcauseItemIes *e2appducontents.RanfunctionIdcauseItemIes) (*C.ProtocolIE_SingleContainer_1547P5_t, error) {
	pIeSC1547P5, err := newRANfunctionIDCauseItemIEs(rfIDcauseItemIes)

	return (*C.ProtocolIE_SingleContainer_1547P5_t)(pIeSC1547P5), err
}

func decodeRicActionToBeSetupItemIesSingleContainer(ratbsIeScC *C.ProtocolIE_SingleContainer_1547P0_t) (*e2appducontents.RicactionToBeSetupItemIes, error) {
	//fmt.Printf("Value %T %v\n", ratbsIeScC, ratbsIeScC)
	switch id := ratbsIeScC.id; id {
	case C.long(v1beta1.ProtocolIeIDRicactionToBeSetupItem):
		return decodeRicActionToBeSetupItemIes(&ratbsIeScC.value)
	default:
		return nil, fmt.Errorf("unexpected id for RicActionToBeSetupItem %v", C.long(id))
	}

}

func decodeRicActionAdmittedItemIesSingleContainer(raaiIeScC *C.ProtocolIE_SingleContainer_1547P1_t) (*e2appducontents.RicactionAdmittedItemIes, error) {
	//fmt.Printf("Value %T %v\n", raaiIeScC, raaiIeScC)
	switch id := raaiIeScC.id; id {
	case C.long(v1beta1.ProtocolIeIDRicactionAdmittedItem):
		return decodeRicActionAdmittedIDItemIes(&raaiIeScC.value)
	default:
		return nil, fmt.Errorf("unexpected id for RicactionAdmittedItemIes %v", C.long(id))
	}

}

func decodeRicActionNotAdmittedItemIesSingleContainer(ranaiIeScC *C.ProtocolIE_SingleContainer_1547P2_t) (*e2appducontents.RicactionNotAdmittedItemIes, error) {
	//fmt.Printf("Value %T %v\n", ranaiIeScC, ranaiIeScC)
	switch id := ranaiIeScC.id; id {
	case C.long(v1beta1.ProtocolIeIDRicactionNotAdmittedItem):
		return decodeRicActionNotAdmittedIDItemIes(&ranaiIeScC.value)
	default:
		return nil, fmt.Errorf("unexpected id for RicactionNotAdmittedItemIes %v", C.long(id))
	}

}

func decodeRanFunctionItemIesSingleContainer(rfiIeScC *C.ProtocolIE_SingleContainer_1547P3_t) (*e2appducontents.RanfunctionItemIes, error) {
	//fmt.Printf("Value %T %v\n", rfiIeScC, rfiIeScC)
	switch id := rfiIeScC.id; id {
	case C.long(v1beta1.ProtocolIeIDRanfunctionItem):
		return decodeRANfunctionItemIes(&rfiIeScC.value)
	default:
		return nil, fmt.Errorf("unexpected id for RanFunctionItem %v", C.long(id))
	}

}

func decodeRanFunctionIDItemIesSingleContainer(rfIDiIeScC *C.ProtocolIE_SingleContainer_1547P4_t) (*e2appducontents.RanfunctionIdItemIes, error) {
	//fmt.Printf("Value %T %v\n", rfIDiIeScC, rfIDiIeScC)
	switch id := rfIDiIeScC.id; id {
	case C.long(v1beta1.ProtocolIeIDRanfunctionIDItem):
		return decodeRANfunctionIDItemIes(&rfIDiIeScC.value)
	default:
		return nil, fmt.Errorf("unexpected id for RanfunctionIDItem %v", C.long(id))
	}

}

func decodeRanFunctionIDCauseItemIesSingleContainer(rfIDciIeScC *C.ProtocolIE_SingleContainer_1547P5_t) (*e2appducontents.RanfunctionIdcauseItemIes, error) {
	//fmt.Printf("Value %T %v\n", rfIDciIeScC, rfIDciIeScC)
	switch id := rfIDciIeScC.id; id {
	case C.long(v1beta1.ProtocolIeIDRanfunctionIeCauseItem):
		return decodeRANfunctionIDCauseItemIes(&rfIDciIeScC.value)
	default:
		return nil, fmt.Errorf("unexpected id for RanfunctionIeCauseItem %v", C.long(id))
	}

}
