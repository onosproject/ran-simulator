// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package utils

import (
	"testing"

	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/util/random"
	"github.com/onosproject/onos-test/pkg/onostest"
	"github.com/stretchr/testify/assert"
)

// CreateSdranRelease creates a helm release for an sd-ran instance
func CreateSdranRelease() (*helm.HelmRelease, error) {
	sdran := helm.Chart("sd-ran", onostest.SdranChartRepo).
		Release("sd-ran").
		Set("import.onos-config.enabled", false).
		Set("import.onos-topo.enabled", false).
		Set("onos-e2t.image.tag", "latest").
		Set("onos-e2sub.image.tag", "latest")

	return sdran, nil
}

// CreateE2Simulator creates a device simulator
func CreateE2Simulator(t *testing.T) *helm.HelmRelease {
	return CreateE2SimulatorWithName(t, random.NewPetName(2))
}

// CreateE2SimulatorWithName creates a device simulator
func CreateE2SimulatorWithName(t *testing.T, name string) *helm.HelmRelease {
	simulator := helm.
		Chart("e2-simulator", onostest.SdranChartRepo).
		Release(name).
		Set("image.tag", "latest")
	err := simulator.Install(true)
	assert.NoError(t, err, "could not install device simulator %v", err)

	return simulator
}
