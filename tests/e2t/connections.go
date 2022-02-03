// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package e2t

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onosproject/ran-simulator/tests/utils"
)

// TestConnections test connectivity between e2 agents and e2t
func (s *TestSuite) TestConnections(t *testing.T) {
	connections, err := utils.GetE2Connections()
	assert.NoError(t, err, "unable to connect to E2T admin service %v", err)
	assert.Equal(t, 2, len(connections), "incorrect connection count")
}
