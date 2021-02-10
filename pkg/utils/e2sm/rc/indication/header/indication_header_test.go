// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package header

import (
	"testing"

	"github.com/onosproject/ran-simulator/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationHeader(t *testing.T) {
	plmnID := types.NewUint24(12345)
	indicationHeader, err := NewIndicationHeader(WithPlmnID(plmnID.Value()),
		WithEutracellIdentity(32)).Build()
	assert.NoError(t, err)

	assert.Equal(t, indicationHeader.GetIndicationHeaderFormat1().Cgi.GetEUtraCgi().PLmnIdentity.Value, plmnID.ToBytes())
	assert.Equal(t, indicationHeader.GetIndicationHeaderFormat1().Cgi.GetEUtraCgi().EUtracellIdentity.Value.Value, uint64(32))

}
