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

package device

import "google.golang.org/grpc"

// ID is a device ID
type ID string

// Type is a device type
type Type string

// Role is a device role
type Role string

// Revision is the device revision number
type Revision uint64

// DeviceServiceClientFactory : Default DeviceServiceClient creation.
var DeviceServiceClientFactory = func(cc *grpc.ClientConn) DeviceServiceClient {
	return NewDeviceServiceClient(cc)
}

// CreateDeviceServiceClient creates and returns a new topo device client
func CreateDeviceServiceClient(cc *grpc.ClientConn) DeviceServiceClient {
	return DeviceServiceClientFactory(cc)
}
