// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package messageformat2

import (
	"encoding/hex"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/pdubuilder"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationMessage(t *testing.T) {

	indicationMessage := NewIndicationMessage(WithUeID("12345"), WithRrcStatus(pdubuilder.CreateRrcStatusConnected()))
	assert.NotNil(t, indicationMessage)
	assert.Equal(t, indicationMessage.ueID, "12345")
	assert.Equal(t, indicationMessage.RrcStatus.Number(), pdubuilder.CreateRrcStatusConnected().Number())

	aper, err := indicationMessage.ToAsn1Bytes()
	assert.NoError(t, err)
	t.Logf("E2SM-MHO-IndicationMessage (Format 2) APER bytes are\n%v", hex.Dump(aper))
}
