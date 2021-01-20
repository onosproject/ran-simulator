// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2t

import (
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConnection tests the basic E2 connection of the simulator to the E2T node(s)
func (s *TestSuite) TestConnection(t *testing.T) {
	cfg, err := manager.LoadConfig("test_cfg.yml")
	assert.NoError(t, err, "unable to load config")

	mgr, err := manager.NewManager(cfg)
	assert.NoError(t, err, "unable to create manager")

	err = mgr.Start()
	assert.NoError(t, err, "unable to start")
}
