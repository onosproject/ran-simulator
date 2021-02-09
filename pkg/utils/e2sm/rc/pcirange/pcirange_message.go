// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package pcirange

import (
	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

// PciRange lower and upper values of Pci
type PciRange struct {
	lowerPci int32
	upperPci int32
}

// WithLowerPci sets lower pci
func WithLowerPci(lowePci int32) func(pciRange *PciRange) {
	return func(pciRange *PciRange) {
		pciRange.lowerPci = lowePci
	}
}

// WithUpperPci sets upper pci
func WithUpperPci(upperPci int32) func(pciRange *PciRange) {
	return func(pciRange *PciRange) {
		pciRange.upperPci = upperPci
	}
}

// NewPciRange creates a new PciRang message
func NewPciRange(options ...func(pciRange *PciRange)) *PciRange {
	pciRange := &PciRange{}
	for _, option := range options {
		option(pciRange)
	}
	return pciRange
}

// Build builds pciRange IE
func (pciRange *PciRange) Build() (*e2smrcpreies.PciRange, error) {
	picRange := &e2smrcpreies.PciRange{
		LowerPci: &e2smrcpreies.Pci{
			Value: pciRange.lowerPci,
		},
		UpperPci: &e2smrcpreies.Pci{
			Value: pciRange.upperPci,
		},
	}

	return picRange, nil
}
