// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "Cause.h"
//#include "CriticalityDiagnostics-IE-List.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newCause(cause *e2apies.Cause) (*C.Cause_t, error) {
	var pr C.Cause_PR

	choiceC := [8]byte{}

	switch causeType := cause.GetCause().(type) {
	case *e2apies.Cause_Misc:
		pr = C.Cause_PR_misc
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(cause.GetMisc()))
	case *e2apies.Cause_Protocol:
		pr = C.Cause_PR_protocol
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(cause.GetProtocol()))
	case *e2apies.Cause_RicRequest:
		pr = C.Cause_PR_ricRequest
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(cause.GetRicRequest()))
	case *e2apies.Cause_RicService:
		pr = C.Cause_PR_ricService
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(cause.GetRicService()))
	case *e2apies.Cause_Transport:
		pr = C.Cause_PR_transport
		binary.LittleEndian.PutUint64(choiceC[0:], uint64(cause.GetTransport()))
	default:
		return nil, fmt.Errorf("unexpected cause type %v", causeType)
	}

	causeC := C.Cause_t{
		present: pr,
		choice:  choiceC,
	}

	return &causeC, nil
}

func decodeCauseBytes(bytes []byte) (*e2apies.Cause, error) {
	causeC := C.Cause_t{
		present: C.Cause_PR(binary.LittleEndian.Uint64(bytes[:8])),
	}
	copy(causeC.choice[:8], bytes[8:16])

	return decodeCause(&causeC)
}

func decodeCause(causeC *C.Cause_t) (*e2apies.Cause, error) {
	cause := new(e2apies.Cause)
	switch causeC.present {
	case C.Cause_PR_misc:
		cause.Cause = &e2apies.Cause_Misc{
			Misc: e2apies.CauseMisc(binary.LittleEndian.Uint64(causeC.choice[:])),
		}
	case C.Cause_PR_protocol:
		cause.Cause = &e2apies.Cause_Protocol{
			Protocol: e2apies.CauseProtocol(binary.LittleEndian.Uint64(causeC.choice[:])),
		}
	case C.Cause_PR_ricRequest:
		cause.Cause = &e2apies.Cause_RicRequest{
			RicRequest: e2apies.CauseRic(binary.LittleEndian.Uint64(causeC.choice[:])),
		}
	case C.Cause_PR_ricService:
		cause.Cause = &e2apies.Cause_RicService{
			RicService: e2apies.CauseRicservice(binary.LittleEndian.Uint64(causeC.choice[:])),
		}
	case C.Cause_PR_transport:
		cause.Cause = &e2apies.Cause_Transport{
			Transport: e2apies.CauseTransport(binary.LittleEndian.Uint64(causeC.choice[:])),
		}
	default:
		return nil, fmt.Errorf("unexpected cause type %v", causeC.present)
	}
	return cause, nil
}
