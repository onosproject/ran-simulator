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

// OperationID is an operation identifier
type OperationID string

// Executor executes primitive operations
type Executor interface {
	// RegisterUnaryOperation registers a unary primitive operation
	RegisterUnaryOperation(id OperationID, callback func([]byte) ([]byte, error))

	// RegisterStreamOperation registers a new primitive operation
	RegisterStreamOperation(id OperationID, callback func([]byte, Stream))

	// GetOperation returns an operation by name
	GetOperation(id OperationID) Operation
}

// Operation is the base interface for primitive operations
type Operation interface{}

// UnaryOperation is a primitive operation that returns a result
type UnaryOperation interface {
	// Execute executes the operation
	Execute(bytes []byte) ([]byte, error)
}

// StreamingOperation is a primitive operation that returns a stream
type StreamingOperation interface {
	// Execute executes the operation
	Execute(bytes []byte, stream Stream)
}

// newExecutor returns a new executor
func newExecutor() Executor {
	return &executor{
		operations: make(map[OperationID]Operation),
	}
}

// executor is an implementation of the Executor interface
type executor struct {
	Executor
	operations map[OperationID]Operation
}

func (e *executor) RegisterUnaryOperation(id OperationID, callback func([]byte) ([]byte, error)) {
	e.operations[id] = &unaryOperation{
		f: callback,
	}
}

func (e *executor) RegisterStreamOperation(id OperationID, callback func([]byte, Stream)) {
	e.operations[id] = &streamingOperation{
		f: callback,
	}
}

func (e *executor) GetOperation(id OperationID) Operation {
	return e.operations[id]
}

// unaryOperation is an implementation of the UnaryOperation interface
type unaryOperation struct {
	f func([]byte) ([]byte, error)
}

func (o *unaryOperation) Execute(bytes []byte) ([]byte, error) {
	return o.f(bytes)
}

// streamingOperation is an implementation of the StreamingOperation interface
type streamingOperation struct {
	f func([]byte, Stream)
}

func (o *streamingOperation) Execute(bytes []byte, stream Stream) {
	o.f(bytes, stream)
}
