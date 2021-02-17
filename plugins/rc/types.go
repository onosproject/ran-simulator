// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package main

import (
	"github.com/onosproject/ran-simulator/pkg/model"
)

type serviceModel struct {
	Pci     map[model.ECGI]uint32 `mapstructure:"pci"`
	PciPool []PciRange            `mapstructure:"pciPool"`
}

type PciRange struct {
	min uint32 `mapstructure:"min"`
	max uint32 `mapstructure:"max"`
}
