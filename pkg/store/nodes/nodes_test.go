// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package nodes

import (
	"context"
	"testing"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/event"

	"github.com/stretchr/testify/assert"
)

func TestNodes(t *testing.T) {
	m := &model.Model{}
	err := model.LoadConfig(m, "../../model/test")
	assert.NoError(t, err)
	t.Log(m)
	ctx := context.Background()

	nodeStore := NewNodeRegistry(m.Nodes)
	node1EnbID := types.EnbID(144472)
	node2EnbID := types.EnbID(144473)
	numNodes, err := nodeStore.Len(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 2, numNodes)

	ch := make(chan event.Event)
	err = nodeStore.Watch(ctx, ch)
	assert.NoError(t, err)
	node1 := &model.Node{
		EnbID:         node1EnbID,
		Controllers:   []string{"controller1"},
		ServiceModels: []string{"kpm"},
		Cells:         []types.ECGI{1234, 4321},
	}

	node2 := &model.Node{
		EnbID:         node2EnbID,
		Controllers:   []string{"controller1"},
		ServiceModels: []string{"kpm"},
		Cells:         []types.ECGI{5678, 8765},
	}
	err = nodeStore.Add(ctx, node1)
	assert.NoError(t, err)
	err = nodeStore.Add(ctx, node2)
	assert.NoError(t, err)

	nodeEvent := <-ch
	assert.Equal(t, Created, nodeEvent.Type.(NodeEvent))
	nodeEvent = <-ch
	assert.Equal(t, Created, nodeEvent.Type.(NodeEvent))

	node1, err = nodeStore.Get(ctx, node1EnbID)
	assert.NoError(t, err)
	assert.Equal(t, node1.EnbID, node1EnbID)
	_, err = nodeStore.Delete(ctx, node1EnbID)
	assert.NoError(t, err)
	nodeEvent = <-ch
	assert.Equal(t, Deleted, nodeEvent.Type.(NodeEvent))

	_, err = nodeStore.Get(ctx, node1EnbID)
	assert.Error(t, err, "node found")

	nodeStore.Clear(ctx)
	ids, _ := nodeStore.List(ctx)
	assert.Equal(t, 0, len(ids), "should be empty")
}
