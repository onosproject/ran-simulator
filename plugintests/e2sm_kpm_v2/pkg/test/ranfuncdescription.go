// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package test

import (
	"encoding/hex"
	"fmt"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/kpm2/payloads"
	"github.com/onosproject/ran-simulator/plugintests/e2sm_kpm_v2/pkg/modelplugins"
)

func newRanFunctionDescription() error {
	modelPluginRegistry := modelplugins.NewMockModelRegistry()
	shortName, version, err := modelPluginRegistry.RegisterModelPlugin("kpmv2")
	if err != nil {
		return err
	}
	fmt.Printf("Mock model plugin %s %s\n", shortName, version)
	perBytes, err := payloads.RanFunctionDescriptionBytes(ransimtypes.PlmnID(123), modelPluginRegistry)
	if err != nil {
		return err
	}
	fmt.Printf("PER Encode bytes %v", hex.Dump(perBytes))

	return nil
}