// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package header

import (
	"testing"

	"github.com/onosproject/ran-simulator/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationHeader(t *testing.T) {
	indicationHeader, err := NewIndicationHeader(WithPlmnID(12345),
		WithEutracellIdentity(32)).Build()
	assert.NoError(t, err)

	assert.Equal(t, indicationHeader.GetIndicationHeaderFormat1().Cgi.GetEUtraCgi().PLmnIdentity.Value, utils.PlmnIDToByteArray(12345))
	assert.Equal(t, indicationHeader.GetIndicationHeaderFormat1().Cgi.GetEUtraCgi().EUtracellIdentity.Value.Value, uint64(32))

}
