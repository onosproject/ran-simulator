// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package pciload

import (
	"bytes"
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
	Cells map[types.ECGI]PciCell `mapstructure:"cells" yaml:"cells"`
}

// PciCell is an auxiliary structure for inport PCI data from YAML configuration
type PciCell struct {
	CellSize string     `mapstructure:"cellSize" yaml:"cellSize"`
	Earfcn   uint32     `mapstructure:"earfcn" yaml:"earfcn"`
	Pci      uint32     `mapstructure:"pci" yaml:"pci"`
	PciPool  []PciRange `mapstructure:"pciPool" yaml:"pciPool"`
}

// PciRange is an auxiliary structure for inport PCI data from YAML configuration
type PciRange struct {
	Min uint32 `mapstructure:"min" yaml:"min"`
	Max uint32 `mapstructure:"max" yaml:"max"`
}

// LoadPCIMetrics loads metrics with data in "metrics" yaml file
func LoadPCIMetrics(store metrics.Store, metricName string) error {
	return LoadPCIMetricsConfig(store, metricName)
}

// LoadPCIMetricsConfig loads metrics with data in the named configuration
func LoadPCIMetricsConfig(store metrics.Store, configName string) error {
	log.Infof("Loading PCI metrics from %s...", configName)

	model.ViperConfigure(configName)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read metrics config: %v", err)
		return err
	}

	return unmarshal(store)
}

// LoadPCIMetricsData loads metrics with data in the specified bytes
func LoadPCIMetricsData(store metrics.Store, metricsData []byte) error {
	log.Infof("Loading PCI metrics from bytes...")

	// Set the file type of the configurations file
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewBuffer(metricsData)); err != nil {
		log.Errorf("Unable to read metrics config: %v", err)
		return err
	}

	return unmarshal(store)
}

func unmarshal(store metrics.Store) error {
	const key string = "rc.pci"
	if !viper.IsSet(key) {
		log.Infof("PCI metrics not found. skipping...")
	}

	pcis := &PciMetrics{}
	if err := viper.UnmarshalKey(key, pcis); err != nil {
		return err
	}

	log.Infof("Storing PCI metrics for %d cells...", len(pcis.Cells))

	ctx := context.Background()
	for ecgi, m := range pcis.Cells {
		id := uint64(ecgi)
		_ = store.Set(ctx, id, "cellSize", m.CellSize)
		_ = store.Set(ctx, id, "earfcn", m.Earfcn)
		_ = store.Set(ctx, id, "pci", m.Pci)
		_ = store.Set(ctx, id, "pcipool", m.PciPool)
	}

	return nil
}
