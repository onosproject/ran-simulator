// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package test

import (
	"gotest.tools/assert"
	"testing"
)

func Test_newRanFunctionDescription(t *testing.T) {
	err := newRanFunctionDescription()

	assert.NilError(t, err)
}