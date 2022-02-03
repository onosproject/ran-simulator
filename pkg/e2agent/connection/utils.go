// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package connection

import (
	"github.com/onosproject/ran-simulator/pkg/e2agent/addressing"

	"time"

	"github.com/cenkalti/backoff/v4"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
)

func (e *e2Connection) getRICAddress(tnlInfo *e2apies.Tnlinformation) addressing.RICAddress {
	tnlAddr := tnlInfo.GetTnlAddress().GetValue()
	var ricAddress addressing.RICAddress

	if tnlInfo.GetTnlPort() != nil {
		tnlPort := tnlInfo.GetTnlPort().GetValue()
		tnlPortLen := tnlInfo.GetTnlPort().GetLen()
		p := &addressing.Port{
			Value: tnlPort,
			Len:   tnlPortLen,
		}
		ricAddress.Port = p.ToUint()

	} else {
		ricAddress.Port = e.ricAddress.Port
	}
	ricAddress.IPAddress = tnlAddr
	return ricAddress
}

const (
	backoffInterval = 10 * time.Millisecond
	maxBackoffTime  = 5 * time.Second
)

func newExpBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = backoffInterval
	// MaxInterval caps the RetryInterval
	b.MaxInterval = maxBackoffTime
	// Never stops retrying
	b.MaxElapsedTime = 0
	return b
}
