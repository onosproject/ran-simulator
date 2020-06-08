// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package northbound

import (
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/northbound/e2"
	"github.com/onosproject/ran-simulator/pkg/northbound/gnmi"
	"github.com/onosproject/ran-simulator/pkg/utils"
)

var log = liblog.GetLogger("northbound")

// NewCellServer - start a new gRPC server per cell
func NewCellServer(ecgi types.ECGI, port uint16, serverParams utils.ServerParams) error {
	s := service.NewServer(service.NewServerConfig(serverParams.CaPath, serverParams.KeyPath, serverParams.CertPath, int16(port), true))
	s.AddService(e2.Service{
		Port:      int(port),
		TowerEcID: ecgi.EcID,
		PlmnID:    ecgi.PlmnID,
	})
	s.AddService(gnmi.Service{
		Port:      int(port),
		TowerEcID: ecgi.EcID,
		PlmnID:    ecgi.PlmnID,
	})

	return s.Serve(func(started string) {
		log.Info("Started E2 & gNMI server for a cell on ", started)
	})
}
