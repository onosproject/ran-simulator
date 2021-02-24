// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package pciload

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"
	"github.com/spf13/viper"
)

var log = logging.GetLogger("pci", "load")

// PCIMetrics is an auxiliary structure for importing PCI data from YAML configuration
type PCIMetrics struct {
	Values map[uint64]uint32 `mapstructure:"pcis"`
}

// LoadPCIMetrics Loads model with data in "pcis" yaml file
func LoadPCIMetrics(store metrics.Store) error {
	var err error

	model.ViperConfigure("pcis")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read metrics config: %v", err)
		return err
	}

	pcis := &PCIMetrics{}
	err = viper.Unmarshal(pcis)
	if err != nil {
		return err
	}

	ctx := context.Background()
	for id, pci := range pcis.Values {
		_ = store.Set(ctx, id, "pci", pci)
	}

	return err
}
