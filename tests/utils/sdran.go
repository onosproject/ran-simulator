// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"

	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/input"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/onosproject/helmit/pkg/util/random"
	"github.com/onosproject/onos-test/pkg/onostest"
)

func getCredentials() (string, string, error) {
	kubClient, err := kubernetes.New()
	if err != nil {
		return "", "", err
	}
	secrets, err := kubClient.CoreV1().Secrets().Get(context.Background(), onostest.SecretsName)
	if err != nil {
		return "", "", err
	}
	username := string(secrets.Object.Data["sd-ran-username"])
	password := string(secrets.Object.Data["sd-ran-password"])

	return username, password, nil
}

// CreateSdranRelease creates a helm release for an sd-ran instance
func CreateSdranRelease(c *input.Context) (*helm.HelmRelease, error) {
	username, password, err := getCredentials()
	registry := c.GetArg("registry").String("")
	if err != nil {
		return nil, err
	}

	sdran := helm.Chart("sd-ran", onostest.SdranChartRepo).
		Release("sd-ran").
		SetUsername(username).
		SetPassword(password).
		Set("global.image.registry", registry).
		Set("import.onos-config.enabled", false)

	return sdran, nil
}

// CreateRanSimulator creates a ran simulator
func CreateRanSimulator(c *input.Context) *helm.HelmRelease {
	return CreateRanSimulatorWithName(c, random.NewPetName(2))
}

// CreateRanSimulatorWithName creates a ran simulator
func CreateRanSimulatorWithName(c *input.Context, name string) *helm.HelmRelease {
	username, password, _ := getCredentials()
	registry := c.GetArg("registry").String("")

	simulator := helm.
		Chart(name, onostest.SdranChartRepo).
		Release(name).
		SetUsername(username).
		SetPassword(password).
		Set("global.image.registry", registry).
		Set("image.tag", "latest")
	return simulator
}
