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

package e2

import (
	"time"

	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

// UpdateTelemetryMetrics ...
func UpdateTelemetryMetrics(m *e2.TelemetryMessage) {
	trafficSimMgr := manager.GetManager()
	switch x := m.S.(type) {
	case *e2.TelemetryMessage_RadioMeasReportPerUE:
		r := x.RadioMeasReportPerUE
		name := crntiToName(r.Crnti)
		ue := trafficSimMgr.UserEquipments[name]
		tower := trafficSimMgr.GetTowerByName(ue.ServingTower)
		reports := r.RadioReportServCells
		cqi := makeCqi(ue.ServingTowerDist, tower.GetTxPowerdB())
		if cqi < reports[0].CqiHist[0] || cqi < reports[1].CqiHist[0] || cqi < reports[2].CqiHist[0] {
			if ue.Metrics.HoReportTimestamp == 0 {
				ue.Metrics.HoReportTimestamp = time.Now().UnixNano()
			}
		}
	}
}

// UpdateControlMetrics ...
func UpdateControlMetrics(in *e2.ControlResponse) {
	trafficSimMgr := manager.GetManager()
	switch x := in.S.(type) {
	case *e2.ControlResponse_HORequest:
		m := x.HORequest
		ue := trafficSimMgr.UserEquipments[crntiToName(m.Crnti)]
		if ue.Metrics.HoReportTimestamp != 0 {
			ue.Metrics.HoLatency = time.Now().UnixNano() - ue.Metrics.HoReportTimestamp
			if trafficSimMgr.InfluxDbAvail {
				go manager.WriteHoMetricToInfluxDb(trafficSimMgr.InfluxDbAddr, m.Crnti, ue.Metrics.HoLatency, ue.Metrics.HoReportTimestamp)
			}
			ue.Metrics.HoReportTimestamp = 0
			log.Infof("Hand-over latency: %d microsec", ue.Metrics.HoLatency/1000)
		}
	}
}
