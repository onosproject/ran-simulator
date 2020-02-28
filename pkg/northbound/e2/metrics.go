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
	"github.com/onosproject/ran-simulator/api/types"
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
		var ue *types.Ue
		var ok bool
		trafficSimMgr.UserEquipmentsLock.RLock()
		if ue, ok = trafficSimMgr.UserEquipments[name]; !ok {
			trafficSimMgr.UserEquipmentsLock.RUnlock()
			return
		}
		trafficSimMgr.UserEquipmentsLock.RUnlock()

		reports := r.RadioReportServCells

		bestCQI := reports[0].CqiHist[0]
		bestStationID := reports[0].Ecgi

		for i := 1; i < len(reports); i++ {
			temp := reports[i].CqiHist[0]
			if bestCQI < temp {
				bestCQI = temp
				bestStationID = reports[i].Ecgi
			}
		}
		trafficSimMgr.TowersLock.RLock()
		servingTower := trafficSimMgr.GetTowerByName(ue.ServingTower)
		trafficSimMgr.TowersLock.RUnlock()

		if servingTower.EcID != bestStationID.Ecid || servingTower.PlmnID != bestStationID.PlmnId {
			trafficSimMgr.UserEquipmentsLock.Lock()
			if ue.Metrics.HoReportTimestamp == 0 {
				ue.Metrics.HoReportTimestamp = time.Now().UnixNano()
			}
			trafficSimMgr.UserEquipmentsLock.Unlock()
		}
	}
}

// UpdateControlMetrics ...
func UpdateControlMetrics(in *e2.ControlResponse) {
	trafficSimMgr := manager.GetManager()
	switch x := in.S.(type) {
	case *e2.ControlResponse_HORequest:
		m := x.HORequest
		trafficSimMgr.UserEquipmentsLock.Lock()
		ue, ok := trafficSimMgr.UserEquipments[crntiToName(m.Crnti)]
		if !ok {
			log.Errorf("ue %s not found", crntiToName(m.Crnti))
			trafficSimMgr.UserEquipmentsLock.Unlock()
			return
		}
		if ue.Metrics.HoReportTimestamp != 0 {
			ue.Metrics.HoLatency = time.Now().UnixNano() - ue.Metrics.HoReportTimestamp
			ue.Metrics.HoReportTimestamp = 0
			trafficSimMgr.LatencyChannel <- ue.Metrics.HoLatency
			log.Infof("%s Hand-over latency: %d microsec", ue.Name, ue.Metrics.HoLatency/1000)
		}
		trafficSimMgr.UserEquipmentsLock.Unlock()
	}
}
