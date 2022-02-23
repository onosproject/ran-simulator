// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"bytes"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/spf13/viper"
)

const configDir = ".onos"

var log = logging.GetLogger()

// ViperConfigure Sets up viper for unmarshalling configuration file
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

// LoadConfig Loads model with data in configuration yaml file
func LoadConfig(model *Model, configname string) error {
	var err error

	ViperConfigure(configname)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read %s config: %v", configname, err)
		return err
	}

	err = viper.Unmarshal(model)

	// Convert the MCC-MNC format into numeric PLMNID
	model.PlmnID = types.PlmnIDFromString(model.Plmn)

	// initialize neighbor's Ocn value - for mlb/handover
	for k, v := range model.Cells {
		v.MeasurementParams.NCellIndividualOffsets = make(map[types.NCGI]int32)
		for _, n := range v.Neighbors {
			v.MeasurementParams.NCellIndividualOffsets[n] = 0
		}
		model.Cells[k] = v
	}

	return err
}

// Load the model configuration.
func Load(model *Model, modelName string) error {
	return LoadConfig(model, modelName)
}

// LoadConfigFromBytes Loads model with data in configuration yaml file
func LoadConfigFromBytes(model *Model, modelData []byte) error {
	var err error
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewBuffer(modelData)); err != nil {
		log.Errorf("Unable to read %s config: %v", modelData, err)
		return err
	}

	err = viper.Unmarshal(model)

	// Convert the MCC-MNC format into numeric PLMNID
	model.PlmnID = types.PlmnIDFromString(model.Plmn)

	// initialize neighbor's Ocn value - for mlb/handover
	for k, v := range model.Cells {
		v.MeasurementParams.NCellIndividualOffsets = make(map[types.NCGI]int32)
		for _, n := range v.Neighbors {
			v.MeasurementParams.NCellIndividualOffsets[n] = 0
		}
		model.Cells[k] = v
	}
	log.Infof("routeEndPoints: %v", model.RouteEndPoints)
	return err
}
