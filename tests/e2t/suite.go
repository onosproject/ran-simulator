// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2t

import (
	"github.com/onosproject/helmit/pkg/test"
)

// TestSuite is the primary onos-e2t test suite
type TestSuite struct {
	test.Suite
}

// SetupTestSuite sets up the onos-e2t test suite
func (s *TestSuite) SetupTestSuite() error {
	sdran, err := CreateSdranRelease()
	if err != nil {
		return err
	}
	return sdran.Install(true)
}
