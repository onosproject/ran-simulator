// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

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
	us := ues.NewUERegistry(1, cs)
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

	driver := NewMobilityDriver(rs, us)
	tickUnit = time.Millisecond // For testing
	driver.Start()

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

	close(ch)
	driver.Stop()
}
