// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

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
	"github.com/onosproject/ran-simulator/pkg/northbound"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/northbound/trafficsim"
	"github.com/onosproject/ran-simulator/pkg/utils"
)

var log = logging.GetLogger("main")

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	googleAPIKey := flag.String("googleAPIKey", "", "your google maps api key")
	zoom := flag.Float64("zoom", 13, "The starting Zoom level")
	fade := flag.Bool("fade", true, "Show map as faded on start")
	showRoutes := flag.Bool("showRoutes", true, "Show routes on start")
	showPower := flag.Bool("showPower", true, "Show power as circle on start")
	locationsScale := flag.Float64("locationsScale", 1.25, "Ratio of random locations diameter to tower grid width")
	maxUEs := flag.Uint("maxUEs", 300, "Max number of UEs for complete simulation")
	minUEs := flag.Uint("minUEs", 3, "Max number of UEs for complete simulation")
	stepDelayMs := flag.Uint("stepDelayMs", 1000, "delay between steps on route")
	metricsPort := flag.Uint("metricsPort", 9090, "port for Prometheus metrics")
	metricsAllHoEvents := flag.Bool("metricsAllHoEvents", true, "Export all HO events in metrics (only historgram if false)")
	topoEndpoint := flag.String("topoEndpoint", "onos-topo:5150", "Endpoint for the onos-topo service")
	addK8sSvcPorts := flag.Bool("addK8sSvcPorts", true, "Add K8S service ports per tower")
	// TODO - remove the following - only after it's been removed from the Helm chart
	flag.String("towerConfigName", "", "unused - config is got from topo")

	flag.Parse()

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
	if mapLayoutParams.MinUes > 1000 {
		log.Fatal("Invalid number for MinUEs - must be between 0 and 1000 inclusive")
	}
	if mapLayoutParams.MinUes*2 > mapLayoutParams.MaxUes {
		log.Fatal("Invalid ratio of MaxUEs:MinUEs - must be at least 2")
	}

	routesParams := manager.RoutesParams{
		APIKey:    *googleAPIKey,
		StepDelay: time.Duration(*stepDelayMs) * time.Millisecond,
	}

	if *stepDelayMs < 100 || *stepDelayMs > 60000 {
		log.Fatal("Invalid step Delay - must be between 100ms and 60000ms inclusive")
	}

	serverParams := utils.ServerParams{
		CaPath:         *caPath,
		KeyPath:        *keyPath,
		CertPath:       *certPath,
		TopoEndpoint:   *topoEndpoint,
		AddK8sSvcPorts: *addK8sSvcPorts,
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
	mgr.Run(mapLayoutParams, routesParams, *topoEndpoint, serverParams, metricsParams, northbound.NewCellServer)

	sigs := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := startServer(*caPath, *keyPath, *certPath); err != nil {
			log.Warnf("Unable to start server %s. Shutting down", err.Error())
			shutdown <- true
		}
	}()
	go func() {
		sig := <-sigs
		log.Warnf("Received the %v signal. Shutting down", sig)
		shutdown <- true
	}()
	<-shutdown // Block here until bool arrives on channel
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
