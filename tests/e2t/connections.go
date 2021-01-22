// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2t

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onosproject/ran-simulator/tests/utils"
)

// TestConnections test connectivity between e2 agents and e2t
func (s *TestSuite) TestConnections(t *testing.T) {
	// Creates an instance of the simulator
	simulator := utils.CreateRanSimulatorWithName(t, "ran-simulator")
	err := simulator.Install(true)
	assert.NoError(t, err, "could not install device simulator %v", err)

	// TODO retrieve list of connections in onos-e2t
}
