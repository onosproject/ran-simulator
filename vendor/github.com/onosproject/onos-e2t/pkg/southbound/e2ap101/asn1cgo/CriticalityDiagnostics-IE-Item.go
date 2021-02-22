// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "CriticalityDiagnostics-IE-Item.h"
import "C"
import (
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newCriticalityDiagnosticsIEItem(cdIeI *e2apies.CriticalityDiagnosticsIeItem) (*C.CriticalityDiagnostics_IE_Item_t, error) {
	c, err := criticalityToC(cdIeI.GetIEcriticality())
	if err != nil {
		return nil, fmt.Errorf("criticalityToC() %s", err.Error())
	}
	p, err := newprotocolIeID(cdIeI.GetIEId())
	if err != nil {
		return nil, fmt.Errorf("protocolIeIDToC() %s", err.Error())
	}
	t, err := newTypeOfError(cdIeI.GetTypeOfError())
	if err != nil {
		return nil, fmt.Errorf("newTypeOfError() %s", err.Error())
	}
	cdIeIC := C.CriticalityDiagnostics_IE_Item_t{
		iECriticality: c,
		iE_ID:         p,
		typeOfError:   t,
	}
	return &cdIeIC, nil
}

func decodeCriticalityDiagnosticsIEItem(cdIC *C.CriticalityDiagnostics_IE_Item_t) (*e2apies.CriticalityDiagnosticsIeItem, error) {
	result := e2apies.CriticalityDiagnosticsIeItem{
		IEcriticality: decodeCriticality(cdIC.iECriticality),
		IEId:          decodeProtocolIeID(cdIC.iE_ID),
		TypeOfError:   decodeTypeOfError(cdIC.typeOfError),
	}

	return &result, nil
}
