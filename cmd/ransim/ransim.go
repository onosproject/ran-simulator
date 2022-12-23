// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

/*
Package trafficsim is the main entry point to the ONOS TrafficSim application.

# Arguments

-caPath <the location of a CA certificate>

-keyPath <the location of a client private key>

-certPath <the location of a client certificate>

See ../../docs/run.md for how to run the application.
*/
package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

var log = logging.GetLogger("main")

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// The main entry point
func main() {
	log.Info("Starting Ran simulator")

	rand.Seed(time.Now().UnixNano())

	ready := make(chan bool)

	var serviceModelPlugins arrayFlags
	flag.Var(&serviceModelPlugins, "serviceModel", "names of service model plugins to load (repeated)")
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	grpcPort := flag.Int("grpcPort", 5150, "GRPC port for e2T server")
	modelName := flag.String("modelName", "model", "RANSim model file/resource name")
	metricName := flag.String("metricName", "", "RANSim metric file/resource name")
	hoLogic := flag.String("hoLogic", "local", "the location of handover logic {local, mho}")
	flag.Parse()

	if *hoLogic != "local" && *hoLogic != "mho" {
		log.Errorf("hoLogic arg should be one of {local, mho}")
		return
	}

	cfg := &manager.Config{
		CAPath:     *caPath,
		KeyPath:    *keyPath,
		CertPath:   *certPath,
		GRPCPort:   *grpcPort,
		ModelName:  *modelName,
		MetricName: *metricName,
		HOLogic:    *hoLogic,
	}

	mgr, err := manager.NewManager(cfg)
	if err == nil {
		mgr.Run()
		<-ready
		mgr.Close()
	}
}
