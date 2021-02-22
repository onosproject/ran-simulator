// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

// #include "Criticality.h"
import "C"
import (
	"fmt"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
)

func criticalityToC(criticality e2ap_commondatatypes.Criticality) (C.Criticality_t, error) {
	var critC C.Criticality_t
	switch criticality {
	case e2ap_commondatatypes.Criticality_CRITICALITY_REJECT:
		critC = C.Criticality_reject
	case e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE:
		critC = C.Criticality_ignore
	case e2ap_commondatatypes.Criticality_CRITICALITY_NOTIFY:
		critC = C.Criticality_notify
	default:
		return C.Criticality_t(-1), fmt.Errorf("unexpected value for criticality %d", criticality)
	}

	return critC, nil
}

func decodeCriticality(criticalityC C.Criticality_t) e2ap_commondatatypes.Criticality {
	return e2ap_commondatatypes.Criticality(criticalityC)
}
