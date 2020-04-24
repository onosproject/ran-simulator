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
	"github.com/onosproject/ran-simulator/pkg/config"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/northbound/e2"
	"github.com/onosproject/ran-simulator/pkg/northbound/trafficsim"
	"github.com/onosproject/ran-simulator/pkg/southbound/kubernetes"
	"github.com/onosproject/ran-simulator/pkg/utils"
)

var log = logging.GetLogger("main")

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	googleAPIKey := flag.String("googleAPIKey", "", "your google maps api key")
	flag.Int("towerRows", 0, "replaced by yaml")        // TODO remove
	flag.Int("towerCols", 0, "replaced by yaml")        // TODO remove
	flag.Float64("mapCenterLat", 0, "replaced by yaml") // TODO remove
	flag.Float64("mapCenterLng", 0, "replaced by yaml") // TODO remove
	zoom := flag.Float64("zoom", 13, "The starting Zoom level")
	fade := flag.Bool("fade", true, "Show map as faded on start")
	showRoutes := flag.Bool("showRoutes", true, "Show routes on start")
	showPower := flag.Bool("showPower", true, "Show power as circle on start")
	flag.Float64("towerSpacingVert", 0, "replaced by yaml")  // TODO remove once removed from helm chart
	flag.Float64("towerSpacingHoriz", 0, "replaced by yaml") // TODO remove
	locationsScale := flag.Float64("locationsScale", 1.25, "Ratio of random locations diameter to tower grid width")
	maxUEs := flag.Uint("maxUEs", 300, "Max number of UEs for complete simulation")
	minUEs := flag.Uint("minUEs", 3, "Max number of UEs for complete simulation")
	stepDelayMs := flag.Uint("stepDelayMs", 1000, "delay between steps on route")
	flag.Int("maxUEsPerTower", 0, "replaced by yaml") // TODO remove
	metricsPort := flag.Uint("metricsPort", 9090, "port for Prometheus metrics")
	metricsAllHoEvents := flag.Bool("metricsAllHoEvents", true, "Export all HO events in metrics (only historgram if false)")
	topoEndpoint := flag.String("topoEndpoint", "onos-topo:5150", "Endpoint for the onos-topo service")
	flag.String("loglevel", "", "replaced by yaml") // TODO remove
	addK8sSvcPorts := flag.Bool("addK8sSvcPorts", true, "Add K8S service ports per tower")
	flag.Float64("avgCellcPerTower", 0, "replaced by yaml") // TODO remove
	towerConfigName := flag.String("towerConfigName", "berlin-honeycomb-4-3.yaml", "the name of a tower configuration")

	flag.Parse()

	towersConfig, err := config.GetTowerConfig(*towerConfigName)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Loaded config. %s. %d towers", *towerConfigName, len(towersConfig.TowersLayout))

	mapLayoutParams := types.MapLayout{
		Zoom:           float32(*zoom),
		ShowRoutes:     *showRoutes,
		Fade:           *fade,
		ShowPower:      *showPower,
		MinUes:         uint32(*minUEs),
		MaxUes:         uint32(*maxUEs),
		LocationsScale: float32(*locationsScale),
	}
	if mapLayoutParams.Zoom < 10 || mapLayoutParams.Zoom > 15 {
		log.Fatal("Invalid Zoom level - must be between 10 and 15 inclusive")
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

	mapLayoutParams.Center = &towersConfig.MapCentre

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

	allGrpcPorts := make([]uint16, 0)
	for _, tower := range towersConfig.TowersLayout {
		for _, sector := range tower.Sectors {
			allGrpcPorts = append(allGrpcPorts, sector.GrpcPort)
			log.Warnf("Handling Sector %s %s %d %d %d", tower.TowerID, sector.EcID, sector.GrpcPort, sector.Azimuth, sector.Arc)
			ecgi := types.ECGI{
				EcID:   sector.EcID,
				PlmnID: tower.PlmnID,
			}
			portNum := sector.GrpcPort
			go func() {
				// Blocks here when server running
				err := e2.NewTowerServer(ecgi, portNum, serverParams)
				if err != nil {
					log.Fatal("Unable to start server ", err)
				}
			}()
		}
	}
	// Add these new ports to the K8s service
	if *addK8sSvcPorts {
		err := kubernetes.AddK8SServicePorts(allGrpcPorts)
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
	mgr.Run(mapLayoutParams, towersConfig, routesParams, *topoEndpoint, serverParams, metricsParams)

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
