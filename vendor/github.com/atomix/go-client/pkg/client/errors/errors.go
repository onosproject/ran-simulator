// Copyright 2020-present Open Networking Foundation.
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

package errors

import (
	"fmt"
	"github.com/atomix/api/proto/atomix/headers"
)

// Type is an error type
type Type int

const (
	// Unknown is an unknown error type
	Unknown Type = iota
	// Canceled indicates a request context was canceled
	Canceled
	// NotFound indicates a resource was not found
	NotFound
	// AlreadyExists indicates a resource already exists
	AlreadyExists
	// Unauthorized indicates access to a resource is not authorized
	Unauthorized
	// Forbidden indicates the operation requested to be performed on a resource is forbidden
	Forbidden
	// Conflict indicates a conflict occurred during concurrent modifications to a resource
	Conflict
	// Invalid indicates a message or request is invalid
	Invalid
	// Unavailable indicates a service is not available
	Unavailable
	// NotSupported indicates a method is not supported
	NotSupported
	// Timeout indicates a request timed out
	Timeout
	// Internal indicates an unexpected internal error occurred
	Internal
)

// TypedError is an typed error
type TypedError struct {
	// Type is the error type
	Type Type
	// Message is the error message
	Message string
}

func (e *TypedError) Error() string {
	return e.Message
}

var _ error = &TypedError{}

// FromHeader creates a typed error from a response header
func FromHeader(header *headers.ResponseHeader) error {
	switch header.Status {
	case headers.ResponseStatus_OK:
		return nil
	case headers.ResponseStatus_ERROR:
		return NewUnknown(header.Message)
	case headers.ResponseStatus_UNKNOWN:
		return NewUnknown(header.Message)
	case headers.ResponseStatus_CANCELED:
		return NewCanceled(header.Message)
	case headers.ResponseStatus_NOT_FOUND:
		return NewNotFound(header.Message)
	case headers.ResponseStatus_ALREADY_EXISTS:
		return NewAlreadyExists(header.Message)
	case headers.ResponseStatus_UNAUTHORIZED:
		return NewUnauthorized(header.Message)
	case headers.ResponseStatus_FORBIDDEN:
		return NewForbidden(header.Message)
	case headers.ResponseStatus_CONFLICT:
		return NewConflict(header.Message)
	case headers.ResponseStatus_INVALID:
		return NewInvalid(header.Message)
	case headers.ResponseStatus_UNAVAILABLE:
		return NewUnavailable(header.Message)
	case headers.ResponseStatus_NOT_SUPPORTED:
		return NewNotSupported(header.Message)
	case headers.ResponseStatus_TIMEOUT:
		return NewTimeout(header.Message)
	case headers.ResponseStatus_INTERNAL:
		return NewInternal(header.Message)
	default:
		return NewUnknown(header.Message)
	}
}

// New creates a new typed error
func New(t Type, msg string, args ...interface{}) error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &TypedError{
		Type:    t,
		Message: msg,
	}
}

// NewUnknown returns a new Unknown error
func NewUnknown(msg string) error {
	return New(Unknown, msg)
}

// NewCanceled returns a new Canceled error
func NewCanceled(msg string) error {
	return New(Canceled, msg)
}

// NewNotFound returns a new NotFound error
func NewNotFound(msg string) error {
	return New(NotFound, msg)
}

// NewAlreadyExists returns a new AlreadyExists error
func NewAlreadyExists(msg string) error {
	return New(AlreadyExists, msg)
}

// NewUnauthorized returns a new Unauthorized error
func NewUnauthorized(msg string) error {
	return New(Unauthorized, msg)
}

// NewForbidden returns a new Forbidden error
func NewForbidden(msg string) error {
	return New(Forbidden, msg)
}

// NewConflict returns a new Conflict error
func NewConflict(msg string) error {
	return New(Conflict, msg)
}

// NewInvalid returns a new Invalid error
func NewInvalid(msg string) error {
	return New(Invalid, msg)
}

// NewUnavailable returns a new Unavailable error
func NewUnavailable(msg string) error {
	return New(Unavailable, msg)
}

// NewNotSupported returns a new NotSupported error
func NewNotSupported(msg string) error {
	return New(NotSupported, msg)
}

// NewTimeout returns a new Timeout error
func NewTimeout(msg string) error {
	return New(Timeout, msg)
}

// NewInternal returns a new Internal error
func NewInternal(msg string) error {
	return New(Internal, msg)
}

// TypeOf returns the type of the given error
func TypeOf(err error) Type {
	if typed, ok := err.(*TypedError); ok {
		return typed.Type
	}
	return Unknown
}

// IsType checks whether the given error is of the given type
func IsType(err error, t Type) bool {
	if typed, ok := err.(*TypedError); ok {
		return typed.Type == t
	}
	return false
}

// IsUnknown checks whether the given error is an Unknown error
func IsUnknown(err error) bool {
	return IsType(err, Unknown)
}

// IsCanceled checks whether the given error is an Canceled error
func IsCanceled(err error) bool {
	return IsType(err, Canceled)
}

// IsNotFound checks whether the given error is a NotFound error
func IsNotFound(err error) bool {
	return IsType(err, NotFound)
}

// IsAlreadyExists checks whether the given error is a AlreadyExists error
func IsAlreadyExists(err error) bool {
	return IsType(err, AlreadyExists)
}

// IsUnauthorized checks whether the given error is a Unauthorized error
func IsUnauthorized(err error) bool {
	return IsType(err, Unauthorized)
}

// IsForbidden checks whether the given error is a Forbidden error
func IsForbidden(err error) bool {
	return IsType(err, Forbidden)
}

// IsConflict checks whether the given error is a Conflict error
func IsConflict(err error) bool {
	return IsType(err, Conflict)
}

// IsInvalid checks whether the given error is an Invalid error
func IsInvalid(err error) bool {
	return IsType(err, Invalid)
}

// IsUnavailable checks whether the given error is an Unavailable error
func IsUnavailable(err error) bool {
	return IsType(err, Unavailable)
}

// IsNotSupported checks whether the given error is a NotSupported error
func IsNotSupported(err error) bool {
	return IsType(err, NotSupported)
}

// IsTimeout checks whether the given error is a Timeout error
func IsTimeout(err error) bool {
	return IsType(err, Timeout)
}

// IsInternal checks whether the given error is an Internal error
func IsInternal(err error) bool {
	return IsType(err, Internal)
}
