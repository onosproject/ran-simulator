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
	"github.com/atomix/go-client/pkg/client/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Status gets the gRPC status for the given error
func Status(err error) *status.Status {
	if err == nil {
		return status.New(codes.OK, "")
	}

	typed, ok := err.(*TypedError)
	if !ok {
		return status.New(codes.Internal, err.Error())
	}

	switch typed.Type {
	case Unknown:
		return status.New(codes.Unknown, typed.Message)
	case Canceled:
		return status.New(codes.Canceled, typed.Message)
	case NotFound:
		return status.New(codes.NotFound, typed.Message)
	case AlreadyExists:
		return status.New(codes.AlreadyExists, typed.Message)
	case Unauthorized:
		return status.New(codes.Unauthenticated, typed.Message)
	case Forbidden:
		return status.New(codes.PermissionDenied, typed.Message)
	case Conflict:
		return status.New(codes.FailedPrecondition, typed.Message)
	case Invalid:
		return status.New(codes.InvalidArgument, typed.Message)
	case Unavailable:
		return status.New(codes.Unavailable, typed.Message)
	case NotSupported:
		return status.New(codes.Unimplemented, typed.Message)
	case Timeout:
		return status.New(codes.DeadlineExceeded, typed.Message)
	case Internal:
		return status.New(codes.Internal, typed.Message)
	default:
		return status.New(codes.Internal, err.Error())
	}
}

// FromStatus creates a typed error from a gRPC status
func FromStatus(status *status.Status) error {
	switch status.Code() {
	case codes.OK:
		return nil
	case codes.Unknown:
		return NewUnknown(status.Message())
	case codes.Canceled:
		return NewCanceled(status.Message())
	case codes.NotFound:
		return NewNotFound(status.Message())
	case codes.AlreadyExists:
		return NewAlreadyExists(status.Message())
	case codes.Unauthenticated:
		return NewUnauthorized(status.Message())
	case codes.PermissionDenied:
		return NewForbidden(status.Message())
	case codes.FailedPrecondition:
		return NewConflict(status.Message())
	case codes.InvalidArgument:
		return NewInvalid(status.Message())
	case codes.Unavailable:
		return NewUnavailable(status.Message())
	case codes.Unimplemented:
		return NewNotSupported(status.Message())
	case codes.DeadlineExceeded:
		return NewTimeout(status.Message())
	case codes.Internal:
		return NewInternal(status.Message())
	default:
		return NewUnknown(status.Message())
	}
}

// FromGRPC creates a typed error from a gRPC error
func FromGRPC(err error) error {
	if err == nil {
		return nil
	}

	stat, ok := status.FromError(err)
	if !ok {
		return New(Unknown, err.Error())
	}

	switch stat.Code() {
	case codes.OK:
		return nil
	case codes.Unknown:
		return New(Unknown, stat.Message())
	case codes.Canceled:
		return New(Canceled, stat.Message())
	case codes.NotFound:
		return New(NotFound, stat.Message())
	case codes.AlreadyExists:
		return New(AlreadyExists, stat.Message())
	case codes.Unauthenticated:
		return New(Unauthorized, stat.Message())
	case codes.PermissionDenied:
		return New(Forbidden, stat.Message())
	case codes.FailedPrecondition:
		return New(Conflict, stat.Message())
	case codes.InvalidArgument:
		return New(Invalid, stat.Message())
	case codes.Unavailable:
		return New(Unavailable, stat.Message())
	case codes.Unimplemented:
		return New(NotSupported, stat.Message())
	case codes.DeadlineExceeded:
		return New(Timeout, stat.Message())
	case codes.Internal:
		return New(Internal, stat.Message())
	default:
		return New(Unknown, stat.Message())
	}
}

// FromAtomix creates a typed error from an Atomix error
func FromAtomix(err error) error {
	if err == nil {
		return nil
	}

	switch errors.TypeOf(err) {
	case errors.Unknown:
		return New(Unknown, err.Error())
	case errors.Canceled:
		return New(Canceled, err.Error())
	case errors.NotFound:
		return New(NotFound, err.Error())
	case errors.AlreadyExists:
		return New(AlreadyExists, err.Error())
	case errors.Unauthorized:
		return New(Unauthorized, err.Error())
	case errors.Forbidden:
		return New(Forbidden, err.Error())
	case errors.Conflict:
		return New(Conflict, err.Error())
	case errors.Invalid:
		return New(Invalid, err.Error())
	case errors.Unavailable:
		return New(Unavailable, err.Error())
	case errors.NotSupported:
		return New(NotSupported, err.Error())
	case errors.Timeout:
		return New(Timeout, err.Error())
	case errors.Internal:
		return New(Internal, err.Error())
	default:
		return New(Unknown, err.Error())
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
func NewUnknown(msg string, args ...interface{}) error {
	return New(Unknown, msg, args...)
}

// NewCanceled returns a new Canceled error
func NewCanceled(msg string, args ...interface{}) error {
	return New(Canceled, msg, args...)
}

// NewNotFound returns a new NotFound error
func NewNotFound(msg string, args ...interface{}) error {
	return New(NotFound, msg, args...)
}

// NewAlreadyExists returns a new AlreadyExists error
func NewAlreadyExists(msg string, args ...interface{}) error {
	return New(AlreadyExists, msg, args...)
}

// NewUnauthorized returns a new Unauthorized error
func NewUnauthorized(msg string, args ...interface{}) error {
	return New(Unauthorized, msg, args...)
}

// NewForbidden returns a new Forbidden error
func NewForbidden(msg string, args ...interface{}) error {
	return New(Forbidden, msg, args...)
}

// NewConflict returns a new Conflict error
func NewConflict(msg string, args ...interface{}) error {
	return New(Conflict, msg, args...)
}

// NewInvalid returns a new Invalid error
func NewInvalid(msg string, args ...interface{}) error {
	return New(Invalid, msg, args...)
}

// NewUnavailable returns a new Unavailable error
func NewUnavailable(msg string, args ...interface{}) error {
	return New(Unavailable, msg, args...)
}

// NewNotSupported returns a new NotSupported error
func NewNotSupported(msg string, args ...interface{}) error {
	return New(NotSupported, msg, args...)
}

// NewTimeout returns a new Timeout error
func NewTimeout(msg string, args ...interface{}) error {
	return New(Timeout, msg, args...)
}

// NewInternal returns a new Internal error
func NewInternal(msg string, args ...interface{}) error {
	return New(Internal, msg, args...)
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
