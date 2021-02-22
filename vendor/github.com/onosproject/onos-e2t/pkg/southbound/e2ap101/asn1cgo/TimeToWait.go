// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "TimeToWait.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

func newTimeToWait(ttw e2apies.TimeToWait) (*C.TimeToWait_t, error) {
	var ret C.TimeToWait_t
	switch ttw {
	case e2apies.TimeToWait_TIME_TO_WAIT_V1S:
		ret = C.TimeToWait_v1s
	case e2apies.TimeToWait_TIME_TO_WAIT_V2S:
		ret = C.TimeToWait_v2s
	case e2apies.TimeToWait_TIME_TO_WAIT_V5S:
		ret = C.TimeToWait_v5s
	case e2apies.TimeToWait_TIME_TO_WAIT_V10S:
		ret = C.TimeToWait_v10s
	case e2apies.TimeToWait_TIME_TO_WAIT_V20S:
		ret = C.TimeToWait_v20s
	case e2apies.TimeToWait_TIME_TO_WAIT_V60S:
		ret = C.TimeToWait_v60s
	default:
		return nil, fmt.Errorf("unexpected TimeToWait %v", ttw)
	}
	return &ret, nil
}

func decodeTimeToWaitBytes(bytes []byte) e2apies.TimeToWait {
	ttwC := C.long(binary.LittleEndian.Uint64(bytes[:8]))
	return decodeTimeToWait(&ttwC)
}

func decodeTimeToWait(ttwC *C.TimeToWait_t) e2apies.TimeToWait {
	return e2apies.TimeToWait(*ttwC)
}
