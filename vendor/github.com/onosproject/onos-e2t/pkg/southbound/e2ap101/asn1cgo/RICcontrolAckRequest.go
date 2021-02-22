// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICcontrolAckRequest.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newRicControlAckRequest(rcar e2apies.RiccontrolAckRequest) (*C.RICcontrolAckRequest_t, error) {
	var ret C.RICcontrolAckRequest_t
	switch rcar {
	case e2apies.RiccontrolAckRequest_RICCONTROL_ACK_REQUEST_NO_ACK:
		ret = C.RICcontrolAckRequest_noAck
	case e2apies.RiccontrolAckRequest_RICCONTROL_ACK_REQUEST_ACK:
		ret = C.RICcontrolAckRequest_ack
	case e2apies.RiccontrolAckRequest_RICCONTROL_ACK_REQUEST_N_ACK:
		ret = C.RICcontrolAckRequest_nAck
	default:
		return nil, fmt.Errorf("unexpected RICcontrolAckRequest %v", rcar)
	}
	return &ret, nil
}

func decodeRicControlAckRequestBytes(bytes []byte) e2apies.RiccontrolAckRequest {
	rcarC := C.long(binary.LittleEndian.Uint64(bytes[:8]))
	return decodeRicControlAckRequest(&rcarC)
}

func decodeRicControlAckRequest(rcarC *C.RICcontrolAckRequest_t) e2apies.RiccontrolAckRequest {
	return e2apies.RiccontrolAckRequest(*rcarC)
}
