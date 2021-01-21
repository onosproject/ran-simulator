// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2t

import (
	"testing"

	"github.com/onosproject/ran-simulator/tests/utils"
)

// TestConnections
func (s *TestSuite) TestConnections(t *testing.T) {
	utils.CreateRanSimulatorWithName(t, "ran-simulator")

}
