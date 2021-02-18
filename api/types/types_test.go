// Copyright 2021-present Open Networking Foundation.
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

package types

import (
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestPlmnID(t *testing.T) {
	plmnID := PlmnID(101)
	ecgi := ToECGI(plmnID, ECI(0))
	assert.Equal(t, plmnID, GetPlmnID(uint64(ecgi)))

	ecgi = ToECGI(plmnID, ECI(0xfffffff))
	assert.Equal(t, plmnID, GetPlmnID(uint64(ecgi)))
}

func TestTypes(t *testing.T) {
	plmnID := PlmnID(221)
	cellID := CellID(192)
	enbID := EnbID(0xf8f8f)

	eci := ToECI(enbID, cellID)
	ecgi := ToECGI(plmnID, eci)
	genbID := ToGEnbID(plmnID, enbID)

	assert.Equal(t, cellID, GetCellID(uint64(ecgi)), "incorrect CID")
	assert.Equal(t, plmnID, GetPlmnID(uint64(ecgi)), "incorrect PLMNID")
	assert.Equal(t, eci, GetECI(uint64(ecgi)), "incorrect ECI")
	assert.Equal(t, enbID, GetEnbID(uint64(ecgi)), "incorrect ECGI EnbID")
	assert.Equal(t, enbID, GetEnbID(uint64(genbID)), "incorrect EnbID")
}

func TestSimValues(t *testing.T) {
	plmnID := PlmnID(314)
	enb1 := EnbID(144470)
	enb2 := EnbID(144471)
	ecgi11 := ToECGI(plmnID, ToECI(enb1, CellID(1)))
	ecgi12 := ToECGI(plmnID, ToECI(enb1, CellID(2)))
	ecgi21 := ToECGI(plmnID, ToECI(enb2, CellID(1)))
	ecgi22 := ToECGI(plmnID, ToECI(enb2, CellID(2)))

	fmt.Printf("%d\n%d\n%d\n%d\n", ecgi11, ecgi12, ecgi21, ecgi22)
}
