// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package pciload

import (
	"context"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"
	"github.com/spf13/viper"
)

var log = logging.GetLogger("pci", "load")

// PciMetrics is an auxiliary structure for importing PCI data from YAML configuration
type PciMetrics struct {
	Cells map[types.ECGI]PciCell `mapstructure:"cells"`
}

// PciCell is an auxiliary structure for inport PCI data from YAML configuration
type PciCell struct {
	CellSize string     `mapstructure:"cellSize"`
	Earfcn   uint32     `mapstructure:"earfcn"`
	Pci      uint32     `mapstructure:"pci"`
	PciPool  []PciRange `mapstructure:"pciPool"`
}

// PciRange is an auxiliary structure for inport PCI data from YAML configuration
type PciRange struct {
	Min uint32 `mapstructure:"min"`
	Max uint32 `mapstructure:"max"`
}

// LoadPCIMetrics Loads model with data in "metrics" yaml file
func LoadPCIMetrics(store metrics.Store) error {
	log.Infof("Loading initial PCI metrics...")
	var err error

	model.ViperConfigure("metrics")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read metrics config: %v", err)
		return err
	}

	const key string = "rc.pci"
	if !viper.IsSet(key) {
		log.Infof("PCI metrics not found. skipping...")
	}

	pcis := &PciMetrics{}
	err = viper.UnmarshalKey(key, pcis)
	if err != nil {
		return err
	}

	log.Infof("Storing initial PCI metrics for %d cells...", len(pcis.Cells))

	ctx := context.Background()
	for ecgi, m := range pcis.Cells {
		id := uint64(ecgi)
		_ = store.Set(ctx, id, "cellSize", m.CellSize)
		_ = store.Set(ctx, id, "earfcn", m.Earfcn)
		_ = store.Set(ctx, id, "pci", m.Pci)
		_ = store.Set(ctx, id, "pcipool", m.PciPool)
	}

	return err
}
