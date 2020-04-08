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

/*
Package trafficsim is the main entry point to the ONOS TrafficSim application.

Arguments

-caPath <the location of a CA certificate>

-keyPath <the location of a client private key>

-certPath <the location of a client certificate>


See ../../docs/run.md for how to run the application.
*/
package main

import (
	"flag"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/northbound/e2"
	"github.com/onosproject/ran-simulator/pkg/northbound/trafficsim"
	"github.com/onosproject/ran-simulator/pkg/southbound/kubernetes"
	"github.com/onosproject/ran-simulator/pkg/utils"
	_ "net/http/pprof"
	"runtime"
	"time"
)

var log = liblog.GetDefaultLogger()

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	googleAPIKey := flag.String("googleAPIKey", "", "your google maps api key")
	towerRows := flag.Int("towerRows", 3, "Number of rows of towers")
	towerCols := flag.Int("towerCols", 3, "Number of columns of towers")
	mapCenterLat := flag.Float64("mapCenterLat", 52.5200, "Map center latitude")
	mapCenterLng := flag.Float64("mapCenterLng", 13.4050, "Map center longitude") // Berlin
	zoom := flag.Float64("zoom", 13, "The starting Zoom level")
	fade := flag.Bool("fade", true, "Show map as faded on start")
	showRoutes := flag.Bool("showRoutes", true, "Show routes on start")
	showPower := flag.Bool("showPower", true, "Show power as circle on start")
	towerSpacingVert := flag.Float64("towerSpacingVert", 0.02, "Tower spacing vert in degrees latitude")
	towerSpacingHoriz := flag.Float64("towerSpacingHoriz", 0.0333, "Tower spacing horiz in degrees longitude")
	locationsScale := flag.Float64("locationsScale", 1.25, "Ratio of random locations diameter to tower grid width")
	maxUEs := flag.Int("maxUEs", 300, "Max number of UEs for complete simulation")
	minUEs := flag.Int("minUEs", 3, "Max number of UEs for complete simulation")
	stepDelayMs := flag.Int("stepDelayMs", 1000, "delay between steps on route")
	maxUEsPerTower := flag.Int("maxUEsPerTower", 5, "Max num of UEs per tower")
	metricsPort := flag.Int("metricsPort", 9090, "port for Prometheus metrics")
	metricsAllHoEvents := flag.Bool("metricsAllHoEvents", true, "Export all HO events in metrics (only historgram if false)")
	topoEndpoint := flag.String("topoEndpoint", "onos-topo:5150", "Endpoint for the onos-topo service")
	loglevel := flag.String("loglevel", "warn", "Initial log level - debug, info, warn, error")
	addK8sSvcPorts := flag.Bool("addK8sSvcPorts", true, "Add K8S service ports per tower")

	flag.Parse()
	setLogLevel(*loglevel)

	mapLayoutParams := types.MapLayout{
		Center: &types.Point{
			Lat: float32(*mapCenterLat),
			Lng: float32(*mapCenterLng),
		},
		Zoom:       float32(*zoom),
		ShowRoutes: *showRoutes,
		Fade:       *fade,
		ShowPower:  *showPower,
		MinUes:     uint32(*minUEs),
		MaxUes:     uint32(*maxUEs),
	}
	if mapLayoutParams.Zoom < 10 || mapLayoutParams.Zoom > 15 {
		log.Fatal("Invalid Zoom level - must be between 10 and 15 inclusive")
	}
	if mapLayoutParams.Center.GetLat() <= -90.0 || mapLayoutParams.Center.GetLat() >= 90.0 {
		log.Fatal("Invalid Map Centre Latitude - must be between -90 and 90 exclusive")
	}
	if mapLayoutParams.Center.GetLng() <= -180.0 || mapLayoutParams.Center.GetLng() >= 180.0 {
		log.Fatal("Invalid Map Centre Longitude - must be between -180 and 180 exclusive")
	}

	if mapLayoutParams.MaxUes < 10 || mapLayoutParams.MaxUes > 1000000 {
		log.Fatal("Invalid number for MaxUEs - must be between 10 and 1000000")
	}
	if mapLayoutParams.MinUes < 2 || mapLayoutParams.MinUes > 1000 {
		log.Fatal("Invalid number for MinUEs - must be between 2 and 1000 inclusive")
	}
	if mapLayoutParams.MinUes*2 > mapLayoutParams.MaxUes {
		log.Fatal("Invalid ratio of MaxUEs:MinUEs - must be at least 2")
	}

	towerParams := types.TowersParams{
		TowerRows:         uint32(*towerRows),
		TowerCols:         uint32(*towerCols),
		TowerSpacingVert:  float32(*towerSpacingVert),
		TowerSpacingHoriz: float32(*towerSpacingHoriz),
		MaxUEsPerTower:    uint32(*maxUEsPerTower),
		LocationsScale:    float32(*locationsScale),
	}
	checkTowerLimits(*towerRows, *towerCols)
	if towerParams.TowerSpacingVert < 0.001 || towerParams.TowerSpacingVert > 1.0 {
		log.Fatal("Invalid vertical tower spacing - must be between 0.001 and 1.0 degree latitude inclusive")
	}
	if towerParams.TowerSpacingHoriz < 0.001 || towerParams.TowerSpacingHoriz > 1.0 {
		log.Fatal("Invalid horizontal tower spacing - must be between 0.001 and 1.0 degree longitude inclusive")
	}

	if towerParams.LocationsScale < 0.1 || towerParams.LocationsScale > 2.0 {
		log.Fatal("Invalid locationsScale - must be between 0.1 and 2.0")
	}

	routesParams := manager.RoutesParams{
		APIKey:    *googleAPIKey,
		StepDelay: time.Duration(*stepDelayMs) * time.Millisecond,
	}

	if *stepDelayMs < 100 || *stepDelayMs > 60000 {
		log.Fatal("Invalid step Delay - must be between 100ms and 60000ms inclusive")
	}

	serverParams := utils.ServerParams{
		CaPath:       *caPath,
		KeyPath:      *keyPath,
		CertPath:     *certPath,
		TopoEndpoint: *topoEndpoint,
	}

	for r := 0; r < *towerRows; r++ {
		for c := 0; c < *towerCols; c++ {
			towerNum := r**towerCols + c + 1 // Start at 1
			go func() {
				// Blocks here when server running
				err := e2.NewTowerServer(towerNum, utils.TestPlmnID, serverParams)
				if err != nil {
					log.Fatal("Unable to start server ", err)
				}
			}()
		}
	}
	// Add these new ports to the K8s service
	rangeStart := utils.GrpcBasePort + 2
	rangeEnd := rangeStart + *towerCols**towerRows
	if *addK8sSvcPorts {
		err := kubernetes.AddK8SServicePorts(int32(rangeStart), int32(rangeEnd))
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	metricsParams := manager.MetricsParams{
		Port:              *metricsPort,
		ExportAllHOEvents: *metricsAllHoEvents,
	}

	log.Info("Starting trafficsim")
	mgr, err := manager.NewManager()
	if err != nil {
		log.Fatal("Unable to load trafficsim ", err)
		return
	}
	mgr.Run(mapLayoutParams, towerParams, routesParams, *topoEndpoint, *metricsPort, serverParams, metricsParams)

	if err = startServer(*caPath, *keyPath, *certPath); err != nil {
		log.Fatal("Unable to start trafficsim ", err)
	}
	mgr.Close()
}

