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
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/simmanager"
)

var log = logging.GetLogger("main")

// The main entry point
func main() {
	log.Info("Starting simulator")
	ready := make(chan bool)

	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	host := flag.String("host", "onos-e2t", "Host name of E2T server")
	grpcPort := flag.Int("grpcPort", 5150, "GRPC port for e2T server")
	sctpPort := flag.Int("sctpPort", 36421, "SCTP port for e2T server")
	flag.Parse()

	cfg := simmanager.Config{
		CAPath:     *caPath,
		KeyPath:    *keyPath,
		CertPath:   *certPath,
		GRPCPort:   *grpcPort,
		E2THost:    *host,
		E2SCTPPort: *sctpPort,
	}

	mgr := simmanager.NewManager(cfg)
	mgr.Run()

	<-ready

	mgr.Close()
}
