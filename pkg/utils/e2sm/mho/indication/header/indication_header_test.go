// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package header

import (
	"encoding/hex"
	"testing"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationHeader(t *testing.T) {
	plmnID := ransimtypes.NewUint24(12345)
	indicationHeader := NewIndicationHeader(WithPlmnID(plmnID.Value()), WithNrcellIdentity([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xE0}))
	assert.NotNil(t, indicationHeader)
	assert.Equal(t, indicationHeader.plmnID.String(), plmnID.String())
	assert.Equal(t, indicationHeader.nrCellIdentity, []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xE0})

	aper, err := indicationHeader.MhoToAsn1Bytes()
	assert.NoError(t, err)
	t.Logf("E2SM-MHO-IndicationHeader APER bytes are\n%v", hex.Dump(aper))
}
