// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package pciload

import (
	"context"
	"testing"

	"github.com/onosproject/ran-simulator/pkg/store/metrics"

	"github.com/stretchr/testify/assert"
)

func TestPCILoad(t *testing.T) {
	ctx := context.TODO()
	store := metrics.NewMetricsStore()
	err := LoadPCIMetrics(store)
	assert.NoError(t, err)

	v, ok := store.Get(ctx, 123, "pci")
	assert.True(t, ok)
	assert.Equal(t, uint32(42), v)

	v, ok = store.Get(ctx, 123, "pcipool")
	assert.True(t, ok)
	pciPool := v.([]PciRange)
	assert.Equal(t, uint32(90), pciPool[1].Max)

	v, ok = store.Get(ctx, 213, "pci")
	assert.True(t, ok)
	assert.Equal(t, uint32(69), v)

	v, ok = store.Get(ctx, 213, "earfcn")
	assert.True(t, ok)
	assert.Equal(t, uint32(7213), v)
}

func TestPCISampleLoad(t *testing.T) {
	ctx := context.TODO()
	store := metrics.NewMetricsStore()
	err := LoadPCIMetricsConfig(store, "sample")
	assert.NoError(t, err)

	v, ok := store.Get(ctx, 21458294227474, "pci")
	assert.True(t, ok)
	assert.Equal(t, uint32(459), v)
}
