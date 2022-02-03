// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package addressing

import (
	"encoding/binary"
	"net"
)

// RICAddress RIC IP and port number
type RICAddress struct {
	IPAddress net.IP
	Port      uint64
}

// Port byte array representation of port
type Port struct {
	Value []byte
	Len   uint32
}

// ToUint converts port in byte array to uint based on given size
func (p *Port) ToUint() uint64 {
	var port uint64
	if p.Len == 16 {
		port = uint64(binary.BigEndian.Uint16(p.Value))
	} else if p.Len == 32 {
		port = uint64(binary.BigEndian.Uint32(p.Value))
	} else if p.Len == 64 {
		port = binary.BigEndian.Uint64(p.Value)
	}

	return port
}
