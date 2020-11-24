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

package util

import (
	"io"
	"reflect"
)

// WriteBool writes a boolean to the given writer
func WriteBool(writer io.Writer, b bool) error {
	if b {
		if _, err := writer.Write([]byte{1}); err != nil {
			return err
		}
	} else {
		if _, err := writer.Write([]byte{0}); err != nil {
			return err
		}
	}
	return nil
}

// WriteVarUint64 writes an unsigned variable length integer to the given writer
func WriteVarUint64(writer io.Writer, x uint64) error {
	for x >= 0x80 {
		if _, err := writer.Write([]byte{byte(x) | 0x80}); err != nil {
			return err
		}
		x >>= 7
	}
	if _, err := writer.Write([]byte{byte(x)}); err != nil {
		return err
	}
	return nil
}

// WriteVarInt64 writes a signed variable length integer to the given writer
func WriteVarInt64(writer io.Writer, x int64) error {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return WriteVarUint64(writer, ux)
}

// WriteVarInt32 writes a signed 32-bit integer to the given writer
func WriteVarInt32(writer io.Writer, i int32) error {
	return WriteVarInt64(writer, int64(i))
}

// WriteVarInt writes a signed integer to the given writer
func WriteVarInt(writer io.Writer, i int) error {
	return WriteVarInt64(writer, int64(i))
}

// WriteVarUint32 writes an unsigned 32-bit integer to the given writer
func WriteVarUint32(writer io.Writer, i uint32) error {
	return WriteVarUint64(writer, uint64(i))
}

// WriteVarUint writes an unsigned integer to the given writer
func WriteVarUint(writer io.Writer, i uint) error {
	return WriteVarUint64(writer, uint64(i))
}

// WriteBytes writes a byte slice to the given writer
func WriteBytes(writer io.Writer, bytes []byte) error {
	if err := WriteVarInt(writer, len(bytes)); err != nil {
		return err
	}
	if _, err := writer.Write(bytes); err != nil {
		return err
	}
	return nil
}

// WriteValue writes the given value to the given writer
func WriteValue(writer io.Writer, value interface{}, f interface{}) error {
	funcVal := reflect.ValueOf(f)
	values := funcVal.Call([]reflect.Value{reflect.ValueOf(value)})
	bytesVal, errVal := values[0], values[1]
	if errVal.Interface() != nil {
		return errVal.Interface().(error)
	}
	return WriteBytes(writer, bytesVal.Interface().([]byte))
}

// WriteSlice writes a slice the the given writer
func WriteSlice(writer io.Writer, s interface{}, f interface{}) error {
	sliceVal := reflect.ValueOf(s)
	funcVal := reflect.ValueOf(f)
	for i := 0; i < sliceVal.Len(); i++ {
		values := funcVal.Call([]reflect.Value{sliceVal.Index(i)})
		bytesVal, errVal := values[0], values[1]
		if errVal.Interface() != nil {
			return errVal.Interface().(error)
		}
		bytes := bytesVal.Interface().([]byte)
		if err := WriteBytes(writer, bytes); err != nil {
			return err
		}
	}
	return nil
}

// WriteMap writes a map to the given writer
func WriteMap(writer io.Writer, m interface{}, f interface{}) error {
	mapVal := reflect.ValueOf(m)
	if err := WriteVarInt(writer, mapVal.Len()); err != nil {
		return err
	}

	funcVal := reflect.ValueOf(f)
	iter := mapVal.MapRange()
	for iter.Next() {
		values := funcVal.Call([]reflect.Value{iter.Key(), iter.Value()})
		bytesVal, errVal := values[0], values[1]
		if errVal.Interface() != nil {
			return errVal.Interface().(error)
		}
		bytes := bytesVal.Interface().([]byte)
		if err := WriteVarInt(writer, len(bytes)); err != nil {
			return err
		}
		if _, err := writer.Write(bytes); err != nil {
			return err
		}
	}
	return nil
}

