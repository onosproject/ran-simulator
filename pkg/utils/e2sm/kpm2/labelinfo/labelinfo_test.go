// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package labelinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLabelInfo(t *testing.T) {
	labelInfo, err := NewLabelInfo(WithFiveQI(200))
	assert.NoError(t, err)
	assert.Equal(t, int32(200), labelInfo.fiveQI)

}

func TestWrongFiveQI(t *testing.T) {
	_, err := NewLabelInfo(WithFiveQI(400))
	assert.Error(t, err)

}
