// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package routes

import (
	"context"
	"testing"

	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/event"
	"github.com/stretchr/testify/assert"
)

func TestRouteRegistry(t *testing.T) {
	ctx := context.Background()
	routes := NewRouteRegistry()
	assert.NotNil(t, routes, "unable to create route registry")
	assert.Equal(t, 0, routes.Len(ctx))

	ch := make(chan event.Event)
	err := routes.Watch(ctx, ch)
	assert.NoError(t, err)

	route1 := &model.Route{
		IMSI:   123456789,
		Points: []*model.Coordinate{{Lat: 1, Lng: 2}, {Lat: 2, Lng: 1}},
		Color:  "green",
	}

	route2 := &model.Route{
		IMSI:   923456781,
		Points: []*model.Coordinate{{Lat: 3, Lng: 2}, {Lat: 3, Lng: 1}},
		Color:  "blue",
	}

	err = routes.Add(ctx, route1)
	assert.NoError(t, err)
	assert.Equal(t, 1, routes.Len(ctx))

	err = routes.Add(ctx, route2)
	assert.NoError(t, err)
	assert.Equal(t, 2, routes.Len(ctx))

	nodeEvent := <-ch
	assert.Equal(t, Created, nodeEvent.Type.(RouteEvent))
	nodeEvent = <-ch
	assert.Equal(t, Created, nodeEvent.Type.(RouteEvent))

	list := routes.List(ctx)
	assert.Equal(t, 2, len(list))

	r1, err := routes.Get(ctx, route1.IMSI)
	assert.NoError(t, err)
	assert.Equal(t, route1.IMSI, r1.IMSI)
	_, err = routes.Delete(ctx, r1.IMSI)
	assert.NoError(t, err)
	nodeEvent = <-ch
	assert.Equal(t, Deleted, nodeEvent.Type.(RouteEvent))

	_, err = routes.Get(ctx, route1.IMSI)
	assert.Error(t, err, "route found")

	routes.Clear(ctx)
	assert.Equal(t, 0, routes.Len(ctx))
}
