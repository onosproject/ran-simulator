// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package manager

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
)

// Load the startup configuration.
func (m *Manager) LoadStartup() error {
	model.ViperConfigure("startup")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read %s config: %v", "startup", err)
		return err
	}

	for enbID, agent := range m.agents.Agents {
		for name, id := range registry.StringToRanFunctionID {
			key := fmt.Sprintf("%d.servicemodels.%s", enbID, name)
			if viper.IsSet(key) {
				sm, _ := agent.GetSM(id)
				plugin := sm.ModelPluginRegistry.ModelPlugins[sm.ModelFullName]
				viper.UnmarshalKey(key, plugin)
				var _ = plugin
			}
		}
	}
	return nil
}
