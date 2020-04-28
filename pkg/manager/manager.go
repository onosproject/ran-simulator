// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package manager is is the main coordinator for the ONOS RAN subsystem.
package manager

import (
	"context"
	"sync"

	"github.com/onosproject/onos-topo/api/device"
	"github.com/onosproject/ran-simulator/pkg/southbound/topo"
	"github.com/onosproject/ran-simulator/pkg/utils"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/config"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"github.com/onosproject/ran-simulator/pkg/northbound/metrics"
)

var log = logging.GetLogger("manager")

var mgr Manager

// Manager single point of entry for the trafficsim system.
type Manager struct {
	MapLayout             types.MapLayout
	Cells                 map[types.ECGI]*types.Cell
	CellsLock             *sync.RWMutex
	Locations             map[string]*Location
	Routes                map[types.Imsi]*types.Route
	UserEquipments        map[types.Imsi]*types.Ue
	UserEquipmentsLock    *sync.RWMutex
	UserEquipmentsMapLock *sync.RWMutex
	Dispatcher            *dispatcher.Dispatcher
	UeChannel             chan dispatcher.Event
	RouteChannel          chan dispatcher.Event
	CellsChannel          chan dispatcher.Event
	googleAPIKey          string
	LatencyChannel        chan metrics.HOEvent
	ResetMetricsChannel   chan bool
	TopoClient            device.DeviceServiceClient
	AspectRatio           float64
}

// MetricsParams for the Prometheus exporter
type MetricsParams struct {
	Port              uint
	ExportAllHOEvents bool
}

// NewManager initializes the RAN subsystem.
func NewManager() (*Manager, error) {
	log.Info("Creating Manager")
	mgr = Manager{
		CellsLock:             &sync.RWMutex{},
		UserEquipmentsLock:    &sync.RWMutex{},
		UserEquipmentsMapLock: &sync.RWMutex{},
		Dispatcher:            dispatcher.NewDispatcher(),
		UeChannel:             make(chan dispatcher.Event),
		RouteChannel:          make(chan dispatcher.Event),
		CellsChannel:          make(chan dispatcher.Event),
		LatencyChannel:        make(chan metrics.HOEvent),
		ResetMetricsChannel:   make(chan bool),
	}
	return &mgr, nil
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run(mapLayoutParams types.MapLayout, towerConfig config.TowerConfig,
	routesParams RoutesParams, topoEndpoint string, serverParams utils.ServerParams,
	metricsParams MetricsParams) {
	log.Infof("Starting Manager with %v %v", mapLayoutParams, routesParams)

	m.MapLayout = mapLayoutParams
	m.CellsLock.Lock()
	m.Cells = NewCells(towerConfig)
	m.CellsLock.Unlock()
	m.Locations = NewLocations(towerConfig, int(mapLayoutParams.MaxUes), mapLayoutParams.LocationsScale)
	m.MapLayout.MinUes = mapLayoutParams.MinUes
	m.MapLayout.MaxUes = mapLayoutParams.MaxUes
	m.googleAPIKey = routesParams.APIKey
	// Compensate for the narrowing of meridians at higher latitudes
	m.AspectRatio = utils.AspectRatio(&towerConfig.MapCentre)

	go m.Dispatcher.ListenUeEvents(m.UeChannel)
	go m.Dispatcher.ListenRouteEvents(m.RouteChannel)
	go m.Dispatcher.ListenCellEvents(m.CellsChannel)

	var err error
	m.Routes, err = m.NewRoutes(mapLayoutParams, routesParams)
	if err != nil {
		log.Fatalf("Error calculating routes %s", err.Error())
	}
	m.UserEquipments, err = m.NewUserEquipments(mapLayoutParams, routesParams)
	if err != nil {
		log.Fatalf("Error creating new UEs %s", err.Error())
	}
	go m.startMoving(routesParams)

	go metrics.RunHOExposer(int(metricsParams.Port), m.LatencyChannel, metricsParams.ExportAllHOEvents, m.ResetMetricsChannel)

	ctx := context.Background()
	m.TopoClient = topo.ConnectToTopo(ctx, topoEndpoint, serverParams)
	go topo.SyncToTopo(ctx, &m.TopoClient, m.Cells)
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	close(m.CellsChannel)
	close(m.UeChannel)
	close(m.RouteChannel)
	close(m.LatencyChannel)
	for r := range m.Routes {
		delete(m.Routes, r)
	}
	for l := range m.Locations {
		delete(m.Locations, l)
	}
	m.CellsLock.Lock()
	for tid := range m.Cells {
		delete(m.Cells, tid)
	}
	m.CellsLock.Unlock()
	// TODO - clean up the topo entries on shutdown
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}
