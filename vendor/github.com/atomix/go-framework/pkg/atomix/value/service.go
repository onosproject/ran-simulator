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
	"github.com/atomix/go-framework/pkg/atomix/errors"
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/util"
	"github.com/golang/protobuf/proto"
	"io"
)

// Service is a state machine for a list primitive
type Service struct {
	primitive.Service
	value   []byte
	version uint64
}

// init initializes the list service
func (v *Service) init() {
	v.RegisterUnaryOperation(opSet, v.Set)
	v.RegisterUnaryOperation(opGet, v.Get)
	v.RegisterStreamOperation(opEvents, v.Events)
}

// Backup takes a snapshot of the service
func (v *Service) Backup(writer io.Writer) error {
	if err := util.WriteVarUint64(writer, v.version); err != nil {
		return err
	}
	if err := util.WriteBytes(writer, v.value); err != nil {
		return err
	}
	return nil
}

// Restore restores the service from a snapshot
func (v *Service) Restore(reader io.Reader) error {
	version, err := util.ReadVarUint64(reader)
	if err != nil {
		return err
	}
	v.version = version
	value, err := util.ReadBytes(reader)
	if err != nil {
		return err
	}
	v.value = value
	return nil
}

// Set sets the value
func (v *Service) Set(bytes []byte) ([]byte, error) {
	request := &SetRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	if request.ExpectVersion > 0 && request.ExpectVersion != v.version {
		return nil, errors.NewConflict("expected version %d does not match actual version %d", request.ExpectVersion, v.version)
	} else if request.ExpectValue != nil && len(request.ExpectValue) > 0 && (v.value == nil || !slicesEqual(v.value, request.ExpectValue)) {
		return nil, errors.NewConflict("expected value %v does not match actual value %v", request.ExpectValue, v.value)
	} else {
		prevValue := v.value
		prevVersion := v.version
		v.value = request.Value
		v.version++

		v.sendEvent(&ListenResponse{
			Type:            ListenResponse_UPDATED,
			PreviousValue:   prevValue,
			PreviousVersion: prevVersion,
			NewValue:        v.value,
			NewVersion:      v.version,
		})

		return proto.Marshal(&SetResponse{
			Version:   v.version,
			Succeeded: true,
		})
	}
}

// Get gets the current value
func (v *Service) Get(bytes []byte) ([]byte, error) {
	request := &GetRequest{}
	if err := proto.Unmarshal(bytes, request); err != nil {
		return nil, err
	}

	return proto.Marshal(&GetResponse{
		Value:   v.value,
		Version: v.version,
	})
}

// Events registers a channel on which to send events
func (v *Service) Events(bytes []byte, stream primitive.Stream) {
	// Keep the stream open for events
}

func (v *Service) sendEvent(event *ListenResponse) {
	bytes, err := proto.Marshal(event)
	for _, session := range v.Sessions() {
		for _, stream := range session.StreamsOf(opEvents) {
			stream.Result(bytes, err)
		}
	}
}

func slicesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for _, i := range a {
		for _, j := range b {
			if i != j {
				return false
			}
		}
	}
	return true
}
