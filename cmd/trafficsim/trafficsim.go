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
	"github.com/OpenNetworkingFoundation/gmap-ran/pkg/manager"
	"github.com/OpenNetworkingFoundation/gmap-ran/pkg/northbound/trafficsim"
	"github.com/OpenNetworkingFoundation/gmap-ran/pkg/service"
	log "k8s.io/klog"
)

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	googleApiKey := flag.String("googleApiKey", "", "your google maps api key")
	towerRows := flag.Int("towerRows", 3, "Number of rows of towers")
	towerCols := flag.Int("towerCols", 3, "Number of columns of towers")
	mapCenterLat := flag.Float64("mapCenterLat", 52.5200, "Map center latitude")
	mapCenterLng := flag.Float64("mapCenterLng", 13.4050, "Map center longitude") // Berlin
	towerSpacingVert := flag.Float64("towerSpacingVert", 0.02, "Tower spacing vert in degrees latitude")
	towerSpacingHoriz := flag.Float64("towerSpacingHoriz", 0.02, "Tower spacing horiz in degrees longitude")
	numLocations := flag.Int("numLocations", 10, "Number of locations")
	numRoutes := flag.Int("numRoutes", 3, "Number of routes")

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

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	log.InitFlags(klogFlags)

	// Sync the glog and klog flags.
	flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			_ = f2.Value.Set(value)
		}
	})

	towerParams := manager.TowersParams{
		MapCenterLat:      float32(*mapCenterLat),
		MapCenterLng:      float32(*mapCenterLng),
		TowerRows:         *towerRows,
		TowerCols:         *towerCols,
		TowerSpacingVert:  float32(*towerSpacingVert),
		TowerSpacingHoriz: float32(*towerSpacingHoriz),
	}

	locationParams := manager.LocationsParams{NumLocations: *numLocations}

	routesParams := manager.RoutesParams{
		NumRoutes: *numRoutes,
		ApiKey:    *googleApiKey,
	}

	log.Info("Starting trafficsim")

	mgr, err := manager.NewManager()
	if err != nil {
		log.Fatal("Unable to load trafficsim ", err)
	} else {
		mgr.Run(towerParams, locationParams, routesParams)
		err = startServer(*caPath, *keyPath, *certPath)
		if err != nil {
			log.Fatal("Unable to start trafficsim ", err)
		}
	}
}

// Creates gRPC server and registers various services; then serves.
func startServer(caPath string, keyPath string, certPath string) error {
	s := service.NewServer(service.NewServerConfig(caPath, keyPath, certPath))
	s.AddService(trafficsim.Service{})

	return s.Serve(func(started string) {
		log.Info("Started NBI on ", started)
	})
}
