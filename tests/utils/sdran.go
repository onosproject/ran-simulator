// SPDX-FileCopyrightText: ${year}-present Open Networking Foundation <info@opennetworking.org>
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/onosproject/helmit/pkg/util/random"
	"github.com/onosproject/onos-test/pkg/onostest"
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
		Set("onos-e2sub.image.tag", "latest")

	return sdran, nil
}

// CreateRanSimulator creates a ran simulator
func CreateRanSimulator() *helm.HelmRelease {
	return CreateRanSimulatorWithName(random.NewPetName(2))
}

// CreateRanSimulatorWithName creates a ran simulator
func CreateRanSimulatorWithName(name string) *helm.HelmRelease {
	username, password, _ := getCredentials()

	simulator := helm.
		Chart(name, onostest.SdranChartRepo).
		Release(name).
		SetUsername(username).
		SetPassword(password).
		Set("image.tag", "latest")
	return simulator
}
