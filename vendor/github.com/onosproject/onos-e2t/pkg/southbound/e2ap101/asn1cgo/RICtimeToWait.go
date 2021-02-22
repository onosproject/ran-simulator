// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICtimeToWait.h"
import "C"
import (
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newRicTimeToWait(rttw e2apies.RictimeToWait) (*C.RICtimeToWait_t, error) {
	var ret C.RICtimeToWait_t
	switch rttw {
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_ZERO:
		ret = C.RICtimeToWait_zero
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W1MS:
		ret = C.RICtimeToWait_w1ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W2MS:
		ret = C.RICtimeToWait_w2ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W5MS:
		ret = C.RICtimeToWait_w5ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W10MS:
		ret = C.RICtimeToWait_w10ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W20MS:
		ret = C.RICtimeToWait_w20ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W30MS:
		ret = C.RICtimeToWait_w30ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W40MS:
		ret = C.RICtimeToWait_w40ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W50MS:
		ret = C.RICtimeToWait_w50ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W100MS:
		ret = C.RICtimeToWait_w100ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W200MS:
		ret = C.RICtimeToWait_w200ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W500MS:
		ret = C.RICtimeToWait_w500ms
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W1S:
		ret = C.RICtimeToWait_w1s
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W2S:
		ret = C.RICtimeToWait_w2s
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W5S:
		ret = C.RICtimeToWait_w5s
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W10S:
		ret = C.RICtimeToWait_w10s
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W20S:
		ret = C.RICtimeToWait_w20s
	case e2apies.RictimeToWait_RICTIME_TO_WAIT_W60S:
		ret = C.RICtimeToWait_w60s
	default:
		return nil, fmt.Errorf("unexpected RICtimeToWait %v", rttw)
	}
	return &ret, nil
}

//func decodeRicTimeToWaitBytes(bytes []byte) e2apies.RictimeToWait {
//	rttwC := C.long(binary.LittleEndian.Uint64(bytes[:8]))
//	return decodeRicTimeToWait(&rttwC)
//}

func decodeRicTimeToWait(rttwC *C.RICtimeToWait_t) e2apies.RictimeToWait {
	return e2apies.RictimeToWait(*rttwC)
}
