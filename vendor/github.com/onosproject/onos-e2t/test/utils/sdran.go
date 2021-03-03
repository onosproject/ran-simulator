// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package utils

import (
	"testing"

	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/onosproject/helmit/pkg/util/random"
	"github.com/onosproject/onos-test/pkg/onostest"
	"github.com/stretchr/testify/assert"
)

func getCredentials() (string, string, error) {
	kubClient, err := kubernetes.New()
	if err != nil {
		return "", "", err
	}
	secrets, err := kubClient.CoreV1().Secrets().Get(onostest.SecretsName)
	if err != nil {
		return "", "", err
	}
	username := string(secrets.Object.Data["sd-ran-username"])
	password := string(secrets.Object.Data["sd-ran-password"])

	return username, password, nil
}

// CreateSdranRelease creates a helm release for an sd-ran instance
func CreateSdranRelease() (*helm.HelmRelease, error) {
	username, password, err := getCredentials()
	if err != nil {
		return nil, err
	}

	sdran := helm.Chart("sd-ran", onostest.SdranChartRepo).
		Release("sd-ran").
		SetUsername(username).
		SetPassword(password).
		Set("import.onos-config.enabled", false).
		Set("import.onos-topo.enabled", false).
		Set("onos-e2t.image.tag", "latest").
		Set("onos-e2sub.image.tag", "latest").
		Set("ran-simulator.image.tag", "latest")

	return sdran, nil
}

// CreateRanSimulator creates a ran simulator
func CreateRanSimulator(t *testing.T) *helm.HelmRelease {
	return CreateRanSimulatorWithName(t, random.NewPetName(2))
}

// CreateRanSimulatorWithNameOrDie creates a simulator and fails the test if the creation returned an error
func CreateRanSimulatorWithNameOrDie(t *testing.T, simName string) *helm.HelmRelease {
	sim := CreateRanSimulatorWithName(t, simName)
	assert.NotNil(t, sim)
	return sim
}

// CreateRanSimulatorWithName creates a ran simulator
func CreateRanSimulatorWithName(t *testing.T, name string) *helm.HelmRelease {
	username, password, err := getCredentials()
	assert.NoError(t, err)

	simulator := helm.
		Chart(name, onostest.SdranChartRepo).
		Release(name).
		SetUsername(username).
		SetPassword(password).
		Set("image.tag", "latest")
	err = simulator.Install(true)
	assert.NoError(t, err, "could not install device simulator %v", err)

	return simulator
}
