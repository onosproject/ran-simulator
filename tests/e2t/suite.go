// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2t

import (
	"github.com/onosproject/helmit/pkg/test"
	"github.com/onosproject/ran-simulator/tests/utils"
)

// TestSuite is the primary ran-simulator test suite
type TestSuite struct {
	test.Suite
}

// SetupTestSuite sets up the ran-simulator test suite
func (s *TestSuite) SetupTestSuite() error {
	sdran, err := utils.CreateSdranRelease()
	if err != nil {
		return err
	}
	return sdran.Install(true)
}
