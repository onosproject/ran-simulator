// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/onosproject/ran-simulator/api/types"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestNodes(t *testing.T) {
	model := Model{}
	bytes, err := ioutil.ReadFile("test.yaml")
	assert.NoError(t, err)
	err = yaml.Unmarshal(bytes, &model)
	assert.NoError(t, err)
	t.Log(model)

	reg := NewNodeRegistry(model.Nodes)
	assert.Equal(t, 2, countNodes(reg))

	ch := make(chan NodeEvent)
	reg.WatchNodes(ch, WatchOptions{Replay: true, Monitor: true})

	event := <-ch
	assert.Equal(t, NONE, event.Type)
	event = <-ch
	assert.Equal(t, NONE, event.Type)

	_, err = reg.GetNode(144472)
	assert.True(t, err != nil, "node should not exist")

	go func() {
		_ = reg.AddNode(&Node{
			EnbID:         144472,
			Controllers:   []string{"controller1"},
			ServiceModels: []string{"kpm"},
			Cells:         make(map[string]Cell),
		})
	}()

	event, ok := <-ch
	assert.True(t, ok)
	assert.Equal(t, ADDED, event.Type)
	assert.Equal(t, 3, countNodes(reg))

	node, err := reg.GetNode(144472)
	assert.NoError(t, err, "node not found")
	assert.Equal(t, types.EnbID(144472), node.EnbID)

	go func() {
		err := reg.UpdateNode(&Node{
			EnbID:         144472,
			Controllers:   []string{"controller2"},
			ServiceModels: []string{"kpm"},
			Cells:         make(map[string]Cell),
		})
		assert.NoError(t, err, "node not updated")
	}()

	event, ok = <-ch
	assert.True(t, ok)
	assert.Equal(t, UPDATED, event.Type)

	go func() {
		n, err := reg.DeleteNode(types.EnbID(144472))
		assert.NoError(t, err, "node not deleted")
		assert.Equal(t, types.EnbID(144472), n.EnbID, "incorrect node deleted")
	}()

	event, ok = <-ch
	assert.True(t, ok)
	assert.Equal(t, DELETED, event.Type)
	assert.Equal(t, 2, countNodes(reg))

	err = reg.AddNode(&Node{
		EnbID:         144471,
		Controllers:   []string{"controller1"},
		ServiceModels: []string{"kpm"},
		Cells:         make(map[string]Cell),
	})
	assert.True(t, err != nil, "node should already exist")
	assert.Equal(t, 2, countNodes(reg))

	err = reg.UpdateNode(&Node{
		EnbID:         144472,
		Controllers:   []string{"controller1"},
		ServiceModels: []string{"kpm"},
		Cells:         make(map[string]Cell),
	})
	assert.True(t, err != nil, "node does not exist")

	_, err = reg.DeleteNode(144472)
	assert.True(t, err != nil, "node does not exist")
	assert.Equal(t, 2, countNodes(reg))

	close(ch)
}

func countNodes(reg NodeRegistry) int {
	c := 0
	ch := make(chan NodeEvent)
	reg.WatchNodes(ch, WatchOptions{Replay: true, Monitor: false})

	for range ch {
		c = c + 1
	}
	return c
}
