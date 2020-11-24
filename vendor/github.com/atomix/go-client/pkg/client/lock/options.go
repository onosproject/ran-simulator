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

package lock

import (
	api "github.com/atomix/api/proto/atomix/lock"
	"time"
)

// LockOption is an option for Lock calls
//nolint:golint
type LockOption interface {
	beforeLock(request *api.LockRequest)
	afterLock(response *api.LockResponse)
}

// WithTimeout sets the lock timeout
func WithTimeout(timeout time.Duration) LockOption {
	return timeoutOption{timeout: timeout}
}

type timeoutOption struct {
	timeout time.Duration
}

func (o timeoutOption) beforeLock(request *api.LockRequest) {
	request.Timeout = &o.timeout
}

func (o timeoutOption) afterLock(response *api.LockResponse) {

}

// UnlockOption is an option for Unlock calls
type UnlockOption interface {
	beforeUnlock(request *api.UnlockRequest)
	afterUnlock(response *api.UnlockResponse)
}

// IsLockedOption is an option for IsLocked calls
type IsLockedOption interface {
	beforeIsLocked(request *api.IsLockedRequest)
	afterIsLocked(response *api.IsLockedResponse)
}

// IfVersion sets the lock version to check
func IfVersion(version uint64) IfVersionOption {
	return IfVersionOption{version: version}
}

// IfVersionOption is a lock option for checking the version
type IfVersionOption struct {
	version uint64
}

func (o IfVersionOption) beforeUnlock(request *api.UnlockRequest) {
	request.Version = o.version
}

func (o IfVersionOption) afterUnlock(response *api.UnlockResponse) {

}

func (o IfVersionOption) beforeIsLocked(request *api.IsLockedRequest) {
	request.Version = o.version
}

func (o IfVersionOption) afterIsLocked(response *api.IsLockedResponse) {

}
