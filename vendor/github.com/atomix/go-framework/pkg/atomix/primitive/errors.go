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

package primitive

import (
	"github.com/atomix/go-framework/pkg/atomix/errors"
)

// getStatus gets the proto status for the given error
func getStatus(err error) SessionResponseStatus {
	if err == nil {
		return SessionResponseStatus_OK
	}

	typed, ok := err.(*errors.TypedError)
	if !ok {
		return SessionResponseStatus_ERROR
	}

	switch typed.Type {
	case errors.Unknown:
		return SessionResponseStatus_UNKNOWN
	case errors.Canceled:
		return SessionResponseStatus_CANCELED
	case errors.NotFound:
		return SessionResponseStatus_NOT_FOUND
	case errors.AlreadyExists:
		return SessionResponseStatus_ALREADY_EXISTS
	case errors.Unauthorized:
		return SessionResponseStatus_UNAUTHORIZED
	case errors.Forbidden:
		return SessionResponseStatus_FORBIDDEN
	case errors.Conflict:
		return SessionResponseStatus_CONFLICT
	case errors.Invalid:
		return SessionResponseStatus_INVALID
	case errors.Unavailable:
		return SessionResponseStatus_UNAVAILABLE
	case errors.NotSupported:
		return SessionResponseStatus_NOT_SUPPORTED
	case errors.Timeout:
		return SessionResponseStatus_TIMEOUT
	case errors.Internal:
		return SessionResponseStatus_INTERNAL
	default:
		return SessionResponseStatus_ERROR
	}
}

// getMessage gets the message for the given error
func getMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
