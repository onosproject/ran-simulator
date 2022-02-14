// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

//go:build !race
// +build !race

package mobility

import (
	"context"
	"fmt"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/routes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDriver(t *testing.T) {
	m := &model.Model{}
	err := model.LoadConfig(m, "../model/test")
	assert.NoError(t, err)

	ns := nodes.NewNodeRegistry(m.Nodes)
	cs := cells.NewCellRegistry(m.Cells, ns)
	us := ues.NewUERegistry(1, cs, "random")
	rs := routes.NewRouteRegistry()

	ctx := context.TODO()
	ch := make(chan event.Event)
	err = us.Watch(ctx, ch, ues.WatchOptions{Replay: true})
	assert.NoError(t, err)

	e := <-ch
	ue := e.Value.(*model.UE)

	route := &model.Route{
		IMSI:     ue.IMSI,
		Points:   []*model.Coordinate{{Lat: 50.0001, Lng: 0.0000}, {Lat: 50.0000, Lng: 0.0000}, {Lat: 50.0000, Lng: 0.0002}},
		SpeedAvg: 40000.0,
	}
	err = rs.Add(ctx, route)
	assert.NoError(t, err)

	driver := NewMobilityDriver(cs, rs, us, "", "local", 15, false, false)
	tickUnit = time.Millisecond // For testing
	driver.Start(ctx)

	c := 0
	for e = range ch {
		ue = e.Value.(*model.UE)
		fmt.Printf("%v: %v\n", ue.Location, ue.Heading)
		c = c + 1
		if c > 10 {
			assert.Equal(t, 50.0, ue.Location.Lat)
			assert.Equal(t, 0.0, ue.Location.Lng)
			assert.Equal(t, uint32(180), ue.Heading)
			break
		} else if c == 6 {
			assert.Equal(t, uint32(270), ue.Heading)
		}
	}

	driver.Stop()
}

func TestRouteGeneration(t *testing.T) {
	m := &model.Model{}
	err := model.LoadConfig(m, "../utils/honeycomb/sample")
	assert.NoError(t, err)

	ns := nodes.NewNodeRegistry(m.Nodes)
	cs := cells.NewCellRegistry(m.Cells, ns)
	us := ues.NewUERegistry(1, cs, "random")
	rs := routes.NewRouteRegistry()

	ctx := context.TODO()
	us.SetUECount(ctx, 100)
	assert.Equal(t, 100, us.Len(ctx))

	driver := NewMobilityDriver(cs, rs, us, "", "local", 15, false, false)
	driver.GenerateRoutes(ctx, 30000, 160000, 20000, nil, false)
	assert.Equal(t, 100, rs.Len(ctx))

	ch := make(chan event.Event)
	err = us.Watch(ctx, ch, ues.WatchOptions{Replay: true})
	assert.NoError(t, err)

	tickUnit = time.Millisecond
	driver.Start(ctx)

	c := 0
	for e := range ch {
		ue := e.Value.(*model.UE)
		//fmt.Printf("%v: %v\n", ue.Location, ue.Heading)
		assert.True(t, 52.41 < ue.Location.Lat && ue.Location.Lat < 52.57, "UE latitude is out of range")
		assert.True(t, 13.29 < ue.Location.Lng && ue.Location.Lng < 13.52, "UE longitude is out of range")
		c = c + 1
		if c > 500 {
			break
		}
	}

	driver.Stop()
}
