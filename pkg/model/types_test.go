// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package model

import (
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

	assert.Equal(t, cellID, GetCellID(uint64(ecgi)), "incorrect CID")
	assert.Equal(t, plmnID, GetPlmnID(uint64(ecgi)), "incorrect PLMNID")
	assert.Equal(t, eci, GetECI(uint64(ecgi)), "incorrect ECI")
	assert.Equal(t, enbID, GetEnbID(uint64(ecgi)), "incorrect EnbID")
}
