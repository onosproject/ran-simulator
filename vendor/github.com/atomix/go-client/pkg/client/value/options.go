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

package value

import (
	api "github.com/atomix/api/proto/atomix/value"
)

// SetOption is an option for Set calls
type SetOption interface {
	beforeSet(request *api.SetRequest)
	afterSet(response *api.SetResponse)
}

// IfValue updates the value if the current value matches the given value
func IfValue(value []byte) SetOption {
	return valueOption{value}
}

type valueOption struct {
	value []byte
}

func (o valueOption) beforeSet(request *api.SetRequest) {
	request.ExpectValue = o.value
}

func (o valueOption) afterSet(response *api.SetResponse) {

}

// IfVersion updates the value if the version matches the given version
func IfVersion(version uint64) SetOption {
	return versionOption{version}
}

type versionOption struct {
	version uint64
}

func (o versionOption) beforeSet(request *api.SetRequest) {
	request.ExpectVersion = o.version
}

func (o versionOption) afterSet(response *api.SetResponse) {

}
