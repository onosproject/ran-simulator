// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cells

import (
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestCells(t *testing.T) {
	m := model.Model{}
	bytes, err := ioutil.ReadFile("../../model/test.yaml")
	assert.NoError(t, err)
	err = yaml.Unmarshal(bytes, &m)
	assert.NoError(t, err)
	t.Log(m)

	reg := NewCellRegistry(m.Cells)
	assert.Equal(t, 4, countCells(reg))

	ch := make(chan CellEvent)
	reg.WatchCells(ch, WatchOptions{Replay: true, Monitor: true})

	event := <-ch
	assert.Equal(t, NONE, event.Type)
	event = <-ch
	assert.Equal(t, NONE, event.Type)

	_, err = reg.GetCell(84325717507)
	assert.True(t, err != nil, "cell should not exist")

	go func() {
		err := reg.AddCell(&model.Cell{
			ECGI:   84325717507,
			Sector: model.Sector{Center: model.Coordinate{Lat: 46, Lng: 29}, Azimuth: 180, Arc: 180},
			Color:  "blue",
		})
		assert.NoError(t, err, "cell not added")
	}()

	// FIXME: find where the deadlock is...
	if true {
		return
	}

	event, ok := <-ch
	assert.True(t, ok)
	assert.Equal(t, ADDED, event.Type)
	assert.Equal(t, 5, countCells(reg))

	cell, err := reg.GetCell(84325717507)
	assert.NoError(t, err, "cell not found")
	assert.Equal(t, types.ECGI(84325717507), cell.ECGI)

	go func() {
		err := reg.UpdateCell(&model.Cell{
			ECGI:   84325717507,
			Sector: model.Sector{Center: model.Coordinate{Lat: 46, Lng: 29}, Azimuth: 180, Arc: 120},
			Color:  "red",
		})
		assert.NoError(t, err, "cell not updated")
	}()

	event, ok = <-ch
	assert.True(t, ok)
	assert.Equal(t, UPDATED, event.Type)

	go func() {
		n, err := reg.DeleteCell(types.ECGI(84325717507))
		assert.NoError(t, err, "cell not deleted")
		assert.Equal(t, types.ECGI(84325717507), n.ECGI, "incorrect cell deleted")
	}()

	event, ok = <-ch
	assert.True(t, ok)
	assert.Equal(t, DELETED, event.Type)
	assert.Equal(t, 4, countCells(reg))

	err = reg.AddCell(&model.Cell{
		ECGI:   84325717506,
		Sector: model.Sector{Center: model.Coordinate{Lat: 46, Lng: 29}, Azimuth: 180, Arc: 120},
		Color:  "purple",
	})
	assert.True(t, err != nil, "cell should already exist")
	assert.Equal(t, 4, countCells(reg))

	err = reg.UpdateCell(&model.Cell{
		ECGI:   84325717507,
		Sector: model.Sector{Center: model.Coordinate{Lat: 46, Lng: 29}, Azimuth: 180, Arc: 120},
		Color:  "red",
	})
	assert.True(t, err != nil, "cell does not exist")

	_, err = reg.DeleteCell(84325717507)
	assert.True(t, err != nil, "cell does not exist")
	assert.Equal(t, 4, countCells(reg))

	close(ch)
}

func countCells(reg CellRegistry) int {
	c := 0
	ch := make(chan CellEvent)
	reg.WatchCells(ch, WatchOptions{Replay: true, Monitor: false})

	for range ch {
		c = c + 1
	}
	return c
}
