// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/mitchellh/go-homedir"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const configDir = ".onos"

var log = logging.GetLogger("manager", "load")

// Load loads the model configuration
func Load(model *Model) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// Set the file name of the configurations file
	viper.SetConfigName("model")

	// Set the path to look for the configurations file
	viper.AddConfigPath("./" + configDir + "/config")
	viper.AddConfigPath(home + "/" + configDir + "/config")
	viper.AddConfigPath("/etc/onos/config")
	viper.AddConfigPath(".")

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read config: %v", err)
		return nil
	}

	// FIXME: coerce Viper to properly unmarshal top level entities, e.g. layout.
	err = viper.Unmarshal(model)
	if err != nil {
		return err
	}

	// Fallback method because Viper is unmarshalling top-level as a map...
	if model.MapLayout.LocationsScale == 0.0 {
		log.Infof("Viper is being stupid! %v", viper.AllSettings())
		bytes, err := ioutil.ReadFile("/etc/onos/config/model.yaml")
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(bytes, model)
		if err != nil {
			return err
		}
	}

	return nil
}
