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
	"sync"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"github.com/onosproject/ran-simulator/pkg/northbound/metrics"
)

var log = logging.GetLogger("manager")

var mgr Manager

// Manager single point of entry for the trafficsim system.
type Manager struct {
	MapLayout          types.MapLayout
	Towers             map[string]*types.Tower
	TowersLock         *sync.RWMutex
	Locations          map[string]*Location
	Routes             map[string]*types.Route
	RoutesLock         *sync.RWMutex
	UserEquipments     map[string]*types.Ue
	UserEquipmentsLock *sync.RWMutex
	Dispatcher         *dispatcher.Dispatcher
	UeChannel          chan dispatcher.Event
	RouteChannel       chan dispatcher.Event
	TowerChannel       chan dispatcher.Event
	googleAPIKey       string
	LatencyChannel     chan metrics.HOEvent
}

// NewManager initializes the RAN subsystem.
func NewManager() (*Manager, error) {
	log.Info("Creating Manager")
	mgr = Manager{
		TowersLock:         &sync.RWMutex{},
		RoutesLock:         &sync.RWMutex{},
		UserEquipmentsLock: &sync.RWMutex{},
		Dispatcher:         dispatcher.NewDispatcher(),
		UeChannel:          make(chan dispatcher.Event),
		RouteChannel:       make(chan dispatcher.Event),
		TowerChannel:       make(chan dispatcher.Event),
		LatencyChannel:     make(chan metrics.HOEvent),
	}
	return &mgr, nil
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run(mapLayoutParams types.MapLayout, towerparams types.TowersParams,
	routesParams RoutesParams, metricsPort int) {
	log.Infof("Starting Manager with %v %v %v", mapLayoutParams, towerparams, routesParams)
	m.MapLayout = mapLayoutParams
	m.TowersLock.Lock()
	m.Towers = NewTowers(towerparams, mapLayoutParams)
	m.TowersLock.Unlock()
	m.Locations = NewLocations(towerparams, mapLayoutParams)
	m.MapLayout.MinUes = mapLayoutParams.MinUes
	m.MapLayout.MaxUes = mapLayoutParams.MaxUes
	m.googleAPIKey = routesParams.APIKey

	go m.Dispatcher.ListenUeEvents(m.UeChannel)
	go m.Dispatcher.ListenRouteEvents(m.RouteChannel)
	go m.Dispatcher.ListenTowerEvents(m.TowerChannel)

	var err error
	m.Routes, err = m.NewRoutes(mapLayoutParams, routesParams)
	if err != nil {
		log.Fatalf("Error calculating routes %s", err.Error())
	}
	m.UserEquipments = m.NewUserEquipments(mapLayoutParams, routesParams)

	go m.startMoving(routesParams)

	go metrics.RunHOExposer(metricsPort, m.LatencyChannel)
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	close(m.TowerChannel)
	close(m.UeChannel)
	close(m.RouteChannel)
	close(m.LatencyChannel)
	for r := range m.Routes {
		delete(m.Routes, r)
	}
	for l := range m.Locations {
		delete(m.Locations, l)
	}
	m.TowersLock.Lock()
	for tid := range m.Towers {
		delete(m.Towers, tid)
	}
	m.TowersLock.Unlock()
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}
