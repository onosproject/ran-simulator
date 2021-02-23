// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package rc

import (
	"github.com/onosproject/ran-simulator/api/types"
)

type RcStore struct {
	RcCells map[types.ECGI]RcCell `mapstructure:"rcCells"`
}

type RcCell struct {
	CellSize uint32     `mapstructure:"cellSize"`
	Earfcn   uint32     `mapstructure:"earfcn"`
	Pci      uint32     `mapstructure:"pci"`
	PciPool  []PciRange `mapstructure:"pciPool"`
}

type PciRange struct {
	min uint32 `mapstructure:"min"`
	max uint32 `mapstructure:"max"`
}
