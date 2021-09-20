// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package tranidpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransIDPool(t *testing.T) {
	pool := NewTransactionIDPool()
	for i := 0; i < 256; i++ {
		id, err := pool.NewID()
		assert.NoError(t, err)
		assert.Equal(t, i, id)
	}
	id, err := pool.NewID()
	assert.NotNil(t, err)
	assert.Equal(t, -1, id)
	pool.Release(0)
	id, err = pool.NewID()
	assert.NoError(t, err)
	assert.Equal(t, 0, id)
}
