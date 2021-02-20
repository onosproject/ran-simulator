// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package ues

import (
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func cellStore(t *testing.T) cells.CellRegistry {
	m := model.Model{}
	bytes, err := ioutil.ReadFile("../../model/test.yaml")
	assert.NoError(t, err)
	err = yaml.Unmarshal(bytes, &m)
	assert.NoError(t, err)
	t.Log(m)
	return cells.NewCellRegistry(m.Cells, nodes.NewNodeRegistry(m.Nodes))
}

func TestUERegistry(t *testing.T) {
	ues := NewUERegistry(16, cellStore(t))
	assert.NotNil(t, ues, "unable to create UE registry")
	assert.Equal(t, 16, len(ues.ListAllUEs()))

	ues.SetUECount(10)
	assert.Equal(t, 10, len(ues.ListAllUEs()))

	ues.SetUECount(200)
	assert.Equal(t, 200, len(ues.ListAllUEs()))
}

func TestMoveUE(t *testing.T) {
	cellStore := cellStore(t)
	ues := NewUERegistry(18, cellStore)
	assert.NotNil(t, ues, "unable to create UE registry")

	// Get a cell ECGI
	ecgi1 := cellStore.GetRandomCell().ECGI

	// Get another cell ECGI; make sure it's different than the first.
	ecgi2 := cellStore.GetRandomCell().ECGI
	for ecgi1 == ecgi2 {
		ecgi2 = cellStore.GetRandomCell().ECGI
	}

	for i, ue := range ues.ListAllUEs() {
		ecgi := ecgi1
		if i%3 == 0 {
			ecgi = ecgi2
		}
		err := ues.MoveUE(ue.IMSI, ecgi, rand.Float64())
		assert.NoError(t, err)
	}

	assert.Equal(t, 12, len(ues.ListUEs(ecgi1)))
	assert.Equal(t, 6, len(ues.ListUEs(ecgi2)))
}
