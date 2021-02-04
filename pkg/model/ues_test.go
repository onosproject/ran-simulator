// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package model

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUERegistry(t *testing.T) {
	ues := NewUERegistry(16)
	assert.NotNil(t, ues, "unable to create UE registry")
	assert.Equal(t, 16, len(ues.ListAllUEs()))

	ues.SetUECount(10)
	assert.Equal(t, 10, len(ues.ListAllUEs()))

	ues.SetUECount(200)
	assert.Equal(t, 200, len(ues.ListAllUEs()))
}

func TestMoveUE(t *testing.T) {
	ues := NewUERegistry(24)
	assert.NotNil(t, ues, "unable to create UE registry")

	ecgi1 := Ecgi("foo")
	ecgi2 := Ecgi("bar")

	for i, ue := range ues.ListAllUEs() {
		ecgi := ecgi1
		if i%3 == 0 {
			ecgi = ecgi2
		}
		ues.MoveUE(ue.Imsi, ecgi, rand.Float64())
	}

	assert.Equal(t, 16, len(ues.ListUEs(ecgi1)))
	assert.Equal(t, 8, len(ues.ListUEs(ecgi2)))
}
