// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/cenkalti/backoff"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
)

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

// Port byte array representation of port
type Port struct {
	value []byte
	len   uint32
}

// ToUint converts port in byte array to uint based on given size
func (p *Port) ToUint() uint64 {
	var port uint64
	if p.len == 16 {
		port = uint64(binary.BigEndian.Uint16(p.value))
	} else if p.len == 32 {
		port = uint64(binary.BigEndian.Uint32(p.value))
	} else if p.len == 64 {
		port = binary.BigEndian.Uint64(p.value)
	}

	return port
}

func (e *e2Instance) getRICAddress(tnlInfo *e2apies.Tnlinformation) RICAddress {
	tnlAddr := tnlInfo.GetTnlAddress().GetValue()
	tnlAddrLen := tnlInfo.GetTnlAddress().GetLen()
	var ricAddress RICAddress
	if tnlInfo.GetTnlPort() != nil {
		tnlPort := tnlInfo.GetTnlPort().GetValue()
		tnlPortLen := tnlInfo.GetTnlPort().GetLen()
		p := &Port{
			value: tnlPort,
			len:   tnlPortLen,
		}
		ricAddress.port = p.ToUint()

	} else {
		ricAddress.port = e.ricAddress.port
	}
	if tnlAddrLen == net.IPv4len*8 {
		ricAddress.ipAddress = net.IPv4(tnlAddr[0], tnlAddr[1], tnlAddr[2], tnlAddr[3])
	}
	return ricAddress
}
