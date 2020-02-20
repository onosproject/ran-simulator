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
	"github.com/onosproject/ran-simulator/pkg/service"
	"time"

	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/northbound/e2"
	"github.com/onosproject/ran-simulator/pkg/northbound/trafficsim"
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
	zoom := flag.Float64("zoom", 12.5, "The starting Zoom level")
	fade := flag.Bool("fade", true, "Show map as faded on start")
	showRoutes := flag.Bool("showRoutes", true, "Show routes on start")
	showPower := flag.Bool("showPower", true, "Show power as circle on start")
	towerSpacingVert := flag.Float64("towerSpacingVert", 0.02, "Tower spacing vert in degrees latitude")
	towerSpacingHoriz := flag.Float64("towerSpacingHoriz", 0.02, "Tower spacing horiz in degrees longitude")
	numLocations := flag.Int("numLocations", 10, "Number of locations")
	numRoutes := flag.Int("numRoutes", 3, "Number of routes")
	stepDelayMs := flag.Int("stepDelayMs", 1000, "delay between steps on route")
	maxUEs := flag.Int("maxUEsPerTower", 5, "Max num of UEs per tower")

	//lines 93-109 are implemented according to
	// https://github.com/kubernetes/klog/blob/master/examples/coexist_glog/coexist_glog.go
	// because of libraries importing glog. With glog import we can't call log.InitFlags(nil) as per klog readme
	// thus the alsologtostderr is not set properly and we issue multiple logs.
	// Calling log.InitFlags(nil) throws panic with error `flag redefined: log_dir`
	err := flag.Set("alsologtostderr", "true")
	if err != nil {
		log.Error("Cant' avoid double Error logging ", err)
	}
	flag.Parse()

	mapLayoutParams := types.MapLayout{
		Center: &types.Point{
			Lat: float32(*mapCenterLat),
			Lng: float32(*mapCenterLng),
		},
		Zoom:       float32(*zoom),
		ShowRoutes: *showRoutes,
		Fade:       *fade,
		ShowPower:  *showPower,
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

	towerParams := types.TowersParams{
		TowerRows:         uint32(*towerRows),
		TowerCols:         uint32(*towerCols),
		TowerSpacingVert:  float32(*towerSpacingVert),
		TowerSpacingHoriz: float32(*towerSpacingHoriz),
		MaxUEs:            uint32(*maxUEs),
	}
	if towerParams.TowerRows < 2 || towerParams.TowerRows > 20 {
		log.Fatal("Invalid number of Tower Rows - must be between 2 and 20 inclusive")
	}
	if towerParams.TowerCols < 2 || towerParams.TowerCols > 20 {
		log.Fatal("Invalid number of Tower Cols - must be between 2 and 20 inclusive")
	}
	if towerParams.TowerSpacingVert < 0.001 || towerParams.TowerSpacingVert > 1.0 {
		log.Fatal("Invalid vertical tower spacing - must be between 0.001 and 1.0 degree latitude inclusive")
	}
	if towerParams.TowerSpacingHoriz < 0.001 || towerParams.TowerSpacingHoriz > 1.0 {
		log.Fatal("Invalid horizontal tower spacing - must be between 0.001 and 1.0 degree longitude inclusive")
	}

	locationParams := manager.LocationsParams{NumLocations: *numLocations}
	if locationParams.NumLocations < 3 || locationParams.NumLocations > 200 {
		log.Fatal("Invalid number of Locations - must be between 3 and 100 inclusive")
	}

	routesParams := manager.RoutesParams{
		NumRoutes: *numRoutes,
		APIKey:    *googleAPIKey,
		StepDelay: time.Duration(*stepDelayMs) * time.Millisecond,
	}
	if routesParams.NumRoutes < 2 || routesParams.NumRoutes > 100 {
		log.Fatal("Invalid number of Routes - must be between 2 and 100 inclusive")
	}
	if locationParams.NumLocations < routesParams.NumRoutes*2 {
		log.Fatal("Invalid number of Location:Routes - must be at least 2")
	}
	if *stepDelayMs < 100 || *stepDelayMs > 60000 {
		log.Fatal("Invalid step Delay - must be between 100ms and 60000ms inclusive")
	}

	log.Info("Starting trafficsim")

	mgr, err := manager.NewManager()
	if err != nil {
		log.Fatal("Unable to load trafficsim ", err)
		return
	}
	mgr.Run(mapLayoutParams, towerParams, locationParams, routesParams)

	if err = startServer(*caPath, *keyPath, *certPath); err != nil {
		log.Fatal("Unable to start trafficsim ", err)
	}
}

// Creates gRPC server and registers various services; then serves.
func startServer(caPath string, keyPath string, certPath string) error {
	s := service.NewServer(service.NewServerConfig(caPath, keyPath, certPath))
	s.AddService(trafficsim.Service{})
	s.AddService(e2.Service{})

	return s.Serve(func(started string) {
		log.Info("Started NBI on ", started)
	})
}
