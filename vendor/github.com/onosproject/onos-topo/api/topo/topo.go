// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topo

import (
	fmt "fmt"
	"strconv"
	"strings"
	"time"

	"github.com/onosproject/onos-topo/api/device"
	grpc "google.golang.org/grpc"
)

// ID ...
type ID string

// NullID ...
const NullID = ""

// Revision is an object revision
type Revision uint64

// Attribute keys
const (
	Address     = "address"
	Target      = "target"
	Type        = "type"
	Version     = "version"
	Timeout     = "timeout"
	Role        = "role"
	Displayname = "displayname"
	User        = "user"
	Password    = "password"
	TLSCaCert   = "tls-ca-cert"
	TLSCert     = "tls-cert"
	TLSKey      = "tls-key"
	TLSPlain    = "tls-plain"
	TLSInsecure = "tls-insecure"
)

// TopoClientFactory : Default EntityServiceClient creation.
var TopoClientFactory = func(cc *grpc.ClientConn) TopoClient {
	return NewTopoClient(cc)
}

// CreateTopoClient creates and returns a new topo device client
func CreateTopoClient(cc *grpc.ClientConn) TopoClient {
	return TopoClientFactory(cc)
}

// ObjectToDevice ...
func ObjectToDevice(obj *Object) *device.Device {
	if obj == nil || obj.Type != Object_ENTITY { // Device is an entity
		return nil
	}

	d := &device.Device{
		ID:        device.ID(obj.ID),
		Protocols: []*device.ProtocolState{},
	}

	d.Address = obj.Attributes[Address]
	d.Target = obj.Attributes[Target]
	d.Version = obj.Attributes[Version]
	t, err := strconv.Atoi(obj.Attributes[Timeout])
	var timeout time.Duration
	if err == nil {
		timeout = time.Duration(t) * time.Second
	} else {
		timeout = time.Duration(0) * time.Second
	}
	d.Timeout = &timeout
	d.Role = device.Role(obj.Attributes[Role])
	d.Displayname = obj.Attributes[Displayname]
	d.Credentials.User = obj.Attributes[User]
	d.Credentials.Password = obj.Attributes[Password]
	d.TLS.CaCert = obj.Attributes[TLSCaCert]
	d.TLS.Cert = obj.Attributes[TLSCert]
	d.TLS.Key = obj.Attributes[TLSKey]
	if strings.ToLower(obj.Attributes[TLSPlain]) == "true" {
		d.TLS.Plain = true
	} else if strings.ToLower(obj.Attributes[TLSPlain]) == "false" {
		d.TLS.Plain = false
	}
	if strings.ToLower(obj.Attributes[TLSInsecure]) == "true" {
		d.TLS.Insecure = true
	} else if strings.ToLower(obj.Attributes[TLSInsecure]) == "false" {
		d.TLS.Insecure = false
	}
	entity := obj.GetEntity()
	if entity != nil {
		d.Protocols = obj.GetEntity().Protocols
		d.Type = device.Type(obj.GetEntity().KindID)
	}

	return d
}

// DeviceToObject ...
func DeviceToObject(d *device.Device) *Object {
	obj := &Object{
		ID:         ID(d.ID),
		Type:       Object_ENTITY,
		Attributes: map[string]string{},
		Obj: &Object_Entity{
			Entity: &Entity{
				Protocols: []*device.ProtocolState{},
			},
		},
	}
	obj.Attributes[Address] = d.Address
	obj.Attributes[Target] = d.Target
	obj.Attributes[Version] = d.Version
	obj.Attributes[Timeout] = fmt.Sprintf("%f", d.Timeout.Seconds())
	obj.Attributes[Role] = string(d.Role)
	obj.Attributes[Displayname] = d.Displayname
	obj.Attributes[User] = d.Credentials.User
	obj.Attributes[Password] = d.Credentials.Password
	obj.Attributes[TLSCaCert] = d.TLS.CaCert
	obj.Attributes[TLSCert] = d.TLS.Cert
	obj.Attributes[TLSKey] = d.TLS.Key
	if d.TLS.Plain {
		obj.Attributes[TLSPlain] = "true"
	} else {
		obj.Attributes[TLSPlain] = "false"
	}
	if d.TLS.Insecure {
		obj.Attributes[TLSInsecure] = "true"
	} else {
		obj.Attributes[TLSInsecure] = "false"
	}
	obj.GetEntity().Protocols = d.Protocols
	obj.GetEntity().KindID = ID(d.Type)

	return obj
}
