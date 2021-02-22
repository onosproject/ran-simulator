// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package ues

import (
	"context"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

func cellStore(t *testing.T) cells.Store {
	m := model.Model{}
	bytes, err := ioutil.ReadFile("../../model/test.yaml")
	assert.NoError(t, err)
	err = yaml.Unmarshal(bytes, &m)
	assert.NoError(t, err)
	t.Log(m)
	return cells.NewCellRegistry(m.Cells, nodes.NewNodeRegistry(m.Nodes))
}

func TestUERegistry(t *testing.T) {
	ctx := context.Background()
	ues := NewUERegistry(16, cellStore(t))
	assert.NotNil(t, ues, "unable to create UE registry")
	assert.Equal(t, 16, ues.Len(ctx))

	ues.SetUECount(ctx, 10)
	assert.Equal(t, 10, ues.Len(ctx))

	ues.SetUECount(ctx, 200)
	assert.Equal(t, 200, ues.Len(ctx))
}

func TestMoveUE(t *testing.T) {
	ctx := context.Background()
	cellStore := cellStore(t)
	ues := NewUERegistry(18, cellStore)
	assert.NotNil(t, ues, "unable to create UE registry")
	// Get a cell ECGI
	cell1, err := cellStore.GetRandomCell()
	assert.NoError(t, err)
	ecgi1 := cell1.ECGI

	// Get another cell ECGI; make sure it's different than the first.
	cell2, err := cellStore.GetRandomCell()
	assert.NoError(t, err)
	ecgi2 := cell2.ECGI
	for ecgi1 == ecgi2 {
		cell2, err = cellStore.GetRandomCell()
		assert.NoError(t, err)
		ecgi2 = cell2.ECGI
	}

	for i, ue := range ues.ListAllUEs(ctx) {
		ecgi := ecgi1
		if i%3 == 0 {
			ecgi = ecgi2
		}
		err := ues.Move(ctx, ue.IMSI, ecgi, rand.Float64())
		assert.NoError(t, err)
	}

	assert.Equal(t, 12, len(ues.ListUEs(ctx, ecgi1)))
	assert.Equal(t, 6, len(ues.ListUEs(ctx, ecgi2)))
}