// Creates gRPC server and registers various services; then serves.
func startServer(caPath string, keyPath string, certPath string) error {
	s := service.NewServer(service.NewServerConfig(caPath, keyPath, certPath, 5150, true))
	s.AddService(trafficsim.Service{})

	return s.Serve(func(started string) {
		log.Info("Started NBI on ", started)
	})
}

func checkTowerLimits(rows int, cols int) {
	if rows < 2 || rows > 64 {
		log.Fatal("Invalid number of Tower Rows - must be between 2 and 64 inclusive")
	}
	if cols < 2 || cols > 64 {
		log.Fatal("Invalid number of Tower Cols - must be between 2 and 64 inclusive")
	}
	if cols*rows > 1024 {
		log.Fatal("Invalid number of Tower (Rows x Cols) - must not exceed 1024")
	}
}

func setLogLevel(loglevel string) {
	initialLogLevel := liblog.WarnLevel
	switch loglevel {
	case "debug":
		initialLogLevel = liblog.DebugLevel
	case "info":
		initialLogLevel = liblog.InfoLevel
	case "warn":
		initialLogLevel = liblog.WarnLevel
	case "error":
		initialLogLevel = liblog.ErrorLevel
	}

	log.Infof("logs level: %s", initialLogLevel)
	runtime.SetMutexProfileFraction(5)
	log.SetLevel(initialLogLevel)
	liblog.GetLogger("northbound").SetLevel(initialLogLevel)
	liblog.GetLogger("northbound", "e2").SetLevel(initialLogLevel)
	liblog.GetLogger("northbound", "trafficsim").SetLevel(initialLogLevel)
	liblog.GetLogger("manager").SetLevel(initialLogLevel)
	liblog.GetLogger("dispatcher").SetLevel(initialLogLevel)
	liblog.GetLogger("southbound", "kubernetes").SetLevel(initialLogLevel)
	liblog.GetLogger("southbound", "topo").SetLevel(initialLogLevel)
}