// ReadBool reads a boolean from the given reader
func ReadBool(reader io.Reader) (bool, error) {
	bytes := make([]byte, 1)
	if _, err := reader.Read(bytes); err != nil {
		return false, err
	}
	return bytes[0] == 1, nil
}

// ReadVarUint64 reads an unsigned variable length integer from the given reader
func ReadVarUint64(reader io.Reader) (uint64, error) {
	var x uint64
	var s uint
	bytes := make([]byte, 1)
	for i := 0; i <= 9; i++ {
		if n, err := reader.Read(bytes); err != nil || n == -1 {
			return 0, err
		}
		b := bytes[0]
		if b < 0x80 {
			if i == 9 && b > 1 {
				return 0, nil
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, nil
}

// ReadVarInt64 reads a signed variable length integer from the given reader
func ReadVarInt64(reader io.Reader) (int64, error) {
	ux, n := ReadVarUint64(reader)
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, n
}

// ReadVarInt32 reads a signed 32-bit integer from the given reader
func ReadVarInt32(reader io.Reader) (int32, error) {
	i, err := ReadVarInt64(reader)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// ReadVarInt reads a signed integer from the given reader
func ReadVarInt(reader io.Reader) (int, error) {
	i, err := ReadVarInt64(reader)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// ReadVarUint32 reads an unsigned 32-bit integer from the given reader
func ReadVarUint32(reader io.Reader) (uint32, error) {
	i, err := ReadVarUint64(reader)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

// ReadVarUint reads an unsigned integer from the given reader
func ReadVarUint(reader io.Reader) (uint, error) {
	i, err := ReadVarUint64(reader)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

// ReadBytes reads a byte slice from the given reader
func ReadBytes(reader io.Reader) ([]byte, error) {
	length, err := ReadVarInt(reader)
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, length)
	if _, err := reader.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// ReadValue reads a value from the given reader
func ReadValue(reader io.Reader, f interface{}) (interface{}, error) {
	bytes, err := ReadBytes(reader)
	if err != nil {
		return nil, err
	}
	funcVal := reflect.ValueOf(f)
	values := funcVal.Call([]reflect.Value{reflect.ValueOf(bytes)})
	valVal, errVal := values[0], values[1]
	if errVal.Interface() != nil {
		return nil, errVal.Interface().(error)
	}
	return valVal.Interface(), nil
}

// ReadSlice reads a fixed length slice from the given reader
func ReadSlice(reader io.Reader, s interface{}, f interface{}) error {
	funcVal := reflect.ValueOf(f)
	sliceVal := reflect.ValueOf(s)
	for i := 0; i < sliceVal.Len(); i++ {
		bytes, err := ReadBytes(reader)
		if err != nil {
			return err
		}
		values := funcVal.Call([]reflect.Value{reflect.ValueOf(bytes)})
		valVal, errVal := values[0], values[1]
		if errVal.Interface() != nil {
			return errVal.Interface().(error)
		}
		sliceVal.Index(i).Set(valVal)
	}
	return nil
}

// ReadMap reads a map from the given reader
func ReadMap(reader io.Reader, m interface{}, f interface{}) error {
	length, err := ReadVarInt(reader)
	if err != nil {
		return err
	}

	funcVal := reflect.ValueOf(f)
	mapVal := reflect.ValueOf(m)
	for i := 0; i < length; i++ {
		length, err := ReadVarInt(reader)
		if err != nil {
			return err
		}
		bytes := make([]byte, length)
		if _, err := reader.Read(bytes); err != nil {
			return err
		}
		values := funcVal.Call([]reflect.Value{reflect.ValueOf(bytes)})
		keyVal, valVal, errVal := values[0], values[1], values[2]
		if errVal.Interface() != nil {
			return errVal.Interface().(error)
		}
		mapVal.SetMapIndex(keyVal, valVal)
	}
	return nil
}
