// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ransim

import (
	"github.com/onosproject/helmit/pkg/input"
	"github.com/onosproject/helmit/pkg/test"
	"github.com/onosproject/ran-simulator/tests/utils"
)

// TestSuite is the primary ran-simulator test suite
type TestSuite struct {
	test.Suite
}

// SetupTestSuite sets up the ran-simulator test suite
func (s *TestSuite) SetupTestSuite(c *input.Context) error {
	sdran, err := utils.CreateSdranRelease(c)
	if err != nil {
		return err
	}
	err = sdran.Install(true)
	if err != nil {
		return err
	}

	// Create an instance of the simulator
	simulator := utils.CreateRanSimulatorWithName(c, "ran-simulator")
	return simulator.Install(true)
}
