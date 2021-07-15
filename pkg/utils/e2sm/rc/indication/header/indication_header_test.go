// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package header

import (
	"encoding/binary"
	"testing"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationHeader(t *testing.T) {
	plmnID := ransimtypes.NewUint24(12345)
	indicationHeader, err := NewIndicationHeader(WithPlmnID(plmnID.Value()),
		WithNRcellIdentity(32)).Build()
	assert.NoError(t, err)

	assert.Equal(t, plmnID.ToBytes(), indicationHeader.GetIndicationHeaderFormat1().Cgi.GetNrCgi().PLmnIdentity.Value)
	assert.Equal(t, uint64(32), binary.LittleEndian.Uint64(indicationHeader.GetIndicationHeaderFormat1().Cgi.GetNrCgi().NRcellIdentity.Value.GetValue()))

}
