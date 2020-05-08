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

package gnmi

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"gotest.tools/assert"
	"testing"
)

func Test_TypedValueToBytes(t *testing.T) {
	dec64Val1 := &gnmi.TypedValue{
		Value: &gnmi.TypedValue_IntVal{
			IntVal: 33,
		},
	}
	bytes, err := gnmiValueToBytes(dec64Val1)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(bytes))
	t.Logf("Bytes %v", bytes)

	decoded, err := bytesToGnmiValue(bytes)
	assert.NilError(t, err)
	assert.Equal(t, "int_val:33 ", decoded.String())
}

func Test_TypedValueToBytes2(t *testing.T) {
	typedValue := &gnmi.TypedValue{
		Value: &gnmi.TypedValue_FloatVal{
			FloatVal: 3.1,
		},
	}
	bytes, err := gnmiValueToBytes(typedValue)
	assert.NilError(t, err)
	assert.Equal(t, 5, len(bytes))
	t.Logf("Bytes %v", bytes)

	decoded, err := bytesToGnmiValue(bytes)
	assert.NilError(t, err)
	assert.Equal(t, "float_val:3.1 ", decoded.String())
}

func Test_TypedValueToBytes3(t *testing.T) {
	typedValue := &gnmi.TypedValue{
		Value: &gnmi.TypedValue_StringVal{
			StringVal: "This is a test",
		},
	}
	bytes, err := gnmiValueToBytes(typedValue)
	assert.NilError(t, err)
	assert.Equal(t, 16, len(bytes))
	t.Logf("Bytes %v", bytes)

	decoded, err := bytesToGnmiValue(bytes)
	assert.NilError(t, err)
	assert.Equal(t, `string_val:"This is a test" `, decoded.String())
}
