// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package header

import (
	"testing"

	"github.com/onosproject/ran-simulator/pkg/utils"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationHeader(t *testing.T) {
	plmnID := ransimtypes.NewUint24(12345)
	indicationHeader, err := NewIndicationHeader(WithPlmnID(plmnID.Value()),
		WithNRcellIdentity(32)).Build()
	assert.NoError(t, err)

	assert.Equal(t, plmnID.ToBytes(), indicationHeader.GetIndicationHeaderFormat1().Cgi.GetNrCgi().PLmnIdentity.Value)
	assert.Equal(t, uint64(32), utils.BitStringToUint64(indicationHeader.GetIndicationHeaderFormat1().Cgi.GetNrCgi().NRcellIdentity.Value.GetValue(), 36))

}
