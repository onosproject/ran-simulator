// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package messageformat2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationMessage(t *testing.T) {

	// TODO
	indicationMessage := NewIndicationMessage()
	assert.NotNil(t, indicationMessage)
}
