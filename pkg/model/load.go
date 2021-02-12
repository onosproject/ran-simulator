// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/spf13/viper"
)

const configDir = ".onos"

var log = logging.GetLogger("manager", "load")

// Set up viper for unmarshalling configuration file
func ViperConfigure(configname string) {
	// Set the file type of the configurations file
	viper.SetConfigType("yaml")

	// Set the file name of the configurations file
	viper.SetConfigName(configname)

	// Set the path to look for the configurations file
	viper.AddConfigPath("./" + configDir + "/config")
	viper.AddConfigPath("$HOME/" + configDir + "/config")
	viper.AddConfigPath("/etc/onos/config")
	viper.AddConfigPath(".")
}

// Load model with data in configuration yaml file
func LoadConfig(model *Model, configname string) error {
	var err error

	ViperConfigure(configname)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read %s config: %v", configname, err)
		return err
	}

	err = viper.Unmarshal(model)

	return err;
}

// Load the model configuration.
func Load(model *Model) error {
	return LoadConfig(model, "model")
}
