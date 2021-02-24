// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package pciload

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"
	"github.com/spf13/viper"
)

var log = logging.GetLogger("pci", "load")

// PciMetrics is an auxiliary structure for importing PCI data from YAML configuration
type PciMetrics struct {
	Cells map[uint64]PciCell `mapstructure:"cells"`
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

// LoadPCIMetrics Loads model with data in "pcis" yaml file
func LoadPCIMetrics(store metrics.Store) error {
	log.Infof("Loading initial PCI metrics...")
	var err error

	model.ViperConfigure("pci")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read metrics config: %v", err)
		return err
	}

	pcis := &PciMetrics{}
	err = viper.Unmarshal(pcis)
	if err != nil {
		return err
	}

	log.Infof("Storing initial PCI metrics for %d cells...", len(pcis.Cells))

	ctx := context.Background()
	for id, m := range pcis.Cells {
		_ = store.Set(ctx, id, "cellSize", m.CellSize)
		_ = store.Set(ctx, id, "earfcn", m.Earfcn)
		_ = store.Set(ctx, id, "pci", m.Pci)

		for i, p := range m.PciPool {
			_ = store.Set(ctx, id, fmt.Sprintf("pci%dMin", i+1), p.Min)
			_ = store.Set(ctx, id, fmt.Sprintf("pci%dMax", i+1), p.Max)
		}
	}

	return err
}
