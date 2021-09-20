// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package tranidpool

import (
	"math/big"

	"github.com/onosproject/onos-lib-go/pkg/errors"
)

// TransactionIDPool is the transaction ID pool.
type TransactionIDPool struct {
	pool *big.Int
	max  int
	min  int
}

// NewTransactionIDPool returns a new transaction ID Pool for given range.
func NewTransactionIDPool() *TransactionIDPool {
	return &TransactionIDPool{
		max:  256,
		min:  0,
		pool: big.NewInt(0),
	}
}

// NewID an id from the pool.
func (p *TransactionIDPool) NewID() (int, error) {
	for i := p.min; i < p.max; i++ {
		if p.pool.Bit(i) == 0 {
			p.pool.SetBit(p.pool, i, 1)
			return i, nil
		}
	}
	return -1, errors.NewUnavailable("all of the transaction IDs are assigned")
}

// Release an id back to the pool.
// Do nothing if the id is outside of the range.
func (p *TransactionIDPool) Release(id int) {
	if id >= p.min && id < p.max {
		p.pool.SetBit(p.pool, id, 0)
	}
}
