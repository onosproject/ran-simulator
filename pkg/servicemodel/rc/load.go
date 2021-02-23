// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package rc

import (
	"fmt"

	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/spf13/viper"
)

// Load the startup configuration.
func LoadRC(enbID types.EnbID, store *RcStore) error {
	model.ViperConfigure("startup")

	if err := viper.ReadInConfig()
    err != nil {
		log.Errorf("Unable to read %s config: %v", "startup", err)
		return err
	}

	key := fmt.Sprintf("servicemodels.rc.%d", enbID)
	if viper.IsSet(key) {
		err := viper.UnmarshalKey(key, store)
		if err != nil {
			log.Errorf("Error unmarshaling RC service model for enbID %d: %v", enbID, err)
			return err
		}
	} else {
	    log.Errorf("KEY NOT FOUND %d", 1)
	}

	return nil
}
