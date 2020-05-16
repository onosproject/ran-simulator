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
	"github.com/onosproject/config-models/modelplugin/e2node-1.0.0/e2node_1_0_0"
	"sync"
	"time"

	"github.com/onosproject/onos-topo/api/device"
	"github.com/onosproject/ran-simulator/pkg/southbound/topo"
	"github.com/onosproject/ran-simulator/pkg/utils"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"github.com/onosproject/ran-simulator/pkg/northbound/metrics"
)

var log = logging.GetLogger("manager")

var mgr Manager

// NewServerHandler a call back function to avoid import cycle
type NewServerHandler func(ecgi types.ECGI, port uint16, serverParams utils.ServerParams) error

// Manager single point of entry for the trafficsim system.
type Manager struct {
	MapLayout             types.MapLayout
	Cells                 map[types.ECGI]*types.Cell
	CellConfigs           map[types.ECGI]*e2node_1_0_0.Device
	CellsLock             *sync.RWMutex
	Locations             map[LocationID]*Location
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
	cellCreateTimer       *time.Timer
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
		Cells:                 make(map[types.ECGI]*types.Cell),
		CellConfigs:           make(map[types.ECGI]*e2node_1_0_0.Device),
		Locations:             make(map[LocationID]*Location),
		Routes:                make(map[types.Imsi]*types.Route),
		UserEquipments:        make(map[types.Imsi]*types.Ue),
		cellCreateTimer:       time.NewTimer(time.Second),
	}
	return &mgr, nil
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run(mapLayoutParams types.MapLayout, routesParams RoutesParams,
	topoEndpoint string, serverParams utils.ServerParams,
	metricsParams MetricsParams, newServerHandler NewServerHandler) {
	log.Infof("Starting Manager with %v %v", mapLayoutParams, routesParams)

	m.MapLayout = mapLayoutParams
	m.CellsLock.Lock()
	for ecgi := range m.Cells {
		m.CellConfigs[ecgi] = &e2node_1_0_0.Device{
			E2Node: &e2node_1_0_0.E2Node_E2Node{
				Intervals: &e2node_1_0_0.E2Node_E2Node_Intervals{},
			},
		}
	}
	m.CellsLock.Unlock()
	m.MapLayout.Center = &types.Point{}
	m.MapLayout.MinUes = mapLayoutParams.MinUes
	m.MapLayout.MaxUes = mapLayoutParams.MaxUes
	m.googleAPIKey = routesParams.APIKey
	// Compensate for the narrowing of meridians at higher latitudes
	m.AspectRatio = utils.AspectRatio(m.MapLayout.Center)

	go m.Dispatcher.ListenUeEvents(m.UeChannel)
	go m.Dispatcher.ListenRouteEvents(m.RouteChannel)
	go m.Dispatcher.ListenCellEvents(m.CellsChannel)

	go metrics.RunHOExposer(int(metricsParams.Port), m.LatencyChannel, metricsParams.ExportAllHOEvents, m.ResetMetricsChannel)

	go m.afterCellCreation(mapLayoutParams, routesParams, serverParams, newServerHandler)

	ctx := context.Background()

	go func() {
		var err error
		m.TopoClient, err = topo.ConnectToTopo(ctx, topoEndpoint, serverParams, CellCreator, CellDeleter)
		if err != nil {
			log.Fatalf("Error connecting to onos-topo %v", err)
		}
	}()
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	if err := m.SetNumberUes(0); err != nil {
		log.Warnf("Unable to set number of UEs to 0 %s", err.Error())
	}

	time.Sleep(time.Second) // Wait for topo, but don't hang around if it can't be done
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
	log.Warn("Closed Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}
