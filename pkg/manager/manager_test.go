// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package manager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("test_cfg.yml")
	assert.NoError(t, err, "unable to load config")
	if err == nil {
		assert.Equal(t, "/foo/ca", cfg.CAPath, "incorrect caCertPath")
		assert.Equal(t, "/foo/cert", cfg.CertPath, "incorrect certPath")
		assert.Equal(t, "/foo/key", cfg.KeyPath, "incorrect keyPath")
		assert.Equal(t, 901, cfg.GRPCPort, "incorrect grpc port")
	}
}

func TestNewManager(t *testing.T) {
	cfg, err := LoadConfig("test_cfg.yml")
	assert.NoError(t, err, "unable to load config")

	_, err = NewManager(cfg)
	assert.NoError(t, err, "unable to create new manager")
}
