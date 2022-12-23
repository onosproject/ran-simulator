// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cells

import (
	"context"
	"os"
	"testing"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/store/event"

	"github.com/onosproject/ran-simulator/pkg/store/nodes"

	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestCells(t *testing.T) {
	m := model.Model{}
	bytes, err := os.ReadFile("../../model/test.yaml")
	assert.NoError(t, err)
	err = yaml.Unmarshal(bytes, &m)
	assert.NoError(t, err)
	t.Log(m)
	ctx := context.Background()

	cellStore := NewCellRegistry(m.Cells, nodes.NewNodeRegistry(m.Nodes))
	ch := make(chan event.Event)
	err = cellStore.Watch(ctx, ch, WatchOptions{Replay: false, Monitor: false})
	assert.NoError(t, err)
	cell, err := cellStore.Get(ctx, 84325717505)
	assert.NoError(t, err)
	assert.Equal(t, types.NCGI(84325717505), cell.NCGI)

	ecgi1 := types.NCGI(84325717507)
	cell1 := &model.Cell{
		NCGI:   ecgi1,
		Sector: model.Sector{Center: model.Coordinate{Lat: 46, Lng: 29}, Azimuth: 180, Arc: 180, Height: 30, Tilt: -10},
		Color:  "blue"}

	err = cellStore.Add(ctx, cell1)
	assert.NoError(t, err)

	cellEvent := <-ch
	assert.Equal(t, Created, cellEvent.Type)

	cell1, err = cellStore.Get(ctx, ecgi1)
	assert.NoError(t, err)
	assert.Equal(t, ecgi1, cell1.NCGI)

	_, err = cellStore.Delete(ctx, ecgi1)
	assert.NoError(t, err)
	cellEvent = <-ch
	assert.Equal(t, Deleted, cellEvent.Type)

	cellStore.Clear(ctx)
	ids, _ := cellStore.List(ctx)
	assert.Equal(t, 0, len(ids), "should be empty")
}
