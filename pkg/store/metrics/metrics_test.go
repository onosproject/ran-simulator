// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"context"
	"testing"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	store := NewMetricsStore()
	ctx := context.Background()
	ch := make(chan event.Event)

	err := store.Watch(ctx, ch)
	assert.NoError(t, err)

	_, ok := store.Get(ctx, 123, "foo")
	assert.False(t, ok)
	entities, err := store.ListEntities(ctx)
	assert.Equal(t, 0, len(entities))
	assert.NoError(t, err)

	_ = store.Set(ctx, 123, "foo", 6.28)
	_ = store.Set(ctx, 123, "bar", 3.14)
	_ = store.Set(ctx, 123, "bah", "42")
	_ = store.Set(ctx, 321, "foo", 2.718)
	_ = store.Set(ctx, 321, "goo", 1.618)

	metricEvent := <-ch
	assert.Equal(t, Updated, metricEvent.Type)
	metricEvent = <-ch
	assert.Equal(t, Updated, metricEvent.Type)
	metricEvent = <-ch
	assert.Equal(t, Updated, metricEvent.Type)
	metricEvent = <-ch
	assert.Equal(t, Updated, metricEvent.Type)
	metricEvent = <-ch
	assert.Equal(t, Updated, metricEvent.Type)

	v, ok := store.Get(ctx, 123, "foo")
	assert.True(t, ok)
	assert.Equal(t, 6.28, v)

	entities, _ = store.ListEntities(ctx)
	assert.Equal(t, 2, len(entities), "incorrect entity count")
	metrics, _ := store.List(ctx, 123)
	assert.Equal(t, 3, len(metrics), "incorrect metric count")

	_ = store.Delete(ctx, 123, "bah")
	metrics, _ = store.List(ctx, 123)
	assert.Equal(t, 2, len(metrics), "incorrect metric count")

	metricEvent = <-ch
	assert.Equal(t, Deleted, metricEvent.Type)

	_, ok = store.Get(ctx, 123, "bah")
	assert.False(t, ok)

	_ = store.DeleteAll(ctx, 123)
	entities, _ = store.ListEntities(ctx)
	assert.Equal(t, 1, len(entities), "incorrect entity count")
	metrics, _ = store.List(ctx, 123)
	assert.Equal(t, 0, len(metrics), "incorrect metric count")

	metricEvent = <-ch
	assert.Equal(t, Deleted, metricEvent.Type)
	metricEvent = <-ch
	assert.Equal(t, Deleted, metricEvent.Type)
	assert.Equal(t, uint64(123), metricEvent.Key.(Key).EntityID)

	name := metricEvent.Key.(Key).Name
	assert.True(t, name == "bar" || name == "foo")

	store.Clear(ctx)
	ids, _ := store.ListEntities(ctx)
	assert.Equal(t, 0, len(ids), "should be empty")

	ctx.Done()
}
