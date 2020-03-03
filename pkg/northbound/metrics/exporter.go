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

package metrics

import (
	"fmt"
	"net/http"

	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HOEvent is a structure for HO event
type HOEvent struct {
	Crnti        string
	ServingTower string
	HOLatency    int64
}

var log = liblog.GetLogger("northbound", "trafficsim")

// RunHOExposer runs Prometheus exposer
func RunHOExposer(port int, latencyChan chan HOEvent) {
	log.Infof("Starting Prometheus agent on http://:%d/metrics", port)
	hoLatencyHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "onosproject",
			Subsystem: "ransimulator",
			Name:      "hometrics",
			Help:      "time (µs) from when RadioMeasReportUE is sent to when Handover is complete",
			Buckets:   prometheus.ExponentialBuckets(1e3, 1.5, 20),
		},
	)
	prometheus.MustRegister(hoLatencyHistogram)
	go func() {
		// block here until a latency measurement is received
		for latency := range latencyChan {
			hoLatencyHistogram.Observe(float64(latency.HOLatency / 1e3))
			exposeAllHOEvents(latency)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("error serving prometheus metrics %s", err.Error())
	}
}

func exposeAllHOEvents(event HOEvent) {
	tmpHOEvent := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "onosproject",
		Subsystem: "ransimulator",
		Name:      "hoevents",
		ConstLabels: prometheus.Labels{
			"crnti":        event.Crnti,
			"servingtower": event.ServingTower,
		},
	})
	tmpHOEvent.Add(float64(event.HOLatency / 1e3))
	prometheus.MustRegister(tmpHOEvent)
}
